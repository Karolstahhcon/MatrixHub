// Copyright The MatrixHub Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	sugar *zap.SugaredLogger
	name  string
}

type LoggerPayload interface {
	LoggerOutput(simple bool) (keysAndValues []any)
}

func NewLogger(name string) *Logger {
	loggerLock.Lock()
	defer loggerLock.Unlock()

	return logger.LoggerWithOptions(name, zap.AddCallerSkip(-1))
}
