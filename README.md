# Temporal project
Project to explore Temporal API with Encore.dev framework.

RESTFUL service to trigger temporal workflows w.r.t to Bills. Ability to start workflow to increase invoice, and await confirmation signal.

# Installation
To run:
1. Install [Encore.dev](https://encore.dev/docs/install)
2. Install [Temporal locally](https://learn.temporal.io/getting_started/typescript/dev_environment/#set-up-a-local-temporal-development-cluster)
3. On a terminal, run `make temporal`.
4. On another, run `make dev` to run encore.

# Test
`make test`

# Overall flow 
## 1. Create Bill
```
curl --location --request POST 'http://127.0.0.1:4000/bill/456'
```

Response:
```
{
    "ID": "456",
    "status": "OPEN",
    "transactions": []
}
```


## 2. Increase Bill, pending confirmation
```
curl --location --request PUT 'http://127.0.0.1:4000/bill/456/' \
--header 'Content-Type: application/json' \
--data '{
    "timestamp": 10000,
    "itemName": "abc",
    "amount": {
        "number": 100,
        "currency": "USD"
    }
}'

```

Response
```
{
    "BillID": "456",
    "WorkflowID": "bill-456-XmrdfVQAtZXJ43BPtd45D"
}
```

The `BillID` and the `WorkflowID` will be used to trigger the confirmation endpoint.

## 3. Send Confirmation
```
curl --location 'http://127.0.0.1:4000/confirm/bill/456/bill-456-XmrdfVQAtZXJ43BPtd45D'
```

Response
```
{
    "message": "invoiced confirmed"
}

```

## 4. Get Bill
Specify currency via query param, e.g. `currency=GEL`. Otherwise, currency defaults to `USD`.
```
curl --location 'http://127.0.0.1:4000/bill/456?currency=GEL'
```

Response:
```
{
    "ID": "456",
    "status": "CLOSED",
    "transactions": [
        {
            "timestamp": 10000,
            "itemName": "abc",
            "amount": {
                "number": "271.00",
                "currency": "GEL"
            }
        }
    ]
}
```

## 5. Close Bill
```
curl --location --request PUT 'http://127.0.0.1:4000/close/bill/456/'
```

Response:
```
{
    "items": [
        "abc"
    ],
    "total": {
        "number": "100",
        "currency": "USD"
    }
}
```
