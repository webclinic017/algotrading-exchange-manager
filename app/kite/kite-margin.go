package kite

// curl https://api.kite.trade/margins/orders \
//     -H 'X-Kite-Version: 3' \
//     -H 'Authorization: token api_key:access_token' \
//     -H 'Content-Type: application/json' \
//     -d '[
//     {
//         "exchange": "NSE",
//         "tradingsymbol": "INFY",
//         "transaction_type": "BUY",
//         "variety": "regular",
//         "product": "CNC",
//         "order_type": "MARKET",
//         "quantity": 1,
//         "price": 0,
//         "trigger_price": 0
//     }
// ]'
