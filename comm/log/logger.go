package log

import (
    "github.com/alecthomas/log4go"
)

var Log = make(log4go.Logger)

//!< 日志等级从低到高: FINEST, FINE, DEBUG, TRACE, INFO, WARNING, ERROR, CRITICAL
func init() {
    file := log4go.NewFileLogWriter("run.log", true)
    file.SetRotateLines(3600)
    file.SetRotateDaily(true)
    Log.AddFilter("file", log4go.DEBUG, file)

    Log.Info("Log Framework Start")
}
