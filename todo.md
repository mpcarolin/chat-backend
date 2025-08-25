# TODO: Setup .NET Aspire 

## Description

This project has two programs: a golang web api (in packages/api) and a vite react app (packages/web). I need to be able to run both of these in a .NET aspire configuration. They should be runnable locally but also support deploying to Azure.

## Acceptance Criteria
- All Aspire configuration files should be placed in `infra/aspire/`
- I should be able to run the web api and the react app simultaneously with this setup
- The same ports and configuration defined in the docker-compose.yml file should be emulated
- I should be able to run the apps locally and see them at the .net aspire dashboard
- I should be able to deploy to azure, preferably with azure container apps (or azure containers for web apps)
- Preferably should support also running the ollama container also defined in the docker-compose.yml. If this is not possible in the .net aspire config directly, tell me.

## Notes
General goals
- Run the entire app locally with one command
- Deploy the entire app with one command





