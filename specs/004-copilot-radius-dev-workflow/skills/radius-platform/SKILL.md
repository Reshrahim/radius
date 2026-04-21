---
name: radius-platform
description: |
  Platform engineering skill for Radius that handles initialization, architecture design, 
  artifact generation, validation, and deployment. Use when setting up Radius for a project, 
  designing cloud architectures, generating IaC artifacts, or deploying applications.
license: MIT
metadata:
  author: radius-platform
  version: "1.0.0"
  mcp_server: true
---

# Radius Platform Engineering Skill

This skill provides end-to-end platform engineering capabilities for Radius projects, from initialization through deployment. It bridges developers and platform engineers by enforcing organizational policies while enabling self-service infrastructure.

## Capabilities

The skill provides the following capabilities:

| Capability | Command | Description |
|------------|---------|-------------|
| Initialize | `@radius-platform init` | Set up Radius with platform constitution |
| Design | `@radius-platform design` | Generate cloud architecture from source code |
| Generate | `@radius-platform generate` | Create deployment artifacts (Bicep, Terraform) |
| Validate | `@radius-platform validate` | Check compliance with platform constitution |
| Deploy | `@radius-platform deploy` | Create PR and initiate deployment workflow |
| List Types | `@radius-platform list-types` | Show available resource types |
| Create Type | `@radius-platform create-type` | Generate a new resource type definition |
| Restart MCP | `@radius-platform restart-mcp` | Restart the MCP server |

## Prerequisites

Before using this skill, ensure:
- GitHub Copilot CLI is installed and authenticated
- Git repository is initialized
- Cloud provider credentials are configured (e.g., `az login`, `aws configure`)
- Kubernetes cluster is available (for compute targets)

## Workflow

### 1. Initialize Radius (`@radius-platform init`)

**Purpose:** Set up Radius for a project by creating the platform constitution and registering the skill.

**Interactive Mode:**
```
> @radius-platform init
```

You'll be prompted for:
1. **Organization name** — your team or company identifier
2. **Cloud providers** — Azure, AWS, or both (with subscription/account details)
3. **Compute targets** — Kubernetes cluster context, ACA, etc.
4. **IaC tooling preference** — Terraform (recommended) or Bicep
5. **Module catalog** — your approved Terraform registry or Git repo
6. **Policies** — import from Azure Policy, AWS Config, or define manually

**Non-Interactive Mode:**
```
# From configuration file
> @radius-platform init --config constitution-config.yaml

# With inline parameters
> @radius-platform init --org "contoso" --cloud azure --subscription "xxxx" --compute kubernetes --iac terraform
```

**Output:**
```
my-app/
├── Platform-Engineering-Constitution.md   # Organization policies
├── .github/
│   ├── skills/radius-platform/SKILL.md    # This skill (auto-registered)
│   └── copilot/mcp.json                   # MCP server registration
└── (your existing source code)
```

**Confirmation:**
```
✓ Platform constitution created
✓ Radius platform skill registered
✓ MCP server started
✓ Ready to design your application architecture!
```

### 2. Design Architecture (`@radius-platform design`)

**Purpose:** Analyze source code and generate a compliant cloud architecture.

**Basic Usage:**
```
> @radius-platform design
```

The skill will:
1. **Scan source code** — detect services, dependencies, and data stores
2. **Propose architecture** — with visual diagram and component mappings
3. **Validate compliance** — check against platform constitution
4. **Show recommendations** — suggest optimizations or alternatives

**Example Output:**
```
Architecture: Multi-service with managed backing services

┌─────────┐     ┌──────────┐     ┌────────────┐
│ Frontend │────>│ Backend  │────>│ PostgreSQL │
│ (React)  │    │ (Go API) │    │ (Azure DB) │
└─────────┘     └──────────┘     └────────────┘
                      │
                      v
               ┌────────────┐
               │   Redis    │
               │(Azure Cache)│
               └────────────┘

Compliance: 8/8 policies pass ✓
  ✓ naming-convention: All resources follow pattern
  ✓ encryption-at-rest: All data stores encrypted
  ✓ approved-modules: Using catalog modules
  ✓ network-isolation: VNet integration enabled
  ✓ backup-retention: 30-day retention configured
  ✓ monitoring: Application Insights enabled
  ✓ cost-tags: Required tags present
  ✓ region-restriction: Using approved region
```

**Refine Design:**
```
> @radius-platform design --refine "use AWS ElastiCache for Redis instead of Azure Cache"
```

**Options:**
- `--cloud <provider>` — Prefer specific cloud provider
- `--compute <target>` — Target specific compute platform
- `--refine "<changes>"` — Modify the proposed architecture
- `--output <file>` — Save architecture graph to file

### 3. Generate Artifacts (`@radius-platform generate`)

**Purpose:** Generate deployment artifacts (Bicep, Terraform, YAML) from the approved architecture.

**Basic Usage:**
```
> @radius-platform generate
```

**Generated Structure:**
```
radius/
├── app.bicep                    # Application definition
├── app-graph.json               # Architecture graph for visualization
├── types/
│   └── postgreSqlDatabases.yaml # Resource type schemas
├── recipes/
│   ├── azure-postgres/
│   │   ├── main.tf              # PostgreSQL recipe (Terraform)
│   │   └── variables.tf
│   └── aws-redis/
│       ├── main.tf              # Redis recipe (Terraform)
│       └── variables.tf
└── environments/
    ├── dev.bicep                # Development environment
    └── prod.bicep               # Production environment
```

**Validation Output:**
```
Compliance Report: 8/8 policies pass
  ✓ naming-convention: All resources follow {org}-{env}-{service} pattern
  ✓ encryption-at-rest: All data stores enable encryption
  ✓ approved-modules: All recipes use catalog modules
  ✓ network-isolation: VNet integration configured
  ✓ backup-retention: Backup policies applied
  ✓ monitoring: Observability configured
  ✓ cost-tags: Required tags present
  ✓ region-restriction: Using approved regions
```

**Options:**
- `--env <name>` — Generate specific environment (dev, staging, prod)
- `--output-dir <path>` — Custom output directory (default: `./radius`)
- `--format <bicep|terraform>` — Preferred IaC format (uses constitution default if not specified)

### 4. Validate (`@radius-platform validate`)

**Purpose:** Validate existing Radius artifacts against the platform constitution.

**Basic Usage:**
```
> @radius-platform validate
```

**Example Output:**
```
Validating: radius/app.bicep
Validating: radius/recipes/azure-postgres/main.tf
Validating: radius/environments/dev.bicep

Results: 7/8 policies pass, 1 warning

✓ naming-convention: PASS
✓ encryption-at-rest: PASS
✓ approved-modules: PASS
✓ network-isolation: PASS
✓ backup-retention: PASS
⚠ monitoring: WARNING - Missing Application Insights for backend service
✓ cost-tags: PASS
✓ region-restriction: PASS

Recommendation: Add Application Insights connection to backend container
```

**Options:**
- `--file <path>` — Validate specific file
- `--policy <name>` — Check specific policy
- `--fix` — Automatically fix common issues
- `--strict` — Fail on warnings

### 5. Deploy (`@radius-platform deploy`)

**Purpose:** Create a pull request with deployment workflow for review and approval.

**Basic Usage:**
```
> @radius-platform deploy
```

The skill will:
1. **Stage artifacts** — Commit generated files to a new branch
2. **Create PR** — With architecture visualization and compliance report
3. **Add workflow** — GitHub Actions workflow for deployment
4. **Request review** — Notify platform team

**Example PR Content:**
```
## Architecture Changes

### Before
(No existing architecture)

### After
[Mermaid diagram showing new architecture]

## Compliance Report
8/8 policies pass ✓

## Resource Summary
- 1 Application (myapp)
- 2 Containers (frontend, backend)
- 1 PostgreSQL Database (Azure)
- 1 Redis Cache (Azure)

## Deployment Plan
1. Create resource types in Radius environment
2. Register recipes for PostgreSQL and Redis
3. Deploy application to Kubernetes
4. Configure connections and secrets

## Review Checklist
- [ ] Architecture aligns with requirements
- [ ] Resource types are appropriate
- [ ] Recipes use approved modules
- [ ] Environment configuration is correct
- [ ] Cost estimate is acceptable
```

**Options:**
- `--branch <name>` — Custom branch name (default: auto-generated)
- `--title "<text>"` — Custom PR title
- `--draft` — Create as draft PR
- `--auto-merge` — Enable auto-merge after approval

### 6. List Resource Types (`@radius-platform list-types`)

**Purpose:** Show available resource types in the Radius environment.

**Basic Usage:**
```
> @radius-platform list-types
```

**Example Output:**
```
Available Resource Types:

Radius.Core:
  - applications

Radius.Compute:
  - containers
  - persistentVolumes
  - routes

Radius.Data:
  - postgreSqlDatabases
  - mySqlDatabases

Radius.Security:
  - secrets

Total: 7 resource types across 4 namespaces
```

**Options:**
- `--namespace <name>` — Filter by namespace
- `--search <query>` — Search resource types by name
- `--json` — Output in JSON format

### 7. Create Resource Type (`@radius-platform create-type`)

**Purpose:** Generate a new custom resource type and recipes.

**Basic Usage:**
```
> @radius-platform create-type --service mongodb
```

This invokes the underlying `radius-new-type` workflow to:
1. Create resource type schema
2. Generate recipes for approved cloud providers
3. Add to resource-types-contrib structure
4. Provide registration instructions

**Options:**
- `--service <name>` — Service name (e.g., mongodb, rabbitmq)
- `--namespace <name>` — Custom namespace (default: auto-detected)
- `--cloud <provider>` — Generate recipes for specific cloud

### 8. Restart MCP Server (`@radius-platform restart-mcp`)

**Purpose:** Restart the Radius MCP server if connection issues occur.

**Basic Usage:**
```
> @radius-platform restart-mcp
```

**Output:**
```
Stopping existing MCP server...
Starting MCP server...
✓ MCP server restarted successfully
✓ Connected to Radius environment
```

## MCP Server Integration

This skill requires an MCP (Model Context Protocol) server to access live Radius environment data:
- Registered resource types
- Available recipes
- Deployed applications
- Resource connections
- Environment configuration

The MCP server is automatically started during `@radius-platform init` and registered in `.github/copilot/mcp.json`.

## Constitution Compliance

All operations respect the platform constitution (`Platform-Engineering-Constitution.md`):
- **Cloud providers** — Only use approved providers
- **Compute platforms** — Only target approved platforms
- **IaC modules** — Only use catalog modules
- **Naming conventions** — Enforce organizational patterns
- **Tagging requirements** — Apply required tags
- **Network policies** — Follow network topology rules
- **Security policies** — Meet encryption and access requirements
- **Cost controls** — Respect budget and quota limits

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| `Constitution not found` | Missing `Platform-Engineering-Constitution.md` | Run `@radius-platform init` first |
| `MCP server not responding` | Server crashed or not started | Run `@radius-platform restart-mcp` |
| `Policy violation` | Generated artifacts don't comply | Review validation errors and run `@radius-platform validate --fix` |
| `Resource type not found` | Requested type doesn't exist | Run `@radius-platform list-types` or `@radius-platform create-type` |
| `Cloud credentials missing` | Not authenticated to cloud provider | Run `az login` or `aws configure` |

## Examples

### Complete Workflow

```bash
# 1. Initialize Radius for your organization
> @radius-platform init --org "contoso" --cloud azure --compute kubernetes --iac terraform

# 2. Design architecture from source code
> @radius-platform design

# 3. Refine the design if needed
> @radius-platform design --refine "add Redis cache for session storage"

# 4. Generate deployment artifacts
> @radius-platform generate

# 5. Validate compliance
> @radius-platform validate

# 6. Create deployment PR
> @radius-platform deploy
```

### Multi-Cloud Setup

```bash
# Initialize with both Azure and AWS
> @radius-platform init --org "contoso" --cloud azure,aws --compute kubernetes --iac terraform

# Design with cloud preference
> @radius-platform design --cloud azure

# Generate artifacts for both clouds
> @radius-platform generate --env dev --cloud azure
> @radius-platform generate --env prod --cloud aws
```

### Custom Resource Type

```bash
# Check if MongoDB type exists
> @radius-platform list-types --search mongo

# Create new MongoDB resource type
> @radius-platform create-type --service mongodb --namespace Radius.Data

# Regenerate application with new type
> @radius-platform generate
```

## References

- [Platform Engineering Constitution Schema](../contracts/constitution-schema.yaml)
- [MCP Server Tools](../contracts/mcp-server.yaml)
- [Resource Types Catalog](https://github.com/radius-project/resource-types-contrib)
- [Radius Documentation](https://radapp.io)

## License

MIT License - See LICENSE file for details
