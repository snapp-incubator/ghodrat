package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger
type zapLogger struct {
	config   *Config
	instance *zap.Logger
}

func NewZap(cfg *Config) Logger {
	zLog := &zapLogger{config: cfg}

	zLog.instance = zap.New(
		zapcore.NewCore(
			zLog.getEncoder(),
			zLog.getWriteSyncer(),
			zLog.getLoggerLevel(),
		),
		zLog.getOptions()...,
	)

	return zLog
}

///////////////////////////////////////////////////////
// OPTIONS
///////////////////////////////////////////////////////

// For mapping config logger to app logger levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *zapLogger) getEncoder() zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig
	if l.config.Development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if l.config.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	return encoder
}

func (l *zapLogger) getWriteSyncer() zapcore.WriteSyncer {
	return zapcore.Lock(os.Stdout)
}

func (l *zapLogger) getLoggerLevel() zap.AtomicLevel {
	level, exist := loggerLevelMap[l.config.Level]
	if !exist {
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	return zap.NewAtomicLevelAt(level)
}

func (l *zapLogger) getOptions() []zap.Option {
	options := make([]zap.Option, 0, 2)

	if !l.config.Development {
		options = append(options, zap.AddCaller())
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	}

	return options
}

///////////////////////////////////////////////////////
// METHODS
///////////////////////////////////////////////////////

func (l *zapLogger) Named(name string) Logger {
	return &zapLogger{instance: l.instance.Named(name)}
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.instance.Debug(msg, convertFields(fields...)...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.instance.Info(msg, convertFields(fields...)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.instance.Warn(msg, convertFields(fields...)...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.instance.Error(msg, convertFields(fields...)...)
}

func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.instance.Panic(msg, convertFields(fields...)...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.instance.Fatal(msg, convertFields(fields...)...)
}

// convertFields converts Field To ZapField
func convertFields(fields ...Field) []zapcore.Field {
	zapFileds := make([]zapcore.Field, 0, len(fields))

	for index := 0; index < len(fields); index++ {
		zapField := convertField(fields[index])
		zapFileds = append(zapFileds, zapField)
	}

	return zapFileds
}

func convertField(field Field) zapcore.Field {
	switch field.Type {
	case AnyType:
		return zap.Any(field.Key, field.Value)
	case BoolType:
		return zap.Bool(field.Key, field.Value.(bool))
	case IntType:
		return zap.Int(field.Key, field.Value.(int))
	case Float64Type:
		return zap.Float64(field.Key, field.Value.(float64))
	case StringType:
		return zap.String(field.Key, field.Value.(string))
	case ErrorType:
		return zap.Error(field.Value.(error))
	}

	return zapcore.Field{}
}
