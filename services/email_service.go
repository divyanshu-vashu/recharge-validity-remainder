package services

import (
    "fmt"
    "net/smtp"
    "time"
    "os"
    "github.com/joho/godotenv"
)

const (
    smtpHost = "smtp.gmail.com"
    smtpPort = "587"
    fromEmail = "vdkalife@gmail.com"
)

var toEmails = []string{
    "vashusingh2005.jan@gmail.com",
    "divyanshusingh@appointy.com",
    "divyanshu.singh2021c@vitstudent.ac.in",
}

type EmailService struct{}

func NewEmailService() *EmailService {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        fmt.Printf("Error loading .env file: %v\n", err)
    }
    return &EmailService{}
}

func (s *EmailService) SendExpiryNotification(simName string, number string, expiryDate time.Time, daysLeft int) error {
    password := os.Getenv("EMAIL_PASSWORD")
    auth := smtp.PlainAuth("", fromEmail, password, smtpHost)
    
    subject := "SIM Recharge Reminder"
    body := fmt.Sprintf(
        "Your SIM card %s (%s) will expire in %d days on %s. Please recharge soon to avoid service interruption.",
        simName,
        number,
        daysLeft,
        expiryDate.Format("2006-01-02"),
    )
    
    msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
    
    return smtp.SendMail(
        smtpHost+":"+smtpPort,
        auth,
        fromEmail,
        toEmails,
        msg,
    )
}