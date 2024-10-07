// Package grpcpool @Author evan_yxw
// @Date 2024/7/10 18:41:00
// @Desc
package grpcpool

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/etcdv3"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"

	//"github.com/evanyxw/monster-go/pkg/gormtool"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/timeutil"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	READ_MYSQL  = "read_mysql"
	WRITE_MYSQL = "write_mysql"
)

type MySQL struct {
	Read struct {
		Addr string `toml:"addr"`
		User string `toml:"user"`
		Pass string `toml:"pass"`
		Name string `toml:"name"`
	}
	Write struct {
		Addr string `toml:"addr"`
		User string `toml:"user"`
		Pass string `toml:"pass"`
		Name string `toml:"name"`
	}
	Base struct {
		MaxOpenConn     int           `toml:"maxOpenConn"`
		MaxIdleConn     int           `toml:"maxIdleConn"`
		ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
	}
}

//func InitMysql(mysqlConf *MySQL, logger *zap.Logger) {
//	// mysql
//
//	mysqlConn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
//		mysqlConf.Read.User, mysqlConf.Read.Pass, mysqlConf.Read.Addr, mysqlConf.Read.Name)
//
//	err := gormtool.NewMySQLFromConnNoPrepareStmt(READ_MYSQL, mysqlConn,
//		mysqlConf.Base.MaxOpenConn, mysqlConf.Base.MaxIdleConn)
//
//	if err != nil {
//		//variables.Log.Error(fmt.Sprintf("加载%s失败:", model.READ_MYSQL), err)
//		logger.Error("加载失败:", zap.Any("mysql", READ_MYSQL), zap.Any("error", err))
//		panic(err)
//	}
//
//	err = gormtool.NewMySQLFromConnNoPrepareStmt(WRITE_MYSQL, mysqlConn,
//		mysqlConf.Base.MaxOpenConn, mysqlConf.Base.MaxIdleConn)
//
//	if err != nil {
//		//variables.Log.Error(fmt.Sprintf("加载%s失败:", model.WRITE_MYSQL), err)
//		logger.Error("加载失败:", zap.Any("mysql", WRITE_MYSQL), zap.Any("error", err))
//		panic(err)
//	}
//}

func InitLog(logFile, service, env string) *zap.Logger {
	log, err := logger.NewJSONLogger(
		logger.WithDisableConsole(),
		//logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.All().Server.Name, env.Active().Value())),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", service, env)),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileP(logFile, service+".log"),
	)
	if err != nil {
		panic(err)
	}

	return log
}

func InitEtcd(addr []string, user, pass string) *etcdv3.Etcd {
	etcd, err := etcdv3.NewEtcd(addr, user, pass,
		"default", 3, nil)
	if err != nil {
		panic(err)
	}

	return etcd
}

func InitEtcdClient(addr []string, user, pass string) *clientv3.Client {
	// 创建 etcd 客户端
	etcdClient, err := clientv3.New(clientv3.Config{
		Username:    user,
		Password:    pass,
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}

	return etcdClient
}

func Shutdown(fn func()) {
	quit := make(chan os.Signal, 1)
	//signal.Notify(quit, os.Interrupt, os.Kill)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fn()
}
