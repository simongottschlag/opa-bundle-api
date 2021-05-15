# opa-bundle-api
Proof-of-concept for an API to produce OPA bundles


## cURL

### Download bundle

```shell
curl localhost:8080/bundle/bundle.tar.gz --output /tmp/bundle.tar.gz
```

### Download bundle (with If-None-Match)

If the header matches the current revision, a status code `304` should be returned.

```shell
curl --header "If-None-Match: 476d1f14d83110241366a81f82753523b850e150f55ed51bf5379f40cabc323d" localhost:8080/bundle/bundle.tar.gz --output /tmp/bundle.tar.gz
```

### Test bundle with OPA

```shell
opa eval --bundle /tmp/bundle.tar.gz --format pretty 'data.rules[i].id == 1; data.rules[i].role'
```

This should output the first role, something like:

```shell
+---+--------------------+
| i | data.rules[i].role |
+---+--------------------+
| 0 | "super_admin"      |
+---+--------------------+
```

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