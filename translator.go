package translator

type Language struct {
	code string
	name string
}

type Translator interface {
	Languages() ([]Language, error)
	Translate(text, from, to string) (string, error)
}
