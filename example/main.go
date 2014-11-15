package main

import (
	"fmt"

	"github.com/st3v/translator/microsoft"
)

func main() {
	translator := microsoft.NewTranslator("globe", "gQwodqYqfffKHRCh/3iudM7k/7I0JoqcvSc8fH4Dpf0=")
	fmt.Println(translator.Translate("Hello World!", "en", "de"))
}
