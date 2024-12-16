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
	"time"

	"github.com/GoogleCloudPlatform/k8s-config-connector/dev/tools/aiframework/api"
)

type Framework struct {
	validators []api.Validator
	tools      api.Tools
	ai         api.AIClient
}

func (f *Framework) ProcessAPIChange(change api.APIChange) error {
	ctx := &api.UpdateContext{
		APIChange:  change,
		MaxRetries: 3,
		State: api.UpdateState{
			Phase: api.PhaseAnalysis,
		},
	}

	for ctx.State.CurrentAttempt < ctx.MaxRetries {
		if err := f.runPhase(ctx); err != nil {
			attempt := api.AttemptResult{
				Phase:     ctx.State.Phase,
				Error:     err,
				Timestamp: time.Now(),
			}

			if ctx.State.Phase == api.PhaseValidation {
				attempt.ValidationResults = ctx.State.LastValidationResults
			}
			if ctx.State.Phase == api.PhaseImplementation {
				attempt.ModificationsAttempted = ctx.State.PlannedModifications
			}

			ctx.History = append(ctx.History, attempt)
			ctx.State.CurrentAttempt++
			continue
		}

		if ctx.State.Phase == api.PhasePRCreation {
			ctx.History = append(ctx.History, api.AttemptResult{
				Phase:     api.PhasePRCreation,
				Timestamp: time.Now(),
			})
			return nil
		}

		ctx.State.Phase = nextPhase(ctx.State.Phase)
	}

	return fmt.Errorf("max retries exceeded after %d attempts", len(ctx.History))
}

func (f *Framework) runPhase(ctx *api.UpdateContext) error {
	switch ctx.State.Phase {
	case api.PhaseAnalysis:
		return f.analyzeChange(ctx)
	case api.PhaseImplementation:
		return f.implementChange(ctx)
	case api.PhaseValidation:
		return f.validateChange(ctx)
	case api.PhasePRCreation:
		return f.createPR(ctx)
	}
	return fmt.Errorf("unknown phase: %s", ctx.State.Phase)
}

func nextPhase(current api.UpdatePhase) api.UpdatePhase {
	switch current {
	case api.PhaseAnalysis:
		return api.PhaseImplementation
	case api.PhaseImplementation:
		return api.PhaseValidation
	case api.PhaseValidation:
		return api.PhasePRCreation
	default:
		return api.PhaseAnalysis
	}
}
