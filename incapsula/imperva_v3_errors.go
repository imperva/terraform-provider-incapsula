package incapsula

type APIErrors struct {
	Status int               `json:"status"`
	Id     string            `json:"id"`
	Code   int               `json:"code"`
	Source map[string]string `json:"source"`
	Title  string            `json:"title"`
	Detail string            `json:"detail"`
}
