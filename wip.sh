echo flushdb | redis-cli -n 2

curl -XGET -D/dev/stderr http://localhost:3000/indices/test

validJson='{
  "id": "test",
  "fields": [{
    "name": "field1",
    "type": "integer/equality"
  }, {
    "name": "field2",
    "type": "integer/inequality"
  }]
}'

curl -XPOST -d"$validJson" -H 'Content-Type: application/json' -D/dev/stderr http://localhost:3000/indices
