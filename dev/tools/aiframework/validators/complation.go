// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validators

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-config-connector/dev/tools/aiframework/api"
)

type CompilationValidator struct {
	// Root directory of the project
	ProjectRoot string
}

func NewCompilationValidator(projectRoot string) *CompilationValidator {
	return &CompilationValidator{
		ProjectRoot: projectRoot,
	}
}

func (v *CompilationValidator) Validate(ctx *api.UpdateContext) (*api.ValidationResult, error) {
	// Change to project root directory
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = v.ProjectRoot

	output, err := cmd.CombinedOutput()

	result := &api.ValidationResult{
		Success: err == nil,
		Logs:    string(output),
	}

	if err != nil {
		// Parse compiler errors into structured format
		result.Errors = v.parseCompilerErrors(string(output))
	}

	return result, nil
}

func (v *CompilationValidator) parseCompilerErrors(output string) []api.CompileError {
	var errors []api.CompileError

	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, ":") {
			// Parse error line into structured format
			// Example: file.go:23:45: undefined: foo
			parts := strings.SplitN(line, ":", 4)
			if len(parts) >= 4 {
				// Convert file path to relative path
				filePath, _ := filepath.Rel(v.ProjectRoot, parts[0])

				errors = append(errors, api.CompileError{
					File:    filePath,
					Line:    parts[1],
					Column:  parts[2],
					Message: strings.TrimSpace(parts[3]),
				})
			}
		}
	}

	return errors
}
