package tx

import (
	"github.com/influxdata/influxdb/client/v2"
	. "github.com/thewayma/suricataM/comm/log"
	. "github.com/thewayma/suricataM/comm/st"
	"github.com/thewayma/suricataM/comm/utils"
	"github.com/thewayma/suricataM/transporter/g"
	"github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/list"
	"time"
)

// send
const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

//!< 半同步发送任务, 消费半异步队列
func startSendTasks() {
	cfg := g.Config()

	// init semaphore
	checkerConcurrent := cfg.Checker.MaxConcurrentConns
	tsdbConcurrent := cfg.Tsdb.MaxConcurrentConns

	if checkerConcurrent < 1 {
		checkerConcurrent = 1
	}

	if tsdbConcurrent < 1 {
		tsdbConcurrent = 1
	}

	// init send go-routines
	for node := range cfg.Checker.Cluster {
		queue := CheckerQueues[node]
		go forward2CheckerTask(queue, node, checkerConcurrent)
	}

	if cfg.Tsdb.Enabled {
		go forward2TsdbTask(tsdbConcurrent)
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

// Tsdb定时任务, 将数据通过http api发送到influxdb
func forward2TsdbTask(concurrent int) {
	batch := g.Config().Tsdb.Batch
	sema := semaphore.NewSemaphore(concurrent)

	for {
		items := TsdbQueue.PopBackBy(batch)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}
		//  同步Call + 有限并发 进行发送
		sema.Acquire()

		go func(itemList []interface{}) {
			defer sema.Release()

			// Make client
			c, err := client.NewHTTPClient(client.HTTPConfig{
				Addr:     g.Config().Tsdb.Address,
				Username: g.Config().Tsdb.UserName,
				Password: g.Config().Tsdb.Password,
			})
			if err != nil {
				Log.Error("Error creating InfluxDB Client: %s", err.Error())
			}
			defer c.Close()

			bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  g.Config().Tsdb.Database,
				Precision: "s",
			})

			for i := 0; i < len(itemList); i++ {
				tsdbitem := itemList[i].(*TsdbItem)

				ti, err := utils.String2Time(utils.UnixTsFormat(tsdbitem.Timestamp))
				if err != nil {
					ti = time.Now()
				}

				pt, err := client.NewPoint(
					tsdbitem.Name,
					tsdbitem.Tags,
					tsdbitem.Field,
					ti,
				)
				if err != nil {
					Log.Error("Formart TsdbItem: %s", err.Error())
					continue
				}

				bp.AddPoint(pt)
			}

			err = c.Write(bp)
			if err != nil {
				Log.Error("TsdbItem Write API err: %s", err.Error())
			}
		}(items)
	}
}
