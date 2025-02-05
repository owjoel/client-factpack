package model

type Client struct {
	ID          uint
	Name        string
	Age         uint
	Nationality string
	Status      string

	CreatedAt string
	UpdatedAt string
}

type StatusRes struct {
	Status string `json:"status"`
}

type GetClientRes struct {
	Name        string
	Age         uint
	Nationality string
}

type CreateClientReq struct {
	Name        string `json:"name"`
	Age         uint   `json:"age"`
	Nationality string `json:"nationality"`
}

type UpdateClientReq struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Age         uint   `json:"age"`
	Nationality string `json:"nationality"`
}

type DeleteClientReq struct {
	ID uint `json:"id"`
}
