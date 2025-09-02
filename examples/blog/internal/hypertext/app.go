package hypertext

type App struct{}

type Article struct {
	Title   string
	Content string
	Error   error
}

func (*App) Article(id int) Article {
	return Article{} // noop, the fake will be used to set values in test cases
}
