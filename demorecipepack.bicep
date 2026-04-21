extension radius

resource demorecipepack 'Radius.Core/recipePacks@2025-08-01-preview' = {
  name: 'demorecipepack'
  location: 'global'
  properties: {
    recipes: {
      'Radius.Compute/containers': {
        recipeKind: 'bicep'
        recipeLocation: 'ghcr.io/reshrahim/recipes/containers:1.0'
      }
      'Radius.Data/mySqlDatabases': {
        recipeKind: 'terraform'
        recipeLocation: ''
        parameters:{
         vpcId: '' 
         subnetIds: ''
        }
      }
      'Radius.Security/secrets': {
        recipeKind: 'bicep'
        recipeLocation: 'ghcr.io/reshrahim/recipes/secrets:1.0'
      }
    }
  }
}
