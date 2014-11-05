package query

import ()

// A query that merges results.
// Represented by a JSON array of generic queries, as the value of an $or key.
type query_or_t struct {
	query_sequence_t
}
