package common

// Convenience method to create an error response with the given
// message only.
//
// msg: message to add to the error response
func MsgErrorRes(msg string) ErrorResponse {
	return ErrorResponse{
		Msg: msg,
	}
}

// Convenience method to create an error response with the given
// message and adding the error message as well.
//
// msg: message to add to the error response
// err: error occurred during process of request.
func ErrorRes(msg string, err error) ErrorResponse {
	return ErrorResponse{
		msg,
		err.Error(),
	}
}

// Supported responses for urchin shared between app and admin-app.
type ErrorResponse struct {
	Msg string `json:"msg"`
	Err string `json:"error,omitempty"`
}
