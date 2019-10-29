package keycloak

import (
	"context"
	"flag"
	_ "fmt"
	_ "net/url"
	"testing"
	"time"

	"github.com/dumacp/utils"
)

var clientId string
var urlserver string
var realm string
var clientID string
var username string
var password string
var clientSecret string

func init() {
	flag.StringVar(&clientSecret, "clientSecret", "0a955e51-7263-4c83-8900-bcb652258f62", "client ID in realms oauth2")
	flag.StringVar(&clientID, "clientID", "iot-devices", "client ID in realms oauth2")
	flag.StringVar(&username, "username", "b-z8-0001", "username in realms oauth2")
	flag.StringVar(&password, "password", "b-z8-0001", "password in realms oauth2")
	flag.StringVar(&urlserver, "url", "http://localhost:8080/auth", "url openid server")
	flag.StringVar(&realm, "realm", "master", "realm in oauth2")
	flag.StringVar(&clientId, "clientId", "clientId", "client in realm")
}

func TestRquest(t *testing.T) {

	ctx := context.Background()
	// Configure an OpenID Connect aware OAuth2 client.
	serverConfig := &ServerConfig{Url: urlserver,
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectUrl:  "",
	}
	server, err := NewConfig(ctx, serverConfig)
	if err != nil {
		t.Fatal(err)
	}

	token, err := server.TokenRequest(ctx, username, password)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Token: %+v: ", token)
	t.Logf("Token: %v: ", utils.PrettyPrint(token))

	for x := 0; x < 2; x++ {
		time.Sleep(3 * time.Minute)
		tokenSource := server.TokenSource(ctx, token)
		if err != nil {
			t.Fatal(err)
		}
		t1, err := tokenSource.Token()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("TokenSource: %v: ", utils.PrettyPrint(t1))

		userInfo, err := server.UserInfo(ctx, tokenSource)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("UserInfo: %+v: ", userInfo)
		t.Logf("UserInfo: %v: ", utils.PrettyPrint(userInfo))
	}

	/**
	client := server.Client(ctx, token)
	t.Logf("URL: %s/introspect", server.Oauth2Config.Endpoint.TokenURL)
	data := make(url.Values)
	data.Add("client_secret", server.Config.ClientSecret)
	data.Add("client_id", server.Config.ClientID)
	data.Add("username", username)
	data.Add("token", token.AccessToken)

	resp, err := client.PostForm(fmt.Sprintf("%s/introspect", server.Oauth2Config.Endpoint.TokenURL), data)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	ch := make(chan []byte, 0)
	go func() {
		defer close(ch)
		for {
			buffer := make([]byte, 128)
			n, err := resp.Body.Read(buffer)
			if n <= 0 {
				t.Log("no hay más datos en el Body")
				return
			}
			if err != nil {
				t.Log(err)
				return
			}
			ch <- buffer[0:n]
		}
	}()

	readAll := make([]byte, 0)
	for v := range ch {
		readAll = append(readAll, v...)
	}
	t.Logf("Header: %s: ", resp.Header)
	t.Logf("UserInfo: %s: ", readAll)
	**/
}
