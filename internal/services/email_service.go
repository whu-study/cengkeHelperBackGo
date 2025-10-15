package services

import (
	"cengkeHelperBackGo/internal/config"
	database "cengkeHelperBackGo/internal/db"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"gopkg.in/gomail.v2"
)

// EmailService 邮箱服务结构体
type EmailService struct{}

// NewEmailService 创建邮箱服务实例
func NewEmailService() *EmailService {
	return &EmailService{}
}

// GenerateVerificationCode 生成6位数字验证码
func (e *EmailService) GenerateVerificationCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += n.String()
	}
	return code
}

// SendVerificationCode 发送验证码邮件
func (e *EmailService) SendVerificationCode(email string) (string, error) {
	// 生成验证码
	code := e.GenerateVerificationCode()

	// 创建邮件
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(config.Conf.Email.Username, config.Conf.Email.FromName))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "【蹭课小助手】邮箱验证码")

	// 邮件内容
	body := fmt.Sprintf(`
<html>
<body>
    <h2>蹭课小助手 - 邮箱验证</h2>
    <p>您好！</p>
    <p>您的验证码是：<strong style="color: #007cff; font-size: 24px;">%s</strong></p>
    <p>验证码有效期为5分钟，请及时使用。</p>
    <p>如果这不是您的操作，请忽略此邮件。</p>
    <br>
    <p>此邮件由系统自动发送，请勿回复。</p>
    <p>蹭课小助手团队</p>
</body>
</html>
	`, code)

	m.SetBody("text/html", body)

	// 配置SMTP
	d := gomail.NewDialer(
		config.Conf.Email.SmtpHost,
		config.Conf.Email.SmtpPort,
		config.Conf.Email.Username,
		config.Conf.Email.Password,
	)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Printf("发送邮件失败: %v", err)
		return "", fmt.Errorf("发送邮件失败: %v", err)
	}

	// 将验证码存储到Redis
	if err := database.SetEmailCode(email, code); err != nil {
		log.Printf("存储验证码到Redis失败: %v", err)
		// 即使Redis存储失败，也返回验证码，以便进行测试
		log.Printf("验证码已生成但未存储到Redis，测试用验证码: %s", code)
	}

	log.Printf("验证码已发送到邮箱: %s", email)
	return code, nil
}

// VerifyEmailCode 验证邮箱验证码
func (e *EmailService) VerifyEmailCode(email, inputCode string) bool {
	// 从Redis获取验证码
	storedCode, err := database.GetEmailCode(email)
	if err != nil {
		log.Printf("从Redis获取验证码失败: %v", err)
		// 如果Redis不可用，使用硬编码验证码作为后备方案
		if inputCode == "134567" {
			log.Printf("使用后备验证码验证成功: %s", email)
			return true
		}
		return false
	}

	// 验证码匹配
	if storedCode == inputCode {
		// 验证成功后删除验证码
		if err := database.DeleteEmailCode(email); err != nil {
			log.Printf("删除验证码失败: %v", err)
		}
		log.Printf("邮箱验证成功: %s", email)
		return true
	}

	log.Printf("邮箱验证失败，验证码不匹配: %s", email)
	return false
}
