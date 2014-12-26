package translator

// The Language struct represents a given language by its
// name and code.
type Language struct {
	Code string
	Name string
}

// The Translator interface represents a translation service.
type Translator interface {
	// Languages returns a slice of language structs that are supported
	// by the given translator.
	Languages() ([]Language, error)

	// Transalate takes a string in a given language and returns its translation
	// to another language. Source and destination languages are specified by their
	// corresponding language codes.
	Translate(text, from, to string) (string, error)
}
