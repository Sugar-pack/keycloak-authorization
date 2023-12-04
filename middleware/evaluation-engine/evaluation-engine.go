package evaluation_engine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
)

type EvaluationEngine struct {
	idOfClient   string
	clientID     string
	clientSecret string
	realm        string
	host         string
	tokenURL     string
	evalClient   EvaluationClient
	httpClient   *http.Client
}

const tokenURITemplate = "%s/realms/%s/protocol/openid-connect/token" //nolint:gosec

type EvaluationClient interface {
	LoginClient(ctx context.Context, clientID string, clientSecret string, realm string) (*gocloak.JWT, error)
	GetResources(ctx context.Context, token string, realm string, idOfClient string, params gocloak.GetResourceParams) ([]*gocloak.ResourceRepresentation, error)
	GetClients(ctx context.Context, token string, realm string, params gocloak.GetClientsParams) ([]*gocloak.Client, error)
}

func NewEvaluationEngine(client EvaluationClient, host, realm, clientID, clientSecret string) (*EvaluationEngine, error) {
	token, err := client.LoginClient(context.Background(), clientID, clientSecret, realm)
	if err != nil {
		return nil, err
	}
	keycloakClient, err := client.GetClients(context.Background(), token.AccessToken, realm, gocloak.GetClientsParams{
		ClientID: gocloak.StringP(clientID),
		Max:      gocloak.IntP(1),
	})
	if err != nil {
		return nil, err
	}
	if keycloakClient != nil && keycloakClient[0] != nil && keycloakClient[0].ID != nil {
		return &EvaluationEngine{
			host:         host,
			clientID:     clientID,
			realm:        realm,
			idOfClient:   *keycloakClient[0].ID,
			evalClient:   client,
			httpClient:   http.DefaultClient,
			clientSecret: clientSecret,
			tokenURL:     fmt.Sprintf(tokenURITemplate, host, realm),
		}, nil
	}

	return nil, fmt.Errorf("keycloak client %s does not have id", clientID)
}

func (a *EvaluationEngine) EvaluateResources(ctx context.Context, userAccessToken string, resourceName, resourceScope string) error {
	return a.evaluate(ctx, userAccessToken, a.tokenURL, resourceName, resourceScope)
}

func (a *EvaluationEngine) evaluate(ctx context.Context, accessToken, tokenURL, resourceName, scope string) error {
	data := []byte(fmt.Sprintf("grant_type=urn:ietf:params:oauth:grant-type:uma-ticket&subject_token=%s&permission=%s#%s&response_mode=decision&audience=%s", accessToken, resourceName, scope, a.clientID))
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(a.clientID, a.clientSecret)

	resp, err := a.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("error evaluating resource %s with scope %s. return code %s", resourceName, scope, resp.Status)
}

func (a *EvaluationEngine) GetResourceNameByURI(ctx context.Context, uri string) (string, error) {
	token, err := a.evalClient.LoginClient(ctx, a.clientID, a.clientSecret, a.realm)
	if err != nil {
		return "", err
	}
	resources, err := a.evalClient.GetResources(ctx, token.AccessToken, a.realm, a.idOfClient, gocloak.GetResourceParams{
		URI:         gocloak.StringP(uri),
		MatchingURI: gocloak.BoolP(true),
	}) // keycloak return only one (ore zero) resource if search by uri
	if err != nil {
		return "", err
	}
	if len(resources) != 1 {
		return "", errors.New("resource not found or more than one resource found")
	}
	if resources[0] == nil || resources[0].URIs == nil || resources[0].ID == nil || resources[0].Name == nil {
		return "", errors.New("resource not valid")
	}
	return *resources[0].Name, nil
}
