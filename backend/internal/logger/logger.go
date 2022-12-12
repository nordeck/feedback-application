/*
 *  Copyright 2022 Nordeck IT + Consulting GmbH
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and  limitations
 *  under the License.
 *
 */

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance Logger
)

type Logger interface {
	Info(...interface{})
	Debug(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	DPanic(...interface{})
	Panic(...interface{})
	OnExit()
}

type AppLogger struct {
	logger *zap.SugaredLogger
}

func Instance() Logger {
	once.Do(func() {
		core := createCore()
		zapLogger := zap.New(core, zap.AddCaller())
		instance = &AppLogger{zapLogger.Sugar()}
	})
	return instance
}

func createCore() zapcore.Core {
	encoder := createEncoder()
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLevelEnablerFunc(zap.InfoLevel)),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLevelEnablerFunc(zap.DebugLevel)),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLevelEnablerFunc(zap.WarnLevel)),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLevelEnablerFunc(zap.ErrorLevel)),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLevelEnablerFunc(zap.FatalLevel)),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLevelEnablerFunc(zap.DPanicLevel)),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLevelEnablerFunc(zap.PanicLevel)),
	)
	return core
}

func createEncoder() zapcore.Encoder {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
	return encoder
}

func getLevelEnablerFunc(level zapcore.Level) zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl == level
	}
}

func (log *AppLogger) Info(args ...interface{}) {
	log.logger.Info(args)
}

func (log *AppLogger) Debug(args ...interface{}) {
	log.logger.Debug(args)
}

func (log *AppLogger) Warn(args ...interface{}) {
	log.logger.Warn(args)
}

func (log *AppLogger) Error(args ...interface{}) {
	log.logger.Error(args)
}
func (log *AppLogger) Fatal(args ...interface{}) {
	log.logger.Fatal(args)
}

func (log *AppLogger) DPanic(args ...interface{}) {
	log.logger.DPanic(args)
}

func (log *AppLogger) Panic(args ...interface{}) {
	log.logger.Panic(args)
}

func (log *AppLogger) OnExit() {
	err := log.logger.Sync()
	if err != nil {
		log.Debug(err)
	}
}
