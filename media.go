package ember_backend_api_gateway

type MediaFile struct {
	Content  []byte
	Filename string
	MimeType string
}

type SendMediaRequest struct {
	Author   string
	Type     string
	Duration float64
	IsActive bool
}
