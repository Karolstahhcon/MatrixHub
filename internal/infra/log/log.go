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
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"k8s.io/klog/v2"
)

var (
	loggerLock     sync.Mutex
	zapLogger      *zap.Logger
	zapSugarLogger *zap.SugaredLogger

	logger *Logger
)

func init() {
	debug := strings.ToLower(os.Getenv("DEBUG")) == "true"
	if err := SetLoggerWithConfig(debug, Config{}); err != nil {
		panic(fmt.Sprintf("log init failed: %s", err))
	}
}

func SetLoggerWithConfig(debug bool, config Config) error {
	loggerLock.Lock()
	defer loggerLock.Unlock()

	sink, err := openSinks(config)
	if err != nil {
		return err
