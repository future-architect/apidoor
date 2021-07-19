package managementapi

type Api struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}
