// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	JwtBearerAuthScopes = "jwtBearerAuth.Scopes"
)

// SignJwsRequest defines model for SignJwsRequest.
type SignJwsRequest struct {
	// Detached In detached mode the payload is signed but NOT included in the returned JWS object. Instead, the space between the first and second dot is empty, like this: "<header>..<signature>". Defaults to false.
	Detached *bool `json:"detached,omitempty"`

	// Headers The map of protected headers
	Headers map[string]interface{} `json:"headers"`

	// Kid Reference to the key ID used for signing the JWS.
	Kid string `json:"kid"`

	// Payload The payload to be signed as bytes. The bytes must be encoded with base64 encoding.
	Payload []byte `json:"payload"`
}

// SignJwtRequest defines model for SignJwtRequest.
type SignJwtRequest struct {
	Claims map[string]interface{} `json:"claims"`
	Kid    string                 `json:"kid"`
}

// SignJwsJSONRequestBody defines body for SignJws for application/json ContentType.
type SignJwsJSONRequestBody = SignJwsRequest

// SignJwtJSONRequestBody defines body for SignJwt for application/json ContentType.
type SignJwtJSONRequestBody = SignJwtRequest

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// SignJws request with any body
	SignJwsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SignJws(ctx context.Context, body SignJwsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SignJwt request with any body
	SignJwtWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SignJwt(ctx context.Context, body SignJwtJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) SignJwsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSignJwsRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SignJws(ctx context.Context, body SignJwsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSignJwsRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SignJwtWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSignJwtRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SignJwt(ctx context.Context, body SignJwtJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSignJwtRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewSignJwsRequest calls the generic SignJws builder with application/json body
func NewSignJwsRequest(server string, body SignJwsJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewSignJwsRequestWithBody(server, "application/json", bodyReader)
}

// NewSignJwsRequestWithBody generates requests for SignJws with any type of body
func NewSignJwsRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/crypto/v1/sign_jws")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewSignJwtRequest calls the generic SignJwt builder with application/json body
func NewSignJwtRequest(server string, body SignJwtJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewSignJwtRequestWithBody(server, "application/json", bodyReader)
}

// NewSignJwtRequestWithBody generates requests for SignJwt with any type of body
func NewSignJwtRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/crypto/v1/sign_jwt")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// SignJws request with any body
	SignJwsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignJwsResponse, error)

	SignJwsWithResponse(ctx context.Context, body SignJwsJSONRequestBody, reqEditors ...RequestEditorFn) (*SignJwsResponse, error)

	// SignJwt request with any body
	SignJwtWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignJwtResponse, error)

	SignJwtWithResponse(ctx context.Context, body SignJwtJSONRequestBody, reqEditors ...RequestEditorFn) (*SignJwtResponse, error)
}

type SignJwsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *struct {
		// Detail A human-readable explanation specific to this occurrence of the problem.
		Detail string `json:"detail"`

		// Status HTTP statuscode
		Status float32 `json:"status"`

		// Title A short, human-readable summary of the problem type.
		Title string `json:"title"`
	}
}

// Status returns HTTPResponse.Status
func (r SignJwsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SignJwsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SignJwtResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *struct {
		// Detail A human-readable explanation specific to this occurrence of the problem.
		Detail string `json:"detail"`

		// Status HTTP statuscode
		Status float32 `json:"status"`

		// Title A short, human-readable summary of the problem type.
		Title string `json:"title"`
	}
}

// Status returns HTTPResponse.Status
func (r SignJwtResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SignJwtResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// SignJwsWithBodyWithResponse request with arbitrary body returning *SignJwsResponse
func (c *ClientWithResponses) SignJwsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignJwsResponse, error) {
	rsp, err := c.SignJwsWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignJwsResponse(rsp)
}

func (c *ClientWithResponses) SignJwsWithResponse(ctx context.Context, body SignJwsJSONRequestBody, reqEditors ...RequestEditorFn) (*SignJwsResponse, error) {
	rsp, err := c.SignJws(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignJwsResponse(rsp)
}

// SignJwtWithBodyWithResponse request with arbitrary body returning *SignJwtResponse
func (c *ClientWithResponses) SignJwtWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignJwtResponse, error) {
	rsp, err := c.SignJwtWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignJwtResponse(rsp)
}

func (c *ClientWithResponses) SignJwtWithResponse(ctx context.Context, body SignJwtJSONRequestBody, reqEditors ...RequestEditorFn) (*SignJwtResponse, error) {
	rsp, err := c.SignJwt(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignJwtResponse(rsp)
}

// ParseSignJwsResponse parses an HTTP response from a SignJwsWithResponse call
func ParseSignJwsResponse(rsp *http.Response) (*SignJwsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SignJwsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest struct {
			// Detail A human-readable explanation specific to this occurrence of the problem.
			Detail string `json:"detail"`

			// Status HTTP statuscode
			Status float32 `json:"status"`

			// Title A short, human-readable summary of the problem type.
			Title string `json:"title"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseSignJwtResponse parses an HTTP response from a SignJwtWithResponse call
func ParseSignJwtResponse(rsp *http.Response) (*SignJwtResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SignJwtResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest struct {
			// Detail A human-readable explanation specific to this occurrence of the problem.
			Detail string `json:"detail"`

			// Status HTTP statuscode
			Status float32 `json:"status"`

			// Title A short, human-readable summary of the problem type.
			Title string `json:"title"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// sign a payload and headers with the private key of the given kid into a JWS object
	// (POST /internal/crypto/v1/sign_jws)
	SignJws(ctx echo.Context) error
	// sign a JWT payload with the private key of the given kid
	// (POST /internal/crypto/v1/sign_jwt)
	SignJwt(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// SignJws converts echo context to params.
func (w *ServerInterfaceWrapper) SignJws(ctx echo.Context) error {
	var err error

	ctx.Set(JwtBearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SignJws(ctx)
	return err
}

// SignJwt converts echo context to params.
func (w *ServerInterfaceWrapper) SignJwt(ctx echo.Context) error {
	var err error

	ctx.Set(JwtBearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SignJwt(ctx)
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

	router.POST(baseURL+"/internal/crypto/v1/sign_jws", wrapper.SignJws)
	router.POST(baseURL+"/internal/crypto/v1/sign_jwt", wrapper.SignJwt)

}

type SignJwsRequestObject struct {
	Body *SignJwsJSONRequestBody
}

type SignJwsResponseObject interface {
	VisitSignJwsResponse(w http.ResponseWriter) error
}

type SignJws200TextResponse string

func (response SignJws200TextResponse) VisitSignJwsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)

	_, err := w.Write([]byte(response))
	return err
}

type SignJwsdefaultJSONResponse struct {
	Body struct {
		// Detail A human-readable explanation specific to this occurrence of the problem.
		Detail string `json:"detail"`

		// Status HTTP statuscode
		Status float32 `json:"status"`

		// Title A short, human-readable summary of the problem type.
		Title string `json:"title"`
	}
	StatusCode int
}

func (response SignJwsdefaultJSONResponse) VisitSignJwsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type SignJwtRequestObject struct {
	Body *SignJwtJSONRequestBody
}

type SignJwtResponseObject interface {
	VisitSignJwtResponse(w http.ResponseWriter) error
}

type SignJwt200TextResponse string

func (response SignJwt200TextResponse) VisitSignJwtResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)

	_, err := w.Write([]byte(response))
	return err
}

type SignJwtdefaultJSONResponse struct {
	Body struct {
		// Detail A human-readable explanation specific to this occurrence of the problem.
		Detail string `json:"detail"`

		// Status HTTP statuscode
		Status float32 `json:"status"`

		// Title A short, human-readable summary of the problem type.
		Title string `json:"title"`
	}
	StatusCode int
}

func (response SignJwtdefaultJSONResponse) VisitSignJwtResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// sign a payload and headers with the private key of the given kid into a JWS object
	// (POST /internal/crypto/v1/sign_jws)
	SignJws(ctx context.Context, request SignJwsRequestObject) (SignJwsResponseObject, error)
	// sign a JWT payload with the private key of the given kid
	// (POST /internal/crypto/v1/sign_jwt)
	SignJwt(ctx context.Context, request SignJwtRequestObject) (SignJwtResponseObject, error)
}

type StrictHandlerFunc func(ctx echo.Context, args interface{}) (interface{}, error)

type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// SignJws operation middleware
func (sh *strictHandler) SignJws(ctx echo.Context) error {
	var request SignJwsRequestObject

	var body SignJwsJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.SignJws(ctx.Request().Context(), request.(SignJwsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "SignJws")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(SignJwsResponseObject); ok {
		return validResponse.VisitSignJwsResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// SignJwt operation middleware
func (sh *strictHandler) SignJwt(ctx echo.Context) error {
	var request SignJwtRequestObject

	var body SignJwtJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.SignJwt(ctx.Request().Context(), request.(SignJwtRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "SignJwt")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(SignJwtResponseObject); ok {
		return validResponse.VisitSignJwtResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}
