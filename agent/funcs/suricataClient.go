package funcs

import (
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/thewayma/suricataM/agent/g"
	"github.com/thewayma/suricataM/comm/log"
	"net"
	"os"
)

type ifStat struct {
	Iface          string
	Pkts           int64
	Drop           int64
	InvaldChecksum int64
}

var (
	protocolMap map[string]string
	buf         = make([]byte, 1024)

	ifaceMap map[int]string //!< portId <-> portName
)

func init() {
	protocolMap = make(map[string]string)
	protocolMap["version"] = `{"version": "0.1"}`
	protocolMap["command"] = `{"command": "%s"}`
	protocolMap["commandArgument"] = `{"command": "%s", "arguments": {"%s": "%s"}}`

	ifaceMap = make(map[int]string)
}

func suriConnect() (net.Conn, err) {
	conn, err := net.Dial("unix", g.Config().Suricata.UnixSockFile)
	if err != nil {
		log.Log.CRITICAL("Unix File %s not found", g.Config().Suricata.UnixSockFile)
		return nil, err
	}

	return conn, nil
}

func suriMakeCommand(com string) string {
	return fmt.Sprintf(protocolMap["command"], com)
}

func suriMakeCommandArgument(com, argKey, argValue string) string {
	return fmt.Sprintf(protocolMap["commandArgument"], com, argKey, argValue)
}

func suriSendVersion(conn net.Conn) {
	//fmt.Printf("SND: %s\n", protocolMap["version"])
	conn.Write([]byte(protocolMap["version"]))

	conn.Read(buf)
	//fmt.Printf("RCV: %s\n", buf)

	//!< TODO: OK, NOK
}

func suriSendCommandGet(conn net.Conn, data string) (interface{}, error) {
	conn.Write([]byte(data))
	//fmt.Printf("SND: %s\n", data)

	conn.Read(buf)
	//fmt.Printf("RCV: %s\n", buf)

	j, _ := jason.NewObjectFromBytes([]byte(buf))

	if res, _ := j.GetString("return"); res == "OK" {
		return j, nil
	} else {
		return -299, fmt.Errorf("%s Command Error", data)
	}
}

func GetUptime() int64 {
	conn, err := suriConnect()
	if err != nil {
		return -1
	}

	defer conn.Close()
	suriSendVersion(conn)

	com := suriMakeCommand("uptime")
	//ret, _ := suriSendCommandGetInt(conn, com)
	ret, _ := suriSendCommandGet(conn, com)
	obj := ret.(*jason.Object)
	uptime, _ := obj.GetInt64("message")

	//fmt.Println("Uptime:", g.GaugeValue("suricata_uptime", uptime)
	//return []*g.MetricValue{g.GaugeValue("suricata_uptime", uptime)}
	return uptime
}

func ShutDown() string {
	conn, err := suriConnect()
	if err != nil {
		return "UnixSockFile Not Found"
	}

	defer conn.Close()
	suriSendVersion(conn)

	com := suriMakeCommand("shutdown")
	ret, _ := suriSendCommandGet(conn, com)
	obj := ret.(*jason.Object)
	str, _ := obj.GetString("message")

	//fmt.Println(str)
	return str
}

func ReloadRules() string {
	conn, err := suriConnect()
	if err != nil {
		return "UnixSockFile Not Found"
	}

	defer conn.Close()
	suriSendVersion(conn)

	com := suriMakeCommand("reload-rules")
	ret, _ := suriSendCommandGet(conn, com)
	obj := ret.(*jason.Object)
	str, _ := obj.GetString("message")

	//fmt.Println(str)
	return str
}

func GetVersion() string {
	conn, err := suriConnect()
	if err != nil {
		return "UnixSockFile Not Found"
	}

	defer conn.Close()
	suriSendVersion(conn)

	com := suriMakeCommand("version")
	ret, _ := suriSendCommandGet(conn, com)
	obj := ret.(*jason.Object)
	str, _ := obj.GetString("message")

	//fmt.Println(str)
	return str
}

func GetRunningMode() string {
	conn, err := suriConnect()
	if err != nil {
		return "UnixSockFile Not Found"
	}

	defer conn.Close()
	suriSendVersion(conn)

	com := suriMakeCommand("running-mode")
	ret, _ := suriSendCommandGet(conn, com)
	obj := ret.(*jason.Object)
	str, _ := obj.GetString("message")

	//fmt.Println(str)
	return str
}

func GetCaptureMode() string {
	conn, err := suriConnect()
	if err != nil {
		return "UnixSockFile Not Found"
	}

	defer conn.Close()
	suriSendVersion(conn)

	com := suriMakeCommand("capture-mode")
	ret, _ := suriSendCommandGet(conn, com)
	obj := ret.(*jason.Object)
	str, _ := obj.GetString("message")

	//fmt.Println(str)
	return str
}

func GetProfilingCouters() []byte {
	conn, err := suriConnect()
	if err != nil {
		return []byte("UnixSockFile Not Found")
	}

	defer conn.Close()
	suriSendVersion(conn)

	com := suriMakeCommand("dump-counters")
	conn.Write([]byte(com))

	buf = make([]byte, 10240)
	conn.Read(buf)

	//fmt.Printf("ProfilingCounters: %s\n", buf)
	return buf
}

func GetAllPortStats() map[string]*ifStat {
	conn, err := suriConnect()
	if err != nil {
		return nil
	}

	defer conn.Close()
	iStat := make(map[string]*ifStat)

	suriSendVersion(conn)

	com := suriMakeCommand("iface-list")
	res, _ := suriSendCommandGet(conn, com)

	obj, _ := res.(*jason.Object)
	messObj, _ := obj.GetObject("message")
	ifaceObj, _ := messObj.GetStringArray("ifaces")

	for index, dataItem := range ifaceObj {
		ifaceMap[index] = dataItem
		com := suriMakeCommandArgument("iface-stat", "iface", dataItem)
		res, _ := suriSendCommandGet(conn, com)

		obj, _ := res.(*jason.Object)
		messObj, _ := obj.GetObject("message")

		pkts, _ := messObj.GetInt64("pkts")
		drop, _ := messObj.GetInt64("drop")
		invalid, _ := messObj.GetInt64("invalid-checksums")

		iStat[dataItem] = &ifStat{dataItem, pkts, drop, invalid}

		//fmt.Printf("Iface:%s, pkt=%d, drop=%d, invalid-checksums=%d\n", dataItem, pkts, drop, invalid)
	}

	return iStat
}
