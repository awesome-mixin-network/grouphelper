package main

import (
	"encoding/json"
	"fmt"
	"errors"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/bot-api-go-client/config"
	"context"
	"time"
)

type Snapshot struct {
	SnapshotId string       `json:"snapshot_id"`
	Amount     string       `json:"amount"`
	Asset      Asset `json:"asset"`
	CreatedAt  time.Time    `json:"created_at"`
	TraceId    string `json:"trace_id"`
	UserId     string `json:"user_id"`
	OpponentId string `json:"opponent_id"`
	Data       string `json:"data"`
}

type Asset struct {
	AssetId string `json:"asset_id"`
	ChainId string `json:"chain_id"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Logo    string `json:"icon_url"`
}


func requestMixinNetwork(ctx context.Context, checkpoint time.Time, limit int) ([]*Snapshot, error) {
	uri := fmt.Sprintf("/network/snapshots?offset=%s&order=ASC&limit=%d", checkpoint.Format(time.RFC3339Nano), limit)
	token, err := bot.SignAuthenticationToken(config.GetConfig().ClientID, config.GetConfig().SessionID, config.GetConfig().PrivateKey, "GET", uri, "")
	if err != nil {
		return nil, err
	}
	body, err := bot.Request(ctx, "GET", uri, nil, token)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data  []*Snapshot `json:"data"`
		Error string      `json:"error"`
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	return resp.Data, nil
}