
# Quickstart: Direct Terraform Module as Recipe Template Path

## What This Feature Does

Allows platform engineers to use any standard Terraform module directly as a recipe's `templatePath` in a RecipePack — without wrapping, republishing, or creating Radius-specific artifacts. The system downloads the module fresh at each deployment, passes recipe parameters as Terraform input variables, and surfaces all module outputs through the standard Radius recipe output mechanism.

## Supported Source Formats

| Format | Example | Notes |
|--------|---------|-------|
| Terraform Registry | `hashicorp/consul/aws` | Standard 3-segment format |
| Git (HTTPS) | `git::https://github.com/org/terraform-aws-vpc.git` | Supports `?ref=` for pinning |
| Git (subdirectory) | `git::https://github.com/org/repo.git//modules/vpc` | `//` separates repo from subdir |
| Git (SSH) | `git::ssh://git@github.com/org/module.git` | Requires SSH key configuration |
| HTTP Archive | `https://example.com/modules/vpc.tar.gz` | Tar/zip archives |

## Usage Examples

### Example 1: Use Public Registry and Git-Hosted Modules in a RecipePack

```bicep
resource recipePack 'Radius.Core/recipePacks@2025-08-01-preview' = {
  name: 'datastore-recipes'
  properties: {
    recipes: {
      'Radius.Data/mySqlDatabases': {
        recipeKind: 'terraform'
        recipeLocation: 'Azure/mysql/azurerm'
        parameters: {
          server_name: 'my-mysql-server'
          sku_name: 'GP_Gen5_2'
          storage_mb: 5120
          administrator_login: 'mysqladmin'
        }
      }
      'Radius.Data/postgreSqlDatabases': {
        recipeKind: 'terraform'
        recipeLocation: 'git::https://github.com/terraform-aws-modules/terraform-aws-rds.git?ref=v6.1.0'
        parameters: {
          identifier: 'my-postgres-db'
          engine: 'postgres'
          engine_version: '15.4'
          instance_class: 'db.t3.micro'
          allocated_storage: 20
        }
      }
    }
  }
}
```

### Example 2: Version Pinning via Environment Recipe

> **Note:** `templateVersion` is available on Environment-level recipe definitions, not on RecipePack `RecipeDefinition`. To pin a version in a RecipePack, embed the version in the `recipeLocation` URL (e.g., `?ref=v3.1.0` for Git sources). For registry modules, use Environment-level recipes:

```bicep
resource env 'Radius.Core/environments@2025-08-01-preview' = {
  name: 'my-env'
  properties: {
    compute: {
      kind: 'kubernetes'
      namespace: 'default'
    }
    recipes: {
      'Radius.Data/mySqlDatabases': {
        mysql: {
          templateKind: 'terraform'
          templatePath: 'Azure/mysql/azurerm'
          templateVersion: '3.1.0'
          parameters: {
            server_name: 'my-mysql'
            sku_name: 'B_Gen5_1'
            storage_mb: 5120
            administrator_login: 'mysqladmin'
          }
        }
      }
    }
  }
}
```

### Example 3: Private Git Repository with Authentication

Configure authentication in the environment's recipe configuration:

```json
{
  "recipeConfig": {
    "terraform": {
      "authentication": {
        "git": {
          "pat": {
            "github.com": {
              "secret": "/planes/radius/local/resourceGroups/default/providers/Radius.Core/secretStores/github-pat"
            }
          }
        }
      }
    }
  }
}
```

Then reference the private module:

```json
{
  "recipeLocation": "git::https://github.com/my-org/private-terraform-modules.git//modules/mysql"
}
```

## How Outputs Work

For direct Terraform modules, ALL module outputs are available as recipe outputs:

- **Non-sensitive outputs** → available in the resource's `Values` map
- **Sensitive outputs** (marked `sensitive = true` in the module) → stored in `Secrets`

Example: If your Terraform MySQL module declares:
```hcl
output "host" {
  value = kubernetes_service.mysql.metadata[0].name
}

output "port" {
  value = 3306
}

output "connection_string" {
  value     = "mysql://${var.username}:${var.password}@${kubernetes_service.mysql.metadata[0].name}:3306/${var.database}"
  sensitive = true
}
```

After deployment:
- `resource.properties.values.host` → `"mysql-service"`
- `resource.properties.values.port` → `3306`
- `resource.properties.secrets.connection_string` → stored securely, accessible via secret reference

## How Parameters Work

Recipe parameters map directly to Terraform input variables:

1. **Environment-level parameters** (in RecipePack `parameters`) → applied to every deployment
2. **Resource-level parameters** (at deploy time) → override environment-level for same keys
3. **No parameter** for optional variables → Terraform uses the module's default value
4. **Missing required variable** → Terraform error surfaced through recipe failure

## Inspecting Module Schema

Use the existing `GetRecipeMetadata` API to discover a module's inputs before deploying:

```bash
rad recipe show --name my-recipe --environment my-env
```

This returns the module's input variables (name, type, description, required/optional, defaults) and outputs.

## Key Behaviors

| Behavior | Details |
|----------|---------|
| **Module download** | Fresh download every deployment (no caching) |
| **Version pinning** | Use `templateVersion` on Environment-level recipes for registry modules, `?ref=` for Git |
| **Provider config** | Uses existing `recipeConfig.terraform.providers` from environment |
| **State management** | Same Kubernetes secret backend as existing recipes |
| **Error handling** | Terraform errors surfaced directly in recipe failure response |
| **Existing recipes** | Zero behavioral changes — fully backward compatible |
| **Fallback** | If direct resolution fails, falls back to OCI/wrapped path |

## Development Workflow

### Building

```bash
make build
```

### Running Tests

```bash
# Unit tests for the source resolver
go test ./pkg/recipes/source/...

# Unit tests for terraform driver changes
go test ./pkg/recipes/driver/terraform/...

# Unit tests for terraform executor changes
go test ./pkg/recipes/terraform/...

# All unit tests
make test
```

### Linting

```bash
make lint
make format-check
```
