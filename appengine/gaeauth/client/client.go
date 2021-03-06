// Copyright 2015 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package client implements OAuth2 authentication for outbound connections
// from Appengine using the application services account.
package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/luci/gae/service/info"
	"github.com/luci/gae/service/urlfetch"

	"github.com/luci/luci-go/common/auth"
	"github.com/luci/luci-go/common/clock"
	"github.com/luci/luci-go/common/data/caching/proccache"
	"github.com/luci/luci-go/common/data/rand/mathrand"
	"github.com/luci/luci-go/common/data/stringset"
	"github.com/luci/luci-go/common/errors"
	"github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/common/transport"
	serverauth "github.com/luci/luci-go/server/auth"
)

// GetAccessToken returns an OAuth access token representing app's service
// account.
//
// If scopes is empty, uses auth.OAuthScopeEmail scope.
//
// Implements a caching layer on top of GAE's GetAccessToken RPC. May return
// transient errors.
func GetAccessToken(c context.Context, scopes []string) (auth.Token, error) {
	scopes, cacheKey := normalizeScopes(scopes)

	// Try to find the token in the local memory first. If it expires soon,
	// refresh it earlier with some probability. That avoids a situation when
	// parallel requests that use access tokens suddenly see the cache expired
	// and rush to refresh the token all at once.
	pcache := proccache.GetCache(c)
	if entry := pcache.Get(cacheKey); entry != nil && !closeToExpRandomized(c, entry.Exp) {
		return entry.Value.(auth.Token), nil
	}

	// The token needs to be refreshed.
	logging.Debugf(c, "Getting an access token for scopes %q", strings.Join(scopes, ", "))
	accessToken, exp, err := info.Get(c).AccessToken(scopes...)
	if err != nil {
		return auth.Token{}, errors.WrapTransient(err)
	}

	// Prematurely expire it to guarantee all returned token live for at least
	// 'expirationMinLifetime'.
	tok := auth.Token{
		AccessToken: accessToken,
		Expiry:      exp.Add(-expirationMinLifetime),
		TokenType:   "Bearer",
	}

	// Store the new token in the cache (overriding what's already there).
	pcache.Put(cacheKey, tok, tok.Expiry)

	return tok, nil
}

// UseServiceAccountTransport injects authenticating transport into
// context.Context. It can be extracted back via transport.Get(c).
//
// If scopes is empty, uses auth.OAuthScopeEmail scope.
//
// TODO(vadimsh): Get rid of this in favor of auth.GetRPCTransport.
func UseServiceAccountTransport(c context.Context, scopes []string) context.Context {
	return transport.SetFactory(c, func(ic context.Context) http.RoundTripper {
		t, err := serverauth.GetRPCTransport(ic, serverauth.AsSelf, serverauth.WithScopes(scopes...))
		if err != nil {
			return failTransport{err}
		}
		return t
	})
}

// UseAnonymousTransport injects non-authenticating GAE transport into
// context.Context. It can be extracted back via transport.Get(c). Use it with
// libraries that search for transport in the context (e.g. common/config),
// since by default they revert to http.DefaultTransport that doesn't work in
// GAE environment.
//
// TODO(vadimsh): Get rid of this in favor of auth.GetRPCTransport.
func UseAnonymousTransport(c context.Context) context.Context {
	return transport.SetFactory(c, func(ic context.Context) http.RoundTripper {
		return urlfetch.Get(ic)
	})
}

//// Internal stuff.

type cacheKey string

const (
	// expirationMinLifetime is minimal possible lifetime of a returned token.
	expirationMinLifetime = 2 * time.Minute
	// expirationRandomization defines how much to randomize expiration time.
	expirationRandomization = 10 * time.Minute
)

func normalizeScopes(scopes []string) ([]string, cacheKey) {
	if len(scopes) == 0 {
		scopes = []string{auth.OAuthScopeEmail}
	} else {
		set := stringset.New(len(scopes))
		for _, s := range scopes {
			if strings.ContainsRune(s, '\n') {
				panic(fmt.Errorf("invalid scope %q", s))
			}
			set.Add(s)
		}
		scopes = set.ToSlice()
		sort.Strings(scopes)
	}
	return scopes, cacheKey(strings.Join(scopes, "\n"))
}

func closeToExpRandomized(c context.Context, exp time.Time) bool {
	switch now := clock.Now(c); {
	case now.After(exp):
		return true // expired already
	case now.Add(expirationRandomization).Before(exp):
		return false // far from expiration
	default:
		// The expiration is close enough. Do the randomization.
		rnd := time.Duration(mathrand.Get(c).Int63n(int64(expirationRandomization)))
		return now.Add(rnd).After(exp)
	}
}

type failTransport struct {
	err error
}

func (f failTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	// http.RoundTripper contract states: "RoundTrip should not modify the
	// request, except for consuming and closing the Body, including on errors"
	if r.Body != nil {
		_, _ = io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}
	return nil, f.err
}
