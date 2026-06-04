package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appconfig "sandbox-game/internal/config"
)

func TestRegisterFrontendStaticRoutes(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	distDir := t.TempDir()
	assetsDir := filepath.Join(distDir, "assets")
	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		t.Fatalf("mkdir assets dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(distDir, "index.html"), []byte("<html>index</html>"), 0o644); err != nil {
		t.Fatalf("write index.html: %v", err)
	}
	if err := os.WriteFile(filepath.Join(assetsDir, "app.js"), []byte("console.log('ok');"), 0o644); err != nil {
		t.Fatalf("write app.js: %v", err)
	}

	engine := gin.New()
	engine.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	registerFrontendStaticRoutes(engine, testFrontendConfig(distDir), zap.NewNop())

	t.Run("serves index on root", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		engine.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", recorder.Code)
		}
		if body := recorder.Body.String(); body != "<html>index</html>" {
			t.Fatalf("unexpected body: %s", body)
		}
	})

	t.Run("serves index on spa route", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/sandbox-game/login", nil)

		engine.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", recorder.Code)
		}
		if body := recorder.Body.String(); body != "<html>index</html>" {
			t.Fatalf("unexpected body: %s", body)
		}
	})

	t.Run("serves built assets", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)

		engine.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", recorder.Code)
		}
		if body := recorder.Body.String(); body != "console.log('ok');" {
			t.Fatalf("unexpected body: %s", body)
		}
	})

	t.Run("keeps api 404 behavior", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/sandbox-game/missing", nil)

		engine.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", recorder.Code)
		}
	})
}

func testFrontendConfig(distDir string) *appconfig.Config {
	return &appconfig.Config{
		Frontend: appconfig.FrontendConfig{
			DistDir: distDir,
		},
	}
}
