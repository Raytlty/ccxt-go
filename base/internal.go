package base

type ExchangeConfig struct {
	ApiKey          string `json:"apiKey"`
    Password        string `json:"password"`
	Secret          string `json:"secret"`
	Timeout         int64  `json:"timeout"`
	EnableRateLimit bool   `json:"enableRateLimit"`
	Test            bool   `json:"test"`
	Verbose         bool   `json:"verbose"`
    Websocket       bool   `json:"websocket"`
    Proxy           string `json:"proxy"`
}

type ConfigOption func(conf *ExchangeConfig)

func SetApiKey(apiKey string) ConfigOption {
    return func(conf *ExchangeConfig) {
        conf.ApiKey = apiKey
    }
}

func SetPassword(password string) ConfigOption {
    return func(conf *ExchangeConfig) {
        conf.Password = password
    }
}

func SetSecret(secret string) ConfigOption {
    return func(conf *ExchangeConfig) {
        conf.Secret = secret
    }
}

func SetTimeout(timeout int64) ConfigOption {
    return func(conf *ExchangeConfig) {
        conf.Timeout = timeout
    }
}

func SetEnableRateLimit(enableRateLimit bool) ConfigOption {
    return func(conf *ExchangeConfig) {
        conf.EnableRateLimit = enableRateLimit
    }
}

func SetTest(test bool) ConfigOption {
    return func(conf *ExchangeConfig) {
        conf.Test = test
    }
}

func SetVerbose(verbose bool) ConfigOption {
    return func(conf *ExchangeConfig) {
        conf.Verbose = verbose
    }
}

func SetWebsocket(websocket bool) ConfigOption {
    return func(conf *ExchangeConfig) {
        conf.Websocket = websocket
    }
}

type HasDescription struct {
	CancelAllOrders      bool `json:"cancelAllOrders"`
	CancelOrder          bool `json:"cancelOrder"`
	CancelOrders         bool `json:"cancelOrders"`
	CORS                 bool `json:"CORS"`
	CreateDepositAddress bool `json:"createDepositAddress"`
	CreateLimitOrder     bool `json:"createLimitOrder"`
	CreateMarketOrder    bool `json:"createMarketOrder"`
	CreateOrder          bool `json:"createOrder"`
	Deposit              bool `json:"deposit"`
	EditOrder            bool `json:"editOrder"`
	FetchBalance         bool `json:"fetchBalance"`
	FetchBidsAsks        bool `json:"fetchBidsAsks"`
	FetchClosedOrders    bool `json:"fetchClosedOrders"`
	FetchCurrencies      bool `json:"fetchCurrencies"`
	FetchDepositAddress  bool `json:"fetchDepositAddress"`
	FetchDeposits        bool `json:"fetchDeposits"`
	FetchFundingFees     bool `json:"fetchFundingFees"`
	FetchL2OrderBook     bool `json:"fetchL2OrderBook"`
	FetchLedger          bool `json:"fetchLedger"`
	FetchMarkets         bool `json:"fetchMarkets"`
	FetchMyTrades        bool `json:"fetchMyTrades"`
	FetchOHLCV           bool `json:"fetchOHLCV"`
	FetchOpenOrders      bool `json:"fetchOpenOrders"`
	FetchOrder           bool `json:"fetchOrder"`
	FetchOrderBook       bool `json:"fetchOrderBook"`
	FetchOrderBooks      bool `json:"fetchOrderBooks"`
	FetchOrders          bool `json:"fetchOrders"`
	FetchTicker          bool `json:"fetchTicker"`
	FetchTickers         bool `json:"fetchTickers"`
	FetchTrades          bool `json:"fetchTrades"`
	FetchTradingFee      bool `json:"fetchTradingFee"`
	FetchTradingFees     bool `json:"fetchTradingFees"`
	FetchTradingLimits   bool `json:"fetchTradingLimits"`
	FetchTransactions    bool `json:"fetchTransactions"`
	FetchWithdrawals     bool `json:"fetchWithdrawals"`
	PrivateApi           bool `json:"privateApi"`
	PublicApi            bool `json:"publicApi"`
	Withdraw             bool `json:"withdraw"`
}

type Balance struct {
    Free   float64  `json:"free"`
    Used   float64  `json:"used"`
    Total  float64  `json:"total"`
}

type OrderBook struct {
    Symbol      string         `json:"symbol"`
    Timestamp   int64          `json:"time"`
    Asks        [][2]float64   `json:"asks"`
    Bids        [][2]float64   `json:"bids"`
}

type OHLCV struct {
	Symbol    string    `json:"symbol"`    // 标
	Timestamp int64     `json:"timestamp"` // 时间
	Open      float64   `json:"open"`      // 开盘价
	High      float64   `json:"high"`      // 最高价
	Low       float64   `json:"low"`       // 最低价
	Close     float64   `json:"close"`     // 收盘价
	Volume    float64   `json:"volume"`    // 量
}

type Direction int32

const (
    BUY   Direction = iota
    SELL
    CLOSE_BUY
    CLOSE_SELL
)

func (d Direction) toString() string {
    switch d {
    case BUY:
        return "Buy"
    case SELL:
        return "Sell"
    case CLOSE_BUY:
        return "CloseBuy"
    case CLOSE_SELL:
        return "CloseSell"
    default:
        return "None"
    }
}

type OrderType int32

const (
    MARKET OrderType = iota
    LIMIT
    STOP_MARKET
    STOP_LIMIT
    TRAILING_STOP_MARKET
)

func (o OrderType) toString() string {
    switch o {
    case MARKET:
        return "Market"
    case LIMIT:
        return "Limit"
    case STOP_MARKET:
        return "StopMarket"
    case STOP_LIMIT:
        return "StopLimit"
    case TRAILING_STOP_MARKET:
        return "TrailingStopMarket"
    default:
        return "None"
    }
}

type Order struct {
	ID           string      `json:"id"`
	ClientOId    string      `json:"client_oid"`
	Symbol       string      `json:"symbol"`
	Timestamp    int64       `json:"time"`
	Price        float64     `json:"price"`
	StopPx       float64     `json:"stop_px"`
	Amount       float64     `json:"amount"`
	AvgPrice     float64     `json:"avg_price"`
	FilledAmount float64     `json:"filled_amount"`
	Direction    Direction   `json:"direction"`
	Type         OrderType   `json:"type"`
	PostOnly     bool        `json:"post_only"`
	ReduceOnly   bool        `json:"reduce_only"`
	Commission   float64     `json:"commission"`
	Pnl          float64     `json:"pnl"`
	UpdateTime   int64       `json:"update_time"`
	Status       OrderStatus `json:"status"`

	ActivatePrice string `json:"activatePrice"`
	PriceRate     string `json:"priceRate"`
	ClosePosition bool   `json:"closePosition"`
}

type OrderStatus int32

const (
	CREATED          OrderStatus = iota
	REJECTED
	NEW
	PARTIALLYFILLED
	FILLED
	CANCELPENDING
	CANCELED
	UNTRIGGERED
	TRIGGERED
)

func (s OrderStatus) toString() string {
	switch s {
	case CREATED:
		return "Created"
	case REJECTED:
		return "Rejected"
	case NEW:
		return "New"
	case PARTIALLYFILLED:
		return "PartiallyFilled"
	case FILLED:
		return "Filled"
	case CANCELPENDING:
		return "CancelPending"
	case CANCELED:
		return "Canceled"
	case UNTRIGGERED:
		return "Untriggered"
	case TRIGGERED:
		return "Triggered"
	default:
		return "None"
	}
}

type Trade struct {
	ID        string    `json:"id"`
	Direction Direction `json:"type"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	Ts        int64     `json:"ts"`
	Symbol    string    `json:"symbol,omitempty"`
}

type Position struct {
	Symbol    string         `json:"symbol"`
	OpenTime  int64          `json:"open_time"`
	OpenPrice float64        `json:"open_price"`
	Size      float64        `json:"size"`
	AvgPrice  float64        `json:"avg_price"`
	Profit    float64        `json:"profit"`

	MarginType       string  `json:"marginType"`
	IsAutoAddMargin  bool    `json:"isAutoAddMargin"`
	IsolatedMargin   float64 `json:"isolatedMargin"`
	Leverage         float64 `json:"leverage"`
	LiquidationPrice float64 `json:"liquidationPrice"`
	MarkPrice        float64 `json:"markPrice"`
	MaxNotionalValue float64 `json:"maxNotionalValue"`
	PositionSide     string  `json:"positionSide"`
}

func (p *Position) Side() Direction {
	if p.Size > 0 {
		return BUY
	} else if p.Size < 0 {
		return SELL
	}
	return BUY
}

func (p *Position) IsOpen() bool {
	return p.Size != 0
}

func (p *Position) IsLong() bool {
	return p.Size > 0
}

func (p *Position) IsShort() bool {
	return p.Size < 0
}
