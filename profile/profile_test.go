// Copyright 2012 The AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package profile

import (
	"appengine/datastore"
	"github.com/scotch/aego/v1/context"
	"github.com/scotch/aego/v1/ds"
	"github.com/scotch/aego/v1/person"
	"testing"
)

func tearDown() {
	context.Close()
}

func TestNewKey(t *testing.T) {
	c := context.NewContext(nil)
	defer tearDown()

	k1 := datastore.NewKey(c, "AuthProfile", "google|12345", 0, nil)
	k2 := newKey(c, "Google", "12345")
	if k1.String() != k2.String() {
		t.Errorf("k2: %q, want %q.", k2, k1)
		t.Errorf("k1:", k1)
		t.Errorf("k2:", k2)
	}
}

func TestGet(t *testing.T) {
	c := context.NewContext(nil)
	defer tearDown()

	// Save it.

	u := New("Google", "http://plus.google.com")
	u.ID = "12345"
	u.Person = &person.Person{
		Name: &person.PersonName{
			GivenName:  "Barack",
			FamilyName: "Obama",
		},
	}
	key := newKey(c, "google", "12345")
	u.Key = key
	err := u.Put(c)
	if err != nil {
		t.Errorf(`err: %q, want nil`, err)
	}

	// Get it.

	u2 := &Profile{}
	id := "google|12345"
	key = datastore.NewKey(c, "AuthProfile", id, 0, nil)
	err = ds.Get(c, key, u2)
	if err != nil {
		t.Errorf(`err: %q, want nil`, err)
	}
	u2, err = Get(c, id)
	if err != nil {
		t.Errorf(`err: %v, want nil`, err)
	}
	if u2.ID != "12345" {
		t.Errorf(`u2.ID: %v, want "1"`, u2.ID)
	}
	if u2.Key.StringID() != "google|12345" {
		t.Errorf(`uKey.StringID(): %v, want "google|12345"`, u2.Key.StringID())
	}
	if x := u2.ProviderURL; x != "http://plus.google.com" {
		t.Errorf(`u2.ProviderURL: %v, want %s`, x, "http://plus.google.com")
	}
	if x := u2.Person.ID; x != "12345" {
		t.Errorf(`u2.Person.ID: %v, want %s`, x, "12345")
	}
	if x := u2.Person.Name.GivenName; x != "Barack" {
		t.Errorf(`u2.Person.Name.GivenName: %v, want %s`, x, "Barack")
	}
	if x := u2.Person.Provider.Name; x != "Google" {
		t.Errorf(`u2.Person.Provider.Name: %v, want %s`, x, "Google")
	}
	if x := u2.Person.Provider.URL; x != "http://plus.google.com" {
		t.Errorf(`u2.Person.Provider.URL: %v, want %s`, x, "http://plus.google.com")
	}
	if x := u2.Person.Kind; x != "google#person" {
		t.Errorf(`u2.Person.Kind: %v, want %s`, x, "google#person")
	}
	if u2.Person.Created != u2.Created.UnixNano()/1000000 {
		t.Errorf(`u2.Created: %v, want %v`, u2.Person.Created,
			u2.Created.UnixNano()/1000000)
	}
	if u2.Person.Updated != u2.Updated.UnixNano()/1000000 {
		t.Errorf(`u2.Updated: %v, want %v`, u2.Person.Updated,
			u2.Updated.UnixNano()/1000000)
	}
}
