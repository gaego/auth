// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package auth/dev provides a developement strategy for testing.
*/
package dev

import (
	"code.google.com/p/gorilla/schema"
	"github.com/scotch/aego/v1/auth/profile"
	"github.com/scotch/aego/v1/person"
	"net/http"
)

type Provider struct {
	Name, URL string
}

// New creates a New provider.
func New() *Provider {
	return &Provider{"Dev", "http://localhost:8080"}
}

// Authenticate process the request and returns a populated Profile.
// If the Authenticate method can not authenticate the User based on the
// request, an error or a redirect URL wll be return.
func (p *Provider) Authenticate(w http.ResponseWriter, r *http.Request) (
	up *profile.Profile, url string, err error) {

	up = profile.New(p.Name, p.URL)
	// Add the User's Unique ID. If an ID is not provided make this
	// value "default"
	up.ID = r.FormValue("ID")
	if up.ID == "" {
		up.ID = "default"
	}

	// Decode the form data and add the resulting Person type to the Profile.
	per := &person.Person{}
	decoder := schema.NewDecoder()
	decoder.Decode(per, r.Form)
	up.Person = per

	return up, "", nil
}
