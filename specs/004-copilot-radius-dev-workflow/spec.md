# Feature Specification: Developer Workflow in Copilot and GitHub Using Radius

**Feature Branch**: `004-copilot-radius-dev-workflow`  
**Created**: 2026-02-20  
**Status**: Draft  
**Input**: User description: "Developer workflow in Copilot and GitHub using Radius for cloud-native application design, build, and deployment — covering initial Radius setup, application development and packaging with Copilot, cloud architecture design using Radius Agent Skill and MCP server, and deployment via Repo Radius through GitHub pull requests"

## Principles

1. **Platform engineering as code**: Organizational policies, guardrails, and infrastructure standards are captured in a machine-readable constitution file that grounds every subsequent decision — from repository creation to deployment approval. The constitution is the single source of truth for compliance.

2. **AI-assisted, human-approved**: Copilot and Radius collaborate to suggest architecture patterns, generate artifacts, and propose deployments, but the developer retains approval authority at every critical gate (topology confirmation, PR review, deployment approval). Automation accelerates; humans decide.

3. **Progressive disclosure of complexity**: The default workflow hides Radius internals (resource types, recipes, Bicep) behind Copilot's conversational interface. Developers who want deeper control can inspect and edit all generated artifacts directly. The path from "describe my app" to "deployed and running" is minutes, not hours.

4. **Composable integration**: The workflow is built from two complementary integration surfaces — Agent Skill (always-on guidance and interactive design sessions) and MCP server (deterministic tool calls). Each can operate independently but together they deliver the full experience.

5. **Repository-centric deployment**: All state, configuration, and deployment artifacts live in the repository. Deployment is triggered by pull request merge, reviewed through architecture visualizations in the PR, and approved through standard GitHub workflows. There is no out-of-band state.

6. **Visual-first, code-behind**: The primary developer experience is through rich visual interfaces — interactive application graphs, architecture diagrams, and deployment plans — not raw IaC files. Files like `app.bicep` are generated artifacts that exist behind the scenes; the agent and skills present them visually. Developers interact with the graph, not the code.

7. **Skills-centric agentic architecture**: Radius Skills are the central mechanism for agentic workloads in Copilot. Skills drive actions like `rad init`, architecture design, and deployment through conversational interfaces, following patterns pioneered by Vercel and database-skills.com. The Skill is the primary entry point; the MCP server provides the deterministic backend.

8. **Application graph as the universal artifact**: The application graph is the core data structure that flows through every stage of the workflow — from initial design through PR review to deployment. It appears in every PR comment, lives in the README, feeds enterprise-level visibility dashboards, and when correlated with commit data and IaC deployment metadata, enables organizational intelligence (e.g., tracking Terraform version usage per team).

## User Journey

This section describes the end-to-end developer experience using the Radius platform skill in GitHub Copilot CLI.

### Persona: Emma, Enterprise Application Developer

Emma is building a new microservices application for her company. She has experience with containers and cloud services but wants to ensure her infrastructure follows her organization's platform standards from day one.

### Journey Steps

#### 1. Project Initialization (5 minutes)

Emma creates a new directory for her project and opens her terminal.

```bash
cd ~/projects
mkdir inventory-service
cd inventory-service
git init
gh copilot
```

In Copilot CLI, she invokes the Radius platform skill:

```
> @radius-platform init
```

The skill guides her through a friendly conversation:

```
Welcome to Radius Platform Engineering Setup!

I'll help you set up your project with organizational standards.

What is your organization name? 
> contoso

Which cloud providers will you use? (azure, aws, both)
> azure

Azure Subscription ID?
> xxxx-yyyy-zzzz

Which compute platform? (kubernetes, aca, both)
> kubernetes

Current kubectl context?
> prod-cluster

IaC tooling preference? (terraform, bicep)
> terraform

Terraform module catalog URL?
> https://registry.terraform.io/namespaces/contoso

Import existing policies? (azure-policy, aws-config, manual, skip)
> azure-policy

Azure Policy assignment path?
> /subscriptions/xxxx-yyyy-zzzz/providers/Microsoft.Authorization/policyAssignments/contoso-standards

✓ Analyzing policies...
✓ Generating constitution...
✓ Creating Radius platform skill...
✓ Starting MCP server...
✓ Registering with Copilot...

Done! Your platform engineering foundation is ready.
```

Emma's repository now contains:
- `Platform-Engineering-Constitution.md` — her organization's policies
- `.github/skills/radius-platform/SKILL.md` — the platform skill definition
- `.github/copilot/mcp.json` — MCP server configuration

#### 2. Application Development (30 minutes)

Emma writes her application code — a Go-based inventory API and a React frontend. She structures her repository:

```
inventory-service/
├── Platform-Engineering-Constitution.md
├── .github/
├── backend/
│   ├── main.go
│   ├── handlers/
│   └── Dockerfile
└── frontend/
    ├── src/
    ├── package.json
    └── Dockerfile
```

Her backend code connects to PostgreSQL and Redis, but she hasn't provisioned any infrastructure yet.

#### 3. Architecture Design (10 minutes)

Emma returns to Copilot CLI and asks the Radius platform skill to design her cloud architecture:

```
> @radius-platform design
```

The skill analyzes her code:

```
Scanning application code...
✓ Found 2 services: backend (Go API), frontend (React)
✓ Detected dependencies: PostgreSQL, Redis
✓ Identified connections: frontend → backend → database/cache

Proposed Architecture:

┌─────────────┐
│  Frontend   │
│   (React)   │
│  Container  │
└──────┬──────┘
       │ HTTP
       v
┌─────────────┐      ┌──────────────┐
│   Backend   │─────>│  PostgreSQL  │
│  (Go API)   │      │ (Azure DB)   │
│  Container  │      └──────────────┘
└──────┬──────┘
       │
       v
┌──────────────┐
│    Redis     │
│(Azure Cache) │
└──────────────┘

Architecture Pattern: Multi-service with managed backing services
Target Compute: Kubernetes (prod-cluster)
Target Cloud: Azure (contoso standards)

Compliance Check: Running validation...
✓ naming-convention: Resources follow contoso-{env}-{service} pattern
✓ encryption-at-rest: All data stores encrypted
✓ network-isolation: VNet integration enabled
✓ backup-retention: 30-day retention configured
✓ monitoring: Application Insights enabled
✓ approved-modules: Using contoso catalog modules
✓ cost-tags: Required tags (cost-center, environment, owner) present
✓ region-restriction: Using approved region (eastus)

Compliance: 8/8 policies pass ✓

Does this look correct? (yes/no/refine)
```

Emma reviews the architecture. She notices Redis is on Azure but wants to ensure it's using a specific cache size:

```
> @radius-platform design --refine "use Premium tier for Redis with 6GB cache"
```

The skill updates the design:

```
Updated Architecture:
  - Redis: Azure Cache for Redis, Premium P1 (6GB)
  
Re-validating compliance...
✓ All policies pass

Design confirmed! Ready to generate deployment artifacts.
```

#### 4. Artifact Generation (3 minutes)

Emma generates the deployment artifacts:

```
> @radius-platform generate
```

The skill creates all necessary files:

```
Generating deployment artifacts...

✓ Resolved resource types:
  - Radius.Compute/containers (community catalog)
  - Radius.Data/postgreSqlDatabases (community catalog)
  - Radius.Data/redisCaches (community catalog)

✓ Generated recipes:
  - azure-postgres (using contoso/azure-postgres v2.1.0)
  - azure-redis (using contoso/azure-redis v1.8.0)

✓ Created application definition:
  - radius/app.bicep

✓ Created environment configuration:
  - radius/environments/dev.bicep
  - radius/environments/prod.bicep

✓ Generated application graph:
  - radius/app-graph.json

Validating against constitution...
✓ All artifacts comply with platform standards

Files created in ./radius directory
```

Emma's repository now has a complete `radius/` directory with all deployment artifacts.

#### 5. Deployment via Pull Request (15 minutes)

Emma is ready to deploy. She invokes:

```
> @radius-platform deploy
```

The skill creates a pull request:

```
Creating deployment pull request...

✓ Committed artifacts to branch: radius-deploy-20260225
✓ Created PR #123: "Deploy inventory-service infrastructure"
✓ Added architecture visualization to PR
✓ Added compliance report to PR

PR URL: https://github.com/contoso/inventory-service/pull/123

Next steps:
1. Review the PR and architecture visualization
2. Get approval from your team
3. Merge the PR to trigger deployment
```

Emma opens the PR in her browser. She sees:

**PR Comment (auto-generated):**

```
## 🎨 Architecture Changes

### Proposed Architecture

[Mermaid diagram showing the complete topology]

### Component Summary

| Component | Type | Cloud Service | Module |
|-----------|------|---------------|--------|
| frontend | Container | AKS (Kubernetes) | native |
| backend | Container | AKS (Kubernetes) | native |
| database | PostgreSQL | Azure Database for PostgreSQL | contoso/azure-postgres v2.1.0 |
| cache | Redis | Azure Cache for Redis (Premium P1, 6GB) | contoso/azure-redis v1.8.0 |

### Compliance Report

✓ 8/8 policies pass

<details>
<summary>View detailed compliance results</summary>

- ✓ naming-convention: All resources follow contoso-{env}-{service} pattern
- ✓ encryption-at-rest: All data stores encrypted with customer-managed keys
- ✓ network-isolation: VNet integration with private endpoints
- ✓ backup-retention: 30-day automated backup retention
- ✓ monitoring: Application Insights configured for all services
- ✓ approved-modules: All recipes use approved contoso catalog modules
- ✓ cost-tags: Required tags present (cost-center: engineering, environment: prod, owner: emma)
- ✓ region-restriction: Using approved region (eastus)

</details>

### Estimated Monthly Cost

- PostgreSQL (General Purpose, 2 vCores): ~$150
- Redis (Premium P1, 6GB): ~$280
- AKS (existing cluster): $0
- Networking: ~$20
- **Total: ~$450/month**

---

✅ Ready for review and deployment
```

Emma's teammate reviews and approves the PR. Emma merges it.

#### 6. Deployment Execution (10 minutes)

After the merge, Repo Radius automatically kicks off:

```
GitHub Actions Workflow: Radius Deployment

✓ Compiling Radius application definition...
✓ Validating environment configuration...
✓ Generating deployment plan...

====================================
Deployment Plan: inventory-service
====================================

Resources to CREATE:
  + PostgreSQL Database (contoso-prod-inventory-db)
    - Region: eastus
    - SKU: General Purpose, 2 vCores
    - Storage: 100GB, encrypted
    - Backup retention: 30 days
    
  + Redis Cache (contoso-prod-inventory-cache)
    - Region: eastus
    - SKU: Premium P1, 6GB
    - Persistence: RDB enabled
    
  + Kubernetes Deployment: backend
    - Replicas: 3
    - Image: contoso.azurecr.io/inventory-backend:sha-a1b2c3d
    - Environment variables: 12 (including DB/Cache connections)
    
  + Kubernetes Deployment: frontend
    - Replicas: 3
    - Image: contoso.azurecr.io/inventory-frontend:sha-a1b2c3d
    - Environment variables: 3

Estimated provisioning time: 8-12 minutes

⚠️ This will create real cloud resources and incur costs.
   Review the plan carefully before approving.

[Approve] [Reject]
```

Emma clicks **[Approve]**.

```
✓ Provisioning PostgreSQL database... (6m 20s)
✓ Provisioning Redis cache... (4m 15s)
✓ Deploying backend to Kubernetes... (1m 10s)
✓ Deploying frontend to Kubernetes... (0m 45s)
✓ Configuring connections... (0m 30s)

====================================
Deployment Complete! 🎉
====================================

Application: inventory-service
Environment: prod
Status: Healthy

Endpoints:
  - Frontend: https://inventory-frontend.contoso.com
  - Backend API: https://inventory-api.contoso.com

Resources Created:
  - PostgreSQL: contoso-prod-inventory-db.postgres.database.azure.com
  - Redis: contoso-prod-inventory-cache.redis.cache.windows.net
  - Kubernetes Deployments: 2
  - Kubernetes Services: 2
  - Ingress Routes: 2

Monitoring:
  - Application Insights: https://portal.azure.com/...
  - Kubernetes Dashboards: kubectl port-forward...

Next Steps:
  - Monitor application health in Application Insights
  - Set up alerts and dashboards
  - Configure CI/CD for application updates
```

A final comment is posted to the PR with the deployment summary and links.

#### 7. Ongoing Operations

Emma's application is now running. The repository README has been automatically updated with the application graph visualization showing the live architecture.

When Emma needs to make changes — add a message queue, scale Redis, or update the database tier — she simply:

1. Invokes `@radius-platform design --refine "add message queue"` 
2. Reviews the updated design
3. Runs `@radius-platform generate` to update artifacts
4. Runs `@radius-platform deploy` to create a new PR

The entire workflow remains consistent, auditable, and compliant with her organization's standards.

### Journey Summary

**Time to Production**: ~75 minutes from empty directory to running application

**Developer Actions**:
- Invoked `@radius-platform` skill 5 times
- Wrote application code in familiar languages (Go, React)
- Never manually authored IaC or worried about cloud provider APIs
- Never left GitHub/Copilot CLI workflow

**Platform Team Confidence**:
- 100% compliance with organizational standards
- Full audit trail through Git history
- Visual architecture documentation always up-to-date
- Reusable patterns captured in Terraform module catalog

**Key Benefits**:
- **Developer velocity**: Infrastructure in minutes, not hours or days
- **Platform governance**: Automated compliance, no manual reviews needed
- **Knowledge capture**: Architecture graphs feed organizational intelligence
- **Consistency**: Same workflow for every application, every developer

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Platform Engineering Constitution Discovery and Generation (Priority: P1)

As an enterprise developer starting a new project, I want to invoke `@radius-platform init` in Copilot CLI so that the Radius skill interactively gathers my organization's platform requirements (cloud providers, compute targets, IaC tooling, approved modules, imported policies) and generates a `Platform-Engineering-Constitution.md` file that captures these requirements as the authoritative reference for all subsequent decisions.

**Why this priority**: The constitution is the foundation of the entire workflow. Every downstream action — repository creation, architecture design, artifact generation, deployment validation — depends on the organizational policies and standards captured here. Without it, no compliance-aware workflow is possible.

**Independent Test**: Can be fully tested by invoking `@radius-platform init` in Copilot CLI in an empty directory, answering the interactive prompts (or providing defaults), and verifying that a valid `Platform-Engineering-Constitution.md` is generated with the configured policies, cloud providers, compute targets, IaC tooling preferences, and imported policy rules.

**Acceptance Scenarios**:

1. **Given** an empty project directory with no existing Radius configuration, **When** the developer invokes `@radius-platform init`, **Then** the system presents an interactive guided conversation that asks about organization name, cloud providers, compute targets, existing policy sources, IaC tooling, and approved module locations.
2. **Given** the developer answers all prompts, **When** the initialization completes, **Then** a `Platform-Engineering-Constitution.md` file is written to the repository root containing all configured policies, provider preferences, naming conventions, and compliance rules in a structured, human-readable format.
3. **Given** the developer indicates they have existing infrastructure policies (e.g., Azure Policy, AWS Config), **When** the import path is provided, **Then** the skill imports and incorporates those policies into the constitution.
4. **Given** the developer skips optional prompts with default responses, **When** initialization completes, **Then** the constitution is generated with sensible defaults and clearly marks which settings use defaults versus explicit choices.
5. **Given** a directory already containing a `Platform-Engineering-Constitution.md`, **When** the developer invokes `@radius-platform init`, **Then** the skill offers to update the existing constitution (preserving previous settings) rather than overwriting it.

---

### User Story 2 — Radius Agent Skill and MCP Server Initialization (Priority: P1)

As an enterprise developer, I want `@radius-platform init` to also generate the Radius platform skill definition, start a local Radius MCP server, and register the MCP server with Copilot CLI so that the skill automatically has access to Radius platform guidance and tools.

**Why this priority**: The platform skill and MCP server are the integration surfaces that connect Radius's knowledge to Copilot. Without them, the developer must manually apply platform standards — defeating the purpose of the AI-assisted workflow. These must be initialized alongside the constitution so the developer experience is seamless from the first Copilot session.

**Independent Test**: Can be tested by invoking `@radius-platform init`, then verifying that (a) the Radius platform skill is registered (visible in Copilot), (b) the MCP server is running and reachable, and (c) Copilot lists the Radius MCP tools as available.

**Acceptance Scenarios**:

1. **Given** `@radius-platform init` has completed successfully, **When** the skill file is inspected, **Then** it exists at `.github/skills/radius-platform/SKILL.md` and contains instructions grounded on the `Platform-Engineering-Constitution.md` policies.
2. **Given** `@radius-platform init` has completed, **When** the MCP server configuration is inspected, **Then** a `.github/copilot/mcp.json` file exists and points to the locally running Radius MCP server.
3. **Given** the MCP server is running, **When** Copilot CLI is used in the repository, **Then** the Radius MCP server is connected and available tools can be invoked by the skill (e.g., `resolve_resource_types`, `search_recipes`, `generate_app_definition`, `validate_constitution`, `get_environments`).
4. **Given** the platform skill is registered, **When** a developer invokes `@radius-platform` commands, **Then** the skill applies Radius platform guidance and uses the MCP server tools for its operations.
5. **Given** the MCP server process is terminated, **When** the developer invokes `@radius-platform restart-mcp`, **Then** the MCP server is restarted and reconnected.

---

### User Story 3 — Application Architecture Design via Copilot and Radius (Priority: P1)

As an enterprise developer with application code in a repository, I want to invoke `@radius-platform design` so that the skill (using the MCP server) analyzes my source code, identifies services and dependencies, detects architecture patterns, proposes an infrastructure design compliant with my platform constitution, and presents it as a rich interactive application graph for my review.

**Why this priority**: Architecture design is the core value proposition of the Radius-Copilot integration. This is where Radius's platform knowledge, the MCP server's deterministic tooling, and Copilot's conversational interface converge to transform an application description into a compliant infrastructure design. The visual application graph — not raw IaC — is the primary interface. Without this, the integration is just documentation.

**Independent Test**: Can be tested by invoking `@radius-platform design` in a repository containing at least two services with infrastructure dependencies (e.g., a frontend referencing Redis and a backend referencing PostgreSQL), and verifying that the skill presents a correct application graph visualization, lists all components with their cloud mappings, and confirms compliance with the constitution.

**Acceptance Scenarios**:

1. **Given** a repository with application code and a `Platform-Engineering-Constitution.md`, **When** the developer invokes `@radius-platform design`, **Then** the skill analyzes the source code and presents a summary of detected services, their dependencies, and inter-service connections.
2. **Given** the skill has analyzed the application, **When** the architecture proposal is presented, **Then** it includes a rich application graph visualization showing services, dependencies, connections, and the target compute and cloud platform. The graph is the primary artifact — not Markdown or raw Bicep.
3. **Given** the architecture pattern is detected (e.g., multi-service with managed backing services), **When** the design is presented, **Then** the pattern name is identified and the design is validated against the platform constitution with a compliance summary (e.g., "12/12 policies pass").
4. **Given** the proposed architecture includes infrastructure components, **When** each component is listed, **Then** it shows the mapped cloud service (e.g., "Azure Cache for Redis"), the approved IaC module from the organization's module catalog, and the module version.
5. **Given** the developer disagrees with a suggested component or mapping, **When** they provide refinement feedback via `@radius-platform design --refine`, **Then** the skill adjusts the architecture proposal accordingly and re-validates against the constitution.
6. **Given** the application graph is generated, **When** the developer confirms the design, **Then** the graph is persisted as `radius/app-graph.json` and a rendered visualization is added to the repository README for ongoing visibility.

---

### User Story 4 — Application Definition and IaC Artifact Generation (Priority: P1)

As an enterprise developer who has confirmed the architecture design, I want to invoke `@radius-platform generate` to create the complete set of Radius deployment artifacts — the application graph, application definition file, resource type definitions, recipe configurations, and environment configuration — so that I have everything needed to deploy my application. The application graph is the primary artifact that drives the developer experience; the `app.bicep` and other IaC files are generated behind the scenes.

**Why this priority**: Artifact generation converts the approved architecture design into deployable code. This is the bridge between design and deployment. The visual graph-first approach means developers interact with the application graph; the Bicep, resource types, and recipes are generated transparently by the skill. Without generated artifacts, the developer must hand-author IaC — the exact manual work this workflow eliminates.

**Independent Test**: Can be tested by invoking `@radius-platform generate` after confirming an architecture design and verifying that the generated files exist in the `radius/` directory, that the application graph is valid and renderable, that the application definition file is syntactically valid, that resource types reference correct schemas, and that recipes reference approved IaC modules from the constitution.

**Acceptance Scenarios**:

1. **Given** the developer confirms the architecture design, **When** artifact generation is triggered via `@radius-platform generate`, **Then** the skill (via the MCP server) resolves resource types — first checking the repository for existing types, then fetching contributed types from the community catalog, and generating custom types only when no match exists.
2. **Given** resource types are resolved, **When** recipe generation runs, **Then** each infrastructure dependency is mapped to an approved IaC module from the organization's module catalog (as specified in the constitution), and the module is wrapped in a Radius recipe configuration.
3. **Given** all artifacts are generated, **When** the output is presented, **Then** the following files exist: `radius/app-graph.json` (application graph — primary artifact), `radius/app.bicep` (application definition — generated behind the scenes), resource type definitions (e.g., `radius/types/*.yaml`), recipe configurations (e.g., `radius/recipes/*/main.tf`), and environment configuration (e.g., `radius/environments/*.bicep`). The developer is shown the visual graph, not the raw files.
4. **Given** artifacts are generated, **When** constitution validation runs, **Then** all generated artifacts are validated against the `Platform-Engineering-Constitution.md` and a compliance report is displayed with zero violations or specific flagged issues with remediation suggestions.
5. **Given** the developer requests changes after seeing the generated artifacts, **When** they provide modification instructions and re-invoke generation, **Then** the skill regenerates only the affected artifacts and re-validates compliance.

---

### User Story 5 — Pull Request-Based Deployment via Repo Radius (Priority: P1)

As an enterprise developer with generated Radius artifacts, I want to invoke `@radius-platform deploy` to create a pull request containing the deployment artifacts so that merging the PR triggers Repo Radius to deploy my application to the target environment, following standard GitHub review and approval workflows.

**Why this priority**: The PR-based deployment flow is the production path for getting applications running. It connects the design workflow to actual infrastructure provisioning through a familiar, auditable, and governed mechanism (pull requests). Without it, generated artifacts have no deployment path.

**Independent Test**: Can be tested by invoking `@radius-platform deploy`, verifying a PR is created with the `radius/` directory contents, and confirming that merging the PR triggers a deployment workflow (which can be validated against a test environment or a dry-run mode).

**Acceptance Scenarios**:

1. **Given** generated Radius artifacts exist in the repository, **When** the developer invokes `@radius-platform deploy`, **Then** the skill creates a pull request on a feature branch containing all files in the `radius/` directory with a descriptive title and summary.
2. **Given** the PR is created, **When** it appears in GitHub, **Then** a Radius-integrated GitHub Action or App adds a comment to the PR with a summary of changes and a before/after architecture visualization of the application graph.
3. **Given** the PR is approved and merged, **When** the merge event fires, **Then** the Repo Radius deployment workflow is triggered and the developer is prompted to approve the deployment.
4. **Given** the deployment is approved, **When** Repo Radius processes the artifacts, **Then** it compiles the Radius application definition, provisions the target environment, and deploys all application components and infrastructure dependencies.
5. **Given** the deployment completes, **When** the status is updated, **Then** a notification is posted to the PR with the deployment status (success/failure), links to deployed resources, and any relevant diagnostics.

---

### User Story 6 — Architecture Graph Visualization in Pull Requests and README (Priority: P2)

As a developer reviewing a PR that includes Radius deployment artifacts, I want to see a rich architecture graph visualization directly in the PR so that I can understand the application topology, see what resources will be deployed, and identify any changes from the previous state. I also want the application graph to be rendered in the repository README for ongoing team visibility.

**Why this priority**: Visual review of architecture changes is essential for confident approvals. Without it, reviewers must read raw Bicep and YAML to understand what will be deployed. The visualization reduces review time and increases the likelihood of catching design issues before deployment. The README graph provides always-on visibility into the application topology.

**Independent Test**: Can be tested by opening a PR containing Radius artifacts in the GitHub UI and verifying that a rendered architecture graph is displayed as a PR comment, showing services, dependencies, connections, and change indicators (added/modified/removed). Additionally verify the README contains a rendered graph.

**Acceptance Scenarios**:

1. **Given** a PR contains new Radius artifacts (first deployment), **When** the PR is opened, **Then** an architecture graph visualization is rendered showing all services, infrastructure dependencies, and connections with their target cloud services.
2. **Given** a PR modifies existing Radius artifacts, **When** the PR is opened, **Then** the visualization shows a before/after diff of the architecture graph, highlighting added, modified, and removed components.
3. **Given** the architecture graph includes infrastructure components, **When** the visualization is rendered, **Then** each component shows its resource type, target cloud service, and the IaC module that provisions it.
4. **Given** a policy violation is detected during PR analysis, **When** the graph is rendered, **Then** the violating component is highlighted with a description of the violation and a suggested remediation.
5. **Given** a PR is merged with Radius artifacts, **When** the merge completes, **Then** the repository README is automatically updated with the latest application graph visualization reflecting the current deployed state.

---

### User Story 7 — Deployment Plan Review and Approval (Priority: P2)

As an enterprise developer, I want to review a detailed deployment plan (including the full resource graph with rich data) before approving the deployment so that I can verify exactly what will be provisioned and catch any issues before resources are created.

**Why this priority**: Deployment approval is a governance gate. Developers and platform teams need confidence that the deployment will do what they expect. A detailed, visual deployment plan reduces risk and supports compliance requirements.

**Independent Test**: Can be tested by merging a PR with Radius artifacts, observing the deployment plan rendered in the GitHub UI, and verifying it lists all resources to be created/modified/deleted with their properties and any flagged issues.

**Acceptance Scenarios**:

1. **Given** a PR has been merged and the deployment workflow starts, **When** the plan is generated, **Then** it is rendered in the PR as a rich visualization showing all resources to be created, their properties, and their relationships.
2. **Given** the deployment plan includes resources that violate platform policies, **When** the plan is displayed, **Then** violations are highlighted with descriptions and suggested fixes, and the deployment can be blocked until violations are resolved.
3. **Given** the developer reviews the deployment plan, **When** they approve the deployment, **Then** Repo Radius proceeds with provisioning and deployment.
4. **Given** the developer rejects the deployment plan, **When** they decline approval, **Then** the deployment is cancelled with no resources created and the developer is guided on how to correct the issues and re-submit.

---

### User Story 8 — Constitution-Compliant Application Scaffolding via Copilot (Priority: P2)

As an enterprise developer, I want Copilot to scaffold my application code, Dockerfiles, and GitHub repository with settings that comply with my platform constitution so that the application is compliant from inception rather than requiring retrofitting.

**Why this priority**: Compliance from the start is cheaper than remediation later. If the repository, container images, and code patterns are constitution-aware from day one, the architecture design and deployment steps encounter fewer violations. However, this story enhances rather than enables the core workflow.

**Independent Test**: Can be tested by asking Copilot to scaffold an application in a repository that has a `Platform-Engineering-Constitution.md`, then verifying that the generated repo settings (visibility, branch protection), Dockerfiles (approved base images), and container image tags match the constitution's policies.

**Acceptance Scenarios**:

1. **Given** a `Platform-Engineering-Constitution.md` exists in the repository, **When** the developer asks Copilot to scaffold application code, **Then** Copilot reads the constitution and applies its policies (e.g., approved base images, naming conventions) to the generated code and Dockerfiles.
2. **Given** the developer asks Copilot to create a GitHub repository, **When** the repository is created, **Then** its visibility, branch protection rules, and default branch settings match the constitution's repository policies.
3. **Given** the developer asks Copilot to build and publish container images, **When** the images are built, **Then** they use approved base images, follow the constitution's tagging policy (e.g., immutable tags, git SHA tags), and are pushed to the approved container registry.
4. **Given** the constitution specifies image scanning and signing requirements, **When** container images are built, **Then** they are scanned for vulnerabilities and signed according to the constitution's supply chain policies before being pushed.

---

### User Story 9 — Deployment Status and Notifications (Priority: P3)

As an enterprise developer who has approved a deployment, I want to receive timely notifications about deployment progress and completion status in the PR and in Copilot so that I can monitor the deployment without leaving my workflow.

**Why this priority**: Deployment visibility reduces anxiety and enables faster issue response. However, this is an enhancement over the core workflow (which already provides deployment completion status).

**Independent Test**: Can be tested by approving a deployment and verifying that progress updates appear in the PR comments and that a final summary with resource links and diagnostics is posted upon completion.

**Acceptance Scenarios**:

1. **Given** a deployment has been approved, **When** deployment progresses through stages (environment setup, resource provisioning, application deployment), **Then** progress updates are posted to the PR.
2. **Given** a deployment completes successfully, **When** the final status is posted, **Then** it includes links to the deployed application endpoint, links to provisioned cloud resources, and a final architecture graph reflecting the live state.
3. **Given** a deployment fails, **When** the failure status is posted, **Then** it includes the specific error, the resource that failed, relevant diagnostic logs, and suggested remediation steps.

---

### User Story 10 — Enterprise Application Graph Database and Organizational Intelligence (Priority: P2)

As a platform engineering team lead at an enterprise, I want Radius to publish application graph data to a centralized enterprise graph database so that I can see every application's topology across the organization, correlate graph data with commit history and IaC deployment metadata, and answer questions like "which teams are using Terraform version X?" or "which applications depend on Redis?".

**Why this priority**: Enterprise visibility into application topologies is a key differentiator. Individual app graphs are valuable, but when aggregated across an organization — and correlated with commit data, IaC deployment data (Terraform/Bicep versions, module usage), and team ownership — they become a powerful intelligence layer. This enables platform teams to make informed decisions about standards, migration planning, and risk assessment.

**Independent Test**: Can be tested by deploying two applications with different architectures, publishing their graphs to the enterprise database, and verifying that queries can return cross-application insights (e.g., "all applications using PostgreSQL", "Terraform version distribution by team").

**Acceptance Scenarios**:

1. **Given** a successful deployment via Repo Radius, **When** the deployment completes, **Then** the application graph is published to the enterprise graph database with metadata including team/org, commit SHA, deployment timestamp, and IaC versions used.
2. **Given** multiple applications have published their graphs, **When** a platform engineer queries the enterprise database, **Then** they can see a unified view of all application topologies, filter by team, cloud provider, resource type, or IaC module.
3. **Given** the enterprise graph database contains deployment history, **When** a platform engineer queries IaC version usage, **Then** they can see which teams are using which Terraform/Bicep versions and module versions, enabling informed migration and standardization decisions.
4. **Given** the enterprise graph database contains cross-application data, **When** a dependency query is run (e.g., "all applications using Redis"), **Then** the system returns a list of applications, their teams, and the specific Redis resource type and recipe they use.
5. **Given** the organization has a GitHub Enterprise instance, **When** Radius is configured at the org level, **Then** graph data is automatically published from all repositories that use Radius, without per-repo configuration.

---

### Edge Cases

- What happens when the developer invokes `@radius-platform init` in a directory that is not a git repository? The system should initialize Radius configuration locally and warn that git-based features (PR deployment, Repo Radius) will not be available until the directory is a git repository.
- How does the system handle a `Platform-Engineering-Constitution.md` that references policies from a cloud provider the developer no longer has access to? The system should flag inaccessible policy sources during initialization and allow the developer to skip or provide alternative sources.
- What happens when the MCP server becomes unreachable during an active Copilot session? The skill should gracefully degrade and inform the developer that live tooling is temporarily unavailable, suggesting `@radius-platform restart-mcp`.
- How does the system handle conflicting policies in the constitution (e.g., one policy requires public endpoints and another prohibits them)? The system should detect conflicts during validation and surface them with descriptions of both conflicting rules.
- What happens when the approved IaC module catalog does not contain a module for a detected dependency? The system should fall back to community-contributed resource types, then to a generated custom resource type, and clearly label the source and confidence level of each.
- What happens when a PR contains Radius artifacts that reference a target environment that does not exist? The deployment workflow should fail with a clear error identifying the missing environment and instructions to create it.
- How does the system handle multiple developers invoking `@radius-platform init` in the same repository with different organizational preferences? The constitution should be treated as repository-level configuration; the second initialization should update (not overwrite) the existing constitution and flag conflicting settings for resolution.
- What happens when the developer's Copilot session context window is exceeded during a complex architecture design? The Radius platform skill provides persistent context grounding (loaded at session start), and MCP server tools are stateless and idempotent — the developer can re-invoke any skill command in a new session to reconstruct the design state. Critical design artifacts (constitution, app-graph.json) are persisted to disk, not held in memory.

## Requirements *(mandatory)*

### Functional Requirements

#### Initialization (`@radius-platform init`)

- **FR-001**: The skill MUST present an interactive guided conversation that collects organization name, cloud providers, compute targets, existing policy import sources, IaC tooling preferences, and approved module catalog locations.
- **FR-002**: The skill MUST generate a `Platform-Engineering-Constitution.md` file in the repository root containing all collected policies, provider preferences, naming conventions, compliance rules, container image policies, and repository policies in a structured, human-readable format.
- **FR-003**: The skill MUST import existing infrastructure policies (e.g., Azure Policy, AWS Config) when the developer provides a source, incorporating them into the constitution.
- **FR-004**: The skill MUST generate a Radius platform skill definition file at `.github/skills/radius-platform/SKILL.md` with instructions grounded on the generated constitution.
- **FR-005**: The skill MUST start a local Radius MCP server and register it with Copilot CLI by writing the server configuration to `.github/copilot/mcp.json`.
- **FR-006**: The MCP server MUST expose the following tools: resource type resolution, recipe search, application definition generation, constitution validation, and environment management.
- **FR-007**: When a `Platform-Engineering-Constitution.md` already exists, the skill MUST offer to update it (preserving existing settings) rather than overwriting.
- **FR-008**: The skill MUST provide sensible defaults for all prompts so that a developer can complete initialization by accepting defaults and still receive a usable constitution.

#### Architecture Design (`@radius-platform design`)

- **FR-009**: The Radius platform skill MUST analyze application source code in the repository to detect services, dependencies, and connections when invoked for architecture design.
- **FR-010**: The system MUST present a visual architecture diagram (text-based for CLI, rendered as Mermaid in GitHub PR comments) showing all detected services, infrastructure dependencies, and their connections. The visual graph is the primary interface; raw IaC files are generated behind the scenes.
- **FR-011**: The system MUST identify the architecture pattern (e.g., multi-service with managed backing services, event-driven, microservices) and display it alongside the design.
- **FR-012**: The system MUST validate the proposed architecture against the `Platform-Engineering-Constitution.md` and display a compliance summary showing how many policies pass and which (if any) fail.
- **FR-013**: The system MUST map each infrastructure dependency to the corresponding cloud service based on the constitution's cloud provider preferences and display the mapping with the approved IaC module reference and version.
- **FR-014**: The system MUST allow the developer to provide feedback on the proposed architecture in natural language, adjust the design accordingly, and re-validate compliance.

#### Artifact Generation (`@radius-platform generate`)

- **FR-015**: The system MUST resolve resource types for each application component by checking (in order): existing types in the repository, contributed types from the community catalog, and generating custom types only when no match exists.
- **FR-016**: The system MUST generate recipe configurations that wrap approved IaC modules from the organization's module catalog for each infrastructure dependency.
- **FR-017**: The system MUST generate a syntactically valid application definition file (`radius/app.bicep`) declaring all services, dependencies, connections, and secret bindings.
- **FR-018**: The system MUST generate environment configuration files (e.g., `radius/environments/*.bicep`) for the target deployment environment.
- **FR-019**: The system MUST validate all generated artifacts against the `Platform-Engineering-Constitution.md` and report compliance status.
- **FR-020**: The system MUST generate an application graph data file (`radius/app-graph.json`) suitable for visualization in pull requests, README embedding, and enterprise graph database publishing.

#### Deployment (`@radius-platform deploy`)

- **FR-021**: The skill MUST support creating a pull request containing all generated Radius artifacts in the `radius/` directory with a descriptive title and summary.
- **FR-022**: A Radius-integrated GitHub mechanism (Action or App) MUST add an architecture graph visualization as a comment on PRs that contain Radius artifacts.
- **FR-023**: The PR visualization MUST show before/after architecture diffs when modifying existing Radius artifacts.
- **FR-024**: Merging a PR containing Radius artifacts MUST trigger the Repo Radius deployment workflow.
- **FR-025**: The deployment workflow MUST generate and display a deployment plan for developer review and approval before provisioning any resources.
- **FR-026**: The deployment plan MUST highlight any policy violations with descriptions and suggested remediations.
- **FR-027**: Upon deployment approval, Repo Radius MUST compile the application definition, provision the target environment, and deploy all application components and infrastructure dependencies.
- **FR-028**: The system MUST post deployment status (progress, completion, failure) and resource links as comments on the originating PR.

#### Application Graph Visibility

- **FR-035**: The system MUST update the repository README with a rendered application graph visualization after each successful deployment or artifact generation.
- **FR-036**: The Radius platform skill MUST build and present the application graph as a rich visual artifact during design sessions, not as raw Markdown or IaC code.
- **FR-037**: The skill MUST support publishing application graph data to an enterprise-level graph database upon deployment, including metadata (team, commit SHA, deployment timestamp, IaC tool versions, module versions).
- **FR-038**: The enterprise graph database MUST support queries across applications for cross-cutting insights (e.g., dependency usage, IaC version distribution by team, migration readiness assessment).
- **FR-039**: The Radius platform skill MUST follow a skills-centric pattern where the skill drives actions like initialization, architecture design, deployment, and graph generation through conversational interfaces in Copilot CLI with clear command syntax (e.g., `@radius-platform <command>`).

#### Constitution-Compliant Scaffolding

- **FR-029**: When Copilot scaffolds application code in a repository with a constitution, it MUST apply the constitution's policies to generated code, Dockerfiles, and repository settings.
- **FR-030**: Container images built through the workflow MUST use approved base images, follow the tagging policy, and be pushed to the approved registry as specified in the constitution.
- **FR-031**: The system MUST scan container images for vulnerabilities and sign them according to the constitution's supply chain policies when such policies are defined.

#### General

- **FR-032**: All generated deployment artifacts MUST be placed in the `radius/` directory within the repository root.
- **FR-033**: The skill MUST support a non-interactive mode that accepts configuration from a file or inline parameters for use in automation.
- **FR-034**: The Radius platform skill MUST provide clear command syntax (e.g., `@radius-platform init`, `@radius-platform design`, `@radius-platform generate`, `@radius-platform deploy`) for all capabilities.

### Key Entities

- **Platform Engineering Constitution** (`Platform-Engineering-Constitution.md`): The authoritative document capturing an organization's platform requirements, cloud provider preferences, compute targets, IaC tooling, approved module catalogs, imported policies, naming conventions, container image policies, and repository policies. All downstream decisions are grounded on this document.
- **Radius Platform Skill** (`.github/skills/radius-platform/SKILL.md`): A GitHub Copilot skill that provides platform engineering capabilities through conversational commands (e.g., `@radius-platform init`, `@radius-platform design`, `@radius-platform deploy`). Grounded on the constitution and uses the MCP server for deterministic operations.
- **Radius MCP Server**: A locally hosted service exposing deterministic tools for resource type resolution, recipe search, application definition generation, constitution validation, and environment management. Called by Copilot during architecture design and artifact generation.
- **Application Definition** (`radius/app.bicep`): A Bicep file declaring all application services, infrastructure dependencies, connections, and secret bindings. The primary deployment artifact.
- **Resource Type**: A schema definition for an infrastructure component (e.g., `Applications.Core/containers`, `Applications.Datastores/redisCaches`). Can be sourced from the repository, community catalog, or generated custom.
- **Recipe**: A configuration that maps an infrastructure dependency to a specific IaC module (Terraform or Bicep) that provisions it. References approved modules from the organization's catalog.
- **Environment Configuration** (`radius/environments/*.bicep`): Files that define the target deployment environment (cloud provider, region, resource group, etc.) for Radius.
- **Application Graph** (`radius/app-graph.json`): A structured data file representing the full application topology suitable for visualization in pull requests, README embedding, and enterprise graph database publishing. The application graph is the primary artifact of the workflow — it drives the visual-first developer experience and is the unit of organizational intelligence when aggregated at enterprise scale.
- **Deployment Plan**: A pre-deployment summary showing all resources to be created, modified, or deleted, with their properties and policy compliance status. Generated by Repo Radius and displayed in the PR for approval.
- **Repo Radius**: The Radius control plane integrated with GitHub that processes repository-hosted Radius artifacts, compiles application definitions, and orchestrates deployment upon PR merge and approval.
- **Enterprise Graph Database**: A centralized store of application graph data across all repositories in an organization. Aggregates graph metadata (team, commit SHA, deployment timestamp, IaC versions, module usage) to enable cross-organizational queries like "which teams use Terraform version X?" or "which applications depend on Redis?". Published to automatically on deployment.

## Assumptions

- The developer has a GitHub account with Copilot CLI access and the ability to create repositories in their organization.
- The Radius platform skill is implemented and available for use in GitHub Copilot CLI.
- The MCP server runs locally on the developer's machine (no hosted infrastructure required for the development workflow). It starts automatically during `@radius-platform init` and can be restarted via `@radius-platform restart-mcp`.
- Repo Radius (the server-side deployment control plane) is pre-configured for the organization's GitHub environment. Setting up Repo Radius itself is outside the scope of this specification.
- The developer's cloud provider credentials are available in the local environment for policy import and container registry access. Credential management is handled by existing tooling (e.g., `az login`, `aws configure`).
- The MCP server's tool responses are deterministic — given the same inputs (source code, constitution, module catalog), they produce the same outputs. This ensures reproducibility.
- The Radius platform skill is compatible with Copilot CLI and provides clear command syntax for all operations.
- The architecture graph visualization in PRs uses a format that renders natively in the GitHub PR UI (e.g., Mermaid diagrams in Markdown or an image generated by the GitHub App/Action).
- For enterprise graph database publishing, an API endpoint or event sink is available for receiving application graph data. The specific backend (e.g., GitHub's internal graph, a third-party graph database, or a Radius-hosted store) is configurable per organization.
- The Radius platform skill follows a skills-centric pattern where the skill drives the workflow through conversational commands with clear syntax (e.g., `@radius-platform <command>`).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A developer can go from an empty directory to a fully initialized Radius workspace (constitution, platform skill, MCP server running) in under 5 minutes using `@radius-platform init` with interactive prompts.
- **SC-002**: A developer with a multi-service repository (2+ services, 2+ infrastructure dependencies) can produce a confirmed architecture design with compliance validation using `@radius-platform design` in under 10 minutes.
- **SC-003**: Generated Radius artifacts (application definition, resource types, recipes, environment configuration) are produced within 2 minutes of architecture confirmation and are fully compliant with the platform constitution on first generation for standard architectures.
- **SC-004**: 100% of generated application definitions are syntactically valid and reference only approved IaC modules from the organization's catalog when such modules exist.
- **SC-005**: Pull requests containing Radius artifacts display an architecture graph visualization within 60 seconds of PR creation.
- **SC-006**: The end-to-end workflow from application description to deployed application (including all interactive review and approval steps) completes in under 30 minutes for a standard multi-service application.
- **SC-007**: Developers report confidence in the architecture design process measured by 85%+ of design proposals being approved on first review (without major revisions).
- **SC-008**: All deployment artifacts pass constitution validation before a PR is created — zero policy violations reach the PR review stage unless the developer explicitly overrides.
