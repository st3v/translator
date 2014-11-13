package translator

type Language struct {
	Code string
	Name string
}

type Translator interface {
	Languages() ([]Language, error)
	Translate(text, from, to string) (string, error)
}
