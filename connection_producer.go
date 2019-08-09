package mockdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/mitchellh/mapstructure"
)

// mockdbConnectionProducer implements ConnectionProducer and provides an
// interface for mockdb databases to make connections.
type githubdkConnectionProducer struct {
	Url     string `json:"url" structs:"url" mapstructure:"url"`
	ApiToken string `json:"apitoken" structs:"apitoken" mapstructure:"apitoken"`


	rawConfig map[string]interface{}

	Initialized bool
	Type        string
	client      string

	sync.Mutex
}

func (i *githubdkConnectionProducer) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := i.Init(ctx, conf, verifyConnection)
	return err
}

func (i *githubdkConnectionProducer) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
	i.Lock()
	defer i.Unlock()

	i.rawConfig = conf

	err := mapstructure.WeakDecode(conf, i)
	if err != nil {
		return nil, err
	}

	if i.Port == "" {
		i.Port = "8086"
	}

	switch {
	case len(i.Url) == 0:
		return nil, fmt.Errorf("url cannot be empty")
	case len(i.ApiToken) == 0:
		return nil, fmt.Errorf("apitoken cannot be empty")
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	i.Initialized = true

	if verifyConnection {
		if _, err := i.Connection(ctx); err != nil {
			return nil, errwrap.Wrapf("error verifying connection: {{err}}", err)
		}
	}

	return conf, nil
}

func (i *githubdkConnectionProducer) Connection(_ context.Context) (interface{}, error) {
	if !i.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	return nil, nil
}

func (i *githubdkConnectionProducer) Close() error {
	// Grab the write lock
	i.Lock()
	defer i.Unlock()

	return nil
}

func (i *githubdkConnectionProducer) secretValues() map[string]interface{} {
	return map[string]interface{}{
		i.Password: "[password]",
	}
}

// SetCredentials uses provided information to set/create a user in the
// database. Unlike CreateUser, this method requires a username be provided and
// uses the name given, instead of generating a name. This is used for creating
// and setting the password of static accounts, as well as rolling back
// passwords in the database in the event an updated database fails to save in
// Vault's storage.
func (i *githubdkConnectionProducer) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticUser dbplugin.StaticUserConfig) (username, password string, err error) {
	return "", "", dbutil.Unimplemented()
}
