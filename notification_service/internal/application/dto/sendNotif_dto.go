package dto

type SendNotificationInputDto struct {
	Email string `json:"email"`
	Name string `json:"name"`
	OperationType string `json:"operation_type"`
	OperationInfo string `json:"operation_info"`
	Status string `json:"status"`
	Msg string `json:"msg"`
}

type SendNotificationOutputDto struct {
	Status string `json:"status"`
	Msg string `json:"msg"`
}