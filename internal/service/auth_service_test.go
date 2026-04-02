package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"

	appconfig "sandbox-game/internal/config"
	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

func TestAuthServiceLoginAndAuthenticateGroupAccount(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	groupID := createAdminIntegrationGroup(t, ctx, tx, "认证测试-玩家组", enum.BusinessStatusNormal, nil)
	username := fmt.Sprintf("auth_group_%d", time.Now().UnixNano())
	accountID := createAuthIntegrationAccount(t, ctx, tx, username, "123456", enum.RoleTypeGroup, &groupID)

	authService := NewAuthService(repository.NewAccountRepository(tx), repository.NewGroupRepository(tx), appconfig.AuthConfig{
		Mode:               "local",
		TokenSecret:        "integration-auth-secret",
		TokenExpireSeconds: 3600,
	})

	result, err := authService.Login(ctx, username, "123456")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if result.AccessToken == "" {
		t.Fatalf("expected access token to be returned")
	}
	if result.TokenType != authTokenTypeBearer {
		t.Fatalf("expected token type %q, got %q", authTokenTypeBearer, result.TokenType)
	}
	if result.User.UserID != accountID || result.User.RoleType != enum.RoleTypeGroup {
		t.Fatalf("unexpected login user info: %#v", result.User)
	}
	if result.DefaultRoute != "/sandbox-game/player/operating?yearNo=0" {
		t.Fatalf("unexpected default route: %s", result.DefaultRoute)
	}

	authenticated, err := authService.AuthenticateAccessToken(ctx, result.AccessToken)
	if err != nil {
		t.Fatalf("authenticate token failed: %v", err)
	}
	if authenticated.UserID != accountID || authenticated.GroupID == nil || *authenticated.GroupID != groupID {
		t.Fatalf("unexpected authenticated identity: %#v", authenticated)
	}

	account, err := repository.NewAccountRepository(tx).GetByID(ctx, accountID)
	if err != nil {
		t.Fatalf("reload account failed: %v", err)
	}
	if account.LastLoginTime == nil {
		t.Fatalf("expected last_login_time to be updated after login")
	}

	currentUser, err := authService.BuildCurrentUser(ctx, *authenticated)
	if err != nil {
		t.Fatalf("build current user failed: %v", err)
	}
	if currentUser.DefaultRoute != "/sandbox-game/player/operating?yearNo=0" {
		t.Fatalf("unexpected current user default route: %s", currentUser.DefaultRoute)
	}
}

func TestAuthServiceRejectsInvalidPassword(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	username := fmt.Sprintf("auth_admin_%d", time.Now().UnixNano())
	createAuthIntegrationAccount(t, ctx, tx, username, "123456", enum.RoleTypeAdmin, nil)

	authService := NewAuthService(repository.NewAccountRepository(tx), repository.NewGroupRepository(tx), appconfig.AuthConfig{
		Mode:               "local",
		TokenSecret:        "integration-auth-secret",
		TokenExpireSeconds: 3600,
	})

	_, err := authService.Login(ctx, username, "654321")
	if !errors.Is(err, ErrAuthInvalidCredentials) {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}

func TestAuthServiceAdminDefaultRouteDependsOnInitialization(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	if err := tx.WithContext(ctx).Where("1 = 1").Delete(&entity.Group{}).Error; err != nil {
		t.Fatalf("clear groups: %v", err)
	}

	username := fmt.Sprintf("auth_admin_setup_%d", time.Now().UnixNano())
	accountID := createAuthIntegrationAccount(t, ctx, tx, username, "123456", enum.RoleTypeAdmin, nil)

	authService := NewAuthService(repository.NewAccountRepository(tx), repository.NewGroupRepository(tx), appconfig.AuthConfig{
		Mode:               "local",
		TokenSecret:        "integration-auth-secret",
		TokenExpireSeconds: 3600,
	})

	setupResult, err := authService.Login(ctx, username, "123456")
	if err != nil {
		t.Fatalf("login before initialization failed: %v", err)
	}
	if setupResult.DefaultRoute != "/sandbox-game/admin/setup" {
		t.Fatalf("expected admin setup route before initialization, got %s", setupResult.DefaultRoute)
	}

	createAdminIntegrationGroup(t, ctx, tx, "认证测试-初始化后", enum.BusinessStatusNormal, nil)
	currentUser, err := authService.BuildCurrentUser(ctx, AuthenticatedUser{
		UserID:   accountID,
		Username: username,
		RoleType: enum.RoleTypeAdmin,
	})
	if err != nil {
		t.Fatalf("build current user after initialization failed: %v", err)
	}
	if currentUser.DefaultRoute != "/sandbox-game/admin/summary" {
		t.Fatalf("expected admin summary route after initialization, got %s", currentUser.DefaultRoute)
	}
}

func createAuthIntegrationAccount(
	t *testing.T,
	ctx context.Context,
	tx *gorm.DB,
	username string,
	password string,
	roleType string,
	groupID *int64,
) int64 {
	t.Helper()

	now := time.Now()
	account := entity.Account{
		Username:     username,
		PasswordHash: hashSHA256Password(password),
		RoleType:     roleType,
		GroupID:      groupID,
		Status:       enum.AccountStatusEnabled,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&account).Error; err != nil {
		t.Fatalf("create auth integration account: %v", err)
	}
	return account.ID
}
