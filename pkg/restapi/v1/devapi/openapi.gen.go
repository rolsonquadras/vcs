// Package devapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package devapi

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// DID Config response.
type DidConfig struct {
	// context.
	Context *string `json:"@context,omitempty"`

	// Presentation in jws(string) or jsonld(object) formats
	LinkedDids *[]interface{} `json:"linked_dids,omitempty"`
}

// DID Config response.
type RequestObject struct {
	// Content of requested object
	Content *string `json:"content,omitempty"`

	// id.
	Id *string `json:"id,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Request request object by uuid
	// (GET /request-object/{uuid})
	RequestObjectByUuid(ctx echo.Context, uuid string) error
	// Request did-config
	// (GET /{profileType}/profiles/{profileID}/well-known/did-config)
	DidConfig(ctx echo.Context, profileType string, profileID string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// RequestObjectByUuid converts echo context to params.
func (w *ServerInterfaceWrapper) RequestObjectByUuid(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "uuid" -------------
	var uuid string

	err = runtime.BindStyledParameterWithLocation("simple", false, "uuid", runtime.ParamLocationPath, ctx.Param("uuid"), &uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter uuid: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RequestObjectByUuid(ctx, uuid)
	return err
}

// DidConfig converts echo context to params.
func (w *ServerInterfaceWrapper) DidConfig(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "profileType" -------------
	var profileType string

	err = runtime.BindStyledParameterWithLocation("simple", false, "profileType", runtime.ParamLocationPath, ctx.Param("profileType"), &profileType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter profileType: %s", err))
	}

	// ------------- Path parameter "profileID" -------------
	var profileID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "profileID", runtime.ParamLocationPath, ctx.Param("profileID"), &profileID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter profileID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DidConfig(ctx, profileType, profileID)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/request-object/:uuid", wrapper.RequestObjectByUuid)
	router.GET(baseURL+"/:profileType/profiles/:profileID/well-known/did-config", wrapper.DidConfig)

}
