# Implementation Plan: Developer Workflow in Copilot and GitHub Using Radius

**Branch**: `004-copilot-radius-dev-workflow` | **Date**: 2026-02-20 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/004-copilot-radius-dev-workflow/spec.md`

## Summary

This feature extends the Radius CLI and introduces two new integration surfaces — a Copilot Agent Skill and a local MCP server — to enable a complete developer workflow from application description to deployment via GitHub pull requests. The `rad init` command is extended to generate a Platform Engineering Constitution, produce an Agent Skill grounded on the constitution, and start an MCP server exposing deterministic tools for resource type resolution, recipe search, artifact generation, and constitution validation. Architecture design sessions in Copilot leverage the Agent Skill for always-on guidance and the MCP server for tool calls, producing application definitions, resource types, recipes, and environment configurations. A PR-based deployment flow through Repo Radius completes the workflow with architecture graph visualizations and deployment plan reviews.

**Visual-First, Skills-Centric Direction** (per Mario Rodriguez / GitHub review): The primary developer experience is through rich visual application graphs, not raw IaC files. `app.bicep` and other deployment files are generated artifacts that exist behind the scenes — the agent and Radius Skills present them visually. Radius Skills are the central mechanism for agentic workloads in Copilot, following patterns pioneered by Vercel and database-skills.com. The application graph appears in every PR comment and the repository README. At enterprise scale, graph data correlated with commit history and IaC deployment metadata (Terraform/Bicep versions, module usage) enables organizational intelligence across teams.

## Technical Context

**Language/Version**: Go (version per go.mod, currently 1.23+) for CLI, MCP server, and control plane; TypeScript/Node.js for Bicep tooling; Bicep for artifact generation  
**Primary Dependencies**: Radius CLI (`cmd/rad`), Radius control plane (`pkg/corerp`, `pkg/ucp`), TypeSpec (API definitions), GitHub Copilot (Agent Skills, MCP protocol), GitHub Actions/Apps (PR integration)  
**Storage**: File system (generated artifacts in `radius/` directory), Git repository (state via PRs), Radius control plane datastore (deployment state)  
**Testing**: `go test` via `make test` (unit), `test/functional/` (functional), MCP protocol tests (JSON-RPC test client), Agent Skill snapshot tests (golden files) — see [research.md](research.md#r7-testing-agent-skills-and-mcp-protocol-compliance)  
**Target Platform**: Developer workstation (macOS, Linux, Windows) for CLI and MCP server; GitHub.com for PR integration and Repo Radius deployment  
**Project Type**: Multi-component (CLI extension + MCP server + Agent Skill + GitHub integration + artifact generators)  
**Performance Goals**: `rad init` completes in <5 min (SC-001), architecture design in <10 min (SC-002), artifact generation in <2 min (SC-003), PR visualization in <60 sec (SC-005), end-to-end in <30 min (SC-006)  
**Constraints**: MCP server runs locally (no hosted infra for dev workflow), deterministic tool responses (same inputs → same outputs), all state in repository (no out-of-band state)  
**Scale/Scope**: Target multi-service apps (2+ services, 2+ infra dependencies), single developer workflow initially, organizational scale via constitution sharing

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Status | Notes |
|---|-----------|--------|-------|
| I | API-First Design | PASS | MCP server tools are API-first with schemas defined in OpenAPI 3.1 (`contracts/mcp-server.yaml`). MCP tools use JSON-RPC 2.0 (not ARM-RPC), so handwritten OpenAPI is appropriate rather than TypeSpec generation. Agent Skill is a declarative API surface. All new Radius resource types follow ARM-RPC patterns with TypeSpec. |
| II | Idiomatic Code Standards | PASS | Go for CLI extension, MCP server, and control plane changes (follows Effective Go, `gofmt`, godoc). TypeScript/Node.js for Bicep tooling (follows TypeScript handbook). Bicep for artifact templates (follows official best practices). Markdown for constitution and Agent Skill. |
| III | Multi-Cloud Neutrality | PASS | Constitution captures multi-cloud preferences. Architecture design maps to cloud-agnostic resource types. Recipes abstract cloud-specific provisioning. MCP server resolves types across providers. No cloud-specific assumptions in the core workflow. |
| IV | Testing Pyramid Discipline | PASS with NEEDS CLARIFICATION | Go unit tests for CLI extension, constitution parser, artifact generator. Functional tests for `rad init` workflow. MCP server integration tests. **NEEDS CLARIFICATION**: No existing test pattern for Agent Skill validation or MCP protocol compliance testing — research required. |
| V | Collaboration-Centric Design | PASS | Constitution bridges platform engineer governance (policies, module catalogs) with developer experience (Copilot conversation, progressive disclosure). Spec explicitly addresses both audiences. |
| VI | Open Source and Community-First | PASS | Spec and plan stored in public design-notes flow. Feature follows issue-first workflow. All design decisions documented. |
| VII | Simplicity Over Cleverness | PASS | Two integration surfaces (Agent Skill, MCP server) each have a single responsibility and can operate independently. `rad init` consolidates setup into one command. Progressive disclosure hides internals by default. |
| VIII | Separation of Concerns | PASS | Clear separation: CLI handles init, MCP server handles tool calls, Agent Skill handles guidance, Repo Radius handles deployment. No circular dependencies. Domain logic (constitution parsing, resource resolution) separated from transport (MCP protocol, CLI UI). |
| IX | Incremental Adoption | PASS | `rad init` is opt-in (run explicitly). Agent Skill auto-activates but doesn't break existing workflows. MCP server is local and optional. Constitution supports gradual policy adoption. Existing `rad` commands unaffected. |
| X | TypeScript/React (Dashboard) | N/A | No dashboard changes in initial scope. Future P2/P3 stories may add graph visualization to dashboard. |
| XI | Frontend Testing (Dashboard) | N/A | No dashboard changes in initial scope. |
| XII | Resource Type Schema Quality | PASS | Generated resource types follow YAML schema standards with property definitions, descriptions, and examples. MCP server validates schemas against catalog. |
| XIII | Recipe Development Standards | PASS | Generated recipes follow Terraform/Bicep module structure with README, parameters with descriptions, secure defaults. References approved modules from constitution. |
| XIV | Documentation Structure | PASS | Feature requires docs updates following Diátaxis: tutorial (quickstart), how-to (init guide, architecture design), reference (MCP tools, constitution format), explanation (workflow overview). |
| XV | Documentation Contribution | PASS | Docs will follow Hugo build, markdownlint, pyspelling. CLI docs auto-generated from Cobra. MCP tool docs generated from TypeSpec. |
| XVI | Repository-Specific Standards | PASS | Follows existing radius repo conventions (Makefile targets, `pkg/` structure, `cmd/` entrypoints). |
| XVII | Polyglot Project Coherence | PASS | MCP server uses same API patterns as control plane. Constitution format consistent across repos. Error handling follows established Go patterns. Observability via existing telemetry infrastructure. |

**Gate Result**: PASS — No blocking violations. One NEEDS CLARIFICATION item (Testing: Agent Skill and MCP protocol test patterns) to be resolved in Phase 0 research.

### Post-Design Re-evaluation (Phase 1 complete)

The NEEDS CLARIFICATION item from principle IV (Testing Pyramid Discipline) has been **resolved** in [research.md](research.md#r7-testing-agent-skills-and-mcp-protocol-compliance):
- **MCP tools**: Standard Go unit tests with `gomock` (follows existing `radinit` test patterns)
- **MCP protocol**: JSON-RPC test client harness (follows existing HTTP handler test patterns in `pkg/`)
- **Agent Skill**: Snapshot/golden file tests comparing generated output to expected Markdown
- **End-to-end**: Functional tests in `test/functional/copilot-workflow/`

All 17 constitution principles now **PASS** with no open clarifications.

**Post-Design Gate Result**: PASS — All constitution principles satisfied. No violations requiring justification.

## Project Structure

### Documentation (this feature)

```text
specs/004-copilot-radius-dev-workflow/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   ├── mcp-server.yaml  # MCP server tool definitions (OpenAPI)
│   └── constitution-schema.yaml  # Constitution file JSON Schema
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Radius CLI extension (Go)
pkg/
├── cli/
│   └── cmd/radinit/     # Extended rad init command (constitution, skill, MCP setup)
├── constitution/        # Constitution parsing, validation, generation
│   ├── parser.go
│   ├── generator.go
│   ├── validator.go
│   └── types.go
├── agentskill/          # Agent Skill generation from constitution
│   ├── generator.go
│   └── templates/
├── mcpserver/           # MCP server implementation
│   ├── server.go        # MCP protocol handler
│   ├── tools/           # Tool implementations
│   │   ├── resolve_resource_types.go
│   │   ├── search_recipes.go
│   │   ├── generate_app_definition.go
│   │   ├── validate_constitution.go
│   │   └── get_environments.go
│   └── config.go        # Server configuration and registration
├── artifacts/           # Deployment artifact generation
│   ├── appdef.go        # Application definition generator (app.bicep)
│   ├── resourcetype.go  # Resource type resolver and generator
│   ├── recipe.go        # Recipe configuration generator
│   ├── environment.go   # Environment configuration generator
│   └── graph.go         # Application graph generator (app-graph.json)
└── reporadius/          # Repo Radius deployment integration
    ├── workflow.go       # GitHub Actions workflow generator
    ├── visualization.go  # Architecture graph visualization for PRs
    ├── deployment.go     # Deployment plan generation
    ├── pullrequest.go    # PR creation and branch management
    ├── readme.go         # README app graph embedding and updating
    └── publisher.go      # Enterprise graph database publisher

# Generated output directory (in user's repo)
radius/
├── app.bicep            # Application definition
├── app-graph.json       # Application graph for PR visualization
├── types/               # Resource type definitions
│   └── *.yaml
├── recipes/             # Recipe configurations
│   └── */main.tf or main.bicep
└── environments/        # Environment configurations
    └── *.bicep

# Agent Skill (in user's repo)
.github/skills/radius-platform/
└── SKILL.md

# MCP server config (in user's repo)
.github/copilot/
└── mcp.json

# GitHub Actions (in user's repo)
.github/workflows/
└── radius-deploy.yml

# Tests
pkg/constitution/        # Unit tests alongside source
pkg/mcpserver/           # Unit tests alongside source
pkg/artifacts/           # Unit tests alongside source
test/functional/
└── copilot-workflow/    # Functional tests for end-to-end workflow
```

**Structure Decision**: Multi-component extension to the existing radius repository structure. New packages under `pkg/` follow the established Go package layout. The MCP server is implemented in Go to share types and logic with the CLI. Generated artifacts follow a convention-based `radius/` directory in the user's repository. No new top-level directories are introduced; all code integrates into existing `pkg/` and `cmd/` hierarchies.

## Complexity Tracking

> No constitution violations requiring justification. The feature introduces multiple new packages but each has a single responsibility and integrates into the existing repository structure without new abstraction layers.

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| MCP server in Go (not TypeScript) | Go implementation | Shares types and logic with CLI and control plane; avoids polyglot build complexity for a single binary distribution |
| Two integration surfaces | Agent Skill + MCP server | Each addresses a distinct concern (guidance + interactive design vs. deterministic tooling); constitution principle VIII (Separation of Concerns) supports this |
| Constitution as Markdown | Human-readable Markdown with structured sections | Aligns with constitution principle VII (Simplicity); Markdown is familiar to developers and renderable in GitHub |
| Visual-first over file-first | Application graph drives UX, IaC files are behind-the-scenes | Per Mario Rodriguez / GitHub review: developers interact with the visual graph. app.bicep is generated transparently. "Go fast to a future that works with UI instead of app.bicep." |
| Skills-centric agentic architecture | Radius Skill is the primary entry point | Skills drive rad init, architecture design, deployment — follows Vercel and database-skills.com patterns for agentic workloads |
| Application graph in every PR + README | Graph is the universal artifact | Every PR gets a graph visualization comment; README embeds the current graph. The graph is the unit of visibility, not Bicep |
| Enterprise graph database publishing | Aggregate app graphs across org | Correlate app graph + commit data + IaC deployment data to enable org-level intelligence (e.g., Terraform version usage per team, dependency impact analysis) |
