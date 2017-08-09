package st

import (
    "fmt"
	"sync"
	"time"
	"math"
	"net/rpc"
	"github.com/toolkits/net"
    ."github.com/thewayma/suricataM/comm/log"
)

type SimpleRpcResponse struct {
    Code int `json:"code"`  //!< 0:success, 1:bad request
}

func (this *SimpleRpcResponse) String() string {
    return fmt.Sprintf("<Code: %d>", this.Code)
}

type TransporterResponse struct {
    Message string
    Total   int
    Invalid int
    Latency int64
}

func (this *TransporterResponse) String() string {
    return fmt.Sprintf(
        "<Total=%v, Invalid:%v, Latency=%vms, Message:%s>",
        this.Total,
        this.Invalid,
        this.Latency,
        this.Message,
    )
}

type RpcClient struct {
	sync.Mutex
	Peer      string //!< debug info: agent => heartbeat, or agent <= transporter
	RpcServer string
	rpcClient *rpc.Client
	Timeout   time.Duration
}

func (this *RpcClient) close() {
	if this.rpcClient != nil {
		this.rpcClient.Close()
		this.rpcClient = nil
	}
}

func (this *RpcClient) serverConn() error {
	if this.rpcClient != nil {
		return nil
	}

	var err error
	var retry int = 1

	for {
		if this.rpcClient != nil {
			return nil
		}

		this.rpcClient, err = net.JsonRpcClient("tcp", this.RpcServer, this.Timeout)
		if err != nil {
			Log.Error("%s, dial %s fail: %v", this.Peer, this.RpcServer, err)
			if retry > 3 {
				return err
			}
			time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)  //!< 指数回退
			retry++
			continue
		}
		return err
	}
}

func (this *RpcClient) Call(method string, args interface{}, reply interface{}) error {
	this.Lock()
	defer this.Unlock()

	err := this.serverConn()
	if err != nil {
		return err
	}

	timeout := time.Duration(10 * time.Second)
	done := make(chan error, 1)

	go func() {
		err := this.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		Log.Error("%s, [WARN] rpc call timeout %v => %v", this.Peer, this.rpcClient, this.RpcServer)
		this.close()
	case err := <-done:
		if err != nil {
			this.close()
			return err
		}
	}

	return nil
}
