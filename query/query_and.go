package query

import ()

// A query that intersects results, progressively transforming the set of IDs.
// Represented by a JSON array of generic queries, as the value of an $and key.
type query_and_t struct {
	query_sequence_t
}
