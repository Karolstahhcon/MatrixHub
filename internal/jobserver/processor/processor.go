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

package processor

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/matrixhub-ai/matrixhub/internal/domain/job"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
)

// Processor identifies a processor implementation family (sync, future cleanup, ...).
type Processor string

const (
	ProcessorSyncPolicy Processor = "syncPolicy"
	ProcessorSyncTask   Processor = "syncTask"
	ProcessorSyncJob    Processor = "syncJob"
)

// ExecuteFn runs one claimed policy execution (insert task/jobs + git work happens inside service).
type ExecuteFn func(ctx context.Context, policyID int, triggerType int) error

// PollDueFn lists due work and performs CAS advance in the service / repo layer.
type PollDueFn func(ctx context.Context, nowMs int64) ([]job.DueJob, error)

// ClaimOneFn atomically claims one item selected by its ID.
type ClaimOneFn func(ctx context.Context, id int) (job.DueJob, bool, error)

// ReleaseClaimFn returns unfinished work to pending so polling can retry it.
type ReleaseClaimFn func(ctx context.Context, id int) error

// PolicyLocker serializes execution per logical policy key (single-replica in-memory).
type PolicyLocker interface {
	TryAcquire(key string, fireAt time.Time) bool
	Release(key string)
}

type memLocker struct {
	m sync.Map // key string -> struct{}
}

// NewMemLocker returns an in-process PolicyLocker (lost on restart).
func NewMemLocker() PolicyLocker {
	return &memLocker{}
}

func (l *memLocker) TryAcquire(key string, _ time.Time) bool {
	_, loaded := l.m.LoadOrStore(key, struct{}{})
	return !loaded
}

func (l *memLocker) Release(key string) {
	l.m.Delete(key)
}

// Adapter is the minimal lifecycle surface exposed to JobServer.
type Adapter interface {
	Processor() Processor
	Trigger(id int) bool
	Start(ctx context.Context)
	Wait()
}

// processor holds the shared poll loop, semaphore, per-policy lock, and execute path for one Processor value.
type processor struct {
	processor Processor
	interval  time.Duration
	sem       chan struct{}
	locker    PolicyLocker
	taskMax   time.Duration
	execute   ExecuteFn
	pollDueFn PollDueFn
	claimOne  ClaimOneFn
	release   ReleaseClaimFn
	runOneFn  func(ctx context.Context, d job.DueJob) // optional override
	triggerCh chan int

	wg   sync.WaitGroup
	done chan struct{}
}

func newProcessor(
	proc Processor,
	interval time.Duration,
	maxConcurrent int,
	taskMax time.Duration,
	execute ExecuteFn,
	pollDueFn PollDueFn,
) *processor {
	if maxConcurrent <= 0 {
		maxConcurrent = 5
	}
	if interval <= 0 {
		interval = 10 * time.Second
	}
	if taskMax <= 0 {
		taskMax = 2 * time.Hour
	}
	return &processor{
		processor: proc,
		interval:  interval,
		sem:       make(chan struct{}, maxConcurrent),
		locker:    NewMemLocker(),
		taskMax:   taskMax,
		execute:   execute,
		pollDueFn: pollDueFn,
		triggerCh: make(chan int, 256),
	}
}

func (b *processor) Processor() Processor { return b.processor }

// Trigger requests immediate processing of one ID. Processors with a targeted
// claim function CAS-claim that row; periodic scanning remains the recovery path.
func (b *processor) Trigger(id int) bool {
	select {
	case b.triggerCh <- id:
		return true
	default:
