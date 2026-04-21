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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/radius-project/radius/pkg/ucp/ucplog"
)

const (
	// reachabilityTimeout is the maximum time to wait for a reachability check.
	reachabilityTimeout = 30 * time.Second

	// terraformRegistryBaseURL is the base URL for the public Terraform registry API.
	terraformRegistryBaseURL = "https://registry.terraform.io/v1/modules"

	// ociScheme is the URI scheme prefix for OCI references.
	ociScheme = "oci://"

	// gitPrefix is the go-getter prefix for Git sources.
	gitPrefix = "git::"

	// s3Prefix is the go-getter prefix for S3 sources.
	s3Prefix = "s3::"

	// gcsPrefix is the go-getter prefix for GCS sources.
	gcsPrefix = "gcs::"
)

// defaultResolver implements the Resolver interface using pattern-matching
// classification and lightweight network probes for reachability.
type defaultResolver struct {
	httpClient  *http.Client
	registryURL string // overridable for testing; defaults to terraformRegistryBaseURL
}

// NewResolver creates a new Resolver with default settings.
func NewResolver() Resolver {
	return &defaultResolver{
		httpClient: &http.Client{
			Timeout: reachabilityTimeout,
		},
		registryURL: terraformRegistryBaseURL,
	}
}

// newResolverWithClient creates a Resolver with a custom HTTP client (for testing).
func newResolverWithClient(client *http.Client) Resolver {
	return &defaultResolver{
		httpClient:  client,
		registryURL: terraformRegistryBaseURL,
	}
}

// Classify determines the source type of a template path without making
// any network calls. Classification is purely based on string pattern matching.
func (r *defaultResolver) Classify(templatePath string) ResolvedSource {
	path := strings.TrimSpace(templatePath)
	if path == "" {
		return ResolvedSource{
			Type:           SourceTypeUnknown,
			OriginalPath:   templatePath,
			IsDirectModule: false,
		}
	}

	sourceType := classifyPath(path)
	isDirect := sourceType != SourceTypeUnknown && sourceType != SourceTypeOCI

	return ResolvedSource{
		Type:           sourceType,
		OriginalPath:   templatePath,
		IsDirectModule: isDirect,
	}
}

// classifyPath applies ordered pattern-matching rules to determine source type.
func classifyPath(path string) SourceType {
	// 1. Explicit go-getter prefixes take highest priority.
	if strings.HasPrefix(path, gitPrefix) {
		return SourceTypeGit
	}
	if strings.HasPrefix(path, s3Prefix) {
		return SourceTypeS3
	}
	if strings.HasPrefix(path, gcsPrefix) {
		return SourceTypeGCS
	}

	// 2. OCI scheme.
	if strings.HasPrefix(path, ociScheme) {
		return SourceTypeOCI
	}

	// 3. HTTP/HTTPS URLs (without git:: prefix).
	if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
		return SourceTypeHTTP
	}

	// 4. Terraform registry: exactly 3 slash-separated segments with no scheme.
	// Examples: "hashicorp/consul/aws", "Azure/mysql/azurerm"
	if isTerraformRegistryPath(path) {
		return SourceTypeTerraformRegistry
	}

	// 5. OCI image reference patterns (e.g., "ghcr.io/org/image:tag").
	// These contain a dot in the first segment (the registry hostname).
	if looksLikeOCIReference(path) {
		return SourceTypeOCI
	}

	return SourceTypeUnknown
}

// isTerraformRegistryPath returns true if the path matches the Terraform
// registry format: exactly 3 slash-separated segments with no scheme or dots
// in the first segment (to distinguish from OCI references like ghcr.io/...).
func isTerraformRegistryPath(path string) bool {
	// Strip query parameters and subdirectory paths for classification.
	cleanPath := path
	if idx := strings.Index(cleanPath, "?"); idx != -1 {
		cleanPath = cleanPath[:idx]
	}
	if idx := strings.Index(cleanPath, "//"); idx != -1 {
		cleanPath = cleanPath[:idx]
	}

	parts := strings.Split(cleanPath, "/")
	if len(parts) != 3 {
		return false
	}

	// All 3 segments must be non-empty.
	for _, p := range parts {
		if p == "" {
			return false
		}
	}

	// The first segment (namespace) must not contain a dot — that would
	// indicate an OCI registry hostname like "ghcr.io".
	if strings.Contains(parts[0], ".") {
		return false
	}

	return true
}

// looksLikeOCIReference returns true if the path looks like an OCI image
// reference (e.g., "ghcr.io/radius-project/recipes/mysql:latest").
func looksLikeOCIReference(path string) bool {
	// OCI references have a hostname with a dot as the first segment.
	firstSlash := strings.Index(path, "/")
	if firstSlash <= 0 {
		return false
	}
	hostname := path[:firstSlash]
	return strings.Contains(hostname, ".")
}

// ValidateReachability performs a lightweight network check to verify
// that the module source is accessible.
func (r *defaultResolver) ValidateReachability(ctx context.Context, templatePath string, templateVersion string) error {
	logger := ucplog.FromContextOrDiscard(ctx)

	resolved := r.Classify(templatePath)

	// Skip validation for unknown/OCI sources.
	if resolved.Type == SourceTypeUnknown || resolved.Type == SourceTypeOCI {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, reachabilityTimeout)
	defer cancel()

	var err error
	switch resolved.Type {
	case SourceTypeTerraformRegistry:
		err = r.validateRegistryReachability(ctx, templatePath, templateVersion)
	case SourceTypeHTTP:
		err = r.validateHTTPReachability(ctx, templatePath)
	case SourceTypeGit, SourceTypeS3, SourceTypeGCS:
		// For Git/S3/GCS, we skip proactive validation — these sources
		// will be validated at terraform init time. The lightweight
		// probes for these source types are complex (git ls-remote requires
		// git binary, S3/GCS require cloud SDK) and not worth the
		// complexity for a best-effort check.
		logger.Info("Skipping proactive reachability check for source type", "sourceType", resolved.Type, "templatePath", templatePath)
		return nil
	}

	if err != nil {
		// Check if this is a transient failure (timeout/context deadline).
		if ctx.Err() != nil {
			logger.Info("Transient failure during reachability check, allowing creation", "templatePath", templatePath, "error", err.Error())
			return nil
		}
		return err
	}

	return nil
}

// validateRegistryReachability checks if a Terraform registry module exists.
func (r *defaultResolver) validateRegistryReachability(ctx context.Context, templatePath string, templateVersion string) error {
	baseURL := r.registryURL
	url := fmt.Sprintf("%s/%s", baseURL, templatePath)
	if templateVersion != "" {
		url = fmt.Sprintf("%s/%s", url, templateVersion)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create registry request for %q: %w", templatePath, err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("terraform module %q not found in registry", templatePath)
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("access denied to terraform module %q in registry (HTTP %d)", templatePath, resp.StatusCode)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("registry returned HTTP %d for terraform module %q", resp.StatusCode, templatePath)
	}

	return nil
}

// validateHTTPReachability checks if an HTTP source is reachable via HEAD request.
func (r *defaultResolver) validateHTTPReachability(ctx context.Context, templatePath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, templatePath, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request for %q: %w", templatePath, err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("HTTP source %q not found (404)", templatePath)
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("access denied to HTTP source %q (HTTP %d)", templatePath, resp.StatusCode)
	}

	return nil
}

// IsDirectModuleSource is a convenience function that classifies the given
// templatePath and returns true if it represents a direct Terraform module
// source (registry, git, HTTP, S3, or GCS) rather than a wrapped/OCI recipe.
//
// This is the primary entry point for the terraform driver to determine
// which output mapping strategy to use.
func IsDirectModuleSource(templatePath string) bool {
	r := &defaultResolver{}
	return r.Classify(templatePath).IsDirectModule
}
