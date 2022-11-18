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
	"fmt"
	"github.com/pressly/goose/v3"
)

type GooseLogger struct {
	logger Logger
}

func GooseLoggerWrapper(logger Logger) goose.Logger {
	return &GooseLogger{
		logger: logger,
	}
}

func (gl *GooseLogger) Fatal(v ...interface{}) {
	gl.logger.Error(v)
}

func (gl *GooseLogger) Fatalf(format string, v ...interface{}) {
	gl.logger.Error(fmt.Sprintf(format, v...))
}

func (gl *GooseLogger) Print(v ...interface{}) {
	gl.logger.Info(v)
}

func (gl *GooseLogger) Println(v ...interface{}) {
	gl.logger.Info(v)
}

func (gl *GooseLogger) Printf(format string, v ...interface{}) {
	gl.logger.Info(fmt.Sprintf(format, v...))
}
