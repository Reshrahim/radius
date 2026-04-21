
# Implementation Plan: Direct Terraform Module as Recipe Template Path

**Branch**: `001-terraform-module-registration` | **Date**: 2025-07-18 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-terraform-module-registration/spec.md`

## Summary

Enable platform engineers to use existing Terraform modules directly as the `templatePath` for a recipe in a RecipePack — without wrapping, republishing, or creating Radius-specific artifacts. The system resolves standard Terraform module sources (registry `namespace/name/provider`, Git `git::https://...`, HTTP URLs) at deployment time, passes recipe parameters as Terraform input variables, and surfaces all module outputs flat into the `Values` map. Resolution uses a "try direct first, fall back to existing OCI/wrapped path" strategy. Validation at RecipePack creation time checks module source reachability.

## Technical Context

**Language/Version**: Go (version per `go.mod`, currently 1.21+)  
**Primary Dependencies**: `hashicorp/terraform-exec`, `hashicorp/go-getter`, `hashicorp/terraform-config-inspect` (`tfconfig`), `github.com/hashicorp/hc-install`  
**Storage**: Kubernetes secrets (Terraform state backend), existing data store for RecipePack resources  
**Testing**: `go test` with `go.uber.org/mock/gomock`, table-driven tests, functional tests via `magpiego`  
**Target Platform**: Linux containers (Kubernetes), development on macOS/Linux  
**Project Type**: Existing Go monorepo — changes span multiple packages  
**Performance Goals**: Deployment within 10% of equivalent wrapped-recipe deployments; fail within 60s on unreachable source  
**Constraints**: No caching between deployments; no new data model fields; backward-compatible with existing recipes  
**Scale/Scope**: ~6 packages modified, ~3 new files, ~400-600 lines of new production code

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. API-First Design | ✅ PASS | No new API endpoints; operates within existing RecipePack schema |
| II. Idiomatic Code Standards | ✅ PASS | Go patterns, explicit error handling, godoc on all exported items |
| III. Multi-Cloud Neutrality | ✅ PASS | Terraform modules are inherently multi-cloud; no cloud-specific logic added |
| IV. Testing Pyramid | ✅ PASS | Unit tests for source detection/output mapping, integration tests for module download, functional tests for full deploy |
| V. Collaboration-Centric | ✅ PASS | Reduces platform engineer toil; developers access outputs via existing Radius interface |
| VII. Simplicity Over Cleverness | ✅ PASS | Leverages existing go-getter/terraform-exec; minimal new abstraction |
| VIII. Separation of Concerns | ✅ PASS | Source resolution isolated in new module; output mapping in driver layer |
| IX. Incremental Adoption | ✅ PASS | Existing recipes unaffected; try-direct-then-fallback strategy ensures backward compat |
| XII. Resource Type Schema Quality | N/A | No new resource types |
| XIII. Recipe Development Standards | ✅ PASS | Extends recipe system without changing conventions |

**Gate Result**: ✅ PASS — No violations. Proceed to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/001-terraform-module-registration/
├── plan.md              # This file
├── research.md          # Phase 0: Research findings
├── data-model.md        # Phase 1: Data model documentation
├── quickstart.md        # Phase 1: Developer quickstart
├── contracts/           # Phase 1: Internal contracts
│   └── source-resolver.go  # Source resolution interface contract
└── tasks.md             # Phase 2: Implementation tasks (generated separately)
```

### Source Code (repository root)

```text
pkg/
├── recipes/
│   ├── source/                          # NEW: Module source resolution
│   │   ├── resolver.go                  # Source format detection & resolution
│   │   ├── resolver_test.go             # Unit tests
│   │   └── types.go                     # SourceType enum, ResolvedSource struct
│   ├── driver/terraform/
│   │   ├── terraform.go                 # MODIFIED: prepareRecipeResponse for direct modules
│   │   └── terraform_test.go            # MODIFIED: New test cases
│   ├── terraform/
│   │   ├── config/
│   │   │   ├── config.go               # MODIFIED: AddOutputs for direct module mode
│   │   │   └── config_test.go          # MODIFIED: New test cases
│   │   ├── execute.go                   # MODIFIED: Direct module execution path
│   │   ├── execute_test.go             # MODIFIED: New test cases
│   │   └── module.go                    # MODIFIED: Source-aware download
│   └── types.go                         # Reference only (RecipeOutput, EnvironmentDefinition)
├── corerp/
│   └── frontend/controller/recipepacks/
│       ├── createorupdaterecipepack.go  # MODIFIED: Add source reachability validation
│       └── createorupdaterecipepack_test.go # MODIFIED: Validation tests
└── ...

test/
└── functional/
    └── shared/resources/
        └── testdata/recipes/
            └── direct-module/           # NEW: Functional test recipes
```

**Structure Decision**: Existing monorepo Go project. New `pkg/recipes/source/` package isolates module source resolution logic. All other changes are modifications to existing packages, following established patterns.

## Complexity Tracking

> No violations to justify — design adheres to all constitution principles.

## Post-Design Constitution Re-Check (after Phase 1)

| Principle | Status | Notes |
|-----------|--------|-------|
| I. API-First Design | ✅ PASS | No new APIs introduced. Existing RecipePack API schema unchanged. |
| II. Idiomatic Code Standards | ✅ PASS | New `source` package follows Go conventions: exported interface with godoc, package-level convenience function, typed constants. |
| III. Multi-Cloud Neutrality | ✅ PASS | Source resolver is cloud-agnostic. Module execution uses existing multi-cloud provider framework. |
| IV. Testing Pyramid | ✅ PASS | Design includes: unit tests for Classify/ValidateReachability, unit tests for output mapping logic, integration tests for module download, functional tests for end-to-end deploy. |
| V. Collaboration-Centric | ✅ PASS | Platform engineers get immediate access to entire TF ecosystem. Developers access outputs through standard Radius resource properties. |
| VII. Simplicity Over Cleverness | ✅ PASS | Single new package (`source/`) with ~3 files. No complex abstractions. Pattern-matching classification is deterministic and testable. |
| VIII. Separation of Concerns | ✅ PASS | Source classification (new package) separated from execution (existing executor) and output mapping (existing driver). Each responsibility in one place. |
| IX. Incremental Adoption | ✅ PASS | Fully backward compatible. Unknown sources fall back to existing behavior. No feature flags needed — direct detection is transparent. |
| XIII. Recipe Development Standards | ✅ PASS | Does not change how recipes are authored. Simply expands what sources are accepted. |

**Post-Design Gate Result**: ✅ PASS — Design is consistent with constitution. No new violations introduced.

## Architecture Decision Summary

| Decision | Approach | Reference |
|----------|----------|-----------|
| Source detection | Pattern-matching heuristics, no network calls | research.md R-001 |
| Fallback strategy | Classify first; unambiguous sources skip fallback, ambiguous try-then-fallback | research.md R-002 |
| Output mapping | Direct modules: flat Values/Secrets; Wrapped: existing `result` parsing | research.md R-003 |
| Creation validation | Lightweight reachability check (registry API, git ls-remote, HTTP HEAD) | research.md R-004 |
| Version pinning | Existing `TemplateVersion` field → `version` in module block | research.md R-005 |
| Mode detection | Runtime re-classification of TemplatePath (no new persisted fields) | research.md R-006 |
| Config generation | All-output forwarding for direct modules vs. single `result` for wrapped | research.md R-007 |
| Metadata retrieval | No changes — existing `downloadAndInspect` works for all sources | research.md R-008 |
| Error handling | Existing `RecipeError` types with Terraform error details | research.md R-009 |
| Sensitive outputs | Sensitive → Secrets map; non-sensitive → Values map | research.md R-010 |
