package admin_app

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

// Supported responses for urchin.

type ErrorResponse struct {
	Msg string `json:"msg"`
	Err string `json:"error"`
}

type PostIdResponse struct {
	Id int `json:"id"`
}

type GetPostResponse struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content"`
}

type ImageIdResponse struct {
	Id string `json:"id"`
}
