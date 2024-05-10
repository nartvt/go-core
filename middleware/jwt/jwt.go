package jwt

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type authKey struct{}

type String string

const (

	// bearerWord the bearer key word for authorization
	bearerWord String = "Bearer"

	// bearerFormat authorization token format
	bearerFormat String = "Bearer %s"

	// authorizationKey holds the key used to store the JWT Token in the request tokenHeader.
	authorizationKey String = "Authorization"

	// reason holds the error reason.
	reason String = "UNAUTHORIZED"
)

var (
	ErrMissingJwtToken        = errors.Unauthorized(string(reason), "JWT token is missing")
	ErrMissingKeyFunc         = errors.Unauthorized(string(reason), "keyFunc is missing")
	ErrTokenInvalid           = errors.Unauthorized(string(reason), "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized(string(reason), "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized(string(reason), "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized(string(reason), "Wrong signing method")
	ErrWrongContext           = errors.Unauthorized(string(reason), "Wrong context for middleware")
	ErrNeedTokenProvider      = errors.Unauthorized(string(reason), "Token provider is missing")
	ErrSignToken              = errors.Unauthorized(string(reason), "Can not sign token.Is the key correct?")
	ErrGetKey                 = errors.Unauthorized(string(reason), "Can not get key while signing token")
)

// Option is jwt option.
type Option func(*options)

// Parser is a jwt parser
type options struct {
	required      bool
	signingMethod jwtlib.SigningMethod
	excludes      []string
	autoParse     bool
	claims        func() jwtlib.Claims
	tokenHeader   map[string]interface{}
	key           string
	keyFunc       jwtlib.Keyfunc
}

func WithRequired(required bool) Option {
	return func(o *options) {
		o.required = required
	}
}

func WithAutoParse(autoParse bool) Option {
	return func(o *options) {
		o.autoParse = autoParse
	}
}

func WithExcludes(excludes string) Option {
	return func(o *options) {
		if len(excludes) > 0 {
			o.excludes = strings.Split(excludes, ",")
		}
	}
}

// WithSigningMethod with signing method option.
func WithSigningMethod(method jwtlib.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// WithClaims with customer claim
// If you use it in Server, f needs to return a new jwt.Claims object each time to avoid concurrent write problems
// If you use it in Client, f only needs to return a single object to provide performance
func WithClaims(f func() jwtlib.Claims) Option {
	return func(o *options) {
		o.claims = f
	}
}

// WithTokenHeader with customer tokenHeader for client side
func WithTokenHeader(header map[string]interface{}) Option {
	return func(o *options) {
		o.tokenHeader = header
	}
}

func WithSecretKey(key string) Option {
	return func(o *options) {
		o.key = key
	}
}

func WithCustomKeyFunc(keyFunc jwtlib.Keyfunc) Option {
	return func(o *options) {
		o.keyFunc = keyFunc
	}
}

// Server is a server auth middleware. Check the token and extract the info from token.
func Server(opts ...Option) middleware.Middleware {
	secretKey := os.Getenv("ENV_JWT_SECRET")
	o := &options{
		signingMethod: jwtlib.SigningMethodHS256,
		key:           secretKey,
		required:      false,
	}
	for _, opt := range opts {
		opt(o)
	}

	if o.keyFunc == nil {
		o.keyFunc = func(t *jwtlib.Token) (interface{}, error) {
			if t.Header["alg"] == o.signingMethod.Alg() {
				return []byte(o.key), nil
			}
			return nil, ErrMissingKeyFunc
		}
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				var (
					tokenInfo *jwtlib.Token
					authErr   error
				)
				if o.autoParse {
					tokenAuth := header.RequestHeader().Get(string(authorizationKey))
					tokenInfo, authErr = parseToken(tokenAuth, o)
					if tokenInfo != nil {
						ctx = NewContext(ctx, tokenInfo.Claims)
						ctx = context.WithValue(ctx, authorizationKey, tokenInfo.Raw)
					}
				}

				if !o.required {
					return handler(ctx, req)
				}

				if len(o.excludes) > 0 {
					for _, exclude := range o.excludes {
						if strings.HasPrefix(header.Operation(), exclude) {
							return handler(ctx, req)
						}
					}
				}

				if authErr != nil {
					return nil, authErr
				}

				if tokenInfo == nil {
					tokenAuth := header.RequestHeader().Get(string(authorizationKey))
					tokenInfo, err := parseToken(tokenAuth, o)
					if err != nil {
						return nil, err
					}

					ctx = NewContext(ctx, tokenInfo.Claims)
					ctx = context.WithValue(ctx, authorizationKey, tokenInfo.Raw)
				}
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

func parseToken(token string, o *options) (*jwtlib.Token, error) {
	auths := strings.SplitN(token, " ", 2)
	if len(auths) != 2 || !strings.EqualFold(auths[0], string(bearerWord)) {
		return nil, ErrMissingJwtToken
	}
	jwtToken := auths[1]
	var (
		tokenInfo *jwtlib.Token
		err       error
	)
	if o.claims != nil {
		tokenInfo, err = jwtlib.ParseWithClaims(jwtToken, o.claims(), o.keyFunc)
	} else {
		tokenInfo, err = jwtlib.Parse(jwtToken, o.keyFunc)
	}
	if err != nil {
		return nil, ErrTokenInvalid
	}
	if !tokenInfo.Valid {
		return nil, ErrTokenInvalid
	}
	if tokenInfo.Method != o.signingMethod {
		return nil, ErrUnSupportSigningMethod
	}
	return tokenInfo, err
}

// Client is a client jwt middleware.
func Client(opts ...Option) middleware.Middleware {
	claims := jwtlib.RegisteredClaims{}
	secretKey := os.Getenv("ENV_JWT_SECRET")
	o := &options{
		signingMethod: jwtlib.SigningMethodHS256,
		key:           secretKey,
		required:      false,
		claims:        func() jwtlib.Claims { return claims },
	}
	for _, opt := range opts {
		opt(o)
	}

	if o.keyFunc == nil {
		o.keyFunc = func(t *jwtlib.Token) (interface{}, error) {
			if t.Header["alg"] == o.signingMethod.Alg() {
				return []byte(o.key), nil
			}
			return nil, ErrMissingKeyFunc
		}
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tokenStr, ok := ctx.Value(authorizationKey).(string); ok {
				if clientContext, ok := transport.FromClientContext(ctx); ok {
					clientContext.RequestHeader().Set(string(authorizationKey), fmt.Sprintf(string(bearerFormat), tokenStr))
					return handler(ctx, req)
				}
				return nil, ErrWrongContext
			}

			token := jwtlib.NewWithClaims(o.signingMethod, o.claims())
			if o.tokenHeader != nil {
				for k, v := range o.tokenHeader {
					token.Header[k] = v
				}
			}
			key, err := o.keyFunc(token)
			if err != nil {
				return nil, ErrGetKey
			}
			tokenStr, err := token.SignedString(key)
			if err != nil {
				return nil, ErrSignToken
			}
			if clientContext, ok := transport.FromClientContext(ctx); ok {
				clientContext.RequestHeader().Set(string(authorizationKey), fmt.Sprintf(string(bearerFormat), tokenStr))
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

// NewContext put auth info into context
func NewContext(ctx context.Context, info jwtlib.Claims) context.Context {
	return context.WithValue(ctx, authKey{}, info)
}

// FromContext extract auth info from context
func FromContext(ctx context.Context) (token jwtlib.Claims, ok bool) {
	token, ok = ctx.Value(authKey{}).(jwtlib.Claims)
	return
}

func GetUserId(ctx context.Context) (string, error) {
	value, ok := FromContext(ctx)

	if !ok {
		return "", nil
	}

	userId, err := value.GetSubject()
	if err != nil {
		return "", err
	}

	return userId, nil
}

func GetEmail(ctx context.Context) string {
	value, ok := FromContext(ctx)

	if !ok {
		return ""
	}

	email := value.(jwtlib.MapClaims)["email"]
	if email != nil {
		return email.(string)
	}

	return ""
}

// Client is a client jwt middleware.
func ClientGrpcAuth(ctx context.Context, opts ...Option) (context.Context, error) {
	meta := metadata.MD{}
	claims := jwtlib.RegisteredClaims{}
	secretKey := os.Getenv("ENV_JWT_SECRET")
	o := &options{
		signingMethod: jwtlib.SigningMethodHS256,
		key:           secretKey,
		required:      false,
		claims:        func() jwtlib.Claims { return claims },
	}
	for _, opt := range opts {
		opt(o)
	}

	if o.keyFunc == nil {
		o.keyFunc = func(t *jwtlib.Token) (interface{}, error) {
			if t.Header["alg"] == o.signingMethod.Alg() {
				return []byte(o.key), nil
			}
			return nil, ErrMissingKeyFunc
		}
	}

	if tokenStr, ok := ctx.Value(authorizationKey).(string); ok {
		meta.Set(string(authorizationKey), fmt.Sprintf(string(bearerFormat), tokenStr))
		return metadata.NewOutgoingContext(ctx, meta), nil
	}

	token := jwtlib.NewWithClaims(o.signingMethod, o.claims())
	if o.tokenHeader != nil {
		for k, v := range o.tokenHeader {
			token.Header[k] = v
		}
	}
	key, err := o.keyFunc(token)
	if err != nil {
		return nil, ErrMissingKeyFunc
	}
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return nil, ErrSignToken
	}
	meta.Set(string(authorizationKey), fmt.Sprintf(string(bearerFormat), tokenStr))
	return metadata.NewOutgoingContext(ctx, meta), nil
}

func GeneratorJwtToken(ctx context.Context, userId string) (context.Context, error) {
	claims := jwtlib.RegisteredClaims{
		Subject: userId,
		ExpiresAt: jwtlib.NewNumericDate(
			time.Now().Add(24 * time.Hour),
		),
	}
	secretKey := os.Getenv("ENV_JWT_SECRET")
	o := &options{
		signingMethod: jwtlib.SigningMethodHS256,
		key:           secretKey,
		required:      false,
		claims:        func() jwtlib.Claims { return claims },
		tokenHeader: map[string]interface{}{
			"typ": "JWT",
			"alg": "HS256",
		},
		keyFunc: func(token *jwtlib.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	}

	token := jwtlib.NewWithClaims(o.signingMethod, o.claims())
	if o.tokenHeader != nil {
		for k, v := range o.tokenHeader {
			token.Header[k] = v
		}
	}
	key, err := o.keyFunc(token)
	if err != nil {
		return ctx, ErrMissingKeyFunc
	}
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return ctx, ErrMissingKeyFunc
	}

	return context.WithValue(ctx, "accessToken", tokenStr), nil
}
