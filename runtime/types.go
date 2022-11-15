package runtime

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorType    string `json:"errorType"`
}
