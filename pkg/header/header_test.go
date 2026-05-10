package header

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVary(t *testing.T) {
	t.Run("empty string returns nil", func(t *testing.T) {
		result := ParseVary("")

		assert.Nil(t, result)
	})

	t.Run("asterisk returns nil", func(t *testing.T) {
		result := ParseVary("*")

		assert.Nil(t, result)
	})

	t.Run("single header is canonicalized", func(t *testing.T) {
		result := ParseVary("accept-encoding")

		assert.Equal(t, []string{"Accept-Encoding"}, result)
	})

	t.Run("multiple headers are trimmed and canonicalized", func(t *testing.T) {
		result := ParseVary("accept-encoding,  accept-language")

		assert.Equal(t, []string{"Accept-Encoding", "Accept-Language"}, result)
	})
}

func TestSplit(t *testing.T) {
	t.Run("empty string returns empty slice", func(t *testing.T) {
		result := Split("")

		assert.Empty(t, result)
	})

	t.Run("single value is lowercased", func(t *testing.T) {
		result := Split("Foo")

		assert.Equal(t, []string{"foo"}, result)
	})

	t.Run("multiple values are trimmed and lowercased", func(t *testing.T) {
		result := Split(" Foo ,  BAR ")

		assert.Equal(t, []string{"foo", "bar"}, result)
	})
}
