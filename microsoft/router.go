package microsoft

const (
	authURL          = "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13"
	serviceURL       = "http://api.microsofttranslator.com/v2/Http.svc/"
	translationURL   = serviceURL + "Translate"
	languageNamesURL = serviceURL + "GetLanguageNames"
	languageCodesURL = serviceURL + "GetLanguagesForTranslate"
	detectURL        = serviceURL + "Detect"
)

// The Router provides the necessary URLs to communicate with
// Microsoft's API.
type Router interface {
	AuthURL() string
	TranslationURL() string
	LanguageNamesURL() string
	LanguageCodesURL() string
	DetectURL() string
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

func (r *router) LanguageNamesURL() string {
	return languageNamesURL
}

func (r *router) LanguageCodesURL() string {
	return languageCodesURL
}

func (r *router) DetectURL() string {
	return detectURL
}
