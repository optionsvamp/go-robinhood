package robinhood

import (
	"strings"
)

type instrumentListRequest struct {
	Instruments string `url:"instruments,omitempty"`
}

type OptionsInstrumentListResponse struct {
	Results []OptionsQuote `json:"results" url:"results"`
}

type OptionsQuote struct {
	AskPrice              string      `json:"ask_price" url:"ask_price"`
	AskSize               int         `json:"ask_size" url:"ask_size"`
	BidPrice              string      `json:"bid_price" url:"bid_price"`
	BidSize               int         `json:"bid_size" url:"bid_size"`
	ImpliedVolatility     string      `json:"implied_volatility" url:"implied_volatility"`
	Vega                  string      `json:"vega" url:"vega"`
	MarkPrice             string      `json:"mark_price" url:"mark_price"`
	PreviousCloseDate     string      `json:"previous_close_date" url:"previous_close_date"`
	ChanceOfProfitShort   string      `json:"chance_of_profit_short" url:"chance_of_profit_short"`
	Delta                 string      `json:"delta" url:"delta"`
	Rho                   string      `json:"rho" url:"rho"`
	OpenInterest          float64     `json:"open_interest" url:"open_interest"`
	Theta                 string      `json:"theta" url:"theta"`
	LowPrice              interface{} // TODO: type
	Gamma                 string      `json:"gamma" url:"gamma"`
	LowFillRateBuyPrice   string      `json:"low_fill_rate_buy_price" url:"low_fill_rate_buy_price"`
	BreakEvenPrice        string      `json:"break_even_price" url:"break_even_price"`
	LastTradeSize         int         `json:"last_trade_size" url:"last_trade_size"`
	PreviousClosePrice    string      `json:"previous_close_price" url:"previous_close_price"`
	HighFillRateSellPrice string      `json:"high_fill_rate_sell_price" url:"high_fill_rate_sell_price"`
	HighFillRateBuyPrice  string      `json:"high_fill_rate_buy_price" url:"high_fill_rate_buy_price"`
	AdjustedMarkPrice     string      `json:"adjusted_mark_price" url:"adjusted_mark_price"`
	HighPrice             interface{} // TODO: type
	Instrument            string      `json:"instrument" url:"instrument"`
	LowFillRateSellPrice  string      `json:"low_fill_rate_sell_price" url:"low_fill_rate_sell_price"`
	LastTradePrice        string      `json:"last_trade_price" url:"last_trade_price"`
	Volume                int         `json:"volume" url:"volume"`
	ChanceOfProfitLong    string      `json:"chance_of_profit_long" url:"chance_of_profit_long"`
	// TODO: MORE
}

type OptionsStrategyQuote struct {
	AdjustedMarkPrice  string                    `json:"adjusted_mark_price" url:"adjusted_mark_price"`
	AskPrice           string                    `json:"ask_price" url:"ask_price"`
	BidPrice           string                    `json:"bid_price" url:"bid_price"`
	Legs               []OptionsStrategyQuoteLeg `json:"legs"`
	MarkPrice          string                    `json:"mark_price" url:"mark_price"`
	PreviousCloseDate  string                    `json:"previous_close_date" url:"previous_close_date"`
	PreviousClosePrice string                    `json:"previous_close_price" url:"previous_close_price"`
}

type OptionsStrategyQuoteLeg struct {
	ID    string `json:"id"`
	Ratio string `json:"ratio"` // e.g. 1.0
	Type  string `json:"type"`  // e.g. "short"/"long"
}

type Quote struct {
	AskPrice                    string `json:"ask_price" url:"ask_price"`
	AskSize                     int    `json:"ask_size" url:"ask_size"`
	BidPrice                    string `json:"bid_price" url:"bid_price"`
	BidSize                     int    `json:"bid_size" url:"bid_size"`
	LastTradePrice              string `json:"last_trade_price" url:"last_trade_price"`
	LastExtendedHoursTradePrice string `json:"last_extended_hours_trade_price" url:"last_extended_hours_trade_price"`
	PreviousClose               string `json:"previous_close" url:"previous_close"`
	AdjustedPreviousClose       string `json:"adjusted_previous_close" url:"adjusted_previous_close"`
	PreviousCloseDate           string `json:"previous_close_date" url:"previous_close_date"`
	Symbol                      string `json:"symbol" url:"symbol"`
	TradingHalted               bool   `json:"trading_halted" url:"trading_halted"`
	HasTraded                   bool   `json:"has_traded" url:"has_traded"`
	LastTradePriceSource        string `json:"last_trade_price_source" url:"last_trade_price_source"`
	UpdatedAt                   string `json:"updated_at" url:"updated_at"`
	Instrument                  string `json:"instrument" url:"instrument"`
}

type Historicals struct {
	Quote        string
	Symbol       string
	Interval     string
	Span         string
	Instrument   string
	Historicals  []Historical `json:"historicals,omitempty"`
	InstrumentID string
}

type Historical struct {
	BeginsAt     string `json:"begins_at"'`
	OpenPrice    string `json:"open_price"'`
	ClosePrice   string `json:"close_price"`
	HighPrice    string `json:"high_price"`
	LowPrice     string `json:"low_price"`
	Volume       int64
	Session      string
	Interpolated bool
}

// :path: /marketdata/options/?instruments=
// note that a limit of 80 applies
func (c *Client) ListOptionsMarketDataByInstrumentURLList(urlList []string) ([]OptionsQuote, error) {
	urlVal := Endpoint + "/marketdata/options/?"
	resp := OptionsInstrumentListResponse{}
	err := c.getJSON(urlVal, instrumentListRequest{Instruments: strings.Join(urlList, ",")}, &resp)
	return resp.Results, err
}

func (c *Client) GetQuoteForSymbol(symbol string) (*Quote, error) {
	urlVal := Endpoint + "/marketdata/quotes/" + symbol + "/?bounds=trading" // TODO: params?
	resp := Quote{}
	err := c.getJSON(urlVal, nil, &resp)
	return &resp, err
}

func (c *Client) GetHistoricalMarketDataForSymbol(symbol, interval, span string) (Historicals, error) {
	urlVal := Endpoint + "/marketdata/historicals/" + symbol + "/?bounds=regular&include_inactive=true&interval=" + interval + "&span=" + span // TODO: params?
	resp := Historicals{}
	err := c.getJSON(urlVal, nil, &resp)
	return resp, err
}

// https://api.robinhood.com/marketdata/historicals/BBBY/?bounds=regular&include_inactive=true&interval=1day&span=week

func (c *Client) GetOptionsStrategyQuote(instruments []string, ratios []string, types []string) (*OptionsStrategyQuote, error) {
	urlValue := Endpoint + "/marketdata/options/strategy/quotes/?instruments=" + strings.Join(instruments, ",") +
		"&ratios=" + strings.Join(ratios, ",") +
		"&types=" + strings.Join(types, ",")

	var resp OptionsStrategyQuote
	err := c.getJSON(urlValue, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
