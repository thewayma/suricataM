package st

import (
    "strings"
    "strconv"
)

// 统一agent,transporter 传输数据格式, 减小内存拷贝
type MetricData struct {
    Endpoint    string              `json:"endpoint"`
    Metric      string              `json:"metric"`
    Value       float64             `json:"value"`
    Step        int64               `json:"step"`
    Type        string              `json:"Type"`
    Tags        map[string]string   `json:"tags"`
    Timestamp   int64               `json:"timestamp"`
}

func (t *MetricData) String() string {
    return fmt.Sprintf("<MetricData Endpoint:%s, Metric:%s, Timestamp:%d, Step:%d, Value:%f, Tags:%v>",
        t.Endpoint, t.Metric, t.Timestamp, t.Step, t.Value, t.Tags)
}

func NewMetric(metric string, v interface{}, dataType string, tags ...string) *MetricData {
    //!< 在agent端判断数据类型, 避免transporter的内存拷贝
    var vv float64
    switch cv := v.(type) {
    case string:
        vv, _ = strconv.ParseFloat(cv, 64)
    case float64:
        vv = cv
    case int64:
        vv = float64(cv)
    }

    mv := MetricData {
        Metric: metric,
        Value:  vv,
        Type:   dataType,
        Tags:   make(map[string]string),
    }

    for _, tag := range tags {
        str := strings.Split(tag, "=")
        mv.Tags[str[0]] = str[1]
    }

    return &mv
}

func GaugeValue(metric string, val interface{}, tags ...string) *MetricData {
    return NewMetric(metric, val, "GAUGE", tags...)		//!< 瞬时型监控值
}

func CounterValue(metric string, val interface{}, tags ...string) *MetricData {
    return NewMetric(metric, val, "COUNTER", tags...)	//!< 累加型监控值
}
