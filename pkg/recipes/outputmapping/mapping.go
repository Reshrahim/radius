/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package outputmapping provides utilities for mapping module outputs to
// resource type read-only properties. It is shared by both Terraform and Bicep
// drivers to apply the RecipePack `outputs` field consistently.
package outputmapping

// OutputValue represents a single output from a module execution, carrying
// both the value and its sensitivity classification.
type OutputValue struct {
	// Value is the output value from the module.
	Value any
	// Sensitive indicates whether this output was marked as sensitive by the module.
	Sensitive bool
}

// Apply maps raw module outputs to resource type read-only properties using the
// provided outputs mapping. The outputsMap keys are resource property names and
// values are module output names.
//
// When outputsMap is nil or empty, all outputs pass through unchanged (values
// go to the values map, sensitive outputs go to the secrets map).
//
// When outputsMap is provided, only explicitly mapped outputs are surfaced
// (acts as an allow-list with key renaming).
//
// Returns two new maps (values and secrets) — inputs are never mutated.
func Apply(rawOutputs map[string]OutputValue, outputsMap map[string]string) (values map[string]any, secrets map[string]any) {
	values = make(map[string]any)
	secrets = make(map[string]any)

	if len(rawOutputs) == 0 {
		return values, secrets
	}

	if len(outputsMap) == 0 {
		// Pass through all outputs
		for name, output := range rawOutputs {
			if output.Sensitive {
				secrets[name] = output.Value
			} else {
				values[name] = output.Value
			}
		}
		return values, secrets
	}

	// For each entry in the outputs map, check if the value references a module
	// output. If it does, use the module output's value. Otherwise, treat the
	// value as a literal (direct value pass-through).
	for resourceProp, moduleOutputOrLiteral := range outputsMap {
		if output, ok := rawOutputs[moduleOutputOrLiteral]; ok {
			// Value matches a module output name — use the output value.
			if output.Sensitive {
				secrets[resourceProp] = output.Value
			} else {
				values[resourceProp] = output.Value
			}
		} else {
			// No matching module output — treat as a literal value.
			values[resourceProp] = moduleOutputOrLiteral
		}
	}

	return values, secrets
}
