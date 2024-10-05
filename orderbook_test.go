package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLimit(t *testing.T){
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true,5 )
	buyOrderB := NewOrder(true,8 )
	buyOrderC := NewOrder(true,10 )

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrderB)
	
	fmt.Println(l)

}


func assert(t *testing.T, a, b any){
	if !reflect.DeepEqual(a,b) {
		t.Errorf("%+v != %+v", a,b)
	}
}


func TestPlaceLimitOrder(t *testing.T){
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 10)
	sellOrderB := NewOrder(false, 5)
	ob.PlaceLimitOrder(10_000, sellOrderA)
	ob.PlaceLimitOrder(9_000, sellOrderB)

	assert(t, len(ob.asks), 2)
}