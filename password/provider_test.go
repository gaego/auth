// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package password

import (
	"github.com/scotch/aego/v1/auth"
	"github.com/scotch/aego/v1/auth/profile"
	"github.com/scotch/aego/v1/context"
	"github.com/scotch/aego/v1/user/email"
	"github.com/scotch/aego/v1/person"
	"github.com/scotch/aego/v1/user"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

var startOnce sync.Once

func setup() *Provider {
	// Register.
	pro := New()
	startOnce.Do(func() {
		auth.Register("password", pro)
	})
	return pro
}

func tearDown() {
	context.Close()
}

func createRequest(v url.Values) *http.Request {
	body := strings.NewReader(v.Encode())
	req, _ := http.NewRequest("POST",
		"http://localhost:8080/-/auth/password", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	return req
}

func TestAuthenticate_Validate(t *testing.T) {
	pro := setup()
	defer tearDown()

	var err error
	var v url.Values
	var r *http.Request

	w := httptest.NewRecorder()

	// Check invalid POST
	v = url.Values{}
	v.Set("Email", "fake")
	v.Set("Password.New", "bad")
	r = createRequest(v)
	// Check.
	if _, _, err = pro.Authenticate(w, r); err == nil {
		t.Errorf(`Invalid Email or Password should result in an error`)
	}
}

// Scenario #1:
// - No User session
// - No Email Saved
// - No Profile Saved
func TestAuthenticate_Scenario1(t *testing.T) {
	pro := setup()
	defer tearDown()

	var pf *profile.Profile
	var uRL string
	var err error
	var v url.Values
	var r *http.Request

	w := httptest.NewRecorder()

	// Post.
	v = url.Values{}
	v.Set("Email", "test@example.org")
	v.Set("Password.New", "secret1")
	v.Set("Name.GivenName", "Barack")
	r = createRequest(v)
	// Check.
	if pf, uRL, err = pro.Authenticate(w, r); uRL != "" || err != nil {
		t.Errorf(`url: %v, want: ""`, uRL)
		t.Errorf(`err: %v, want: %v`, err, nil)
	}
	if x := pf.Person.Name.GivenName; x != "Barack" {
		t.Errorf(`pf.Person.Name.GivenName: %v, want %v`, x, "Barack")
	}
	if x := pf.UserID; x == "" {
		t.Errorf(`pf.UserID: %v, want %v`, x, "Some ID")
	}
	if x := pf.ID; x == "" {
		t.Errorf(`pf.ID: %v, want %v`, x, "Some ID")
	}
}

// Scenario #2:
// - No User session
// - Yes Email Saved
// - Yes Profile Saved
func TestAuthenticate_Scenario2(t *testing.T) {
	pro := setup()
	defer tearDown()

	var pf *profile.Profile
	var uRL string
	var err error
	var v url.Values
	var r *http.Request

	c := context.NewContext(nil)
	w := httptest.NewRecorder()

	// Profile Not found
	v = url.Values{}
	v.Set("Email", "test@example.org")
	v.Set("Password.Current", "secret1")
	r = createRequest(v)
	// Check.
	if pf, uRL, err = pro.Authenticate(w, r); uRL != "" || err != ErrProfileNotFound {
		t.Errorf(`url: %v, want: ""`, uRL)
		t.Errorf(`err: %v, want: %v`, err, ErrProfileNotFound)
	}

	// Setup.
	pf = profile.New("Password", "")
	pf.UserID = "1"
	pf.ID = "1"
	passHash, _ := GenerateFromPassword([]byte("secret1"))
	pf.Auth = passHash
	pf.SetKey(c)
	pf.Person = &person.Person{
		Name: &person.PersonName{
			GivenName:  "Barack",
			FamilyName: "Obama",
		},
	}
	_ = pf.Put(c)
	e := email.New()
	e.UserID = "1"
	e.SetKey(c, "test@example.org")
	_ = e.Put(c)

	// 1. Login
	// a. Correct password.
	v = url.Values{}
	v.Set("Email", "test@example.org")
	v.Set("Password.Current", "secret1")
	v.Set("Name.GivenName", "Berry")
	r = createRequest(v)
	// Check.
	if pf, uRL, err = pro.Authenticate(w, r); uRL != "" || err != nil {
		t.Errorf(`url: %v, want: ""`, uRL)
		t.Fatalf(`err: %v, want: %v`, err, nil)
	}
	if x := pf.Person.Name.GivenName; x != "Barack" {
		t.Errorf(`.Person should not be updated on login`)
	}
	// b. In-Correct password.
	v = url.Values{}
	v.Set("Email", "test@example.org")
	v.Set("Password.Current", "fakepass")
	r = createRequest(v)
	// Check.
	if _, _, err = pro.Authenticate(w, r); err != ErrPasswordMismatch {
		t.Errorf(`err: %v, want: %v`, err, ErrPasswordMismatch)
	}
	// 2. Update
	// a. Correct password.
	v = url.Values{}
	v.Set("Email", "test@example.org")
	v.Set("Password.Current", "secret1")
	v.Set("Password.New", "secret2")
	v.Set("Name.GivenName", "Berry")
	r = createRequest(v)
	// Check.
	if pf, uRL, err = pro.Authenticate(w, r); uRL != "" || err != nil {
		t.Errorf(`url: %v, want: ""`, uRL)
		t.Errorf(`err: %v, want: %v`, err, nil)
	}
	if x := pf.Person.Name.GivenName; x != "Berry" {
		t.Errorf(`pf.Person should be updated on update`)
	}
	if x := pf.UserID; x != "1" {
		t.Errorf(`pf.UserID: %v, want %v`, x, "1")
	}
	if err := CompareHashAndPassword(pf.Auth, []byte("secret2")); err != nil {
		t.Errorf(`Password was not changed`)
	}
	// b. In-Correct password.
	v = url.Values{}
	v.Set("Email", "test@example.org")
	v.Set("Password.Current", "fakepass")
	v.Set("Password.New", "hacked")
	v.Set("Name.GivenName", "Bob")
	r = createRequest(v)
	// Check.
	if _, _, err = pro.Authenticate(w, r); err != ErrPasswordMismatch {
		t.Errorf(`err: %v, want: %v`, err, ErrPasswordMismatch)
	}
	// 2. Create - Should login user
	// a. Correct password.
	v = url.Values{}
	v.Set("Email", "test@example.org")
	v.Set("Password.New", "secret1")
	v.Set("Name.GivenName", "Bob1")
	r = createRequest(v)
	// Check.
	if pf, uRL, err = pro.Authenticate(w, r); uRL != "" || err != nil {
		t.Errorf(`url: %v, want: ""`, uRL)
		t.Errorf(`err: %v, want: %v`, err, nil)
	}
	if x := pf.Person.Name.GivenName; x != "Bob1" {
		t.Errorf(`.Person should be updated on update`)
	}
	if x := pf.UserID; x != "1" {
		t.Errorf(`pf.UserID: %v, want %v`, x, "1")
	}
	if err := CompareHashAndPassword(pf.Auth, []byte("secret1")); err != nil {
		t.Errorf(`Password was not changed`)
	}
	// b. In-Correct password.
	v = url.Values{}
	v.Set("Email", "test@example.org")
	v.Set("Password.New", "fakepass")
	v.Set("Name.GivenName", "Bob2")
	r = createRequest(v)
	// Check.
	if _, _, err = pro.Authenticate(w, r); err != ErrPasswordMismatch {
		t.Errorf(`err: %v, want: %v`, err, ErrPasswordMismatch)
	}
}

// Scenario #3:
// - Yes User session
// - No Email Saved
// - No Profile Saved
func TestAuthenticate_Scenario3(t *testing.T) {
	pro := setup()
	defer tearDown()

	var pf *profile.Profile
	var uRL string
	var err error
	var v url.Values
	var r *http.Request

	//c := context.NewContext(nil)
	w := httptest.NewRecorder()

	// Setup.
	v = url.Values{}
	v.Set("Email", "test@example.org")
	v.Set("Password.New", "secret1")
	v.Set("Name.GivenName", "Bob")
	r = createRequest(v)
	_ = user.CurrentUserSetID(w, r, "1001")

	// Check.
	if pf, uRL, err = pro.Authenticate(w, r); uRL != "" || err != nil {
		t.Errorf(`url: %v, want: ""`, uRL)
		t.Errorf(`err: %v, want: %v`, err, nil)
	}
	if x := pf.Person.Name.GivenName; x != "Bob" {
		t.Errorf(`pf.Person should be updated`)
	}
	if x := pf.UserID; x != "1001" {
		t.Errorf(`pf.UserID: %v, want %v`, x, "1001")
	}
	if pf.UserID != "1001" {
		t.Errorf(`pf.UserID: %v, want: %v`, pf.UserID, "1001")
	}
}
