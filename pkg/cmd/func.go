package cmd

import (
	"fmt"
	"log"

	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/server/engine"
	"github.com/phuhao00/sugar"
	"net/http"
)

func Run(servername string) {
	var serverKernel engine.IServerKernel
	instance := engine.MakeInstance(servername)
	if instance == nil {
		log.Fatalf(fmt.Sprintf("找不到[%v]的服务,请通过 engine.Register 进行注册", servername))
	}

	serverKernel = instance()
	//内部服务启动
	serverKernel.Run()
	defer func() {
		//内部服务关闭
		serverKernel.Destroy()
		release(servername)
	}()

	sugar.WaitSignal(serverKernel.OnSystemSignal)
}

func pprofs(addr string) {
	http.ListenAndServe(addr, nil)
}

func release(servername string) {
	logger.Info(fmt.Sprintf("【 %s 】Stopping server...", servername))

	//mysql.DBRepo.DbRClose()
	//mysql.DBRepo.DbWClose()
	//redis.RedisManagers.Close()
	//for _, etcdServer := range etcdServerArr {
	//	etcdServer.Stop()
	//}

	logger.Info(fmt.Sprintf("【 %s 】server is stoped", servername))

}
