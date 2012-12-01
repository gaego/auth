// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appengine_openid

import (
	"github.com/scotch/aego/v1/auth"
	"github.com/scotch/aego/v1/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setUp() {}

func tearDown() {
	context.Close()
}

func TestAuthenticate(t *testing.T) {
	setUp()
	defer tearDown()

	w := httptest.NewRecorder()

	// Register.

	pro := New()
	auth.Register("appengine_openid", pro)

	// Round 1: Now User.

	req, _ := http.NewRequest("GET",
		"http://localhost:8080/-/auth/appengine_openid?provider=gmail.com", nil)

	// Process.

	_, url, err := pro.Authenticate(w, req)

	if url == "" {
		exampleURL :=
			"/_ah/login?continue=http%3A//127.0.0.1%3A51002/-/auth/appengine_openid/callback"
		t.Errorf(`url: %v, want: %v`, url, exampleURL)
	}
	if err != nil {
		t.Errorf(`err: %v, want: %v`, err, nil)
	}
	// TODO: appenginetesting does not allow headers to passed to the
	// request. This will have to go non tested for the time being.

	// 	// Round 2: Mock User.
	// 
	// 	req, _ = http.NewRequest("GET",
	// 		"http://localhost:8080/-/auth/appengine_openid/callback", nil)
	// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	// 
	// 	req.Header.Set("X-AppEngine-Inbound-User-Email", "test@example.org")
	// 	req.Header.Set("X-AppEngine-Inbound-User-Federated-Identity", "gmail.com")
	// 	req.Header.Set("X-AppEngine-Inbound-User-Federated-Provider", "google")
	// 	req.Header.Set("X-AppEngine-Inbound-User-Id", "12345")
	// 	req.Header.Set("X-AppEngine-Inbound-User-Is-Admin", "0")
	// 
	// 	// Process.
	// 
	// 	up = user_profile.New()
	// 	url, err = pro.Authenticate(w, req, up)
	// 
	// 	// Check.
	// 
	// 	t.Fatalf(`up.Person: %v`, up.Person)
	// 
	// 	if x := up.ProviderURL; x != "gmail.com" {
	// 		t.Errorf(`ProviderURL: %q, want %v`, x, "gmail.com")
	// 	}
	// 	if x := up.Person.Emails[0].Value; x != "test@example.org" {
	// 		t.Errorf(`Email.Value: %v, want %v`, x, "test@example.org")
	// 	}
}
