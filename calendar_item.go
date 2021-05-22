package gomeican

import (
	"context"
	"net/url"
	"time"
)

type calendarItem struct {
	TargetTime int    `json:"targetTime"` // unix time, (ms)
	Title      string `json:"title"`
	UserTab    struct {
		UniqueID string `json:"uniqueId"`
	} `json:"userTab"`
	OpeningTime struct {
		Name      string `json:"name"` // 午餐、晚餐
		OpenTime  string `json:"openTime"`
		CloseTime string `json:"closeTime"`
	} `json:"openingTime"`
}

func (meican *Meican) getTodayCalendarItems(ctx context.Context) ([]calendarItem, error) {
	return meican.getCalendarItems(ctx, time.Now())
}

func (meican *Meican) getCalendarItems(ctx context.Context, day time.Time) ([]calendarItem, error) {
	path := "/calendaritems/list"
	targetPath := path + "?" + makeGetCalendarItemsParams(day)

	var calendarInfo struct {
		DateList []struct {
			CalendarItemList []calendarItem `json:"calendarItemList"`
		} `json:"dateList"`
	}
	if err := meican.getAndUnmarshal(ctx, targetPath, &calendarInfo); err != nil {
		return nil, err
	}

	if list := calendarInfo.DateList; len(list) == 0 {
		return []calendarItem{}, nil
	} else {
		// 应该是多个日期的话，按照日期进行排序，我们只取第一个日期的清单
		return list[0].CalendarItemList, nil
	}
}

func makeGetCalendarItemsParams(day time.Time) string {
	values := url.Values{}
	date := day.Format("2006-01-02") // want: 2021-05-25
	values.Add("beginDate", date)
	values.Add("endDate", date)
	values.Encode()

	return values.Encode()
}
