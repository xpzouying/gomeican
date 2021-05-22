package gomeican

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestMeican() *Meican {

	token := os.Getenv("MEICAN_TOKEN")
	if token == "" {
		panic("empty token. please set MEICAN_TOKEN")
	}

	return NewMeican(token)
}

func TestGetOrderList(t *testing.T) {
	meican := newTestMeican()

	d := time.Date(2021, time.Month(5), 25, 0, 0, 0, 0, time.Local)
	orders, err := meican.GetOrderList(context.Background(), d)
	assert.NoError(t, err)

	for _, order := range orders {
		t.Logf("order: %+v", order)
	}

}
