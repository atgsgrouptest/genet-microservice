package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/atgsgrouptest/genet-microservice/AI-testing/models"
	//"github.com/go-playground/locales/root"
)

// dfsPathCollector recursively traverses the flowchart to collect all unique paths to leaf nodes.
func dfsPathCollector(node models.Node, currentPath []string, allPaths *[][]string) {
	// Build string with name, description, and flags
	nodeInfo := fmt.Sprintf(
		"%s [desc: %s, input: %t%s]",
		node.Name,
		node.Description,
		node.RequiresInput,
		func() string {
			//if node.ExternalLink != "" {
			//	return ", link: " + node.ExternalLink
			//}
			return ""
		}(),
	)

	newPath := append(currentPath, nodeInfo)

	// If it's a leaf node (no children), we've found a complete flow
	if len(node.Children) == 0 {
		pathCopy := make([]string, len(newPath))
		copy(pathCopy, newPath)
		*allPaths = append(*allPaths, pathCopy)
		return
	}

	// Recursively call for each child
	for _, child := range node.Children {
		dfsPathCollector(child, newPath, allPaths)
	}
}


// ProcessFlowchartJSON cleans an escaped JSON string, performs a DFS, prints individual flows,
// and returns them as a slice of string slices.
func ProcessFlowchartJSON(escapedJSON string) ([][]string, error) {
	// Clean the escaped JSON string
    
	root, err := ParseEscapedJSON[models.Node](escapedJSON)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	var allFlowPaths [][]string
	fmt.Println(root)
	// Start DFS to collect paths. Initial path is empty.
	dfsPathCollector(root, []string{}, &allFlowPaths)

	/*fmt.Println("--- Individual Flow Paths ---")
	for i, path := range allFlowPaths {
		fmt.Printf("Flow %d: %s\n", i+1, strings.Join(path, " -> "))
	}*/

	return allFlowPaths, nil
}

func CleanJSON(escaped string) ([]string,error) {
	

	flows, err := ProcessFlowchartJSON(escaped)
	
    if err != nil {
		return nil, err
	}
	var cleanedFlows []string
	for _, flow := range flows {
		cleanedFlows = append(cleanedFlows, strings.Join(flow, " -> "))
	}
	return cleanedFlows, nil
}

func ParseEscapedJSON[T any](input string) (T, error) {
	var result T

	// Step 1: Unescape the JSON string (e.g., turns \"key\":\"value\" into proper JSON)
	unquoted, err := strconv.Unquote(input)
	if err != nil {
		return result, fmt.Errorf("failed to unquote string: %w", err)
	}

	// Step 2: Parse JSON into result of type T
	if err := json.Unmarshal([]byte(unquoted), &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}