# Application Assembly

* **Author**: Reshma Abdul Rahim (@reshrahim)
* **Date**: June 2026

## Summary

**Application Assembly generates a deployable Radius application graph from source code, removing the need for developers or platform engineers to manually author it.**

Developers have always needed help getting applications to the cloud. Infrastructure-as-code tools like Terraform, Bicep, and CloudFormation made provisioning programmable, but still required developers to think in infrastructure terms. Kubernetes standardized how containerized applications are deployed, and tools like Helm, Kustomize, GitOps workflows, and internal platform templates made Kubernetes more reusable across teams. But these abstractions often shifted complexity to platform specialists, who had to interpret application needs, author deployment assets, and keep everything aligned with organizational standards.

Radius introduced application-centric deployment. Developers describe applications in terms of components, connections, and dependencies, while platform teams define approved infrastructure patterns through Resource Types, Recipes, and Environments. Radius then resolves application intent into repeatable, reliable, enterprise-compliant deployments. But Radius introduced its own adoption barrier: platform engineers still need to teach the Radius resource model, help developers author application graphs, and guide teams through onboarding. The model is right, but the path into it is still too manual.

Application Assembly removes much of that cost. Given a source code repository, the system identifies application components, dependencies, and connections, then represents them as a Radius application graph. The graph is generated from approved platform patterns, so the output reflects organizational standards for provisioning, naming, policy, and deployment readiness rather than arbitrary infrastructure choices. Users can review the graph, understand the deployment shape, confirm or adjust the system’s inferences, and deploy through the approved workflow.

This also opens a path for business application builders: technical users outside core engineering who can build applications with AI-assisted tools but lack cloud deployment expertise. Assembly gives them a governed path to cloud deployment without bypassing organizational standards.

Assembly integrates with the tools developers already use, including GitHub Copilot, VS Code, GitHub Actions, CI/CD pipelines, internal developer portals, and GitOps workflows. The core capability is the same across form factors: generate a deployable Radius application graph from source code. In Project Lattice, that experience is surfaced through the GitHub Copilot app. In self-hosted Radius, it can be surfaced through Radius-native workflows such as the CLI, repository workflows, CI/CD, internal portals, or GitOps.

Application Assembly turns Radius from a model teams must learn first into a model they can adopt through the workflows they already use.

### Goals

1. Lower the barrier to Radius adoption - Developers can start from an existing repository and get a deployable Radius application graph without learning the Radius model first.
2. Make the deployment model understandable - Developers can visually review the generated graph, understand components, dependencies, connections, and deployment assumptions, and adjust the output before use.
3. Make platform standards the default - Generated graphs are built from known application signals and approved platform rules, so they follow the organization’s standards for provisioning, naming, policy, tags, and deployment environments.
4. Turn platform investments into self-service - Resource Types, Recipes, Recipe Packs, and Environments become easier for users to consume without much guidance from platform engineers.

### Non-goals

1. Replacing platform engineers - Platform engineers still own standards, Infrastructure as Code, Environments, policy, and governance. Assembly makes those standards easier to consume; it does not define them independently.
2. Generating raw or unapproved infrastructure code - Assembly generates Radius application definitions that reference existing approved Recipes and Resource Types. It does not create arbitrary Terraform, Bicep, or CloudFormation modules.
3. Making independent infrastructure decisions - Assembly generates definitions based on the platform’s approved patterns, not individual developer preferences or model-generated assumptions.
4. Executing deployment on behalf of the developer - Assembly can validate that the generated application definition is ready for deployment, but it stops before deploying. The developer or existing workflow still triggers deployment through Radius.

## User profile and challenges

### User persona(s)

**Primary user: Application developer**
Builds and owns application code in an engineering team. They understand their service and its dependencies, but are not expected to know the organization’s infrastructure modules, policy rules, naming conventions, or Radius model before getting started.

**Secondary user: Business application builder**
A technical or semi-technical user outside core engineering, such as sales, marketing, operations, field, or internal productivity teams. They can create prototypes, automations, demos, or lightweight apps with AI-assisted tools, but do not have cloud deployment expertise.

**Key stakeholder: Platform engineer**
Defines and maintains the approved deployment path, including Infrastructure as Code, Resource Types, Recipes, Recipe Packs, Environments, policy, naming conventions, and governance.

### Challenge(s) faced by the user

Developers can increasingly create applications, but getting those applications deployed through the approved platform path still requires too much translation. Developers and business application builders need to move from repository or prototype to deployable Radius application graph without learning infrastructure details or the full Radius model. Platform engineers need their approved standards to be consumed consistently without much guidance, manual authoring, or repeated review for every application.

Current options do not solve this problem effectively. Generic AI assistants can generate code or infrastructure, but they are not grounded in the organization’s approved infrastructure abstractions and policies. Templates and examples help with repeatability, but they still require users to know which pattern applies and how to adapt it correctly.

### Positive user outcome

Developers can start from an existing repository or AI-built application and get a Radius application graph they can understand, review, validate, and deploy through the approved workflow. The visual graph helps them understand deployment nuances: what components exist, how dependencies connect, what platform resources will be used, and what must be true before deployment. Platform engineers remain in control of infrastructure standards and policy, while their approved patterns become easier for a broader set of users to consume safely.

## Key scenarios

### Scenario 1: Generate a Radius application graph from an existing repository

A developer starts with an existing application repository. The system inspects source code, Dockerfiles, dependency manifests, configuration files, and local development files to detect application components, dependencies, and connections, then represents them as a Radius application graph using approved platform patterns.

### Scenario 2: Review inferred application structure in a visual experience

The developer sees a clear, visually appealing review experience that shows the application graph: detected components, dependencies, connections, evidence, deployment assumptions, and mappings to Radius Resource Types and Recipes. The experience helps the developer understand how the application will be deployed, including which platform resources will be used, how connections are wired, and what constraints must be satisfied. The developer can accept, modify, or reject each inference before the graph is finalized.

### Scenario 3: Validate deployment readiness before deploy

The system checks the generated application graph against platform rules, approved Resource Types and Recipes, required tags, naming conventions, and supported Environments. The developer gets a clear readiness result before deployment, but deployment is still triggered through Radius or the team’s existing workflow.

### Scenario 4: Make platform standards consumable through existing workflows

Platform engineers define approved infrastructure patterns, policy rules, naming conventions, and environment constraints. The system applies those standards during graph generation and validation, so developers can consume the approved path through the tools they already use, without one-off guidance from platform engineers.

## Key dependencies and risks

### Dependencies

1. **Approved catalog** — A catalog of approved Resource Types and Recipes is required to represent the infrastructure patterns developers are allowed to consume. Without it, the system can detect application needs but cannot produce a useful Radius application graph.

2. **Target Environment configuration** — The system needs access to the target Radius Environment, or enough environment metadata, to validate whether the generated graph is ready for deployment.

3. **Platform standards definition** — Platform rules need to be available in a machine-readable form, including approved Recipes, naming conventions, required tags, supported environments, allowed regions, and policy constraints.

### Risks

1. **Low trust in output** — Users may distrust generated graphs if detections are wrong, unclear, or feel like generic AI-generated definitions. *Mitigation:* Use deterministic generation from approved platform catalog and run benchmark analysis to ensure accuracy.

2. **Platform standards are hard to codify** — Standards may be scattered across wikis, IaC repos, templates, CI checks, and platform engineer knowledge. *Mitigation:* Define a very crisp Radius-specific platform standards format that captures all the basic requirements and allow customizations where needed.

3. **Coverage gaps in Resource Types and Recipes** — The system may detect a dependency without an approved Resource Type or Recipe.
   *Mitigation:* Start with common dependencies, provide starter type definitions (e.g. a Resource Types schema with just the basic connection metadata), and clearly flag unsupported resources to improve catalog coverage.

4. **Deployment-readiness validation is incomplete** — The application graph may pass static checks but still fail deployment because of missing environment configuration, secrets, provider setup, or permissions.
   *Mitigation:* Tight integrations with the deployment systems like Repo Radius or other deployment tools to ensure that all necessary configuration and permissions are in place before deployment.

## Key assumptions to test and questions to answer

1. **The current adoption path is a barrier to Radius usage.** - Users may see value in Radius’s application-centric model, but learning the model and manually translating application needs into a Radius definition creates enough friction to slow adoption.

2. **Application assembly works for all common application patterns.** - The assembly process should be able to handle the most common application architectures and dependencies.

3. **Platform engineers will support self-service generation if they retain control.** - Platform engineers are willing to enable self-service generation if the output stays within approved standards and does not bypass governance.

## Key investments

### Feature 1: Resource Types and Recipe catalog

Expand the Resource Types and Recipes that represent platform-approved provisioning logic for application-centric resources. This includes common dependencies such as databases, caches, message queues, object storage, and identity across priority cloud environments. The goal is to ensure the system can map application needs to supported Radius resource types and generate complete application graphs using validated Recipes.

### Feature 2: Dependency detection and Resource Type mapping

Build the deterministic mapping layer that analyzes the repository structure to identify and map to appropriate Radius Resource Types. Signals may include dependency manifests, Dockerfiles, environment variables, exposed ports, configuration files, and Docker Compose services. The goal is to identify common application dependencies reliably and map them to approved application-centric resources.

### Feature 3: Application graph visualization

Generate a deployable Radius application graph from detected application components, dependencies, connections, matched Resource Types, approved Recipes, and platform constraints. The graph should be presented in a clear visual review experience that shows what was detected, why it was detected, how components connect, what platform resources are selected, and what deployment assumptions were made. Developers should be able to understand the deployment nuances, accept or modify inferences, and review readiness before the graph is finalized.

### Feature 4: Custom Platform standards definition

Define machine-readable platform standards format that captures the platform rules required for generation and validation. This should include approved Resource Types and Recipes, naming conventions, required tags, supported Environments, allowed regions, and policy constraints.

### Feature 5: Validation and deployment readiness checks

Provide validation that checks generated definitions against platform standards and target Environment requirements. This includes static validation for naming, tags, approved resources, and Recipes, plus readiness checks for Environment configuration where available. The system should confirm whether a definition is ready to deploy.

### Feature 6: Developer workflow integration

Make the capability available where devleopers already work, including GitHub Copilot, VS Code, Radius CLI, GitHub Actions, CI/CD pipelines, internal developer portals, and GitOps workflows. The goal is to make application definition generation part of the existing path from repository to deployment, not a separate destination.
