package ember_backend_api_gateway

type Media struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ContentType string `json:"content"`
	OwnerId     string `json:"owner_id"`
}

type CreateMediaRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ContentType string `json:"content"`
	OwnerId     string `json:"owner_id"`
}

type GetMediaRequest struct {
	ID string `json:"id"`
}

type UpdateMediaRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type DeleteMediaRequest struct {
	ID string `json:"id"`
}

type ListMediaRequest struct {
	OwnerId string `json:"owner_id"`
	Limit   int32  `json:"limit"`
}
