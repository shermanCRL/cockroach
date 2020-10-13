// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package cloudimpl

import "net/url"

// RedactKMSURI redacts the Master Key ID and the ExternalStorage secret
// credentials.
func RedactKMSURI(kmsURI *url.URL) *url.URL {
	sanitizedKMSURI := SanitizeExternalStorageURI(kmsURI, nil)
	// Redact the path which contains the KMS Master Key identifier.
	sanitizedKMSURI.Path = "/redacted"

	return sanitizedKMSURI
}
