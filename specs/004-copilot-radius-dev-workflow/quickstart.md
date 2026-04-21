# Quickstart: Developer Workflow in Copilot and GitHub Using Radius

This guide walks through the end-to-end developer workflow from project initialization to deployed application using GitHub Copilot CLI with Radius integration.

## Prerequisites

- [GitHub Copilot CLI](https://docs.github.com/copilot) with an active subscription
- Git repository initialized
- Cloud provider credentials configured (e.g., `az login`, `aws configure`)
- Kubernetes cluster available (e.g., via `k3d` or `kind`)

## Step 1: Initialize Radius with Platform Constitution

Open Copilot CLI in your project directory:

```bash
# Navigate to your project repository
cd my-app

# Open Copilot CLI
gh copilot
```

Then invoke the Radius platform skill to initialize:

```
> @radius-platform init
```

The skill will guide you through an interactive conversation to gather:
1. **Organization name** — your team or company identifier
2. **Cloud providers** — Azure, AWS, or both (with subscription/account details)
3. **Compute targets** — Kubernetes cluster context, ACA, etc.
4. **IaC tooling preference** — Terraform (recommended) or Bicep
5. **Module catalog** — your approved Terraform registry or Git repo
6. **Policies** — import from Azure Policy, AWS Config, or define manually

The skill will generate and create the following in your repository:

```
my-app/
├── Platform-Engineering-Constitution.md   # Organization policies
├── .github/
│   ├── skills/
│   │   └── radius-platform/SKILL.md       # Radius platform skill
│   └── copilot/mcp.json                   # MCP server registration
└── (your existing source code)
```

The skill will also start the MCP server automatically and confirm:
```
✓ Platform constitution created
✓ Radius platform skill registered
✓ MCP server started
✓ Ready to design your application architecture!
```

## Step 2: Design Architecture with Radius Platform Skill

Invoke the design capability:

```
> @radius-platform design
```

The skill (with MCP server connected) will:
1. **Scan your source code** — detect services, dependencies, and connections
2. **Propose an architecture** — with a visual diagram and component mappings
3. **Validate compliance** — check the design against your constitution

Example output:
```
Architecture: Multi-service with managed backing services

┌─────────┐     ┌──────────┐     ┌────────────┐
│ Frontend │────>│ Backend  │────>│ PostgreSQL  │
│ (React)  │    │ (Go API) │    │ (Azure DB)  │
└─────────┘     └──────────┘     └────────────┘
                      │
                      v
               ┌────────────┐
               │   Redis    │
               │(Azure Cache)│
               └────────────┘

Compliance: 8/8 policies pass ✓
```

Refine the design through conversation:
```
> @radius-platform design --refine "use AWS ElastiCache for Redis instead of Azure Cache"
```

## Step 3: Generate Deployment Artifacts

Once you confirm the architecture:

```
> @radius-platform generate
```

The skill generates:
```
radius/
├── app.bicep                    # Application definition
├── app-graph.json               # Architecture graph for PR visualization
├── types/
│   └── postgreSqlDatabases.yaml # Resource type schema
├── recipes/
│   ├── azure-postgres/main.tf   # PostgreSQL recipe (Terraform)
│   └── aws-redis/main.tf        # Redis recipe (Terraform)
└── environments/
    └── dev.bicep                # Development environment
```

All artifacts are validated against the constitution:
```
Compliance Report: 8/8 policies pass
  ✓ naming-convention: All resources follow {org}-{env}-{service} pattern
  ✓ encryption-at-rest: All data stores enable encryption
  ✓ approved-modules: All recipes use catalog modules
  ...
```

## Step 4: Create Pull Request and Review

```
> @radius-platform deploy
```

The skill creates a PR with:
- All files in the `radius/` directory
- Architecture graph visualization (Mermaid diagram in PR comment)
- Compliance report summary

The PR displays a before/after architecture diff for reviewers.

## Step 5: Deploy via Repo Radius

1. **Review** the PR and architecture visualization
2. **Merge** the PR — triggers the Radius deployment workflow
3. **Review the deployment plan** — see all resources to be created
4. **Approve** the deployment
5. **Monitor** progress via PR status comments

## Non-Interactive Mode

For CI/CD automation or batch setup, use the skill with configuration files:

```
> @radius-platform init --config constitution-config.yaml
```

Or provide all details in a single invocation:

```
> @radius-platform init --org "contoso" --cloud azure --subscription "xxxx-xxxx" --compute kubernetes --iac terraform
```

The skill will parse your parameters and create all necessary artifacts without further interaction.

## Troubleshooting

| Issue | Solution |
|-------|----------|
| MCP server not connecting | Run: `@radius-platform restart-mcp` |
| Agent Skill not loading | Verify `.github/skills/radius-platform/SKILL.md` exists, then reload VS Code |
| Constitution validation failures | Run: `@radius-platform validate` |
| Missing resource types | Run: `@radius-platform list-types` or `@radius-platform create-type --service [name]` |

## Next Steps

- [Radius Platform Skill Reference](./skills/radius-platform/SKILL.md) — Complete skill capabilities and command reference
- [Constitution Reference](./contracts/constitution-schema.yaml) — Full schema for constitution files
- [MCP Server Tools](./contracts/mcp-server.yaml) — Available MCP tool definitions
- [Data Model](./data-model.md) — Entity definitions and relationships
