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
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-config-connector/dev/tools/controllerbuilder/pkg/codegen"
	"github.com/GoogleCloudPlatform/k8s-config-connector/dev/tools/controllerbuilder/pkg/options"

	"k8s.io/klog/v2"
)

// EnhanceWithIgnoreFields is an enhancer that reads ignored fields from config
type EnhanceWithIgnoreFields struct {
	ignoredFields map[string][]string
}

// NewEnhanceWithIgnoreFields creates a new EnhanceWithIgnoreFields
func NewEnhanceWithIgnoreFields() (*EnhanceWithIgnoreFields, error) {
	repoRoot, err := options.RepoRoot()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join(repoRoot, "dev", "tools", "controllerbuilder", "config")
	configs, err := codegen.LoadAllConfigs(configDir)
	if err != nil {
		return nil, err
	}

	ignoredFields := make(map[string][]string)
	for _, config := range configs {
		for _, resource := range config.Resources {
			fqn := fmt.Sprintf("%s.%s", config.Service, resource.ProtoName)
			ignoredFields[fqn] = resource.IgnoredFields
		}
	}

	return &EnhanceWithIgnoreFields{
		ignoredFields: ignoredFields,
	}, nil
}

var _ Enhancer = &EnhanceWithIgnoreFields{}

// EnhanceDataPoint enhances the data point by adding ignored fields from config
func (x *EnhanceWithIgnoreFields) EnhanceDataPoint(ctx context.Context, p *DataPoint) error {
	if p.Type != "fuzz-gen" {
		return nil
	}

	protoMsg := p.Input["proto.message"]
	if protoMsg == "" {
		return nil
	}

	ignoredFields, ok := x.ignoredFields[protoMsg]
	if !ok {
		return nil
	}

	p.SetInput("unimplemented.fields", strings.Join(ignoredFields, ","))
	klog.Infof("Found unimplemented fields for %q: %v", protoMsg, ignoredFields)
	return nil
}
