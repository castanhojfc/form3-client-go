package form3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Run("can create client with API url", func(*testing.T) {
		client := NewClient()

		assert.NotNil(t, client)
	})
}
