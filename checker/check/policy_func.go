package check

import (
	"fmt"
	"github.com/thewayma/suricataM/comm/log"
	"github.com/thewayma/suricataM/comm/st"
	"math"
)

type Function interface {
	Compute(L *SafeLinkedList) (vs []*st.HistoryData, leftValue float64, isTriggered bool, isEnough bool)
}

type MaxFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (this MaxFunction) Compute(L *SafeLinkedList) (vs []*st.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	max := vs[0].Value
	for i := 1; i < this.Limit; i++ {
		if max < vs[i].Value {
			max = vs[i].Value
		}
	}

	leftValue = max
	isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
	return
}

type MinFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (this MinFunction) Compute(L *SafeLinkedList) (vs []*st.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	min := vs[0].Value
	for i := 1; i < this.Limit; i++ {
		if min > vs[i].Value {
			min = vs[i].Value
		}
	}

	leftValue = min
	isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
	return
}

type AllFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (this AllFunction) Compute(L *SafeLinkedList) (vs []*st.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	isTriggered = true
	for i := 0; i < this.Limit; i++ {
		isTriggered = checkIsTriggered(vs[i].Value, this.Operator, this.RightValue)
		if !isTriggered {
			break
		}
	}

	leftValue = vs[0].Value
	return
}

type SumFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (this SumFunction) Compute(L *SafeLinkedList) (vs []*st.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	sum := 0.0
	for i := 0; i < this.Limit; i++ {
		sum += vs[i].Value
	}

	leftValue = sum
	isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
	return
}

type AvgFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (this AvgFunction) Compute(L *SafeLinkedList) (vs []*st.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(this.Limit)
	if !isEnough {
		return
	}

	sum := 0.0
	for i := 0; i < this.Limit; i++ {
		sum += vs[i].Value
	}

	leftValue = sum / float64(this.Limit)
	isTriggered = checkIsTriggered(leftValue, this.Operator, this.RightValue)
	return
}

func ParseFuncFromString(cycle int, funcs string, operator string, rightValue float64) (fn Function, err error) {
	if cycle < 0 || cycle > 20 {
		log.Log.Error("PolicyChecker, Invalid Cycle=%d", cycle)
		cycle = 3
	}

	switch funcs {
	case "max":
		fn = &MaxFunction{Limit: cycle, Operator: operator, RightValue: rightValue}
	case "min":
		fn = &MinFunction{Limit: cycle, Operator: operator, RightValue: rightValue}
	case "all":
		fn = &AllFunction{Limit: cycle, Operator: operator, RightValue: rightValue}
	case "sum":
		fn = &SumFunction{Limit: cycle, Operator: operator, RightValue: rightValue}
	case "avg":
		fn = &AvgFunction{Limit: cycle, Operator: operator, RightValue: rightValue}
	default:
		err = fmt.Errorf("not_supported_method")
	}

	return
}

func checkIsTriggered(leftValue float64, operator string, rightValue float64) (isTriggered bool) {
	switch operator {
	case "=", "==":
		isTriggered = math.Abs(leftValue-rightValue) < 0.0001
	case "!=":
		isTriggered = math.Abs(leftValue-rightValue) > 0.0001
	case "<":
		isTriggered = leftValue < rightValue
	case "<=":
		isTriggered = leftValue <= rightValue
	case ">":
		isTriggered = leftValue > rightValue
	case ">=":
		isTriggered = leftValue >= rightValue
	}

	return
}
