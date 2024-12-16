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

package api

import "time"

type UpdateContext struct {
	APIChange  APIChange
	State      UpdateState
	History    []AttemptResult
	MaxRetries int
}

type APIChange struct {
	MessageName string
	FieldName   string
}

type UpdateState struct {
	Phase                 UpdatePhase
	CurrentAttempt        int
	LastError             error
	Analysis              *AnalysisResult
	ModifiedFiles         []string
	PlannedModifications  []Modification
	LastValidationResults []*ValidationResult
}

type UpdatePhase string

const (
	PhaseAnalysis       UpdatePhase = "Analysis"
	PhaseImplementation UpdatePhase = "Implementation"
	PhaseValidation     UpdatePhase = "Validation"
	PhasePRCreation     UpdatePhase = "PRCreation"
)

type AttemptResult struct {
	Phase                  UpdatePhase
	Error                  error
	ValidationResults      []*ValidationResult
	ModificationsAttempted []Modification
	Timestamp              time.Time
}

type AnalysisResult struct {
	AffectedFiles   []string
	RequiredChanges []Change
}

type Change struct {
	Type        ChangeType
	Description string
	FilePath    string
}

type ChangeType string

const (
	ChangeTypeKRM          ChangeType = "KRM"
	ChangeTypeMapper       ChangeType = "Mapper"
	ChangeTypeModification ChangeType = "Modification"
)

type ImplementationPlan struct {
	Steps []ImplementationStep
}

type ImplementationStep struct {
	Type          StepType
	FilePath      string
	TypesOptions  *GenerateTypesOptions
	MapperOptions *GenerateMapperOptions
	Modifications []Modification
}

type StepType string

const (
	StepGenerateKRM    StepType = "GenerateKRM"
	StepGenerateMapper StepType = "GenerateMapper"
	StepModifyFile     StepType = "ModifyFile"
)

type Modification struct {
	StartLine  int
	EndLine    int
	NewContent string
}

type ValidationResult struct {
	Success bool
	Errors  []CompileError
	Logs    string
}

type CompileError struct {
	File    string
	Line    string
	Column  string
	Message string
}

type ReviewDecision struct {
	Success       bool
	Modifications []Modification
}

type GenerateTypesOptions struct {
	Service    string
	Resource   string
	APIVersion string
}

type GenerateMapperOptions struct {
	Service    string
	APIVersion string
	OutputDir  string
}
