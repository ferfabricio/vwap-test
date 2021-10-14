# VWAP calculation engine

## TODO

### 1. Retrieve data

Retrieve a data feed from the coinbase websocket and subscribe to the matches channel. Pull data for
the following three trading pairs:

- BTC-USD
- ETH-USD
- ETH-BTC

### 2. Calculate VWAP

Calculate the VWAP per trading pair using a sliding window of 200 data points. Meaning, when a new
data point arrives through the websocket feed the oldest data point will fall off and the new one will be
added such that no more than 200 data points are included in the calculation.

- The first 200 updates will have less than 200 data points included. That’s fine for this project.

### 3. Stream result VWAP

Stream the resulting VWAP values on each websocket update.

- Print to stdout or file is ok. Usually you would send them off through a message broker but a
simple print is perfect for this project.

## Knowledge base

### VWAP

The volume-weighted average price (VWAP) is a trading benchmark used by traders that gives the average price a security has traded at throughout the day, based on both volume and price. VWAP is important because it provides traders with insight into both the trend and value of a security.

## Integration details

### Coinbase

- Structure of subscribe message

```json
{
    "type": "subscribe",
    "product_ids": [
        "ETH-USD",
        "ETH-EUR"
    ],
    "channels": ["matches"]
}
```

- Return of subscription call

```json
{
    "type": "subscriptions",
    "channels": [
        {
            "name": "matches",
            "product_ids": [
                "ETH-USD",
                "ETH-EUR"
            ]
        }
    ]
}
```

- Error

```json
{
    "type": "error",
    "message": "Failed to subscribe",
    "reason": "Type has to be either subscribe or unsubscribe"
}
```

- Match type message

```json
{
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
}
```
