package dto

type HandleUserRequestBodyDto struct {
	UserRequestId uint64 `json:"userRequestId"`
	Status        int    `json:"status"`
	Message       string `json:"message"`
}
