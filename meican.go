package gomeican

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	prefix = "https://meican.com/preorder/api/v2.1"
)

// DinnerOrder 餐的订单信息，包含时间点和多家可选餐厅信息
type DinnerOrder struct {
	TimeInfo            *calendarItem
	RestaurantFoodInfos []RestaurantFood
}

// RestaurantFood 一家餐厅的餐的信息
type RestaurantFood struct {
	RestaurantInfo *Restaurant
	FoodList       []FoodInfo
}

type Meican struct {
	token string

	client *http.Client
}

func NewMeican(token string) *Meican {
	return &Meican{
		token:  token,
		client: &http.Client{},
	}
}

// GetTodayList 获取今天的可预定列表。
// 返回值列表表示：午餐、晚餐
func (meican *Meican) GetTodayOrderList(ctx context.Context) ([]DinnerOrder, error) {

	return meican.GetOrderList(ctx, time.Now())
}

// GetOrderList 获取指定日期的可预定列表。
func (meican *Meican) GetOrderList(ctx context.Context, d time.Time) ([]DinnerOrder, error) {

	// 周末calList为空。
	// 工作日为午餐、晚餐
	calList, err := meican.getCalendarItems(ctx, d)
	panicError(err)

	if len(calList) == 0 {
		panic("no valid order")
	}

	orders := make([]DinnerOrder, len(calList))
	for i := 0; i < len(calList); i++ {

		orders[i], err = meican.getTimePeriodOrderInfo(ctx, calList[i], d)
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}

// 获取指定时间段的餐厅的订餐信息
func (meican *Meican) getTimePeriodOrderInfo(ctx context.Context, calItem calendarItem, d time.Time) (orderInfo DinnerOrder, err error) {
	orderInfo.TimeInfo = &calItem

	targetTime := fmt.Sprintf("%s+%s", d.Format("2006-01-02"), calItem.OpeningTime.CloseTime)
	tabID := calItem.UserTab.UniqueID

	// 获取该时间段的餐厅列表
	var rests []*Restaurant
	if rests, err = meican.listRestaurant(ctx, tabID, targetTime); err != nil {
		return
	}

	// 获取每一家餐厅的所有餐
	restOrderInfos := make([]RestaurantFood, len(rests))
	for i, rest := range rests {
		restOrderInfos[i].RestaurantInfo = rest

		restOrderInfos[i].FoodList, err = meican.getRestaurantFood(ctx, tabID, targetTime, rest.ID)
		if err != nil {
			return orderInfo, err
		}
	}

	orderInfo.RestaurantFoodInfos = restOrderInfos
	return
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func (meican *Meican) getAndUnmarshal(ctx context.Context, path string, v interface{}) error {

	data, err := meican.get200(ctx, path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

func (meican *Meican) get200(ctx context.Context, path string) ([]byte, error) {
	return meican.send(ctx, http.MethodGet, path, http.StatusOK, nil)
}

func (meican *Meican) send(ctx context.Context, method, path string, wantStatus int, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, prefix+path, body)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "remember", Value: meican.token})

	resp, err := meican.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != wantStatus {
		return nil, errors.Errorf("HTTP %s: %s (expected %v)", method, path, wantStatus)
	}

	return data, nil
}
