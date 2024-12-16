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

package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-config-connector/dev/tools/aiframework/api"
)

type Tools struct {
	projectRoot string
	gitRepo     string
}

func New(projectRoot, gitRepo string) *Tools {
	return &Tools{
		projectRoot: projectRoot,
		gitRepo:     gitRepo,
	}
}

func (t *Tools) ReadFile(path string) ([]byte, error) {
	fullPath := filepath.Join(t.projectRoot, path)
	return os.ReadFile(fullPath)
}

func (t *Tools) WriteFile(path string, content []byte) error {
	fullPath := filepath.Join(t.projectRoot, path)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	return os.WriteFile(fullPath, content, 0644)
}

func (t *Tools) FindRelatedFiles(resource string) ([]string, error) {
	var files []string

	// Look for files in common locations
	patterns := []string{
		fmt.Sprintf("*/%s_types.go", strings.ToLower(resource)),
		fmt.Sprintf("*/%s_mapping.go", strings.ToLower(resource)),
		fmt.Sprintf("*/%s_controller.go", strings.ToLower(resource)),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(t.projectRoot, pattern))
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			// Convert to relative path
			relPath, _ := filepath.Rel(t.projectRoot, match)
			files = append(files, relPath)
		}
	}

	return files, nil
}

func (t *Tools) GenerateKRMTypes(opts *api.GenerateTypesOptions) error {
	cmd := exec.Command("go", "run", "./dev/tools/controllerbuilder",
		"generate-types",
		"--service", opts.Service,
		"--resource", opts.Resource,
		"--api-version", opts.APIVersion)
	cmd.Dir = t.projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("generating KRM types: %w\nOutput: %s", err, output)
	}

	return nil
}

func (t *Tools) GenerateMapper(opts *api.GenerateMapperOptions) error {
	cmd := exec.Command("go", "run", "./dev/tools/controllerbuilder",
		"generate-mapper",
		"--service", opts.Service,
		"--api-version", opts.APIVersion,
		"--output-dir", opts.OutputDir)
	cmd.Dir = t.projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("generating mapper: %w\nOutput: %s", err, output)
	}

	return nil
}

func (t *Tools) CreatePR(branch string, files []string, description string) error {
	// Create and checkout new branch
	if err := t.gitCommand("checkout", "-b", branch); err != nil {
		return err
	}

	// Add files
	for _, file := range files {
		if err := t.gitCommand("add", file); err != nil {
			return err
		}
	}

	// Commit changes
	if err := t.gitCommand("commit", "-m", description); err != nil {
		return err
	}

	// Push branch
	if err := t.gitCommand("push", "origin", branch); err != nil {
		return err
	}

	// Create PR using gh CLI if available
	if err := t.createGitHubPR(description); err != nil {
		return fmt.Errorf("creating PR: %w", err)
	}

	return nil
}

func (t *Tools) gitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = t.projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s failed: %w\nOutput: %s",
			strings.Join(args, " "), err, output)
	}
	return nil
}

func (t *Tools) createGitHubPR(description string) error {
	cmd := exec.Command("gh", "pr", "create",
		"--title", description,
		"--body", description)
	cmd.Dir = t.projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("creating GitHub PR: %w\nOutput: %s", err, output)
	}

	return nil
}
