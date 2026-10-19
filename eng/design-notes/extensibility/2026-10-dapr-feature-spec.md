# Topic: Dapr Support in Radius

- **Author**: Reshma Abdul Rahim (@Reshrahim)

## Terms and definitions

| Term               | Definition                                                                                                                                                      |
|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Dapr**           | [Distributed Application Runtime](https://dapr.io). A sidecar that gives applications HTTP and gRPC APIs for state, messaging, secrets, and service invocation. |
| **Dapr component** | A YAML or ARM object telling the sidecar which infrastructure backs a building block. Has a name, a type such as `state.redis`, and metadata.                   |
| **Building block** | A Dapr API such as state management or publish/subscribe, backed by one or more components.                                                                     |
| **App ID**         | The identifier a sidecar uses for the application it serves. Used for service invocation and component scoping.                                                 |
| **Recipe Pack**    | A bundle of Recipes registered on an Environment as one unit. See [Recipe Packs](../recipes/2025-08-recipe-packs.md).                                           |

## Topic Summary

Radius supports Dapr today through four hard-coded resource types in `Applications.Dapr`. Provisioning is Go code in `pkg/daprrp`, it targets only Kubernetes, and developers write Dapr component internals in their application file.

This spec sets the direction to replace that with Recipe-backed resource types in a new `Radius.Dapr` namespace, consistent with [compute extensibility](2025-06-compute-extensibility-feature-spec.md).

The work is three things:

1. **Resource types** in the `Radius.Dapr` namespace.
2. **Recipes** that create the Dapr component in the form each compute platform expects.
3. **Dapr Recipe Packs** that bundle those Recipes. MVP ships two: Kubernetes and Azure Container Apps.

Recipes and Packs live in [resource-types-contrib](https://github.com/radius-project/resource-types-contrib).

### Top level goals

- Let developers use Dapr building blocks without writing component YAML or naming backing infrastructure.
- Let platform engineers choose the backing infrastructure for each Environment by registering a Recipe Pack.
- Make Dapr work on compute platforms other than Kubernetes by moving platform-specific logic into Recipes.
- Keep the Dapr component name predictable so application code does not change between Environments.
- Stop adding Dapr provisioning logic to Applications RP, so that Dapr support grows by writing Recipes rather than Go code.

### Non-goals (out of scope)

- **Installing Dapr.** The platform engineer installs the Dapr control plane before registering Dapr Recipes.
- **Building blocks deferred past MVP.** Configuration stores, distributed lock, cryptography, and the conversation API. Workflows are deferred but expected next, so MVP must not block them. See Building block scope.
- **Dapr policy configuration.** Dapr `Configuration`, `Resiliency`, and mTLS settings belong to the platform engineer and stay outside the developer schema. Declarative `Subscription` resources are also out of scope, so subscribers register topics programmatically, which is how Dapr applications most commonly do it.
- **Changing Dapr itself.** Radius consumes upstream Dapr.
- **Replacing the Dapr CLI or dashboard.** Developers keep using `dapr` tooling locally.
- **Scale to zero.** A container autoscaling concern. See Appendix 3.
- **ECS Recipes.** The mechanism is in Appendix 2 so the schema does not block it, but no ECS Recipe ships in MVP.

## User profile and challenges

### User persona(s)

**Application developer.** Writes services that call the Dapr API. Wants state and messaging without learning Redis, Azure Service Bus, or AWS SNS.

**Platform engineer.** Owns the Dapr installation, the backing infrastructure, and the security posture. Decides that development uses Redis and production uses managed cloud services, without editing application files.

### Challenge(s) faced by the user

- **Developers write infrastructure details today.** `Applications.Dapr/stateStores` requires `type`, `version`, and a `metadata` map holding a host and a password reference. This is the knowledge Dapr exists to hide.
- **The application file cannot move between Environments.** Component type and metadata live in the application file, so changing the backing store means editing the application.
- **Dapr only works on Kubernetes in Radius.** There is no path to Azure Container Apps, which has Dapr built in, or to ECS.
- **Platform engineers have no enforcement point.** A developer can set any component type and any metadata.
- **Two mental models.** Dapr uses `Applications.Dapr` while everything else moves to Recipe-backed `Radius.*` types.

### Positive user outcome

Developers declare the Dapr building blocks they need and nothing more. Platform engineers own the backing infrastructure for every team through Recipes, on any platform where Dapr runs.

### Building block scope

Scope follows what Dapr users run in production. A user running Dapr on Azure Container Apps reported active use of service invocation, publish/subscribe, state management, bindings, actors, and secret management, and is evaluating workflows next. All six are long-established building blocks.

| Building block                                              | In MVP | What it needs from Radius                                                       |
|-------------------------------------------------------------|--------|---------------------------------------------------------------------------------|
| Service invocation                                          | Yes    | Sidecar with an app ID. No resource type.                                       |
| State management                                            | Yes    | `Radius.Dapr/stateStores`                                                       |
| Publish/subscribe                                           | Yes    | `Radius.Dapr/pubSubBrokers`                                                     |
| Secret management                                           | Yes    | `Radius.Dapr/secretStores`                                                      |
| Bindings                                                    | Yes    | `Radius.Dapr/bindings`                                                          |
| Actors                                                      | Yes    | A state store marked as the actor state store, plus the Dapr placement service. |
| Workflows                                                   | No     | Nothing new. Workflows build on actors and a state store.                       |
| Configuration, distributed lock, cryptography, conversation | No     | New resource types later, with no change to the MVP types.                      |

Deferring bindings and actors would leave that user unable to adopt Radius, because a partial Dapr surface forces them to keep their existing deployment tooling alongside it.

## Key scenarios

### Scenario 1: Developer uses state management

A developer declares a state store, connects a container to it, enables the Dapr sidecar, and calls the Dapr state API. They do not know which backing store the Environment uses.

### Scenario 2: Developer uses publish/subscribe between two services

A developer declares a pub/sub broker and connects a publisher and a subscriber container to it. Both containers run a sidecar.

### Scenario 3: Platform engineer changes backing infrastructure per Environment

A platform engineer registers a Dapr Recipe Pack that uses in-cluster infrastructure for development and managed cloud services for production. Applications deploy to both unchanged.

### Scenario 4: Platform engineer runs the same applications on Azure Container Apps

A platform engineer registers the ACA Dapr Recipe Pack. Applications written against the Kubernetes Environment deploy unchanged.

### Scenario 5: Developer uses bindings and actors

A developer declares a binding to exchange events with an external system and marks a state store as the actor state store. The Recipe maps each to the right component for the platform.

## Key dependencies and risks

- **Dependency: compute extensibility and the containers Recipe.** Sidecar enablement is done by the containers Recipe for each platform, so a Dapr Recipe Pack is only useful once that platform has a containers Recipe. With ACA in MVP this is a hard gate, not a sequencing preference: the ACA Dapr Pack cannot ship before the ACA containers Recipe. Mitigation: the Kubernetes Pack is independent and ships first, and the schema does not change when the ACA Pack follows.
- **Dependency: Dapr installed on the target platform.** On Kubernetes the platform engineer installs it. On Azure Container Apps it is part of the managed environment. Mitigation: each Recipe checks the precondition its platform has, the Dapr `Component` CRD on Kubernetes and the managed environment on Azure Container Apps, and fails with a message naming the missing dependency.
- **Dependency: Dapr placement service for actors.** A standard Kubernetes install and the ACA managed environment both provide it. On ECS the platform engineer runs it.
- **Risk: component scoping, and what it means for actors on ACA.** A Dapr component is reachable by every Dapr application sharing its scope unless the component sets `scopes`. Radius cannot derive `scopes` from connections in MVP, because the Dapr resource is provisioned before the containers that connect to it. On Kubernetes the blast radius is a single application, because Radius creates one namespace per application. On Azure Container Apps every application in a managed environment shares one scope, so two applications that each mark a state store as the actor state store would place two actor state stores in the same sidecar's view, which Dapr does not accept. Because ACA and actors are both in MVP, this is an MVP work item rather than a follow-up: either scoping is solved, or the ACA Pack ships with actors unsupported and says so. Options are in the open questions below.
- **Risk: secret handling differs per platform.** A Kubernetes Dapr `Component` resolves `secretKeyRef` through a Dapr secret store, either the default Kubernetes one or one named on the component. Azure Container Apps `daprComponents` do not use `secretKeyRef`; they carry a `secrets` array referenced with `secretRef`. Mitigation: the Recipe owns this per platform. The Dapr resource types expose no `secrets` property, because only the sidecar reads the credential, never the application.
- **Risk: migration.** `Applications.Dapr` users have `type` and `metadata` in their application files and those properties do not exist in the new schema. Mitigation: keep `Applications.Dapr` working until the new types reach GA and publish the mapping in Appendix 1.

## Key assumptions to test and questions to answer

- **Assumption:** developers do not need to set the component type or metadata once a Recipe owns them. Validation: review the metadata properties current `Applications.Dapr` users set and confirm each is a Recipe decision or a Recipe parameter.
- **Assumption:** one Recipe per building block per platform is enough. Validation: an Environment can register multiple named Recipes for one type, so a developer selects one if this is wrong.
- **Question:** is a platform-neutral `kind` on bindings workable, or does it collapse into a passthrough of the Dapr component type? Validation: add a third kind and confirm two platforms implement it without a schema change.
- **Question:** how do we scope components to the applications that connect to them? This blocks actors on Azure Container Apps in MVP. Options: the containers Recipe patches `scopes` after the component exists, Radius passes connected app IDs into the Recipe, or the Recipe sets `scopes` from a parameter the platform engineer supplies.
- **Question:** Dapr accepts only one state store flagged as the actor state store within a sidecar's scope. Does one actor state store per application hold up, or do teams need more than one? Validation: confirm with an application using several actor types.
- **Question:** should the component name stay equal to the resource name, or become writable? MVP uses the resource name because it is predictable from the application file.

## Current state

Radius has four Dapr resource types in `Applications.Dapr`: `stateStores`, `pubSubBrokers`, `secretStores`, and `configurationStores`. All share the schema in `typespec/Applications.Dapr/common.tsp`: `componentName` (read-only), `type`, `version`, `metadata`, and `auth.secretStore`. They support `manual` provisioning, where the developer supplies the component type and metadata, and `recipe` provisioning.

Provisioning is Go code under `pkg/daprrp`. The processors build a Kubernetes `Component` in the application namespace. `pkg/daprrp/datamodel/daprdetector.go` checks for the `components.dapr.io` CRD and is Kubernetes-only.

The sidecar is enabled by the `daprSidecar` extension on containers. `pkg/corerp/renderers/daprextension/renderer.go` writes the `dapr.io/enabled`, `dapr.io/app-id`, `dapr.io/app-port`, and `dapr.io/protocol` annotations onto the Deployment. The `Radius.Compute/containers` manifest already carries `extensions.daprSidecar` with `appId`, `appPort`, and `config`.

## Details of user problem

A developer picks Dapr so their code does not depend on a specific broker or database. Modeling that service in Radius today requires writing the thing Dapr hides:

```bicep
resource stateStore 'Applications.Dapr/stateStores@2023-10-01-preview' = {
  name: 'orders'
  properties: {
    environment: environment
    application: app.id
    resourceProvisioning: 'manual'
    type: 'state.redis'
    version: 'v1'
    metadata: {
      redisHost: { value: '${redis.properties.host}:${redis.properties.port}' }
      redisPassword: { secretKeyRef: { name: redis.listSecrets().password, key: 'password' } }
    }
  }
}
```

The application file now names Redis, knows its port, and knows how its authentication works. That is why it cannot move between Environments or to a platform Radius does not emit components for.

## Desired user experience outcome

The component type, version, metadata, and credentials all move to the Recipe. The application file carries no infrastructure knowledge, so it does not change between Environments or compute platforms, and every Dapr component in an Environment comes from a Recipe the platform engineer controls.

### Detailed user experience

The Recipe sources below are illustrative. Final module paths are set in `resource-types-contrib`.

#### Platform engineer: prepare the Environment

**Step 1. Install Dapr on the platform.** On Kubernetes, the Dapr Helm chart or `dapr init -k`. On Azure Container Apps, Dapr is part of the managed environment and nothing is installed.

**Step 2. Define the Dapr Recipe Pack.** One pack per platform, covering all four types, so Dapr is configured in one step:

```bicep
extension radius

resource daprRecipes 'Radius.Core/recipePacks@2025-08-01-preview' = {
  name: 'dapr-kubernetes'
  properties: {
    recipes: {
      'Radius.Dapr/stateStores': {
        kind: 'terraform'
        source: 'git::https://github.com/radius-project/resource-types-contrib.git//dapr/statestores/kubernetes'
        parameters: {
          metadataOverrides: {
            enableTLS: 'true'
          }
        }
      }
      'Radius.Dapr/pubSubBrokers': {
        kind: 'terraform'
        source: 'git::https://github.com/radius-project/resource-types-contrib.git//dapr/pubsubbrokers/kubernetes'
      }
      'Radius.Dapr/secretStores': {
        kind: 'terraform'
        source: 'git::https://github.com/radius-project/resource-types-contrib.git//dapr/secretstores/kubernetes'
      }
      'Radius.Dapr/bindings': {
        kind: 'terraform'
        source: 'git::https://github.com/radius-project/resource-types-contrib.git//dapr/bindings/kubernetes'
      }
    }
  }
}
```

**Step 3. Attach the pack to the Environment** alongside the compute pack, since sidecar enablement belongs to the containers Recipe:

```bash
rad env update myenv --recipe-packs kubernetes-compute,dapr-kubernetes --preview
```

**Step 4. Set Recipe parameters** where a default is wrong for the organization. The example above turns on TLS to the Redis state store, which Dapr leaves off by default and which a security review would require. Component scoping is not settled in MVP, so restricting which app IDs may use a component is an open question rather than a parameter to reach for here. The Dapr `Configuration` applied to every sidecar is a parameter on the compute pack, because it configures the sidecar rather than a component.

#### Developer: declare and deploy the application

**Step 1. Write code against the Dapr API**, using the resource names as component names.

**Step 2. Declare the resources.** The developer states what the application needs. No component types, no versions, no metadata:

```bicep
extension radius

param environment string

resource app 'Radius.Core/applications@2025-08-01-preview' = {
  name: 'checkout'
  properties: { environment: environment }
}

resource orders 'Radius.Dapr/stateStores@2026-10-01-preview' = {
  name: 'orders'
  properties: {
    environment: environment
    application: app.id
  }
}

resource events 'Radius.Dapr/pubSubBrokers@2026-10-01-preview' = {
  name: 'events'
  properties: {
    environment: environment
    application: app.id
  }
}

resource receipts 'Radius.Dapr/bindings@2026-10-01-preview' = {
  name: 'receipts'
  properties: {
    environment: environment
    application: app.id
    direction: 'output'
    kind: 'queue'
  }
}

resource frontend 'Radius.Compute/containers@2025-08-01-preview' = {
  name: 'frontend'
  properties: {
    environment: environment
    application: app.id
    containers: {
      frontend: {
        image: 'ghcr.io/myorg/frontend:latest'
        ports: { web: { containerPort: 8080 } }
      }
    }
    connections: {
      orders: { source: orders.id }
      events: { source: events.id }
      receipts: { source: receipts.id }
    }
    extensions: {
      daprSidecar: {
        appId: 'frontend'
        appPort: 8080
        appProtocol: 'http'
      }
    }
  }
}
```

**Step 3. Deploy.** `rad deploy app.bicep` runs the Dapr Recipes, then the containers Recipe. The Recipes provision the backing infrastructure and create the Dapr components named `orders`, `events`, and `receipts`.

**Step 4. Read the component names from the environment.** Each connection injects `CONNECTION_<NAME>_<PROPERTY>` into every container in the resource, so the container receives `CONNECTION_ORDERS_COMPONENTNAME=orders`, `CONNECTION_EVENTS_COMPONENTNAME=events`, and `CONNECTION_RECEIPTS_COMPONENTNAME=receipts`. Code can read the component name instead of hardcoding it. This uses the existing connection mechanism, which the containers Recipe implements; `componentName` is a normal read-only property and needs no Dapr-specific code.

Input bindings work the other way round: Dapr calls the application at a route named after the component, so the name is part of the application's routing rather than an argument to a Dapr call. A `cron` binding with `direction: 'input'` is declared and connected the same way, and the connection records the dependency in the graph.

**Step 5. Inspect the application.** `rad app graph` shows `frontend` connected to `orders`, `events`, and `receipts`, because a connection to a Dapr resource is an ordinary dependency in the graph.

#### Developer: use actors

The developer marks the state store that holds actor state. Nothing else changes:

```bicep
resource orders 'Radius.Dapr/stateStores@2026-10-01-preview' = {
  name: 'orders'
  properties: {
    environment: environment
    application: app.id
    actorStateStore: true
  }
}
```

The Recipe sets the component metadata Dapr requires for actor state and must select a backing store that supports it. The platform engineer is responsible for the placement service being reachable.

#### Developer: deploy to Azure Container Apps

`rad deploy app.bicep -e aca-prod`. The application file above is unchanged, including the binding and the `daprSidecar` extension.

What differs is entirely inside the Recipes. The containers Recipe sets `properties.configuration.dapr` on the Container App instead of writing pod annotations. The Dapr Recipes create `Microsoft.App/managedEnvironments/daprComponents` named `orders`, `events`, and `receipts` on the managed environment, backed by Cosmos DB and Azure Service Bus instead of Redis. The component names the application calls are the same, so no code changes.

The developer's only requirement is that the target Environment has a Dapr Recipe Pack registered. If it does not, deployment fails when the first Dapr resource has no Recipe for its type.

## Key investments

### Feature 1: `Radius.Dapr` resource types

Four Recipe-backed resource types sharing a common base. The developer declares the need; the Recipe owns the component type, version, and metadata.

**Common properties**

| Property        | Type   | Required  | Description                                                                                |
|-----------------|--------|-----------|--------------------------------------------------------------------------------------------|
| `environment`   | string | Yes       | The Radius Environment ID. Set by the `rad` CLI.                                           |
| `application`   | string | Yes       | The Radius Application ID. Dapr components are application-scoped in MVP.                  |
| `componentName` | string | Read-only | The name the application passes to the Dapr API. Always equal to the Radius resource name. |

`Radius.Dapr/pubSubBrokers` and `Radius.Dapr/secretStores` use the base unchanged. Neither adds a developer property, because every remaining knob is a platform decision.

`Radius.Dapr/secretStores` does not overlap `Radius.Security/secrets`. `Radius.Security/secrets` holds secret values supplied by the developer and handed to other Radius resources at deployment time. A Dapr secret store is a runtime API: the application asks its sidecar for a secret by name and the sidecar reads it from the backing vault, so values are never in the application file and can change without redeploying. An application can use both.

**`Radius.Dapr/stateStores`** adds one property for actors:

| Property          | Type    | Required | Description                                                                                  |
|-------------------|---------|----------|----------------------------------------------------------------------------------------------|
| `actorStateStore` | boolean | No       | Set when the application uses Dapr actors and this store holds actor state. Default `false`. |

Actors need no resource type. Dapr requires that exactly one state store component within a sidecar's scope is flagged as the actor state store. The flag is a developer property because only the application knows it uses actors; choosing a backing store that supports actor state stays with the platform engineer.

**`Radius.Dapr/bindings`** adds two properties:

| Property    | Type   | Required | Description                                                                                |
|-------------|--------|----------|--------------------------------------------------------------------------------------------|
| `direction` | string | Yes      | `input` delivers events from the external system to the application. `output` sends to it. |
| `kind`      | string | Yes      | The kind of external system in platform-neutral terms, such as `cron` or `queue`.          |

Bindings are the one building block where the developer must say more than the resource name. A scheduled trigger and an external queue are different application requirements, not interchangeable infrastructure, so the Recipe cannot choose between them. The Recipe maps `kind` to a concrete component type per platform, so platform engineers add kinds without a schema change.

Not every pair is valid: `cron` is input only. Valid pairs belong to the Recipe contract rather than the schema, and a Recipe rejects an unsupported pair with a message naming it.

### Feature 2: Dapr sidecar configuration on containers

`extensions.daprSidecar` on `Radius.Compute/containers` is the single way a developer asks for a sidecar. Changes from the current manifest:

- **`appId` becomes required** when `daprSidecar` is set. A sidecar without an app ID cannot be addressed. It must be unique among Dapr applications sharing a scope: the application namespace on Kubernetes, the managed environment on Azure Container Apps. Radius does not enforce uniqueness in MVP.
- **`appProtocol` is added**, `http` or `grpc`, default `http`. `Applications.Core` has this setting today but calls it `protocol`; the new name matches `appId` and `appPort` and Dapr's own `app-protocol`. gRPC applications need it, and it is absent from the `Radius.Compute/containers` manifest.
- **`config` is removed** from the developer schema. The Dapr `Configuration` resource holds tracing, mTLS, and middleware settings, which are platform decisions. It becomes a containers Recipe parameter so it applies consistently across the Environment.

### Feature 3: Recipes and Dapr Recipe Packs

Each Recipe does three things: provision or reference the backing infrastructure, create the Dapr component in the form the platform expects, and output `componentName`.

**MVP ships two Packs, Kubernetes and Azure Container Apps**, each covering all four types. Shipping both from the start is what proves the design: a resource type that only ever has one Recipe has not been shown to be platform-neutral.

Backing infrastructure is a Recipe decision and is not fixed by this spec. The table records the expected starting point and, more importantly, the one constraint actors place on it:

| Type                        | Kubernetes                           | Azure Container Apps                               |
|-----------------------------|--------------------------------------|----------------------------------------------------|
| `Radius.Dapr/stateStores`   | `state.redis`                        | `state.azure.cosmosdb`                             |
| `Radius.Dapr/pubSubBrokers` | `pubsub.redis`                       | `pubsub.azure.servicebus.topics`                   |
| `Radius.Dapr/secretStores`  | `secretstores.kubernetes`            | `secretstores.azure.keyvault`                      |
| `Radius.Dapr/bindings`      | `bindings.cron`, `bindings.rabbitmq` | `bindings.cron`, `bindings.azure.servicebusqueues` |

The state store Recipe must use a store that supports Dapr actor state, which requires transactions and ETag. Redis qualifies on Kubernetes. On Azure Container Apps, **Azure Table Storage does not support actors because it has no transaction support**, so the ACA Recipe uses Cosmos DB. Picking Table Storage would deploy successfully and fail only when an application used actors.

Recipe parameters give the platform engineer control without forking: at minimum the component version and a metadata override map for tuning. Whether `scopes` is also a parameter depends on how component scoping is resolved, which is an open question.

### Feature 4: Retire `Applications.Dapr`

`Applications.Dapr/stateStores`, `pubSubBrokers`, and `secretStores` are deprecated when `Radius.Dapr` reaches GA, following the same path as other `Applications.*` to `Radius.*` moves. Migration guidance is in Appendix 1.

`Applications.Dapr/configurationStores` has no replacement in MVP and stays until `Radius.Dapr/configurationStores` ships. Because of that, `pkg/daprrp` is not deleted in MVP. The Kubernetes-only Dapr detection in `pkg/daprrp/datamodel/daprdetector.go` is replaced by the equivalent check in the Recipes for the new types.

## Appendix 1: Migration from `Applications.Dapr`

| `Applications.Dapr` property   | `Radius.Dapr` equivalent                                                                        |
|--------------------------------|-------------------------------------------------------------------------------------------------|
| `componentName`                | Unchanged. Still read-only and equal to the resource name.                                      |
| `type`                         | Removed. The Recipe decides the component type.                                                 |
| `version`                      | Removed. Set by a Recipe parameter.                                                             |
| `metadata`                     | Removed. Set by the Recipe, with a metadata override parameter for per-Environment tuning.      |
| `auth.secretStore`             | Removed. The Recipe wires credentials using the mechanism its platform supports.                |
| `resourceProvisioning: manual` | Removed. The equivalent is a Recipe that accepts the component type and metadata as parameters. |

Application code does not change, because the component name is still the resource name.

`Radius.Dapr/bindings` has no `Applications.Dapr` predecessor. Teams using Dapr bindings today declare them outside Radius, so adopting it is a move onto Radius rather than a migration.

## Appendix 2: How Dapr runs on each compute platform

| Platform             | How the sidecar runs                                                                                                 | How a component is created                                                                        | Who installs Dapr                                                                   |
|----------------------|----------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------|
| Kubernetes           | The sidecar injector adds `daprd` based on the `dapr.io/*` pod annotations set by the containers Recipe.             | `dapr.io/v1alpha1` `Component` in the application namespace.                                      | Platform engineer, with the Dapr Helm chart or `dapr init -k`.                      |
| Azure Container Apps | Built into the managed environment. The containers Recipe sets `properties.configuration.dapr` on the Container App. | `Microsoft.App/managedEnvironments/daprComponents`, which accept a `scopes` list of app IDs.      | Azure, as part of the managed environment.                                          |
| AWS ECS              | No injector exists. The containers Recipe adds a `daprd` sidecar container to the task definition.                   | `daprd` reads component YAML from a path, so the Recipe needs shared storage for component files. | Platform engineer, including the placement service and a name resolution component. |

ECS is out of MVP because it needs an ECS containers Recipe that injects a sidecar, plus shared storage and a name resolution strategy, since mDNS discovery does not work there. The schema in Feature 1 holds no platform-specific property, so ECS can be added with Recipes only.

## Appendix 3: Adjacent requests that are not Dapr work

Scale to zero for HTTP services is a frequent request from Dapr users, because serverless container platforms make it attractive. It belongs to the containers resource type and its autoscaling property, is driven by container start time rather than by Dapr, and is tracked with compute extensibility.
