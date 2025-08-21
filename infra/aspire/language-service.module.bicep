@description('The location for the resource(s) to be deployed.')
param location string = resourceGroup().location

resource language 'Microsoft.CognitiveServices/accounts@2024-10-01' = {
  name: take('language-${uniqueString(resourceGroup().id)}', 64)
  location: location
  kind: 'TextAnalytics'
  sku: {
    name: 'F0'
  }
}