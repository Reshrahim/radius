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

package paramresolver

import (
	"testing"

	"github.com/radius-project/radius/pkg/recipes"
	"github.com/radius-project/radius/pkg/recipes/recipecontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testContext() *recipecontext.Context {
	return &recipecontext.Context{
		Resource: recipecontext.Resource{
			ResourceInfo: recipecontext.ResourceInfo{
				Name: "my-resource",
				ID:   "/planes/radius/local/resourceGroups/test/providers/Applications.Datastores/redisCaches/my-resource",
			},
			Type: "Applications.Datastores/redisCaches",
			Properties: map[string]any{
				"size":     "s",
				"replicas": 3,
				"enabled":  true,
				"nested": map[string]any{
					"key": "nested-value",
				},
			},
		},
		Application: recipecontext.ResourceInfo{
			Name: "my-app",
			ID:   "/planes/radius/local/resourceGroups/test/providers/Applications.Core/applications/my-app",
		},
		Environment: recipecontext.ResourceInfo{
			Name: "my-env",
			ID:   "/planes/radius/local/resourceGroups/test/providers/Applications.Core/environments/my-env",
		},
		Runtime: recipes.RuntimeConfiguration{
			Kubernetes: &recipes.KubernetesRuntime{
				Namespace:            "app-ns",
				EnvironmentNamespace: "env-ns",
			},
		},
		Azure: &recipecontext.ProviderAzure{
			ResourceGroup: recipecontext.AzureResourceGroup{
				Name: "my-rg",
				ID:   "/subscriptions/sub-id/resourceGroups/my-rg",
			},
			Subscription: recipecontext.AzureSubscription{
				SubscriptionID: "sub-id",
			},
		},
		AWS: &recipecontext.ProviderAWS{
			Region:  "us-east-1",
			Account: "123456789",
		},
	}
}

func TestResolveParameters(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]any
		ctx      *recipecontext.Context
		expected map[string]any
	}{
		{
			name:     "nil params",
			params:   nil,
			ctx:      testContext(),
			expected: nil,
		},
		{
			name:     "nil context",
			params:   map[string]any{"key": "value"},
			ctx:      nil,
			expected: map[string]any{"key": "value"},
		},
		{
			name:     "empty map",
			params:   map[string]any{},
			ctx:      testContext(),
			expected: map[string]any{},
		},
		{
			name: "simple expression resolves to string type",
			params: map[string]any{
				"name": "{{context.resource.name}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"name": "my-resource",
			},
		},
		{
			name: "expression embedded in text resolves to string",
			params: map[string]any{
				"name": "prefix-{{context.resource.name}}-suffix",
			},
			ctx: testContext(),
			expected: map[string]any{
				"name": "prefix-my-resource-suffix",
			},
		},
		{
			name: "multiple expressions in one value",
			params: map[string]any{
				"tag": "{{context.application.name}}-{{context.environment.name}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"tag": "my-app-my-env",
			},
		},
		{
			name: "unrecognized expression left as-is",
			params: map[string]any{
				"val": "{{context.unknown.path}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"val": "{{context.unknown.path}}",
			},
		},
		{
			name: "non-string values pass through unchanged",
			params: map[string]any{
				"port":    8080,
				"enabled": true,
				"ratio":   3.14,
			},
			ctx: testContext(),
			expected: map[string]any{
				"port":    8080,
				"enabled": true,
				"ratio":   3.14,
			},
		},
		{
			name: "nested map traversal",
			params: map[string]any{
				"config": map[string]any{
					"namespace": "{{context.runtime.kubernetes.namespace}}",
					"static":    "unchanged",
				},
			},
			ctx: testContext(),
			expected: map[string]any{
				"config": map[string]any{
					"namespace": "app-ns",
					"static":    "unchanged",
				},
			},
		},
		{
			name: "slice traversal",
			params: map[string]any{
				"tags": []any{"{{context.resource.name}}", "static"},
			},
			ctx: testContext(),
			expected: map[string]any{
				"tags": []any{"my-resource", "static"},
			},
		},
		{
			name: "context.resource.properties resolves existing property",
			params: map[string]any{
				"size": "{{context.resource.properties.size}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"size": "s",
			},
		},
		{
			name: "context.resource.properties resolves typed int value",
			params: map[string]any{
				"replicas": "{{context.resource.properties.replicas}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"replicas": 3,
			},
		},
		{
			name: "context.resource.properties resolves typed bool value",
			params: map[string]any{
				"enabled": "{{context.resource.properties.enabled}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"enabled": true,
			},
		},
		{
			name: "context.resource.properties nested key",
			params: map[string]any{
				"val": "{{context.resource.properties.nested.key}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"val": "nested-value",
			},
		},
		{
			name: "context.resource.properties missing property left as-is",
			params: map[string]any{
				"val": "{{context.resource.properties.nonexistent}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"val": "{{context.resource.properties.nonexistent}}",
			},
		},
		{
			name: "property embedded in text uses string representation",
			params: map[string]any{
				"label": "replicas-{{context.resource.properties.replicas}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"label": "replicas-3",
			},
		},
		{
			name: "kubernetes context",
			params: map[string]any{
				"ns":    "{{context.runtime.kubernetes.namespace}}",
				"envNs": "{{context.runtime.kubernetes.environmentNamespace}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"ns":    "app-ns",
				"envNs": "env-ns",
			},
		},
		{
			name: "azure context",
			params: map[string]any{
				"rg":  "{{context.azure.resourceGroup.name}}",
				"sub": "{{context.azure.subscription.subscriptionId}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"rg":  "my-rg",
				"sub": "sub-id",
			},
		},
		{
			name: "aws context",
			params: map[string]any{
				"region":  "{{context.aws.region}}",
				"account": "{{context.aws.account}}",
			},
			ctx: testContext(),
			expected: map[string]any{
				"region":  "us-east-1",
				"account": "123456789",
			},
		},
		{
			name: "malformed expression (unclosed) left as-is",
			params: map[string]any{
				"val": "{{context.resource.name",
			},
			ctx: testContext(),
			expected: map[string]any{
				"val": "{{context.resource.name",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolveParameters(tt.params, tt.ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolveParameters_DoesNotMutateInput(t *testing.T) {
	original := map[string]any{
		"name": "{{context.resource.name}}",
		"nested": map[string]any{
			"ns": "{{context.runtime.kubernetes.namespace}}",
		},
	}

	// Deep copy for comparison
	originalCopy := map[string]any{
		"name": "{{context.resource.name}}",
		"nested": map[string]any{
			"ns": "{{context.runtime.kubernetes.namespace}}",
		},
	}

	result := ResolveParameters(original, testContext())
	require.NotNil(t, result)

	// Original must be unchanged
	assert.Equal(t, originalCopy, original)
	// Result should be different
	assert.Equal(t, "my-resource", result["name"])
}

func TestTernaryExpression(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]any
		ctx      *recipecontext.Context
		expected map[string]any
	}{
		{
			name: "ternary true branch",
			params: map[string]any{
				"sku": `{{context.resource.properties.size == "s" ? "B_Standard_B1ms" : "GP_Standard_D2s_v3"}}`,
			},
			ctx: testContext(),
			expected: map[string]any{
				"sku": "B_Standard_B1ms",
			},
		},
		{
			name: "ternary false branch",
			params: map[string]any{
				"sku": `{{context.resource.properties.size == "l" ? "B_Standard_B1ms" : "GP_Standard_D2s_v3"}}`,
			},
			ctx: testContext(),
			expected: map[string]any{
				"sku": "GP_Standard_D2s_v3",
			},
		},
		{
			name: "ternary with unresolvable condition left as-is",
			params: map[string]any{
				"sku": `{{context.resource.properties.unknown == "s" ? "small" : "large"}}`,
			},
			ctx: testContext(),
			expected: map[string]any{
				"sku": `{{context.resource.properties.unknown == "s" ? "small" : "large"}}`,
			},
		},
		{
			name: "nested/chained ternary left as-is (V1 limitation)",
			params: map[string]any{
				"sku": `{{context.resource.properties.size == "s" ? "small" : context.resource.properties.size == "m" ? "medium" : "large"}}`,
			},
			ctx: testContext(),
			expected: map[string]any{
				"sku": `{{context.resource.properties.size == "s" ? "small" : context.resource.properties.size == "m" ? "medium" : "large"}}`,
			},
		},
		{
			name: "ternary mixed with literal text",
			params: map[string]any{
				"label": `prefix-{{context.resource.properties.size == "s" ? "small" : "large"}}-suffix`,
			},
			ctx: testContext(),
			expected: map[string]any{
				"label": "prefix-small-suffix",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolveParameters(tt.params, tt.ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestShallowMergeParameters(t *testing.T) {
	tests := []struct {
		name     string
		base     map[string]any
		override map[string]any
		expected map[string]any
	}{
		{
			name:     "both nil",
			base:     nil,
			override: nil,
			expected: nil,
		},
		{
			name:     "nil base",
			base:     nil,
			override: map[string]any{"key": "value"},
			expected: map[string]any{"key": "value"},
		},
		{
			name:     "nil override",
			base:     map[string]any{"key": "value"},
			override: nil,
			expected: map[string]any{"key": "value"},
		},
		{
			name:     "disjoint keys",
			base:     map[string]any{"a": 1},
			override: map[string]any{"b": 2},
			expected: map[string]any{"a": 1, "b": 2},
		},
		{
			name:     "overlapping keys - override wins",
			base:     map[string]any{"a": 1, "b": 2},
			override: map[string]any{"b": 99},
			expected: map[string]any{"a": 1, "b": 99},
		},
		{
			name: "nested object replaced entirely (not deep-merged)",
			base: map[string]any{
				"config": map[string]any{"x": 1, "y": 2},
			},
			override: map[string]any{
				"config": map[string]any{"z": 3},
			},
			expected: map[string]any{
				"config": map[string]any{"z": 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShallowMergeParameters(tt.base, tt.override)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestShallowMergeParameters_DoesNotMutateInput(t *testing.T) {
	base := map[string]any{"a": 1, "b": 2}
	override := map[string]any{"b": 99, "c": 3}

	baseCopy := map[string]any{"a": 1, "b": 2}
	overrideCopy := map[string]any{"b": 99, "c": 3}

	result := ShallowMergeParameters(base, override)
	require.NotNil(t, result)

	assert.Equal(t, baseCopy, base)
	assert.Equal(t, overrideCopy, override)
	assert.Equal(t, 99, result["b"])
}
