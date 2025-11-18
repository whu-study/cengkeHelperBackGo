package chat

import (
	"bytes"
	"cengkeHelperBackGo/internal/config"
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type ChatHandler struct {
	//
	//
	//
	httpClient *http.Client
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		httpClient: &http.Client{}, //
	}
}

// ChatStreamHandler
func (h *ChatHandler) ChatStreamHandler(c *gin.Context) {
	// 1.
	_, exists := c.Get("userId")
	if !exists {
		vo.RespondError(c, http.StatusUnauthorized, config.CodeUnauthorized, "用户未授权", nil)
		return
	}

	// 2.
	var reqDTO dto.ChatRequestDTO
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeBadRequest, "请求参数无效: "+err.Error(), nil)
		return
	}

	// 3.
	//
	agentReqBody, err := json.Marshal(reqDTO) //
	if err != nil {
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "创建 Agent 请求失败", err)
		return
	}

	// 4.
	//
	agentReq, err := http.NewRequest("POST", config.Conf.Agent.ServiceURL, bytes.NewBuffer(agentReqBody))
	if err != nil {
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "创建 Agent HTTP 请求失败", err)
		return
	}
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// (B)
		agentReq.Header.Set("Authorization", authHeader)
	}
	agentReq.Header.Set("Content-Type", "application/json")
	agentReq.Header.Set("Accept", "text/plain") //

	// 5.
	agentResp, err := h.httpClient.Do(agentReq)
	if err != nil {
		vo.RespondError(c, http.StatusServiceUnavailable, config.CodeServerError, "Agent 服务连接失败", err)
		return
	}
	defer agentResp.Body.Close()

	if agentResp.StatusCode != http.StatusOK {
		vo.RespondError(c, agentResp.StatusCode, config.CodeServerError, "Agent 服务返回错误", nil)
		return
	}

	// 6.
	//
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	//
	// Gin
	c.Stream(func(w io.Writer) bool {
		//
		_, err := io.Copy(w, agentResp.Body)
		if err != nil {
			//
			return false
		}
		//
		return false
	})
}
