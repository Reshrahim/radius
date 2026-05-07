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

package source

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClassify(t *testing.T) {
	resolver := NewResolver()

	tests := []struct {
		name           string
		templatePath   string
		expectedType   SourceType
		expectedDirect bool
	}{
		// Terraform Registry sources
		{
			name:           "registry: hashicorp/consul/aws",
			templatePath:   "hashicorp/consul/aws",
			expectedType:   SourceTypeTerraformRegistry,
			expectedDirect: true,
		},
		{
			name:           "registry: Azure/mysql/azurerm",
			templatePath:   "Azure/mysql/azurerm",
			expectedType:   SourceTypeTerraformRegistry,
			expectedDirect: true,
		},
		{
			name:           "registry: terraform-aws-modules/vpc/aws",
			templatePath:   "terraform-aws-modules/vpc/aws",
			expectedType:   SourceTypeTerraformRegistry,
			expectedDirect: true,
		},

		// Git sources
		{
			name:           "git: https URL",
			templatePath:   "git::https://github.com/org/repo.git",
			expectedType:   SourceTypeGit,
			expectedDirect: true,
		},
		{
			name:           "git: with ref and subdirectory",
			templatePath:   "git::https://github.com/org/repo.git//modules/vpc?ref=v1.0.0",
			expectedType:   SourceTypeGit,
			expectedDirect: true,
		},
		{
			name:           "git: ssh URL",
			templatePath:   "git::ssh://git@github.com/org/repo.git",
			expectedType:   SourceTypeGit,
			expectedDirect: true,
		},

		// HTTP sources
		{
			name:           "http: https archive URL",
			templatePath:   "https://example.com/module.tar.gz",
			expectedType:   SourceTypeHTTP,
			expectedDirect: true,
		},
		{
			name:           "http: http URL",
			templatePath:   "http://example.com/modules/vpc.zip",
			expectedType:   SourceTypeHTTP,
			expectedDirect: true,
		},

		// S3 sources
		{
			name:           "s3: bucket source",
			templatePath:   "s3::https://s3-eu-west-1.amazonaws.com/bucket/module.zip",
			expectedType:   SourceTypeS3,
			expectedDirect: true,
		},

		// GCS sources
		{
			name:           "gcs: bucket source",
			templatePath:   "gcs::https://www.googleapis.com/storage/v1/modules/module.zip",
			expectedType:   SourceTypeGCS,
			expectedDirect: true,
		},

		// OCI sources (existing wrapped recipes)
		{
			name:           "oci: explicit oci:// scheme",
			templatePath:   "oci://registry.example.com/module",
			expectedType:   SourceTypeOCI,
			expectedDirect: false,
		},
		{
			name:           "oci: ghcr.io image reference",
			templatePath:   "ghcr.io/radius-project/recipes/mysql:latest",
			expectedType:   SourceTypeOCI,
			expectedDirect: false,
		},
		{
			name:           "oci: docker.io image reference",
			templatePath:   "docker.io/library/nginx:latest",
			expectedType:   SourceTypeOCI,
			expectedDirect: false,
		},
		{
			name:           "oci: ACR image reference",
			templatePath:   "myregistry.azurecr.io/recipes/terraform/mysql:1.0",
			expectedType:   SourceTypeOCI,
			expectedDirect: false,
		},

		// Unknown / edge cases
		{
			name:           "unknown: empty string",
			templatePath:   "",
			expectedType:   SourceTypeUnknown,
			expectedDirect: false,
		},
		{
			name:           "unknown: whitespace only",
			templatePath:   "   ",
			expectedType:   SourceTypeUnknown,
			expectedDirect: false,
		},
		{
			name:           "unknown: two segments",
			templatePath:   "two/segments",
			expectedType:   SourceTypeUnknown,
			expectedDirect: false,
		},
		{
			name:           "unknown: four segments",
			templatePath:   "four/segments/too/many",
			expectedType:   SourceTypeUnknown,
			expectedDirect: false,
		},
		{
			name:           "unknown: single word",
			templatePath:   "module",
			expectedType:   SourceTypeUnknown,
			expectedDirect: false,
		},
		{
			name:           "unknown: empty segments",
			templatePath:   "a//b",
			expectedType:   SourceTypeUnknown,
			expectedDirect: false,
		},
		{
			name:           "unknown: leading slash",
			templatePath:   "/a/b/c",
			expectedType:   SourceTypeUnknown,
			expectedDirect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.Classify(tt.templatePath)
			require.Equal(t, tt.expectedType, result.Type, "source type mismatch")
			require.Equal(t, tt.expectedDirect, result.IsDirectModule, "IsDirectModule mismatch")
			require.Equal(t, tt.templatePath, result.OriginalPath, "OriginalPath should be unmodified")
		})
	}
}

func TestIsDirectModuleSource(t *testing.T) {
	tests := []struct {
		name         string
		templatePath string
		expected     bool
	}{
		{"registry module", "hashicorp/consul/aws", true},
		{"git module", "git::https://github.com/org/repo.git", true},
		{"http module", "https://example.com/module.tar.gz", true},
		{"s3 module", "s3::https://bucket/module.zip", true},
		{"gcs module", "gcs::https://bucket/module.zip", true},
		{"oci reference", "ghcr.io/radius-project/recipes/mysql:latest", false},
		{"oci scheme", "oci://registry.example.com/module", false},
		{"unknown", "some-random-string", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, IsDirectModuleSource(tt.templatePath))
		})
	}
}

func TestSourceTypeString(t *testing.T) {
	tests := []struct {
		sourceType SourceType
		expected   string
	}{
		{SourceTypeUnknown, "Unknown"},
		{SourceTypeTerraformRegistry, "TerraformRegistry"},
		{SourceTypeGit, "Git"},
		{SourceTypeHTTP, "HTTP"},
		{SourceTypeS3, "S3"},
		{SourceTypeGCS, "GCS"},
		{SourceTypeOCI, "OCI"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.sourceType.String())
		})
	}
}

func TestValidateReachability_Registry(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		expectError bool
		errContains string
	}{
		{
			name:        "registry module exists",
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "registry module not found",
			statusCode:  http.StatusNotFound,
			expectError: true,
			errContains: "not found in registry",
		},
		{
			name:        "registry access denied",
			statusCode:  http.StatusForbidden,
			expectError: true,
			errContains: "access denied",
		},
		{
			name:        "registry server error",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
			errContains: "registry returned HTTP 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			resolver := &defaultResolver{
				httpClient:  server.Client(),
				registryURL: server.URL,
			}

			ctx := context.Background()
			err := resolver.validateRegistryReachability(ctx, "test/module/aws")

			if tt.expectError {
				require.Error(t, err)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateReachability_HTTP(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		expectError bool
		errContains string
	}{
		{
			name:        "HTTP source reachable",
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "HTTP source not found",
			statusCode:  http.StatusNotFound,
			expectError: true,
			errContains: "not found (404)",
		},
		{
			name:        "HTTP source access denied",
			statusCode:  http.StatusUnauthorized,
			expectError: true,
			errContains: "access denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			resolver := newResolverWithClient(server.Client())
			ctx := context.Background()

			err := resolver.ValidateReachability(ctx, server.URL+"/module.tar.gz")

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateReachability_SkipForUnknownAndOCI(t *testing.T) {
	resolver := NewResolver()
	ctx := context.Background()

	// Unknown source — no validation
	err := resolver.ValidateReachability(ctx, "some-random-string")
	require.NoError(t, err)

	// OCI source — no validation
	err = resolver.ValidateReachability(ctx, "ghcr.io/radius-project/recipes/mysql:latest")
	require.NoError(t, err)

	// OCI scheme — no validation
	err = resolver.ValidateReachability(ctx, "oci://registry.example.com/module")
	require.NoError(t, err)
}

func TestValidateReachability_SkipForGitS3GCS(t *testing.T) {
	resolver := NewResolver()
	ctx := context.Background()

	// Git source — skip proactive check
	err := resolver.ValidateReachability(ctx, "git::https://github.com/org/repo.git")
	require.NoError(t, err)

	// S3 source — skip proactive check
	err = resolver.ValidateReachability(ctx, "s3::https://bucket/module.zip")
	require.NoError(t, err)

	// GCS source — skip proactive check
	err = resolver.ValidateReachability(ctx, "gcs::https://bucket/module.zip")
	require.NoError(t, err)
}
