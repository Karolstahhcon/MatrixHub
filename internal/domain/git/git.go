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
