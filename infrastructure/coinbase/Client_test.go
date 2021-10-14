package coinbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

// TODO: this test will connect in the coinbase, can be an anti-pattern
func TestCreateNewClient(t *testing.T) {
	_, err := NewClient()
	if err != nil {
		t.Error(err)
	}
}

// TODO: Improve tests to use test cases
func TestConfiguration(t *testing.T) {
	t.Run("With success", func(t *testing.T) {
		m := SubscriptionResult{
			Type: "subscriptions",
			Channels: []ChannelSubscription{
				{
					Name: "matches",
					ProductIds: []string{
						"ETH-USD",
						"ETH-EUR",
					},
				},
			},
		}
		s := createTestWs(t, wsMockIgnoreReceived(m))
		c := Client{
			conn: s,
		}
		err := c.Configure([]string{"ETH-USD", "ETH-EUR"})
		if err != nil {
			t.Error(err)
		}
		defer s.Close()
	})

	t.Run("With error", func(t *testing.T) {
		m := GenericResult{
			Type:    "error",
			Message: "Failed to subscribe",
			Reason:  "Type has to be either subscribe or unsubscribe",
		}
		s := createTestWs(t, wsMockIgnoreReceived(m))
		c := Client{
			conn: s,
		}
		err := c.Configure([]string{"ETH-USD", "ETH-EUR"})
		if err.Error() != m.Message {
			t.Fail()
		}
		defer s.Close()
	})

	t.Run("With connection error", func(t *testing.T) {
		s := createTestWs(t, func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Log(err)
			}
			defer c.Close()
		})
		c := Client{
			conn: s,
		}
		err := c.Configure([]string{"ETH-USD", "ETH-EUR"})
		if err.Error() != "error to configure Coinbase WS subscription" {
			t.Error(err)
		}
		defer s.Close()
	})
}

func TestGetData(t *testing.T) {
	expSuccess := `{
		"type": "match",
		"trade_id": 165850138,
		"maker_order_id": "41e0e6f8-7aa7-4f99-9f1a-fbaf99f9a66b",
		"taker_order_id": "0d611470-1c11-457c-9c92-02ad4c026268",
		"side": "sell",
		"size": "0.05",
		"price": "3438.23",
		"product_id": "ETH-USD",
		"sequence": 21661080566,
		"time": "2021-10-13T14:23:38.231856Z"
	}`

	expError := `{
		"type": "error",
		"message": "some error message"
	}`

	sc := []struct {
		exp string
		tp  string
	}{
		{expSuccess, "match"},
		{expError, "error"},
	}

	for _, sce := range sc {
		m := MatchMessage{}
		err := json.Unmarshal([]byte(sce.exp), &m)
		if err != nil {
			t.Fail()
		}
		s := createTestWs(t, wsMockSendMessages(m))
		defer s.Close()
		c := Client{
			conn: s,
		}
		ch := make(chan MatchMessage)
		go c.GetData(ch)
		r := <-ch
		if r.Type != sce.tp {
			t.Fail()
		}
	}
}

var upgrader = websocket.Upgrader{}

// mock websocket for tests source: https://stackoverflow.com/questions/47637308/create-unit-test-for-ws-in-golang/47637670
func createTestWs(t *testing.T, r http.HandlerFunc) *websocket.Conn {
	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(r))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	return ws
}

func wsMockIgnoreReceived(re interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}

			err = c.WriteJSON(re)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}
}

func wsMockSendMessages(re interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			err = c.WriteJSON(re)
			if err != nil {
				break
			}
		}
	}
}
