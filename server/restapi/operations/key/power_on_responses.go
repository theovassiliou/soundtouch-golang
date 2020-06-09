// Code generated by go-swagger; DO NOT EDIT.

package key

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/theovassiliou/soundtouch-golang/server/models"
)

// PowerOnOKCode is the HTTP code returned for type PowerOnOK
const PowerOnOKCode int = 200

/*PowerOnOK The returned status will be true if the SoundTouch is turned on. The returned status will be false if the SoundTouch was already turned on.

swagger:response powerOnOK
*/
type PowerOnOK struct {

	/*
	  In: Body
	*/
	Payload *models.BStatus `json:"body,omitempty"`
}

// NewPowerOnOK creates PowerOnOK with default headers values
func NewPowerOnOK() *PowerOnOK {

	return &PowerOnOK{}
}

// WithPayload adds the payload to the power on o k response
func (o *PowerOnOK) WithPayload(payload *models.BStatus) *PowerOnOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the power on o k response
func (o *PowerOnOK) SetPayload(payload *models.BStatus) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PowerOnOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PowerOnDefault generic error response

swagger:response powerOnDefault
*/
type PowerOnDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPowerOnDefault creates PowerOnDefault with default headers values
func NewPowerOnDefault(code int) *PowerOnDefault {
	if code <= 0 {
		code = 500
	}

	return &PowerOnDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the power on default response
func (o *PowerOnDefault) WithStatusCode(code int) *PowerOnDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the power on default response
func (o *PowerOnDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the power on default response
func (o *PowerOnDefault) WithPayload(payload *models.Error) *PowerOnDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the power on default response
func (o *PowerOnDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PowerOnDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
