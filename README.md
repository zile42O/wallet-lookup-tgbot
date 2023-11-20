# wallet-lookup-tgbot
Telegram bot for lookup crypto wallet by address with custom API

* This bot was created to look up any crypto wallet by address.
* The bot will automatically detect the coin type by providing a wallet address to check.
* Bot uses custom API for data
  `https://api.zile42o.dev/wallet_lookup/api.php`
  * To buy this API please contact me via telegram [@Zile42O](https://t.me/zile42O)

* Try bot now  [@WalletLookupBot](https://t.me/WalletLookupBot)

### API Example
GET `https://api.zile42o.dev/wallet_lookup/api.php?key=walletlookup:1234567890&address=1ZU8E66ylUsPr8Eau8qeXNTQHJSv8yCnRj`
#### Response
```json
{
    "data": {
        "coin": "bitcoin",
        "market_price_usd": 37338,
        "address": "1ZU8E66ylUsPr8Eau8qeXNTQHJSv8yCnRj",
        "total_transactions": 0,
        "transactions_received": 0,
        "transactions_sent": 0,
        "receieved_value": 0,
        "receieved_value_usd": 0,
        "sent_value": 0,
        "sent_value_usd": 0,
        "balance": 0,
        "balance_usd": 0,
        "first_transaction": "1970-01-01",
        "last_transaction": "1970-01-01",
        "transaction_list": {
            "in_transactions": [],
            "out_transactions": []
        }
    },
    "cached_ttl": 1700439295,
    "author": "zile42O"
}
```
