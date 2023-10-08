package logger

var Nil = NilLogger{}

type NilLogger struct{}

func (n NilLogger) Debug(msg string)                    {}
func (n NilLogger) DebugF(msg string, v ...interface{}) {}
func (n NilLogger) Info(msg string)                     {}
func (n NilLogger) InfoF(msg string, v ...interface{})  {}
func (n NilLogger) Warn(msg string)                     {}
func (n NilLogger) WarnF(msg string, v ...interface{})  {}
func (n NilLogger) Error(err error)                     {}
func (n NilLogger) ErrorC(err error, errCode int)       {}
