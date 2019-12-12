package keycloak

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type ServerConfig struct {
	Url          string
	Realm        string
	ClientID     string
	ClientSecret string
	RedirectUrl  string
}

type keycloak struct {
	oauth2Config *oauth2.Config
	provider     *oidc.Provider
}

type Keycloak interface {
	TokenRequest(ctx context.Context, username string, password string) (*oauth2.Token, error)
	TokenSource(context.Context, *oauth2.Token) oauth2.TokenSource
	UserInfo(context.Context, oauth2.TokenSource) (map[string]interface{}, error)
	Client(context.Context, *oauth2.Token) *http.Client
}

func NewClientContext(ctx context.Context, client *http.Client) context.Context {
	return oidc.ClientContext(ctx, client)
}

func NewConfig(ctx context.Context, config *ServerConfig) (Keycloak, error) {
	// Configure an OpenID Connect aware OAuth2 client.

	server := &keycloak{}

	provider, err := oidc.NewProvider(ctx, fmt.Sprintf("%s/realms/%s", config.Url, config.Realm))
	if err != nil {
		return nil, err
	}
	log.Printf("Endpoints: %+v: ", provider.Endpoint())

	server.provider = provider

	server.oauth2Config = &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID},
		//Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
		RedirectURL: config.RedirectUrl,
	}
	return server, nil
}

func (s *keycloak) TokenRequest(ctx context.Context, username, password string) (*oauth2.Token, error) {

	//log.Printf("Endpoints: %+v: ", provider.Endpoint())

	token, err := s.oauth2Config.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, err
	}

	return token, nil
}

/**/

// TokenSource returns a TokenSource that returns t until t expires,
// automatically refreshing it as necessary using the provided context.
//
// Most users will use Config.Client instead.
func (s *keycloak) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {

	return s.oauth2Config.TokenSource(ctx, t)
}

func (s *keycloak) UserInfo(ctx context.Context, ts oauth2.TokenSource) (map[string]interface{}, error) {

	userinfo, err := s.provider.UserInfo(ctx, ts)
	if err != nil {
		return nil, err
	}
	v := make(map[string]interface{})
	if err := userinfo.Claims(&v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *keycloak) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	return s.oauth2Config.Client(ctx, t)
}

/**/
