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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	EnvMatrixHubSSHHost = "MATRIXHUB_SSH_HOST"
	EnvMatrixHubSSHPort = "MATRIXHUB_SSH_PORT"
	DefaultSSHHost      = "127.0.0.1"
	DefaultSSHPort      = 2222
)

// SSHClientFixture holds an isolated key, known_hosts file, and SSH config for
// Git subprocesses.
type SSHClientFixture struct {
	Host           string
	Port           int
	PrivateKeyPath string
	PublicKey      string
	ConfigPath     string
}

func GetSSHHost() string {
	if host := os.Getenv(EnvMatrixHubSSHHost); host != "" {
		return host
	}
	return DefaultSSHHost
}

func GetSSHPort() (int, error) {
	value := os.Getenv(EnvMatrixHubSSHPort)
	if value == "" {
		return DefaultSSHPort, nil
	}
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid %s %q", EnvMatrixHubSSHPort, value)
	}
	return port, nil
}

// NewSSHClientFixture generates a client key and pins the ephemeral server's
// host key in a temporary known_hosts file.
func NewSSHClientFixture(ctx context.Context, root string) (*SSHClientFixture, error) {
	port, err := GetSSHPort()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, err
	}

	privateKeyPath := filepath.Join(root, "id_ed25519")
	keygen := RunCommand(ctx, root, nil, "ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-f", privateKeyPath)
	if keygen.Err != nil {
		return nil, errorsFromCommand(keygen)
	}
	publicKeyBytes, err := os.ReadFile(privateKeyPath + ".pub")
	if err != nil {
		return nil, err
	}
