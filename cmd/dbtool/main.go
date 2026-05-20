package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	driver "github.com/go-sql-driver/mysql"

	appconfig "sandbox-game/internal/config"
)

const (
	actionInitSingle       = "init-single"
	actionResetSingle      = "reset-single"
	actionInitCompetition  = "init-competition"
	actionResetCompetition = "reset-competition"
)

func main() {
	var configPath string
	var action string

	flag.StringVar(&configPath, "config", "configs/local-single.yaml", "配置文件路径")
	flag.StringVar(&action, "action", actionInitSingle, "执行动作：init-single/reset-single/init-competition/reset-competition")
	flag.Parse()

	if err := run(configPath, action); err != nil {
		fmt.Fprintf(os.Stderr, "dbtool failed: %v\n", err)
		os.Exit(1)
	}
}

func run(configPath string, action string) error {
	cfg, err := appconfig.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	dsnConfig, err := driver.ParseDSN(cfg.MySQL.DSN)
	if err != nil {
		return fmt.Errorf("parse mysql dsn: %w", err)
	}
	if strings.TrimSpace(dsnConfig.DBName) == "" {
		return fmt.Errorf("mysql dsn missing database name")
	}

	targetDBName := dsnConfig.DBName
	adminDSN := cloneDSNWithoutDatabase(dsnConfig)

	adminDB, err := sql.Open("mysql", adminDSN)
	if err != nil {
		return fmt.Errorf("open admin mysql: %w", err)
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		return fmt.Errorf("ping admin mysql: %w", err)
	}

	switch action {
	case actionInitSingle:
		if err := ensureDatabase(adminDB, targetDBName); err != nil {
			return err
		}
	case actionInitCompetition:
		if err := ensureDatabase(adminDB, targetDBName); err != nil {
			return err
		}
	case actionResetSingle:
		if err := recreateDatabase(adminDB, targetDBName); err != nil {
			return err
		}
	case actionResetCompetition:
		if err := recreateDatabase(adminDB, targetDBName); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}

	targetDB, err := sql.Open("mysql", cfg.MySQL.DSN)
	if err != nil {
		return fmt.Errorf("open target mysql: %w", err)
	}
	defer targetDB.Close()

	if err := targetDB.Ping(); err != nil {
		return fmt.Errorf("ping target mysql: %w", err)
	}

	projectRoot, err := resolveProjectRoot(configPath)
	if err != nil {
		return err
	}

	migrationFiles, err := resolveMigrationFiles(projectRoot, action)
	if err != nil {
		return err
	}

	for _, file := range migrationFiles {
		if err := executeSQLFile(targetDB, file); err != nil {
			return err
		}
	}

	fmt.Printf("database ready: %s (%s)\n", targetDBName, action)
	return nil
}

func resolveMigrationFiles(projectRoot string, action string) ([]string, error) {
	commonFiles := []string{
		filepath.Join(projectRoot, "migrations", "mysql", "0001_init.sql"),
		filepath.Join(projectRoot, "migrations", "mysql", "0003_notice_adjustment.sql"),
		filepath.Join(projectRoot, "migrations", "mysql", "0005_order_admin.sql"),
		filepath.Join(projectRoot, "migrations", "mysql", "0006_order_generation_refactor.sql"),
	}

	switch action {
	case actionInitSingle, actionResetSingle:
		return append(commonFiles, filepath.Join(projectRoot, "migrations", "mysql", "0002_seed_single_group.sql")), nil
	case actionInitCompetition, actionResetCompetition:
		return append(commonFiles, filepath.Join(projectRoot, "migrations", "mysql", "0004_seed_competition_admin.sql")), nil
	default:
		return nil, fmt.Errorf("unsupported migration action: %s", action)
	}
}

func cloneDSNWithoutDatabase(cfg *driver.Config) string {
	cloned := *cfg
	cloned.DBName = ""
	return cloned.FormatDSN()
}

func ensureDatabase(db *sql.DB, dbName string) error {
	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS `" + dbName + "` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"); err != nil {
		return fmt.Errorf("create database %s: %w", dbName, err)
	}
	return nil
}

func recreateDatabase(db *sql.DB, dbName string) error {
	if _, err := db.Exec("DROP DATABASE IF EXISTS `" + dbName + "`"); err != nil {
		return fmt.Errorf("drop database %s: %w", dbName, err)
	}
	return ensureDatabase(db, dbName)
}

func resolveProjectRoot(configPath string) (string, error) {
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return "", fmt.Errorf("resolve config path: %w", err)
	}
	return filepath.Dir(filepath.Dir(absConfigPath)), nil
}

func executeSQLFile(db *sql.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read sql file %s: %w", path, err)
	}

	statements := splitSQLStatements(removeLineComments(content))
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("execute sql file %s: %w", path, err)
		}
	}

	return nil
}

func removeLineComments(content []byte) string {
	normalized := strings.ReplaceAll(string(bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))), "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		filtered = append(filtered, line)
	}
	return strings.Join(filtered, "\n")
}

func splitSQLStatements(content string) []string {
	var (
		statements []string
		builder    strings.Builder
		inSingle   bool
		inDouble   bool
		inBacktick bool
		escaped    bool
	)

	flush := func() {
		stmt := strings.TrimSpace(builder.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
		builder.Reset()
	}

	for _, r := range content {
		builder.WriteRune(r)

		if escaped {
			escaped = false
			continue
		}

		switch r {
		case '\\':
			if inSingle || inDouble {
				escaped = true
			}
		case '\'':
			if !inDouble && !inBacktick {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle && !inBacktick {
				inDouble = !inDouble
			}
		case '`':
			if !inSingle && !inDouble {
				inBacktick = !inBacktick
			}
		case ';':
			if !inSingle && !inDouble && !inBacktick {
				flush()
			}
		}
	}

	flush()
	return statements
}
