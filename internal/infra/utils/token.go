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

package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	TokenPrefix      = "mh_"
	RobotTokenPrefix = "robot_"
	tokenRandBytes   = 24 // 192 bits of entropy → 32 base64url chars

	// MaxTokenRetries caps the retry loop in the (astronomically unlikely) event of
	// a collision. In practice the loop never executes more than once.
	MaxTokenRetries = 3
