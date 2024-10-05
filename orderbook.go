package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}
type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

type Orders []*Order

// Len implements sort.Interface.
func (o Orders) Len() int {
	panic("unimplemented")
}

// Less implements sort.Interface.
func (o Orders) Less(i int, j int) bool {
	panic("unimplemented")
}

// Swap implements sort.Interface.
func (o Orders) Swap(i int, j int) {
	panic("unimplemented")
}

func (o Orders) len() int           { return len(o) }
func (o Orders) swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

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

func (o *Order) isFilled() bool {
	return o.Size == 0.0
}

type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

type ByBestAsk struct{ Limits }

func (a ByBestAsk) len() int           { return len(a.Limits) }
func (a ByBestAsk) swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsk) less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestBid struct{ Limits }

func (b ByBestBid) len() int           { return len(b.Limits) }
func (b ByBestBid) swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestBid) less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

// func (l *Limit) String() string {
// 	return fmt.Sprintf("[price: %.2f, total volume: %.2f]", l.Price, l.TotalVolume)
// }

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size

}

func (l *Limit) DeleteOrder(o *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}

	o.Limit = nil
	l.TotalVolume -= o.Size

	sort.Sort(l.Orders)
	// TODO; Reset the whole resting orders
}


func (l *Limit) Fill(o *Order) []Match {
	matches := []Match{}
	for _, order := range l.Orders {
		match := l.fillOrder(order, o)
		matches = append(matches, match)

		if o.isFilled() {
			break
		}
	}

	return matches
}

func (l *Limit) fillOrder(a, b *Order) Match {
	
	var (
		bid *Order
		ask *Order
		sizeFilled float64
	)

	if a.Bid{
		bid = a
		ask = b
	} else {
		bid = b
		ask = a
	}
	
	if a.Size >= b.Size {
		a.Size -= b.Size
		sizeFilled = b.Size
		b.Size = 0.0
	} else {
		b.Size -= a.Size
		sizeFilled = a.Size
		a.Size = 0.0
	}
	
	return Match {
		Bid: bid,
		Ask: ask,
		SizeFilled: sizeFilled,
		Price: l.Price,
	}
}


type Orderbook struct {
	asks []*Limit
	bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		asks: []*Limit{},
		bids: []*Limit{},

		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *Orderbook) placeMarketOrder(o *Order) []Match {
	matches := []Match{}
	
	if o.Bid {
		if o.Size > ob.AskTotalVolume() {
			panic("Not enough volume to fill order")
		}

		for _, limit := range ob.Asks() {
			limitMatches := limit.Fill(o)
			matches := append(matches, limitMatches...)

		}
	} else{
		
	}

	return matches
	
	 
}

func (ob *Orderbook) PlaceLimitOrder(price float64, o *Order) {
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
			ob.bids = append(ob.bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.asks = append(ob.asks, limit)
			ob.AskLimits[price] = limit
		}
	}
}


func (ob *Orderbook) BidTotalVolume() float64 {
	totalVolume := 0.0
	for i := 0; i < len(ob.bids); i++ {
		totalVolume += ob.bids[i].TotalVolume
	}

	return totalVolume
}

func (ob *Orderbook) AskTotalVolume() float64 {
	totalVolume := 0.0
	for i := 0; i < len(ob.asks); i++ {
		totalVolume += ob.asks[i].TotalVolume
	}

	return totalVolume
}

func (ob *Orderbook) Asks() []*Limit {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks	
}

func (ob *Orderbook) Bids() []*Limit {
	sort.Sort(ByBestBids{ob.bids})
	return ob.bids
}
