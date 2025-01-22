package server

type AddRequest struct {
	ImageName string `json:"image_name"`
}

type RemoveRequest struct {
	ImageName string `json:"image_name"`
}
