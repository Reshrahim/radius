# Research: Developer Workflow in Copilot and GitHub Using Radius

**Branch**: `004-copilot-radius-dev-workflow` | **Date**: 2026-02-20

## R1: MCP Server Implementation in Go

**Decision**: Implement the Radius MCP server in Go as part of the `rad` binary, using stdio transport.

**Rationale**: The MCP server needs access to constitution parsing, resource type resolution, recipe catalogs, and application definition generation — all of which are Go types in `pkg/`. Implementing in Go eliminates cross-language serialization and allows the server to be distributed as part of the single `rad` binary. The MCP protocol (JSON-RPC 2.0 over stdio) is language-agnostic; Go libraries exist for JSON-RPC serving. The existing `.vscode/mcp.json` in the repo shows the pattern for registering MCP servers with Copilot — the Radius server will follow the same configuration pattern, writing to `.github/copilot/mcp.json` in the user's repo.

**Alternatives considered**:
- **TypeScript MCP server**: Would require a separate build pipeline, Node.js runtime dependency, and type duplication between Go and TypeScript. Rejected because it increases complexity without clear benefit.
- **Hosted MCP server**: Would require infrastructure provisioning for the dev workflow. Rejected per spec constraint that MCP server runs locally with no hosted infra.

**Implementation pattern**: `rad mcp serve` subcommand starts the stdio-based MCP server. `rad init` writes the `mcp.json` config pointing to `rad mcp serve` as the command. Copilot spawns the process on demand.

## R2: Agent Skill Generation Pattern

**Decision**: Generate the Agent Skill as a Markdown file (`.github/skills/radius-platform/SKILL.md`) from a Go template, parameterized by constitution content.

**Rationale**: No existing Agent Skill patterns exist in the Radius workspace. Agent Skills are Markdown files with structured instructions that Copilot loads as context. The skill file must be grounded on the constitution — meaning its content varies per organization. Using a Go `text/template` with constitution fields as template variables follows the existing pattern in `pkg/cli/setup/application.go` where `app.bicep`, `rad.yaml`, and `bicepconfig.json` are generated from Go string templates.

**Alternatives considered**:
- **Static skill file**: Would not adapt to organization-specific policies. Rejected because constitution-aware guidance is a core requirement (FR-004).
- **Dynamic skill via MCP tool**: Agent Skills are loaded at session start, not dynamically generated per query. The MCP server provides dynamic tooling; the skill provides static context. Both are needed.

**Key content sections for the generated skill**:
1. Organization context (from constitution: name, cloud providers, compute targets)
2. Platform policies summary (naming conventions, security rules, approved modules)
3. Radius workflow instructions (how to invoke architecture design, artifact generation)
4. MCP tool descriptions (what tools are available, when to use each)
5. Constitution compliance reminders (nudge toward validated patterns)

## R3: Constitution File Format and Schema

**Decision**: Use structured Markdown with YAML frontmatter for the constitution file, plus a JSON Schema for programmatic validation.

**Rationale**: The design-notes repo (`design-notes/features/2026-02-application-assembly.md`) contains a draft constitution format using Markdown with structured sections (Terraform Registry, Naming Convention, Labels, Authentication Preferences, Security Policies). The spec requires the file to be "structured, human-readable" (FR-002). Markdown with YAML frontmatter provides both human readability and machine parseability. The YAML frontmatter contains structured data (cloud providers, module catalog URLs, policy rules); the Markdown body provides human-readable documentation and rationale.

**Alternatives considered**:
- **Pure YAML/JSON**: Machine-readable but poor developer experience for reading and editing. Rejected per constitution principle VII (Simplicity).
- **Pure Markdown without schema**: Readable but not reliably parseable. Rejected because FR-012 requires programmatic validation of policies against architecture designs.
- **TOML**: Less common in the Kubernetes/cloud-native ecosystem. Rejected for familiarity reasons.

**Schema design**:
- YAML frontmatter keys: `organization`, `version`, `cloudProviders[]`, `computeTargets[]`, `iacTooling`, `moduleCatalog`, `containerRegistry`, `policies[]`
- Each `policy` has: `name`, `source` (manual/imported), `rule`, `severity` (error/warning), `remediation`
- JSON Schema published at `contracts/constitution-schema.yaml` for validation tooling

## R4: Extending `rad init` — Architecture and Approach

**Decision**: Extend the existing `rad init` command by adding new phases to the `enterInitOptions()` chain, following the established multi-file command pattern.

**Rationale**: The current `rad init` implementation at `pkg/cli/cmd/radinit/` is well-structured with a clear extension pattern: each concern (cluster, environment, cloud providers, application, recipes) is a separate file with its own prompt functions. The command chains them in `enterInitOptions()`. Constitution generation, Agent Skill generation, and MCP server setup are three new phases that fit naturally into this chain. The existing `--full` flag pattern supports toggling extended prompts.

**Implementation approach**:
1. New file: `constitution.go` — Prompts for organization info, cloud preferences, policy imports, module catalog. Generates `Platform-Engineering-Constitution.md`.
2. New file: `agentskill.go` — Generates `.github/skills/radius-platform/SKILL.md` from constitution data.
3. New file: `mcpserver.go` — Writes `.github/copilot/mcp.json` configuration. Optionally starts the MCP server.
4. Chain order: constitution → agent skill → MCP server → (existing) cluster → environment → cloud providers → application → recipes.
5. New flag: `--non-interactive` — Accepts constitution config from a file for automation (FR-033).

**Testing approach**: Follow existing test patterns — `constitution_test.go`, `agentskill_test.go`, `mcpserver_test.go` with table-driven `ValidateInput` tests and `gomock` mocks for file system operations, prompter, and MCP server lifecycle.

## R5: Resource Type Resolution Strategy

**Decision**: Three-tier resolution: local repository → community catalog (resource-types-contrib) → generated custom types.

**Rationale**: FR-015 requires checking in order: existing types in the repo, contributed types from the community catalog, and generating custom types only when no match exists. The `resource-types-contrib` repo at `../resource-types-contrib` contains community-contributed resource type schemas organized by category (Compute/containers, Data/postgreSqlDatabases, etc.). Each type has a `*.yaml` manifest file. The `bicep-tools/` module already parses these YAML manifests and converts them to Bicep types. The resolution engine needs to:
1. Scan the local `radius/types/` directory for existing type definitions.
2. Query the community catalog (via bundled index or HTTP) for matching types by dependency name.
3. Fall back to generating a minimal custom type from the detected dependency pattern.

**Alternatives considered**:
- **Online-only catalog lookup**: Would fail offline. Rejected per constraint that MCP server runs locally.
- **Pre-bundled full catalog**: Would bloat the binary. Rejected in favor of a lightweight index with on-demand fetch.

**Implementation pattern**: `pkg/artifacts/resourcetype.go` with a `Resolver` interface: `Resolve(dependency string) (ResourceType, Source, error)` where `Source` is one of `Local`, `Community`, `Generated`.

## R6: Artifact Generation Pipeline

**Decision**: Template-based generation for application definition, environment configuration, and recipes, with the existing `bicep-tools` pipeline for resource type compilation.

**Rationale**: The existing `pkg/cli/setup/application.go` uses Go string templates for `app.bicep` generation. The `bicep-tools/` module handles YAML manifest → Bicep type conversion. For the new workflow, we need: (1) rich `app.bicep` generation with multiple services, connections, and secret bindings — extending the existing template approach; (2) environment Bicep files from constitution preferences; (3) recipe configurations that wrap approved IaC modules; (4) an `app-graph.json` for PR visualization.

**Pipeline stages**:
1. **Input**: Confirmed architecture design (services, dependencies, connections, cloud mappings)
2. **Resource type resolution**: Via R5 resolver
3. **Recipe matching**: Map each dependency to an approved module from the constitution's catalog
4. **Template rendering**: Generate `app.bicep`, environment files, recipe configs
5. **Graph generation**: Build `app-graph.json` from the dependency graph
6. **Validation**: Check all artifacts against constitution policies

## R7: Testing Agent Skills and MCP Protocol Compliance

**Decision**: Unit test the MCP server tools individually (standard Go unit tests), integration test the MCP protocol via a test client harness, and validate Agent Skill content via schema and snapshot tests.

**Rationale**: This was the "NEEDS CLARIFICATION" item from the Constitution Check. No existing test pattern for Agent Skills or MCP servers exists in the codebase. However:
- **MCP tools** are Go functions with defined inputs/outputs — standard unit testing applies.
- **MCP protocol compliance** can be tested with a JSON-RPC test client that sends well-formed and malformed requests to the server and validates responses. This follows the pattern of existing HTTP handler tests in `pkg/`.
- **Agent Skill content** is generated Markdown. Snapshot tests (comparing generated output to golden files) validate that constitution changes produce expected skill updates. The `testify` assertions library supports this.

**Test structure**:
- `pkg/mcpserver/tools/*_test.go` — Unit tests for each tool function
- `pkg/mcpserver/server_test.go` — Protocol compliance tests (JSON-RPC request/response pairs)
- `pkg/agentskill/generator_test.go` — Snapshot tests for skill generation from sample constitutions
- `test/functional/copilot-workflow/` — End-to-end test: `rad init` → verify generated files → start MCP server → send tool requests → validate responses

## R8: Repo Radius and PR Deployment Architecture

**Decision**: Implement Repo Radius as a GitHub Actions workflow generated by `rad init`, with a companion GitHub App for PR visualization and deployment approval.

**Rationale**: The existing Radius codebase has 32 GitHub Actions workflows and detailed workflow design guidelines (`design-notes/tools/2025-03-workflow-changes.md`). A generated `.github/workflows/radius-deploy.yml` follows established patterns. The GitOps design (`design-notes/tools/2025-01-gitops-technical-design.md`) uses a Flux Controller approach — Repo Radius extends this with PR-native integration. The PR visualization (architecture graph as a comment) and deployment approval (environment protection rules) use GitHub-native features.

**Alternatives considered**:
- **Flux Controller only**: Covers GitOps but lacks PR-native visualization and approval. Rejected because the spec requires in-PR architecture graphs and deployment plan review.
- **GitHub App only**: Would handle everything but is harder to customize. Rejected in favor of a hybrid: Actions for workflow execution, App (optional) for rich PR comments.
- **ArgoCD integration**: External dependency, not GitHub-native. Rejected per simplicity principle.

**Implementation approach**:
1. `pkg/reporadius/workflow.go` — Generates a GitHub Actions workflow file that triggers on PR merge, runs `rad deploy`, and posts status comments.
2. `pkg/reporadius/visualization.go` — Generates Mermaid diagrams from `app-graph.json` for PR comments. Mermaid renders natively in GitHub Markdown.
3. `pkg/reporadius/deployment.go` — Generates deployment plans by diffing current vs. proposed `app-graph.json` states.

## R9: Application Discovery and Source Code Analysis

**Decision**: Leverage the existing `rad app discover` design from the application assembly feature, extending it with Copilot-interactive refinement.

**Rationale**: The design note `design-notes/features/2026-02-application-assembly.md` defines a comprehensive approach: scan repository for service manifests (Dockerfiles, docker-compose, Kubernetes YAML, package.json), detect service dependencies (connection strings, environment variables, import statements), and produce a structured application model. This is exactly what US-3 (Architecture Design) needs. The difference is that in the Copilot workflow, the developer confirms and refines the discovery results conversationally rather than just accepting CLI output.

**Implementation**: The MCP server tool `resolve_resource_types` wraps the discovery logic. The Agent Skill instructs Copilot to present discovery results conversationally and accept feedback. The MCP server's `generate_app_definition` tool accepts the confirmed architecture as input.

## R10: Architecture Visualization Format

**Decision**: Use Mermaid diagrams for architecture graph visualization in PRs, with `app-graph.json` as the structured data source.

**Rationale**: The spec assumption states "The architecture graph visualization in PRs uses a format that renders natively in the GitHub PR UI (e.g., Mermaid diagrams in Markdown)." Mermaid renders natively in GitHub Markdown (PR comments, issue bodies, README). The `app-graph.json` file stores the full graph data; the Mermaid diagram is generated from it for display purposes.

**Alternatives considered**:
- **SVG image generation**: Requires a rendering service or local tool. More complex, harder to diff.
- **ASCII art**: Limited expressiveness for complex architectures.
- **GitHub App with custom rendering**: Higher implementation cost, separate deployment.

**Graph structure in `app-graph.json`**:
```json
{
  "version": "1.0",
  "application": { "name": "myapp" },
  "nodes": [
    { "id": "frontend", "type": "Applications.Core/containers", "properties": {} },
    { "id": "redis", "type": "Applications.Datastores/redisCaches", "recipe": "azure-redis", "properties": {} }
  ],
  "edges": [
    { "from": "frontend", "to": "redis", "type": "connection" }
  ]
}
```
