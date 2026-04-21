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

// Package source provides module source classification and validation for
// Terraform recipe template paths. It determines whether a templatePath
// refers to a direct Terraform module source (registry, Git, HTTP) or an
// existing OCI/wrapped recipe path, enabling the recipe system to apply
// appropriate execution and output mapping strategies.
package source

import "context"

// SourceType classifies the format of a Terraform module source path.
type SourceType int

const (
	// SourceTypeUnknown indicates the source format could not be classified.
	// The system should fall back to existing OCI/wrapped recipe resolution.
	SourceTypeUnknown SourceType = iota

	// SourceTypeTerraformRegistry indicates a standard Terraform registry source.
	// Format: "namespace/name/provider" (exactly 3 slash-separated segments, no scheme).
	// Example: "hashicorp/consul/aws", "Azure/cosmosdb/azurerm"
	SourceTypeTerraformRegistry

	// SourceTypeGit indicates a Git-hosted module source.
	// Format: "git::https://..." or "git::ssh://..."
	// Supports ref specifiers (?ref=v1.0.0) and subdirectories (//modules/vpc).
	SourceTypeGit

	// SourceTypeHTTP indicates an HTTP/HTTPS archive source.
	// Format: "https://example.com/module.tar.gz" (without git:: prefix)
	SourceTypeHTTP

	// SourceTypeS3 indicates an S3-hosted module source.
	// Format: "s3::https://..." or "s3::bucket-name/key"
	SourceTypeS3

	// SourceTypeGCS indicates a GCS-hosted module source.
	// Format: "gcs::https://..." or "gcs::bucket-name/key"
	SourceTypeGCS

	// SourceTypeOCI indicates an OCI registry source (existing wrapped recipe path).
	// Format: contains "oci://" or matches OCI image reference patterns.
	SourceTypeOCI
)

// String returns the human-readable name of a SourceType.
func (s SourceType) String() string {
	switch s {
	case SourceTypeTerraformRegistry:
		return "TerraformRegistry"
	case SourceTypeGit:
		return "Git"
	case SourceTypeHTTP:
		return "HTTP"
	case SourceTypeS3:
		return "S3"
	case SourceTypeGCS:
		return "GCS"
	case SourceTypeOCI:
		return "OCI"
	default:
		return "Unknown"
	}
}

// ResolvedSource contains the classification result for a template path.
type ResolvedSource struct {
	// Type is the classified source type.
	Type SourceType

	// OriginalPath is the unmodified templatePath value.
	OriginalPath string

	// IsDirectModule is true when the source is a direct Terraform module
	// (not a wrapped/OCI recipe). This determines output mapping strategy.
	IsDirectModule bool
}

// Resolver classifies and validates Terraform module source paths.
type Resolver interface {
	// Classify determines the source type of a template path without making
	// any network calls. Classification is purely based on string pattern matching.
	//
	// Returns a ResolvedSource with Type set to the detected source type.
	// If the format is not recognized, Type is SourceTypeUnknown and
	// IsDirectModule is false (indicating fallback to existing behavior).
	Classify(templatePath string) ResolvedSource

	// ValidateReachability performs a lightweight network check to verify
	// that the module source is accessible. This is called at RecipePack
	// creation time per FR-014.
	//
	// For registry modules: HTTP GET to registry API
	// For Git sources: git ls-remote
	// For HTTP sources: HTTP HEAD request
	//
	// Returns nil if the source is reachable, or an error describing why
	// it could not be reached. The check has a 30-second timeout.
	// Transient failures (timeouts) return nil with a logged warning.
	// Definitive failures (404, auth denied) return an error.
	//
	// If the source type is SourceTypeUnknown or SourceTypeOCI, this
	// method returns nil (no validation for fallback paths).
	ValidateReachability(ctx context.Context, templatePath string, templateVersion string) error
}

// ModuleOutputInfo describes a single output declared by a Terraform module.
type ModuleOutputInfo struct {
	// Name is the output variable name as declared in the module.
	Name string

	// Sensitive is true if the output is marked sensitive = true in the module.
	Sensitive bool
}
