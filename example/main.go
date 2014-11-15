package main

import (
	"fmt"

	"github.com/st3v/translator/microsoft"
)

func main() {
	translator := microsoft.NewTranslator("client-id", "client-secret")
	fmt.Println(translator.Translate("Hello World!", "en", "de"))
}
