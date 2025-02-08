package emails

import (
	"fmt"
	"net/http"
	"nganterin-cs/emails/dto"
	"nganterin-cs/pkg/exceptions"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(data dto.EmailRequest) *exceptions.Exception {
	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	server := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")

	i, err := strconv.Atoi(smtpPort)
	if err != nil {
		return exceptions.NewException(http.StatusInternalServerError, exceptions.ErrInternalServer)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", data.Email)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", data.Body)

	d := gomail.NewDialer(server, i, email, password)

	if err := d.DialAndSend(m); err != nil {
		return exceptions.NewException(http.StatusBadGateway, exceptions.ErrEmailSendFailed)
	}

	return nil
}

func SendAgentAccountEmail(data dto.EmailAgentAccount) *exceptions.Exception {
	body := fmt.Sprintf(`
		<html lang="en">
		<head>
		    <meta charset="UTF-8" />
		    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
		    <title>Customer Service Account Created</title>

		    <style>
		    @import url("https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap");

		    body {
		        font-family: "Poppins", sans-serif;
		        display: flex;
		        flex-direction: column;
		        justify-content: center;
		        align-items: center;
		        min-height: 100vh;
		        margin: 0;
		        padding: 0;
		        background-color: #f6f6f6;
		    }
		    .container {
		        max-width: 600px;
		        margin: 0 auto;
		        padding: 30px;
		        border-radius: 20px;
		        background-color: #ffffff;
		        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
		    }
		    .header {
		        text-align: center;
		        margin-bottom: 20px;
		    }
		    .title {
		        font-size: 28px;
		        font-weight: 600;
		        color: #171717;
		        margin-bottom: 10px;
		    }
		    .message {
		        font-size: 15px;
		        line-height: 1.6;
		        color: #171717;
		        margin-bottom: 20px;
		        text-align: center;
		    }
		    .credentials-box {
		        background-color: #f8f9fa;
		        border-radius: 10px;
		        padding: 20px;
		        margin: 20px 0;
		        text-align: left;
		    }
		    .credentials-item {
		        margin: 10px 0;
		        font-size: 15px;
		        color: #171717;
		    }
		    .warning-message {
		        font-size: 14px;
		        color: rgb(255, 85, 85);
		        margin-top: 20px;
		        text-align: center;
		    }
		    </style>
		</head>
		<body>
		    <div class="container">
		        <div class="header">
		            <p class="title">Customer Service Account Created</p>
		        </div>
		        <p class="message" style="font-weight: 600">Dear Customer Service Agent,</p>
		        <p class="message">
		            Your customer service account has been created successfully. Below are your account credentials:
		        </p>
		        <div class="credentials-box">
		            <p class="credentials-item"><strong>Email:</strong> %s</p>
		            <p class="credentials-item"><strong>Username:</strong> %s</p>
		            <p class="credentials-item"><strong>Password:</strong> %s</p>
		        </div>
		        <p class="message">
		            Please change your password upon your first login for security purposes.
		        </p>
		        <p class="warning-message">
		            If you did not expect to receive these credentials, please contact the administrator immediately.
		        </p>
		    </div>
		    <p class="footer" style="font-size: 14px">
		        © 2024 Nganterin. All Rights Reserved.,<br />Best regards, Nganterin
		    </p>
		</body>
		</html>
	`, data.Email, data.Username, data.Password)

	emailData := dto.EmailRequest{
		Email:   data.Email,
		Subject: "Nganterin - Customer Service Account Credentials",
		Body:    body,
	}

	err := SendEmail(emailData)
	if err != nil {
		return err
	}

	return nil
}
