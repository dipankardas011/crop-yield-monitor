package main

import "fmt"

type ImageUpload struct {
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

type AuthResponse struct {
	Stdout  string `json:"stdout"`
	Error   string `json:"error"`
	Account any
}

func (r AuthResponse) String() string {
	return fmt.Sprintf("{ Err: %s, Stdout: %s, Account: %v }", r.Error, r.Stdout, r.Account)
}
