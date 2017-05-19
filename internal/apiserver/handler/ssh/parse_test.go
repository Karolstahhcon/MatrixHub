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

package ssh

import (
	"testing"

	"github.com/matrixhub-ai/hfd/pkg/authenticate"
	"github.com/matrixhub-ai/hfd/pkg/repository"
	"github.com/stretchr/testify/require"
	cryptossh "golang.org/x/crypto/ssh"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name           string
		command        string
		wantService    string
		wantRepository string
		wantOperation  string
		wantError      bool
	}{
		{name: "upload pack quoted", command: "git-upload-pack '/org/model.git'", wantService: repository.GitUploadPack, wantRepository: "/org/model.git"},
		{name: "receive pack", command: "git-receive-pack /org/model.git", wantService: repository.GitReceivePack, wantRepository: "/org/model.git"},
		{name: "lfs download", command: "git-lfs-authenticate '/org/model.git' download", wantService: repository.GitLFSAuthenticate, wantRepository: "/org/model.git", wantOperation: "download"},
		{name: "lfs transfer", command: "git-lfs-transfer /org/model.git upload", wantService: repository.GitLFSTransfer, wantRepository: "/org/model.git", wantOperation: "upload"},
		{name: "empty", wantError: true},
		{name: "missing repository", command: "git-upload-pack", wantError: true},
		{name: "missing lfs operation", command: "git-lfs-authenticate /org/model.git", wantError: true},
		{name: "unsupported", command: "git-status /org/model.git", wantError: true},
	}
