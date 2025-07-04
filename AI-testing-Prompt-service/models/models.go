package models


type Node struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	RequiresInput bool   `json:"requires_input"`
	//ExternalLink  string `json:"external_link"`
	Children      []Node `json:"children"`
}

type  NegativeCaseResult struct {
			NegativeCases []string `json:"negative_cases"`
		}

type Outer struct {
	Response string `json:"response"`
  }