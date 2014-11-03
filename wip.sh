flush() {
  echo flushdb | redis-cli -n 2
}

indexJson='{
  "id": "test",
  "fields": [{
    "name": "field1",
    "type": "integer/equality"
  }, {
    "name": "field2",
    "type": "integer/inequality"
  }]
}'

recordJson='{
  "id": 1337,
  "field1": 1234,
  "field2": 4567
}'

searchJson='{
  "filter": [
    ["field1", "neq", 4321],
    ["field1", "neq", 6789],
    ["field2", "gte", 1111],
    ["field2", "lte", 9999]
  ],
  "order": [
    ["field2", "asc"]
  ]
}'

host='http://localhost:3000'

curl -XGET -D/dev/stderr ${host}/indices/test

flush

curl -XPOST -d"$indexJson" -H 'Content-Type: application/json' -D/dev/stderr ${host}/indices

curl -XGET -D/dev/stderr ${host}/indices

curl -XPOST -d"$recordJson" -H 'Content-Type: application/json' -D/dev/stderr ${host}/indices/test/records

curl -XDELETE -D/dev/stderr ${host}/indices/test/records/1337
