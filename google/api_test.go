package google

import "github.com/st3v/translator"

type mockLanguageProvider struct {
	languagesFunc func() ([]translator.Language, error)
	detectFunc    func(text, version string) (string, error)
}

func (m *mockLanguageProvider) languages() ([]translator.Language, error) {
	return m.languagesFunc()
}

func (m *mockLanguageProvider) detect(text, version string) (string, error) {
	return m.detectFunc(text, version)
}

type mockTranslationProvider struct {
	translateFunc func(text, from, to, version string) (string, error)
}

func (m *mockTranslationProvider) translate(text, from, to, version string) (string, error) {
	return m.translateFunc(text, from, to, version)
}
