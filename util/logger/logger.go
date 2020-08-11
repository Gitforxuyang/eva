package logger

import (
	"context"
	"github.com/Gitforxuyang/eva/util/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type EvaLogger interface {
	Debug(ctx context.Context, msg string, fields Fields) error
	Info(ctx context.Context, msg string, fields Fields) error
	Error(ctx context.Context, msg string, fields Fields) error
	Warn(ctx context.Context, msg string, fields Fields) error
	Panic(ctx context.Context, msg string, fields Fields) error
}

type evaLogger struct {
	appId string
	log   *zap.Logger
}

func (m *evaLogger) getFields(ctx context.Context, fields Fields) []zap.Field {
	fids := make([]zap.Field, 0, len(fields)+3)
	fids = append(fids, zap.Any("timestamp", utils.FormatTime(time.Now(), "2006-01-01 15:04:05")))
	fids = append(fids, zap.Any("traceId", utils.GetTraceId(ctx)))
	fids = append(fids, zap.Any("appId", m.appId))
	for k, v := range fields {
		fids = append(fids, zap.Any(k, v))
	}
	return fids
}

func (m *evaLogger) Info(ctx context.Context, msg string, fields Fields) error {
	fids := m.getFields(ctx, fields)
	m.log.Info(msg, fids...)
	return nil
}
func (m *evaLogger) Debug(ctx context.Context, msg string, fields Fields) error {
	fids := m.getFields(ctx, fields)
	m.log.Debug(msg, fids...)
	return nil
}
func (m *evaLogger) Error(ctx context.Context, msg string, fields Fields) error {
	fids := m.getFields(ctx, fields)
	m.log.Error(msg, fids...)
	return nil
}
func (m *evaLogger) Warn(ctx context.Context, msg string, fields Fields) error {
	fids := m.getFields(ctx, fields)
	m.log.Warn(msg, fids...)
	return nil
}
func (m *evaLogger) Panic(ctx context.Context, msg string, fields Fields) error {
	fids := m.getFields(ctx, fields)
	m.log.Panic(msg, fids...)
	return nil
}

type Fields map[string]interface{}

var (
	m *evaLogger
)

func Init(appId string) error {
	m = &evaLogger{appId: appId}
	var err error
	m.log, err = zap.NewProduction(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.WarnLevel))
	//d := Demo{Name: "123123", A: &Animal{Animal: "dog"}}
	//m.Info(context.TODO(), "msg", Fields{"d": d, "key": "value"})
	return err
}
func GetLogger() EvaLogger {
	if m == nil {
		panic("logger没有初始化")
	}
	return m
}
