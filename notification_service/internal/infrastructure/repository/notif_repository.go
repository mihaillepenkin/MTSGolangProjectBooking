package repository

import (
	"fmt"
	"net/smtp"
	"notification_service/internal/domain/entity"
)

func NewRep(fromEmail string, password string, smtpHost string, smtpPort string) *NotificationRepository {
		notRepo := new(NotificationRepository)
		notRepo.password = password
		notRepo.smtpHost = smtpHost
		notRepo.smtpPort = smtpPort
		return notRepo
	}

type NotificationRepository struct {
	fromEmail string
	password string
	smtpHost string
	smtpPort string
}

func (n *NotificationRepository) SendNotification(not *entity.Notification) error {
	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s",
		n.fromEmail,
		not.Email,
		"Notification From MtsBookingService",
		not.TextOfNotification)
		
	auth := smtp.PlainAuth("", n.fromEmail, n.password, n.smtpHost)
	err := smtp.SendMail(
		n.smtpHost+":"+n.smtpPort,
		auth,
		n.fromEmail,
		[]string{not.Email},
		[]byte(message),
	)
	return err
}