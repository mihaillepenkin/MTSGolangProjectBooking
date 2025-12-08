package repository

import "notification_service/internal/domain/entity"

type NotificationRepository interface {
	SendNotification(*entity.Notification) error
}