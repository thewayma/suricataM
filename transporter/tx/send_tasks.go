package tx

import (
	. "github.com/thewayma/suricataM/comm/log"
	. "github.com/thewayma/suricataM/comm/st"
	"github.com/thewayma/suricataM/transporter/g"
	"github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/list"
	"time"
)

const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

//!< 半同步发送任务, 消费半异步队列
func startSendTasks() {
	cfg := g.Config()

	// init semaphore
	checkerConcurrent := cfg.Checker.MaxConcurrentConns
	influxdbConcurrent := cfg.InfluxDB.MaxConcurrentConns

	if checkerConcurrent < 1 {
		checkerConcurrent = 1
	}

	if influxdbConcurrent < 1 {
		influxdbConcurrent = 1
	}

	// init send go-routines
	for node := range cfg.Checker.Cluster {
		queue := CheckerQueues[node]
		go forward2CheckerTask(queue, node, checkerConcurrent)
	}

	if cfg.InfluxDB.Enabled {
		go forward2InfluxDBTask(influxdbConcurrent)
	}
}

func forward2CheckerTask(Q *list.SafeListLimited, node string, concurrent int) {
	batch := g.Config().Checker.Batch // 一次发送,最多batch条数据
	addr := g.Config().Checker.Cluster[node]
	sema := semaphore.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		checkerItems := make([]*CheckerItem, count)
		for i := 0; i < count; i++ {
			checkerItems[i] = items[i].(*CheckerItem)
		}

		//	同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(addr string, checkerItems []*CheckerItem, count int) {
			defer sema.Release()

			resp := &SimpleRpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = CheckerConnPools.Call(addr, "Checker.Send", checkerItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				Log.Error("forward2CheckerTask send checker %s:%s fail: %v", node, addr, err)
				//proc.SendToCheckerFailCnt.IncrBy(int64(count))
			} else {
				//proc.SendToCheckerCnt.IncrBy(int64(count))
			}
		}(addr, checkerItems, count)
	}
}

func forward2InfluxDBTask(concurrent int) {
	batch := g.Config().InfluxDB.Batch
	sema := semaphore.NewSemaphore(concurrent)

	for {
		items := InfluxDBQueue.PopBackBy(batch)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}
		//  同步Call + 有限并发 进行发送
		sema.Acquire()

		go func(itemList []interface{}) {
			defer sema.Release()

			var r InfluxDBInstance
			err := r.InitDB(g.Config().InfluxDB.Address, g.Config().InfluxDB.UserName, g.Config().InfluxDB.Password, g.Config().InfluxDB.Database, "s")
			if err != nil {
				Log.Error("Error Creating InfluxDB Client: %s", err.Error())
				return
			}
			r.DeferClose()

			err = r.AddPoints(itemList)
			if err != nil {
				Log.Error("Error Attach Point to BatchPoints: %s", err.Error())
				return
			}

			err = r.Write()
			if err != nil {
				Log.Error("Error Write BatchPoints API: %s", err.Error())
				return
			}
		}(items)
	}
}
