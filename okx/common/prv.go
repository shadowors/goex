package common

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shadowors/goex/v2/httpcli"
	"github.com/shadowors/goex/v2/logger"
	"github.com/shadowors/goex/v2/model"
	"github.com/shadowors/goex/v2/options"
	"github.com/shadowors/goex/v2/util"
)

type Prv struct {
	*OKxV5
	apiOpts options.ApiOptions
}

func (prv *Prv) GetAccount(coin string) (map[string]model.Account, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetAccountUri)
	params := url.Values{}
	params.Set("ccy", coin)
	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}
	acc, err := prv.UnmarshalOpts.GetAccountResponseUnmarshaler(data)
	return acc, responseBody, err
}

func (prv *Prv) CreateOrder(pair model.CurrencyPair, qty, price float64, side model.OrderSide, orderTy model.OrderType, opts ...model.OptionParameter) (*model.Order, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.NewOrderUri)
	params := url.Values{}

	params.Set("instId", pair.Symbol)
	//params.Set("tdMode", "cash")
	//params.Set("posSide", "")
	params.Set("ordType", adaptOrderTypeToSym(orderTy))
	params.Set("px", util.FloatToString(price, pair.PricePrecision))
	params.Set("sz", util.FloatToString(qty, pair.QtyPrecision))

	side2, posSide := adaptOrderSideToSym(side)
	params.Set("side", side2)
	if posSide != "" {
		params.Set("posSide", posSide)
	}

	util.MergeOptionParams(&params, opts...)
	AdaptOrderClientIDOptionParameter(&params)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if err != nil {
		logger.Errorf("[CreateOrder] response body =%s", string(responseBody))
		return nil, responseBody, err
	}

	ord, err := prv.UnmarshalOpts.CreateOrderResponseUnmarshaler(data)
	if err != nil {
		return nil, responseBody, err
	}

	ord.Pair = pair
	ord.Price = price
	ord.Qty = qty
	ord.Side = side
	ord.OrderTy = orderTy
	ord.Status = model.OrderStatus_Pending

	return ord, responseBody, err
}

func (prv *Prv) GetOrderInfo(pair model.CurrencyPair, id string, opt ...model.OptionParameter) (*model.Order, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetOrderUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	params.Set("ordId", id)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	ord, err := prv.UnmarshalOpts.GetOrderInfoResponseUnmarshaler(data[1 : len(data)-1])
	if err != nil {
		return nil, responseBody, err
	}

	ord.Pair = pair

	return ord, responseBody, nil
}

func (prv *Prv) GetPendingOrders(pair model.CurrencyPair, opt ...model.OptionParameter) ([]model.Order, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetPendingOrdersUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	orders, err := prv.UnmarshalOpts.GetPendingOrdersResponseUnmarshaler(data)
	return orders, responseBody, err
}

func (prv *Prv) GetHistoryOrders(pair model.CurrencyPair, opt ...model.OptionParameter) ([]model.Order, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetHistoryOrdersUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	params.Set("limit", "50")

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	orders, err := prv.UnmarshalOpts.GetHistoryOrdersResponseUnmarshaler(data)
	return orders, responseBody, err
}

func (prv *Prv) CancelOrder(pair model.CurrencyPair, id string, opt ...model.OptionParameter) ([]byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.CancelOrderUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	params.Set("ordId", id)
	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if data != nil && len(data) > 0 {
		return responseBody, prv.UnmarshalOpts.CancelOrderResponseUnmarshaler(data)
	}

	return responseBody, err
}

// GetAssetValuation 获取资产估值
// currency: 币种(USD、USDT、BTC等)，如果为空字符串，则默认使用账户设置的币种
func (prv *Prv) GetAssetValuation(currency string) (*model.AssetValuation, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetAssetValuationUri)
	params := url.Values{}
	if currency != "" {
		params.Set("ccy", currency)
	}

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	valuation, err := prv.UnmarshalOpts.GetAssetValuationResponseUnmarshaler(data)
	return valuation, responseBody, err
}

// GetAssetBalances 获取资产余额
// currency: 币种(BTC等)，如果为空字符串，则获取所有币种余额
func (prv *Prv) GetAssetBalances(currency string) (map[string]model.AssetBalance, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetAssetBalancesUri)
	params := url.Values{}
	if currency != "" {
		params.Set("ccy", currency)
	}

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	balances, err := prv.UnmarshalOpts.GetAssetBalancesResponseUnmarshaler(data)
	return balances, responseBody, err
}

// GetAssetBills 获取账户账单明细
// params:
//   - currency: 币种，如BTC，不填则返回所有币种
//   - type: 账单类型，1:充值，2:提现，13:买入，14:卖出，不填则返回所有类型
//   - startTime: 开始时间，Unix时间戳的毫秒数格式
//   - endTime: 结束时间，Unix时间戳的毫秒数格式
//   - limit: 分页返回的结果集数量，默认为100，最大为100
//   - before: 请求此id之前（更旧的数据）的分页内容
//   - after: 请求此id之后（更新的数据）的分页内容
func (prv *Prv) GetAssetBills(params url.Values) ([]model.AssetBill, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetAssetBillsUri)

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	bills, err := prv.UnmarshalOpts.GetAssetBillsResponseUnmarshaler(data)
	return bills, responseBody, err
}

// GetAssetCurrencies 获取所有币种资产信息
// currency: 币种，如BTC，不填则返回所有币种
func (prv *Prv) GetAssetCurrencies(currency string) ([]model.AssetCurrency, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetAssetCurrenciesUri)
	params := url.Values{}
	if currency != "" {
		params.Set("ccy", currency)
	}

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	currencies, err := prv.UnmarshalOpts.GetAssetCurrenciesResponseUnmarshaler(data)
	return currencies, responseBody, err
}

func (prv *Prv) DoSignParam(httpMethod, apiUri, apiSecret, reqBody string) (signStr, timestamp string) {
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z") //iso time style
	payload := fmt.Sprintf("%s%s%s%s", timestamp, strings.ToUpper(httpMethod), apiUri, reqBody)
	signStr, _ = util.HmacSHA256Base64Sign(apiSecret, payload)
	return
}

func (prv *Prv) DoAuthRequest(httpMethod, reqUrl string, params *url.Values, headers map[string]string) ([]byte, []byte, error) {
	var (
		reqBodyStr string
		reqUri     string
	)

	if http.MethodGet == httpMethod {
		reqUrl += "?" + params.Encode()
	}

	if http.MethodPost == httpMethod {
		params.Set("tag", "86d4a3bf87bcBCDE")
		reqBody, _ := util.ValuesToJson(*params)
		reqBodyStr = string(reqBody)
	}

	_url, _ := url.Parse(reqUrl)
	reqUri = _url.RequestURI()
	signStr, timestamp := prv.DoSignParam(httpMethod, reqUri, prv.apiOpts.Secret, reqBodyStr)
	logger.Debugf("[DoAuthRequest] sign base64: %s, timestamp: %s", signStr, timestamp)

	headers = map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
		//"Accept":               "application/json",
		"OK-ACCESS-KEY":        prv.apiOpts.Key,
		"OK-ACCESS-PASSPHRASE": prv.apiOpts.Passphrase,
		"OK-ACCESS-SIGN":       signStr,
		"OK-ACCESS-TIMESTAMP":  timestamp}

	respBody, err := httpcli.Cli.DoRequest(httpMethod, reqUrl, reqBodyStr, headers)
	if err != nil {
		return nil, respBody, err
	}
	logger.Debugf("[DoAuthRequest] response body: %s", string(respBody))

	var baseResp BaseResp
	err = prv.OKxV5.UnmarshalOpts.ResponseUnmarshaler(respBody, &baseResp)
	if err != nil {
		return nil, respBody, err
	}

	if baseResp.Code != 0 {
		var errData []ErrorResponseData
		err = prv.OKxV5.UnmarshalOpts.ResponseUnmarshaler(baseResp.Data, &errData)
		if err != nil {
			logger.Errorf("unmarshal error data error: %s", err.Error())
			return nil, respBody, errors.New(string(respBody))
		}
		if len(errData) > 0 {
			return nil, respBody, errors.New(errData[0].SMsg)
		}
		return nil, respBody, errors.New(baseResp.Msg)
	} // error response process

	return baseResp.Data, respBody, nil
}

func NewPrvApi(opts ...options.ApiOption) *Prv {
	var api = new(Prv)
	for _, opt := range opts {
		opt(&api.apiOpts)
	}
	return api
}
