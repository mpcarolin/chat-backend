param accounts_chat_language_test_name string = 'chat-language-test'
param searchServices_chatlanguagetest_asztsuuubz6g44k_name string = 'chatlanguagetest-asztsuuubz6g44k'

resource searchServices_chatlanguagetest_asztsuuubz6g44k_name_resource 'Microsoft.Search/searchServices@2025-05-01' = {
  name: searchServices_chatlanguagetest_asztsuuubz6g44k_name
  location: 'East US'
  sku: {
    name: 'free'
  }
  properties: {
    replicaCount: 1
    partitionCount: 1
    endpoint: 'https://${searchServices_chatlanguagetest_asztsuuubz6g44k_name}.search.windows.net'
    hostingMode: 'Default'
    computeType: 'Default'
    publicNetworkAccess: 'Enabled'
    networkRuleSet: {
      ipRules: []
      bypass: 'None'
    }
    encryptionWithCmk: {}
    disableLocalAuth: false
    authOptions: {
      apiKeyOnly: {}
    }
    dataExfiltrationProtections: []
    semanticSearch: 'disabled'
    upgradeAvailable: 'notAvailable'
  }
}

resource accounts_chat_language_test_name_resource 'Microsoft.CognitiveServices/accounts@2025-06-01' = {
  name: accounts_chat_language_test_name
  location: 'eastus'
  sku: {
    name: 'F0'
  }
  kind: 'TextAnalytics'
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    apiProperties: {
      qnaAzureSearchEndpointId: searchServices_chatlanguagetest_asztsuuubz6g44k_name_resource.id
    }
    customSubDomainName: accounts_chat_language_test_name
    networkAcls: {
      defaultAction: 'Allow'
      virtualNetworkRules: []
      ipRules: []
    }
    allowProjectManagement: false
    publicNetworkAccess: 'Enabled'
  }
}
