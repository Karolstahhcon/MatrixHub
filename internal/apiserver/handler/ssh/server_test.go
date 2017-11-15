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

package ssh_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/matrixhub-ai/hfd/pkg/authenticate"
	"github.com/matrixhub-ai/hfd/pkg/permission"
	"github.com/matrixhub-ai/hfd/pkg/receive"
	"github.com/matrixhub-ai/hfd/pkg/storage"
	"github.com/stretchr/testify/require"
	cryptossh "golang.org/x/crypto/ssh"

	backend "github.com/matrixhub-ai/matrixhub/internal/apiserver/handler/ssh"
)

func generateHostKey(t *testing.T) cryptossh.Signer {
	t.Helper()
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	signer, err := cryptossh.NewSignerFromKey(privateKey)
	require.NoError(t, err)
	return signer
}

func runGit(t *testing.T, dir string, environment []string, arguments ...string) string {
	t.Helper()
	command := exec.CommandContext(t.Context(), "git", arguments...)
	command.Dir = dir
	command.Env = append(os.Environ(), environment...)
	output, err := command.CombinedOutput()
	require.NoError(t, err, "git %s failed: %s", strings.Join(arguments, " "), output)
	return string(output)
}

func startSSHServer(t *testing.T, server *backend.Server) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, listener.Close()) })
	go func() { _ = server.Serve(t.Context(), listener) }()
	return listener.Addr().String()
}

func TestSSHGitCloneAndPush(t *testing.T) {
	dataDir := t.TempDir()
	clientDir := t.TempDir()
	repositoryStorage := storage.NewStorage(storage.WithRootDir(dataDir))
	repositoryName := "org/model.git"
	repositoryPath := filepath.Join(repositoryStorage.RepositoriesDir(), repositoryName)
	require.NoError(t, os.MkdirAll(filepath.Dir(repositoryPath), 0o750))
	runGit(t, "", nil, "init", "--bare", "--initial-branch=main", repositoryPath)

	var (
		operations      []permission.Operation
		operationsMutex sync.Mutex
	)
	permissionHook := func(_ context.Context, operation permission.Operation, repo string, _ permission.Context) (bool, error) {
		if strings.TrimPrefix(repo, "/") != repositoryName {
			return false, fmt.Errorf("unexpected repository %q", repo)
		}
		operationsMutex.Lock()
		defer operationsMutex.Unlock()
		operations = append(operations, operation)
		return true, nil
	}
	preReceive := make(chan []receive.RefUpdate, 1)
	postReceive := make(chan []receive.RefUpdate, 1)
	server := backend.NewServer(
		backend.WithHostKey(generateHostKey(t)),
		backend.WithStorage(repositoryStorage),
		backend.WithPermissionHookFunc(permissionHook),
		backend.WithServices(nil),
		backend.WithMirror(nil),
		backend.WithPreReceiveHookFunc(func(_ context.Context, repo string, updates []receive.RefUpdate) (bool, error) {
			if strings.TrimPrefix(repo, "/") != repositoryName {
				return false, fmt.Errorf("unexpected repository %q", repo)
			}
			preReceive <- updates
			return true, nil
		}),
		backend.WithPostReceiveHookFunc(func(_ context.Context, repo string, updates []receive.RefUpdate) error {
			if strings.TrimPrefix(repo, "/") != repositoryName {
				return fmt.Errorf("unexpected repository %q", repo)
			}
			postReceive <- updates
			return nil
		}),
	)
	address := startSSHServer(t, server)
	port := address[strings.LastIndex(address, ":")+1:]
	environment := []string{
