package g

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/toolkits/file"
)

var (
	LocalIp    string
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

type SuricataConfig struct {
	Bin                        string
	Conf                       string
	Ifaces                     []string
	RulesDir                   string
	UnixSockFile               string
	ControlCommandPollInterval int
	MonitorMetricPollInterval  int
}

type HeartbeatConfig struct {
	Enabled  bool
	RpcAddr  string
	HttpAddr string
	Interval int
	Timeout  int
}

type TransporterConfig struct {
	Enabled  bool
	Addrs    []string
	Interval int //!< 监控项采集周期
	Timeout  int
}

type HttpConfig struct {
	Enabled bool
	Listen  string
}

type AgentConfig struct {
	Hostname string
	Ip       string
	Http     *HttpConfig
}

type GlobalConfig struct {
	Agent       *AgentConfig
	Suricata    *SuricataConfig
	Heartbeat   *HeartbeatConfig
	Transporter *TransporterConfig
	DefaultTags map[string]string
}

func InitLocalIp() {
	if Config().Transporter.Enabled {
		conn, err := net.DialTimeout("tcp", Config().Transporter.Addrs[0], time.Second*10)
		if err != nil {
			log.Println("get local addr failed, err:", err.Error())
		} else {
			LocalIp = strings.Split(conn.LocalAddr().String(), ":")[0]
			conn.Close()
		}
	} else {
		log.Println("hearbeat is not enabled, can't get localip")
	}
}

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func Hostname() string {
	hostname := Config().Agent.Hostname
	if hostname != "" {
		return hostname
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
		return "all-default-hostname"
	}
	return hostname
}

func IP() string {
	ip := Config().Agent.Ip
	if ip != "" {
		// use ip in configuration
		return ip
	}

	if len(LocalIp) > 0 {
		ip = LocalIp
	}

	return ip
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Println("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Println("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Println("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Println("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	config = &c
	lock.Unlock()

	InitLocalIp()

	log.Println("read config file:", cfg, "successfully", "LocalIp=", IP())
}
