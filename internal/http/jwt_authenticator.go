package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt"
	middleware "github.com/oapi-codegen/echo-middleware"
	echo2 "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

// JWSValidator is used to validate JWS payloads and return a JWT if they're
// valid
type JWSValidator interface {
	ValidateJWS(jws string) (jwt.Token, error)
}

const JWTClaimsContextKey = "subject"

var userIDKey = "subject"

var (
	ErrNoAuthHeader      = errors.New("Authorization header is missing")
	ErrInvalidAuthHeader = errors.New("Authorization header is malformed")
	ErrClaimsInvalid     = errors.New("Provided claims do not match expected scopes")
)

// GetJWSFromRequest extracts a JWS string from an Authorization: Bearer <jws> header
func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	// Check for the Authorization header.
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	// We expect a header value of the form "Bearer <token>", with 1 space after
	// Bearer, per spec.
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

func NewAuthenticator(v JWSValidator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(v, ctx, input)
	}
}

// Authenticate uses the specified validator to ensure a JWT is valid, then makes
// sure that the claims provided by the JWT match the scopes as required in the API.
func Authenticate(v JWSValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	// Our security scheme is named BearerAuth, ensure this is the case
	if input.SecuritySchemeName != "BearerAuth" {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid security scheme")
	}

	// Now, we need to get the JWS from the request, to match the request expectations
	// against request contents.
	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid JWS")
	}

	// if the JWS is valid, we have a JWT, which will contain a bunch of claims.
	token, err := v.ValidateJWS(jws)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid JWS")
	}

	// We've got a valid token now, and we can look into its claims to see whether
	// they match. Every single scope must be present in the claims.
	//err = CheckTokenClaims(input.Scopes, token)

	//if err != nil {
	//	return fmt.Errorf("token claims don't match: %w", err)
	//}

	// Set the property on the echo context so the handler is able to
	// access the claims data we generate in here.

	eCtx := middleware.GetEchoContext(ctx)
	eCtx.Set(userIDKey, token.Subject())

	//input.RequestValidationInput.Request.Context()

	//sub := token.Subject()
	//rCtx := eCtx.Request().Context()
	//ctx = context.WithValue(rCtx, userIDKey, sub)
	////input.RequestValidationInput.Request.WithContext(ctx)
	//
	//eCtx.SetRequest(eCtx.Request().WithContext(ctx))
	//
	//v1 := eCtx.Request().Context().Value(userIDKey)
	//fmt.Println(v1)

	return nil
}

var userIDKeyStruct = struct{}{}

func StrictMiddlewareUserIDTransfer(f echo2.StrictEchoHandlerFunc, operationID string) echo2.StrictEchoHandlerFunc {
	return func(ctx echo.Context, request interface{}) (response interface{}, err error) {
		value := ctx.Get(userIDKey)
		if value == nil {
		} else {
			if _, ok := value.(string); !ok {
				return nil, fmt.Errorf("user_id is not a string")
			}
		}

		rCtx := ctx.Request().Context()
		rCtx = context.WithValue(rCtx, userIDKeyStruct, value)
		ctx.SetRequest(ctx.Request().WithContext(rCtx))

		return f(ctx, request)
	}
}
