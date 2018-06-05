package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/ninjadotorg/handshake-dispatcher/services"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"fmt"
	"github.com/ninjadotorg/handshake-dispatcher/config"
)

type VerifierController struct{}

func (s VerifierController) SendPhoneVerification(c *gin.Context) {
	phone := c.DefaultQuery("phone", "")
	countryCode := c.DefaultQuery("country", "")
	locale := c.DefaultQuery("locale", "en")

	twilioClient := services.TwilioService{}
	success, err := twilioClient.SendVerification(countryCode, phone, locale)
	if err != nil {
		resp := JsonResponse{0, "Send verification failed", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{0, "Send verification failed", nil}
	if success {
		resp = JsonResponse{1, "", nil}
	}

	c.JSON(http.StatusOK, resp)
}

func (s VerifierController) CheckPhoneVerification(c *gin.Context) {
	phone := c.DefaultQuery("phone", "")
	countryCode := c.DefaultQuery("country", "")
	code := c.DefaultQuery("code", "")

	twilioClient := services.TwilioService{}
	success, err := twilioClient.CheckVerification(countryCode, phone, code)
	if err != nil {
		resp := JsonResponse{0, "Check verification failed", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{0, "Phone verified failed", nil}
	if success {
		resp = JsonResponse{1, "", nil}
	}

	c.JSON(http.StatusOK, resp)
}

func (s VerifierController) SendEmailVerification(c *gin.Context) {
	email := c.DefaultQuery("email", "")
	// locale := c.DefaultQuery("locale", "en")

	bytes, err := bcrypt.GenerateFromPassword([]byte(email), 14)
	token := string(bytes)
	token = url.QueryEscape(token)

	cfg := config.GetConfig()
	verificationUrl := cfg.GetString("email_verification_url")
	verificationUrl = fmt.Sprintf("%s?token=%s", verificationUrl, token)

	mailClient := services.MailService{}

	subject := "Email verification"
	content := fmt.Sprintf(EMAIL_VERIFICATION_TEMPLATE, verificationUrl)
	success, err := mailClient.Send(" shake@shake.ninja ", email, subject, content)

	if err != nil || !success {
		resp := JsonResponse{0, "Send verification failed", nil}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := JsonResponse{1, "", nil}
	c.JSON(http.StatusOK, resp)
}

func (s VerifierController) CheckEmailVerification(c *gin.Context) {
	email := c.DefaultQuery("email", "")
	token := c.DefaultQuery("token", "")

	err := bcrypt.CompareHashAndPassword([]byte(token), []byte(email))

	resp := JsonResponse{0, "Email verified failed", nil}
	if err == nil {
		resp = JsonResponse{1, "", nil}
	}
	c.JSON(http.StatusOK, resp)
}

const EMAIL_VERIFICATION_TEMPLATE = `<html>
<body>
<p>
    Hi,
</p>
<p>
    All you need to do is verify your email address. Click the following <a href="%s">link</a> to get started.
</p>
<p>
    If you have any questions just send an email to <a href="mailto:shake@shake.ninja">shake@shake.ninja</a>
</p>
<p>
    Thanks,<br/>
    Aliesha
</p>
<p>
    Shake | Ninja
</p>
<p>
    Join the conversation at <a href="https://t.me/ninjadotorg">https://t.me/ninjadotorg</a>
</p>
</body>
</html>
`