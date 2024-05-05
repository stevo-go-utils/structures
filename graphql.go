package structures

import "encoding/json"

type GraphqlPayload struct {
	OperationName string         `json:"operationName,omitempty"`
	Query         string         `json:"query,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`
}

func NewGraphqlPayload(operationName string, query string, variables map[string]any) *GraphqlPayload {
	return &GraphqlPayload{
		OperationName: operationName,
		Query:         query,
		Variables:     variables,
	}
}

func (g *GraphqlPayload) MustJson() string {
	jsonBytes, _ := json.Marshal(g)
	return string(jsonBytes)
}
