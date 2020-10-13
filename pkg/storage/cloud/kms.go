// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package cloud

import (
	"context"
	"net/url"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/testutils/lint/passes/fmtsafe/testdata/src/github.com/cockroachdb/errors"
)

// KMS provides an API to interact with a KMS service.
type KMS interface {
	// MasterKeyID will return the identifier used to reference the master key
	// associated with the KMS object.
	MasterKeyID() (string, error)
	// Encrypt returns the ciphertext version of data after encrypting it using
	// the KMS.
	Encrypt(ctx context.Context, data []byte) ([]byte, error)
	// Decrypt returns the plaintext version of data after decrypting it using the
	// KMS.
	Decrypt(ctx context.Context, data []byte) ([]byte, error)
	// Close may be used to perform the necessary cleanup and shutdown of the
	// KMS connection.
	Close() error
}

// KMSEnv is the environment in which a KMS is configured and used.
type KMSEnv interface {
	ClusterSettings() *cluster.Settings
	KMSConfig() *base.ExternalIODirConfig
}

// KMSFromURIFactory describes a factory function for KMS given a URI.
type KMSFromURIFactory func(uri *url.URL, env KMSEnv) (KMS, error)

// Mapping from KMS scheme to its registered factory method.
var kmsFactoryMap = make(map[string]KMSFromURIFactory)

// RegisterKMSFromURIFactory is used by every concrete KMS implementation to
// register its factory method.
func RegisterKMSFromURIFactory(factory KMSFromURIFactory, scheme string) {
	if _, ok := kmsFactoryMap[scheme]; ok {
		panic("factory method for " + scheme + " has already been registered")
	}
	kmsFactoryMap[scheme] = factory
}

// KMSFromURI is the method used to create a KMS instance from the provided URI.
func KMSFromURI(uri *url.URL, env KMSEnv) (KMS, error) {
	// TODO(sherman): this method previously used url.ParseRequestURI instead of url.Parse. Was this important?

	// Find the factory method for the given KMS scheme.
	factory, ok := kmsFactoryMap[uri.Scheme]
	if !ok {
		return nil, errors.Newf("no factory method found for scheme %s", uri.Scheme)
	}

	return factory(uri, env)
}
