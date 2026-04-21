# Data Model: Direct Terraform Module as Recipe Template Path

## Overview

This feature operates entirely within the existing data model. No new entities, fields, or resource types are introduced (per FR-012). The changes are behavioral — expanding what `templatePath` values are accepted and how outputs are mapped.

## Existing Entities (Behavior Extended)

### RecipeDefinition (unchanged schema)

**Location**: `pkg/corerp/datamodel/recipepack.go`

```go
type RecipeDefinition struct {
    RecipeKind     string         `json:"recipeKind"`     // "terraform" or "bicep"
    RecipeLocation string         `json:"recipeLocation"` // Template source path/URL
    Parameters     map[string]any `json:"parameters,omitempty"`
    PlainHTTP      bool           `json:"plainHTTP,omitempty"`
}
```

**Behavioral Change**: `RecipeLocation` now additionally accepts:
- Terraform Registry: `hashicorp/consul/aws`, `Azure/cosmosdb/azurerm`
- Git URLs: `git::https://github.com/org/terraform-aws-vpc.git`
- Git with ref: `git::https://github.com/org/module.git?ref=v2.0.0`
- Git with subdirectory: `git::https://github.com/org/repo.git//modules/vpc`
- HTTP archives: `https://example.com/modules/vpc.tar.gz`

**Validation Rules** (new, at creation time):
- Source must be reachable (lightweight probe, 30s timeout)
- Source format must be classifiable by the resolver
- If unclassifiable, accepted without validation (fallback to existing behavior)

### EnvironmentDefinition (unchanged schema)

**Location**: `pkg/recipes/types.go`

```go
type EnvironmentDefinition struct {
    Name            string         // Recipe name
    Driver          string         // "terraform" or "bicep"
    ResourceType    string         // Portable resource type
    Parameters      map[string]any // Default recipe parameters
    TemplatePath    string         // Module source URL/path (expanded behavior)
    TemplateVersion string         // Module version (used for registry pinning)
    PlainHTTP       bool           // Allow insecure connections
}
```

**Behavioral Change**: Same as RecipeDefinition — `TemplatePath` accepts broader source formats.

### RecipeOutput (unchanged schema, different population logic)

**Location**: `pkg/recipes/types.go`

```go
type RecipeOutput struct {
    Resources []string       // Deployed resource IDs (from TF state)
    Secrets   map[string]any // Sensitive output values
    Values    map[string]any // Non-sensitive output values
    Status    *rpv1.RecipeStatus
}
```

**Behavioral Change for Direct Modules**:
- `Values`: Populated with ALL non-sensitive Terraform module outputs (flat key→value)
- `Secrets`: Populated with ALL sensitive Terraform module outputs
- `Resources`: Populated from Terraform state resource inspection (existing behavior)
- **No `result` output parsing** — outputs are taken directly from the module's declared outputs

**Behavioral Change for Wrapped Recipes** (unchanged):
- Existing logic: looks for `result` output, parses into Resources/Secrets/Values

## New Internal Types (not persisted)

### SourceType Enum

**Location**: `pkg/recipes/source/types.go` (new file)

```go
type SourceType int

const (
    SourceTypeUnknown            SourceType = iota // Unclassified — use fallback
    SourceTypeTerraformRegistry                    // e.g., "hashicorp/consul/aws"
    SourceTypeGit                                  // e.g., "git::https://..."
    SourceTypeHTTP                                 // e.g., "https://example.com/module.tar.gz"
    SourceTypeS3                                   // e.g., "s3::bucket/key"
    SourceTypeGCS                                  // e.g., "gcs::bucket/key"
    SourceTypeOCI                                  // Existing OCI/wrapped recipe path
)
```

### ResolvedSource

**Location**: `pkg/recipes/source/types.go` (new file)

```go
type ResolvedSource struct {
    Type          SourceType // Classified source type
    OriginalPath  string     // Original templatePath value
    IsDirectModule bool      // True if this is a direct TF module (not wrapped)
}
```

### ModuleOutputInfo

**Location**: `pkg/recipes/source/types.go` (new file)

```go
// ModuleOutputInfo captures output metadata from module inspection
type ModuleOutputInfo struct {
    Name        string // Output name
    Description string // Output description (from module)
    Sensitive   bool   // Whether output is marked sensitive
}
```

## State Transitions

### Recipe Deployment Lifecycle (with direct module)

```
┌─────────────────┐
│  RecipePack     │ ← templatePath validated at creation
│  Created        │
└────────┬────────┘
         │ Deploy resource using recipe
         ▼
┌─────────────────┐
│  Source          │ ← Classify templatePath
│  Classification │
└────────┬────────┘
         │ Direct module detected
         ▼
┌─────────────────┐
│  Module          │ ← terraform get (fresh download, no cache)
│  Download        │
└────────┬────────┘
         │ Success
         ▼
┌─────────────────┐
│  Module          │ ← Extract variables, outputs, providers
│  Inspection     │
└────────┬────────┘
         │ Generate config with all-output forwarding
         ▼
┌─────────────────┐
│  Terraform       │ ← init + apply
│  Execution      │
└────────┬────────┘
         │ Success
         ▼
┌─────────────────┐
│  Output          │ ← All outputs → Values/Secrets (flat)
│  Mapping        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Resource        │ ← Outputs accessible via Radius API
│  Deployed       │
└─────────────────┘
```

### Delete Lifecycle (unchanged flow)

```
Resource Delete Request
         │
         ▼
┌─────────────────┐
│  Source          │ ← Same classification as deploy
│  Classification │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Terraform       │ ← init + destroy (using stored state)
│  Destroy        │
└────────┬────────┘
         │ Success
         ▼
┌─────────────────┐
│  State Cleanup  │ ← Delete K8s secret with TF state
└─────────────────┘
```

## Relationships

```
RecipePack (1) ──contains──▶ (N) RecipeDefinition
     │                              │
     │                         RecipeLocation (templatePath)
     │                              │
     │                    ┌─────────┴──────────┐
     │                    │                    │
     │              Direct Module         Wrapped/OCI Recipe
     │              (new behavior)        (existing behavior)
     │                    │                    │
     │              ┌─────┴─────┐              │
     │              │           │              │
     │         Registry      Git/HTTP          │
     │              │           │              │
     │              └─────┬─────┘              │
     │                    │                    │
     ▼                    ▼                    ▼
Environment ──uses──▶ TerraformDriver ◀──uses── Environment
                          │
                    ┌─────┴─────┐
                    │           │
              Direct Mode   Wrapped Mode
              (flat output) (result output)
```

## Validation Rules

| Field | Rule | When Applied |
|-------|------|--------------|
| `RecipeLocation` | Must be non-empty string | Always (existing) |
| `RecipeLocation` | If classifiable as direct module, source must be reachable | RecipePack create/update |
| `RecipeLocation` | Format must match one of: registry, git, http, s3, gcs, or OCI | Soft validation (unknown = fallback) |
| `RecipeKind` | Must be "terraform" for direct module sources | RecipePack create/update |
| `TemplateVersion` | If specified for registry module, version must exist | RecipePack create/update |
| `Parameters` | Keys must match module input variable names | At terraform apply time (Terraform validates) |
