package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"test_iam/generated/swagger/models"
	"test_iam/generated/swagger/restapi/operations/accounts"
)

func AccountNameByID() accounts.GetAccountNameHandlerFunc {
	return func(params accounts.GetAccountNameParams, a *models.Principal) middleware.Responder {
		return accounts.NewGetAccountNameOK().WithPayload("test name")
	}
}
