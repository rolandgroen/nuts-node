/*
 * Nuts node
 * Copyright (C) 2021 Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package v1

import (
	"github.com/google/uuid"
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-node/vcr/presentation"
	"github.com/nuts-foundation/nuts-node/vcr/proof"
	"net/http"
	"time"

	ssi "github.com/nuts-foundation/go-did"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-node/core"
	"github.com/nuts-foundation/nuts-node/vcr"
	"github.com/nuts-foundation/nuts-node/vcr/concept"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
	"github.com/nuts-foundation/nuts-node/vdr/types"
)

var _ ServerInterface = (*Wrapper)(nil)
var _ ErrorStatusCodeResolver = (*Wrapper)(nil)

// Wrapper implements the generated interface from oapi-codegen
type Wrapper struct {
	R  vcr.VCR
	CR concept.Reader
}

func (w *Wrapper) CreateVerifiablePresentation(ctx echo.Context) error {
	request := &CreateVerifiablePresentationRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if len(request.VerifiableCredentials) == 0 {
		return core.InvalidInputError("verifiableCredentials needs at least 1 item")
	}

	method, err := ssi.ParseURI(request.VerificationMethod)
	if err != nil {
		return core.InvalidInputError("could not parse the verificationMethod: %w", err)
	}
	if method.Fragment == "" {
		return core.InvalidInputError("missing key identifier in verificationMethod")
	}

	for _, c := range request.VerifiableCredentials {
		if err := w.R.Validate(c, false, true, nil); err != nil {
			return core.InvalidInputError("invalid credential: %w", err)
		}
	}

	// Set created to now
	created := time.Now()

	var (
		expires *time.Time
		nonce   *string
	)

	if request.Expires != nil {
		if request.Expires.Before(created) {
			return core.InvalidInputError("expires can not lay in the past")
		}
		expires = request.Expires

		// nonce Needs to be unique so the verifier can store it during the duration the presentation is valid.
		// Only set nonce when request.expires is set, otherwise the verifier must store the nonce indefinitely.
		uuidNonce := uuid.New().String()
		nonce = &uuidNonce
	}

	proofOptions := proof.ProofOptions{
		Created:   created,
		Domain:    request.Domain,
		Nonce:     nonce,
		Challenge: request.Challenge,
		Type:      proof.JsonWebSignature2020,
		Expires:   expires,
	}

	vp := presentation.VerifiablePresentation{
		Context:              []string{"https://www.w3.org/2018/credentials/v1"},
		Type:                 []string{"VerifiablePresentation"},
		VerifiableCredential: &request.VerifiableCredentials,
	}
	err = w.R.GeneratePresentationEmbeddedProof(&vp, proofOptions, request.VerificationMethod)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, vp)
}

// ResolveStatusCode maps errors returned by this API to specific HTTP status codes.
func (w *Wrapper) ResolveStatusCode(err error) int {
	return core.ResolveStatusCode(err, map[error]int{
		concept.ErrUnknownConcept: http.StatusNotFound,
		vcr.ErrNotFound:           http.StatusNotFound,
		vcr.ErrRevoked:            http.StatusConflict,
		credential.ErrValidation:  http.StatusBadRequest,
		types.ErrNotFound:         http.StatusBadRequest,
		types.ErrKeyNotFound:      http.StatusBadRequest,
		vcr.ErrInvalidCredential:  http.StatusNotFound,
	})
}

// Preprocess is called just before the API operation itself is invoked.
func (w *Wrapper) Preprocess(operationID string, context echo.Context) {
	context.Set(core.StatusCodeResolverContextKey, w)
	context.Set(core.OperationIDContextKey, operationID)
	context.Set(core.ModuleNameContextKey, "VCR")
}

// Routes registers the handler to the echo router
func (w *Wrapper) Routes(router core.EchoRouter) {
	RegisterHandlers(router, w)
}

// Search finds concepts. Concepts are mapped to VCs. This is primarily used for finding DIDs.
func (w *Wrapper) Search(ctx echo.Context, conceptName string, requestParams SearchParams) error {
	sr := new(SearchRequest)
	if err := ctx.Bind(sr); err != nil {
		return err
	}

	params := make(map[string]string, len(sr.Params))
	for _, pair := range sr.Params {
		params[pair.Key] = pair.Value
	}

	untrusted := false
	if requestParams.Untrusted != nil {
		untrusted = *requestParams.Untrusted
	}

	results, err := w.R.Search(ctx.Request().Context(), conceptName, untrusted, params)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, results)
}

// Revoke a credential
func (w *Wrapper) Revoke(ctx echo.Context, id string) error {
	idURI, err := ssi.ParseURI(id)

	// return 400 for malformed input
	if err != nil {
		return core.InvalidInputError("failed to parse credential ID: %w", err)
	}

	r, err := w.R.Revoke(*idURI)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// Create a Verifiable credential
func (w *Wrapper) Create(ctx echo.Context) error {
	requestedVC := IssueVCRequest{}
	if err := ctx.Bind(&requestedVC); err != nil {
		return core.InvalidInputError("invalid input: %w", err)
	}

	if len(requestedVC.Type) != 1 {
		return core.InvalidInputError("type should contain 1 value")
	}
	if _, err := did.ParseDID(requestedVC.Issuer.String()); err != nil {
		return core.InvalidInputError("invalid issuer: %w", err)
	}

	vcCreated, err := w.R.Issue(requestedVC)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, vcCreated)
}

func (w *Wrapper) Resolve(ctx echo.Context, id string, params ResolveParams) error {
	idURI, err := ssi.ParseURI(id)
	// return 400 for malformed input
	if err != nil {
		return core.InvalidInputError("failed to parse credential ID: %w", err)
	}

	// resolve time
	var at *time.Time
	if params.ResolveTime != nil {
		parsedTime, err := time.Parse(time.RFC3339, *params.ResolveTime)
		if err != nil {
			return core.InvalidInputError("failed to parse query parameter 'at': %w", err)
		}
		at = &parsedTime
	}

	// id is given with fragment
	vc, err := w.R.Resolve(*idURI, at)
	if vc == nil && err != nil {
		return err
	}

	// transform VC && error
	result := ResolutionResult{
		CurrentStatus:        ResolutionResultCurrentStatusTrusted,
		VerifiableCredential: *vc,
	}

	switch err {
	case vcr.ErrUntrusted:
		result.CurrentStatus = ResolutionResultCurrentStatusUntrusted
	case vcr.ErrRevoked:
		result.CurrentStatus = ResolutionResultCurrentStatusRevoked
	}

	return ctx.JSON(http.StatusOK, result)
}

func (w *Wrapper) TrustIssuer(ctx echo.Context) error {
	return changeTrust(ctx, func(cType ssi.URI, issuer ssi.URI) error {
		return w.R.Trust(cType, issuer)
	})
}

func (w *Wrapper) UntrustIssuer(ctx echo.Context) error {
	return changeTrust(ctx, func(cType ssi.URI, issuer ssi.URI) error {
		return w.R.Untrust(cType, issuer)
	})
}

func (w *Wrapper) ListTrusted(ctx echo.Context, credentialType string) error {
	uri, err := parseCredentialType(credentialType)
	if err != nil {
		return err
	}

	trusted, err := w.R.Trusted(*uri)
	if err != nil {
		return err
	}
	result := make([]string, len(trusted))
	for i, t := range trusted {
		result[i] = t.String()
	}

	return ctx.JSON(http.StatusOK, result)
}

func (w *Wrapper) ListUntrusted(ctx echo.Context, credentialType string) error {
	uri, err := parseCredentialType(credentialType)
	if err != nil {
		return err
	}

	untrusted, err := w.R.Untrusted(*uri)
	if err != nil {
		return err
	}

	result := make([]string, len(untrusted))
	for i, t := range untrusted {
		result[i] = t.String()
	}

	return ctx.JSON(http.StatusOK, result)
}

func parseCredentialType(credentialType string) (*ssi.URI, error) {
	uri, err := ssi.ParseURI(credentialType)
	if err != nil {
		return nil, core.InvalidInputError("malformed credential type: %w", err)
	}
	return uri, nil
}

type trustChangeFunc func(ssi.URI, ssi.URI) error

func changeTrust(ctx echo.Context, f trustChangeFunc) error {
	var icc = new(CredentialIssuer)

	if err := ctx.Bind(icc); err != nil {
		return err
	}

	d, err := ssi.ParseURI(icc.Issuer)
	if err != nil {
		return core.InvalidInputError("failed to parse issuer: %w", err)
	}

	cType, err := parseCredentialType(icc.CredentialType)
	if err != nil {
		return err
	}

	if err = f(*cType, *d); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
