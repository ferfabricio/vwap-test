package coinbase

import (
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
		if err.Error() != "error to receive data from Coinbase WS" {
			t.Error(err)
		}
		defer s.Close()
	})
}

var upgrader = websocket.Upgrader{}

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
