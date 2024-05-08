package vault

import (
	"context"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/userpass"
)

type vault struct {
	addr      string
	paths     []string
	user      string
	pass      string
	namespace string
}

func NewSource() config.Source {
	paths := strings.Split(os.Getenv("VAULT_PATHS"), ",")
	addr := os.Getenv("VAULT_HOST")
	user := os.Getenv("VAULT_USER")
	password := os.Getenv("VAULT_PASSWORD")
	namespace := os.Getenv("VAULT_NAMESPACE")
	return &vault{addr: addr, paths: paths, user: user, pass: password, namespace: namespace}
}

func (e *vault) Load() (kv []*config.KeyValue, err error) {
	return e.load()
}

func (e *vault) load() ([]*config.KeyValue, error) {
	ctx := context.Background()
	vaultConfig := vaultapi.DefaultConfig()

	vaultConfig.Address = e.addr

	client, err := vaultapi.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}

	userpass, err := userpass.NewUserpassAuth(e.user, &userpass.Password{FromString: e.pass})
	if err != nil {
		return nil, err
	}
	auth, err := client.Auth().Login(ctx, userpass)
	if err != nil {
		return nil, err
	}

	client.SetToken(auth.Auth.ClientToken)

	// read the secret
	configs := make(map[string]interface{})
	for _, path := range e.paths {
		s, err := client.KVv1("stg").Get(ctx, path)
		if err != nil {
			return nil, err
		}
		for k, v := range s.Data {
			configs[k] = v
		}
	}

	var kv []*config.KeyValue
	for k, v := range configs {
		if len(k) != 0 {
			kv = append(kv, &config.KeyValue{
				Key:   k,
				Value: []byte(v.(string)),
			})
		}
	}
	return kv, nil
}

func (e *vault) Watch() (config.Watcher, error) {
	w, err := NewWatcher()
	if err != nil {
		return nil, err
	}
	return w, nil
}
