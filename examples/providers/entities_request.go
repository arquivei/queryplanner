package providers

// Request represents the request struct that will be passed to the queryPlanner.
// It should contain a method that returns the fields that the user wants filled in their request.
// The request might still contain information needed for your query execution like pagination, filters, etc.
type Request struct {
	Limit  int // custom field that represents the max quantity of documents that will be returned in the request.
	Fields []string
}

// GetRequestedFields returns the fields that the request wants filled in the response.
func (r *Request) GetRequestedFields() []string {
	return r.Fields
}
