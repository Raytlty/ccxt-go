package bitmex

import (
    . "github.com/Raytlty/ccxt-go/base"
    "testing"
    //"time"
    //"fmt"
)

func testExchange() *BitMex {
    config := &ExchangeConfig {
        ApiKey: "",
        Secret: "",
        Timeout: 15,
        Test: true,
        Proxy: "http://localhost:8123",
    }
    ex := NewBitMEX(config)
    return ex
}

func TestBitMEX_GetOrderBook(t *testing.T) {
	ex := testExchange()
	ob, err := ex.GetOrderBook("XBTUSD", 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ob)
}
