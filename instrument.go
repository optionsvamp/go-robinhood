package robinhood

import (
	"fmt"
	"net/url"
	"strings"
)

type Instrument struct {
	MarginInitialRatio float64 `json:"margin_initial_ratio,string"`
	RHSTradability     string  `json:"rhs_tradability"`
	ID                 string
	MarketURL          string   `json:"market"`
	SimpleName         string   `json:"simple_name"`
	MinTickSize        *float64 `json:"min_tick_size,string"`
	MaintenanceRatio   float64  `json:"maintenance_ratio,string"`
	Tradability        string
	State              string
	Type               string
	Tradeable          bool
	FundamentalsURL    string `json:"fundamentals"`
	QuoteURL           string `json:"quote"`
	Symbol             string
	ChainSymbol        string  `json:"chain_symbol"`
	DayTradeRatio      float64 `json:"day_trade_ratio,string"`
	Name               string
	TradableChainID    *string `json:"tradable_chain_id"`
	SplitsURL          string  `json:"splits"`
	URL                string
	BloombergUnique    string `json:"bloomberg_unique"`
	ListDate           string `json:"list_date"`
	Country            string
}

// ListAllInstruments will take a long time to list all known instruments
func (c *Client) ListAllInstruments() ([]*Instrument, error) {
	url := Endpoint + "/instruments/"
	var result []*Instrument
	for url != "" {
		var resp struct {
			Results []*Instrument
			Next    string
		}

		if err := c.getJSON(url, nil, &resp); err != nil {
			return nil, err
		}

		result = append(result, resp.Results...)
		url = resp.Next
		if len(resp.Results) == 0 {
			break
		}
	}

	return result, nil
}

type getInstrumentsRequest struct {
	Symbol string `url:",omitempty"`
}

// ListInstrumentsByIDList will fetch a list of instruments in a single query
func (c *Client) ListInstrumentsByIDList(idList []string) ([]*Instrument, error) {
	url := Endpoint + "/instruments/?ids=" + strings.Join(idList, ",")
	var resp struct {
		Results []*Instrument
		Next    string
	}

	err := c.getJSON(url, nil, &resp)
	return resp.Results, err
}

// Get info for a particular symbol.
func (c *Client) ListInstrumentsForSymbol(symbol string) ([]*Instrument, error) {
	url := Endpoint + "/instruments/"
	req := &getInstrumentsRequest{symbol}
	var result []*Instrument
	for url != "" {
		var resp struct {
			Results []*Instrument
			Next    string
		}

		if err := c.getJSON(url, req, &resp); err != nil {
			return nil, err
		}

		result = append(result, resp.Results...)
		url = resp.Next
		if len(resp.Results) == 0 {
			break
		}
	}

	return result, nil
}

// GetInstrument fetches an instrument by its ID
func (c *Client) GetInstrument(id string) (*Instrument, error) {
	instruments, err := c.ListInstrumentsByIDList([]string{id})
	if err != nil {
		return nil, err
	}
	if len(instruments) == 0 {
		return nil, nil
	}
	return instruments[0], nil
}

// Helper function to extract the instrument ID from an instrument URL.
func ParseInstrumentID(instrumentURL string) (string, error) {
	urlParsed, err := url.Parse(instrumentURL)
	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.Trim(urlParsed.Path, "/"), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid instrument URL: %v", instrumentURL)
	}

	return parts[1], nil
}

func GetInstrumentURL(instrumentID string) string {
	return Endpoint + "/instruments/" + instrumentID + "/"
}
