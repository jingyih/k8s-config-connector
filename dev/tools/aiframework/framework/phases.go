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

package framework

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-config-connector/dev/tools/aiframework/api"
)

func (f *Framework) analyzeChange(ctx *api.UpdateContext) error {
	analysis, err := f.ai.AnalyzeChange(ctx)
	if err != nil {
		return fmt.Errorf("analyzing change: %w", err)
	}

	ctx.State.Analysis = analysis
	return nil
}

func (f *Framework) implementChange(ctx *api.UpdateContext) error {
	plan, err := f.ai.PlanImplementation(ctx)
	if err != nil {
		return fmt.Errorf("planning implementation: %w", err)
	}

	for _, step := range plan.Steps {
		if err := f.executeStep(ctx, step); err != nil {
			return fmt.Errorf("executing step: %w", err)
		}
	}

	return nil
}

func (f *Framework) executeStep(ctx *api.UpdateContext, step api.ImplementationStep) error {
	switch step.Type {
	case api.StepGenerateKRM:
		if err := f.tools.GenerateKRMTypes(step.TypesOptions); err != nil {
			return err
		}
		ctx.State.ModifiedFiles = append(ctx.State.ModifiedFiles, step.FilePath)

	case api.StepGenerateMapper:
		if err := f.tools.GenerateMapper(step.MapperOptions); err != nil {
			return err
		}
		ctx.State.ModifiedFiles = append(ctx.State.ModifiedFiles, step.FilePath)

	case api.StepModifyFile:
		content, err := f.tools.ReadFile(step.FilePath)
		if err != nil {
			return err
		}

		newContent := applyModifications(content, step.Modifications)
		if err := f.tools.WriteFile(step.FilePath, newContent); err != nil {
			return err
		}
		ctx.State.ModifiedFiles = append(ctx.State.ModifiedFiles, step.FilePath)
	}

	return nil
}

func (f *Framework) validateChange(ctx *api.UpdateContext) error {
	var results []*api.ValidationResult

	for _, validator := range f.validators {
		result, err := validator.Validate(ctx)
		if err != nil {
			return fmt.Errorf("validation error: %w", err)
		}
		results = append(results, result)
	}

	ctx.State.LastValidationResults = results

	decision, err := f.ai.ReviewValidationResult(ctx, results[0]) // TODO: Handle multiple results
	if err != nil {
		return fmt.Errorf("reviewing validation results: %w", err)
	}

	if !decision.Success {
		ctx.State.PlannedModifications = decision.Modifications
		return fmt.Errorf("validation failed, modifications planned")
	}

	return nil
}

func (f *Framework) createPR(ctx *api.UpdateContext) error {
	description, err := f.ai.PreparePRDescription(ctx)
	if err != nil {
		return fmt.Errorf("preparing PR description: %w", err)
	}

	branch := fmt.Sprintf("ai-update/%s/%s",
		strings.ToLower(ctx.APIChange.MessageName),
		ctx.APIChange.FieldName)

	return f.tools.CreatePR(branch, ctx.State.ModifiedFiles, description)
}

func applyModifications(content []byte, mods []api.Modification) []byte {
	lines := strings.Split(string(content), "\n")

	// Apply modifications in reverse order to handle line numbers correctly
	for i := len(mods) - 1; i >= 0; i-- {
		mod := mods[i]
		newLines := strings.Split(mod.NewContent, "\n")
		lines = append(lines[:mod.StartLine-1], append(newLines, lines[mod.EndLine:]...)...)
	}

	return []byte(strings.Join(lines, "\n"))
}
