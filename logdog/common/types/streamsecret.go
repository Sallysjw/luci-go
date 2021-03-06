// Copyright 2015 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package types

import (
	"crypto/rand"
	"fmt"
)

const (
	// PrefixSecretLength is the size, in bytes, of the stream secret.
	//
	// This value was chosen such that it is:
	// - Sufficiently large to avoid collisions.
	// - Can be expressed as a Base64 string without ugly padding.
	PrefixSecretLength = 36
)

// PrefixSecret is the prefix secret value. It is used to assert ownership of
// a prefix space.
//
// The Prefix secret is generated by the Coordinator at prefix registration,
// and is included by the Butler to prove that it is the entity that registered
// the stream. The secret is asserted by microservices and the Coordinator
// during Butler-initiated stream operations.
type PrefixSecret []byte

// NewPrefixSecret generates a new, default-length secret parameter.
func NewPrefixSecret() (PrefixSecret, error) {
	buf := make([]byte, PrefixSecretLength)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	value := PrefixSecret(buf)
	if err := value.Validate(); err != nil {
		panic(err)
	}
	return value, nil
}

// Validate confirms that this prefix secret is conformant.
//
// Note that this does not scan the byte contents of the secret for any
// security-related parameters.
func (s PrefixSecret) Validate() error {
	if len(s) != PrefixSecretLength {
		return fmt.Errorf("invalid prefix secret length (%d != %d)", len(s), PrefixSecretLength)
	}
	return nil
}
