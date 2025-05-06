package model

import (
	"time"
)

type OrderType string
type OrderSide string
type KlinePeriod string

type OrderStatus int

func (s OrderStatus) String() string {
	switch s {
	case 1:
		return "pending"
	case 2:
		return "finished"
	case 3:
		return "canceled"
	case 4:
		return "part-finished"
	}
	return "unknown-status"
}

// OptionParameter is api option parameter
type OptionParameter struct {
	Key   string
	Value string
}

func (OptionParameter) OrderClientID(cid string) OptionParameter {
	return OptionParameter{
		Key:   Order_Client_ID__Opt_Key, // 内部根据Order_Client_ID__Opt_Key来做适配
		Value: cid,
	}
}

type CurrencyPair struct {
	Symbol               string  `json:"symbol,omitempty"`          //交易对
	BaseSymbol           string  `json:"base_symbol,omitempty"`     //币种
	QuoteSymbol          string  `json:"quote_symbol,omitempty"`    //交易区：usdt/usdc/btc ...
	PricePrecision       int     `json:"price_precision,omitempty"` //价格小数点位数
	QtyPrecision         int     `json:"qty_precision,omitempty"`   //数量小数点位数
	MinQty               float64 `json:"min_qty,omitempty"`
	MaxQty               float64 `json:"max_qty,omitempty"`
	MarketQty            float64 `json:"market_qty,omitempty"`
	ContractVal          float64 `json:"contract_val,omitempty"`           //1张合约价值
	ContractValCurrency  string  `json:"contract_val_currency,omitempty"`  //合约面值计价币
	SettlementCurrency   string  `json:"settlement_currency,omitempty"`    //结算币
	ContractAlias        string  `json:"contract_alias,omitempty"`         //交割合约alias
	ContractDeliveryDate int64   `json:"contract_delivery_date,omitempty"` //合约交割日期
}

//func (pair CurrencyPair) String() string {
//	return pair.Symbol
//}

//type FuturesCurrencyPair struct {
//	CurrencyPair
//	DeliveryDate int64   //结算日期
//	OnboardDate  int64   //上线日期
//	MarginAsset  float64 //保证金资产
//}

type Ticker struct {
	Pair      CurrencyPair `json:"pair"`
	Last      float64      `json:"l"`
	Buy       float64      `json:"b"`
	Sell      float64      `json:"s"`
	High      float64      `json:"h"`
	Low       float64      `json:"lw"`
	Vol       float64      `json:"v"`
	Percent   float64      `json:"percent"`
	Timestamp int64        `json:"t"`
}

type DepthItem struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

type DepthItems []DepthItem

func (dr DepthItems) Len() int {
	return len(dr)
}

func (dr DepthItems) Swap(i, j int) {
	dr[i], dr[j] = dr[j], dr[i]
}

func (dr DepthItems) Less(i, j int) bool {
	return dr[i].Price < dr[j].Price
}

type Depth struct {
	Pair  CurrencyPair `json:"pair"`
	UTime time.Time    `json:"ut"`
	Asks  DepthItems   `json:"asks"`
	Bids  DepthItems   `json:"bids"`
}

type Kline struct {
	Pair      CurrencyPair `json:"pair"`
	Timestamp int64        `json:"t"`
	Open      float64      `json:"o"`
	Close     float64      `json:"s"`
	High      float64      `json:"h"`
	Low       float64      `json:"l"`
	Vol       float64      `json:"v"`
}

type Order struct {
	Pair        CurrencyPair `json:"pair,omitempty"`
	Id          string       `json:"id,omitempty"`       //订单ID
	CId         string       `json:"c_id,omitempty"`     //客户端自定义ID
	Side        OrderSide    `json:"side,omitempty"`     //交易方向: sell,buy
	OrderTy     OrderType    `json:"order_ty,omitempty"` //类型: limit , market , ...
	Status      OrderStatus  `json:"status,omitempty"`   //状态
	Price       float64      `json:"price,omitempty"`
	Qty         float64      `json:"qty,omitempty"`
	ExecutedQty float64      `json:"executed_qty,omitempty"`
	PriceAvg    float64      `json:"price_avg,omitempty"`
	Fee         float64      `json:"fee,omitempty"`
	FeeCcy      string       `json:"fee_ccy,omitempty"` //收取交易手续费币种
	CreatedAt   int64        `json:"created_at,omitempty"`
	FinishedAt  int64        `json:"finished_at,omitempty"` //订单完成时间
	CanceledAt  int64        `json:"canceled_at,omitempty"`
}

type Account struct {
	Coin             string  `json:"coin,omitempty"`
	Balance          float64 `json:"balance,omitempty"`
	AvailableBalance float64 `json:"available_balance,omitempty"`
	FrozenBalance    float64 `json:"frozen_balance,omitempty"`
}

type FuturesPosition struct {
	Pair     CurrencyPair `json:"pair,omitempty"`
	PosSide  OrderSide    `json:"pos_side,omitempty"`  //开仓方向
	Qty      float64      `json:"qty,omitempty"`       // 持仓数量
	AvailQty float64      `json:"avail_qty,omitempty"` //可平仓数量
	AvgPx    float64      `json:"avg_px,omitempty"`    //开仓均价
	LiqPx    float64      `json:"liq_px,omitempty"`    // 爆仓价格
	Upl      float64      `json:"upl,omitempty"`       //盈亏
	UplRatio float64      `json:"upl_ratio,omitempty"` // 盈亏率
	Lever    float64      `json:"lever,omitempty"`     //杠杆倍数
}

type FuturesAccount struct {
	Coin      string  `json:"coin,omitempty"` //币种
	Eq        float64 `json:"eq,omitempty"`   //总权益
	AvailEq   float64 `json:"avail_eq,omitempty"`
	FrozenBal float64 `json:"frozen_bal,omitempty"`
	MgnRatio  float64 `json:"mgn_ratio,omitempty"`
	Upl       float64 `json:"upl,omitempty"`
	RiskRate  float64 `json:"risk_rate,omitempty"`
}

type FundingRate struct {
	Symbol string  `json:"symbol"`
	Rate   float64 `json:"rate"`
	Tm     int64   `json:"tm"` //资金费收取时间
}

type AssetValuation struct {
	TotalBal       float64 `json:"total_bal"`       // 总资产折合，单位USD
	TotalEquity    float64 `json:"total_equity"`    // 净资产折合，单位USD
	IsolatedEquity float64 `json:"isolated_equity"` // 逐仓仓位权益，单位USD
	UpdateTime     int64   `json:"update_time"`     // 更新时间，单位毫秒
	Ccy            string  `json:"ccy"`             // 币种，如USD
}

type AssetBalance struct {
	AvailBal   float64 `json:"avail_bal"`   // 可用余额
	Bal        float64 `json:"bal"`         // 余额
	Ccy        string  `json:"ccy"`         // 币种
	FrozenBal  float64 `json:"frozen_bal"`  // 冻结（不可用）
	UpdateTime int64   `json:"update_time"` // 更新时间，单位毫秒
}

// AssetBill 资产账单明细
type AssetBill struct {
	BillId  string  `json:"bill_id"`  // 账单ID
	Ccy     string  `json:"ccy"`      // 币种
	Type    string  `json:"type"`     // 账单类型
	Amount  float64 `json:"amount"`   // 金额
	Ts      int64   `json:"ts"`       // 创建时间，Unix时间戳的毫秒数格式
	SubType string  `json:"sub_type"` // 账单子类型
	InstId  string  `json:"inst_id"`  // 产品ID，如BTC-USDT-SWAP
	FromAct string  `json:"from_act"` // 转出账户
	ToAct   string  `json:"to_act"`   // 转入账户
	Notes   string  `json:"notes"`    // 备注
}

// AssetCurrency 币种资产信息
type AssetCurrency struct {
	Ccy             string  `json:"ccy"`               // 币种名称，如BTC
	Name            string  `json:"name"`              // 币种中文名称，不显示则无中文名称
	Chain           string  `json:"chain"`             // 链信息
	CanDep          bool    `json:"can_dep"`           // 是否可充值，true 或 false
	CanWd           bool    `json:"can_wd"`            // 是否可提币，true 或 false
	CanInternal     bool    `json:"can_internal"`      // 是否可内部转账，true 或 false
	MinDep          float64 `json:"min_dep"`           // 最小充值量
	MinWd           float64 `json:"min_wd"`            // 最小提币量
	MaxWd           float64 `json:"max_wd"`            // 最大提币量
	WdFee           float64 `json:"wd_fee"`            // 提币固定手续费
	WdAll           bool    `json:"wd_all"`            // 是否可全部提币，true 或 false
	DepQuotaFixed   float64 `json:"dep_quota_fixed"`   // 充值固定限额
	DepQuoteDynamic float64 `json:"dep_quota_dynamic"` // 充值动态限额
}
