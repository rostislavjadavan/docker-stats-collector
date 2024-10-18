package lib

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"io"
)

// EchoLogger adapts zerolog to Echo's logger interface
type EchoLogger struct {
	ZeroLog zerolog.Logger
}

func (l *EchoLogger) Output() io.Writer {
	return l.ZeroLog
}

func (l *EchoLogger) SetOutput(w io.Writer) {
	l.ZeroLog = l.ZeroLog.Output(w)
}

func (l *EchoLogger) Prefix() string {
	return ""
}

func (l *EchoLogger) SetPrefix(p string) {}

func (l *EchoLogger) Level() log.Lvl {
	return log.INFO
}

func (l *EchoLogger) SetLevel(v log.Lvl) {}

func (l *EchoLogger) SetHeader(h string) {}

func (l *EchoLogger) Print(i ...interface{}) {
	l.ZeroLog.Print(i...)
}

func (l *EchoLogger) Printf(format string, args ...interface{}) {
	l.ZeroLog.Printf(format, args...)
}

func (l *EchoLogger) Printj(j log.JSON) {
	l.ZeroLog.Info().Fields(j).Msg("")
}

func (l *EchoLogger) Debug(i ...interface{}) {
	l.ZeroLog.Debug().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Debugf(format string, args ...interface{}) {
	l.ZeroLog.Debug().Msgf(format, args...)
}

func (l *EchoLogger) Debugj(j log.JSON) {
	l.ZeroLog.Debug().Fields(j).Msg("")
}

func (l *EchoLogger) Info(i ...interface{}) {
	l.ZeroLog.Info().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Infof(format string, args ...interface{}) {
	l.ZeroLog.Info().Msgf(format, args...)
}

func (l *EchoLogger) Infoj(j log.JSON) {
	l.ZeroLog.Info().Fields(j).Msg("")
}

func (l *EchoLogger) Warn(i ...interface{}) {
	l.ZeroLog.Warn().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Warnf(format string, args ...interface{}) {
	l.ZeroLog.Warn().Msgf(format, args...)
}

func (l *EchoLogger) Warnj(j log.JSON) {
	l.ZeroLog.Warn().Fields(j).Msg("")
}

func (l *EchoLogger) Error(i ...interface{}) {
	l.ZeroLog.Error().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Errorf(format string, args ...interface{}) {
	l.ZeroLog.Error().Msgf(format, args...)
}

func (l *EchoLogger) Errorj(j log.JSON) {
	l.ZeroLog.Error().Fields(j).Msg("")
}

func (l *EchoLogger) Fatal(i ...interface{}) {
	l.ZeroLog.Fatal().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Fatalj(j log.JSON) {
	l.ZeroLog.Fatal().Fields(j).Msg("")
}

func (l *EchoLogger) Fatalf(format string, args ...interface{}) {
	l.ZeroLog.Fatal().Msgf(format, args...)
}

func (l *EchoLogger) Panic(i ...interface{}) {
	l.ZeroLog.Panic().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Panicj(j log.JSON) {
	l.ZeroLog.Panic().Fields(j).Msg("")
}

func (l *EchoLogger) Panicf(format string, args ...interface{}) {
	l.ZeroLog.Panic().Msgf(format, args...)
}
