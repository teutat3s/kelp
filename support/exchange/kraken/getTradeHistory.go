package kraken

import (
	"strconv"

	"github.com/lightyeario/kelp/api"
	"github.com/lightyeario/kelp/model"
)

// GetTradeHistory impl.
func (k krakenExchange) GetTradeHistory(maybeCursorStart interface{}, maybeCursorEnd interface{}) (*api.TradeHistoryResult, error) {
	var mcs *int64
	if maybeCursorStart != nil {
		i := maybeCursorStart.(int64)
		mcs = &i
	}

	var mce *int64
	if maybeCursorEnd != nil {
		i := maybeCursorEnd.(int64)
		mce = &i
	}

	return k.getTradeHistory(mcs, mce)
}

func (k krakenExchange) getTradeHistory(maybeCursorStart *int64, maybeCursorEnd *int64) (*api.TradeHistoryResult, error) {
	input := map[string]string{}
	if maybeCursorStart != nil {
		input["start"] = strconv.FormatInt(*maybeCursorStart, 10)
	}
	if maybeCursorEnd != nil {
		input["end"] = strconv.FormatInt(*maybeCursorEnd, 10)
	}

	resp, e := k.api.Query("TradesHistory", input)
	if e != nil {
		return nil, e
	}
	krakenResp := resp.(map[string]interface{})
	krakenTrades := krakenResp["trades"].(map[string]interface{})

	res := api.TradeHistoryResult{Trades: []model.Trade{}}
	for _, v := range krakenTrades {
		m := v.(map[string]interface{})
		_txid := m["ordertxid"].(string)
		_time := m["time"].(float64)
		ts := model.MakeTimestamp(int64(_time))
		_type := m["type"].(string)
		_ordertype := m["ordertype"].(string)
		_price := m["price"].(string)
		_vol := m["vol"].(string)
		_cost := m["cost"].(string)
		_fee := m["fee"].(string)
		_pair := m["pair"].(string)
		pair, e := model.TradingPairFromString(4, k.assetConverter, _pair)
		if e != nil {
			return nil, e
		}

		res.Trades = append(res.Trades, model.Trade{
			Order: model.Order{
				Pair:        pair,
				OrderAction: model.OrderActionFromString(_type),
				OrderType:   model.OrderTypeFromString(_ordertype),
				Price:       model.MustFromString(_price, k.precision),
				Volume:      model.MustFromString(_vol, k.precision),
				Timestamp:   ts,
			},
			TransactionID: model.MakeTransactionID(_txid),
			Cost:          model.MustFromString(_cost, k.precision),
			Fee:           model.MustFromString(_fee, k.precision),
		})
	}
	return &res, nil
}