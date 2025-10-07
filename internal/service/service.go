package service

type Service struct {
	repo RepositoryI
}

func New(r RepositoryI) *Service {
	return &Service{repo: r}
}
