package message

type Status string

const (
	StatusOK             Status = "ok"
	StatusError          Status = "error"
	BookingOperationType        = "postgres"
)

type Message struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	OperationType string `json:"operation_type"`
	OperationInfo string `json:"operation_info"`
	Status        Status `json:"status"`
	Msg           string `json:"msg"`
}

func NewMessage(email, name, operationInfo, msg string, status Status) *Message {
	return &Message{Email: email, Name: name, OperationInfo: operationInfo, Status: status, Msg: msg, OperationType: BookingOperationType}
}
