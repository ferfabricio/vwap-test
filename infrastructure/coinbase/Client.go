package coinbase

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	errorType      = "error"
	matchesChannel = "matches"
	matchType      = "match"
	subscribeType  = "subscribe"
	url            = "wss://ws-feed.exchange.coinbase.com"
)

type ChannelSubscription struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

type SubscribeMessage struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type SubscriptionResult struct {
	Type     string                `json:"type"`
	Channels []ChannelSubscription `json:"channels"`
}
type GenericResult struct {
	Type     string                `json:"type"`
	Message  string                `json:"message"`
	Reason   string                `json:"reason"`
	Channels []ChannelSubscription `json:"channels"`
}

type MatchMessage struct {
	Type         string    `json:"type"`
	TradeId      int32     `json:"trade_id"`
	MakerOrderId uuid.UUID `json:"maker_order_id"`
	TakerOrderId uuid.UUID `json:"taker_order_id"`
	Side         string    `json:"side"`
	Size         string    `json:"size"`
	Price        string    `json:"price"`
	ProductId    string    `json:"product_id"`
	Sequence     int64     `json:"sequence"`
	Time         time.Time `json:"time"`
	Message      string    `json:"message"`
}

type Client struct {
	conn *websocket.Conn
}

func (c Client) Configure(p []string) error {
	err := c.conn.WriteJSON(SubscribeMessage{
		Type:       subscribeType,
		ProductIds: p,
		Channels:   []string{matchesChannel},
	})
	if err != nil {
		log.Println(err)
		return err
	}

	r := GenericResult{}
	if err = c.conn.ReadJSON(&r); err != nil {
		log.Println(err)
		return errors.New("error to configure Coinbase WS subscription")
	}

	if r.Type == "error" {
		err = errors.New(r.Message)
	}

	return err
}

func (c Client) GetData(ch chan MatchMessage) {
	defer c.conn.Close()
	for {
		m := MatchMessage{}
		err := c.conn.ReadJSON(&m)
		if err != nil {
			log.Println(err)
			close(ch)
			break
		}

		ch <- m
	}
}

func NewClient() (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}
