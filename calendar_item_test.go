package gomeican

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMakeGetCalendarItemsParams(t *testing.T) {
	day := time.Date(2021, time.Month(05), 21, 0, 0, 0, 0, time.Local)
	got := makeGetCalendarItemsParams(day)

	assert.Equal(t, "beginDate=2021-05-21&endDate=2021-05-21", got)
}

func TestGetCalendarItems(t *testing.T) {

	day := time.Date(2021, time.Month(5), 25, 0, 0, 0, 0, time.Local)

	meican := NewMeican("", "")

	got, err := meican.getCalendarItems(context.Background(), day)
	assert.NoError(t, err)

	for _, item := range got {
		t.Logf("%+v", item)
	}
}
