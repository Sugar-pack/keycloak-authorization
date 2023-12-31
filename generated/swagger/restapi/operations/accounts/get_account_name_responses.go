// Code generated by go-swagger; DO NOT EDIT.

package accounts

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"test_iam/generated/swagger/models"
)

// GetAccountNameOKCode is the HTTP code returned for type GetAccountNameOK
const GetAccountNameOKCode int = 200

/*
GetAccountNameOK Success

swagger:response getAccountNameOK
*/
type GetAccountNameOK struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetAccountNameOK creates GetAccountNameOK with default headers values
func NewGetAccountNameOK() *GetAccountNameOK {

	return &GetAccountNameOK{}
}

// WithPayload adds the payload to the get account name o k response
func (o *GetAccountNameOK) WithPayload(payload string) *GetAccountNameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get account name o k response
func (o *GetAccountNameOK) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAccountNameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetAccountNameDefault Unexpected error.

swagger:response getAccountNameDefault
*/
type GetAccountNameDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAccountNameDefault creates GetAccountNameDefault with default headers values
func NewGetAccountNameDefault(code int) *GetAccountNameDefault {
	if code <= 0 {
		code = 500
	}

	return &GetAccountNameDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get account name default response
func (o *GetAccountNameDefault) WithStatusCode(code int) *GetAccountNameDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get account name default response
func (o *GetAccountNameDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get account name default response
func (o *GetAccountNameDefault) WithPayload(payload *models.Error) *GetAccountNameDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get account name default response
func (o *GetAccountNameDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAccountNameDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
