import radius as radius

param environment string

resource app 'Applications.Core/applications@2022-03-15-privatepreview' = {
  name: 'daprrp-rs-component-name-conflict'
  location: 'global'
  properties: {
    environment: environment
  }
}

// Dapr Component #1
resource pubsub 'Applications.Dapr/pubSubBrokers@2022-03-15-privatepreview' = {
  name: 'dapr-component'
  location: 'global'
  properties: {
    environment: environment
    application: app.id
    resourceProvisioning: 'manual'
    type: 'pubsub.azure.servicebus'
    version: 'v1'
    metadata: {
      name: 'test'
    }
  }
}

// Dapr Component #2
resource secretstore 'Applications.Dapr/secretStores@2022-03-15-privatepreview' = {
  name: 'dapr-component'
  location: 'global'
  properties: {
    environment: environment
    application: app.id
    resourceProvisioning: 'manual'
    type: 'secretstores.kubernetes'
    metadata: {
      vaultName: 'test'
    }
    version: 'v1'
  }
}