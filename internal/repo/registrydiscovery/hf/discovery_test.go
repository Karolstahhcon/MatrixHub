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

package hf

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matrixhub-ai/matrixhub/internal/domain/registry"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registrydiscovery"
)

func TestDiscovery_ListRepositories_Models(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/models" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("author") != "google" {
			t.Errorf("unexpected author: %s", r.URL.Query().Get("author"))
		}
		if r.URL.Query().Get("limit") != "100" {
			t.Errorf("unexpected limit: %s", r.URL.Query().Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"id":"google/bert-base"},{"id":"google/t5-base"}]`)
	}))
	defer ts.Close()

	d := New()

	reg := &registry.Registry{URL: ts.URL}
	repos, err := d.ListRepositories(context.Background(), reg, registrydiscovery.Filter{
		Namespace:    "google",
		ResourceType: "model",
	})
	if err != nil {
		t.Fatalf("ListRepositories error = %v", err)
	}
	if len(repos) != 2 {
		t.Fatalf("got %d repos, want 2", len(repos))
	}
	if repos[0].Namespace != "google" || repos[0].Name != "bert-base" || repos[0].ResourceType != "model" {
		t.Errorf("repos[0] = %+v, want google/bert-base model", repos[0])
	}
	if repos[1].Namespace != "google" || repos[1].Name != "t5-base" || repos[1].ResourceType != "model" {
		t.Errorf("repos[1] = %+v, want google/t5-base model", repos[1])
	}
}

func TestDiscovery_ListRepositories_Datasets(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/datasets" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"id":"google/wiki"}]`)
	}))
	defer ts.Close()

	d := New()

	reg := &registry.Registry{URL: ts.URL}
	repos, err := d.ListRepositories(context.Background(), reg, registrydiscovery.Filter{
		Namespace:    "google",
		ResourceType: "dataset",
	})
	if err != nil {
		t.Fatalf("ListRepositories error = %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("got %d repos, want 1", len(repos))
	}
	if repos[0].Name != "wiki" {
		t.Errorf("repos[0].Name = %s, want wiki", repos[0].Name)
	}
}

func TestDiscovery_ListRepositories_WithCredential(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer secret" {
			t.Errorf("unexpected auth header: %s", auth)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"id":"google/bert-base"}]`)
	}))
	defer ts.Close()

	d := New()

	reg := &registry.Registry{
		URL:            ts.URL,
		CredentialType: registry.CredentialTypeBasic,
		AuthInfo:       `{"username":"token","password":"secret"}`,
	}
	repos, err := d.ListRepositories(context.Background(), reg, registrydiscovery.Filter{
		Namespace:    "google",
		ResourceType: "model",
	})
	if err != nil {
		t.Fatalf("ListRepositories error = %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("got %d repos, want 1", len(repos))
	}
}

func TestEndpointBase(t *testing.T) {
	tests := []struct {
		name    string
		reg     *registry.Registry
		want    string
		wantErr bool
	}{
		{
			name: "registry url is used as-is",
