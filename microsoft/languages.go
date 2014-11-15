package microsoft

import (
	"encoding/xml"
	"io/ioutil"
	"strings"

	"github.com/st3v/translator"
)

func (a *api) Languages() ([]translator.Language, error) {
	if a.languages == nil {
		codes, err := a.languageCodes()
		if err != nil {
			return nil, err
		}

		names, err := a.languageNames(codes)
		if err != nil {
			return nil, err
		}

		for i := range codes {
			a.languages = append(
				a.languages,
				translator.Language{
					Code: codes[i],
					Name: names[i],
				})
		}
	}
	return a.languages, nil
}

// Return a list of language names that correspond to a given list of language codes.
func (a *api) languageNames(codes []string) ([]string, error) {
	payload, _ := xml.Marshal(newArrayOfStrings(codes))
	uri := a.router.LanguageNamesUrl() + "?locale=en"

	response, err := a.sendRequest("POST", uri, strings.NewReader(string(payload)), "text/xml")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	result := &arrayOfStrings{}
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Strings, nil
}

// Return a list of language codes supported by the API.
func (a *api) languageCodes() ([]string, error) {
	response, err := a.sendRequest("GET", a.router.LanguageCodesUrl(), nil, "text/plain")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	result := &arrayOfStrings{}
	if err = xml.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Strings, nil
}
