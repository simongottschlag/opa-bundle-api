# opa-bundle-api
Proof-of-concept for an API to produce OPA bundles


## cURL

### Read All Rules

```shell
curl localhost:8080/rules
```

### Read Rule

```shell
curl localhost:8080/rules/1
```

### Create Rule

```shell
DATA='{"country": "Iceland", "city": "Reykjavik", "building": "Branch", "role": "user", "device_type": "Printer", "action": "allow"}'
curl -X POST --header "Content-Type: application/json" --data $DATA localhost:8080/rules
```

### Update Rule

```shell
DATA='{"country": "Iceland", "city": "Reykjavik", "building": "Branch", "role": "user", "device_type": "Printer", "action": "allow"}'
curl -X PUT --header "Content-Type: application/json" --data $DATA localhost:8080/rules/1
```

### Delete Rule

```shell
curl -X DELETE localhost:8080/rules/1
```