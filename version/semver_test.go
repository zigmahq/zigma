package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/version"
)

func TestCompare(t *testing.T) {
	v1 := version.Version{Number: "1.0.0"}
	v2 := version.Version{Number: "1.0.1"}
	assert.Equal(t, -1, v1.Compare(v2))
	assert.Equal(t, 1, v2.Compare(v1))

	v1.Number = "1.0.1"
	assert.Equal(t, 0, v1.Compare(v2))
	assert.Equal(t, 0, v2.Compare(v1))

	v2.Number = "0.9.1"
	assert.Equal(t, 1, v1.Compare(v2))
	assert.Equal(t, -1, v2.Compare(v1))
}

func TestNewerThan(t *testing.T) {
	v1 := version.Version{Number: "1.0.0"}
	v2 := version.Version{Number: "1.0.1"}
	assert.False(t, v1.NewerThan(v2))
	assert.True(t, v2.NewerThan(v1))

	v1.Number = "1.0.1"
	assert.False(t, v1.NewerThan(v2))
	assert.False(t, v2.NewerThan(v1))

	v2.Number = "0.9.1"
	assert.True(t, v1.NewerThan(v2))
	assert.False(t, v2.NewerThan(v1))
}

func TestOlderThan(t *testing.T) {
	v1 := version.Version{Number: "1.0.0"}
	v2 := version.Version{Number: "1.0.1"}
	assert.True(t, v1.OlderThan(v2))
	assert.False(t, v2.OlderThan(v1))

	v1.Number = "1.0.1"
	assert.False(t, v1.OlderThan(v2))
	assert.False(t, v2.OlderThan(v1))

	v2.Number = "0.9.1"
	assert.False(t, v1.OlderThan(v2))
	assert.True(t, v2.OlderThan(v1))
}
