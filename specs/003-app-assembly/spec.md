# Feature Specification: Application Assembly

**Feature Branch**: `003-app-assembly`  
**Created**: 2026-02-17  
**Status**: Draft  
**Input**: User description: "Application Assembly — rad app discover and rad app model workflow for assembling apps from source code and configuration to deploy on Radius"

## Principles

1. **Developer-first**: Optimize for a seamless developer experience while honoring platform engineering standards. The default path should be simple for application teams and have basic compliance with organizational policies.

2. **Cognitive load reduction**: Use developer language and hide internal concepts (Resource Types, Recipes, Bicep) by default. Users should get value without learning the Radius resource model.

3. **Trust and transparency**: Always show what was inferred, from where, and with what confidence. Require review before deploy, and make it easy to accept, edit, or reject each inference.

4. **Progressive adoption**: Generate a working starting point that is easy to extend and customize as users learn more. The initial experience is low-friction, with a clear path to deeper control.

5. **Time to value**: Move from "existing app" to "running on Radius" in minutes, but never at the expense of correctness or safety.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Repository Scanning and Application Topology Discovery (Priority: P1)

As an enterprise application developer, I want to run `rad app discover` in my repository so that Radius automatically detects my services, their dependencies, and their connections, producing a human-readable architecture diagram and an editable manifest (`radius/app.yaml`).

**Why this priority**: Discovery is the entry point for the entire assembly workflow. Without accurate detection of services and dependencies, no downstream modeling or deployment is possible. This story delivers immediate value by showing developers what Radius "sees" in their codebase.

**Independent Test**: Can be fully tested by running `rad app discover` against a repository containing at least one Dockerfile and one external dependency reference. The output is a valid `radius/app.yaml` file and an architecture diagram rendered to the terminal.

**Acceptance Scenarios**:

1. **Given** a repository with one or more Dockerfiles and package manifests, **When** the developer runs `rad app discover`, **Then** the system detects all services, their runtimes, exposed ports, and inter-service connections and displays them in a component table.
2. **Given** the system has detected components, **When** the component table is presented, **Then** the developer can accept, edit, remove, or add components interactively before proceeding.
3. **Given** the developer has confirmed the topology, **When** discovery completes, **Then** a `radius/app.yaml` file is written containing the full application manifest with all confirmed services and dependencies.
4. **Given** a repository with no Dockerfiles or recognizable service artifacts, **When** the developer runs `rad app discover`, **Then** the system displays a clear message indicating no services were detected and suggests next steps (e.g., adding a Dockerfile or compose file).
5. **Given** a monorepo with multiple services in subdirectories, **When** the developer runs `rad app discover`, **Then** all services across subdirectories are detected and correctly associated with their respective Dockerfiles and configuration files.

---

### User Story 2 — Secret and Auth Detection and Configuration (Priority: P1)

As an enterprise application developer, I want `rad app discover` to detect secrets and authentication patterns in my codebase so that I can configure how each secret is sourced (secret store, environment variable, or skipped) and how each infrastructure dependency authenticates.

**Why this priority**: Secrets and auth are critical for security and compliance. Misconfigured secrets are a top deployment blocker and security risk. This story is co-equal with topology discovery because the manifest is incomplete without it.

**Independent Test**: Can be tested by running `rad app discover` against a repository containing `.env` files, source code with SDK client constructors, and configuration files referencing credentials. The output `radius/app.yaml` includes a correctly populated `secrets` section.

**Acceptance Scenarios**:

1. **Given** a repository with `.env` files, configuration files, and source code referencing secret values, **When** the developer runs `rad app discover`, **Then** the system detects all secret usage patterns and presents them with evidence (file path, line number, pattern matched).
2. **Given** detected secrets are presented, **When** the developer configures each secret, **Then** each secret can be assigned to a secret store (with a named key), kept as an environment variable, or skipped entirely.
3. **Given** infrastructure dependencies are detected, **When** the auth model step is presented, **Then** the developer can choose an auth method (e.g., managed identity, connection string, workload identity, API key) for each dependency.
4. **Given** a secret is detected only in source code and not in any `.env` file, **When** the secrets table is displayed, **Then** the system notes the absence of a `.env` entry and still allows the developer to configure it.

---

### User Story 3 — Platform Artifacts and Environment Configuration (Priority: P2)

As an enterprise application developer, I want `rad app discover` to locate and index my platform's IaC modules (from a local path or private registry), platform documentation, and deployment artifacts so that modeling can match dependencies to approved recipes.

**Why this priority**: Platform artifact integration ensures generated Bicep aligns with organizational standards. Without it, the developer would need to manually wire IaC modules, reducing the time-to-value promise.

**Independent Test**: Can be tested by running `rad app discover` against a repository with a local `./infra/terraform/` directory containing Terraform modules. The output `radius/app.yaml` includes a `platform` section referencing the discovered modules.

**Acceptance Scenarios**:

1. **Given** a local directory containing IaC modules, **When** the developer selects "Local path" and provides the path, **Then** the system scans and lists all discovered modules with their resource types.
2. **Given** IaC modules are hosted in a private registry, **When** the developer selects "Private registry" and authenticates, **Then** the system fetches and lists available modules with their versions.
3. **Given** the developer provides a platform documentation URL, **When** the URL is provided during discovery, **Then** the system indexes the documentation for use during modeling.
4. **Given** a deployment artifacts path is provided, **When** the path contains resource limits, replica counts, or ingress rules, **Then** these values are incorporated into the manifest's service definitions.
5. **Given** the developer confirms environment settings (app name, container registry, target environment), **When** a naming convention is selected, **Then** the manifest reflects the chosen naming pattern and Kubernetes labels.

---

### User Story 4 — Bicep Generation from Manifest (Priority: P1)

As an enterprise application developer, I want to run `rad app model` to read `radius/app.yaml` and generate a deployable `radius/app.bicep` along with resource type definitions and recipe configurations, so that I have a complete, reviewable deployment artifact.

**Why this priority**: Bicep generation is the culmination of the assembly workflow — it transforms the human-editable manifest into a deployable artifact. Without it, discovery alone does not enable deployment.

**Independent Test**: Can be tested by providing a valid `radius/app.yaml` and running `rad app model`. The output is a syntactically correct `radius/app.bicep` file plus supporting artifacts (`types.yaml`, `recipe-pack.yaml`).

**Acceptance Scenarios**:

1. **Given** a valid `radius/app.yaml` with services and dependencies, **When** the developer runs `rad app model`, **Then** the system presents the abstraction level choice (portable or custom resource types) and generates schema property tables for each resource.
2. **Given** schema properties are displayed for each resource, **When** the developer is prompted to edit, **Then** each property shows its current value and source, and the developer can accept all or edit individual properties.
3. **Given** the developer has confirmed schema properties, **When** recipe selection is presented, **Then** each dependency offers options to use an existing Terraform module, a community-contributed resource type, or an Azure Verified Module (where applicable).
4. **Given** all recipes are confirmed, **When** generation completes, **Then** the system writes `radius/app.bicep`, `radius/types.yaml`, and `radius/recipe-pack.yaml` and displays an application graph (Mermaid diagram) summarizing the topology.
5. **Given** the `radius/app.yaml` is missing required fields (e.g., no services defined), **When** the developer runs `rad app model`, **Then** the system displays a clear validation error listing the missing fields and does not generate any artifacts.

---

### User Story 5 — Editing Manifest Between Discovery and Modeling (Priority: P2)

As an enterprise application developer, I want to manually edit `radius/app.yaml` between discovery and modeling so that I can add services, change auth models, rename secrets, adjust dependencies, or modify environment settings before generating Bicep.

**Why this priority**: The manifest is the single source of truth for the assembly workflow. Editability ensures developers are never locked into inferences and can correct or extend anything between steps.

**Independent Test**: Can be tested by modifying `radius/app.yaml` (e.g., adding a new service entry or changing a secret source) and running `rad app model`. The generated Bicep should reflect the manual edits.

**Acceptance Scenarios**:

1. **Given** a `radius/app.yaml` generated by `rad app discover`, **When** the developer adds a new service entry manually, **Then** `rad app model` includes the new service in the generated Bicep.
2. **Given** a `radius/app.yaml` with a secret sourced from a secret store, **When** the developer changes the source to `env-var`, **Then** `rad app model` generates environment variable bindings instead of secret store references.
3. **Given** a `radius/app.yaml` with auto-detected naming conventions, **When** the developer changes the naming pattern, **Then** all generated resource names in Bicep follow the new pattern.

---

### User Story 6 — Re-running Discovery on an Existing Manifest (Priority: P3)

As an enterprise application developer, I want to re-run `rad app discover` on a repository that already has a `radius/app.yaml` so that new services or dependencies are detected and merged without losing my previous edits.

**Why this priority**: Applications evolve. Re-running discovery must be non-destructive to encourage iterative use as codebases change.

**Independent Test**: Can be tested by running `rad app discover` on a repository with an existing `radius/app.yaml`, adding a new Dockerfile to the repo, and verifying the new service appears in the updated manifest while previously confirmed entries remain unchanged.

**Acceptance Scenarios**:

1. **Given** a repository with an existing `radius/app.yaml`, **When** the developer runs `rad app discover`, **Then** the system detects new services not already in the manifest and presents them for confirmation alongside existing entries.
2. **Given** previously confirmed services exist in `radius/app.yaml`, **When** re-discovery runs, **Then** existing entries are preserved and marked as "previously confirmed" — they are not removed or overwritten.
3. **Given** a previously detected service has been removed from the repository, **When** re-discovery runs, **Then** the system flags the orphaned entry and asks the developer whether to keep or remove it.

---

### Edge Cases

- What happens when the repository contains multiple `Dockerfile` variants (e.g., `Dockerfile.dev`, `Dockerfile.prod`)? The system should present all variants and let the developer choose which to use per service.
- How does the system handle a `docker-compose.yaml` that defines services not backed by a local Dockerfile (e.g., pulling a pre-built image)? These should be detected as dependencies rather than services.
- What happens when an `.env` file contains a value that looks like a secret but is actually a non-sensitive configuration value (e.g., `LOG_LEVEL=debug`)? The system should use heuristics (naming patterns, value entropy) to classify and allow the developer to reclassify.
- How does the system handle authentication to a private registry that requires MFA or short-lived tokens? The system should support token-based auth and display a clear error with remediation steps if authentication fails.
- What happens when `rad app model` is run without a prior `rad app discover` (i.e., the developer hand-writes `radius/app.yaml`)? The system should accept any valid manifest regardless of how it was created.
- What happens when the repository contains no recognizable infrastructure dependencies? The system should generate a manifest with services only and skip the secrets/auth and recipe selection steps.
- How does the system handle conflicting information (e.g., a port defined in Dockerfile differs from the port in docker-compose.yaml)? The system should flag the conflict, show both sources, and ask the developer to choose.

## Requirements *(mandatory)*

### Functional Requirements

#### Discovery (`rad app discover`)

- **FR-001**: The system MUST scan the repository for Dockerfiles, compose files, package manifests, environment variables, connection strings, SDK imports, and configuration files in a single pass.
- **FR-002**: The system MUST detect services (containerized workloads) and classify each with a name, runtime, Dockerfile path, exposed port, and public/private designation.
- **FR-003**: The system MUST detect infrastructure dependencies (databases, caches, AI services, message brokers) from connection strings, SDK imports, and environment variable naming patterns.
- **FR-004**: The system MUST detect inter-service connections (one service calling another) from environment variables, configuration files, and source code references.
- **FR-005**: The system MUST render an ASCII architecture diagram showing services, dependencies, and their connections.
- **FR-006**: The system MUST present a component table listing all detected services and dependencies with evidence (source file and line number).
- **FR-007**: The system MUST allow the developer to accept, edit, remove, or add components interactively before finalizing the topology.
- **FR-008**: The system MUST scan for secret usage patterns across `.env` files, configuration files, source code (SDK calls, client constructors, auth helpers), and Docker/compose files.
- **FR-009**: The system MUST present detected secrets with evidence (file path, line number, matched pattern) and allow the developer to assign each to a secret store (with a named key), keep as environment variable, or skip.
- **FR-010**: The system MUST present auth model choices for each infrastructure dependency (e.g., managed identity, connection string, workload identity, API key) and record the selection.
- **FR-011**: The system MUST support IaC module discovery from a local directory path, scanning for Terraform or Bicep modules.
- **FR-012**: The system MUST support IaC module discovery from a private registry (e.g., Terraform Cloud, Artifactory), including authentication.
- **FR-013**: The system MUST accept a platform documentation URL and index the content for use during modeling.
- **FR-014**: The system MUST accept a deployment artifacts path and extract resource limits, replica counts, and ingress rules.
- **FR-015**: The system MUST detect environment settings (application name, container registry, target environment) from repository context and allow the developer to confirm or edit.
- **FR-016**: The system MUST support configurable naming conventions for generated resources (e.g., `{env}-{app}-{component}`).
- **FR-017**: The system MUST write all discovery results to `radius/app.yaml` in a documented, human-editable YAML format.
- **FR-018**: When re-running discovery on a repository with an existing `radius/app.yaml`, the system MUST preserve previously confirmed entries and present only new or changed detections for review.

#### Modeling (`rad app model`)

- **FR-019**: The system MUST read and validate `radius/app.yaml` before generating any artifacts, reporting clear errors for missing or malformed fields.
- **FR-020**: The system MUST present abstraction level choices (portable resource types or custom resource types) and apply the selection to all generated resources.
- **FR-021**: The system MUST display resolved schema properties for each resource (inputs and outputs) with the current value and its source, allowing the developer to accept all or edit individual properties.
- **FR-022**: The system MUST present recipe selection options for each dependency, including existing IaC modules (from the discovered platform artifacts), community-contributed resource types, and Azure Verified Modules.
- **FR-023**: The system MUST generate a syntactically valid `radius/app.bicep` file that declares all services, dependencies, connections, and secret bindings.
- **FR-024**: The system MUST generate `radius/types.yaml` containing all resource type definitions used by the application.
- **FR-025**: The system MUST generate `radius/recipe-pack.yaml` containing all recipe definitions and module references.
- **FR-026**: The system MUST display an application graph (Mermaid diagram) summarizing the generated topology after artifact generation.
- **FR-027**: The system MUST display a summary including application name, service count, dependency count, resource type count, recipe count, and secret count after generation.

#### General

- **FR-028**: All generated artifacts MUST be placed in the `radius/` directory within the repository root.
- **FR-029**: The system MUST require `rad init` to have been run before `rad app discover` can execute, and display a clear error message if initialization is missing.
- **FR-030**: The system MUST provide a non-interactive mode (e.g., `--yes` flag) that accepts all defaults for use in automation and CI pipelines.

### Key Entities

- **Service**: A containerized workload detected from a Dockerfile or compose file. Key attributes: name, runtime (language/framework), Dockerfile path, exposed port, public/private designation, resource limits.
- **Dependency**: An external infrastructure resource required by a service (e.g., database, cache, AI endpoint). Key attributes: name, type (database, cache, AI, message broker), detected source (env var, SDK import, config file).
- **Connection**: A relationship between two components (service-to-service or service-to-dependency). Key attributes: source, target, protocol, port.
- **Secret**: A sensitive value detected in the repository. Key attributes: environment variable name, evidence (file + line), source type (secret-store, env-var, skip), secret store key name.
- **Platform Artifact**: An IaC module or deployment artifact provided by the platform team. Key attributes: module name, source (local path or registry URL), IaC language, version.
- **Application Manifest** (`radius/app.yaml`): The central, human-editable YAML file that captures the full application topology, services, dependencies, secrets, auth model, platform artifacts, and environment settings.
- **Recipe**: A mapping from a dependency to a specific IaC module that provisions it. Key attributes: dependency name, resource type, module reference, module version.

## Assumptions

- `rad init` has been run prior to `rad app discover`, establishing the Radius workspace context.
- The repository contains at least one Dockerfile or compose file for service detection to produce results.
- Secret detection uses heuristic pattern matching (variable naming conventions, value entropy, known secret patterns) and may produce false positives — the interactive review step is the mitigation.
- Platform documentation indexing is best-effort and used as supplementary context during modeling, not as a strict enforcement mechanism.
- The `radius/app.yaml` schema is stable and versioned; breaking changes to the schema will include a migration path.
- Private registry authentication uses existing credential stores (e.g., Terraform Cloud tokens, Docker config) rather than requiring new credential management flows.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A developer with a multi-service repository (2+ services, 2+ dependencies) can go from `rad app discover` to a generated `radius/app.bicep` in under 10 minutes, including interactive review steps.
- **SC-002**: 90% of services and dependencies in a standard repository (Dockerfile + `.env` + common SDKs) are correctly detected without manual additions on the first discovery run.
- **SC-003**: 100% of detected secrets are presented for review before being written to the manifest — no secret value is ever written to `radius/app.yaml` in plaintext.
- **SC-004**: The generated `radius/app.bicep` is syntactically valid and deploys without policy violations when the platform team's IaC modules are correctly referenced.
- **SC-005**: Developers report understanding what was inferred and feeling confident modifying the manifest, measured by 80%+ task completion rate on a "modify and re-model" test task.
- **SC-006**: Re-running `rad app discover` on a previously discovered repository preserves 100% of previously confirmed entries (no data loss on re-discovery).
- **SC-007**: The assembly workflow reduces the time to first successful Radius deployment from hours (manual Bicep authoring) to under 15 minutes for a typical enterprise application.
