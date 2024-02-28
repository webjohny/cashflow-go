package request

func SuccessResponse() Response {
	return BuildResponse(true, "OK", nil)
}
