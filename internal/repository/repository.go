package repository

type Repository interface {
}
type pg struct {
}

func New() Repository {
	return pg{}
}
