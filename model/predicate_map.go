package model

// PredicateMap is a map of boolean valued functions with single string inputs.
// These can be used to represent "lookups" of allowed values.
type PredicateMap map[string]func(string) bool

// Keys returns a new slice of the PredicateMap keys
func (a PredicateMap) Keys() []string {
	keys := make([]string, 0)
	for s := range a {
		keys = append(keys, s)
	}
	return keys
}
