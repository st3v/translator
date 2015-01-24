package translator

import "testing"

// Make sure nobody breaks the interface.
func TestTranslatorInterface(t *testing.T) {
	var translator Translator = &testTranslator{}
	translator.Translate("", "", "")
}

type testTranslator struct{}

func (t *testTranslator) Languages() ([]Language, error) {
	return nil, nil
}

func (t *testTranslator) Translate(text, from, to string) (string, error) {
	return "", nil
}

func (t *testTranslator) Detect(text string) (string, error) {
	return "", nil
}
