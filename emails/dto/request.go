package dto

type EmailRequest struct {
	Email   string
	Subject string
	Body    string
}

type EmailAgentAccount struct {
	Email    string
	Username string
	Password string
}
