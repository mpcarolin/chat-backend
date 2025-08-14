# Chat Backend

A Go-based FAQ chatbot backend with support for multiple chat providers including Azure Q&A, Ollama, and mock responses.

## Architecture

```mermaid
graph TB
    %% External clients
    Client[HTTP Client]
    
    %% Main application layer
    Client --> Echo[Echo Web Server<br/>:8090]
    
    %% Middleware layer
    Echo --> MW[Middleware Layer]
    MW --> Logger[Logger Middleware]
    MW --> RateLimit[Rate Limit Middleware]
    
    %% Handler layer
    MW --> Routes{Routes}
    Routes -->|GET /status| StatusHandler[Status Handler]
    Routes -->|POST /api/faq| FAQHandler[FAQ Handler]
    
    %% Application context
    StatusHandler --> AppContext[App Context]
    FAQHandler --> AppContext
    
    %% Chat provider interface
    AppContext --> ChatProvider[Chat Provider Interface]
    
    %% Provider implementations
    ChatProvider --> MockProvider[Mock Provider]
    ChatProvider --> AzureProvider[Azure Q&A Provider]
    ChatProvider --> OllamaProvider[Ollama Provider]
    
    %% Mock provider components
    MockProvider --> TSVData[TSV Data File<br/>sample-data.tsv]
    
    %% Azure provider components
    AzureProvider --> AzureClient[Azure HTTP Client]
    AzureClient --> AzureAPI[Azure Cognitive Services<br/>Q&A API]
    
    %% Ollama provider components
    OllamaProvider --> OllamaClient[Ollama HTTP Client]
    OllamaClient --> OllamaAPI[Ollama Local API<br/>:11434]
    
    %% Environment configuration
    EnvVars[Environment Variables] --> AppContext
    
    %% Styling
    classDef external fill:#e1f5fe
    classDef server fill:#f3e5f5
    classDef middleware fill:#fff3e0
    classDef handler fill:#e8f5e8
    classDef provider fill:#fff8e1
    classDef client fill:#fce4ec
    classDef api fill:#f1f8e9
    
    class Client,TSVData,AzureAPI,OllamaAPI external
    class Echo server
    class MW,Logger,RateLimit middleware
    class Routes,StatusHandler,FAQHandler handler
    class AppContext,ChatProvider,MockProvider,AzureProvider,OllamaProvider provider
    class AzureClient,OllamaClient client
    class EnvVars api
```

## Features

- **Multiple Chat Providers**: Support for Azure Q&A, Ollama, and mock responses
- **Streaming Support**: Ollama provider supports streaming responses
- **Rate Limiting**: Built-in rate limiting middleware
- **Structured Logging**: Using Go's structured logging
- **Environment-based Configuration**: Easy provider switching via environment variables

## Configuration

The application uses the `CHAT_PROVIDER` environment variable to determine which provider to use:

### Mock Provider (Default)
```bash
CHAT_PROVIDER=mock
```
- Uses local TSV data file for responses
- Fallback responses when TSV file is unavailable

### Azure Q&A Provider
```bash
CHAT_PROVIDER=azure-qa
AZURE_QNA_ENDPOINT=https://your-service.cognitiveservices.azure.com
AZURE_QNA_API_KEY=your-api-key
AZURE_QNA_PROJECT_NAME=your-project
AZURE_QNA_DEPLOYMENT_NAME=your-deployment
```

### Ollama Provider
```bash
CHAT_PROVIDER=ollama
OLLAMA_BASE_URL=http://localhost:11434  # Optional, defaults to localhost:11434
OLLAMA_MODEL=mistral                    # Optional, defaults to mistral
```

