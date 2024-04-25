package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"
	"os"

	"github.com/AlissonBarbosa/shylockgo/models"

)

func QueryGetPrometheus(query string) models.QueryResult {
  prometheus_url := fmt.Sprintf("%s:%s/api/v1/query?query=%s", os.Getenv("PROMETHEUS_URL"), os.Getenv("PROMETHEUS_PORT"), query)
  fmt.Println(prometheus_url)
  resp, err := http.Get(prometheus_url)
  if err != nil {
    return models.QueryResult{Query: query, Data: nil, Error: fmt.Errorf("Error: %v", err)}
  }
  defer resp.Body.Close()
  
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return models.QueryResult{Query: query, Data: nil, Error: fmt.Errorf("Error: %v", err)}
  }

  var response models.QueryResponse
  err = json.Unmarshal(body, &response)
  if err != nil {
    return models.QueryResult{Query: query, Data: nil, Error: fmt.Errorf("Error: %v", err)}
  }

  return models.QueryResult{Query: query, Data: response, Error: nil}
}
