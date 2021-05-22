package gomeican

import (
	"context"
	"fmt"
)

type Restaurant struct {
	Name string `json:"name"`
	ID   string `json:"uniqueId"`
	Tel  string `json:"tel"`

	DishLimit          int `json:"dishLimit"`
	AvailableDishCount int `json:"availableDishCount"`
}

// listRestaurant 获取可点餐的全部餐厅。依赖于 getTodayCalendarItems()获取时间段信息。
// 下列的请求参数从calendar的函数中获取，
// tabid - 时间段的tab id, `56a3c9ce-0a4d-4c35-8166-821174433d3a`
// target-time date+closeTime, `2021-05-25+09:30`
func (meican *Meican) listRestaurant(ctx context.Context, timeTabID, targetTime string) ([]*Restaurant, error) {
	// 不能url encode，直接拼接字符串
	const path = "/restaurants/list"
	target := fmt.Sprintf("%s?tabUniqueId=%s&targetTime=%s", path, timeTabID, targetTime)

	var res struct {
		RestaurantList []*Restaurant `json:"restaurantList"`
	}
	if err := meican.getAndUnmarshal(ctx, target, &res); err != nil {
		return nil, err
	}

	return res.RestaurantList, nil
}

type FoodInfo struct {
	Name string
}

// tabid - 时间段的tab id, `56a3c9ce-0a4d-4c35-8166-821174433d3a`
// target-time date+closeTime, `2021-05-25+09:30`
func (meican *Meican) getRestaurantFood(ctx context.Context, timeTabID, targetTime, restUID string) ([]FoodInfo, error) {
	// 不能url encode，直接拼接字符串
	const path = "/restaurants/show"
	target := fmt.Sprintf("%s?tabUniqueId=%s&targetTime=%s&restaurantUniqueId=%s", path, timeTabID, targetTime, restUID)

	var res struct {
		DishList []FoodInfo `json:"dishList"`
	}
	if err := meican.getAndUnmarshal(ctx, target, &res); err != nil {
		return nil, err
	}

	return res.DishList, nil
}
