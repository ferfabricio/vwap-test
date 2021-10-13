package coinbase

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

const (
	url            = "wss://ws-feed.exchange.coinbase.com"
	matchesChannel = "matches"
	subscribeType  = "subscribe"
)

type SubscribeMessage struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type ChannelSubscription struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
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

// type orderPairMessage struct {
// 	Type         string    `json:"type"`
// 	TradeId      int32     `json:"trade_id"`
// 	MakerOrderId uuid.UUID `json:"maker_order_id"`
// 	TakerOrderId uuid.UUID `json:"taker_order_id"`
// 	Side         string    `json:"side"`
// 	Size         string    `json:"size"`
// 	Price        string    `json:"price"`
// 	ProductId    string    `json:"product_id"`
// 	Sequence     int64     `json:"sequence"`
// 	Time         time.Time `json:"time"`
// }

// func Connect() string {
// 	return url
// 	// conn, _, err := websocket.DefaultDialer.Dial(url, nil)
// 	// if err != nil {
// 	// 	fmt.Println("Error on websocket connection", err)
// 	// }

// 	// err = conn.WriteJSON(subscribeMessage{
// 	// 	Type:       "subscribe",
// 	// 	ProductIds: []string{"ETH-BTC"},
// 	// 	Channels:   []string{channel},
// 	// })

// 	// if err != nil {
// 	// 	fmt.Println("Error to subscribe on matches", err)
// 	// }

// 	// s := subscription{}
// 	// err = conn.ReadJSON(&s)

// 	// if err != nil {
// 	// 	fmt.Println("Error to subscribe on matches", err)
// 	// }
// 	// fmt.Println(s)
// }

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
		return errors.New("error to receive data from Coinbase WS")
	}

	if r.Type == "error" {
		err = errors.New(r.Message)
	}

	return err
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
