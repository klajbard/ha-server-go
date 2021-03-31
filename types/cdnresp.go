package types

type ResponseBody struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type AWSResponse struct {
	AcceptRanges  string
	LastModified  string
	ContentLength int
	ETag          string
	ContentType   string
	Metadata      interface{}
	Body          ResponseBody
}
