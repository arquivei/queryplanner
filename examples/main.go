package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arquivei/queryplanner"
	"github.com/arquivei/queryplanner/examples/providers"
)

// Represents the mock of the covid database
var covidDatabase = map[string]bool{
	"44452427138": true,
	"85022625806": false,
}

// Represents the mock of the gov database
var govDatabase = map[string]providers.PersonalInfo{
	"44452427138": {
		Name: "João",
	},
	"85022625806": {
		Name: "Maria",
	},
}

// Request represents the request struct that will be passed to the queryPlanner.
// It should contain a method that returns the fields that the user wants filled in their request.
// The request might still contain information needed for your query execution like pagination, filters, etc.
type Request struct {
	fields []string
}

// GetRequestedFields returns the fields that the request wants filled in the response.
func (r *Request) GetRequestedFields() []string {
	return r.fields
}

func main() {
	// Check the providers/cpfIndexer.go for explanation on IndexProvider
	cpfIndexer := providers.NewCPFIndexer()

	// Check providers/govDatabase.go for explanation on FieldProviders.
	govDatabase := providers.NewGovDatabaseProvider(govDatabase)
	covidDatabase := providers.NewCovidDatabaseProvider(covidDatabase)

	// QueryPlanner takes an IndexerProvider and multiple FieldProviders.
	// The indexer provider is the first to be executed and will set the response struct that will be modified by the other providers.
	// The FieldProviders can interact with databases, services or even just calculate attributes.
	planner, err := queryplanner.NewQueryPlanner(cpfIndexer, govDatabase, covidDatabase)
	if err != nil {
		panic(err)
	}

	fmt.Printf("------ Retrieving all fields ------\n\n")
	// Creates a plan of execution from the query planner and the request.
	plan := planner.NewPlan(&Request{
		fields: []string{
			"CPF",
			"Name",
			"HadCovid",
		},
	})

	// Execution of the query.
	payload, err := plan.Execute(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n")
	// Now we can retrieve the response struct from the payload
	printAsJson(payload)

	fmt.Printf("\n\n------ Retrieving only the name field ------\n\n")

	// If we want for instance to retrieve only the Name of the people:
	payload, err = planner.NewPlan(&Request{
		fields: []string{
			"Name",
		},
	}).Execute(context.Background())

	fmt.Printf("\n\n")
	// Now only the name should be returned and only the gov database provider needed to be executed
	printAsJson(payload)
}

func retrievePersonSlice(docs []queryplanner.Document) []providers.Person {
	people := make([]providers.Person, len(docs))
	for i, d := range docs {
		people[i] = d.(providers.Person)
	}
	return people
}

func printAsJson(entity interface{}) {
	j, _ := json.MarshalIndent(entity, "", " ")
	fmt.Println(string(j))
}
