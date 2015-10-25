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

package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// CommandResult captures a black-box CLI invocation without including its
// environment, which may contain access tokens.
type CommandResult struct {
	Name   string
	Args   []string
	Stdout string
	Stderr string
	Err    error
}

// FailureMessage returns diagnostics suitable for a Ginkgo assertion.
func (r CommandResult) FailureMessage() string {
	return fmt.Sprintf("command %q failed: %v\nstdout:\n%s\nstderr:\n%s",
		append([]string{r.Name}, r.Args...), r.Err, r.Stdout, r.Stderr)
}

// RunCommand executes a command directly, without a shell. Environment
// overrides are merged by key so host credentials cannot accidentally win
