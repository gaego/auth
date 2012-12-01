// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"appengine/datastore"
	"errors"
	"github.com/scotch/aego/v1/auth/dev"
	"github.com/scotch/aego/v1/auth/profile"
	"github.com/scotch/aego/v1/context"
	"github.com/scotch/aego/v1/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() {
	BaseURL = "/-/auth/"
	LoginURL = "/-/auth/login"
	LogoutURL = "/-/auth/logout"
	SuccessURL = "/"
}

func teardown() {
	context.Close()
}

type TestProvider struct {
	dev.Provider
}

func (p *TestProvider) Authenticate(w http.ResponseWriter, r *http.Request) (
	up *profile.Profile, url string, err error) {
	return nil, "/redirect-to-url", nil
}

type TPRedirect struct {
	dev.Provider
}

func (p *TPRedirect) Authenticate(w http.ResponseWriter, r *http.Request) (
	up *profile.Profile, url string, err error) {
	return nil, "/redirect-to-url", nil
}

type TPError struct {
	dev.Provider
}

func (p *TPError) Authenticate(w http.ResponseWriter, r *http.Request) (
	up *profile.Profile, url string, err error) {
	err = errors.New("Mock error")
	return nil, "", err
}

type TPComplete struct {
	dev.Provider
}

func (p *TPComplete) Authenticate(w http.ResponseWriter, r *http.Request) (
	up *profile.Profile, url string, err error) {
	up = profile.New("Example", "example.com")
	up.ID = "1"
	return up, "", nil
}

func TestNew(t *testing.T) {
	setup()
	defer teardown()

	x := &TestProvider{}
	x.Name = "Example"
	x.URL = "http://exaple.com"
	// Confirm that it implements authenticater
	var y interface{} = x
	p, ok := y.(authenticater)
	if !ok {
		t.Errorf(`p = %q,"`, p)
	}
}

func TestBreakURL(t *testing.T) {
	setup()
	defer teardown()

	url1 := "http://localhost:8080/-/auth/example1"
	n := breakURL(url1)
	if n != "example1" {
		t.Errorf(`n: %q, want example1`, n)
	}
	url2 := "http://localhost:8080/-/auth/example2/callback?some=crazy[stuff]"
	n2 := breakURL(url2)
	if n2 != "example2" {
		t.Errorf(`n: %q, want example2`, n2)
	}
	// Change the BaseURL
	BaseURL = "/changed/"
	url3 := "http://localhost:8080/changed/example3/callback?some=crazy[stuff]"
	n3 := breakURL(url3)
	if n3 != "example3" {
		t.Errorf(`n: %q, want example3`, n3)
	}
}

func TestRedirect(t *testing.T) {
	setup()
	defer teardown()

	// Register

	p := &TPRedirect{}
	Register("example2", p)
	r, _ := http.NewRequest("GET", "http://localhost:8080/-/auth/example2", nil)
	w := httptest.NewRecorder()

	// Run it through the auth handler.

	handler(w, r)

	// Inspected the redirect.

	hdr := w.Header()
	if hdr["Location"][0] != "/redirect-to-url" {
		t.Errorf(`hdr["Location"]: %q, want "/redirect-to-url"`, hdr["Location"])
	}
}

func TestError(t *testing.T) {
	setup()
	defer teardown()

	// Register

	p := &TPError{}
	Register("example3", p)
	r, _ := http.NewRequest("GET", "http://localhost:8080/-/auth/example3", nil)
	w := httptest.NewRecorder()

	// Run it through the auth handler.

	handler(w, r)

	// Inspected the redirect.

	hdr := w.Header()
	if hdr["Location"][0] != LoginURL {
		t.Errorf(`hdr["Location"]: %q, want %q`, hdr["Location"], LoginURL)
	}
}

func Test_handler(t *testing.T) {
	setup()
	defer teardown()
	_ = context.NewContext(nil)

	// Register the Provider

	p := &TPComplete{}
	Register("example5", p)
	r, _ := http.NewRequest("GET", "http://localhost:8080/-/auth/example5", nil)
	w := httptest.NewRecorder()

	// Run it through the auth handler.

	handler(w, r)

	// Inspected the redirect.

	hdr := w.Header()
	if hdr["Location"][0] != SuccessURL {
		t.Errorf(`hdr["Location"]: %q, want %q`, hdr["Location"][0], SuccessURL)
		t.Errorf(`w: %q`, w)
		t.Errorf(`hdr: %q`, hdr)
	}
}

func Test_CreateAndLogin(t *testing.T) {
	setup()
	defer teardown()
	c := context.NewContext(nil)

	up := profile.New("Example", "example.com")
	r, _ := http.NewRequest("GET", "http://localhost:8080/-/auth/example4", nil)
	w := httptest.NewRecorder()

	// Round 1: No User | No Profile

	// Confirm.

	q := datastore.NewQuery("User")
	if cnt, _ := q.Count(c); cnt != 0 {
		t.Errorf(`User cnt: %v, want 0`, cnt)
	}
	q = datastore.NewQuery("Profile")
	if cnt, _ := q.Count(c); cnt != 0 {
		t.Errorf(`Profile cnt: %v, want 0`, cnt)
	}
	u, err := user.Current(r)
	if err != user.ErrNoLoggedInUser {
		t.Errorf(`err: %v, want %v`, err, user.ErrNoLoggedInUser)
	}

	// Create.

	up.ID = "1"
	up.ProviderName = "Example"
	up.SetKey(c)
	u, err = CreateAndLogin(w, r, up)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}

	if u.Key.StringID() != "1" {
		t.Errorf(`u.Key.StringID(): %v, want 1`, u.Key.StringID())
	}
	if up.Key.StringID() != "example|1" {
		t.Errorf(`up.Key.StringID(): %v, want "example|1"`, up.Key.StringID())
	}
	if up.UserID != u.Key.StringID() {
		t.Errorf(`up.UserID: %v, want %v`, up.UserID, u.Key.StringID())
	}

	// Confirm Profile.

	rup, err := profile.Get(c, "example|1")
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if rup.ID != "1" {
		t.Errorf(`rup.ID: %v, want "1"`, rup.ID)
	}
	if rup.Key.StringID() != "example|1" {
		t.Errorf(`rup.Key.StringID(): %v, want "example|1"`, rup.Key.StringID())
	}
	if rup.UserID != u.Key.StringID() {
		t.Errorf(`rup.UserID: %v, want %v`, rup.UserID, u.Key.StringID())
	}

	// Confirm User.

	ru, err := user.Get(c, "1")
	if err != nil {
		t.Fatalf(`err: %v, want nil`, err)
	}
	if ru.AuthIDs[0] != "example|1" {
		t.Errorf(`ru.AuthIDs[0]: %v, want "example|1"`, ru.AuthIDs[0])
	}
	if ru.Key.StringID() != "1" {
		t.Errorf(`ru.Key.StringID(): %v, want 1`, ru.Key.StringID())
	}
	q2 := datastore.NewQuery("User")
	q4 := datastore.NewQuery("AuthProfile")

	// Confirm Logged in User.

	u, err = user.Current(r)
	if err != nil {
		t.Errorf(`err: %v, want %v`, err, nil)
	}
	if u.Key.StringID() != "1" {
		t.Errorf(`u.Key.StringID(): %v, want 1`, u.Key.StringID())
	}
	if len(u.AuthIDs) != 1 {
		t.Errorf(`len(u.AuthIDs): %v, want 1`, len(u.AuthIDs))
		t.Errorf(`u.AuthIDs: %v`, u.AuthIDs)
		t.Errorf(`u: %v`, u)
	}

	// Round 2: Logged in User | Second Profile

	// Create.

	up = profile.New("AnotherExample", "anotherexample.com")
	up.ID = "2"
	up.SetKey(c)
	u, err = CreateAndLogin(w, r, up)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}

	// Confirm Profile.

	rup, err = profile.Get(c, "anotherexample|2")
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if rup.ID != "2" {
		t.Errorf(`rup.ID: %v, want "2"`, rup.ID)
	}
	if rup.Key.StringID() != "anotherexample|2" {
		t.Errorf(`rup.Key.StringID(): %v, want "anotherexample|2"`, rup.Key.StringID())
	}
	if rup.UserID != u.Key.StringID() {
		t.Errorf(`rup.UserID: %v, want %v`, rup.UserID, u.Key.StringID())
	}

	// Confirm Logged in User hasn't changed.

	u, err = user.Current(r)
	if err != nil {
		t.Errorf(`err: %v, want %v`, err, nil)
	}
	if u.Key.StringID() != "1" {
		t.Errorf(`u.Key.StringID(): %v, want 1`, u.Key.StringID())
	}
	if len(u.AuthIDs) != 2 {
		t.Errorf(`len(u.AuthIDs): %v, want 2`, len(u.AuthIDs))
		t.Errorf(`u.AuthIDs: %v`, u.AuthIDs)
		t.Errorf(`u: %v`, u)
	}
	if u.AuthIDs[0] != "example|1" {
		t.Errorf(`u.AuthIDs[0]: %v, want "example|1"`, u.AuthIDs[0])
	}
	if u.AuthIDs[1] != "anotherexample|2" {
		t.Errorf(`u.AuthIDs[1]: %v, want "anotherexample|2"`, u.AuthIDs[1])
	}

	// Confirm Counts

	q2 = datastore.NewQuery("User")
	if cnt, _ := q2.Count(c); cnt != 1 {
		t.Errorf(`User cnt: %v, want 1`, cnt)
	}
	q4 = datastore.NewQuery("AuthProfile")

	// Round 3: Logged out User | Existing Profile

	err = user.Logout(w, r)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}

	// Confirm Logged out User.

	u, err = user.Current(r)
	if err != user.ErrNoLoggedInUser {
		t.Errorf(`err: %q, want %q`, err, user.ErrNoLoggedInUser)
	}

	// Login.

	up = profile.New("Example", "example.com")
	up.ID = "1"
	up.SetKey(c)
	u, err = CreateAndLogin(w, r, up)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}

	// Confirm.

	q2 = datastore.NewQuery("User")
	if cnt, _ := q2.Count(c); cnt != 1 {
		t.Errorf(`User cnt: %v, want 1`, cnt)
	}
	q4 = datastore.NewQuery("AuthProfile")
	if cnt, _ := q4.Count(c); cnt != 2 {
		t.Errorf(`Profile cnt: %v, want 1`, cnt)
	}

	// Confirm Logged in User hasn't changed.

	u, err = user.Current(r)
	if err != nil {
		t.Errorf(`err: %v, want %v`, err, nil)
	}
	if u.Key.StringID() != "1" {
		t.Errorf(`u.Key.StringID(): %v, want "1"`, u.Key.StringID())
	}
	if len(u.AuthIDs) != 2 {
		t.Errorf(`len(u.AuthIDs): %v, want 2`, len(u.AuthIDs))
		t.Errorf(`u.AuthIDs: %s`, u.AuthIDs)
		t.Errorf(`u: %v`, u)
	}
}
