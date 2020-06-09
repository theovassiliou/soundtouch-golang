// Code generated by go-swagger; DO NOT EDIT.

package key

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/theovassiliou/soundtouch-golang/server/models"
)

// PlayPauseNoContentCode is the HTTP code returned for type PlayPauseNoContent
const PlayPauseNoContentCode int = 204

/*PlayPauseNoContent empty on success

swagger:response playPauseNoContent
*/
type PlayPauseNoContent struct {

	/*empty on success
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewPlayPauseNoContent creates PlayPauseNoContent with default headers values
func NewPlayPauseNoContent() *PlayPauseNoContent {

	return &PlayPauseNoContent{}
}

// WithPayload adds the payload to the play pause no content response
func (o *PlayPauseNoContent) WithPayload(payload string) *PlayPauseNoContent {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the play pause no content response
func (o *PlayPauseNoContent) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PlayPauseNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(204)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*PlayPauseDefault generic error response

swagger:response playPauseDefault
*/
type PlayPauseDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPlayPauseDefault creates PlayPauseDefault with default headers values
func NewPlayPauseDefault(code int) *PlayPauseDefault {
	if code <= 0 {
		code = 500
	}

	return &PlayPauseDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the play pause default response
func (o *PlayPauseDefault) WithStatusCode(code int) *PlayPauseDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the play pause default response
func (o *PlayPauseDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the play pause default response
func (o *PlayPauseDefault) WithPayload(payload *models.Error) *PlayPauseDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the play pause default response
func (o *PlayPauseDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PlayPauseDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}