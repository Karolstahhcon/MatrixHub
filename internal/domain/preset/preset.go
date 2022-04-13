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

package preset

import "context"

type RegistrySpec struct {
	ID          int
	Name        string
	Description string
	URL         string
}

type ProjectSpec struct {
	Name         string
	Organization string
	Registry     string
}

type ModelSpec struct {
	Project        string
	Name           string
	Recommended    bool
	Size           int64
	ParameterCount int64
	Labels         []LabelSpec
}

type LabelSpec struct {
	Category string
	Name     string
}

type Manifest struct {
	Registries []RegistrySpec
	Projects   []ProjectSpec
	Models     []ModelSpec
}
