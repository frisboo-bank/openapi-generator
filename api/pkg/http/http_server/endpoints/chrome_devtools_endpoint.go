package endpoints

import (
	"net/http"
	"os"

	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/routing"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

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
	logger loggerContracts.Logger,
	root contracts.RouteGroup,
) routing.Endpoint {
	validation.AssertValidEnum("env", env)
	validation.AssertNotNil("logger", logger)
	validation.AssertNotNil("root", root)

	return &chromeDevtoolEndpoint{
		Env:    env,
		Logger: logger,
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

		config := &WorkspaceConfig{
			Workspace: struct {
				Root string `json:"root"`
				UUID string `json:"uuid"`
			}{
				Root: projectRoot,
				UUID: workspaceUUID,
			},
		}

		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

		routing.JSON(w, http.StatusOK, config)
	}
}
