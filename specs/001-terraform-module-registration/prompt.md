# Feature: Direct Terraform Module as Recipe Template Path

Enable platform engineers to use an existing Terraform module directly as the `templatePath` for a recipe in a RecipePack — without having to wrap the module, publish it to a separate registry, or create a Radius-specific recipe artifact.

Today, to use an existing Terraform module with Radius recipes, engineers must wrap the module into a Radius recipe format, publish it, and then reference it via `templatePath`. This feels like double work. The feature should let the recipe's `templatePath` point directly at an existing Terraform module (from a Terraform registry, Git repo, or local path) and have it work as-is — no wrapping, no republishing.

Now the other thing to work through, how do we match the module's input variables to the resource properties and the module's outputs to the computed properties in Radius? This is a key part of the feature because it allows the Terraform module to be used seamlessly within Radius without requiring manual mapping or additional configuration. The system should automatically discover the module's input variables and outputs, and then map them to the resource type's properties.
