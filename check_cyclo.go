package queryplanner

import (
	"sort"

	"github.com/arquivei/foundationkit/errors"
)

const (
	cycleDetectionInStack = iota + 1
	cycleDetectionDone
)

// checkCycle uses a DFS algorithm to ensure that there is no cyclic
// dependency in the field provider graph
func checkCycle(fieldToProviderMap fieldProviderByName) error {
	const op = errors.Op("checkCycle")
	visitedNodes := map[FieldName]int{}

	// sorting keys to make the function deterministic
	keys := make([]string, 0, fieldToProviderMap.Length())
	keys = append(keys, fieldToProviderMap.GetFieldNames()...)
	sort.Strings(keys)

	for _, key := range keys {
		startingNode := FieldName(key)
		if visitedNodes[startingNode] == cycleDetectionDone {
			continue
		}
		cycle := getCycleFromNode(startingNode, visitedNodes, fieldToProviderMap)
		if cycle != nil {
			return errors.E(
				op,
				"cycle found in field dependency",
				errors.KV("cycle", getCycleString(cycle)),
			)
		}
	}
	return nil
}

func getCycleFromNode(
	node FieldName,
	visitedNodes map[FieldName]int,
	fieldToProviderMap fieldProviderByName,
) []FieldName {
	if visitedNodes[node] == cycleDetectionInStack {
		return []FieldName{node}
	}
	visitedNodes[node] = cycleDetectionInStack
	defer func() {
		visitedNodes[node] = cycleDetectionDone
	}()

	provider, providerExists := fieldToProviderMap.GetByName(node)
	if !providerExists {
		return nil
	}

	for _, child := range provider.DependsOn() {
		if visitedNodes[child] == cycleDetectionDone {
			continue
		}
		cycle := getCycleFromNode(child, visitedNodes, fieldToProviderMap)
		if cycle == nil {
			continue
		}
		return append(cycle, node)
	}

	return nil
}

func getCycleString(cycle []FieldName) string {
	cycleString := ""
	for i := len(cycle) - 1; i >= 0; i-- {
		cycleString += " -> " + string(cycle[i])
	}
	return cycleString
}

func checkMethodsFromFieldProvider(fieldProvider FieldProvider) error {
	const op = errors.Op("checkMethodsFromFieldProvider")
	for _, field := range fieldProvider.Provides() {
		if field.Fill == nil {
			return errors.E(op, "there is no `fill' method for field", errors.KV("fieldName", field.Name))
		}
		if field.Clear == nil {
			return errors.E(op, "there is no `clear` method for field", errors.KV("fieldName", field.Name))
		}
	}
	return nil
}
