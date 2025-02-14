// Copyright 2025 Google LLC
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

package toolbot

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/klog/v2"
)

type mapperFunc struct {
	FilePath   string
	Definition []string
}

// EnhanceWithMappers is an enhancer that finds mapper functions listed under "// mappers:" in fuzzer files
type EnhanceWithMappers struct {
	srcDirectory string
	mappers      map[string][]*mapperFunc
}

// NewEnhanceWithMappers creates a new EnhanceWithMappers
func NewEnhanceWithMappers(srcDirectory string) (*EnhanceWithMappers, error) {
	x := &EnhanceWithMappers{
		srcDirectory: srcDirectory,
		mappers:      make(map[string][]*mapperFunc),
	}
	return x, nil
}

var _ Enhancer = &EnhanceWithMappers{}

// EnhanceDataPoint enhances the data point by adding matching mapper function definitions
func (x *EnhanceWithMappers) EnhanceDataPoint(ctx context.Context, p *DataPoint) error {
	if p.Type != "fuzz-gen" { // Only enhance if this is a fuzz-gen tool
		return nil
	}

	mappersStr := p.Input["mappers"]
	if mappersStr == "" {
		return nil
	}

	// Split the mappers string into individual mapper names
	var mapperNames []string
	for _, name := range strings.Split(mappersStr, ",") {
		name = strings.TrimSpace(name)
		if name != "" {
			mapperNames = append(mapperNames, name)
		}
	}

	if len(mapperNames) == 0 {
		return nil
	}

	// Find the mapper functions
	var definitions []string
	for i, mapperName := range mapperNames {
		if i > 0 {
			definitions = append(definitions, "") // Add blank line between mappers
		}
		if mapperFuncs := x.findMapperFunction(mapperName); len(mapperFuncs) > 0 {
			definitions = append(definitions, mapperFuncs...)
		} else {
			klog.Infof("unable to find mapper function %q", mapperName)
			return fmt.Errorf("unable to find mapper function %q", mapperName)
		}
	}

	if len(definitions) > 0 {
		p.SetInput("mappers.definition", strings.Join(definitions, "\n"))
	}

	return nil
}

// findMapperFunction finds a specific mapper function definition
func (x *EnhanceWithMappers) findMapperFunction(mapperName string) []string {
	var definition []string
	filepath.WalkDir(x.srcDirectory, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(p) != ".go" {
			return nil
		}

		b, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("reading file %q: %w", p, err)
		}
		r := bytes.NewReader(b)
		br := bufio.NewReader(r)

		for {
			line, err := br.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("scanning file %q: %w", p, err)
			}
			line = strings.TrimSuffix(line, "\n")

			if strings.Contains(line, fmt.Sprintf("func %s", mapperName)) {
				definition = append(definition, line)

				indent := 0
				for _, r := range line {
					if r == '{' {
						indent++
					}
				}

				// Read until we find the matching closing brace
				for indent > 0 {
					line, err = br.ReadString('\n')
					if err != nil {
						if err == io.EOF {
							break
						}
						return fmt.Errorf("scanning file %q: %w", p, err)
					}
					line = strings.TrimSuffix(line, "\n")
					definition = append(definition, line)

					for _, r := range line {
						if r == '{' {
							indent++
						}
						if r == '}' {
							indent--
						}
					}
				}

				return io.EOF // Found what we're looking for, stop walking
			}
		}
		return nil
	})
	return definition
}
