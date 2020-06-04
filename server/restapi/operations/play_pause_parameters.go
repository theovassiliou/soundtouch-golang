// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// NewPlayPauseParams creates a new PlayPauseParams object
// no default values defined in spec.
func NewPlayPauseParams() PlayPauseParams {

	return PlayPauseParams{}
}

// PlayPauseParams contains all the bound params for the play pause operation
// typically these are obtained from a http.Request
//
// swagger:parameters playPause
type PlayPauseParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	SpeakerName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPlayPauseParams() beforehand.
func (o *PlayPauseParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rSpeakerName, rhkSpeakerName, _ := route.Params.GetOK("speakerName")
	if err := o.bindSpeakerName(rSpeakerName, rhkSpeakerName, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindSpeakerName binds and validates parameter SpeakerName from path.
func (o *PlayPauseParams) bindSpeakerName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.SpeakerName = raw

	return nil
}
