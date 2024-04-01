package request

func SuccessResponse(response interface{}) Response {
	return BuildResponse(true, "OK", response)
}
