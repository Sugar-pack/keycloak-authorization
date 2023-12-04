package middleware

import (
	"context"
	"net/http"
	"strings"

	"test_iam/apierrors"
	"test_iam/generated/swagger/models"
)

const (
	forbiddenMessage = "Forbidden"
)

type EvaluationEngine interface {
	GetResourceNameByURI(ctx context.Context, uri string) (string, error)
	EvaluateResources(ctx context.Context, userAccessToken string, resourceName string, resourceAction string) error
}

type Authorizer struct {
	engine EvaluationEngine
}

func NewAuthorizer(evaluationEngine EvaluationEngine) *Authorizer {
	return &Authorizer{
		engine: evaluationEngine,
	}
}

func (a *Authorizer) Authorize(r *http.Request, auth interface{}) error {
	ctx := r.Context()
	// check token info
	principal, ok := auth.(*models.Principal)
	if !ok {
		return apierrors.NewAPIError(http.StatusForbidden, forbiddenMessage)
	}
	if principal == nil {
		return apierrors.NewAPIError(http.StatusForbidden, forbiddenMessage)
	}
	// get request url and process it
	uri := formatURI(r.URL.Path)
	// get resource
	requestResource, err := a.engine.GetResourceNameByURI(ctx, uri)
	if err != nil {
		return apierrors.NewAPIError(http.StatusForbidden, forbiddenMessage)
	}
	err = a.engine.EvaluateResources(ctx, string(*principal), requestResource, r.Method)
	if err != nil {
		return apierrors.NewAPIError(http.StatusForbidden, forbiddenMessage)
	}
	return nil
}

func formatURI(uri string) string {
	if strings.HasSuffix(uri, "/") {
		return uri[:len(uri)-1]
	}
	return uri
}
