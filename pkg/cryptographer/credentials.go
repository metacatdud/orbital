package cryptographer

import (
	"bytes"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

// Credential generator and cryptographically secure and verifiable
// credentials for multiple purposes:
// - storage paths
// - local credentials
// - token and keys for scoped actions (storage:write, someApp:listPosts)

// CredsV1 version control for credential format
const CredsV1 = "v1"

// Generate a base32 encoded key without the = symbol
var b32 = base32.StdEncoding.WithPadding(base32.NoPadding)

// CredentialsScopeID for credentials
// this can be reproduced by anyone given they provide the same data
// used to generate an identification for the owner in regard to a target
// ownerPubKey (machine or user) + devicePubKey + service (provided by the target) + label (arbitrary string) + `version` (system provided)
func CredentialsScopeID(ownerPubKey, devicePubKey []byte, service, label, version string) string {
	key := "orbital/scope:%s"

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(key, version))
	buf.WriteByte('|')
	buf.Write(ownerPubKey)
	buf.Write(devicePubKey)
	buf.WriteByte('|')
	buf.WriteString(service)
	buf.WriteByte('|')
	buf.WriteString(label)

	sum := sha256.Sum256(buf.Bytes())
	return b32.EncodeToString(sum[:10])
}

// CredentialsRoot create a root key for a namespace
// This will allow a device to generate a predictable LOCAL secret for a public key (scope)
func CredentialsRoot(ownerPubKey, deviceSecretKey []byte, scopeID, version, epoch string) ([]byte, error) {
	key := "orbital/root:%s"

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(key, version))
	buf.WriteByte('|')
	buf.WriteString(scopeID)
	buf.WriteByte('|')
	buf.WriteString(epoch)

	r := hkdf.New(sha256.New, deviceSecretKey, ownerPubKey, buf.Bytes())
	root := make([]byte, 32)
	_, err := io.ReadFull(r, root)
	if err != nil {
		return nil, err
	}
	return root, nil
}

// CredentialsDerive create a derived key for a namespace
// This will allow a device to generate a predictable LOCAL secret for a public key (scope)
// TODO: add a version control for the 'purpose'?
func CredentialsDerive(root []byte, purpose string, length int) ([]byte, error) {
	key := "orbital/cred:%s"

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(key, purpose))

	r := hkdf.New(sha256.New, root, nil, buf.Bytes())
	out := make([]byte, length)
	_, err := io.ReadFull(r, out)
	return out, err
}
