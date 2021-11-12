package robinhood

import (
	"fmt"
	"net/url"
	"strings"
)

type OptionInstrument struct {
	ChainID        string          `json:"chain_id" url:"chain_id"`
	ChainSymbol    string          `json:"chain_symbol" url:"chain_symbol"`
	CreatedAt      string          `json:"created_at" url:"created_at"`
	ExpirationDate string          `json:"expiration_date" url:"expiration_date"`
	ID             string          `json:"id" url:"id"`
	IssueDate      string          `json:"issue_date" url:"issue_date"`
	MinTicks       OptionChainTick `json:"min_ticks" url:"min_ticks"`
	RHSTradability string          `json:"rhs_tradability" url:"rhs_tradability"` // e.g. "untradable"
	State          string          `json:"state" url:"state"`                     // e.g. "active"
	StrikePrice    string          `json:"strike_price" url:"strike_price"`
	Tradability    string          `json:"tradability" url:"tradability"` // e.g. "tradable"
	Type           string          `json:"type" url:"type"`               // e.g. "call"
	UpdatedAt      string          `json:"updated_at" url:"updated_at"`
	Url            string          `json:"url" url:"url"`
}

type OptionChain struct {
	ID                    string          `json:"id"`
	Symbol                string          `json:"symbol"`
	CanOpenPosition       bool            `json:"can_open_position"`
	CashComponent         *string         `json:"cash_component,omitempty"`
	ExpirationDates       []string        `json:"expiration_dates"`
	TradeValueMultiplier  string          `json:"trade_value_multiplier"`
	UnderlyingInstruments []Instrument    `json:"underlying_instruments"`
	MinTicks              OptionChainTick `json:"min_ticks"`
}

type OptionChainTick struct {
	AboveTick   string `json:"above_tick"`
	BelowTick   string `json:"below_tick"`
	CutoffPrice string `json:"cutoff_price"`
}

type OptionsPositionsResult struct {
	Results []OptionsPosition
}

type OptionsOrdersResult struct {
	Results []OptionsOrder
	Next    string
}

type OptionsLeg struct {
	ID             string `json:"id"`
	Position       string `json:"position"`      // "https://api.robinhood.com/options/positions/{id}/"
	PositionType   string `json:"position_type"` // e.g. "long"
	Option         string `json:"option"`        // "https://api.robinhood.com/options/instruments/{id}/"
	RatioQuantity  int    `json:"ratio_quantity"`
	ExpirationDate string `json:"expiration_date"` // e.g. "2020-03-27"
	StrikePrice    string `json:"strike_price"`
	OptionType     string `json:"option_type"` // e.g. "put"
}

type OptionsPosition struct {
	ID                       string       `json:"id"`
	Chain                    string       `json:"chain"` // "https://api.robinhood.com/options/chains/{id}/"
	Symbol                   string       `json:"symbol"`
	Strategy                 string       `json:"strategy"`
	AverageOpenPrice         string       `json:"average_open_price"`
	Legs                     []OptionsLeg `json:"legs"`
	Quantity                 string       `json:"quantity"`
	IntradayAverageOpenPrice string       `json:"intraday_average_open_price"`
	IntradayQuantity         string       `json:"intraday_quantity"`
	Direction                string       `json:"direction"` // e.g. "debit"/"credit"
	IntradayDirection        string       `json:"intraday_direction"`
	TradeValueMultiplier     string       `json:"trade_value_multiplier"`
	CreatedAt                string       `json:"created_at"` // e.g. "2020-03-20T19:49:03.862417Z"
	UpdatedAt                string       `json:"updated_at"`
}

type OptionsOrderLeg struct {
	//executions: []
	ID             string `json:"id"`
	Option         string `json:"option"`
	PositionEffect string `json:"position_effect"` // e.g. "close
	RatioQuantity  int    `json:"ratio_quantity"`
	Side           string `json:"sell"`
}

type OptionsOrder struct {
	CancelUrl         string            `json:"cancel_url"`
	CanceledQuantity  string            `json:"canceled_quantity"`
	ChainID           string            `json:"chain_id"`
	ChainSymbol       string            `json:"chain_symbol"`
	ClosingStrategy   string            `json:"closing_strategy"` // e.g. "long_call_spread" - set on "close" orders
	CreatedAt         string            `json:"created_at"`       // e.g. "2020-03-20T19:49:03.862417Z"
	Direction         string            `json:"direction"`        // e.g. "debit"/"credit"
	ID                string            `json:"id"`
	Legs              []OptionsOrderLeg `json:"legs"`
	OpeningStrategy   string            `json:"opening_strategy"` // set on "open" order
	PendingQuantity   string            `json:"pending_quantity"`
	Premium           string            `json:"premium"`
	Price             string            `json:"price"`
	ProcessedPremium  string            `json:"processed_premium"`
	ProcessedQuantity string            `json:"processed_quantity"`
	Quantity          string            `json:"quantity"`
	RefID             string            `json:"ref_id"`
	ResponseCategory  interface{}       `json:"response_category"` // TODO: type unknown
	State             string            `json:"state"`
	StopPrice         string            `json:"stop_price"`
	TimeInForce       string            `json:"time_in_force"` // e.g. "gtc"
	Trigger           string            `json:"trigger"`       // e.g. "immediate"
	Type              string            `json:"type"`
	UpdatedAt         string            `json:"updated_at"`
}

type ListOptionsInstrumentsParameters struct {
	OptionInstrument
	ExpirationDates string `json:"expiration_dates" url:"expiration_dates"`
}

// Get info for a particular set of parameters
func (c *Client) ListOptionsInstruments(
	parameters ListOptionsInstrumentsParameters,
) ([]*OptionInstrument, error) {
	url := Endpoint + "/options/instruments/"
	var req interface{} = &parameters
	var result []*OptionInstrument
	for url != "" {
		var resp struct {
			Results []*OptionInstrument
			Next    string
		}

		if err := c.getJSON(url, req, &resp); err != nil {
			return nil, err
		}

		req = nil
		result = append(result, resp.Results...)
		url = resp.Next
		if len(resp.Results) == 0 {
			break
		}
	}

	return result, nil
}

func (c *Client) GetAggregateOptionPositions(nonZero bool) (*OptionsPositionsResult, error) {
	urlValue := Endpoint + "/options/aggregate_positions/?"

	if nonZero {
		urlValue += "nonzero=true"
	}

	var resp OptionsPositionsResult
	err := c.getJSON(urlValue, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) GetOptionOrders(cursor *string) (*OptionsOrdersResult, error) {
	urlValue := Endpoint + "/options/orders/"

	if cursor != nil && len(*cursor) > 0 {
		urlValue += "?cursor=" + *cursor
	}

	var resp OptionsOrdersResult
	err := c.getJSON(urlValue, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) PlaceOptionsOrder(order OptionsOrder) (*OptionsOrder, error) {
	urlValue := Endpoint + "/options/orders/"

	var resp OptionsOrder
	err := c.postForm(urlValue, &order, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type getOptionsRequest struct {
	IDs string `url:"ids,omitempty"`
}

func (c *Client) GetOptionChains(chainIds []string) ([]*OptionChain, error) {
	urlValue := Endpoint + "/options/chains/?"

	var result []*OptionChain
	for urlValue != "" {
		var resp struct {
			Results []*OptionChain `json:"results"`
			Next    string         `json:"next"`
		}

		err := c.getJSON(urlValue, getOptionsRequest{
			IDs: strings.Join(chainIds, ","),
		}, &resp)
		if err != nil {
			return nil, err
		}

		result = append(result, resp.Results...)

		if resp.Next == "" || resp.Next == urlValue {
			break
		}

		urlValue = resp.Next
	}

	return result, nil
}

func ParseChainID(instrumentURL string) (string, error) {
	urlParsed, err := url.Parse(instrumentURL)
	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.Trim(urlParsed.Path, "/"), "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid instrument URL: %v", instrumentURL)
	}

	return parts[2], nil
}
