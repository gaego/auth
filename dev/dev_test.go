// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dev

import (
	"github.com/scotch/aego/v1/auth"
	"github.com/scotch/aego/v1/context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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
	auth.Register("dev", pro)

	// Post.
	v := url.Values{}
	v.Set("ID", "1")
	v.Set("Gender", "male")
	v.Set("Name.GivenName", "Barack")
	v.Set("Name.FamilyName", "Obama")
	v.Set("AboutMe", "This is a bio about me.")
	body := strings.NewReader(v.Encode())

	req, _ := http.NewRequest("POST", "http://localhost:8080/-/auth/dev", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")

	// TODO(kylefinley) for some reason is if this isn't call here the form will
	// be empty in the Authentication method? Perhaps this is a bug.
	if id := req.FormValue("ID"); id != "1" {
		t.Errorf(`req.FormValue("ID") = %q, want "1"`, id)
	}

	// Process.

	up, url, err := pro.Authenticate(w, req)

	// Check.

	if url != "" {
		t.Errorf(`url: %v, want: ""`, url)
	}
	if err != nil {
		t.Errorf(`err: %v, want: %v`, err, nil)
	}

	per := up.Person

	if x := per.Name.GivenName; x != "Barack" {
		t.Errorf(`per.Name.GivenName: %q, want %v`, x, "Barack")
	}
	if x := per.Name.FamilyName; x != "Obama" {
		t.Errorf(`per.Name.FamilyName: %q, want %v`, x, "Obama")
	}
}
