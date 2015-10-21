// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package auth

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/luci/luci-go/server/auth/identity"
	. "github.com/smartystreets/goconvey/convey"
)

func TestState(t *testing.T) {
	Convey("Check empty ctx", t, func() {
		ctx := context.Background()
		So(GetState(ctx), ShouldBeNil)
		So(CurrentUser(ctx).Identity, ShouldEqual, identity.AnonymousIdentity)
		So(CurrentIdentity(ctx), ShouldEqual, identity.AnonymousIdentity)
	})

	Convey("Check non-empty ctx", t, func() {
		s := state{
			user:      &User{Identity: "user:abc@example.com"},
			peerIdent: "user:abc@example.com",
		}
		ctx := context.WithValue(context.Background(), stateContextKey(0), &s)
		So(GetState(ctx), ShouldNotBeNil)
		So(GetState(ctx).Method(), ShouldBeNil)
		So(GetState(ctx).PeerIdentity(), ShouldEqual, identity.Identity("user:abc@example.com"))
		So(GetState(ctx).PeerIP(), ShouldBeNil)
		So(CurrentUser(ctx).Identity, ShouldEqual, identity.Identity("user:abc@example.com"))
		So(CurrentIdentity(ctx), ShouldEqual, identity.Identity("user:abc@example.com"))
	})
}
