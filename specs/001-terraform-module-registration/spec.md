# Feature Specification: Direct Terraform Module as Recipe Template Path

**Feature Branch**: `001-terraform-module-registration`  
**Created**: 2025-07-17  
**Updated**: 2025-07-18  
**Status**: Draft  
**Input**: User description: "Enable platform engineers to use an existing Terraform module directly as the templatePath for a recipe in a RecipePack — without having to wrap the module, publish it to a separate registry, or create a Radius-specific recipe artifact."

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Use a Terraform Registry Module Directly as templatePath (Priority: P1)

As a platform engineer, I want to set the `templatePath` of a recipe in my RecipePack to point directly at a standard Terraform registry module (e.g., `hashicorp/consul/aws`) and have it work as-is — no wrapping, no republishing to a different location, and no Radius-specific recipe artifact creation.

**Today's workflow** requires: (1) find or write a Terraform module, (2) wrap it to conform to Radius recipe conventions (e.g., specific output structure, required variables), (3) publish the wrapped version to a registry or OCI artifact store, (4) reference the wrapped version as `templatePath`. **With this feature**: (1) set `templatePath` directly to the Terraform module source, (2) the system accepts it and provisions infrastructure using the module as-is.

**Why this priority**: This is the core value proposition — eliminating the "double work" of wrapping and republishing. It unlocks access to the entire Terraform module ecosystem for Radius recipes without any intermediate steps.

**Independent Test**: Can be fully tested by creating a RecipePack with a recipe whose `templatePath` points to a publicly available Terraform registry module (e.g., `hashicorp/consul/aws`), deploying a resource that uses that recipe, and verifying infrastructure is provisioned successfully.

**Acceptance Scenarios**:

1. **Given** a RecipePack definition with a recipe whose `templatePath` is set to a Terraform registry module (e.g., `hashicorp/consul/aws`), **When** the RecipePack is created or updated, **Then** the system accepts the `templatePath` as valid and persists the recipe definition.
2. **Given** a RecipePack with a recipe pointing to a Terraform registry module with a `templateVersion` specified, **When** the RecipePack is created, **Then** the system records and uses that specific version for deployments.
3. **Given** a `templatePath` pointing to a Terraform registry module that does not exist or is inaccessible, **When** a resource using this recipe is deployed, **Then** the system returns a clear error message indicating the module cannot be fetched.
4. **Given** a RecipePack with a recipe whose `templatePath` points to a registry module using standard Terraform source syntax (e.g., `namespace/name/provider`), **When** the recipe is resolved at deployment time, **Then** the system downloads and executes the module using standard Terraform module retrieval mechanisms.

---

### User Story 2 — Use a Git-Hosted Terraform Module Directly as templatePath (Priority: P1)

As a platform engineer, I want to set the `templatePath` of a recipe to a Git repository URL containing a Terraform module (e.g., `git::https://github.com/org/terraform-aws-vpc.git`) so that I can use modules from Git without wrapping or republishing them.

**Why this priority**: Many organizations store Terraform modules in Git repositories rather than registries. Supporting Git sources alongside registry sources covers the majority of real-world module hosting scenarios.

**Independent Test**: Can be fully tested by creating a recipe with `templatePath` set to a public Git URL containing a known Terraform module, deploying a resource, and verifying the module is fetched from Git and infrastructure is provisioned.

**Acceptance Scenarios**:

1. **Given** a recipe with `templatePath` set to `git::https://github.com/org/module.git`, **When** a resource using this recipe is deployed, **Then** the system clones the module from Git and executes it successfully.
2. **Given** a recipe with `templatePath` pointing to a Git URL with a ref specifier (e.g., `git::https://github.com/org/module.git?ref=v2.0.0`), **When** deployed, **Then** the system uses the specified Git ref (tag, branch, or commit).
3. **Given** a recipe with `templatePath` pointing to a Git URL with a subdirectory (e.g., `git::https://github.com/org/repo.git//modules/vpc`), **When** deployed, **Then** the system navigates to the subdirectory and executes the module found there.
4. **Given** a recipe with `templatePath` pointing to a private Git repository, **When** credentials are configured via the existing secret store mechanism, **Then** the system authenticates and fetches the module successfully.

---

### User Story 3 — Automatic Input Variable Passthrough as Recipe Parameters (Priority: P1)

As a platform engineer, I want the system to automatically pass recipe parameters through to the Terraform module's input variables when I use a direct module reference as `templatePath`, so that application developers can supply values for the module's inputs without me having to create a wrapping layer.

**Why this priority**: Without parameter passthrough, direct module references would be useless — there would be no way to supply inputs to the module. This is what makes the feature functional rather than just cosmetic.

**Independent Test**: Can be fully tested by creating a recipe whose `templatePath` points to a Terraform module with known required variables (e.g., `region`, `name`), deploying a resource that provides those variables as parameters, and verifying the module receives them and uses them during execution.

**Acceptance Scenarios**:

1. **Given** a recipe whose `templatePath` points to a Terraform module with required input variables, **When** a resource is deployed providing values for those variables via recipe parameters, **Then** the Terraform module receives those values as its input variables and executes successfully.
2. **Given** a Terraform module with optional input variables (those with default values), **When** a resource is deployed without providing values for optional variables, **Then** the module uses its own default values for those variables.
3. **Given** a recipe with environment-level parameters defined in the RecipePack and resource-level parameters supplied at deploy time, **When** the resource is deployed, **Then** both sets of parameters are passed to the module, with resource-level parameters taking precedence for any overlaps.
4. **Given** a Terraform module requires a variable that is not supplied via any parameter source, **When** deployment is attempted, **Then** the system surfaces the Terraform error clearly indicating which required variable is missing.

---

### User Story 4 — Terraform Module Outputs Available as Recipe Outputs (Priority: P1)

As an application developer, I want the Terraform module's output values to be available as the recipe's output after deployment, so that I can access provisioned resource details (endpoints, connection strings, IDs) through Radius without checking Terraform state directly.

**Why this priority**: Outputs are how provisioned infrastructure communicates back to the application. Without output surfacing, developers would have no way to use the provisioned resources from their application code.

**Independent Test**: Can be fully tested by deploying a recipe that uses a Terraform module with known outputs (e.g., `endpoint`, `id`), and verifying those output values are accessible through the deployed resource's properties in Radius.

**Acceptance Scenarios**:

1. **Given** a Terraform module with output values (e.g., `output "endpoint" { value = ... }`), **When** a resource using this module is successfully deployed, **Then** the module's outputs are surfaced in the recipe output and accessible on the deployed resource.
2. **Given** a Terraform module with sensitive outputs (marked `sensitive = true`), **When** the resource is deployed, **Then** sensitive outputs are stored securely and not exposed in plain text through standard resource queries.
3. **Given** a Terraform module with no output values defined, **When** the resource is deployed successfully, **Then** the deployment succeeds with an empty output set — the absence of outputs does not cause failure.

---

### User Story 5 — Inspect Module Schema Before Deployment (Priority: P2)

As a platform engineer, I want to retrieve the input variables and outputs of a Terraform module referenced by `templatePath` before deploying any resources, so that I can understand what parameters to configure and what outputs to expect.

**Why this priority**: Schema inspection aids discoverability and reduces trial-and-error during configuration, but is not required for the core deploy workflow. Engineers can consult module documentation as a workaround.

**Independent Test**: Can be fully tested by pointing the metadata retrieval mechanism at a known Terraform module and verifying it returns the module's variables (name, type, description, default, required) and outputs (name, type, description).

**Acceptance Scenarios**:

1. **Given** a recipe with `templatePath` pointing to a Terraform module, **When** the platform engineer requests the recipe's metadata, **Then** the system returns the module's input variables (name, type, description, default value, required/optional) and output values (name, type, description).
2. **Given** a Terraform module with complex input types (maps, lists, objects), **When** metadata is retrieved, **Then** the type information is presented accurately reflecting the module's type constraints.
3. **Given** a `templatePath` that points to an inaccessible or invalid module, **When** metadata retrieval is attempted, **Then** the system returns a clear error without crashing or timing out silently.

---

### User Story 6 — Version Pinning for Direct Module References (Priority: P2)

As a platform engineer, I want to pin the version of a Terraform module referenced by `templatePath` using the existing `templateVersion` field, so that deployments are reproducible and I control when to adopt newer module versions.

**Why this priority**: Version pinning is essential for production stability but engineers can initially work around it by specifying version in the source URL or accepting latest.

**Independent Test**: Can be fully tested by creating a recipe with `templatePath` pointing to a registry module and `templateVersion` set to a specific version, deploying, and verifying the pinned version is used.

**Acceptance Scenarios**:

1. **Given** a recipe with `templatePath` set to a registry module and `templateVersion` set to `1.2.0`, **When** a resource is deployed, **Then** the system uses version `1.2.0` of the module regardless of what newer versions exist.
2. **Given** a recipe with `templatePath` set to a registry module and no `templateVersion` specified, **When** a resource is deployed, **Then** the system uses the latest available version of the module.
3. **Given** a `templateVersion` that does not exist for the specified module, **When** deployment is attempted, **Then** the system returns a clear error indicating the version is not found.

---

### Edge Cases

- What happens when a Terraform module has no input variables? The recipe deploys successfully with no parameters required — the module runs with its internal defaults.
- What happens when a Terraform module has no outputs? The recipe deploys successfully — the recipe output set is empty, and success is determined by the module execution completing without error.
- What happens when the module source becomes unavailable after initial successful deployments (e.g., registry down, Git repo deleted)? Existing deployed resources are unaffected. New deployments fail with a clear error indicating the source is unreachable.
- How does the system handle Terraform modules that expect specific provider configurations? The existing Terraform driver's provider configuration mechanism applies. Providers are configured through the environment's recipe configuration, consistent with how wrapped recipes work today.
- What happens when a parameter name supplied to the recipe does not match any input variable in the module? Terraform's standard behavior applies — unrecognized variables produce a clear error from `terraform apply`.
- What happens when the `templatePath` format is ambiguous (could be a registry module or a wrapped recipe path)? The system uses Terraform's standard source resolution — registry-style paths (`namespace/name/provider`) are resolved as registry modules; URLs with explicit schemes are resolved accordingly.
- What happens when a module's source has breaking changes between versions? The platform engineer must update `templateVersion` deliberately. No automatic version bumping occurs.
- What happens when a module uses Terraform features not yet supported by the Radius Terraform driver (e.g., specific provider plugins)? The module execution fails with the Terraform error surfaced through the standard recipe failure path. The system does not silently swallow provider-related errors.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept a standard Terraform module source (registry path, Git URL, HTTP URL, S3 URL, or GCS URL) as the `templatePath` value in a `RecipeDefinition` within a RecipePack, without requiring the module to be wrapped, republished, or conformed to a Radius-specific format.
- **FR-002**: System MUST resolve and download the Terraform module referenced by `templatePath` fresh at each deployment time using standard Terraform module retrieval mechanisms (registry protocol, Git clone, HTTP download). No module caching is performed between deployments.
- **FR-003**: System MUST pass recipe parameters (both environment-level from the RecipePack and resource-level from the deployment request) through to the Terraform module as input variables during execution.
- **FR-004**: System MUST surface Terraform module output values as recipe outputs after successful deployment. Non-sensitive outputs are placed into the `Values` map; sensitive outputs (those marked `sensitive = true` in the Terraform module) are placed into the `Secrets` map. No automatic `Resources` detection is performed — outputs are split only between `Values` and `Secrets` based on the module's sensitivity annotation.
- **FR-005**: System MUST support the `templateVersion` field to pin a specific version of a registry-hosted Terraform module, using the pinned version for all deployments of that recipe.
- **FR-006**: System MUST support Git-hosted modules via `git::` prefixed URLs in `templatePath`, including ref specifiers and subdirectory paths.
- **FR-007**: System MUST support authentication for private module sources (private registries and private Git repositories) using the existing secret store and credential configuration mechanisms.
- **FR-008**: System MUST execute a `terraform destroy` operation when a resource deployed via a direct module recipe is deleted, cleaning up all infrastructure provisioned by the module.
- **FR-009**: System MUST surface Terraform execution errors (missing variables, provider failures, permission errors) as recipe deployment failures with actionable error messages that include the relevant Terraform error details.
- **FR-010**: System MUST provide module metadata retrieval (input variables and outputs) for recipes using direct module references, through the existing `GetRecipeMetadata` mechanism.
- **FR-011**: System MUST handle Terraform modules with no input variables (deploy with no parameters) and modules with no outputs (succeed with empty output set) without errors.
- **FR-012**: System MUST NOT require any changes to the existing `RecipePack`, `RecipeDefinition`, or `EnvironmentDefinition` data models — the feature operates within the existing schema by expanding what `templatePath` values are accepted.
- **FR-013**: System MUST ensure that the existing recipe-based workflows (wrapped recipes, OCI-published recipes) continue to function identically — zero behavioral changes to existing deployments.
- **FR-014**: System MUST perform best-effort validation that a `templatePath` pointing to a direct Terraform module source is reachable at RecipePack creation time, using lightweight probes (registry API query, `git ls-remote`, HTTP HEAD) with a 30-second timeout. Definitive failures (HTTP 404, authentication denied, malformed URL) MUST reject the RecipePack. Transient failures (timeouts, temporary network errors) MUST be logged as warnings but MUST NOT block creation.
- **FR-015**: System MUST resolve `templatePath` by first attempting direct Terraform module resolution, then falling back to the existing OCI/wrapped recipe path if direct resolution fails. No new fields or explicit flags are required to distinguish direct modules from wrapped recipes.

### Key Entities

- **RecipeDefinition (extended behavior)**: The existing entity that defines a recipe within a RecipePack. With this feature, the `RecipeLocation` (templatePath) field now additionally accepts standard Terraform module sources — registry paths (`namespace/name/provider`), Git URLs (`git::https://...`), HTTP URLs, S3 URLs (`s3::https://...`), and GCS URLs (`gcs::https://...`) — alongside the existing wrapped recipe references. No structural changes to the entity itself.
- **Recipe Parameters (passthrough)**: The existing `Parameters` field on `RecipeDefinition` and resource-level parameters supplied at deploy time. With this feature, these parameters are passed directly as Terraform input variables to the referenced module, with no intermediate mapping layer required.
- **Recipe Output (from module outputs)**: The existing recipe output mechanism. With this feature, Terraform module `output` blocks are automatically surfaced through the standard recipe output path, making them available on the deployed resource.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A platform engineer can configure a recipe to use an existing Terraform module in under 1 minute by setting `templatePath` directly to the module source — zero wrapping or republishing steps required.
- **SC-002**: End-to-end provisioning time for a recipe using a direct Terraform module reference is within 10% of equivalent wrapped-recipe deployments — no significant overhead from the direct reference path.
- **SC-003**: 100% of a Terraform module's output values are automatically accessible as recipe outputs after successful deployment — zero manual output mapping configuration required.
- **SC-004**: Deleting a resource deployed via a direct module recipe fully destroys the underlying infrastructure with zero orphaned cloud resources.
- **SC-005**: Deployment of a recipe with an inaccessible `templatePath` fails within 60 seconds with a clear error message — no silent hangs or unclear failures.
- **SC-006**: Existing recipe-based workflows (wrapped recipes, OCI-published recipes) continue to function with zero behavioral changes after this feature is introduced.
- **SC-007**: Any Terraform module from the public Terraform Registry that does not require custom Radius-specific conventions can be used directly as a recipe template without modification to the module.
- **SC-008**: A platform engineer can inspect the input variables and outputs of a direct module recipe through the existing metadata retrieval mechanism before deploying any resources.

## Assumptions

- **A-001**: The feature works within the existing `RecipePack`/`RecipeDefinition`/`EnvironmentDefinition` data model. No new resource types, API endpoints, or "registration" concepts are introduced. The `templatePath` field is simply broadened in what source formats it accepts.
- **A-002**: Terraform module input variable types are passed through as recipe parameters without type transformation. Type checking is delegated to Terraform at execution time, which produces clear errors for type mismatches.
- **A-003**: This feature operates alongside the existing recipe workflow. Wrapped/published recipes continue to work exactly as before. The system resolves `templatePath` by trying direct Terraform module resolution first, then falling back to the existing OCI/wrapped path — no explicit flag or field distinguishes the two.
- **A-009**: No new observability signals (logs, metrics) are introduced for the direct module path. The existing Terraform driver logs and recipe engine metrics provide sufficient visibility.
- **A-004**: Infrastructure state management uses the existing Terraform state storage mechanism in the Radius Terraform driver. No new state backend is introduced.
- **A-005**: The existing Terraform driver's provider configuration (via `recipes.Configuration`) applies to direct module references. No additional provider setup mechanism is needed.
- **A-006**: Module introspection for metadata retrieval uses the existing `GetRecipeMetadata` path, which already runs `terraform get` to fetch module information. The feature extends this to work with broader source formats.
- **A-007**: Local filesystem paths as `templatePath` remain out of scope for the initial version (currently explicitly rejected in validation). The initial scope covers registry modules and Git-hosted modules.
- **A-008**: Standard Terraform execution timeouts and retry behavior apply. The feature does not introduce new timeout or retry configuration.
