/*
 *  Licensed to Wikifeat under one or more contributor license agreements.
 *  See the LICENSE.txt file distributed with this work for additional information
 *  regarding copyright ownership.
 *
 *  Redistribution and use in source and binary forms, with or without
 *  modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *  this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright
 *  notice, this list of conditions and the following disclaimer in the
 *  documentation and/or other materials provided with the distribution.
 *  * Neither the name of Wikifeat nor the names of its contributors may be used
 *  to endorse or promote products derived from this software without
 *  specific prior written permission.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */
package auth_service

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	. "github.com/rhinoman/wikifeat/common/auth"
	"github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/common/util"
	"net/http"
)

/**
 * Standard Authenticator - we authenticate against CouchDB users directly
 */
type StandardAuthenticator struct{}

// Very simple for the standard auth, the AuthSession cookie value *is* the session id
func (sta StandardAuthenticator) GetSessionId(req *http.Request) (string, error) {
	//Ok, check for a session cookie
	sessCookie, err := req.Cookie("AuthSession")
	if err != nil {
		return "", UnauthenticatedError()
	}
	sessionToken := sessCookie.Value
	//csrfErr := sta.checkCsrf(req)
	if sessionToken == "" {
		//Bad user, no cookie
		return "", UnauthenticatedError()
	}
	//Return the sessionId
	return sessionToken, nil
}

func (sta StandardAuthenticator) SetAuth(rw http.ResponseWriter, cAuth couchdb.Auth) {
	authData := cAuth.GetUpdatedAuth()
	if authData == nil {
		return
	}
	if val, ok := authData["AuthSession"]; ok {
		authCookie := http.Cookie{
			Name:     "AuthSession",
			Value:    val,
			Path:     "/",
			HttpOnly: true,
		}
		//Create a CSRF cookie
		csrfCookie := http.Cookie{
			Name:     "CsrfToken",
			Value:    util.GenHashString(val),
			Path:     "/",
			HttpOnly: false,
		}
		rw.Header().Add("Set-Cookie", authCookie.String())
		rw.Header().Add("Set-Cookie", csrfCookie.String())
	}
}

// Create a new session by validating against a user's couchdb credentials
func (sta StandardAuthenticator) CreateSession(username, password string) (*Session, error) {
	//Let's validate the user's credentials against CouchDB
	ba := &couchdb.BasicAuth{
		Username: username,
		Password: password,
	}
	authInfo, err := services.Connection.GetAuthInfo(ba)
	if err != nil || authInfo == nil {
		return nil, UnauthenticatedError()
	} else if !authInfo.Ok || authInfo.UserCtx.Name == "" {
		return nil, UnauthenticatedError()
	}
	return NewSession(&authInfo.UserCtx, "standard"), nil
}

// The StandardAuthenticator need do nothing here.
func (sta StandardAuthenticator) DestroySession(sessionId string) error {
	return nil
}

func (sta StandardAuthenticator) checkCsrf(request *http.Request) error {
	//Check CSRF token
	csrfCookie, err := request.Cookie("CsrfToken")
	ourCsrf := request.Header.Get("X-Csrf-Token")
	if err != nil || ourCsrf == "" {
		return UnauthenticatedError()
	} else if token := csrfCookie.Value; token != ourCsrf {
		return UnauthenticatedError()
	} else {
		return nil
	}
}
