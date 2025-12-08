package usecase

import (
	"notification_service/internal/application/dto"
	"notification_service/internal/application/interface/repository"
	"notification_service/internal/domain/entity"
)

func NewService(repo repository.NotificationRepository) *NotificationUseCase {
	serv := new(NotificationUseCase)
	serv.repo = repo
	return serv
}

type NotificationUseCase struct {
	repo repository.NotificationRepository
}

func (n *NotificationUseCase) SendNotification(input *dto.SendNotificationInputDto) dto.SendNotificationOutputDto {
	if (input.Email == "") {
		return dto.SendNotificationOutputDto{Status: "error", Msg : "Email is required to send a message"}
	}
	notification := new(entity.Notification)
	notification.Email = input.Email
	notification.Name = input.Name
	var notString string
	switch input.OperationType {
	case "booking":
		notString = input.Name
		switch input.Status {
		case "error":
			notString += (" Error booking the next room: " + input.OperationInfo)
		case "ok":
			notString += (" successful booking of the next room: " + input.OperationType)
		default:	
			return dto.SendNotificationOutputDto{Status: "error", Msg : "Incorrect operation status"}
		}
	case "change_hotel_info":
	// на данный момент пока еще не решили, будем ли реализовывать
	default:
		return dto.SendNotificationOutputDto{Status: "error", Msg : "Incorrect operation type"}	
	}
	notification.TextOfNotification = notString
	err := n.repo.SendNotification(notification)
	if (err != nil) {
		return dto.SendNotificationOutputDto{Status: "error", Msg : err.Error()}	
	}
	return dto.SendNotificationOutputDto{Status: "ok", Msg : ""}	
}