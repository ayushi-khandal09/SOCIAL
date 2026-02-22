package mailer

import "embed"

const (
	FromName            = "GopherSocial"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitations.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}
