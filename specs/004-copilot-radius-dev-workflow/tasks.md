# Tasks: Developer Workflow in Copilot and GitHub Using Radius

**Feature**: `004-copilot-radius-dev-workflow`
**Generated from**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/)
**User Stories**: 10 (P1: US1–US5, P2: US6–US8, US10, P3: US9)
**Total Tasks**: 119

---

## Phase 1: Setup

Goal: Initialize project structure, Go packages, and build targets for all new components.

- [X] T001 Create package directory structure for all new packages (`pkg/constitution/`, `pkg/agentskill/`, `pkg/mcpserver/`, `pkg/mcpserver/tools/`, `pkg/artifacts/`, `pkg/reporadius/`, `test/functional/copilot-workflow/`)
- [X] T002 [P] Add Makefile targets for new packages (test, lint, build) in `build/*.mk`
- [X] T003 [P] Create `rad mcp serve` subcommand scaffold in `cmd/rad/cmd/mcp/serve.go`

---

## Phase 2: Foundational

Goal: Implement core shared types that all user stories depend on. MUST complete before user story phases.

- [X] T004 Define PlatformConstitution, CloudProvider, ComputeTarget, IaCTooling, ModuleCatalog, Policy, NamingConvention, ContainerImagePolicy, RepositoryPolicy types in `pkg/constitution/types.go`
- [X] T005 [P] Define ApplicationGraph, GraphNode, GraphEdge, ModuleReference, ApplicationDefinition, ServiceDefinition, RecipeConfiguration, EnvironmentConfiguration types in `pkg/artifacts/types.go`
- [X] T006 [P] Define MCPToolDefinition and MCP protocol request/response types in `pkg/mcpserver/types.go`
- [X] T007 [P] Define DeploymentPlan, DeploymentAction, PolicyResult types in `pkg/reporadius/types.go`
- [X] T008 Define shared interfaces (ConstitutionParser, ConstitutionValidator, ConstitutionGenerator, ResourceTypeResolver, RecipeSearcher, ArtifactGenerator) in `pkg/constitution/interfaces.go` and `pkg/artifacts/interfaces.go`

---

## Phase 3: User Story 1 — Platform Engineering Constitution Discovery and Generation (P1)

**Story goal**: A developer can run `rad init` and get a valid `Platform-Engineering-Constitution.md` generated from interactive prompts capturing organization policies, cloud providers, compute targets, IaC tooling, and module catalogs.

**Independent test**: Run `rad init` in an empty directory, answer prompts (or use defaults), verify a valid `Platform-Engineering-Constitution.md` is generated with configured policies, cloud providers, compute targets, IaC tooling, and imported policy rules.

**FRs covered**: FR-001, FR-002, FR-003, FR-007, FR-008, FR-033

### Implementation

- [X] T009 [US1] Implement constitution YAML frontmatter parser (extract structured data from Markdown+YAML) in `pkg/constitution/parser.go`
- [X] T010 [US1] Implement constitution Markdown body parser (extract human-readable sections) in `pkg/constitution/parser.go`
- [X] T011 [US1] Implement constitution generator using Go `text/template` in `pkg/constitution/generator.go`
- [X] T012 [P] [US1] Create constitution Markdown+YAML generation templates in `pkg/constitution/templates/constitution.md.tmpl`
- [X] T013 [US1] Implement constitution validator against JSON Schema (`contracts/constitution-schema.yaml`) in `pkg/constitution/validator.go`
- [X] T014 [US1] Implement policy import logic for Azure Policy and AWS Config sources in `pkg/constitution/import.go`
- [X] T015 [US1] Add constitution interactive prompts (org name, cloud providers, compute targets, IaC tooling, module catalog, policies) in `pkg/cli/cmd/radinit/constitution.go`
- [X] T016 [US1] Implement update-existing-constitution logic (detect existing file, offer merge, preserve settings) in `pkg/cli/cmd/radinit/constitution.go`
- [X] T017 [US1] Add sensible defaults for all constitution prompts (skip-friendly initialization) in `pkg/cli/cmd/radinit/constitution.go`
- [X] T018 [US1] Add `--non-interactive` flag support accepting constitution config from file in `pkg/cli/cmd/radinit/init.go`
- [X] T019 [US1] Wire constitution phase into `enterInitOptions()` chain (before cluster phase) in `pkg/cli/cmd/radinit/init.go`
- [X] T020 [P] [US1] Add unit tests for constitution YAML frontmatter and Markdown parsing in `pkg/constitution/parser_test.go`
- [X] T021 [P] [US1] Add unit tests for constitution generation with various template inputs in `pkg/constitution/generator_test.go`
- [X] T022 [P] [US1] Add unit tests for constitution validation (valid, invalid, edge cases) in `pkg/constitution/validator_test.go`
- [X] T023 [US1] Add unit tests for radinit constitution prompts and `enterInitOptions()` integration in `pkg/cli/cmd/radinit/constitution_test.go`

### Checkpoint: US1

Verify: `rad init` generates a valid `Platform-Engineering-Constitution.md` with all configured fields. The file round-trips through parser → validator without errors. Defaults produce a usable constitution.

---

## Phase 4: User Story 2 — Radius Agent Skill and MCP Server Initialization (P1)

**Story goal**: `rad init` also generates a Radius Agent Skill (`.github/skills/radius-platform/SKILL.md`) grounded on the constitution, writes MCP server config (`.github/copilot/mcp.json`), and `rad mcp serve` starts a functional stdio JSON-RPC server.

**Independent test**: Run `rad init`, verify SKILL.md exists with constitution-derived content, mcp.json points to `rad mcp serve`, and `rad mcp serve` starts and responds to `initialize` and `tools/list` JSON-RPC requests.

**FRs covered**: FR-004, FR-005, FR-006, FR-034

### Implementation

- [X] T024 [US2] Create Agent Skill Markdown template with constitution-parameterized sections (org context, policies, workflow instructions, MCP tool descriptions) in `pkg/agentskill/templates/skill.md.tmpl`
- [X] T025 [US2] Implement Agent Skill generator from PlatformConstitution data in `pkg/agentskill/generator.go`
- [X] T026 [US2] Add Agent Skill generation step to `rad init` in `pkg/cli/cmd/radinit/agentskill.go`
- [X] T027 [US2] Implement MCP server configuration writer (generate mcp.json pointing to `rad mcp serve`) in `pkg/mcpserver/config.go`
- [X] T028 [US2] Add MCP server registration step to `rad init` in `pkg/cli/cmd/radinit/mcpserver.go`
- [X] T029 [US2] Implement `rad mcp serve` command with stdio JSON-RPC 2.0 transport in `cmd/rad/cmd/mcp/serve.go`
- [X] T030 [US2] Implement MCP protocol handler (initialize, tools/list, tools/call dispatch) in `pkg/mcpserver/server.go`
- [X] T031 [US2] Wire Agent Skill and MCP server phases into `enterInitOptions()` chain (after constitution, before cluster) in `pkg/cli/cmd/radinit/init.go`
- [X] T032 [P] [US2] Add snapshot/golden file tests for Agent Skill generation from sample constitutions in `pkg/agentskill/generator_test.go`
- [X] T033 [P] [US2] Add unit tests for MCP config writer (mcp.json structure, paths) in `pkg/mcpserver/config_test.go`
- [X] T034 [P] [US2] Add MCP protocol compliance tests (initialize, tools/list, tools/call, error handling) in `pkg/mcpserver/server_test.go`
- [X] T035 [US2] Add unit tests for radinit Agent Skill and MCP registration prompts in `pkg/cli/cmd/radinit/agentskill_test.go` and `pkg/cli/cmd/radinit/mcpserver_test.go`

### Checkpoint: US2

Verify: `rad init` generates SKILL.md at `.github/skills/radius-platform/`, mcp.json at `.github/copilot/`, and `rad mcp serve` starts successfully, responding to `initialize` and `tools/list` JSON-RPC messages.

---

## Phase 5: User Story 3 — Application Architecture Design via Copilot and Radius (P1)

**Story goal**: The MCP server exposes tools that enable Copilot to analyze source code, resolve resource types, search recipes, validate architecture against the constitution, and list available environments.

**Independent test**: Start MCP server, send `tools/call` for `resolve_resource_types` with test dependencies and `search_recipes` with a resource type, verify correct deterministic responses matching the constitution and catalog.

**FRs covered**: FR-009, FR-010, FR-011, FR-012, FR-013, FR-014 *(FR-014 is fulfilled by Copilot's conversational ability to re-invoke MCP tools with a modified graph; no dedicated task required)*

### Implementation

- [X] T036 [US3] Implement `resolve_resource_types` MCP tool with 3-tier resolution (local → community → generated) in `pkg/mcpserver/tools/resolve_resource_types.go`
- [X] T037 [US3] Implement local resource type scanner (scan `radius/types/` directory) in `pkg/artifacts/resourcetype.go`
- [X] T038 [US3] Implement community catalog fetcher (parse resource-types-contrib YAML index) in `pkg/artifacts/resourcetype.go`
- [X] T039 [US3] Implement custom resource type generator fallback (minimal type from dependency pattern) in `pkg/artifacts/resourcetype.go`
- [X] T040 [US3] Implement `search_recipes` MCP tool (filter by resource type, provider, IaC tooling) in `pkg/mcpserver/tools/search_recipes.go`
- [X] T041 [US3] Implement recipe catalog search against constitution's module catalog URL in `pkg/artifacts/recipe.go`
- [X] T042 [US3] Implement `validate_constitution` MCP tool (validate artifacts or graph against constitution policies) in `pkg/mcpserver/tools/validate_constitution.go`
- [X] T043 [US3] Implement `get_environments` MCP tool (list available Radius environments) in `pkg/mcpserver/tools/get_environments.go`
- [X] T099 [US3] Implement `discover_application` MCP tool wrapping `rad app discover` logic (scan for services, dependencies, connections) in `pkg/mcpserver/tools/discover_application.go`
- [X] T100 [US3] Implement architecture diagram text renderer (generate Mermaid diagram from ApplicationGraph for design-phase display) in `pkg/artifacts/diagram.go`
- [X] T101 [US3] Implement architecture pattern detector (classify topology: multi-service, event-driven, microservices, monolith) in `pkg/artifacts/pattern.go`
- [X] T044 [P] [US3] Add unit tests for `resolve_resource_types` tool (all 3 tiers, confidence scores) in `pkg/mcpserver/tools/resolve_resource_types_test.go`
- [X] T045 [P] [US3] Add unit tests for `search_recipes` tool (filters, approved flag) in `pkg/mcpserver/tools/search_recipes_test.go`
- [X] T046 [P] [US3] Add unit tests for `validate_constitution` tool (pass, fail, warning scenarios) in `pkg/mcpserver/tools/validate_constitution_test.go`
- [X] T047 [P] [US3] Add unit tests for `get_environments` tool in `pkg/mcpserver/tools/get_environments_test.go`
- [X] T102 [P] [US3] Add unit tests for `discover_application` tool in `pkg/mcpserver/tools/discover_application_test.go`
- [X] T103 [P] [US3] Add unit tests for architecture diagram renderer and pattern detector in `pkg/artifacts/diagram_test.go` and `pkg/artifacts/pattern_test.go`

### Checkpoint: US3

Verify: MCP server responds correctly to all 5 design-phase tool calls (resolve, search, validate, environments, discover) with deterministic outputs. Resource type resolution follows the 3-tier priority. Recipe search filters by constitution preferences. Pattern detector classifies known topologies. Diagram renderer produces valid Mermaid.

---

## Phase 6: User Story 4 — Application Definition and IaC Artifact Generation (P1)

**Story goal**: The `generate_app_definition` MCP tool produces a complete set of deployment artifacts — `app.bicep`, resource type definitions, recipe configurations, environment configs, and `app-graph.json` — all validated against the constitution.

**Independent test**: Call `generate_app_definition` with a test ApplicationGraph (2 services, 2 infrastructure dependencies), verify all output files exist, are syntactically valid, reference approved modules, and pass constitution validation.

**FRs covered**: FR-015, FR-016, FR-017, FR-018, FR-019, FR-020, FR-032

### Implementation

- [X] T048 [US4] Implement application definition generator (multi-service `app.bicep` with connections and secret bindings) in `pkg/artifacts/appdef.go`
- [X] T049 [P] [US4] Create `app.bicep` Go template with parameterized services, connections, and secret bindings in `pkg/artifacts/templates/app.bicep.tmpl`
- [X] T050 [US4] Implement recipe configuration generator (wrap approved IaC modules from constitution catalog) in `pkg/artifacts/recipeconfig.go`
- [X] T051 [US4] Implement environment configuration generator (`radius/environments/*.bicep` from constitution preferences) in `pkg/artifacts/environment.go`
- [X] T052 [US4] Implement application graph generator (`app-graph.json` from ApplicationGraph type) in `pkg/artifacts/graph.go`
- [X] T053 [US4] Implement `generate_app_definition` MCP tool orchestrating the full artifact pipeline in `pkg/mcpserver/tools/generate_app_definition.go`
- [X] T054 [US4] Implement artifact compliance validation (validate all generated files against constitution policies) in `pkg/artifacts/validate.go`
- [X] T055 [P] [US4] Add unit tests for application definition generator (single-service, multi-service, connections) in `pkg/artifacts/appdef_test.go`
- [X] T056 [P] [US4] Add unit tests for recipe configuration generator (Terraform, Bicep, approved/unapproved modules) in `pkg/artifacts/recipeconfig_test.go`
- [X] T057 [P] [US4] Add unit tests for environment configuration generator in `pkg/artifacts/environment_test.go`
- [X] T058 [P] [US4] Add unit tests for application graph generator (nodes, edges, metadata) in `pkg/artifacts/graph_test.go`
- [X] T059 [US4] Add unit tests for `generate_app_definition` MCP tool (full pipeline, compliance report) in `pkg/mcpserver/tools/generate_app_definition_test.go`

### Checkpoint: US4

Verify: Calling `generate_app_definition` produces `radius/app.bicep`, `radius/app-graph.json`, `radius/types/*.yaml`, `radius/recipes/*/main.tf`, and `radius/environments/*.bicep`, all passing constitution validation with zero violations for standard architectures.

---

## Phase 7: User Story 5 — Pull Request-Based Deployment via Repo Radius (P1)

**Story goal**: Radius artifacts can be submitted as a PR that includes a deployment workflow. Merging the PR triggers Repo Radius to deploy the application.

**Independent test**: Generate a GitHub Actions workflow file, verify it has the correct triggers (PR merge to main with `radius/` path filter), runs `rad deploy`, and posts status comments.

**FRs covered**: FR-021, FR-024, FR-027

### Implementation

- [X] T060 [US5] Implement GitHub Actions workflow generator (radius-deploy.yml with merge trigger, rad deploy, status comments) in `pkg/reporadius/workflow.go`
- [X] T061 [P] [US5] Create workflow template (`radius-deploy.yml.tmpl`) with PR merge trigger and `radius/` path filter in `pkg/reporadius/templates/radius-deploy.yml.tmpl`
- [X] T062 [US5] Implement PR creation helper (feature branch, commit artifacts, PR title/summary) in `pkg/reporadius/pullrequest.go`
- [X] T063 [US5] Implement deployment orchestration logic (compile app definition, provision environment, deploy) in `pkg/reporadius/deployment.go`
- [X] T064 [P] [US5] Add unit tests for workflow generator (template rendering, trigger config) in `pkg/reporadius/workflow_test.go`
- [X] T065 [P] [US5] Add unit tests for PR creation helper (branch naming, commit content, PR body) in `pkg/reporadius/pullrequest_test.go`
- [X] T066 [US5] Add unit tests for deployment orchestration logic in `pkg/reporadius/deployment_test.go`

### Checkpoint: US5

Verify: Generated `.github/workflows/radius-deploy.yml` is valid YAML, triggers on push to main with `radius/**` path filter, and includes `rad deploy` step with status comment posting.

---

## Phase 8: User Story 6 — Architecture Graph Visualization in Pull Requests and README (P2)

**Story goal**: PRs containing Radius artifacts display a rich architecture graph visualization (Mermaid diagram) as a PR comment, with before/after diffs for modifications. The repository README is automatically updated with the current application graph after each deployment.

**Independent test**: Generate a Mermaid diagram from a test `app-graph.json`, verify it includes all nodes (services and infrastructure) and edges (connections), and renders correctly in Markdown. Verify README is updated with graph.

**FRs covered**: FR-022, FR-023, FR-035, FR-036

### Implementation

- [X] T067 [US6] Implement Mermaid diagram generator from ApplicationGraph (nodes with types, edges with labels) in `pkg/reporadius/visualization.go`
- [X] T068 [US6] Implement before/after architecture diff (compare two ApplicationGraphs, mark added/modified/removed) in `pkg/reporadius/visualization.go`
- [X] T069 [US6] Implement PR comment formatter (graph + compliance summary + component details) in `pkg/reporadius/visualization.go`
- [X] T070 [US6] Add visualization step to `radius-deploy.yml` workflow template (generate and post PR comment) in `pkg/reporadius/templates/radius-deploy.yml.tmpl`
- [X] T105 [US6] Implement README graph updater (embed/update Mermaid diagram between markers in README.md) in `pkg/reporadius/readme.go`
- [X] T106 [US6] Add README update step to deployment workflow (after successful deployment, update README with latest graph) in `pkg/reporadius/templates/radius-deploy.yml.tmpl`
- [X] T071 [P] [US6] Add unit tests for Mermaid diagram generation (various topologies, node types) in `pkg/reporadius/visualization_test.go`
- [X] T072 [P] [US6] Add unit tests for architecture diff logic (additions, modifications, deletions) in `pkg/reporadius/visualization_test.go`
- [X] T107 [P] [US6] Add unit tests for README graph updater (insert new, update existing, preserve surrounding content) in `pkg/reporadius/readme_test.go`

### Checkpoint: US6

Verify: Given an `app-graph.json`, the system produces a valid Mermaid diagram showing all nodes and edges. Given two graphs, the diff correctly shows added, modified, and removed components. README is updated with latest graph between markers.

---

## Phase 9: User Story 7 — Deployment Plan Review and Approval (P2)

**Story goal**: A detailed deployment plan is generated before provisioning, showing all resources to be created/modified/deleted with policy compliance, and supporting approval/rejection through GitHub environment protection rules.

**Independent test**: Generate a deployment plan from current vs. proposed `app-graph.json`, verify all resources are listed with their actions (create/update/delete), policy results are included, and the plan renders as a PR comment.

**FRs covered**: FR-025, FR-026

### Implementation

- [X] T073 [US7] Implement deployment plan generator (diff current vs. proposed ApplicationGraph, produce DeploymentAction list) in `pkg/reporadius/deployment.go`
- [X] T074 [US7] Implement policy violation highlighting in deployment plan (flag violations with remediation suggestions) in `pkg/reporadius/deployment.go`
- [X] T075 [US7] Implement deployment plan PR comment renderer (rich table with resources, actions, properties, policy status) in `pkg/reporadius/deployment.go`
- [X] T076 [US7] Add approval/rejection flow using GitHub environment protection rules to workflow template in `pkg/reporadius/workflow.go`
- [X] T077 [P] [US7] Add unit tests for deployment plan generation (create-only, mixed actions, no-change) in `pkg/reporadius/deployment_test.go`
- [X] T078 [P] [US7] Add unit tests for approval flow logic and plan rendering in `pkg/reporadius/workflow_test.go`

### Checkpoint: US7

Verify: Deployment plan lists all create/update/delete actions with resource properties, includes policy compliance results, and the workflow template includes environment protection with approval gate.

---

## Phase 10: User Story 8 — Constitution-Compliant Application Scaffolding (P2)

**Story goal**: Copilot scaffolds application code, Dockerfiles, and repository settings that comply with the platform constitution from inception.

**Independent test**: Request scaffolding in a repo with a constitution specifying approved base images and tagging policy, verify generated Dockerfiles use those images and tags.

**FRs covered**: FR-029, FR-030, FR-031

### Implementation

- [X] T079 [US8] Implement constitution-aware Dockerfile generator (select approved base images, apply tagging policy) in `pkg/artifacts/scaffold.go`
- [X] T080 [US8] Implement repository settings validator against constitution (visibility, branch protection, default branch) in `pkg/artifacts/scaffold.go`
- [X] T081 [US8] Implement container image policy enforcement (scanning and signing requirement checks) in `pkg/artifacts/scaffold.go`
- [X] T082 [US8] Add scaffolding instructions and constitution-compliance guidance to Agent Skill template in `pkg/agentskill/templates/skill.md.tmpl`
- [X] T083 [P] [US8] Add unit tests for constitution-aware scaffolding (Dockerfile, repo settings, image policies) in `pkg/artifacts/scaffold_test.go`

### Checkpoint: US8

Verify: Generated Dockerfiles reference approved base images and tagging policies from the constitution. Repository settings validator flags non-compliant configurations.

---

## Phase 11: User Story 9 — Deployment Status and Notifications (P3)

**Story goal**: Developers receive timely deployment progress updates and completion/failure notifications in the PR comments.

**Independent test**: Simulate deployment stages, verify progress update comments and final summary (with resource links for success, diagnostics for failure) are correctly formatted.

**FRs covered**: FR-028

### Implementation

- [X] T084 [US9] Implement deployment progress reporter (stage-based PR comment updates: environment setup → provisioning → deployment) in `pkg/reporadius/notifications.go`
- [X] T085 [US9] Implement deployment completion summary (success: resource links + final graph; failure: error, failed resource, diagnostics, remediation) in `pkg/reporadius/notifications.go`
- [X] T086 [US9] Add notification steps to `radius-deploy.yml` workflow template (post at each stage, post final summary) in `pkg/reporadius/templates/radius-deploy.yml.tmpl`
- [X] T087 [P] [US9] Add unit tests for notification formatting (progress, success, failure scenarios) in `pkg/reporadius/notifications_test.go`

### Checkpoint: US9

Verify: Progress updates and final deployment status comments are correctly formatted with resource links (success) or failure diagnostics (failure).

---

## Phase 11b: User Story 10 — Enterprise Application Graph Database and Organizational Intelligence (P2)

**Story goal**: Application graph data is published to a centralized enterprise graph database upon deployment. The database supports cross-organizational queries correlating graph data with commit history, IaC versions, and team ownership.

**Independent test**: Publish two sample application graphs with different metadata (teams, IaC versions), verify cross-application queries return correct results (e.g., "all applications using PostgreSQL", "Terraform version distribution by team").

**FRs covered**: FR-037, FR-038

### Implementation

- [X] T108 [US10] Define enterprise graph record types (GraphPublishRecord with app graph, team, commit SHA, deployment timestamp, IaC versions, module usage) in `pkg/reporadius/types.go`
- [X] T109 [US10] Implement graph publisher interface and HTTP client (POST app graph record to enterprise endpoint) in `pkg/reporadius/publisher.go`
- [X] T110 [US10] Implement graph metadata collector (extract IaC versions, module usage, team/org from constitution and deployment context) in `pkg/reporadius/publisher.go`
- [X] T111 [US10] Add graph publishing step to deployment workflow (after successful deployment, publish to enterprise endpoint) in `pkg/reporadius/templates/radius-deploy.yml.tmpl`
- [X] T112 [US10] Implement `publish_app_graph` MCP tool for manual/on-demand graph publishing in `pkg/mcpserver/tools/publish_app_graph.go`
- [X] T113 [P] [US10] Add unit tests for graph publisher (successful publish, error handling, metadata collection) in `pkg/reporadius/publisher_test.go`
- [X] T114 [P] [US10] Add unit tests for `publish_app_graph` MCP tool in `pkg/mcpserver/tools/publish_app_graph_test.go`

### Checkpoint: US10

Verify: After deployment, application graph data with full metadata is published to the enterprise endpoint. MCP tool supports on-demand publishing. Published records include team, commit SHA, IaC versions, and module usage.

---

## Phase 11c: Skills-Centric Agent Enhancement (Cross-cutting)

**Story goal**: Enhance the Radius Agent Skill to follow Skills-centric agentic patterns (Vercel / database-skills.com style) where the Skill drives actions conversationally, and the application graph is presented visually rather than as raw IaC.

**FRs covered**: FR-036, FR-039

### Implementation

- [X] T115 [US2,US3] Update Agent Skill template to include Skills-centric action directives (rad init, architecture design, deploy, graph generation) in `pkg/agentskill/templates/skill.md.tmpl`
- [X] T116 [US3] Update Agent Skill template to instruct Copilot to present application graphs visually (Mermaid rendering) rather than raw Bicep/YAML in `pkg/agentskill/templates/skill.md.tmpl`
- [X] T117 [US3] Add skill-driven workflow examples to Agent Skill template (conversational rad init, describe-then-deploy, graph-first review patterns) in `pkg/agentskill/templates/skill.md.tmpl`
- [X] T118 [P] Update Agent Skill snapshot tests for new Skills-centric content in `pkg/agentskill/generator_test.go`
- [X] T119 [P] Add unit tests for visual graph presentation directives in Agent Skill output in `pkg/agentskill/generator_test.go`

### Checkpoint: Skills Enhancement

Verify: Generated Agent Skill includes Skills-centric action directives, graph-first visual instructions, and conversational workflow examples. Skill output passes snapshot tests.

---

## Phase 12: Polish & Cross-Cutting Concerns

Goal: End-to-end testing, documentation, performance validation, and operational hardening.

- [X] T088 Add functional test for end-to-end `rad init` workflow (constitution + skill + MCP) in `test/functional/copilot-workflow/init_test.go`
- [X] T089 [P] Add functional test for MCP server tool call round-trips in `test/functional/copilot-workflow/mcp_test.go`
- [X] T090 [P] Add functional test for artifact generation pipeline in `test/functional/copilot-workflow/artifacts_test.go`
- [X] T091 Add CLI reference documentation for `rad init` constitution options (auto-generated from Cobra) in `docs/`
- [X] T092 [P] Add CLI reference documentation for `rad mcp serve` command in `docs/`
- [X] T093 [P] Add reference documentation for `Platform-Engineering-Constitution.md` file format in `docs/`
- [X] T094 [P] Add reference documentation for MCP server tools (input/output schemas) in `docs/`
- [X] T095 Add tutorial guide for end-to-end developer workflow (matches quickstart.md flow) in `docs/`
- [X] T096 Implement MCP server graceful degradation (timeout handling, reconnection, error responses) in `pkg/mcpserver/server.go`
- [X] T097 [P] Add telemetry/observability instrumentation to MCP server tools (request timing, error rates) in `pkg/mcpserver/tools/`
- [X] T098 Performance validation: verify `rad init` completes in <5 min (SC-001) and artifact generation in <2 min (SC-003) as integration benchmarks
- [X] T104 Performance validation: verify architecture design completes in <10 min (SC-002), PR visualization renders in <60 sec (SC-005), and end-to-end workflow in <30 min (SC-006) — document exclusions for metrics dependent on external services (Copilot, GitHub API latency)

---

## Dependencies & Execution Order

### Story Dependency Graph

```
Phase 1: Setup ──────────────────────┐
Phase 2: Foundational ───────────────┤ (blocks ALL stories)
                                     │
Phase 3: US1 (Constitution) ─────────┤
    │                                │
    ├── Phase 4: US2 (Skill + MCP) ──┤ (depends on US1: constitution data → skill content)
    │       │                        │
    │       ├── Phase 5: US3 (Design)┤ (depends on US2: MCP server running)
    │       │       │                │
    │       │       └── Phase 6: US4 (Artifacts) ── (depends on US3: tool implementations)
    │       │               │
    │       │               └── Phase 7: US5 (PR Deploy) ── (depends on US4: artifacts exist)
    │       │                       │
    │       │                       ├── Phase 8: US6 (Graph Viz + README) ── (depends on US5: PR workflow)
    │       │                       ├── Phase 9: US7 (Plan Review) ── (depends on US5: deployment flow)
    │       │                       ├── Phase 11: US9 (Notifications) ── (depends on US5: deployment flow)
    │       │                       └── Phase 11b: US10 (Enterprise Graph) ── (depends on US5: deployment + graph)
    │       │
    │       ├── Phase 10: US8 (Scaffolding) ── (depends on US2: skill template)
    │       └── Phase 11c: Skills Enhancement ── (depends on US2: skill template + US3: design tools)
    │
    └── (no other direct dependents)

Phase 12: Polish ── (after all stories)
```

### Critical Path

```
Setup → Foundational → US1 → US2 → US3 → US4 → US5 → Polish
```

### Parallel Opportunities (per phase)

**Phase 1 (Setup)**:
- T002, T003 can run in parallel (different files, no dependencies)

**Phase 2 (Foundational)**:
- T005, T006, T007 can run in parallel (different packages)
- T008 depends on T004 (interfaces reference constitution types)

**Phase 3 (US1)**:
- T012 parallel with T011 (template file vs. generator logic)
- T020, T021, T022 parallel (test files for different source files)

**Phase 4 (US2)**:
- T032, T033, T034 parallel (test files for different packages)

**Phase 5 (US3)**:
- T044, T045, T046, T047, T102, T103 all parallel (test files for independent tools)

**Phase 6 (US4)**:
- T049 parallel with T048 (template file vs. generator logic)
- T055, T056, T057, T058 all parallel (test files for independent generators)

**Phase 7 (US5)**:
- T061 parallel with T060 (template file vs. generator logic)
- T064, T065 parallel (test files for different source files)

**Phase 8 (US6)**:
- T071, T072 parallel (test files for different functions)

**Phase 9 (US7)**:
- T077, T078 parallel (test files for different source files)

**Phase 12 (Polish)**:
- T089, T090 parallel (functional tests for independent subsystems)
- T092, T093, T094 parallel (independent documentation pages)

### Cross-Story Parallelism

Once US5 (Phase 7) is complete:
- **US6** (Phase 8), **US7** (Phase 9), **US9** (Phase 11), and **US10** (Phase 11b) can proceed **in parallel** — they enhance the PR/deployment workflow independently.

Once US2 (Phase 4) is complete:
- **US8** (Phase 10) can proceed **in parallel** with US3–US5 — it extends the skill template independently.
- **Skills Enhancement** (Phase 11c) can proceed **in parallel** with US3–US5 — it updates the skill template with Skills-centric patterns.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL — blocks all stories)
3. Complete Phase 3: User Story 1 (Constitution Generation)
4. **STOP and VALIDATE**: Run `rad init`, verify constitution is generated and valid
5. Demo/review if ready — this alone delivers value (FR-001, FR-002, FR-003, FR-007, FR-008)

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. **US1** → Constitution generation → Test independently → **MVP!**
3. **US2** → Agent Skill + MCP init → Test independently → Developer can start Copilot sessions
4. **US3** → Architecture design tools → Test independently → Copilot can analyze and propose
5. **US4** → Artifact generation → Test independently → Full design-to-artifacts pipeline
6. **US5** → PR deployment → Test independently → Complete end-to-end workflow
7. **US6 + US7 + US9 + US10** (parallel) → PR enhancements, enterprise graph → Polish the deployment experience
8. **US8 + Skills Enhancement** (parallel) → Scaffolding compliance + Skills-centric agent → Additional value
9. Polish → Docs, performance, observability

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - **Developer A**: US1 (Constitution) → US2 (Skill + MCP init)
   - **Developer B**: Reviews US1, then starts US3 (Design tools) once US2 is ready
3. Once US2 is done:
   - **Developer A**: US4 (Artifact generation)
   - **Developer B**: US3 (Architecture design tools)
   - **Developer C**: US8 (Scaffolding — only depends on US2)
4. Once US5 is done:
   - **Developer A**: US6 (Graph visualization)
   - **Developer B**: US7 (Deployment plan review)
   - **Developer C**: US9 (Notifications)

---

## Notes

- `[P]` tasks = different files, no dependencies on incomplete tasks within the same phase
- `[US#]` label maps task to a specific user story for traceability
- Task IDs within a phase are grouped by concern (implementation then tests) for readability; per constitution Principle IV (Red-Green-Refactor), tests SHOULD be written before or during implementation execution
- Each user story should be independently completable and testable at its checkpoint
- Commit after each task or logical group of tasks
- Stop at any checkpoint to validate the story independently before proceeding
- All file paths are relative to the repository root (`/radius/`)
- Testing follows existing patterns: `testify` + `gomock`, table-driven tests, `radcli.SharedCommandValidation()` for CLI commands
- MCP server tests use JSON-RPC test client harness per research decision R7
- Agent Skill tests use snapshot/golden file comparison per research decision R7
