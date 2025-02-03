package main

import (
    "fmt"
    "net/smtp"
)

func main() {
    // Email configuration
    from := "vdkalife@gmail.com"
    password := "" //gmail app password
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    // Recipients
    to := []string{
        "vashusingh2005.jan@gmail.com",
        "divyanshusingh@appointy.com",
        "divyanshu.singh2021c@vitstudent.ac.in",
    }

    // Email content
    subject := "Test Email - Mobile Recharge Reminder"
    body := "This is a test email from the Mobile Recharge Reminder application.\n\n" +
            "If you receive this email, the email service is working correctly."

    message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

    // Authentication
    auth := smtp.PlainAuth("", from, password, smtpHost)

    // Send email
    err := smtp.SendMail(
        smtpHost+":"+smtpPort,
        auth,
        from,
        to,
        message,
    )

    if err != nil {
        fmt.Printf("Error sending email: %v\n", err)
        return
    }

    fmt.Println("Test email sent successfully!")
}
