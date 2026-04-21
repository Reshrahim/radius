# Data Model: Developer Workflow in Copilot and GitHub Using Radius

**Branch**: `004-copilot-radius-dev-workflow` | **Date**: 2026-02-20

## Entities

### 1. PlatformConstitution

The authoritative document capturing an organization's platform requirements. Stored as `Platform-Engineering-Constitution.md` at the repository root.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| organization | string | yes | Organization name |
| version | semver string | yes | Constitution version (e.g., "1.0.0") |
| cloudProviders | CloudProvider[] | yes | Configured cloud providers (min 1) |
| computeTargets | ComputeTarget[] | yes | Deployment targets (e.g., Kubernetes, ACA) |
| iacTooling | IaCTooling | yes | Preferred Infrastructure-as-Code tooling |
| moduleCatalog | ModuleCatalog | yes | Approved IaC module catalog location |
| containerRegistry | ContainerRegistry | no | Approved container registry configuration |
| namingConvention | NamingConvention | no | Resource naming pattern and rules |
| labels | Label[] | no | Required labels/tags for resources |
| policies | Policy[] | no | Platform policies (manual or imported) |
| repositoryPolicies | RepositoryPolicy | no | GitHub repository configuration policies |
| containerImagePolicies | ContainerImagePolicy | no | Container image build/publish policies |
| authenticationPreferences | string[] | no | Ordered authentication method preferences |

**Validation rules**:
- `cloudProviders` must contain at least one entry
- `computeTargets` must contain at least one entry
- `version` must be valid semver
- If `containerImagePolicies.scanning` is enabled, `containerRegistry` must be set
- Policy `severity` must be one of: `error`, `warning`, `info`

**State transitions**: Draft → Active → Updated (via `rad init` re-run preserving existing settings)

### 2. CloudProvider

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | enum: azure, aws, gcp | yes | Cloud provider identifier |
| subscriptionId | string | conditional | Azure subscription (required if azure) |
| region | string | yes | Default deployment region |
| resourceGroup | string | conditional | Azure resource group (required if azure) |
| accountId | string | conditional | AWS account ID (required if aws) |

### 3. ComputeTarget

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| type | enum: kubernetes, aca, ecs, appservice | yes | Compute platform type |
| name | string | yes | Human-readable target name |
| context | string | conditional | Kubernetes context (required if kubernetes) |
| namespace | string | no | Default namespace |

### 4. IaCTooling

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| preferred | enum: terraform, bicep | yes | Preferred IaC language for recipes |
| terraformVersion | string | no | Minimum Terraform version |
| bicepVersion | string | no | Minimum Bicep version |

### 5. ModuleCatalog

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| type | enum: terraform-registry, git, local | yes | Catalog source type |
| url | string | yes | Catalog URL or path |
| authMethod | string | no | Authentication method for private catalogs |
| namespace | string | no | Module namespace filter |

### 6. Policy

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | yes | Policy name (unique within constitution) |
| source | enum: manual, azure-policy, aws-config, opa | yes | Policy origin |
| importPath | string | conditional | Path/URL for imported policies |
| rule | string | yes | Policy rule description |
| severity | enum: error, warning, info | yes | Violation severity |
| remediation | string | no | Suggested fix for violations |
| category | string | no | Policy category (security, naming, cost, etc.) |

### 7. ApplicationGraph

The structured representation of an application's topology. Stored as `radius/app-graph.json`.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| version | string | yes | Schema version (e.g., "1.0") |
| application | ApplicationMeta | yes | Application metadata |
| nodes | GraphNode[] | yes | Services and infrastructure components |
| edges | GraphEdge[] | yes | Connections between nodes |
| metadata | GraphMetadata | no | Generation metadata (timestamp, tool version) |

### 8. GraphNode

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | string | yes | Unique node identifier (e.g., "frontend") |
| type | string | yes | Radius resource type (e.g., "Applications.Core/containers") |
| name | string | yes | Human-readable name |
| source | enum: local, community, generated | no | Resource type resolution source |
| recipe | string | no | Recipe name (for infrastructure nodes) |
| cloudService | string | no | Mapped cloud service (e.g., "Azure Cache for Redis") |
| module | ModuleReference | no | IaC module reference |
| properties | map[string]any | no | Node-specific properties |

### 9. GraphEdge

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| from | string | yes | Source node ID |
| to | string | yes | Target node ID |
| type | enum: connection, dependency, secret | yes | Edge type |
| port | int | no | Connection port |

### 10. ModuleReference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| source | string | yes | Module source (registry URL or path) |
| name | string | yes | Module name |
| version | string | yes | Module version |

### 11. ApplicationDefinition

The Bicep file declaring all application services and infrastructure. Stored as `radius/app.bicep`.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| applicationName | string | yes | Radius application name |
| services | ServiceDefinition[] | yes | Application services |
| connections | ConnectionDefinition[] | no | Inter-service connections |
| secretBindings | SecretBinding[] | no | Secret references |

### 12. ServiceDefinition

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | yes | Service name |
| resourceType | string | yes | Radius resource type |
| image | string | conditional | Container image (for compute services) |
| ports | PortMapping[] | no | Exposed ports |
| environment | map[string]string | no | Environment variables |

### 13. RecipeConfiguration

IaC module wrapper for infrastructure provisioning. Stored in `radius/recipes/`.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | yes | Recipe name |
| resourceType | string | yes | Target Radius resource type |
| templateKind | enum: terraform, bicep | yes | IaC language |
| modulePath | string | yes | Module source path |
| moduleVersion | string | yes | Module version |
| parameters | map[string]any | no | Default parameter values |
| secureByDefault | bool | yes | Whether secure defaults are applied |

### 14. EnvironmentConfiguration

Target deployment environment definition. Stored in `radius/environments/`.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | yes | Environment name |
| computeTarget | ComputeTarget | yes | Target compute platform |
| cloudProvider | CloudProvider | yes | Target cloud provider |
| namespace | string | no | Kubernetes namespace |
| recipes | RecipeRegistration[] | yes | Registered recipes |

### 15. DeploymentPlan

Pre-deployment summary generated by Repo Radius.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | string | yes | Unique plan identifier |
| timestamp | datetime | yes | Plan generation time |
| applicationName | string | yes | Target application |
| environment | string | yes | Target environment |
| actions | DeploymentAction[] | yes | Resources to create/modify/delete |
| policyResults | PolicyResult[] | yes | Constitution compliance results |
| status | enum: pending, approved, rejected, executed, failed | yes | Plan state |

### 16. DeploymentAction

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| resourceName | string | yes | Resource identifier |
| resourceType | string | yes | Radius resource type |
| action | enum: create, update, delete, no-change | yes | Planned action |
| properties | map[string]any | no | Resource properties |
| diff | PropertyDiff[] | no | Changed properties (for updates) |

### 17. PolicyResult

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| policyName | string | yes | Constitution policy name |
| status | enum: pass, fail, warning, skip | yes | Evaluation result |
| resource | string | no | Affected resource (if specific) |
| message | string | no | Result description |
| remediation | string | no | Suggested fix (for failures) |

### 18. MCPToolDefinition

Metadata for an MCP server tool exposed to Copilot.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | yes | Tool name (e.g., "resolve_resource_types") |
| description | string | yes | Human-readable tool description |
| inputSchema | JSONSchema | yes | JSON Schema for tool input parameters |
| outputSchema | JSONSchema | yes | JSON Schema for tool output |

## Entity Relationships

```
PlatformConstitution
├── CloudProvider[] ─────────────┐
├── ComputeTarget[] ─────────────┤
├── Policy[] ────────────────────┤
├── ModuleCatalog ─── validates ──> RecipeConfiguration.modulePath
└── ContainerImagePolicy         │
                                 │
ApplicationGraph                 │
├── GraphNode[] ─── typed by ──> ResourceType (from resolution)
│   ├── recipe ──── references ─> RecipeConfiguration
│   └── module ──── sourced from > ModuleCatalog
├── GraphEdge[] ─── connects ──> GraphNode pairs
└── generates ─────────────────> ApplicationDefinition (app.bicep)
                                 EnvironmentConfiguration
                                 RecipeConfiguration[]
                                 
DeploymentPlan
├── compares ──> ApplicationGraph (current vs. proposed)
├── validates ─> PlatformConstitution.policies
└── produces ──> DeploymentAction[]
                 PolicyResult[]
```

## State Machine: Deployment Lifecycle

```
[PR Created] → [Architecture Visualization Added] → [PR Reviewed] → [PR Merged]
    → [Deployment Plan Generated] → [Plan Reviewed]
        → [Approved] → [Deploying] → [Success] / [Failed]
        → [Rejected] → [Cancelled]
```
