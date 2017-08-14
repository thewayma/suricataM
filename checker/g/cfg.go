package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
)

type RpcConfig struct {
	Enabled bool
	Listen  string
}

type HbsConfig struct {
	Server   string
	Timeout  int64
	Interval int64
}

type RedisConfig struct {
	Address      string
	MaxIdle      int
	ConnTimeout  int
	ReadTimeout  int
	WriteTimeout int
}

type AlarmConfig struct {
	Enabled      bool
	MinInterval  int64
	QueuePattern string
	Redis        *RedisConfig
}

type GlobalConfig struct {
	MaxLinklistNum int
	Rpc            *RpcConfig
	Hbs            *HbsConfig
	Alarm          *AlarmConfig
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
