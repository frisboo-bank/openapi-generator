package endpoints

import (
	"net/http"
	"os"

	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/routing"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"github.com/google/uuid"
)

type WorkspaceConfig struct {
	Workspace struct {
		Root string `json:"root"`
		UUID string `json:"uuid"`
	} `json:"workspace"`
}

var _ routing.Endpoint = (*chromeDevtoolEndpoint)(nil)

type chromeDevtoolEndpoint struct {
	Env    environmentEnum.Environment
	Logger loggerContracts.Logger
	Root   contracts.RouteGroup
}

func NewChromeDevtoolEndpoint(
	env environmentEnum.Environment,
	log loggerContracts.Logger,
	root contracts.RouteGroup,
) routing.Endpoint {
	return &chromeDevtoolEndpoint{
		Env:    env,
		Logger: log,
		Root:   root,
	}
}

func (ep *chromeDevtoolEndpoint) MapEndpoint() {
	ep.Root.GET("/.well-known/appspecific/com.chrome.devtools.json", ep.handler())
}

func (ep *chromeDevtoolEndpoint) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ep.Env.IsDevelopment() {
			routing.JSON(w, http.StatusOK, "{}")
			return
		}

		projectRoot, err := os.Getwd()
		if err != nil {
			routing.Error(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		workspaceUUID := uuid.New().String()

		var config WorkspaceConfig
		config.Workspace.Root = projectRoot
		config.Workspace.UUID = workspaceUUID

		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

		routing.JSON(w, http.StatusOK, config)
	}
}
