package logger

const (
	LogLevelDebug = iota
	LogLevelTrace
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

const (
	LogSplitType = "size"    // 默认使用大小进行切分
	LogSplitSize = 104857600 // 100M
	LogSplitTime = "hour"
)
