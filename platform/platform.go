package platform

type Info struct {
	Title       string
	Platform    string
	Description string
}

type Platform interface {
	Match(url string) bool
	GetName() string
}
