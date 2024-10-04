package main

import (
	"fmt"
	"sort"
	"time"
)


type Match struct {
	Ask *Order
	Bid *Order
	SizeFilled float64
	Price float64

}
type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}


type Orders []*Order
func (o Orders) len() int 				{ return len(o) }
func (o Orders) swap(i, j int) int 		{  o[i], o[j] = o[j], o[i] }
func (o Orders) less(i, j int) bool 	{ return o[i].Timestamp< o[j].Timestamp }



func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("[size: %.2f]", o.Size)
}

type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

type ByBestAsk struct { Limits }
func (a ByBestAsk) len() int { return len(a.Limits) }
func (a ByBestAsk) swap(i, j int) int { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsk) less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestBids struct { Limits }
func (b ByBestBids) len() int { return len(b.Limits) }
func (b ByBestBids) swap(i, j int) int { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestBids) less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

// func (l *Limit) String() string {
// 	return fmt.Sprintf("[price: %.2f, total volume: %.2f]", l.Price, l.TotalVolume)
// }

func (l *Limit) AddOrder(o *Order){
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size

}



func (l *Limit) DeleteOrder(o *Order) {
	for i :=0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}

	o.Limit = nil
	l.TotalVolume -= o.Size

	sort.Sort[l.Orders]
	// TODO; Reset the whole resting orders
}


type Orderbook struct {
	Asks []*Limit
	Bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Asks: []*Limit{},
		Bids: []*Limit{},

		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}


func (ob *Orderbook) PlaceOrder(price float64, o *Order) []Match {
	// 1. Try to match the orders
	// matching logic

	// 2. add the rest of the order to the orderbook

	if o.Size > 0.0 {
		ob.add(price, o)
	}

	return []Match{}
}

// BUY 10 BTC  => 15k

func (ob *Orderbook) add(price float64, o *Order){
	var limit *Limit

	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		limit.AddOrder(o)
		if o.Bid {
			ob.Bids = append(ob.Bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.Asks = append(ob.Asks, limit)
			ob.AskLimits[price] = limit
		}
	}
}