var builder = DistributedApplication.CreateBuilder(args);

// Go API service (publicly accessible), used both locally and remotely
var api = builder.AddDockerfile("api", "../../../packages/api")
    .WithHttpEndpoint(targetPort: 8090, port: 8090, name: "http")
    .WithEnvironment("CHAT_PROVIDER", builder.Configuration["CHAT_PROVIDER"] ?? "mock")
    .WithEnvironment("OLLAMA_MODEL", "gemma3:1b")
    .WithEnvironment("AZURE_QNA_ENDPOINT", builder.Configuration["AZURE_QNA_ENDPOINT"] ?? "")
    .WithEnvironment("AZURE_QNA_API_KEY", builder.Configuration["AZURE_QNA_API_KEY"] ?? "")
    .WithEnvironment("AZURE_QNA_PROJECT_NAME", builder.Configuration["AZURE_QNA_PROJECT_NAME"] ?? "")
    .WithEnvironment("AZURE_QNA_DEPLOYMENT_NAME", builder.Configuration["AZURE_QNA_DEPLOYMENT_NAME"] ?? "");

if (builder.ExecutionContext.IsRunMode)
{
    // React web app *dev server*. In cloud, react app prof bundle is served by api directly,
    // so we only need this locally
    var web = builder.AddDockerfile("web", "../../../packages/web")
        .WithHttpEndpoint(targetPort: 5173, port: 5173, name: "http")
        .WithEnvironment("VITE_API_URL", api.GetEndpoint("http"));

    // ollama is also only for local testing 
    var ollama = builder.AddDockerfile("ollama", "../../../", "Dockerfile.ollama")
        .WithHttpEndpoint(targetPort: 11434, port: 11434, name: "http")
        .WithEnvironment("OLLAMA_KEEP_ALIVE", "24h")
        .WithEnvironment("OLLAMA_HOST", "0.0.0.0");

    // configuring api for ollama
    api.WithEnvironment("OLLAMA_BASE_URL", ollama.GetEndpoint("http"));
}
else if (builder.ExecutionContext.IsPublishMode)
{
    // Azure Language Service (includes Custom Question Answering) - only provision when deploying
    var languageService = builder.AddAzureInfrastructure("language-service", infra =>
    {
        var cognitiveService = new Azure.Provisioning.CognitiveServices.CognitiveServicesAccount("language")
        {
            Sku = new Azure.Provisioning.CognitiveServices.CognitiveServicesSku()
            {
                Name = "F0" // Free tier
            },
            Kind = "TextAnalytics" // Language Service kind
        };
        infra.Add(cognitiveService);
    });
}


builder.Build().Run();
