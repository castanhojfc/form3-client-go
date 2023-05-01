//go:build unit

package form3_test

import (
	"testing"

	"github.com/castanhojfc/form3-client-go/form3"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Run("can create client with API url", func(*testing.T) {
		t.Parallel()

		client, _ := form3.NewClient()

		assert.NotNil(t, client)
	})
}
