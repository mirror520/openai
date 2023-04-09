package chat

type Repository interface {
	Store(*Chat) error
	Find(ChatID) (*Chat, error)
	Close() error
}
