package microsoft

const (
	authUrl          = "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13"
	serviceUrl       = "http://api.microsofttranslator.com/v2/Http.svc/"
	translationUrl   = serviceUrl + "Translate"
	languageNamesUrl = serviceUrl + "GetLanguageNames"
	languageCodesUrl = serviceUrl + "GetLanguagesForTranslate"
)

type Router interface {
	AuthUrl() string
	TranslationUrl() string
	LanguageNamesUrl() string
	LanguageCodesUrl() string
}

type router struct{}

func NewRouter() Router {
	return &router{}
}

func (r *router) AuthUrl() string {
	return authUrl
}

func (r *router) TranslationUrl() string {
	return translationUrl
}

func (r *router) LanguageNamesUrl() string {
	return languageNamesUrl
}

func (r *router) LanguageCodesUrl() string {
	return languageCodesUrl
}
