package main

import (
	"fmt"
	"sync"

	"github.com/st3v/translator"
	"github.com/st3v/translator/microsoft"
)

// Instantiates a translator that is backed by the Microsoft Translation API and passes it to helloWorld.
// Get your own clientId and clientSecret by registering an app with Microsoft's Azure DataMarket (see http://goo.gl/DJtCGy)
func main() {
	translator := microsoft.NewTranslator("your-client-id", "your-client-secret")
	helloWorld(translator)
}

// Fetches all supported languages and triggers concurrent translations of 'Hello World' for each of them.
func helloWorld(t translator.Translator) {
	languages, err := t.Languages()
	if err != nil {
		fmt.Printf("Error retrieving languages: %s\n", err)
		return
	}

	fmt.Printf("%d Supported Languages:\n", len(languages))
	fmt.Println("-----------------------")

	translations := make([]<-chan string, len(languages))
	for i, language := range languages {
		translations[i] = translate(t, "Hello World!", "en", language)
	}

	for n := range mergeChannels(translations) {
		fmt.Println(n)
	}
}

// Starts a go routine to translate text for a particular language. Returns a channel that will be
// used to send either the translation or an error string if something went wrong.
func translate(t translator.Translator, text, from string, to translator.Language) <-chan string {
	out := make(chan string)
	go func() {
		translation, err := t.Translate(text, from, to.Code)
		if err != nil {
			out <- fmt.Sprintf("Error during translation for %s: %s", to.Name, err)
		} else {
			out <- fmt.Sprintf("%s [%s]: %s", to.Name, to.Code, translation)
		}
		close(out)
	}()
	return out
}

// Merges a slice of incoming channels of strings into a single incoming channel of strings.
func mergeChannels(cs []<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
