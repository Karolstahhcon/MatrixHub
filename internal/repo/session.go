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

package repo

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/config"
	infradb "github.com/matrixhub-ai/matrixhub/internal/infra/db"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

type sessionRepo struct {
	*scs.SessionManager
	user.SessionConfig
	stopCleanup func()
	closeOnce   sync.Once
}

type sessionStore interface {
	scs.Store
	StopCleanup()
}

func (s *sessionRepo) LoadSession(ctx context.Context) (context.Context, error) {
	token := utils.GetCookieFromContext(ctx)
	return s.Manager().Load(ctx, token)
}

func (s *sessionRepo) Manager() *scs.SessionManager {
	return s.SessionManager
}

func (s *sessionRepo) GetSessionCookie() scs.SessionCookie {
	return s.Cookie
}

func (s *sessionRepo) GetSessionConfig() user.SessionConfig {
	return s.SessionConfig
}

func (s *sessionRepo) WriteSessionCookie(ctx context.Context, token string, expiry time.Time) error {
	sessionCookie := s.GetSessionCookie()
	cookie := &http.Cookie{
		Value:       token,
		Name:        sessionCookie.Name,
		Domain:      sessionCookie.Domain,
		HttpOnly:    sessionCookie.HttpOnly,
		Path:        sessionCookie.Path,
		SameSite:    sessionCookie.SameSite,
