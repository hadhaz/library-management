package domain

type Author struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	BirthDate   string `json:"birth_date"`
	Nationality string `json:"nationality"`
}

type AuthorResponse struct {
	Authors []Author `json:"authors"`
}
