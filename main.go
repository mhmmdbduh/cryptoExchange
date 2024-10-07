package main

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/mhmmdbduh/cryptoExchange/orderbook"
)

func main() {

	e:= echo.New()
	ex := NewExchange()

	e.POST("/order", ex.handlePlaceOrder)

	e.Start(":3000")


}	
type OrderType string

const (
	MarketOrder OrderType = "MARKET"
	LimitOrder OrderType = "LIMIT"
)

type Market string

const (
	MarketETH Market = "ETH"
)

type Exchange struct {
	orderbooks map[Market]*orderbook.Orderbook

}

func NewExchange() *Exchange {
	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderbook()

	return &Exchange{
		orderbooks: orderbooks,
	}
}

type PlaceOrderRequest struct {
	Type OrderType //Limit or Market
	Bid bool
	Size float64
	Price float64
	Market Market
}

func (ex *Exchange) handlePlaceOrder(c echo.Context) error{
	var placeOrderData PlaceOrderRequest
	if err:= json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}	
	
	market := Market(placeOrderData.Market)
	ob := ex.orderbooks[market]
	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size)

	ob.PlaceLimitOrder(placeOrderData.Price, order)
	
	return c.JSON(200, map[string]any{"msg": "order placed"})
}



