package gomeican

import (
	"context"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

type calendarItem struct {
	TargetTime int // unix time, (ms)
	Title      string
	UserTab    struct {
		UniqueID string
	}
	OpeningTime struct {
		Name      string // 午餐、晚餐
		OpenTime  string
		CloseTime string
	}
}

func (meican *Meican) getTodayCalendarItems(ctx context.Context) ([]calendarItem, error) {
	return meican.getCalendarItems(ctx, time.Now())
}

func (meican *Meican) getCalendarItems(ctx context.Context, day time.Time) ([]calendarItem, error) {
	path := "/calendaritems/list"
	targetPath := path + "?" + makeGetCalendarItemsParams(day) + "?client_id=Xqr8w0Uk4ciodqfPwjhav5rdxTaYepD&client_secret=vD11O6xI9bG3kqYRu9OyPAHkRGxLh4E"

	logrus.Infof("target path: %s", targetPath)

	var calendarInfo struct {
		DateList []struct {
			Date             string         `json:"date"`
			CalendarItemList []calendarItem `json:"calendarItemList"`
		} `json:"dateList"`
	}
	if err := meican.get200AndUnmarshal(ctx, targetPath, &calendarInfo); err != nil {
		return nil, err
	}

	if dateList := calendarInfo.DateList; len(dateList) == 0 {
		return []calendarItem{}, nil
	} else {
		// 应该是多个日期的话，按照日期进行排序，我们只取第一个日期的清单
		return dateList[0].CalendarItemList, nil
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
