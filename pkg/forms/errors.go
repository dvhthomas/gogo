package forms

// Define a new errors type. The name of the form field will hold
// validate errors for that field.
// Note that this is *not* a struct, but we will still add methods
type errors map[string][]string

// Implement an Add() method to add errors for a give field to the map
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get the first error message for a given field
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
