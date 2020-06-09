// Code generated by go-swagger; DO NOT EDIT.

package key

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/theovassiliou/soundtouch-golang/server/models"
)

// PlayNoContentCode is the HTTP code returned for type PlayNoContent
const PlayNoContentCode int = 204

/*PlayNoContent empty on success

swagger:response playNoContent
*/
type PlayNoContent struct {

	/*empty on success
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewPlayNoContent creates PlayNoContent with default headers values
func NewPlayNoContent() *PlayNoContent {

	return &PlayNoContent{}
}

// WithPayload adds the payload to the play no content response
func (o *PlayNoContent) WithPayload(payload string) *PlayNoContent {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the play no content response
func (o *PlayNoContent) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PlayNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(204)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*PlayDefault generic error response

swagger:response playDefault
*/
type PlayDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPlayDefault creates PlayDefault with default headers values
func NewPlayDefault(code int) *PlayDefault {
	if code <= 0 {
		code = 500
	}

	return &PlayDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the play default response
func (o *PlayDefault) WithStatusCode(code int) *PlayDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the play default response
func (o *PlayDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the play default response
func (o *PlayDefault) WithPayload(payload *models.Error) *PlayDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the play default response
func (o *PlayDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PlayDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
