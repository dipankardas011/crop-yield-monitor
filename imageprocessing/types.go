package main

type ImageGet struct {
	Uuid        string `json:"uuid"`
	AccessToken string `json:"token"`
}

type ImageUpload struct {
	Uuid     string `json:"uuid"`
	RawImage []byte `json:"raw_image_bytes"`
	Format   string `json:"image_format"`
}

type ImageGetResp struct {
	RawImage []byte `json:"raw_image_bytes"`
}

type Response struct {
	Error  string `json:"errors"`
	Stdout string `json:"stdout"`
	Image  any
}
