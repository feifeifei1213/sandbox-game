package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appconfig "sandbox-game/internal/config"
)

type frontendStaticAssets struct {
	distDir   string
	indexPath string
}

func registerFrontendStaticRoutes(engine *gin.Engine, cfg *appconfig.Config, logger *zap.Logger) {
	assets, err := newFrontendStaticAssets(cfg.Frontend.DistDir)
	if err != nil {
		logger.Warn("frontend static serving disabled",
			zap.String("distDir", strings.TrimSpace(cfg.Frontend.DistDir)),
			zap.Error(err),
		)
		return
	}

	logger.Info("frontend static serving enabled",
		zap.String("distDir", assets.distDir),
	)

	engine.NoRoute(func(c *gin.Context) {
		assets.handleNoRoute(c)
	})
}

func newFrontendStaticAssets(distDir string) (*frontendStaticAssets, error) {
	resolvedDistDir, err := resolveFrontendDistDir(distDir)
	if err != nil {
		return nil, err
	}

	indexPath := filepath.Join(resolvedDistDir, "index.html")
	info, err := os.Stat(indexPath)
	if err != nil {
		return nil, fmt.Errorf("stat frontend index file: %w", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("frontend index path is a directory: %s", indexPath)
	}

	return &frontendStaticAssets{
		distDir:   resolvedDistDir,
		indexPath: indexPath,
	}, nil
}

func resolveFrontendDistDir(distDir string) (string, error) {
	trimmed := strings.TrimSpace(distDir)
	if trimmed == "" {
		trimmed = "frontend/dist"
	}

	if !filepath.IsAbs(trimmed) {
		workingDir, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get working directory: %w", err)
		}
		trimmed = filepath.Join(workingDir, trimmed)
	}

	resolvedDistDir, err := filepath.Abs(trimmed)
	if err != nil {
		return "", fmt.Errorf("resolve frontend dist dir: %w", err)
	}

	info, err := os.Stat(resolvedDistDir)
	if err != nil {
		return "", fmt.Errorf("stat frontend dist dir: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("frontend dist path is not a directory: %s", resolvedDistDir)
	}

	return resolvedDistDir, nil
}

func (assets *frontendStaticAssets) handleNoRoute(c *gin.Context) {
	if isReservedBackendRoute(c.Request.URL.Path) || !isStaticFileMethod(c.Request.Method) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if filePath, ok := assets.resolveRequestFile(c.Request.URL.Path); ok {
		c.File(filePath)
		return
	}

	c.File(assets.indexPath)
}

func isReservedBackendRoute(path string) bool {
	return path == "/healthz" || strings.HasPrefix(path, "/api/")
}

func isStaticFileMethod(method string) bool {
	return method == http.MethodGet || method == http.MethodHead
}

func (assets *frontendStaticAssets) resolveRequestFile(requestPath string) (string, bool) {
	trimmedPath := strings.TrimPrefix(requestPath, "/")
	if trimmedPath == "" {
		return "", false
	}

	cleanPath := filepath.Clean(filepath.FromSlash(trimmedPath))
	candidatePath := filepath.Join(assets.distDir, cleanPath)

	relativePath, err := filepath.Rel(assets.distDir, candidatePath)
	if err != nil {
		return "", false
	}
	if relativePath == "." || strings.HasPrefix(relativePath, "..") {
		return "", false
	}

	info, err := os.Stat(candidatePath)
	if err != nil || info.IsDir() {
		return "", false
	}

	return candidatePath, true
}
