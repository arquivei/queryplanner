package providers

// Person describes the document that we want retrieve with our query.
type Person struct {
	CPF      *string // CPF is an identifier that we will retrieve from the index provider
	Name     *string // Name is retrieved from the government database
	HadCovid *bool   // HasCovid is retrieved from the covid database
}
