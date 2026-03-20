package pkg_error_response

type ErrorResponse struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}
