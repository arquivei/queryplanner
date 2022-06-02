package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arquivei/queryplanner"
	"github.com/arquivei/queryplanner/examples/providers"
)

/*
	- Utilizar _ em uma das dependências. (Criar provider que modifica formato do cpf)
*/

// Represents the mock of the covid database
var covidDatabase = map[string]bool{
	"44452427138": true,
	"85022625806": false,
	"20662340442": true,
}

// Represents the mock of the gov database
var govDatabase = map[string]providers.PersonalInfo{
	"44452427138": {
		Name: "João",
		Sex:  "Male",
	},
	"85022625806": {
		Name: "Maria",
		Sex:  "Female",
	},
	"20662340442": {
		Name: "José",
		Sex:  "Male",
	},
}

var indexerDatabase = []string{"85022625806", "20662340442", "44452427138"}

func main() {
	// Check the providers/cpfIndexer.go for explanation on IndexProvider
	cpfIndexer := providers.NewCPFIndexer(indexerDatabase)

	// Check providers/govDatabase.go for explanation on FieldProviders.
	govDatabase := providers.NewGovDatabaseProvider(govDatabase)
	covidDatabase := providers.NewCovidDatabaseProvider(covidDatabase)
	cpfFormatter := &providers.CPFFormatterProvider{}

	// QueryPlanner takes an IndexerProvider and multiple FieldProviders.
	// The indexer provider is the first to be executed and will set the response struct that will be modified by the other providers.
	// The FieldProviders can interact with databases, services or even just calculate attributes.
	// The query's are executed in order.
	planner, err := queryplanner.NewQueryPlanner(cpfIndexer, govDatabase, covidDatabase, cpfFormatter)
	if err != nil {
		panic(err)
	}

	fmt.Printf("------ Retrieving all fields (limit 2) ------\n\n")
	// Creates a plan of execution from the query planner and the request.
	plan := planner.NewPlan(&providers.Request{
		Fields: []string{
			"CPF",
			"Name",
			"HadCovid",
			"Sex",
		},
		Limit: 2,
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
	payload, err = planner.NewPlan(&providers.Request{
		Fields: []string{
			"Name",
		},
	}).Execute(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n")
	// Now only the name should be returned and only the gov database provider needed to be executed
	printAsJson(payload)
}

func printAsJson(entity interface{}) {
	j, _ := json.MarshalIndent(entity, "", " ")
	fmt.Println(string(j))
}
