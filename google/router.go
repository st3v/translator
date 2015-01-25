package google

// The Router provides necessary URLs to communicate with Google's Translate API
type Router interface {
	LanguagesURL() string
	TranslateURL() string
	DetectURL() string
}

const baseURL = "https://www.googleapis.com/language/translate/v2/"

type router struct {
	baseURL      string
	languagesURL string
	translateURL string
	detectURL    string
}

func newRouter() Router {
	return &router{
		languagesURL: baseURL + "languages",
		detectURL:    baseURL + "detect",
		translateURL: baseURL,
	}
}

func (r *router) LanguagesURL() string {
	return r.languagesURL
}

func (r *router) TranslateURL() string {
	return r.translateURL
}

func (r *router) DetectURL() string {
	return r.detectURL
}
