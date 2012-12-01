// Copyright 2012 GAEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package auth/facebook provides Facebook authentication
*/
package facebook

import (
	"github.com/gaego/auth/oauth2"
)

const (
	CLIENT_ID     = "343417275669983"
	CLIENT_SECRET = "fec59504f33b238a5d7b5f3b35bd958a"
	PROFILE_URL   = "https://graph.facebook.com/me"
)

type Provider struct {
	oauth2.Provider
}

func New(clientID, clientSecret, scope string) *Provider {
	return &Provider{
		Provider: oauth2.Provider{
			Name:         "Facebook",
			URL:          "http://facebook.com",
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scope:        scope,
			AuthURL:      "https://graph.facebook.com/oauth/authorize",
			TokenURL:     "https://graph.facebook.com/oauth/access_token",
		},
	}
}
