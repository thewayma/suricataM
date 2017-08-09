package g

import (
    "runtime"
)

const (
    GAUGE        = "GAUGE"
    COUNTER      = "COUNTER"
    DERIVE       = "DERIVE"
    DEFAULT_STEP = 60
    MIN_STEP     = 30
)

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU())
}
