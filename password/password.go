// Copyright 2012 GAEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package password

import (
	"appengine"
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"github.com/gaego/auth/profile"
	"github.com/gaego/person"
	"github.com/gaego/user"
	"github.com/gaego/user/email"
)

var (
	PasswordLengthMin = 4
	PasswordLengthMax = 31
	BryptCost         = 12
)

var (
	ErrPasswordMismatch = errors.New("auth/password: passwords do not match")
	ErrPasswordLength   = errors.New("auth/password: passwords must be between 4 and 31 charaters")
)

type Password struct {
	New     string `json:"new,omitempty"`
	Current string `json:"current,omitempty"`
	IsSet   bool   `json:"isSet"`
	Email   string `json:"email"`
}

// validatePasswordLength returns true if the supplied string is
// between 4 and 31 character.
func Validate(p string) error {
	if len(p) < PasswordLengthMin {
		return ErrPasswordLength
	}
	if len(p) > PasswordLengthMax {
		return ErrPasswordLength
	}
	return nil
}

func (p *Password) Validate() (err error) {
	// Validate pasword
	if p.New != "" {
		if err = Validate(p.New); err != nil {
			return
		}
	}
	if p.Current != "" {
		if err = Validate(p.Current); err != nil {
			return
		}
	}
	// Validate email
	if err = email.Validate(p.Email); err != nil {
		return
	}
	return
}

func GenerateFromPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, BryptCost)
}

func CompareHashAndPassword(hash, password []byte) error {
	if bcrypt.CompareHashAndPassword(hash, password) != nil {
		return ErrPasswordMismatch
	}
	return nil
}

func authenticate(c appengine.Context, pass *Password, pers *person.Person, userID string) (
	pf *profile.Profile, err error) {

	if err = pass.Validate(); err != nil {
		return nil, err
	}
	if pass.New != "" && pass.Current != "" {
		pf, err = update(c, pass.Current, pass.New, userID, pers)
		return
	}
	if pass.New != "" {
		// if we have a user ID check for a profile
		if userID != "" {
			if pf, err = login(c, pass.New, userID); err == ErrProfileNotFound {
				pf, err = create(c, pass.New, pers, userID)
				return
			}
			if err != nil {
				return
			}
		}
		pf, err = create(c, pass.New, pers, "")
		return
	}
	if pass.Current != "" {
		pf, err = login(c, pass.Current, userID)
		return
	}
	return pf, nil
}

func create(c appengine.Context, pass string, pers *person.Person, userID string) (
	pf *profile.Profile, err error) {

	var id string
	if userID == "" {
		u := user.New()
		u.SetKey(c)
		if err = u.Put(c); err != nil {
			return
		}
		id = u.Key.StringID()
	} else {
		id = userID
	}
	pf = profile.New("Password", "")
	pf.ID = id
	pf.UserID = id
	pf.Auth, _ = GenerateFromPassword([]byte(pass))
	pf.Person = pers
	return
}

func login(c appengine.Context, pass string, userID string) (
	pf *profile.Profile, err error) {

	if userID == "" {
		return nil, ErrProfileNotFound
	}
	pid := profile.GenAuthID("Password", userID)
	if pf, err = profile.Get(c, pid); err != nil {
		return nil, ErrProfileNotFound
	}
	if err := CompareHashAndPassword(pf.Auth, []byte(pass)); err != nil {
		return nil, err
	}
	return pf, nil
}

func update(c appengine.Context, passCurrent, passNew string, userID string, pers *person.Person) (
	pf *profile.Profile, err error) {

	if pf, err = login(c, passCurrent, userID); err != nil {
		return
	}
	pf.Auth, _ = GenerateFromPassword([]byte(passNew))
	pf.Person = pers
	return pf, nil
}
