package microsoft

const (
	authURL          = "https://api.cognitive.microsoft.com/sts/v1.0/issueToken"
	serviceURL       = "https://api.cognitive.microsofttranslator.com/"
	translationURL   = serviceURL + "Translate"
	detectURL        = serviceURL + "Detect"
	languageNamesURL = serviceURL + "Languages"
	languageCodesURL = serviceURL + "Languages"
	apiVersion       = "3.0"
)

// The Router provides necessary URLs to communicate with
// Microsoft's API.
type Router interface {
	AuthURL() string
	TranslationURL() string
	DetectURL() string
	LanguageNamesURL() string
	LanguageCodesURL() string
	ApiVersion() string
}

type router struct{}

func newRouter() Router {
	return &router{}
}

func (r *router) AuthURL() string {
	return authURL
}

func (r *router) TranslationURL() string {
	return translationURL
}

func (r *router) DetectURL() string {
	return detectURL
}

func (r *router) LanguageNamesURL() string {
	return languageNamesURL
}

func (r *router) LanguageCodesURL() string {
	return languageCodesURL
}

func (r *router) ApiVersion() string {
	return apiVersion
}
