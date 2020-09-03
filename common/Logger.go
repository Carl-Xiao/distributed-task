package common

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Logger 我们在项目用都使用这个日志对象
var Logger *zap.Logger

//InitLogger 初始化
func InitLogger(cfg *LogConfig) (err error) {
	// 做日志切割第三方包
	ws := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	// 日志输出的格式
	encoder := getEncoder()
	var level = new(zapcore.Level)
	err = level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}
	//添加控制台输出
	var writeSyncer zapcore.WriteSyncer
	if *level == zapcore.DebugLevel {
		zapcore.AddSync(os.Stdout)
		writeSyncer = zapcore.NewMultiWriteSyncer(ws, os.Stdout)
	} else {
		writeSyncer = ws
	}
	core := zapcore.NewCore(encoder, writeSyncer, level)
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return
}

//TimeEncoder 时间格式
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}
	layout := "2006-01-02 15:04:05"
	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}
	enc.AppendString(t.Format(layout))
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = TimeEncoder // 时间字符串
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // 函数调用
	return zapcore.NewJSONEncoder(encoderConfig)            // JSON格式
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

//Debug 调试
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...) // logger.go
}

//Info 调试
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

//Warn 调试
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

//Error 调试
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

//With 调试
func With(fields ...zap.Field) *zap.Logger {
	return Logger.With(fields...)
}
