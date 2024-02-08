package configs

import (
	"bytes"
	_ "embed"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/evanyxw/monster-go/pkg/env"
	"github.com/evanyxw/monster-go/pkg/file"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config = new(Config)
var NotifyChan = make(chan struct{})

// Config 配置
type Config struct {
	MySQL struct {
		Read struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"read"`
		Write struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"write"`
		Base struct {
			MaxOpenConn     int           `toml:"maxOpenConn"`
			MaxIdleConn     int           `toml:"maxIdleConn"`
			ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
		} `toml:"base"`
	} `toml:"mysql"`

	Redis struct {
		Addr         string `toml:"addr"`
		Pass         string `toml:"pass"`
		Db           int    `toml:"db"`
		MaxRetries   int    `toml:"maxRetries"`
		PoolSize     int    `toml:"poolSize"`
		MinIdleConns int    `toml:"minIdleConns"`
		IsCluster    bool   `toml:"isCluster"`
	} `toml:"redis"`

	Mail struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
		User string `toml:"user"`
		Pass string `toml:"pass"`
		To   string `toml:"to"`
	} `toml:"mail"`

	HashIds struct {
		Secret string `toml:"secret"`
		Length int    `toml:"length"`
	} `toml:"hashids"`

	Language struct {
		Local string `toml:"local"`
	} `toml:"language"`

	Server struct {
		Address      string `toml:"address"`
		MaxConnNum   int    `toml:"maxConnNum"`
		BuffSize     int    `toml:"buffSize"`
		Pprof        bool   `toml:"prof"`
		PprofAddress string `toml:"profAddress"`
	} `toml:"server"`

	Rpc struct {
		Address string `toml:"address"`
	} `toml:"rpc"`

	Etcd struct {
		Address   []string `toml:"address"`
		Namespace string   `toml:"namespace"`
		Secret    string   `toml:"secret"`
	}
}

var (
	//go:embed dev_configs.toml
	devConfigs []byte

	//go:embed fat_configs.toml
	fatConfigs []byte

	//go:embed uat_configs.toml
	uatConfigs []byte

	//go:embed pro_configs.toml
	proConfigs []byte
)

// Get 获取config
func Get() Config {
	return *config
}

func Init() {
	var r io.Reader
	switch env.Active().Value() {
	case "dev":
		r = bytes.NewReader(devConfigs)
	case "fat":
		r = bytes.NewReader(fatConfigs)
	case "uat":
		r = bytes.NewReader(uatConfigs)
	case "pro":
		r = bytes.NewReader(proConfigs)
	default:
		r = bytes.NewReader(fatConfigs)
	}

	viper.SetConfigType("toml")

	if err := viper.ReadConfig(r); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	viper.SetConfigName(env.Active().Value() + "_configs")
	viper.AddConfigPath("./configs")

	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	//viper.AutomaticEnv()
	//d := viper.GetStringSlice("etcd.address")
	//fmt.Println(d)

	configFile := "./configs/" + env.Active().Value() + "_configs.toml"
	_, ok := file.IsExists(configFile)
	if !ok {
		if err := os.MkdirAll(filepath.Dir(configFile), 0766); err != nil {
			panic(err)
		}

		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := viper.WriteConfig(); err != nil {
			panic(err)
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		NotifyChan <- struct{}{}
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
}
