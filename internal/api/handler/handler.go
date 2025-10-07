package handler

type Handler struct {
	service ServiceI
}

func New(s ServiceI) *Handler {
	return &Handler{service: s}
}

type CreateRequest struct {
	URL          string `json:"url"       validate:"required"`
	UserShortURL string `json:"user_short_url" validate:"-"`
}
