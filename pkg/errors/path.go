// Copyright 2025 s3-mcp Authors
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

package errors

import (
	"runtime"
	"strings"
)

var (
	goPath    string
	goPathLen int
)

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return
	}
	size := len(file)
	suffixLen := len("pkg/errors/path.go") // also trim project directory here
	goPath = file[:size-suffixLen]
	goPathLen = len(goPath)
}

func trimGOPATH(path string) string {
	if strings.HasPrefix(path, goPath) {
		return path[goPathLen:]
	}
	return path
}
