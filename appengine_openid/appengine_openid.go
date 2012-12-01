// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package auth/appengine_openid provides authentication through 
Google App Engine OpenID
*/
package appengine_openid

import (
	aeuser "appengine/user"
	"github.com/scotch/aego/v1/auth/profile"
	"github.com/scotch/aego/v1/context"
	"github.com/scotch/aego/v1/person"
	"net/http"
)

type Provider struct {
	Name, URL string
}

// New creates a New provider.
func New() *Provider {
	return &Provider{
		"AppEngineOpenID",
		"appengine.google.com",
	}
}

// Authenticate process the request and returns a populated UserProfile.
// If the Authenticate method can not authenticate the User based on the
// request, an error or a redirect URL wll be return.
func (p *Provider) Authenticate(w http.ResponseWriter, r *http.Request) (
	up *profile.Profile, redirectURL string, err error) {

	c := context.NewContext(r)

	url := r.FormValue("provider")
	// Set provider info.
	up = profile.New(p.Name, url)

	// Check for current User.

	u := aeuser.Current(c)

	if u == nil {
		redirectURL := r.URL.Path + "/callback"
		loginUrl, err := aeuser.LoginURLFederated(c, redirectURL, url)
		return up, loginUrl, err
	}

	if u.FederatedIdentity != "" {
		up.ID = u.FederatedIdentity
	} else {
		up.ID = u.ID
	}

	per := new(person.Person)
	per.Email = u.Email
	per.Emails = []*person.PersonEmails{
		&person.PersonEmails{true, "home", u.Email},
	}
	per.URL = u.FederatedIdentity
	up.Person = per

	return up, "", nil
}
