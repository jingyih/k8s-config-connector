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

package client

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-config-connector/dev/tools/aiframework/api"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiAIClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiAIClient() (*GeminiAIClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}

	return &GeminiAIClient{
		client: client,
		model:  client.GenerativeModel("gemini-pro"),
	}, nil
}

func (c *GeminiAIClient) AnalyzeChange(ctx *api.UpdateContext) (*api.AnalysisResult, error) {
	prompt := fmt.Sprintf(`Analyze the following API change:
Message: %s
New Field: %s

Review the history of previous attempts if any:
%s

Determine which files need to be modified and what changes are required.`,
		ctx.APIChange.MessageName,
		ctx.APIChange.FieldName,
		formatHistory(ctx.History))

	resp, err := c.generateContent(prompt)
	if err != nil {
		return nil, err
	}

	// Parse AI response into AnalysisResult
	// This is a simplified example - you'd want more structured response parsing
	return &api.AnalysisResult{
		AffectedFiles:   parseAffectedFiles(resp),
		RequiredChanges: parseRequiredChanges(resp),
	}, nil
}

func (c *GeminiAIClient) PlanImplementation(ctx *api.UpdateContext) (*api.ImplementationPlan, error) {
	prompt := fmt.Sprintf(`Plan implementation for API change:
Message: %s
Field: %s

Analysis results:
%v

Previous attempts and their outcomes:
%s

Create a detailed implementation plan with specific steps.`,
		ctx.APIChange.MessageName,
		ctx.APIChange.FieldName,
		ctx.State.Analysis,
		formatHistory(ctx.History))

	resp, err := c.generateContent(prompt)
	if err != nil {
		return nil, err
	}

	return parsePlan(resp), nil
}

func (c *GeminiAIClient) ReviewValidationResult(ctx *api.UpdateContext, result *api.ValidationResult) (*api.ReviewDecision, error) {
	if result.Success {
		return &api.ReviewDecision{Success: true}, nil
	}

	prompt := fmt.Sprintf(`Review validation failures:
Errors:
%v

Previous attempts:
%s

Suggest modifications to fix these errors.`,
		formatErrors(result.Errors),
		formatHistory(ctx.History))

	resp, err := c.generateContent(prompt)
	if err != nil {
		return nil, err
	}

	return parseReviewDecision(resp), nil
}

func (c *GeminiAIClient) PreparePRDescription(ctx *api.UpdateContext) (string, error) {
	prompt := fmt.Sprintf(`Create PR description for:
API Change: %s.%s
Modified Files: %v
Implementation History: %s`,
		ctx.APIChange.MessageName,
		ctx.APIChange.FieldName,
		ctx.State.ModifiedFiles,
		formatHistory(ctx.History))

	resp, err := c.generateContent(prompt)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (c *GeminiAIClient) generateContent(prompt string) (string, error) {
	ctx := context.Background()
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	var result strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			result.WriteString(string(text))
		}
	}

	return result.String(), nil
}

// Helper functions for parsing AI responses
func formatHistory(history []api.AttemptResult) string {
	var result strings.Builder
	for i, attempt := range history {
		result.WriteString(fmt.Sprintf("Attempt %d (%s):\n", i+1, attempt.Phase))
		if attempt.Error != nil {
			result.WriteString(fmt.Sprintf("Error: %v\n", attempt.Error))
		}
		// Add more details as needed
	}
	return result.String()
}

func formatErrors(errors []api.CompileError) string {
	var result strings.Builder
	for _, err := range errors {
		result.WriteString(fmt.Sprintf("%s:%s:%s: %s\n",
			err.File, err.Line, err.Column, err.Message))
	}
	return result.String()
}

// Helper functions for parsing AI responses
func parseAffectedFiles(aiResponse string) []string {
	// TODO: Implement more sophisticated parsing
	// This is a basic implementation - should be enhanced based on actual AI response format
	var files []string
	lines := strings.Split(aiResponse, "\n")
	for _, line := range lines {
		if strings.Contains(line, "File:") {
			file := strings.TrimSpace(strings.TrimPrefix(line, "File:"))
			files = append(files, file)
		}
	}
	return files
}

func parseRequiredChanges(aiResponse string) []api.Change {
	// TODO: Implement more sophisticated parsing
	var changes []api.Change
	lines := strings.Split(aiResponse, "\n")
	currentChange := api.Change{}

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "Type:"):
			currentChange.Type = api.ChangeType(strings.TrimSpace(strings.TrimPrefix(line, "Type:")))
		case strings.HasPrefix(line, "File:"):
			currentChange.FilePath = strings.TrimSpace(strings.TrimPrefix(line, "File:"))
		case strings.HasPrefix(line, "Description:"):
			currentChange.Description = strings.TrimSpace(strings.TrimPrefix(line, "Description:"))
			// Assume this completes a change entry
			changes = append(changes, currentChange)
			currentChange = api.Change{}
		}
	}
	return changes
}

func parsePlan(aiResponse string) *api.ImplementationPlan {
	// TODO: Implement more sophisticated parsing
	var plan api.ImplementationPlan
	lines := strings.Split(aiResponse, "\n")
	var currentStep api.ImplementationStep

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "Step Type:"):
			if currentStep.Type != "" {
				plan.Steps = append(plan.Steps, currentStep)
			}
			currentStep = api.ImplementationStep{
				Type: api.StepType(strings.TrimSpace(strings.TrimPrefix(line, "Step Type:"))),
			}
		case strings.HasPrefix(line, "File:"):
			currentStep.FilePath = strings.TrimSpace(strings.TrimPrefix(line, "File:"))
		case strings.HasPrefix(line, "TypesOptions:"):
			// Parse TypesOptions
			currentStep.TypesOptions = &api.GenerateTypesOptions{}
			// TODO: Parse detailed options
		case strings.HasPrefix(line, "MapperOptions:"):
			// Parse MapperOptions
			currentStep.MapperOptions = &api.GenerateMapperOptions{}
			// TODO: Parse detailed options
		case strings.HasPrefix(line, "Modification:"):
			// Parse Modification
			mod := parseModification(strings.TrimPrefix(line, "Modification:"))
			currentStep.Modifications = append(currentStep.Modifications, mod)
		}
	}

	// Add the last step if any
	if currentStep.Type != "" {
		plan.Steps = append(plan.Steps, currentStep)
	}

	return &plan
}

func parseModification(modText string) api.Modification {
	// TODO: Implement more sophisticated parsing
	var mod api.Modification
	parts := strings.Split(modText, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch {
		case strings.HasPrefix(part, "StartLine:"):
			fmt.Sscanf(strings.TrimPrefix(part, "StartLine:"), "%d", &mod.StartLine)
		case strings.HasPrefix(part, "EndLine:"):
			fmt.Sscanf(strings.TrimPrefix(part, "EndLine:"), "%d", &mod.EndLine)
		case strings.HasPrefix(part, "Content:"):
			mod.NewContent = strings.TrimPrefix(part, "Content:")
		}
	}
	return mod
}

func parseReviewDecision(aiResponse string) *api.ReviewDecision {
	// TODO: Implement more sophisticated parsing
	decision := &api.ReviewDecision{}

	// Simple success/failure check
	if strings.Contains(strings.ToLower(aiResponse), "success") {
		decision.Success = true
		return decision
	}

	// Parse modifications if any
	lines := strings.Split(aiResponse, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Modification:") {
			mod := parseModification(strings.TrimPrefix(line, "Modification:"))
			decision.Modifications = append(decision.Modifications, mod)
		}
	}

	return decision
}
