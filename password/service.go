// Copyright 2012 GAEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package password

import (
	"github.com/gaego/auth"
	"github.com/gaego/auth/profile"
	"github.com/gaego/context"
	"github.com/gaego/person"
	"github.com/gaego/user"
	"net/http"
)

type Service struct{}

type Args struct {
	Password *Password
	Person   *person.Person
}

func (s *Service) Authenticate(w http.ResponseWriter, r *http.Request,
	args *Args, reply *Args) (err error) {

	c := context.NewContext(r)
	args.Person.Email = args.Password.Email
	userID, _ := user.CurrentUserIDByEmail(r, args.Password.Email)
	pf, err := authenticate(c, args.Password, args.Person, userID)
	if err != nil {
		return err
	}
	if _, err = auth.CreateAndLogin(w, r, pf); err != nil {
		return err
	}
	reply.Person = pf.Person
	return nil
}

// Current returns the current users password object minus the password
func (s *Service) Current(w http.ResponseWriter, r *http.Request,
	args *Args, reply *Args) (err error) {

	c := context.NewContext(r)
	var isSet bool
	userID, _ := user.CurrentUserID(r)
	_, err = profile.Get(c, profile.GenAuthID("Password", userID))
	if err == nil {
		isSet = true
	}
	reply.Password = &Password{IsSet: isSet}
	return nil
}
