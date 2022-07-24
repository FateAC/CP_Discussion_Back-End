package mail

import (
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpServer   = "smtp.gmail.com"
	smtpPort     = ":587"
	smtpUsername = "cpdiscussion.ta@gmail.com"
	smtpPassword = ""
)

func SendMail(to string, subject string, text string) error {
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	e := &email.Email{
		From:    "CPDiscussion <cpdiscussion.ta@gmail.com>",
		To:      []string{to},
		Subject: subject,
		Text:    []byte(text),
	}
	err := e.Send(smtpServer+smtpPort, auth)
	return err
}

func ResetPWDContent(to string, token string, url string) string {
	str := "哈囉, " +
		to +
		"\r\n您目前正在進行密碼重置操作，請點擊以下連結重置 CP_Discussion 的密碼。" +
		"\r\nLink: " +
		url +
		"resetpwd?token=" +
		token +
		"\r\n此連結有效時間 10 分鐘。\r\n若您並未嘗試重置密碼，請忽略此訊息。\r\n此郵件為系統自動發送，請勿回覆。\r\n\r\n" +
		"Hi, " +
		to +
		"\r\nYou are resetting your password of CP_Discussion.\r\n" +
		"Please click the following link to reset CP_Discussion password." +
		"\r\nLink: " +
		url +
		"resetpwd?token=" +
		token +
		"\r\nThis link will expire in 10 minutes.\r\nIf you didn't try to reset password, please ignore this message.\r\nThis is an automated reply from a system mailbox. Please do not reply to this email.\r\n\r\nBest Regards,\r\n\r\n" +
		"---\r\n" +
		"CP_Discussion\r\ncpdiscussion.ta@gmail.com\r\n"
	return str
}

func ResetPWDSuccess(to string) string {
	str := "哈囉, " +
		to +
		"\r\n您的密碼已被修改，如果您並未修改您的密碼，請盡速聯絡 CP_Discussion 管理員。\r\n" +
		"此郵件為系統自動發送，請勿回覆。\r\n\r\n" +
		"Hi, " +
		to +
		"\r\nYour password had been reset. If you didn't reset your password, please contact CP_Discussion administrator ASAP.\r\n" +
		"This is an automated reply from a system mailbox. Please do not reply to this email.\r\n\r\nBest Regards,\r\n\r\n" +
		"---\r\n" +
		"CP_Discussion\r\ncpdiscussion.ta@gmail.com\r\n"
	return str
}
