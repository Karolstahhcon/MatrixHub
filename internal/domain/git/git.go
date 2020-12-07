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

package git

import (
	"context"
	"io"
	"time"
)

// Git-related types for model version control

// Revision represents a Git reference (branch or tag).
type Revision struct {
	Name string `json:"name"`
}

// Revisions contains branches and tags.
type Revisions struct {
	Branches []*Revision `json:"branches"`
	Tags     []*Revision `json:"tags"`
}

type CommitOperationType string

const (
	// CommitOperationAdd adds or updates a file.
	CommitOperationAdd CommitOperationType = "add"
	// CommitOperationDelete deletes a file.
	CommitOperationDelete CommitOperationType = "delete"
)

// CommitOperation represents a single operation in a commit.
type CommitOperation struct {
	Type    CommitOperationType
	Path    string
	Content []byte // file content for add operations
}

// Commit represents a Git commit.
type Commit struct {
	ID             string    `json:"id"`
	Message        string    `json:"message"`
	AuthorName     string    `json:"authorName"`
	AuthorEmail    string    `json:"authorEmail"`
	AuthorDate     time.Time `json:"authorDate"`
	CommitterName  string    `json:"committerName"`
	CommitterEmail string    `json:"committerEmail"`
	CommitterDate  time.Time `json:"committerDate"`
	Diff           string    `json:"diff"`
	ParentCommit   string
	CreatedAt      time.Time `json:"createdAt"`
}

// Diff represents a file change in a commit.
type Diff struct {
	Diff    string `json:"diff"`
	Deleted bool   `json:"deleted"`
	NewPath string `json:"newPath"`
	OldPath string `json:"oldPath"`
}

// FileType represents the type of file in the Git tree.
type FileType int

const (
	FileTypeDir  FileType = 0
	FileTypeFile FileType = 1
)

// TreeEntry represents a file or directory in the Git tree.
type TreeEntry struct {
	Name   string   `json:"name"`
	Type   FileType `json:"type"`
	Size   int64    `json:"size"`
	Path   string   `json:"path"`
	Hash   string   `json:"hash"`
	IsLFS  bool     `json:"isLFS"`
	URL    string   `json:"url,omitempty"`
	Commit *Commit  `json:"commit,omitempty"`
}

// RepoMetadataFiles holds raw metadata-related files loaded from a git repo.
// File interpretation stays in the model domain; git only loads bytes.
type RepoMetadataFiles struct {
	ReadmeContent        []byte
	ConfigJSON           []byte
	SafetensorsIndexJSON []byte
	SafetensorsFiles     map[string][]byte
	// SafetensorsSizes holds LFS pointer sizes for safetensors files whose
	// headers could not be read. It lets the model domain fall back to a
	// size-based estimate instead of reporting no parameter count at all.
	SafetensorsSizes map[string]int64
	Size             int64
}

// OrphanedRepo represents an orphaned Git repository on disk
// that has no corresponding record in the database.
type OrphanedRepo struct {
	Path         string
	Type         string
	ProjectName  string
	ResourceName string
	SizeBytes    int64
}

// OrphanedLFS represents an orphaned LFS object on disk
// that is not referenced by any Git repository.
type OrphanedLFS struct {
	OID       string
	SizeBytes int64
	Path      string
}

// BasicCredential holds username/password for remote git authentication.
type BasicCredential struct {
	Username string
	Password string
}

type GitRepository struct {
	RemoteRegistryURL  string
	RemoteProjectName  string
	RemoteResourceName string
