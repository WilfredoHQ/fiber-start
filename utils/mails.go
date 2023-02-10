package utils

import (
	"context"
	"fmt"
	"time"

	sendinblue "github.com/sendinblue/APIv3-go-library/v2/lib"
	"github.com/wilfredohq/fiber-start/config"
)

func sendEmail(emailTo string, templateId int, params map[string]interface{}) {
	if !config.Config.EmailsEnabled {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := sendinblue.NewConfiguration()
	cfg.AddDefaultHeader("api-key", config.Config.EmailsApiKey)
	cfg.AddDefaultHeader("partner-key", config.Config.EmailsApiKey)

	sib := sendinblue.NewAPIClient(cfg)

	body := sendinblue.SendSmtpEmail{
		TemplateId: int64(templateId),
		To:         []sendinblue.SendSmtpEmailTo{{Email: emailTo}},
		Params:     params,
	}

	_, _, err := sib.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		return
	}
}

func SendWelcomeEmail(emailTo string, fullName string) {
	params := map[string]interface{}{
		"projectName": config.Config.ProjectName,
		"fullName":    fullName,
		"link":        config.Config.ClientUrl,
	}

	sendEmail(emailTo, 5, params)
}

func SendResetPasswordEmail(emailTo string) {
	tokenString, err := GetJwt(emailTo, config.Config.PasswordResetTokenExpirationMinutes)
	if err != nil {
		return
	}

	params := map[string]interface{}{
		"projectName":  config.Config.ProjectName,
		"validMinutes": config.Config.PasswordResetTokenExpirationMinutes,
		"link":         fmt.Sprintf("%s/restablecer?token=%s", config.Config.ClientUrl, tokenString),
	}

	sendEmail(emailTo, 6, params)
}
