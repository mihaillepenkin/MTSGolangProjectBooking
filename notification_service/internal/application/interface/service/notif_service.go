package service

import "notification_service/internal/application/dto"


type NotificationService interface {
	SendNotification(*dto.SendNotificationInputDto) dto.SendNotificationOutputDto
}