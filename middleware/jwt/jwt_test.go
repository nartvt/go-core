package jwt

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/go-kratos/kratos/v2/transport"
)

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string { return http.Header(hc).Get(key) }

func (hc headerCarrier) Set(key string, value string) { http.Header(hc).Set(key, value) }

func (hc headerCarrier) Add(key string, value string) { http.Header(hc).Add(key, value) }

func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice value associated with the passed key.
func (hc headerCarrier) Values(key string) []string {
	return http.Header(hc).Values(key)
}

type Transport struct {
	kind      transport.Kind
	endpoint  string
	operation string
	reqHeader transport.Header
}

func (tr *Transport) Kind() transport.Kind {
	return tr.kind
}

func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

func (tr *Transport) Operation() string {
	return tr.operation
}

func (tr *Transport) RequestHeader() transport.Header {
	return tr.reqHeader
}

func (tr *Transport) ReplyHeader() transport.Header {
	return nil
}

var (
	requireAuthKey = "x-auth"
	requireAuthVal = "1"
)

func init() {
	os.Setenv("IND_JWT_SECRET", "YOURSECRETKEYGOESHERE")
}

func getActiveToken() string {
	secretKey := []byte(os.Getenv("IND_JWT_SECRET"))
	claim := jwtlib.RegisteredClaims{
		ExpiresAt: jwtlib.NewNumericDate(time.Now().AddDate(0, 0, 1)),
		IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		Subject:   "123456789",
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claim)
	ss, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}

	return ss

}

func getExpiredToken() string {
	secretKey := []byte(os.Getenv("IND_JWT_SECRET"))
	claim := jwtlib.RegisteredClaims{
		ExpiresAt: jwtlib.NewNumericDate(time.Now().AddDate(0, 0, -1)),
		IssuedAt:  jwtlib.NewNumericDate(time.Now().AddDate(0, 0, -2)),
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claim)
	ss, _ := token.SignedString(secretKey)
	return ss
}

func TestSever_WithoutAuth(t *testing.T) {

	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
		return nil, nil
	}
	hc := &headerCarrier{}
	//hc.Set(authKey, authVal)
	ctx := transport.NewServerContext(context.Background(), &Transport{reqHeader: hc})

	// const md
	_, err := Server()(hs)(ctx, "foo")
	require.Nil(t, err)
}

func TestSever_WithAuthErrNoKey(t *testing.T) {

	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
		return nil, nil
	}
	hc := headerCarrier{}
	//hc.Set(authKey, authVal)
	ctx := transport.NewServerContext(context.Background(), &Transport{reqHeader: hc})
	// const md
	_, err := Server(WithRequired(true))(hs)(ctx, "foo")
	require.Equal(t, err, ErrMissingJwtToken)
}

func TestSever_WithAuthTokenSuccess(t *testing.T) {
	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
		return nil, nil
	}
	hc := headerCarrier{}
	hc.Set("Authorization", "Bearer "+getActiveToken())
	ctx := transport.NewServerContext(context.Background(), &Transport{reqHeader: hc})
	// const md
	_, err := Server(WithRequired(true))(hs)(ctx, "foo")
	require.Nil(t, err)
}

func TestSever_WithAuthTokenNoBearer(t *testing.T) {
	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
		return nil, nil
	}
	hc := headerCarrier{}
	hc.Set("Authorization", getActiveToken())
	ctx := transport.NewServerContext(context.Background(), &Transport{reqHeader: hc})
	// const md
	_, err := Server(WithRequired(true))(hs)(ctx, "foo")
	require.Equal(t, err, ErrMissingJwtToken)
}

func TestSever_WithAuthWrongToken(t *testing.T) {
	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
		return nil, nil
	}
	hc := headerCarrier{}
	hc.Set("Authorization", "Bearer 123")
	ctx := transport.NewServerContext(context.Background(), &Transport{reqHeader: hc})
	// const md
	_, err := Server(WithRequired(true))(hs)(ctx, "foo")
	require.Equal(t, err, ErrTokenInvalid)
}

func TestSever_WithAuthTokenInActive(t *testing.T) {
	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
		return nil, nil
	}
	hc := headerCarrier{}
	hc.Set("Authorization", "Bearer "+getExpiredToken())
	ctx := transport.NewServerContext(context.Background(), &Transport{reqHeader: hc})
	// const md
	_, err := Server(WithRequired(true))(hs)(ctx, "foo")
	require.Equal(t, err, ErrTokenInvalid)
}

// func TestSever_WithAuthXUserSuccess(t *testing.T) {
// 	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
// 		return nil, nil
// 	}
// 	hc := headerCarrier{}
// 	exp := time.Now().Unix() + 10000
// 	iat := time.Now().Unix() - 10000
// 	hc.Set("X-User", fmt.Sprintf("{\"active\": true, \"exp\": %d, \"iat\": %d}", exp, iat))
// 	ctx := transport.NewServerContext(context.Background(), &testTransport{hc})
// 	// const md
// 	_, err := Server(Required(true))(hs)(ctx, "foo")
// 	require.Nil(t, err)
// }

// func TestSever_WithAuthXUserInActive(t *testing.T) {
// 	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
// 		return nil, nil
// 	}
// 	hc := headerCarrier{}
// 	exp := time.Now().Unix() + 10000
// 	iat := time.Now().Unix() - 10000
// 	hc.Set("X-User", fmt.Sprintf("{\"active\": false, \"exp\": %d,  \"iat\": %d}", exp, iat))
// 	ctx := transport.NewServerContext(context.Background(), &testTransport{hc})
// 	// const md
// 	_, err := Server(Required(true))(hs)(ctx, "foo")
// 	require.Equal(t, err, ErrTokenInvalid)
// }

// func TestSever_WithAuthForceTestFromClientError(t *testing.T) {
// 	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
// 		return nil, nil
// 	}
// 	hc := headerCarrier{}
// 	//exp := time.Now().Unix() + 10000
// 	//iat := time.Now().Unix() - 10000
// 	hc.Set("X-Auth", "1")
// 	//hc.Set("X-User", fmt.Sprintf("{\"active\": true, \"exp\": %d, \"iat\": %d}", exp, iat))
// 	ctx := transport.NewServerContext(context.Background(), &testTransport{hc})
// 	// const md
// 	_, err := Server(WithTest(true))(hs)(ctx, "foo")
// 	require.Equal(t, err, ErrMissingToken)
// }

// func TestSever_WithAuthForceTestFromClientSuccess(t *testing.T) {
// 	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
// 		return nil, nil
// 	}
// 	hc := headerCarrier{}
// 	exp := time.Now().Unix() + 10000
// 	iat := time.Now().Unix() - 10000
// 	hc.Set("X-Auth", "1")
// 	hc.Set("X-User", fmt.Sprintf("{\"active\": true, \"exp\": %d, \"iat\": %d}", exp, iat))
// 	ctx := transport.NewServerContext(context.Background(), &testTransport{hc})
// 	// const md
// 	_, err := Server(WithTest(true))(hs)(ctx, "foo")
// 	require.Nil(t, err)
// }

// func TestClient_SendUserCtxToOther(t *testing.T) {
// 	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
// 		authInfo, ok := FromContext(ctx)
// 		if !ok {
// 			return nil, errors.New("Missing auth info")
// 		}
// 		return authInfo, nil
// 	}
// 	hc := headerCarrier{}
// 	exp := time.Now().Unix() + 10000
// 	iat := time.Now().Unix() - 10000
// 	//hc.Set("X-Auth", "1")
// 	//hc.Set("X-User", fmt.Sprintf("{\"active\": true, \"exp\": %d, \"iat\": %d}", exp, iat))
// 	ctx := transport.NewClientContext(context.Background(), &testTransport{hc})
// 	active := true
// 	tokenInfo := &models.OAuth2TokenIntrospection{
// 		Active: &active,
// 		Exp:    exp,
// 		Iat:    iat,
// 	}
// 	ctx = NewContext(ctx, tokenInfo)
// 	// const md
// 	tokenResult, err := Client()(hs)(ctx, "foo")
// 	require.Nil(t, err)
// 	require.Equal(t, tokenResult, tokenInfo)
// }
