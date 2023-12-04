package config

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type PublicKey struct {
	Realm           string `json:"realm"`
	PublicKey       string `json:"public_key"`
	TokenService    string `json:"token-service"`
	AccountService  string `json:"account-service"`
	TokensNotBefore int    `json:"tokens-not-before"`
}

func GetPublicKey(host, realm string) (string, error) {
	client := http.Client{}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, host+"/realms/"+realm, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New("unexpected status code")
	}

	var publicKey PublicKey
	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(raw, &publicKey)
	if err != nil {
		return "", err
	}
	return publicKey.PublicKey, nil
}
