package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	appconfig "sandbox-game/internal/config"
	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

const (
	authTokenTypeBearer   = "Bearer"
	sha256PasswordPrefix  = "{sha256}"
	defaultAuthSecret     = "sandbox-game-local-secret"
	defaultAuthExpireSecs = int64(28800)
)

var (
	ErrAuthInvalidCredentials = errors.New("auth: invalid credentials")
	ErrAuthAccountDisabled    = errors.New("auth: account disabled")
	ErrAuthInvalidToken       = errors.New("auth: invalid token")
	ErrAuthTokenExpired       = errors.New("auth: token expired")
)

type AuthService struct {
	accountRepo        *repository.AccountRepository
	groupRepo          *repository.GroupRepository
	tokenSecret        []byte
	tokenExpireSeconds int64
}

type AuthLoginUser struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	RoleType string `json:"roleType"`
	GroupID  *int64 `json:"groupId"`
}

type AuthLoginResult struct {
	AccessToken  string        `json:"accessToken"`
	TokenType    string        `json:"tokenType"`
	ExpiresIn    int64         `json:"expiresIn"`
	User         AuthLoginUser `json:"user"`
	DefaultRoute string        `json:"defaultRoute"`
}

type AuthCurrentUserResult struct {
	UserID       int64  `json:"userId"`
	Username     string `json:"username"`
	RoleType     string `json:"roleType"`
	GroupID      *int64 `json:"groupId"`
	DefaultRoute string `json:"defaultRoute"`
}

type AuthLogoutResult struct {
	LoggedOut bool `json:"loggedOut"`
}

type AuthenticatedUser struct {
	UserID   int64
	Username string
	RoleType string
	GroupID  *int64
}

type authTokenClaims struct {
	UserID    int64  `json:"userId"`
	Username  string `json:"username"`
	RoleType  string `json:"roleType"`
	GroupID   *int64 `json:"groupId"`
	ExpiresAt int64  `json:"expiresAt"`
}

func NewAuthService(accountRepo *repository.AccountRepository, groupRepo *repository.GroupRepository, cfg appconfig.AuthConfig) *AuthService {
	secret := strings.TrimSpace(cfg.TokenSecret)
	if secret == "" {
		secret = defaultAuthSecret
	}
	expireSeconds := cfg.TokenExpireSeconds
	if expireSeconds <= 0 {
		expireSeconds = defaultAuthExpireSecs
	}
	return &AuthService{
		accountRepo:        accountRepo,
		groupRepo:          groupRepo,
		tokenSecret:        []byte(secret),
		tokenExpireSeconds: expireSeconds,
	}
}

func (s *AuthService) Login(ctx context.Context, username string, password string) (*AuthLoginResult, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return nil, ErrAuthInvalidCredentials
	}

	account, err := s.accountRepo.GetByUsername(ctx, username)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, ErrAuthInvalidCredentials
	case err != nil:
		return nil, fmt.Errorf("load account by username: %w", err)
	}
	if !isAccountEnabled(account.Status) {
		return nil, ErrAuthAccountDisabled
	}
	if err := validateAccount(account); err != nil {
		return nil, fmt.Errorf("validate account: %w", err)
	}
	if !verifyPassword(password, account.PasswordHash) {
		return nil, ErrAuthInvalidCredentials
	}

	now := time.Now()
	claims := authTokenClaims{
		UserID:    account.ID,
		Username:  account.Username,
		RoleType:  account.RoleType,
		GroupID:   cloneInt64Pointer(account.GroupID),
		ExpiresAt: now.Add(time.Duration(s.tokenExpireSeconds) * time.Second).Unix(),
	}
	accessToken, err := s.buildAccessToken(claims)
	if err != nil {
		return nil, fmt.Errorf("build access token: %w", err)
	}
	if err := s.accountRepo.UpdateLastLoginTime(ctx, account.ID, now); err != nil {
		return nil, fmt.Errorf("update last login time: %w", err)
	}

	defaultRoute, err := s.resolveDefaultRoute(ctx, account.RoleType)
	if err != nil {
		return nil, err
	}

	return &AuthLoginResult{
		AccessToken: accessToken,
		TokenType:   authTokenTypeBearer,
		ExpiresIn:   s.tokenExpireSeconds,
		User: AuthLoginUser{
			UserID:   account.ID,
			Username: account.Username,
			RoleType: account.RoleType,
			GroupID:  cloneInt64Pointer(account.GroupID),
		},
		DefaultRoute: defaultRoute,
	}, nil
}

func (s *AuthService) AuthenticateAccessToken(ctx context.Context, accessToken string) (*AuthenticatedUser, error) {
	claims, err := s.parseAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	account, err := s.accountRepo.GetByID(ctx, claims.UserID)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, ErrAuthInvalidToken
	case err != nil:
		return nil, fmt.Errorf("load account by id: %w", err)
	}
	if !isAccountEnabled(account.Status) {
		return nil, ErrAuthInvalidToken
	}
	if err := validateAccount(account); err != nil {
		return nil, ErrAuthInvalidToken
	}
	if account.Username != claims.Username || account.RoleType != claims.RoleType || !sameNullableInt64(account.GroupID, claims.GroupID) {
		return nil, ErrAuthInvalidToken
	}

	return &AuthenticatedUser{
		UserID:   account.ID,
		Username: account.Username,
		RoleType: account.RoleType,
		GroupID:  cloneInt64Pointer(account.GroupID),
	}, nil
}

func (s *AuthService) BuildCurrentUser(ctx context.Context, identity AuthenticatedUser) (*AuthCurrentUserResult, error) {
	defaultRoute, err := s.resolveDefaultRoute(ctx, identity.RoleType)
	if err != nil {
		return nil, err
	}
	return &AuthCurrentUserResult{
		UserID:       identity.UserID,
		Username:     identity.Username,
		RoleType:     identity.RoleType,
		GroupID:      cloneInt64Pointer(identity.GroupID),
		DefaultRoute: defaultRoute,
	}, nil
}

func (s *AuthService) BuildLogoutResult() *AuthLogoutResult {
	return &AuthLogoutResult{LoggedOut: true}
}

func (s *AuthService) buildAccessToken(claims authTokenClaims) (string, error) {
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signature := s.signEncodedPayload(encodedPayload)
	return encodedPayload + "." + signature, nil
}

func (s *AuthService) parseAccessToken(accessToken string) (authTokenClaims, error) {
	trimmedToken := strings.TrimSpace(accessToken)
	if trimmedToken == "" {
		return authTokenClaims{}, ErrAuthInvalidToken
	}

	encodedPayload, signature, found := strings.Cut(trimmedToken, ".")
	if !found || encodedPayload == "" || signature == "" {
		return authTokenClaims{}, ErrAuthInvalidToken
	}
	if !compareEncodedSignature(signature, s.signEncodedPayload(encodedPayload)) {
		return authTokenClaims{}, ErrAuthInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return authTokenClaims{}, ErrAuthInvalidToken
	}

	var claims authTokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return authTokenClaims{}, ErrAuthInvalidToken
	}
	if claims.UserID <= 0 || strings.TrimSpace(claims.Username) == "" || !enum.IsValidRoleType(claims.RoleType) {
		return authTokenClaims{}, ErrAuthInvalidToken
	}
	if claims.RoleType == enum.RoleTypeGroup && claims.GroupID == nil {
		return authTokenClaims{}, ErrAuthInvalidToken
	}
	if claims.ExpiresAt <= time.Now().Unix() {
		return authTokenClaims{}, ErrAuthTokenExpired
	}
	return claims, nil
}

func (s *AuthService) signEncodedPayload(encodedPayload string) string {
	mac := hmac.New(sha256.New, s.tokenSecret)
	_, _ = mac.Write([]byte(encodedPayload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func validateAccount(account *entity.Account) error {
	if !enum.IsValidRoleType(account.RoleType) {
		return fmt.Errorf("unsupported role type %q", account.RoleType)
	}
	if account.RoleType == enum.RoleTypeGroup && account.GroupID == nil {
		return errors.New("group account missing group id")
	}
	return nil
}

func isAccountEnabled(status string) bool {
	return strings.EqualFold(strings.TrimSpace(status), enum.AccountStatusEnabled)
}

func verifyPassword(rawPassword string, storedHash string) bool {
	if !strings.HasPrefix(storedHash, sha256PasswordPrefix) {
		return false
	}
	sum := sha256.Sum256([]byte(rawPassword))
	calculatedHash := hex.EncodeToString(sum[:])
	expectedHash := strings.ToLower(strings.TrimPrefix(storedHash, sha256PasswordPrefix))
	return subtle.ConstantTimeCompare([]byte(calculatedHash), []byte(expectedHash)) == 1
}

func compareEncodedSignature(actual string, expected string) bool {
	actualBytes, err := base64.RawURLEncoding.DecodeString(actual)
	if err != nil {
		return false
	}
	expectedBytes, err := base64.RawURLEncoding.DecodeString(expected)
	if err != nil {
		return false
	}
	return hmac.Equal(actualBytes, expectedBytes)
}

func sameNullableInt64(left *int64, right *int64) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	return *left == *right
}

func cloneInt64Pointer(value *int64) *int64 {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}

func (s *AuthService) resolveDefaultRoute(ctx context.Context, roleType string) (string, error) {
	if roleType != enum.RoleTypeAdmin {
		return "/sandbox-game/player/operating?yearNo=0", nil
	}
	return s.resolveAdminDefaultRoute(ctx)
}

func (s *AuthService) resolveAdminDefaultRoute(ctx context.Context) (string, error) {
	if s.groupRepo == nil {
		return "/sandbox-game/admin/summary", nil
	}

	groupCount, err := s.groupRepo.CountAll(ctx)
	if err != nil {
		return "", fmt.Errorf("count groups for admin default route: %w", err)
	}
	if groupCount > 0 {
		return "/sandbox-game/admin/summary", nil
	}
	return "/sandbox-game/admin/setup", nil
}
