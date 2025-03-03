package e2e

import (
	"testing"

	"github.com/stefanoranzini/assessment/cart-service/internal/server"
)

func TestOrder(t *testing.T) {
	server.New(0) //TODO: provide empty random port or model server to handle random port
}
