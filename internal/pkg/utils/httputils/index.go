package httputils

import (
	"encoding/json"
	"io/ioutil"
	"kyberswap_user_monitor/pkg/context"
	"net/http"
)
type Token struct {
	Address  string
	Name     string
	Symbol   string
	Decimals int
	CkgID    string `json:"cgkId"`
}
type TokenInfo struct {
	Address  string
	Name     string
	Symbol   string
	Decimals int
	Price    float64
	Type     string
	CkgID    string   `json:"cgkId"`
	Tokens   []*Token `json:"tokens"`
}


func GetTokenInfo(ctx context.Context, api, address string) (*TokenInfo, error) {
	api += "/tokens"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		ctx.Errorf("failed to prepare client request, err: %v", err)
		return nil, err
	}
	q := req.URL.Query()
	q.Add("ids", address)
	req.URL.RawQuery = q.Encode()

	info := make(map[string]*TokenInfo)
	if err := Process(ctx, req, &info); err != nil {
		ctx.Errorf("failed to call price api, err: %v", err)
		return nil, err
	}

	return info[address], nil
}


func Process(ctx context.Context, req *http.Request, dest interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.Errorf("failed to post request, err: %v", err)
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.Errorf("failed to read response body, err: %v", err)
		return err
	}
	ctx.Debugf("price http response: %v", string(respBody))

	err = json.Unmarshal(respBody, dest)
	if err != nil {
		ctx.Errorf("failed to unmarshal response data, err: %v", err)
		return err
	}

	return nil
}

