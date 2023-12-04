package main

import (
	"log"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-openapi/loads"

	"test_iam/generated/swagger/restapi"
	"test_iam/generated/swagger/restapi/operations"
	"test_iam/middleware"
	evaluationengine "test_iam/middleware/evaluation-engine"

	"test_iam/config"
	"test_iam/handlers"
)

func main() {
	// config
	conf, err := config.GetAppConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	// swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalf("Error loading swagger spec: %s", err)
	}
	api := operations.NewTestKeycloakAbacAPI(swaggerSpec)
	api.UseSwaggerUI()

	// handlers

	api.SystemGetHealthHandler = handlers.Health()
	api.AccountsGetAccountNameHandler = handlers.AccountNameByID()

	// middleware
	auth := middleware.NewAuthenticator(conf.Keycloak.PublicKey)
	api.AuthAuth = auth.Authenticate

	keycloakClient := gocloak.NewClient(conf.Keycloak.Url)
	ee, err := evaluationengine.NewEvaluationEngine(keycloakClient, conf.Keycloak.Url, conf.Keycloak.Realm,
		conf.Keycloak.ClientID, conf.Keycloak.ClientSecret)
	if err != nil {
		log.Fatalf("Error creating evaluation engine: %s", err)
	}
	authorizer := middleware.NewAuthorizer(ee)
	api.APIAuthorizer = authorizer

	api.ServeError = middleware.ServeError()

	// run server
	api.Init()
	server := restapi.NewServer(api)

	server.EnabledListeners = []string{"http"}
	server.Host = conf.Server.Host
	server.Port = conf.Server.Port

	server.SetHandler(
		api.Serve(nil),
	)
	// Swagger servers handles signals and gracefully shuts down by itself
	if err = server.Serve(); err != nil {
		log.Fatalf("Error serving: %s", err)
	}

	if errShutdown := server.Shutdown(); errShutdown != nil {
		log.Fatalf("Error shutting down: %s", errShutdown)
	}
}
