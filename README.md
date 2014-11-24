Work in progress. At this point this is me trying to become a Gopher.

> **Klask** (/klɑsk/, in Breton): to search for someone or something with care.

Klask aims to become a fast, generic search engine over HTTP, with a simplified API comparable to those found in ElasticSearch or MongoDB, with a combination of design aims:

- *In-memory first*: most of today's searchable datasets (e.g. transactional databases for e-commerce platforms) easily fit in a modern machine's memory. They rarely exceed a few 10s of gigabytes. This pratically means that the performance hit from using disk-first backends (InnoDB, PostgreSQL, MongoDB, etc) has become irrelevant.
- *Dataless index*: unless explicitely requested, Klask doesn't make the original indexed data accessible: you can `POST` a resource but you can't `GET` it back. It's a _search engine_, not a transactional datastore, so it focused on returning IDs for matching records in a specifid order. Foregoing storage of the original data can dramatically improve index sizes (e.g. for full-text indices).
- *Aggressive caching*: Any part of a query in Klask can be cached and reused for further queries.
- *Redis all the things*: Redis is a fast, generic in-memory storage engine with optional persistence, reliability, and scaling-out features. We don't want to reinvent all this, so we leverage its numerous primitives and extensibility to provide search functionality.
