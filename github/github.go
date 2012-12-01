// Copyright 2012 GAEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package auth/github provides Github authentication
*/
package github

import (
	"github.com/gaego/auth/oauth2"
)

const (
	CLIENT_ID     = "dbac99a147b10e6bc813"
	CLIENT_SECRET = "5f6e11429eeef14d0fe79721ee53459963e306f5"
)

type Provider struct {
	oauth2.Provider
}

func New(clientID, clientSecret, scope string) *Provider {
	return &Provider{
		Provider: oauth2.Provider{
			Name:         "Github",
			URL:          "http://github.com",
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scope:        "",
			AuthURL:      "https://github.com/login/oauth/authorize",
			TokenURL:     "https://github.com/login/oauth/access_token",
		},
	}
}
