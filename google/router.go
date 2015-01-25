package google

// The Router provides necessary URLs to communicate with Google's Translate API
type Router interface {
	LanguagesURL() string
}

const baseURL = "https://www.googleapis.com/language/translate/v2/"

type router struct {
	baseURL      string
	languagesURL string
}

func newRouter() Router {
	return &router{
		languagesURL: baseURL + "languages",
	}
}

func (r *router) LanguagesURL() string {
	return r.languagesURL
}
