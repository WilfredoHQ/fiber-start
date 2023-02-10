package utils

import (
	"context"
	"fmt"
	"time"

	sendinblue "github.com/sendinblue/APIv3-go-library/v2/lib"
	"github.com/wilfredohq/fiber-start/configs"
)

func sendEmail(emailTo string, templateId int, params map[string]interface{}) {
	if !configs.Env.EmailsEnabled {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := sendinblue.NewConfiguration()
	cfg.AddDefaultHeader("api-key", configs.Env.EmailsApiKey)
	cfg.AddDefaultHeader("partner-key", configs.Env.EmailsApiKey)

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
		"projectName": configs.Env.ProjectName,
		"fullName":    fullName,
		"link":        configs.Env.ClientUrl,
	}

	sendEmail(emailTo, 5, params)
}

func SendResetPasswordEmail(emailTo string) {
	tokenString, err := GetJwt(emailTo, configs.Env.PasswordResetTokenExpirationMinutes)
	if err != nil {
		return
	}

	params := map[string]interface{}{
		"projectName":  configs.Env.ProjectName,
		"validMinutes": configs.Env.PasswordResetTokenExpirationMinutes,
		"link":         fmt.Sprintf("%s/restablecer?token=%s", configs.Env.ClientUrl, tokenString),
	}

	sendEmail(emailTo, 6, params)
}
