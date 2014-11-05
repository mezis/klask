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

searchJson1='{
  "filters": [
    ["field1", "neq", 4321],
    ["field1", "neq", 6789]
  ]
}'

searchJson2='{
  "$or": [{
    "field1": { "$neq": [4321, 6789] },
    "field2": { "$gt": 1111 }
  },{
    "$and": [{
      "field1": { "$neq": 8766 }
    }, {
      "field2": { "$lte": 4567 }
    }]
  }],
  "field1": { "$neq": 3771 },
  "$by": "+field2"
}'

searchJson3='{
  "$sort

}'

# operators
# $gt $lt $eq $in $ni $neq

host='http://localhost:3000'

curl -XGET -D/dev/stderr ${host}/indices/test

flush

curl -XPOST -d"$indexJson" -H 'Content-Type: application/json' -D/dev/stderr ${host}/indices

curl -XGET -D/dev/stderr ${host}/indices

curl -XPOST -d"$recordJson" -H 'Content-Type: application/json' -D/dev/stderr ${host}/indices/test/records

curl -XDELETE -D/dev/stderr ${host}/indices/test/records/1337

curl -XGET -d"$searchJson1" -H "Content-Type: application/json" -D/dev/stderr ${host}/indices/test/records
