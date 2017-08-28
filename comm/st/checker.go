package st

import (
	"fmt"
	"github.com/thewayma/suricataM/comm/utils"
)

type CheckerItem struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	Type      string            `json:"judgeType"`
	Tags      map[string]string `json:"tags"`
}

func (this *CheckerItem) String() string {
	return fmt.Sprintf("<Endpoint:%s, Metric:%s, Value:%f, Timestamp:%d, JudgeType:%s Tags:%v>",
		this.Endpoint,
		this.Metric,
		this.Value,
		this.Timestamp,
		this.Type,
		this.Tags)
}

func (this *CheckerItem) PrimaryKey() string {
	return utils.Md5(utils.PK(this.Endpoint, this.Metric, this.Tags))
}

type HistoryData struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type StrategiesRequest struct {
	Ip       string
	Hostname string
}

type Strategy struct {
	Id          int     `json:"id"`
	Cycle       int     `json:"cycle"`
	Metric      string  `json:"item"`
	Func        string  `json:"calc"`   // e.g. max(#3) all(#3)
	Operator    string  `json:"opt"`    // e.g. < !=
	RightValue  float64 `json:"value1"` // critical value
	RightValue2 float64 `json:"value2"` // critical value
	Action      string  `json:"action"`
	Priority    int     `json:"priority"`
	MaxStep     int     `json:"maxstep"`
}

func (this *Strategy) String() string {
	return fmt.Sprintf(
		"<Id=%d, Metric=%s, %s(#%s) %s %s, Priority=%d, MaxStep=%d",
		this.Id,
		this.Metric,
		this.Func,
		this.Cycle,
		this.Operator,
		utils.ReadableFloat(this.RightValue),
		this.Priority,
		this.MaxStep,
	)
}

type StrategiesResponse struct {
	Version    string     `json:"version"`
	Strategies []Strategy `json:"Policies"`
}
