# Tasks: Direct Terraform Module as Recipe Template Path

**Spec**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md) | **Branch**: `001-terraform-module-registration`

## Task Dependency Graph

```
T1 (source types)
 └─► T2 (resolver impl) ─► T3 (resolver tests)
                            └─► T4 (config generation)
                                 └─► T5 (executor integration)
                                      └─► T6 (driver output mapping)
                                           └─► T7 (module download)
                                                └─► T8 (RecipePack validation)
                                                     └─► T9 (functional tests)
```

---

## T1: Create source package types

**File**: `pkg/recipes/source/types.go` (NEW)
**Depends on**: —
**Implements**: FR-001, R-001, R-006
**Description**: Create the `pkg/recipes/source/` package with type definitions from the contract (`specs/001-terraform-module-registration/contracts/source-resolver.go`). Define `SourceType` enum (Unknown, TerraformRegistry, Git, HTTP, S3, GCS, OCI), `ResolvedSource` struct, and `Resolver` interface. Include godoc on all exported items.

**Acceptance**:
- Package compiles
- All types match the contract exactly
- `go vet ./pkg/recipes/source/...` passes

---

## T2: Implement source resolver

**File**: `pkg/recipes/source/resolver.go` (NEW)
**Depends on**: T1
**Implements**: FR-001, FR-015, R-001, R-002
**Description**: Implement the `Resolver` interface with pattern-matching `Classify()` (no network calls) and `ValidateReachability()` (lightweight probes with 30s timeout). Also implement the package-level `IsDirectModuleSource()` convenience function.

Classification rules:
- **TerraformRegistry**: exactly 3 slash-separated segments, no `://` scheme (e.g., `Azure/mysql/azurerm`)
- **Git**: `git::` prefix
- **HTTP**: `https://` or `http://` without `git::` prefix, not matching OCI patterns
- **S3**: `s3::` prefix
- **GCS**: `gcs::` prefix
- **OCI**: contains `oci://` or matches OCI image reference patterns (e.g., `ghcr.io/...`)
- **Unknown**: anything else → falls back to existing behavior

Reachability checks (FR-014, R-004):
- Registry: HTTP GET to `https://registry.terraform.io/v1/modules/{namespace}/{name}/{provider}`
- Git: `git ls-remote` (via exec or go-getter probe)
- HTTP: HTTP HEAD request
- S3/GCS: HEAD request to bucket URL
- Unknown/OCI: return nil (no validation)
- 30-second timeout; definitive failures (404, auth denied) return error; transient failures (timeout) return nil with logged warning

**Acceptance**:
- `Classify()` is purely deterministic, no network
- `ValidateReachability()` respects 30s timeout
- `IsDirectModuleSource()` delegates to default resolver

---

## T3: Source resolver unit tests

**File**: `pkg/recipes/source/resolver_test.go` (NEW)
**Depends on**: T2
**Implements**: Constitution IV (Testing Pyramid)
**Description**: Table-driven unit tests for `Classify()` covering all source types and edge cases. Tests for `ValidateReachability()` using HTTP test server for registry/HTTP probes. Tests for `IsDirectModuleSource()` convenience function.

Test cases for `Classify()`:
- `hashicorp/consul/aws` → TerraformRegistry
- `Azure/mysql/azurerm` → TerraformRegistry
- `git::https://github.com/org/repo.git` → Git
- `git::https://github.com/org/repo.git//modules/vpc?ref=v1.0` → Git
- `https://example.com/module.tar.gz` → HTTP
- `s3::bucket/key` → S3
- `gcs::bucket/key` → GCS
- `ghcr.io/radius-project/recipes/mysql:latest` → OCI
- `oci://registry.example.com/module` → OCI
- `some-random-string` → Unknown
- Empty string → Unknown
- `two/segments` → Unknown (not 3 segments)
- `four/segments/too/many` → Unknown

Test cases for `ValidateReachability()`:
- Registry module with mock 200 response → nil
- Registry module with mock 404 → error
- HTTP source with mock HEAD 200 → nil
- Timeout scenario → nil (transient, logged warning)
- Unknown/OCI source → nil (skip validation)

**Acceptance**:
- ≥90% coverage on `resolver.go`
- All edge cases covered
- `go test ./pkg/recipes/source/... -v` passes

---

## T4: Config generation for direct modules

**File**: `pkg/recipes/terraform/config/config.go` (MODIFY)
**Test file**: `pkg/recipes/terraform/config/config_test.go` (MODIFY)
**Depends on**: T2
**Implements**: FR-003, FR-004, R-003, R-007, R-010
**Description**: Add an `AddAllOutputs(moduleName string, outputs []ModuleOutputInfo)` method to the TF config generator. For direct modules, the generated config must:
1. Forward ALL module outputs (not just a single `result` block)
2. Mark outputs with `sensitive = true` when the module declares them sensitive

Existing `AddOutputs` for wrapped recipes must remain unchanged. Use `source.IsDirectModuleSource()` or accept a boolean flag to determine which path to use.

Also add `ModuleOutputInfo` struct if not already in types (name string, sensitive bool).

**Acceptance**:
- Wrapped recipes still use existing `result` output block
- Direct modules get individual output blocks for each module output
- Sensitive outputs correctly marked
- New unit tests for `AddAllOutputs`

---

## T5: Executor integration for direct modules

**File**: `pkg/recipes/terraform/execute.go` (MODIFY)
**Test file**: `pkg/recipes/terraform/execute_test.go` (MODIFY)
**Depends on**: T4
**Implements**: FR-002, FR-008, R-003
**Description**: Modify the executor's `Deploy` and `Destroy` methods to handle direct module sources:

1. In `Deploy`: after `terraform apply`, check if this is a direct module source. If yes, read all outputs individually (not just `result`). Map non-sensitive → Values, sensitive → Secrets.
2. In `Destroy`: no changes expected — `terraform destroy` works identically for direct and wrapped modules.
3. Fresh download each deploy (FR-002) — ensure no caching logic is added.

The executor should use `source.IsDirectModuleSource(templatePath)` to determine the output parsing strategy.

**Acceptance**:
- Direct module deploy returns all outputs split into Values/Secrets
- Wrapped module deploy continues to parse `result` output block
- Destroy works for both modes
- New test cases for direct module output parsing

---

## T6: Driver output mapping for direct modules

**File**: `pkg/recipes/driver/terraform/terraform.go` (MODIFY)
**Test file**: `pkg/recipes/driver/terraform/terraform_test.go` (MODIFY)
**Depends on**: T5
**Implements**: FR-004, FR-009, FR-011, R-003, R-009, R-010
**Description**: Modify `prepareRecipeResponse()` in the Terraform driver to handle direct module output format:

1. If direct module: outputs are already split into Values/Secrets by the executor — pass through directly to `RecipeOutput`
2. If wrapped module: existing parsing of `result` block with `Resources`, `Secrets`, `Values` sub-fields
3. Handle edge cases: module with no outputs (FR-011) → empty Values map, no error
4. Surface Terraform errors with actionable messages (FR-009, R-009)

**Acceptance**:
- Direct module outputs correctly mapped to RecipeOutput
- Wrapped module behavior unchanged
- Empty output set handled gracefully
- Error messages include Terraform error details
- New test cases for both paths

---

## T7: Source-aware module download

**File**: `pkg/recipes/terraform/module.go` (MODIFY)
**Depends on**: T2
**Implements**: FR-001, FR-002, FR-006, FR-007, R-001
**Description**: Modify the module download logic to handle direct Terraform module sources alongside existing OCI/wrapped paths. The existing code likely uses go-getter or terraform-exec for module fetching. Ensure:

1. Registry sources resolve via Terraform registry protocol
2. Git sources with `git::` prefix, `?ref=`, and `//subdir` work correctly
3. HTTP archive sources download and extract properly
4. S3/GCS sources work via go-getter's built-in detectors
5. Private auth uses existing credential/secret store mechanisms (FR-007)
6. No caching between deployments (FR-002)

**Acceptance**:
- All source types download correctly
- Version pinning works (registry version, git ref)
- Subdirectory paths resolve correctly
- Existing OCI download path unchanged

---

## T8: RecipePack creation-time validation

**File**: `pkg/corerp/frontend/controller/recipepacks/createorupdaterecipepack.go` (MODIFY)
**Test file**: `pkg/corerp/frontend/controller/recipepacks/createorupdaterecipepack_test.go` (MODIFY)
**Depends on**: T2
**Implements**: FR-014, R-004
**Description**: Add source reachability validation to the RecipePack create/update controller. When a recipe has `recipeKind: terraform`:

1. Classify the `recipeLocation` using `source.Classify()`
2. If it's a direct module source, call `source.ValidateReachability()`
3. Definitive failures (404, auth denied) → reject the RecipePack with descriptive error
4. Transient failures (timeout) → log warning, allow creation
5. Unknown/OCI sources → skip validation (existing behavior)

**Acceptance**:
- Unreachable registry module (404) → creation rejected
- Reachable registry module → creation succeeds
- Timeout → creation succeeds with warning log
- OCI/wrapped paths → no validation added
- New unit tests with mock HTTP for reachability checks

---

## T9: Functional tests

**Directory**: `test/functional/shared/resources/testdata/recipes/direct-module/` (NEW)
**Depends on**: T5, T6, T7, T8
**Implements**: All FRs, Constitution IV
**Description**: End-to-end functional tests that deploy a RecipePack with a direct Terraform module source, verify the resource is created, outputs are mapped correctly, and deletion cleans up. Use a simple, publicly available Terraform module (e.g., `hashicorp/null/null` or a test module in the repo).

Test scenarios:
1. Deploy with registry module source → verify outputs in Values
2. Deploy with git module source → verify outputs
3. Deploy module with sensitive output → verify in Secrets
4. Deploy module with no outputs → verify empty Values, no error
5. Delete resource → verify cleanup
6. Create RecipePack with unreachable source → verify rejection

**Acceptance**:
- All scenarios pass in CI
- Tests are self-contained (no external dependencies that could flake)
- Tests follow existing functional test patterns in the repo

---

## Summary

| Task | Files | Type | LOC Est. |
|------|-------|------|----------|
| T1 | `pkg/recipes/source/types.go` | New | ~60 |
| T2 | `pkg/recipes/source/resolver.go` | New | ~180 |
| T3 | `pkg/recipes/source/resolver_test.go` | New | ~250 |
| T4 | `pkg/recipes/terraform/config/config.go` + test | Modify | ~80 |
| T5 | `pkg/recipes/terraform/execute.go` + test | Modify | ~60 |
| T6 | `pkg/recipes/driver/terraform/terraform.go` + test | Modify | ~80 |
| T7 | `pkg/recipes/terraform/module.go` | Modify | ~40 |
| T8 | `pkg/corerp/.../createorupdaterecipepack.go` + test | Modify | ~60 |
| T9 | `test/functional/...` | New | ~150 |
| **Total** | | | **~960** |
