package chat

type Repository interface {
	Store(*Context) error
	Find(ChatID) (*Context, error)
}
