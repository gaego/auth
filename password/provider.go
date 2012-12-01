// Copyright 2012 GAEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package auth/password provides a password strategy using bcrypt.

auth/password stategy takes a POST with the following keys:

  Email (required)
  Password.New (required/optional)
  Password.Current (required/optional)
  Name.GivenName
  Name.FamilyName
  * (Any other Person attributes)

Based on the supplied attributes auth/password will do one of three things:

1. To Create new User and log them in. POST:
  - "Email"
  - "Password.New" (present)
  - "Password.Current" (NOT present)
  - + Person attributes, E.g. "Name.GivenName", "Name.FamilyName"

2. To Login User or return error if password does not match. POST:
  - "Email"
  - "Password.Current" (present)
  - "Password.New" (NOT present)

3. To Update Password / Person details. POST:
  - "Email"
  - "Password.New" (present)
  - "Password.Current" (present)
  - + Person attributes, E.g. "Name.GivenName", "Name.FamilyName"

*/
package password

import (
	"github.com/gorilla/schema"
	"errors"
	"github.com/gaego/auth/profile"
	"github.com/gaego/context"
	"github.com/gaego/person"
	"github.com/gaego/user"
	"net/http"
)

var (
	ErrProfileNotFound = errors.New("auth/password: profile not found for email address")
)

// Provider represents the auth.Provider
type Provider struct {
	Name, URL string
}

// New creates a New provider.
func New() *Provider {
	return &Provider{"Password", ""}
}

func decodePerson(r *http.Request) *person.Person {
	// Decode the form data and add the resulting Person type to the Profile.
	p := &person.Person{}
	decoder := schema.NewDecoder()
	decoder.Decode(p, r.Form)
	return p
}

// Authenticate process the request and returns a populated Profile.
// If the Authenticate method can not authenticate the User based on the
// request, an error or a redirect URL wll be return.
func (p *Provider) Authenticate(w http.ResponseWriter, r *http.Request) (
	pf *profile.Profile, url string, err error) {

	p.URL = r.URL.Host
	pf = profile.New(p.Name, p.URL)

	pass := &Password{
		New:     r.FormValue("Password.New"),
		Current: r.FormValue("Password.Current"),
		Email:   r.FormValue("Email"),
	}
	c := context.NewContext(r)
	userID, _ := user.CurrentUserIDByEmail(r, pass.Email)
	pers := decodePerson(r)
	pf, err = authenticate(c, pass, pers, userID)
	return pf, "", err
}
