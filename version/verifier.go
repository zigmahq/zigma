package version

import (
	"crypto/ed25519"
	"encoding/hex"
)

// Verifier for version package
var Verifier verifier

// Verifier holds a PublicKey and Algorithm to verify signatures
type verifier []byte

// WithPub initializes a new verifier with a different public key
func (v verifier) WithPub(pub ed25519.PublicKey) verifier {
	return verifier(pub)
}

// Verify reports whether sig is a valid signature of message by publicKey.
func (v verifier) Verify(message, sig []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(v), message, sig)
}

// String returns verifier public key in string
func (v verifier) String() string {
	return hex.EncodeToString(v)
}

// Sign signs the version with a release private key
func (v Version) Sign(pri ed25519.PrivateKey) ([]byte, error) {
	b, err := v.Marshal()
	if err != nil {
		return nil, err
	}
	sig := ed25519.Sign(pri, b)
	v.Signature = sig

	return sig, nil
}

// Verify reports whether the version is signed
func (v Version) Verify(pub ...ed25519.PublicKey) bool {
	ver := Version{
		Number: v.Number,
		Name:   v.Name,
	}

	b, err := ver.Marshal()
	if err != nil {
		return false
	}

	vf := Verifier
	if len(pub) > 0 {
		vf = vf.WithPub(pub[0])
	}
	return vf.Verify(b, v.Signature)
}
