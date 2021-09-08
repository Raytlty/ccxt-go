package bitmex

import (
    . "github.com/Raytlty/ccxt-go/base"
    "github.com/frankrap/bitmex-api"
    "strings"
    "github.com/frankrap/bitmex-api/swagger"
    "sort"
    "time"
    "errors"
)

type BitMex struct {
	client *bitmex.BitMEX
	config *ExchangeConfig
	symbol string
}

func NewBitMEX(config *ExchangeConfig) *BitMex {
	baseURI := "www.bitmex.com"
	if config.Test {
		baseURI = "testnet.bitmex.com"
	}
	client := bitmex.New(
        nil,
		baseURI,
        config.ApiKey,
        config.Secret,
        config.Verbose,
    )
	if strings.HasPrefix(config.Proxy, "socks5") {
		socks5Proxy := strings.ReplaceAll(config.Proxy, "socks5:", "")
		client.SetProxy(socks5Proxy)
	} else if strings.HasPrefix(config.Proxy, "http://") {
		client.SetHttpProxy(config.Proxy)
	}
	if config.Websocket {
		client.StartWS()
	}
	return &BitMex{
		client: client,
		config: config,
	}
}

func (b BitMex) GetSymbol() (symbol string, err error) {
    return b.symbol, nil
}

func (b *BitMex) SetSymbol(symbol string) (err error) {
    b.symbol = symbol
    return
}

func (b *BitMex) GetName() (string) {
    return "bitmex"
}

func (b *BitMex) GetTime() (tm int64, err error) {
    version, _, err := b.client.GetVersion()
    if err != nil {
        return
    } else {
        tm = version.Timestamp
    }
    return
}

func (b *BitMex) GetOrderBook(symbol string, depth int) (result *OrderBook, err error) {
    result = &OrderBook{}
    ret, err := b.client.GetOrderBook(depth, symbol)
    if err != nil {
        return
    }
    for _, v := range ret.Asks {
        result.Asks = append(result.Asks, [2]float64 {
            v.Price,
            v.Amount,
        })
    }
    for _, v := range ret.Bids {
        result.Bids = append(result.Bids, [2]float64 {
            v.Price,
            v.Amount,
        })
    }
    result.Symbol = symbol
    result.Timestamp = ret.Timestamp.Unix()
    return
}

func (b *BitMex) GetOHLCV(
    symbol string,
    period string,
    start int64,
    end int64,
    limit int64,
) (results []*OHLCV, err error) {
    var binSize string;
    if strings.HasPrefix(period, "m") ||
        strings.HasPrefix(period, "h") ||
        strings.HasPrefix(period, "d") {
            binSize = period
        } else {
        binSize = period + "m"
    }
    o, err := b.client.GetBucketed(
        symbol,
        binSize,
        false,
        "",
        "",
        float32((limit)),
        -1,
        false,
        time.Unix(start, 0),
        time.Unix(end, 0))
    if err != nil {
        return
    }
    for _, v := range o {
        results = append(results, &OHLCV{
            Symbol:     v.Symbol,
            Timestamp:  v.Timestamp.Unix(),
            Open:       v.Open,
            High:       v.High,
            Close:      v.Close,
            Low:        v.Low,
            Volume:     float64(v.Volume),
        })
    }
    return
}

func (b *BitMex) PlaceOrder(
    symbol string,
    direction Direction,
    orderType OrderType,
    price float64,
	size float64,
) (result *Order, err error) {
    var side string
    var otype string
    if direction == BUY {
        side = bitmex.SIDE_BUY
    } else {
        side = bitmex.SIDE_SELL
    }
    if orderType == LIMIT {
        otype = bitmex.ORD_TYPE_LIMIT
    } else if orderType == MARKET {
        otype = bitmex.ORD_TYPE_MARKET
    } else if orderType == STOP_LIMIT {
        otype = bitmex.ORD_TYPE_STOP_LIMIT
    } else if orderType == STOP_MARKET {
        otype = bitmex.ORD_TYPE_STOP
    }

    order, err := b.client.PlaceOrder(
        side,
        otype,
        price,
        price,
        int32(size),
        "",
        "",
        symbol,
    )
	if err != nil {
		return
	}
	result = b.parseOrder(&order)
	return
}

func (b *BitMex) parseOrder(order *swagger.Order) (result *Order) {
	result = &Order{}
	result.ID = order.OrderID
	result.Symbol = order.Symbol
	result.Price = order.Price
	result.StopPx = order.StopPx
	result.Amount = float64(order.OrderQty)
	result.Direction = b.convertDirection(order.Side)
	result.Type = b.convertOrderType(order.OrdType)
	result.AvgPrice = order.AvgPx
	result.FilledAmount = float64(order.CumQty)
	if strings.Contains(order.ExecInst, "ParticipateDoNotInitiate") {
		result.PostOnly = true
	}
	if strings.Contains(order.ExecInst, "ReduceOnly") {
		result.ReduceOnly = true
	}
	result.Status = b.convertOrderStatus(order)
	return
}

func (b BitMex) convertDirection(side string) Direction {
	switch side {
	case bitmex.SIDE_BUY:
		return BUY
	case bitmex.SIDE_SELL:
		return BUY
	default:
		return BUY
    }
}

func (b BitMex) convertOrderType(orderType string) OrderType {
	switch orderType {
	case bitmex.ORD_TYPE_LIMIT:
		return LIMIT
	case bitmex.ORD_TYPE_MARKET:
		return MARKET
	case bitmex.ORD_TYPE_STOP_LIMIT:
		return STOP_LIMIT
	case bitmex.ORD_TYPE_STOP:
		return STOP_MARKET
	default:
		return LIMIT
	}
}

func (b *BitMex) convertOrderStatus(order *swagger.Order) OrderStatus {
    switch order.OrdStatus {
	case bitmex.OS_NEW:
		return NEW
	case bitmex.OS_PARTIALLY_FILLED:
		return PARTIALLYFILLED
	case bitmex.OS_FILLED:
		return FILLED
	case bitmex.OS_CANCELED:
		return CANCELED
	case bitmex.OS_REJECTED:
		return REJECTED
	default:
		return CREATED
	}
}

func (b *BitMex) convertPosition(position *swagger.Position) (result *Position) {
	result = &Position{}
	result.Symbol = position.Symbol
	result.OpenTime = (time.Time{}).Unix()
	result.OpenPrice = position.AvgEntryPrice
	result.Size = float64(position.CurrentQty)
	result.AvgPrice = position.AvgCostPrice
	return
}

func (b *BitMex) SubscribeTrades(market string, callback func(trades []*Trade)) error {
	if !b.config.Websocket {
		return errors.New("websocket disabled")
	}
	b.client.On(bitmex.BitmexWSTrade, func(trades []*swagger.Trade, action string) {
		var data []*Trade
		for _, v := range trades {
			var direction Direction
			if v.Side == bitmex.SIDE_BUY {
				direction = BUY
			} else if v.Side == bitmex.SIDE_SELL {
				direction = SELL
			}
			data = append(data, &Trade{
				ID:        v.TrdMatchID,
				Direction: direction,
				Price:     v.Price,
				Amount:    float64(v.Size),
				Ts:        v.Timestamp.UnixNano() / int64(time.Millisecond),
				Symbol:    v.Symbol,
			})
		}
		callback(data)
	})
	subscribeInfos := []bitmex.SubscribeInfo{
		{Op: bitmex.BitmexWSTrade, Param: market},
	}
	err := b.client.Subscribe(subscribeInfos)
	return err
}

func (b *BitMex) SubscribeLevel2Snapshots(market string, callback func(ob *OrderBook)) error {
	if !b.config.Websocket {
		return errors.New("websocket disable")
	}
	b.client.On(bitmex.BitmexWSOrderBookL2, func(m bitmex.OrderBookDataL2, symbol string) {
		var ob OrderBook

		ob.Symbol = symbol
		ob.Timestamp = m.Timestamp.Unix()

		for _, v := range m.RawData {
			switch v.Side {
			case "Buy":
				ob.Bids = append(ob.Bids, [2]float64{
					v.Price,
					float64(v.Size),
				})
			case "Sell":
				ob.Asks = append(ob.Asks, [2]float64{
					v.Price,
					float64(v.Size),
				})
			}
		}

		sort.Slice(ob.Bids, func(i, j int) bool {
			return ob.Bids[i][0] > ob.Bids[j][0]
		})

		sort.Slice(ob.Asks, func(i, j int) bool {
			return ob.Asks[i][0] < ob.Asks[j][0]
		})

		callback(&ob)
	})
	subscribeInfos := []bitmex.SubscribeInfo{
		{Op: bitmex.BitmexWSOrderBookL2, Param: market},
	}
	err := b.client.Subscribe(subscribeInfos)
	return err
}

func (b *BitMex) SubscribeOrders(market string, callback func(orders []*Order)) error {
	if !b.config.Websocket {
		return errors.New("websocket disable")
	}
	b.client.On(bitmex.BitmexWSOrder, func(m []*swagger.Order, action string) {
		var orders []*Order
		for _, v := range m {
			order := b.parseOrder(v)
			orders = append(orders, order)
		}
		callback(orders)
	})
	subscribeInfos := []bitmex.SubscribeInfo{
		{Op: bitmex.BitmexWSOrder, Param: market},
	}
	err := b.client.Subscribe(subscribeInfos)
	return err
}

func (b *BitMex) SubscribePositions(market string, callback func(positions []*Position)) error {
	if !b.config.Websocket {
		return errors.New("websocket disable")
	}
	b.client.On(bitmex.BitmexWSPosition, func(m []*swagger.Position, action string) {
		var positions []*Position
		for _, v := range m {
			positions = append(positions, b.convertPosition(v))
		}
		callback(positions)
	})
	subscribeInfos := []bitmex.SubscribeInfo{
		{Op: bitmex.BitmexWSPosition, Param: market},
	}
	err := b.client.Subscribe(subscribeInfos)
	return err
}
