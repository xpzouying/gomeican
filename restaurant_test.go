package gomeican

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListStores(t *testing.T) {
	tabid := "4fd6c5fb-3376-4f49-aa44-95d9ed5e6d86"
	targetTime := "2021-05-25+17:00"

	meican := newTestMeican()
	rests, err := meican.listRestaurant(context.Background(), tabid, targetTime)
	assert.NoError(t, err)

	for _, rest := range rests {
		t.Logf("%+v", rest)
	}
}
