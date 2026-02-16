package httpfetch

type Result struct {
	StatusCode  int
	ContentType string
	Body        []byte
	FinalURL    string
}
