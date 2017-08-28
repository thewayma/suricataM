package st

import (
	"fmt"
	"github.com/thewayma/suricataM/comm/utils"
)

// 机器监控和实例监控都会产生Event，共用这么一个struct
type Event struct {
	Id          string            `json:"id"`
	Strategy    *Strategy         `json:"strategy"`
	Status      string            `json:"status"` // OK or PROBLEM
	Endpoint    string            `json:"endpoint"`
	LeftValue   float64           `json:"leftValue"`
	CurrentStep int               `json:"currentStep"`
	EventTime   int64             `json:"eventTime"`
	PushedTags  map[string]string `json:"pushedTags"`
}

func (this *Event) FormattedTime() string {
	return utils.UnixTsFormat(this.EventTime)
}

func (this *Event) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Status:%s, Strategy:%v, LeftValue:%s, CurrentStep:%d, PushedTags:%v, TS:%s>",
		this.Endpoint,
		this.Status,
		this.Strategy,
		utils.ReadableFloat(this.LeftValue),
		this.CurrentStep,
		this.PushedTags,
		this.FormattedTime(),
	)
}

func (this *Event) StrategyId() int {
	return this.Strategy.Id
}

func (this *Event) Priority() int {
	return this.Strategy.Priority
}

func (this *Event) Metric() string {
	return this.Strategy.Metric
}

func (this *Event) RightValue() float64 {
	return this.Strategy.RightValue
}

func (this *Event) Operator() string {
	return this.Strategy.Operator
}

func (this *Event) Func() string {
	return this.Strategy.Func
}

func (this *Event) MaxStep() int {
	return this.Strategy.MaxStep
}

func (this *Event) Counter() string {
	return fmt.Sprintf("%s/%s %s", this.Endpoint, this.Metric(), utils.SortedTags(this.PushedTags))
}
