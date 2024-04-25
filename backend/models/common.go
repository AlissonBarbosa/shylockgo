package models

type QueryResponse struct {
  Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name   string `json:"__name__"`
				Domain string `json:"domain,omitempty"`
			} `json:"metric,omitempty"`
			Value []interface{} `json:"value,omitempty"`
		} `json:"result,omitempty"`
  } `json:"data"`
}

type QueryResult struct {
  Query string
  Data interface{}
  Error error
}
