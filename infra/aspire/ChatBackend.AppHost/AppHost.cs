var builder = DistributedApplication.CreateBuilder(args);

// Configure Container Apps Environment with VNet
builder.AddAzureProvisioning().ConfigureInfrastructure(context =>
{
    // Create VNet and subnet for Container Apps
    var vnet = context.Infrastructure.Add(context.Provisioning.CreateVirtualNetwork("chat-vnet")
        .ConfigureConstruct(vnet =>
        {
            vnet.AddressSpace = new() { AddressPrefixes = { "10.0.0.0/16" } };
        }));
    
    var subnet = context.Infrastructure.Add(context.Provisioning.CreateSubnet("container-apps-subnet")
        .ConfigureConstruct(subnet =>
        {
            subnet.Parent = vnet;
            subnet.AddressPrefix = "10.0.1.0/24";
        }));

    // Configure Container Apps Environment to use the VNet
    var containerAppsEnvironment = context.GetProvisionableResources().OfType<ContainerAppsEnvironment>().Single();
    containerAppsEnvironment.ConfigureConstruct(env =>
    {
        env.VnetConfiguration = new()
        {
            Internal = true, // Makes it internal (not publicly accessible)
            InfrastructureSubnetId = subnet.Id
        };
    });
});

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

// Go API service (publicly accessible)
var api = builder.AddDockerfile("api", "../../../packages/api")
    .WithHttpEndpoint(targetPort: 8090, port: 8090, name: "http")
    .WithEnvironment("CHAT_PROVIDER", builder.Configuration["CHAT_PROVIDER"] ?? "mock")
    .WithEnvironment("OLLAMA_MODEL", "gemma3:1b")
    .WithEnvironment("OLLAMA_BASE_URL", ollama.GetEndpoint("http"))
    .WithEnvironment("AZURE_QNA_ENDPOINT", builder.Configuration["AZURE_QNA_ENDPOINT"] ?? "")
    .WithEnvironment("AZURE_QNA_API_KEY", builder.Configuration["AZURE_QNA_API_KEY"] ?? "")
    .WithEnvironment("AZURE_QNA_PROJECT_NAME", builder.Configuration["AZURE_QNA_PROJECT_NAME"] ?? "")
    .WithEnvironment("AZURE_QNA_DEPLOYMENT_NAME", builder.Configuration["AZURE_QNA_DEPLOYMENT_NAME"] ?? "")
    .PublishAsAzureContainerApp()
    .ConfigureConstruct(containerApp =>
    {
        // Make API publicly accessible from VNET
        containerApp.Configuration!.Ingress!.External = true;
    });

// React web app
var web = builder.AddDockerfile("web", "../../../packages/web")
    .WithHttpEndpoint(targetPort: 5173, port: 5173, name: "http")
    .WithEnvironment("VITE_API_URL", api.GetEndpoint("http"))
    .ExcludeFromManifest(); // in cloud, react app prod bundle is served by api 

builder.Build().Run();
