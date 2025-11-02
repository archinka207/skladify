package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		wantErr     bool
		errContains string
		checkConfig func(t *testing.T, cfg *Config)
	}{
		{
			name: "Успешная загрузка конфигурации",
			envVars: map[string]string{
				"DB_USER":     "test_user",
				"DB_PASSWORD": "test_pass",
				"DB_NAME":     "test_db",
			},
			wantErr: false,
			checkConfig: func(t *testing.T, cfg *Config) {
				assert.Contains(t, cfg.DatabaseURL, "host=localhost")
				assert.Contains(t, cfg.DatabaseURL, "user=test_user")
				assert.Contains(t, cfg.DatabaseURL, "password=test_pass")
				assert.Contains(t, cfg.DatabaseURL, "dbname=test_db")
				assert.Equal(t, "8080", cfg.Port)
			},
		},
		{
			name: "Отсутствует DB_USER",
			envVars: map[string]string{
				"DB_PASSWORD": "test_pass",
				"DB_NAME":     "test_db",
			},
			wantErr:     true,
			errContains: "DB_USER не установлена",
		},
		{
			name: "Отсутствует DB_PASSWORD",
			envVars: map[string]string{
				"DB_USER": "test_user",
				"DB_NAME": "test_db",
			},
			wantErr:     true,
			errContains: "DB_PASSWORD не установлена",
		},
		{
			name: "Отсутствует DB_NAME",
			envVars: map[string]string{
				"DB_USER":     "test_user",
				"DB_PASSWORD": "test_pass",
			},
			wantErr:     true,
			errContains: "DB_NAME не установлена",
		},
		{
			name: "Кастомные настройки",
			envVars: map[string]string{
				"DB_USER":     "custom_user",
				"DB_PASSWORD": "custom_pass",
				"DB_NAME":     "custom_db",
				"DB_HOST":     "custom_host",
				"DB_PORT":     "5433",
				"PORT":        "9090",
			},
			wantErr: false,
			checkConfig: func(t *testing.T, cfg *Config) {
				assert.Contains(t, cfg.DatabaseURL, "host=custom_host")
				assert.Contains(t, cfg.DatabaseURL, "port=5433")
				assert.Equal(t, "9090", cfg.Port)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем и очищаем env vars
			originalEnv := make(map[string]string)
			for key := range tt.envVars {
				originalEnv[key] = os.Getenv(key)
				os.Unsetenv(key)
			}
			defer func() {
				for key, value := range originalEnv {
					if value != "" {
						os.Setenv(key, value)
					} else {
						os.Unsetenv(key)
					}
				}
			}()

			// Устанавливаем тестовые env vars
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg, err := Load()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cfg)
				tt.checkConfig(t, cfg)
			}
		})
	}
}
