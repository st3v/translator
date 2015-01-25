package google

const baseURL = "https://www.googleapis.com/language/translate/v2/"

type router struct {
	baseURL           string
	languagesEndpoint string
	translateEndpoint string
	detectEndpoint    string
}

func newRouter() *router {
	return &router{
		languagesEndpoint: baseURL + "languages",
		detectEndpoint:    baseURL + "detect",
		translateEndpoint: baseURL,
	}
}

func (r *router) languagesURL() string {
	return r.languagesEndpoint
}

func (r *router) translateURL() string {
	return r.translateEndpoint
}

func (r *router) detectURL() string {
	return r.detectEndpoint
}
