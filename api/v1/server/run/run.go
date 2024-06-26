package run

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"

	"github.com/hatchet-dev/hatchet/api/v1/server/authn"
	"github.com/hatchet-dev/hatchet/api/v1/server/authz"
	apitokens "github.com/hatchet-dev/hatchet/api/v1/server/handlers/api-tokens"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/events"
	githubapp "github.com/hatchet-dev/hatchet/api/v1/server/handlers/github-app"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/ingestors"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/logs"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/metadata"
	slackapp "github.com/hatchet-dev/hatchet/api/v1/server/handlers/slack-app"
	stepruns "github.com/hatchet-dev/hatchet/api/v1/server/handlers/step-runs"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/tenants"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/users"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/workers"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/workflows"
	hatchetmiddleware "github.com/hatchet-dev/hatchet/api/v1/server/middleware"
	"github.com/hatchet-dev/hatchet/api/v1/server/middleware/populator"
	"github.com/hatchet-dev/hatchet/api/v1/server/oas/gen"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

type apiService struct {
	*users.UserService
	*tenants.TenantService
	*events.EventService
	*logs.LogService
	*workflows.WorkflowService
	*workers.WorkerService
	*metadata.MetadataService
	*apitokens.APITokenService
	*stepruns.StepRunService
	*githubapp.GithubAppService
	*ingestors.IngestorsService
	*slackapp.SlackAppService
}

func newAPIService(config *server.ServerConfig) *apiService {
	return &apiService{
		UserService:      users.NewUserService(config),
		TenantService:    tenants.NewTenantService(config),
		EventService:     events.NewEventService(config),
		LogService:       logs.NewLogService(config),
		WorkflowService:  workflows.NewWorkflowService(config),
		WorkerService:    workers.NewWorkerService(config),
		MetadataService:  metadata.NewMetadataService(config),
		APITokenService:  apitokens.NewAPITokenService(config),
		StepRunService:   stepruns.NewStepRunService(config),
		GithubAppService: githubapp.NewGithubAppService(config),
		IngestorsService: ingestors.NewIngestorsService(config),
		SlackAppService:  slackapp.NewSlackAppService(config),
	}
}

type APIServer struct {
	config *server.ServerConfig
}

func NewAPIServer(config *server.ServerConfig) *APIServer {
	return &APIServer{
		config: config,
	}
}

func (t *APIServer) Run() (func() error, error) {
	oaspec, err := gen.GetSwagger()
	if err != nil {
		return nil, err
	}

	e := echo.New()

	// application middleware
	populatorMW := populator.NewPopulator(t.config)

	populatorMW.RegisterGetter("tenant", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		tenant, err := config.APIRepository.Tenant().GetTenantByID(id)

		if err != nil {
			return nil, "", err
		}

		return tenant, "", nil
	})

	populatorMW.RegisterGetter("api-token", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		apiToken, err := config.APIRepository.APIToken().GetAPITokenById(id)

		if err != nil {
			return nil, "", err
		}

		// at the moment, API tokens should have a tenant id, because there are no other types of
		// API tokens. If we add other types of API tokens, we'll need to pass in a parent id to query
		// for.
		tenantId, ok := apiToken.TenantID()

		if !ok {
			return nil, "", fmt.Errorf("api token has no tenant id")
		}

		return apiToken, tenantId, nil
	})

	populatorMW.RegisterGetter("tenant-invite", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		tenantInvite, err := config.APIRepository.TenantInvite().GetTenantInvite(id)

		if err != nil {
			return nil, "", err
		}

		return tenantInvite, tenantInvite.TenantID, nil
	})

	populatorMW.RegisterGetter("slack", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		slackWebhook, err := config.APIRepository.Slack().GetSlackWebhookById(id)

		if err != nil {
			return nil, "", err
		}

		return slackWebhook, slackWebhook.TenantID, nil
	})

	populatorMW.RegisterGetter("alert-email-group", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		emailGroup, err := config.APIRepository.TenantAlertingSettings().GetTenantAlertGroupById(id)

		if err != nil {
			return nil, "", err
		}

		return emailGroup, emailGroup.TenantID, nil
	})

	populatorMW.RegisterGetter("sns", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		snsIntegration, err := config.APIRepository.SNS().GetSNSIntegrationById(id)

		if err != nil {
			return nil, "", err
		}

		return snsIntegration, snsIntegration.TenantID, nil
	})

	populatorMW.RegisterGetter("workflow", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		workflow, err := config.APIRepository.Workflow().GetWorkflowById(id)

		if err != nil {
			return nil, "", err
		}

		return workflow, workflow.TenantID, nil
	})

	populatorMW.RegisterGetter("workflow-run", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		workflowRun, err := config.APIRepository.WorkflowRun().GetWorkflowRunById(parentId, id)

		if err != nil {
			return nil, "", err
		}

		return workflowRun, workflowRun.TenantID, nil
	})

	populatorMW.RegisterGetter("step-run", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		stepRun, err := config.APIRepository.StepRun().GetStepRunById(parentId, id)

		if err != nil {
			return nil, "", err
		}

		return stepRun, stepRun.TenantID, nil
	})

	populatorMW.RegisterGetter("event", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		event, err := config.APIRepository.Event().GetEventById(id)

		if err != nil {
			return nil, "", err
		}

		return event, event.TenantID, nil
	})

	populatorMW.RegisterGetter("worker", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		worker, err := config.APIRepository.Worker().GetWorkerById(id)

		if err != nil {
			return nil, "", err
		}

		return worker, worker.TenantID, nil
	})

	populatorMW.RegisterGetter("gh-installation", func(config *server.ServerConfig, parentId, id string) (result interface{}, uniqueParentId string, err error) {
		ghInstallation, err := config.APIRepository.Github().ReadGithubAppInstallationByID(id)

		if err != nil {
			return nil, "", err
		}

		return ghInstallation, "", nil
	})

	authnMW := authn.NewAuthN(t.config)
	authzMW := authz.NewAuthZ(t.config)

	mw, err := hatchetmiddleware.NewMiddlewareHandler(oaspec)

	if err != nil {
		return nil, err
	}

	mw.Use(populatorMW.Middleware)
	mw.Use(authnMW.Middleware)
	mw.Use(authzMW.Middleware)

	allHatchetMiddleware, err := mw.Middleware()

	if err != nil {
		return nil, err
	}

	loggerMiddleware := middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogError:     true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURIPath:   true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			statusCode := v.Status

			// note that the status code is not set yet as it gets picked up by the global err handler
			// see here: https://github.com/labstack/echo/issues/2310#issuecomment-1288196898
			if v.Error != nil {
				statusCode = 500
			}

			var e *zerolog.Event

			switch {
			case statusCode >= 500:
				e = t.config.Logger.Error().Err(v.Error)
			case statusCode >= 400:
				e = t.config.Logger.Warn()
			default:
				e = t.config.Logger.Info()
			}

			e.
				Dur("latency", v.Latency).
				Int("status", statusCode).
				Str("method", v.Method).
				Str("uri", v.URI).
				Str("user_agent", v.UserAgent).
				Str("remote_ip", v.RemoteIP).
				Str("host", v.Host).
				Msg("API")

			return nil
		},
	})

	// register echo middleware
	e.Use(
		loggerMiddleware,
		middleware.Recover(),
		allHatchetMiddleware,
	)

	service := newAPIService(t.config)

	myStrictApiHandler := gen.NewStrictHandler(service, []gen.StrictMiddlewareFunc{})

	gen.RegisterHandlers(e, myStrictApiHandler)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", t.config.Runtime.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	cleanup := func() error {
		return e.Shutdown(context.Background())
	}

	return cleanup, nil
}
