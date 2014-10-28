
## Indexing

- typed fields (integer, float, date, time, word, text)
- each field has an index type (see "querying" below)
- each record stored in a hash, possibly except text fields

GIS support?
- not for MVP; lat/lng can still be queried

## Querying

- a list of filters

  - equality (hash index)
  - number inequality (zset index),
  - text inequality (using reversed zset, keys contain the ID)
  - text prefix (like inequality?)
  - fulltext (tf-idf index)
  - fuzzy fulltext (trigrams + tf-idf)

  clients should order their filters for performance
  (no query optimiser)

- a list of orders (first order most important; repeated sorting from last to
  first order)

- limit and offset
  (probably keep around query results with pagination and return cursor-like
  objects)

## TF-IDF

TF(term,document) = 0.5 +
  0.5 *
  Card({w = term, w in document) /
  max(Card(w, w in document))

IDF(term) = log Card(documents) - log Card(doc : term in doc)

### Index

term frequency

    tf:{term} -> zset( {doc-id}, {term-frequency} ... )

number of documents containing term

    doc:{term} -> {doc-count}

all documents

    docs -> {doc-id}

### Algorithm

each term in query has a weight (e.g. from fuzzy matching of terms)

    for each {term} in query
      idf(term) = SCARD(docs) - log SCARD(doc:{term})

result is incremented with the TF-IDF for each term at each step

    ZUNIONSTORE result length(query)
      tf:{term} ...
      WEIGHTS (idf*weight(term)) ...
    ZREVRANGEBYSCORE result +inf -inf WITHSCORES LIMIT 0 10

returns a list of {doc-id}, {score}


on import
double term frequency in e.g. headings
allow boosting?


# Why

https://news.ycombinator.com/item?id=7614577

> Nice to have but I'd stick with Solr/ES/Lucene for any serious string
searching/autocomplete work

