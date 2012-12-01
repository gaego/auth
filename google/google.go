// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright: 2011 Google Inc. All Rights Reserved.
// license: Apache Software License, see LICENSE for details.

/*
Package auth/google provides Google authentication
*/
package google

import (
	"github.com/scotch/aego/v1/auth/oauth2"
)

type Provider struct {
	oauth2.Provider
}

func New(clientID, clientSecret, scope string) *Provider {
	return &Provider{
		Provider: oauth2.Provider{
			Name:         "Google",
			URL:          "https://plus.google.com",
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scope:        scope,
			AuthURL:      "https://accounts.google.com/o/oauth2/auth",
			TokenURL:     "https://accounts.google.com/o/oauth2/token",
		},
	}
}

// func init() {
// 	defaultCnfg = map[string]string{
// 		"BaseURL":    baseURL,
// 		"LoginURL":   loginURL,
// 		"LogoutURL":  loginURL,
// 		"SuccessURL": successURL,
// 	}
// }
// 
// // setConfig retrieves the global config for the "auth" key and sets
// // local variables based on the response.
// func setConfig(c appengine.Context) {
// 	// TODO(kylefinley): This doesn't work. The cnft needs to be set before
// 	// a context/request is available.
// 	cnfg, err := config.GetOrInsert(c, "auth", defaultCnfg)
// 	if err != nil {
// 		panic("hal/auth: an error occured while setting the config")
// 	}
// 	baseURL = cnfg.Values["BaseURL"]
// 	loginURL = cnfg.Values["LoginURL"]
// 	logoutURL = cnfg.Values["LogoutURL"]
// 	successURL = cnfg.Values["SuccessURL"]
// }

// func (p *UserProfile) PersonRaw(c appengine.Context) interface{} {
// 
// 	// There's a bug where Google Plus doesn"t return an email address.
// 	// So we'll retrieve it the old way and inject it into res.
// 	// We're also checking to se if this account is a legacy account,
// 	// in which case we"ll perform the legacy user lookup.
//   is_legacy = False
//   res = {}
//   try:
//   res = self.service().people().get(userId="me").execute(self.http())
//   except Exception, e:
//   is_legacy = True
//   if is_legacy or "emails" not in res:
//   service = self.service(name="oauth2", version="v1")
//   legacy_res = service.userinfo().get().execute(self.http())
// 
//   email = {
//   "value": legacy_res.get("email"),
//   "primary": True,
//   "verified": legacy_res.get("verified_email")}
//   res["emails"] = [email]
// 
//   if "displayName" not in res:
//   res["displayName"] = legacy_res.get("name")
// 
//   if "name" not in res:
//   res["name"] = {
//   "givenName": legacy_res.get("given_name"),
//   "familyName": legacy_res.get("family_name"),
//   }
// 
//   if "url" not in res:
//   res["url"] = legacy_res.get("link")
// 
//   if "image" not in res:
//   res["image"] = {"url": legacy_res.get("picture")}
// 
//   if "locale" not in res:
//   res["locale"] = legacy_res.get("locale")
//   return res
// 
// }
