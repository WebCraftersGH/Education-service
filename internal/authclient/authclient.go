package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/WebCraftersGH/Education-service/pkg/logging"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("unauthorized")

type AuthResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     logging.Logger
}

func New(baseURL string, logger logging.Logger) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		logger: logger,
	}
}

func (c *Client) Check(ctx context.Context, token string) (uuid.UUID, error) {
	token = strings.TrimSpace(token)

	// Логируем начало с частичным токеном (безопаснее)
	tokenPreview := token
	if len(token) > 20 {
		tokenPreview = token[:20] + "..."
	}
	c.logger.WithField("token_preview", tokenPreview).Info("check auth start")

	if token == "" {
		return uuid.Nil, ErrUnauthorized
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/auth/check",
		nil,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("authclient: build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	// Детальное логгирование запроса
	c.logger.WithFields(map[string]interface{}{
		"method":         req.Method,
		"url":            req.URL.String(),
		"headers":        req.Header,
		"content_length": req.ContentLength,
		"host":           req.Host,
	}).Info("📤 sending request to auth service")

	// Логируем curl-команду для тестирования
	curlCmd := fmt.Sprintf("curl -X %s '%s' -H 'Authorization: Bearer %s'",
		req.Method, req.URL.String(), tokenPreview)
	c.logger.WithField("curl", curlCmd).Debug("equivalent curl command")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.WithError(err).Error("❌ http request failed")
		return uuid.Nil, fmt.Errorf("authclient: do request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.WithError(err).Error("read response body error")
		return uuid.Nil, fmt.Errorf("authclient: read response: %w", err)
	}

	// Логируем ответ
	c.logger.WithFields(map[string]interface{}{
		"status_code":      resp.StatusCode,
		"status_text":      http.StatusText(resp.StatusCode),
		"response_body":    string(bodyBytes),
		"response_headers": resp.Header,
		"content_length":   len(bodyBytes),
	}).Info("📥 received response from auth service")

	// Восстанавливаем тело для декодирования
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Если ответ пустой
	if len(bodyBytes) == 0 {
		c.logger.Warn("empty response body received")
		if resp.StatusCode == http.StatusOK {
			userID, err := c.extractUserIDFromToken(token)
			if err != nil {
				return uuid.Nil, fmt.Errorf("authclient: extract user_id: %w", err)
			}
			return userID, nil
		}
		return uuid.Nil, fmt.Errorf("authclient: empty response with status %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		c.logger.WithFields(map[string]interface{}{
			"error":    err.Error(),
			"raw_body": string(bodyBytes),
		}).Error("❌ failed to decode JSON response")
		return uuid.Nil, fmt.Errorf("authclient: decode response: %w", err)
	}

	c.logger.WithField("check_response", authResp).Info("✅ decoded auth response")

	switch resp.StatusCode {
	case http.StatusOK:
		if authResp.Error != "" {
			c.logger.WithField("error_msg", authResp.Error).Warn("auth error in successful response")
			return uuid.Nil, ErrUnauthorized
		}

		userID, err := c.extractUserIDFromToken(token)
		if err != nil {
			c.logger.WithError(err).Error("failed to extract user ID from token")
			return uuid.Nil, fmt.Errorf("authclient: extract user_id from token: %w", err)
		}

		c.logger.WithField("user_id", userID).Info("✅ authentication successful")
		return userID, nil

	case http.StatusUnauthorized, http.StatusForbidden:
		c.logger.Warn("authentication failed - unauthorized")
		return uuid.Nil, ErrUnauthorized

	default:
		c.logger.WithField("status", resp.StatusCode).Error("unexpected status code")
		return uuid.Nil, fmt.Errorf(
			"authclient: unexpected status=%d body=%s",
			resp.StatusCode,
			strings.TrimSpace(string(bodyBytes)),
		)
	}
}

func (c *Client) extractUserIDFromToken(tokenString string) (uuid.UUID, error) {
	// Парсим JWT без верификации (токен уже проверен auth-service)
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims type")
	}

	// Извлекаем uid из claims (как показано в логе)
	var uidStr string

	// Пробуем разные варианты названий полей
	if uid, exists := claims["uid"]; exists {
		uidStr = fmt.Sprintf("%v", uid)
	} else if userID, exists := claims["user_id"]; exists {
		uidStr = fmt.Sprintf("%v", userID)
	} else if sub, exists := claims["sub"]; exists {
		uidStr = fmt.Sprintf("%v", sub)
	} else {
		return uuid.Nil, fmt.Errorf("no uid/user_id/sub found in token claims")
	}

	if uidStr == "" {
		return uuid.Nil, fmt.Errorf("uid is empty")
	}

	// Парсим в UUID
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %s, error: %w", uidStr, err)
	}

	return userID, nil
}
