// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PressKeyHandlerFunc turns a function with the right signature into a press key handler
type PressKeyHandlerFunc func(PressKeyParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PressKeyHandlerFunc) Handle(params PressKeyParams) middleware.Responder {
	return fn(params)
}

// PressKeyHandler interface for that can handle valid press key params
type PressKeyHandler interface {
	Handle(PressKeyParams) middleware.Responder
}

// NewPressKey creates a new http.Handler for the press key operation
func NewPressKey(ctx *middleware.Context, handler PressKeyHandler) *PressKey {
	return &PressKey{Context: ctx, Handler: handler}
}

/*PressKey swagger:route GET /{speakerName}/key/{keyId} pressKey

Presses and releases a key on selected deviceId

*/
type PressKey struct {
	Context *middleware.Context
	Handler PressKeyHandler
}

func (o *PressKey) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPressKeyParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
