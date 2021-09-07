package gateway

type APIForwarding struct {
	APIKey     string `dynamo:"api_key"`
	Path       string `dynamo:"path"`
	ForwardURL string `dynamo:"forward_url"`
}

func (af APIForwarding) Field() Field {
	return Field{
		Template: *NewURITemplate(af.Path),
		Path:     *NewURITemplate(af.ForwardURL),
		Num:      5,
		Max:      10,
	}
}
