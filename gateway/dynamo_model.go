package gateway

type APIForwarding struct {
	APIKey     string `dynamo:"api_key"`
	Path       string `dynamo:"path"`
	ForwardURL string `dynamo:"forward_url"`
}
