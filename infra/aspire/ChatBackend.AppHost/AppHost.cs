var builder = DistributedApplication.CreateBuilder(args);

// Azure Language Service (includes Custom Question Answering) - free tier
var languageService = builder.AddAzureCognitiveServices("language-service")
    .ConfigureConstruct(construct =>
    {
        construct.Kind = "TextAnalytics"; // Language Service kind  
        construct.Sku = new() { Name = "F0" }; // Free tier
    });

// Ollama container (local only)
var ollama = builder.AddDockerfile("ollama", "../../../", "Dockerfile.ollama")
    .WithHttpEndpoint(targetPort: 11434, port: 11434, name: "http")
    .WithEnvironment("OLLAMA_KEEP_ALIVE", "24h")
    .WithEnvironment("OLLAMA_HOST", "0.0.0.0")
    .ExcludeFromManifest();

// Go API service
var api = builder.AddDockerfile("api", "../../../packages/api")
    .WithHttpEndpoint(targetPort: 8090, port: 8090, name: "http")
    .WithEnvironment("CHAT_PROVIDER", builder.Configuration["CHAT_PROVIDER"] ?? "mock")
    .WithEnvironment("OLLAMA_MODEL", "gemma3:1b")
    .WithEnvironment("OLLAMA_BASE_URL", ollama.GetEndpoint("http"))
    .WithEnvironment("AZURE_QNA_ENDPOINT", builder.Configuration["AZURE_QNA_ENDPOINT"] ?? "")
    .WithEnvironment("AZURE_QNA_API_KEY", builder.Configuration["AZURE_QNA_API_KEY"] ?? "")
    .WithEnvironment("AZURE_QNA_PROJECT_NAME", builder.Configuration["AZURE_QNA_PROJECT_NAME"] ?? "")
    .WithEnvironment("AZURE_QNA_DEPLOYMENT_NAME", builder.Configuration["AZURE_QNA_DEPLOYMENT_NAME"] ?? "")
    .PublishAsAzureContainerApp();

// React web app
var web = builder.AddDockerfile("web", "../../../packages/web")
    .WithHttpEndpoint(targetPort: 5173, port: 5173, name: "http")
    .WithEnvironment("VITE_API_URL", api.GetEndpoint("http"))
    .ExcludeFromManifest(); // in cloud, react app prod bundle is served by api 

builder.Build().Run();
