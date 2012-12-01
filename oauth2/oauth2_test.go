// Copyright 2012 GAEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth2

import (
	//"github.com/gaego/auth/provider"
	"net/url"
	"testing"
)

func setup() {

}
func tearDown() {

}

type ExampleProvider struct {
	OAuth2Provider
	Provider
	Config
	ClientID     string
	ClientSecret string
	Scope        string
	AuthURL      string
	TokenURL     string
	RedirectURL  string
}

func TestProvider(clientID, clientSecret, scope string) *ExampleProvider {
	return &ExampleProvider{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
		AuthURL:      "http://example.com/auth",
		TokenURL:     "http://example.com/token",
	}
}

func NewExamplePrvider(clientID, clientSecret, scope string) *ExampleProvider {
	return &ExampleProvider{
		OAuth2Provider.ClientID:     clientID,
		OAuth2Provider.ClientSecret: clientSecret,
		OAuth2Provider.Scope:        scope,
		OAuth2Provider.AuthURL:      "http://example.com/auth",
		OAuth2Provider.TokenURL:     "http://example.com/token",
	}
}

func (p *ExampleProvider) GetProviderName() interface{} {
	return "Good Bye"
}

func TestConfig(t *testing.T) {
	//e := new(ExampleProvider)
	u := &url.URL{
		Host:   "test.com",
		Scheme: "http",
	}
	p := NewProvider("Test", "http://test.com", "123", "abc", "email",
		"http://example.com/auth", "http://example.com/token")
	p.Config(u)

	if x := p.RedirectURL; x == "http://test.com/-/auth/test/callback" {
		t.Errorf(`RedirctURL: %v, want %v`, x, "http://test.com/-/auth/test/callback")
	}
	e := ExampleProvider{
		OAuth2Provider: OAuth2Provider{
			Name:         "Google",
			ClientID:     "12345",
			ClientSecret: "password",
			Scope:        "email",
		},
	}
	e := ExampleProvider{
		Provider: Provider{
			Name: "Google",
		},
		Config: Config{
			ClientID:     "12345",
			ClientSecret: "password",
			Scope:        "email",
		},
	}
	e := &ExampleProvider{Provider: Provider{Name: "Google"}}
	e.OAuth2Provider.ClientSecret = "1234"
	c := e.Config.Get("test.com")
	t.Errorf(` %q"`, e.Name)
	c := e.SayName()
	t.Errorf(` %q"`, c)
	ep := NewExamplePrvider("12345", "password", "email")
	if ep.AuthURL != "http://example.com/auth" {
		t.Errorf(`ep.AuthURL = %q, want "http://example.com/auth"`,
			x.AuthURL)
	}
}

func TestAuthenticate(t *testing.T) {
	setUp()
	defer tearDown()

	w := httptest.NewRecorder()

	// Register.

	pro := New()
	auth.Register("dev", pro)

	// Post.
	v := url.Values{}
	v.Set("ID", "1")
	v.Set("Gender", "male")
	v.Set("Name.GivenName", "Barack")
	v.Set("Name.FamilyName", "Obama")
	v.Set("AboutMe", "This is a bio about me.")
	body := strings.NewReader(v.Encode())

	req, _ := http.NewRequest("POST", "http://localhost:8080/", body)

	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	// TODO(kylefinley) for some reason is if this isn't call here the form will
	// be empty in the Authentication method? Perhaps this is a bug.
	if id := req.FormValue("ID"); id != "1" {
		t.Errorf(`req.FormValue("ID") = %q, want "1"`, id)
	}

	// Process.

	up := profile.New()
	url, err := pro.Authenticate(w, req, up)

	// Check.

	if url != "" {
		t.Errorf(`url: %v, want: ""`, url)
	}
	if err != nil {
		t.Errorf(`err: %v, want: %v`, err, nil)
	}

	per, err := up.Person()

	if err != nil {
		t.Errorf(`err: %v, want: %v`, err, nil)
	}
	if x := per.Name.GivenName; x != "Barack" {
		t.Errorf(`per.Name.GivenName: %q, want %v`, x, "Barack")
	}
	if x := per.Name.FamilyName; x != "Obama" {
		t.Errorf(`per.Name.FamilyName: %q, want %v`, x, "Obama")
	}
}
