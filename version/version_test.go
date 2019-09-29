package version_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/version"
)

func TestCurrent(t *testing.T) {
	v := version.Current

	assert.NotEmpty(t, v.Number)
	assert.NotEmpty(t, v.Name)
	assert.NotEmpty(t, v.Signature)
	assert.True(t, v.Verify())
}

func TestSignVerifyVersion(t *testing.T) {
	pub, pri, err := ed25519.GenerateKey(rand.Reader)
	assert.Nil(t, err)

	ver := version.Version{
		Number: "1.0.0",
		Name:   "winter-sunset",
	}
	sig, err := ver.Sign(pri)
	assert.Nil(t, err)

	ver.Signature = sig
	assert.NotEmpty(t, sig)
	assert.NotEmpty(t, ver.Signature)

	valid := ver.Verify(pub)
	assert.True(t, valid)
}
