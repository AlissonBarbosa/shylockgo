package models

type ServerData struct {
  ID string `json:"id"`
  Name string `json:"name"`
  ProjectID string `json:"project_id"`
  HostID string `json:"host_id"`
  Domain string `json:"domain"`
  MemoryUsage int64 `json:"memory_usage"`
}

type DomainQueryResponse struct {
  Status string `json:"status"`
    Data   struct {
        ResultType string `json:"resultType"`
        Result     []struct {
            Metric struct {
                Domain string `json:"domain"`
            } `json:"metric"`
        } `json:"result"`
    } `json:"data"`
}
