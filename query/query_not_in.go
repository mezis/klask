package query

import ()

// An exclusion filter (returns only values not in the list)
type query_filter_not_in_t struct {
	query_filter_membership_t
}
