package datasource

type Routing struct {
	APIKey     string `dynamo:"api_key"`
	Path       string `dynamo:"path"`
	ForwardURL string `dynamo:"forward_url"`
	ContractID int    `dynamo:"contract_id"`
}
