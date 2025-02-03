package services

import (
    "time"
    "gorm.io/gorm"
    "mobilerecharge/models"
)

type NotificationService struct {
    db *gorm.DB
    emailService *EmailService
}

func NewNotificationService(db *gorm.DB) *NotificationService {
    return &NotificationService{
        db: db,
        emailService: NewEmailService(),
    }
}

func (s *NotificationService) CheckAndSendNotifications() error {
    var sims []models.Sim
    if err := s.db.Find(&sims).Error; err != nil {
        return err
    }

    now := time.Now()
    for _, sim := range sims {
        daysUntilExpiry := int(sim.RechargeValidity.Sub(now).Hours() / 24)

        if daysUntilExpiry == 2 || daysUntilExpiry == 1 {
            err := s.emailService.SendExpiryNotification(
                sim.Name,
                sim.Number,
                sim.RechargeValidity,
                daysUntilExpiry,
            )
            if err != nil {
                return err
            }
        }
    }
    return nil
}