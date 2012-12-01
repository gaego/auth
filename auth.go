// Copyright 2012 GAEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package auth provides multi-provider Authentication.

Example Usage:

  import (
    "github.com/gaego/auth"
    "github.com/gaego/auth/google"
  )

  // Register the Google Provider.
  googleProvider := google.Provider.New("12345", "ABCD")
  Register("google", &googleProvider)
  // Register additional providers.
  // ...


*/
package auth

import (
	"github.com/gaego/auth/profile"
	"github.com/gaego/context"
	"github.com/gaego/user"
	"net/http"
	"strings"
)

var (
	// BaseURL represents the base url to be used for providers. For
	// example if the base url is /auth/ all provider urls would be at
	// /auth/<provider name>
	BaseURL = "/-/auth/"
	// LoginURL is a string representing the URL to be redirected to on
	// errors.
	LoginURL = "/-/auth/login"
	// LogoutURL is a string representing the URL to be used to remove
	// the auth cookie.
	LogoutURL = "/-/auth/logout"
	// SuccessURL is a string representing the URL to be direct to on a
	// successful login.
	SuccessURL = "/"
)

var providers = make(map[string]authenticater)

type authenticater interface {
	Authenticate(http.ResponseWriter, *http.Request) (*profile.Profile, string, error)
}

// Register adds an Authenticater for the auth service.
//
// It takes a string which is used for the url, and a pointer to an
// authentication provider that implements Authenticater.
// E.g.
//
//   googleProvider := google.Provider.New("12345", "ABCD")
//   Register("google", &googleProvider)
//
func Register(key string, auth authenticater) {
	providers[key] = auth
	// Set the start url e.g. /-/auth/google to be handled by the handler.
	http.HandleFunc(BaseURL+key, handler)
	// Set the callback url e.g. /-/auth/google/callback to be handled by the handler.
	http.HandleFunc(BaseURL+key+"/callback", handler)
}

// breakURL parse an url and returns the provider key. If the URL is
// invalid it returns and empty string "".
func breakURL(url string) (name string) {
	if p := strings.Split(url, BaseURL); len(p) > 1 {
		name = strings.Split(p[1], "/")[0]
	}
	return
}

// CreateAndLogin does the following:
//
//  - Search for an existing user - session -> Profile -> email address
//  - Saves the Profile to the datastore
//  - Creates a User or appends the AuthID to the Requesting user's account
//  - Logs in the User
//  - Adds the admin role to the User if they are an GAE Admin.
func CreateAndLogin(w http.ResponseWriter, r *http.Request,
	p *profile.Profile) (u *user.User, err error) {
	c := context.NewContext(r)
	if u, err = p.UpdateUser(w, r); err != nil {
		return
	}
	if err = user.CurrentUserSetID(w, r, p.UserID); err != nil {
		return
	}
	err = p.Put(c)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	var url string
	var err error
	var up *profile.Profile
	k := breakURL(r.URL.Path)
	p := providers[k]
	if up, url, err = p.Authenticate(w, r); err != nil {
		// TODO: set error message in session.
		http.Redirect(w, r, LoginURL, http.StatusFound)
		return
	}
	// If we have a url the Provider wants to make a redirect before
	// proceeding.
	if url != "" {
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	// If we don't have a URL or an error then the user has been authenticated.
	// Check the Profile for an ID and Provider.
	if up.ID == "" || up.ProviderName == "" {
		panic(`auth: The Profile's "ID" or "ProviderName" is empty.` +
			`A Key can not be created.`)
	}
	if _, err = CreateAndLogin(w, r, up); err != nil {
		// TODO: set error message in session.
		http.Redirect(w, r, LoginURL, http.StatusFound)
		return
	}
	// If we've made it this far redirect to the SuccessURL
	http.Redirect(w, r, SuccessURL, http.StatusFound)
	return
}
