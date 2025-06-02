package auth

import (
	"cengkeHelperBackGo/internal/config"
	database "cengkeHelperBackGo/internal/db" // 如果 checkUser 没有返回完整的用户 DTO，可能需要这个
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"cengkeHelperBackGo/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt" // 如果 UserRegisterHandler 在同一个文件并且使用它，则保留
	"log"
	"net/http"
	"strings"
	"time"
)

func UserLoginHandler(c *gin.Context) {
	var req struct {
		// 你的 LoginPage.txt 前端似乎使用 'email' 作为登录标识
		Email    string `json:"email"` // 注意：这里从 Username 改为 Email 以匹配前端
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("请求参数错误: "+err.Error()))
		return
	}

	// services.checkUser 理想情况下应该在成功验证后返回完整的 dto.User 对象。
	// 如果它目前只返回布尔值或有限信息，你可能需要修改它，
	// 或者在密码检查成功后单独获取用户详细信息。
	user, ok := checkUser(req.Email, req.Password) // 将 req.Email 传递给 checkUser
	if !ok {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("邮箱或密码错误"))
		return
	}

	// 确保敏感数据（如密码哈希）不会包含在响应中。
	// dto.User 结构体最好是为安全的 API 响应而设计的，
	// 或者你可以创建一个特定的 dto.UserLoginResponse 结构体。
	// 此处我们假设 checkUser 返回的 user 对象已经是安全的。
	// 如果不是，你需要在这里填充一个新的结构体，例如：
	// safeUserResponse := dto.User{
	// 	Id:       user.Id,
	// 	Username: user.Username,
	// 	Email:    user.Email,
	// 	Avatar:   user.Avatar,
	// 	Bio:      user.Bio,
	// 	UserRole: user.UserRole,
	//  CreatedAt: user.CreatedAt, // 如果前端需要创建时间
	// }
	// 由于 checkUser 现在应该处理密码字段的剥离，可以直接使用 user

	expirationTime := time.Now().Add(5 * 24 * time.Hour) // 5 天有效期
	// GenerateUserToken 可能需要用户名或唯一标识符，以及角色
	// 确保 user.Username 是正确的，或者如果主要用 Email 标识，则传递 user.Email
	tokenIdentifierForJWT := user.Username               // 或者 user.Email，取决于你的JWT策略和用户模型
	if tokenIdentifierForJWT == "" && user.Email != "" { // 如果Username可能为空，但Email总是有值
		tokenIdentifierForJWT = user.Email
	}

	token, err := utils.GenerateUserToken(tokenIdentifierForJWT, user.UserRole, fmt.Sprintf("%d", user.Id))
	if err != nil {
		log.Printf("为用户 %s 生成令牌时出错: %v", tokenIdentifierForJWT, err)
		c.JSON(http.StatusInternalServerError, vo.NewBadResp("生成令牌失败，请联系管理员"))
		return
	}

	log.Printf("用户登录成功: %s (Email: %s), 角色: %d", user.Username, user.Email, user.UserRole)

	c.JSON(http.StatusOK, vo.RespData{
		Code: config.CodeSuccess,
		Data: gin.H{
			"token":        token,
			"expirationAt": expirationTime.Unix(),
			"userInfo":     user, // 在这里包含用户对象
			// "user":         safeUserResponse, // 如果你创建了一个独立的安全响应结构体
		},
		Msg: "登录鉴权成功！",
	})
}

// UserRegisterHandler 保持与你提供的 handler.go 文件中的一致，但我们也将使其返回用户信息和token
func UserRegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest // 确保 dto.RegisterRequest 包含前端发送的所有字段，例如 email, password, emailCode
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("请求参数错误: "+err.Error()))
		return
	}

	// 1. 检查邮箱唯一性
	var emailCount int64
	if err := database.Client.Model(&dto.User{}).
		Where("email = ?", req.Email).
		Count(&emailCount).Error; err != nil {
		c.JSON(http.StatusBadRequest, vo.NewBadResp(err.Error()))
		return
	}
	if emailCount > 0 {
		c.JSON(http.StatusBadRequest, vo.RespData{
			Code: config.CodeEmailExists, // 确保这个常量在 config 中定义
			Msg:  "该邮箱已被注册",
		})
		return
	}

	// 检查验证码
	if !checkEmailCode(req.Email, req.EmailCode) {
		c.JSON(http.StatusBadRequest, vo.NewBadResp("验证码错误"))
		return
	}

	// (可选) 检查用户名唯一性，如果你的系统中有独立的用户名概念并且在注册时收集
	// 你的 RegisterPage.txt 主要收集邮箱。你需要决定如何处理用户名：
	// - 是否与邮箱相同？
	// - 是否从邮箱派生？
	// - 是否在注册表单中单独收集？(如果是，dto.RegisterRequest 和前端都需要修改)
	// 假设 dto.RegisterRequest 中可能有一个 Username 字段，如果前端传了就用，否则可以考虑使用邮箱作为用户名。
	var usernameToCheck = req.Username // 假设 req.Username 来自 dto.RegisterRequest
	if usernameToCheck == "" {
		// 如果前端没有提供用户名，可以考虑将邮箱作为默认用户名（需要确保其唯一性，如果Username字段有唯一约束）
		// 或者生成一个唯一的用户名。为简单起见，如果dto.User的Username字段允许为空或不作唯一要求，可以不强制。
		// 但通常，用户名是需要的。
		usernameToCheck = strings.Split(req.Email, "@")[0] // 简单示例：邮箱前缀作为用户名，需进一步处理唯一性
	}

	var usernameCount int64
	if err := database.Client.Model(&dto.User{}).
		Where("username = ?", usernameToCheck).
		Count(&usernameCount).Error; err != nil {
		c.JSON(http.StatusBadRequest, vo.NewBadResp(err.Error()))
		return
	}
	if usernameCount > 0 {
		c.JSON(http.StatusBadRequest, vo.RespData{
			Code: config.CodeUsernameExists, // 确保这个常量在 config 中定义
			Msg:  "该用户名已存在",
		})
		return
	}
	// 如果你的系统强制要求用户名，但前端没传，这里应该返回错误，或赋予一个默认值（如邮箱）
	// 为了演示，如果dto.RegisterRequest中没有Username字段，
	// 我们可以决定在dto.User中将Username设置为Email。

	// 2. 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewBadResp(err.Error()))
		return
	}

	// 3. 设置默认值
	avatar := req.Avatar
	if avatar == "" {
		avatar = "default_avatar.png" // 你的默认头像路径
	}

	// 4. 创建用户记录
	newUser := dto.User{
		Username: usernameToCheck,        // 使用检查过或默认的用户名
		Password: string(hashedPassword), // 存储加密后的密码
		Email:    req.Email,
		UserRole: dto.UserRoleCommon, // 假设普通用户角色定义为 dto.UserRoleCommon
		Avatar:   avatar,
		Bio:      req.Bio, // 来自请求，可能为空
		// CreatedAt, UpdatedAt 会由 GORM 的 gorm.Model 自动处理（如果你的 dto.User 嵌入了它）
	}

	if result := database.Client.Create(&newUser); result.Error != nil {
		log.Printf("创建用户时出错: %v", result.Error)
		c.JSON(http.StatusBadRequest, vo.NewBadResp(err.Error()))
		return
	}

	// 用户创建成功，现在为自动登录生成令牌
	expirationTime := time.Now().Add(5 * 24 * time.Hour)

	token, err := utils.GenerateUserToken(newUser.Username, newUser.UserRole, fmt.Sprintf("%d", newUser.Id))
	if err != nil {
		log.Printf("为新注册用户 %+v 生成令牌时出错: %v", newUser, err)
		// 注册本身是成功的，但令牌生成失败。
		// 你仍然可以返回用户数据，但不带令牌，或者带一个关于令牌的错误消息。
		// 为简单起见，我们继续，但令牌可能为空。
	}

	// 准备响应的用户数据，dto里排除密码等敏感字段。
	// GORM Create 操作后，newUser 对象会包含数据库生成的 ID 和时间戳。
	log.Printf("用户注册成功: 用户名='%s', 邮箱='%s'", newUser.Username, newUser.Email)

	c.JSON(http.StatusCreated, vo.RespData{
		Code: config.CodeSuccess, // 使用统一的成功码
		Data: gin.H{
			"token":        token, // 发送令牌以便立即登录
			"expirationAt": expirationTime.Unix(),
			"userInfo":     newUser, // 返回新创建的用户信息
		},
		Msg: "注册成功，并已自动登录！",
	})
}
