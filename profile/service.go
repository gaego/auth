// Copyright 2012 GAEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package profile

import (
	"github.com/gaego/context"
	"github.com/gaego/person"
	"github.com/gaego/user"
	"net/http"
)

type Service struct{}

type Args struct{}

type Reply struct {
	Profiles []*person.Person
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request,
	args *Args, reply *Reply) (err error) {

	c := context.NewContext(r)
	u, err := user.Current(r)
	if err != nil {
		return err
	}
	if reply.Profiles, err = GetPersonMulti(c, u.AuthIDs); err != nil {
		return err
	}
	return nil
}
