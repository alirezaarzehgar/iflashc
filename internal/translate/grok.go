package translate

type grok struct {
	To string
}

func (g grok) Translate(text string) (string, error) {
	panic("not implemented")
	return "", nil
}
