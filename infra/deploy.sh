# Guide: https://learn.microsoft.com/en-us/azure/container-apps/tutorial-code-to-cloud?tabs=bash%2Cgo&pivots=docker-local

######### Azure CLI authentication and setup
az login
az upgrade
az extension add --name containerapp --upgrade

az provider register --namespace Microsoft.App
az provider register --namespace Microsoft.OperationalInsights

RESOURCE_GROUP="album-containerapps"
LOCATION="canadacentral"
ENVIRONMENT="env-album-containerapps"
API_NAME="album-api"
FRONTEND_NAME="album-ui"
GITHUB_USERNAME="mpcarolin"

# registry name
ACR_NAME="acaalbums"$GITHUB_USERNAME

######### Resource Group  
az group create --name $RESOURCE_GROUP --location "$LOCATION"

######### Container Registry   
az acr create --resource-group $RESOURCE_GROUP --location $LOCATION --name $ACR_NAME --sku Basic

# check if ARM tokens are allowed
az acr config authentication-as-arm show --registry "$ACR_NAME"
# enable it if it's not enabled
az acr config authentication-as-arm update --registry "$ACR_NAME" --status enabled


######### Azure Identity 
IDENTITY="<YOUR_IDENTITY_NAME>"
az identity create --name $IDENTITY --resource-group $RESOURCE_GROUP
IDENTITY_ID=$(az identity show --name $IDENTITY --resource-group $RESOURCE_GROUP --query id --output tsv)

######### Docker Image Build and Push 
# build the api and tag it with our registry
docker compose build api --tag $ACR_NAME.azurecr.io/$API_NAME .

# login to the azure container registry
az acr login --name $ACR_NAME

# push docker image to registry
docker push $ACR_NAME.azurecr.io/$API_NAME


######### Container Apps Resources
az containerapp env create --name $ENVIRONMENT --resource-group $RESOURCE_GROUP --location "$LOCATION"

# deploy the image to container app
az containerapp create --name $API_NAME \
    --resource-group $RESOURCE_GROUP \
    --environment $ENVIRONMENT \
    --image $ACR_NAME.azurecr.io/$API_NAME \
    --target-port 8080 \
    --ingress external \
    --registry-server $ACR_NAME.azurecr.io \
    --user-assigned "$IDENTITY_ID" \
    --registry-identity "$IDENTITY_ID" \
    --query properties.configuration.ingress.fqdn

######### Cleanup 
# Manually run this when you want to clean up the whole resource group
# az group delete --name $RESOURCE_GROUP
