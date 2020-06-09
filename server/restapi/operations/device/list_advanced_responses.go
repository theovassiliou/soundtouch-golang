// Code generated by go-swagger; DO NOT EDIT.

package device

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/theovassiliou/soundtouch-golang/server/models"
)

// ListAdvancedOKCode is the HTTP code returned for type ListAdvancedOK
const ListAdvancedOKCode int = 200

/*ListAdvancedOK a JSON object with the found soundtouch devices on your network.

swagger:response listAdvancedOK
*/
type ListAdvancedOK struct {

	/*
	  In: Body
	*/
	Payload models.BObject `json:"body,omitempty"`
}

// NewListAdvancedOK creates ListAdvancedOK with default headers values
func NewListAdvancedOK() *ListAdvancedOK {

	return &ListAdvancedOK{}
}

// WithPayload adds the payload to the list advanced o k response
func (o *ListAdvancedOK) WithPayload(payload models.BObject) *ListAdvancedOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list advanced o k response
func (o *ListAdvancedOK) SetPayload(payload models.BObject) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListAdvancedOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*ListAdvancedDefault generic error response

swagger:response listAdvancedDefault
*/
type ListAdvancedDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewListAdvancedDefault creates ListAdvancedDefault with default headers values
func NewListAdvancedDefault(code int) *ListAdvancedDefault {
	if code <= 0 {
		code = 500
	}

	return &ListAdvancedDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the list advanced default response
func (o *ListAdvancedDefault) WithStatusCode(code int) *ListAdvancedDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the list advanced default response
func (o *ListAdvancedDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the list advanced default response
func (o *ListAdvancedDefault) WithPayload(payload *models.Error) *ListAdvancedDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list advanced default response
func (o *ListAdvancedDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListAdvancedDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
