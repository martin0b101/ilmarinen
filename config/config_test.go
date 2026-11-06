package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateConfiguration(t *testing.T) {

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "happy case",
			config: Config{
				Address: "localhost",
				Port:    3000,
				Database: DatabaseConfig{
					Host:         "database-host",
					DatabaseName: "database-name",
					Port:         5678,
					Username:     "testUser",
					Password:     "test-pwd",
				},
				Worker: WorkerConfig{
					Counter: 20,
				},
				Job: JobConfig{
					Timeout:       10 * time.Second,
					MaxAttempts:   10,
					RetryDelay:    10 * time.Second,
					LeaseDuration: 10 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "empty database host",
			config: Config{
				Address: "localhost",
				Port:    3000,
				Database: DatabaseConfig{
					Host:         "",
					DatabaseName: "database-name",
					Port:         5678,
					Username:     "testUser",
					Password:     "test-pwd",
				},
				Worker: WorkerConfig{
					Counter: 20,
				},
				Job: JobConfig{
					Timeout:       10 * time.Second,
					MaxAttempts:   10,
					RetryDelay:    10 * time.Second,
					LeaseDuration: 10 * time.Second,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestLoadConfig(t *testing.T) {

	tests := []struct {
		name         string
		envVariables map[string]string
		wantConfig   Config
		wantErr      bool
	}{
		{
			name: "all env variables configurable",
			envVariables: map[string]string{

				string(Address):              "localhost",
				string(Port):                 "3300",
				string(PostgresDatabaseHost): "localhost",
				string(PostgresDB):           "database-name",
				string(PostgresPort):         "3000",
				string(PostgresUser):         "admin",
				string(PostgresPassword):     "password",

				string(WorkerCount): "4",

				string(JobTimeout):    "30s",
				string(RetryDelay):    "60s",
				string(MaxAttempts):   "5",
				string(LeaseDuration): "30m",
			},
			wantConfig: Config{
				Address: "localhost",
				Port:    3300,
				Database: DatabaseConfig{
					Host:         "localhost",
					DatabaseName: "database-name",
					Port:         3000,
					Username:     "admin",
					Password:     "password",
				},
				Worker: WorkerConfig{
					Counter: 4,
				},
				Job: JobConfig{
					Timeout:       30 * time.Second,
					RetryDelay:    60 * time.Second,
					MaxAttempts:   5,
					LeaseDuration: 30 * time.Minute,
				},
			},
			wantErr: false,
		},
		{
			name: "default values",
			envVariables: map[string]string{
				string(PostgresDatabaseHost): "localhost",
				string(PostgresDB):           "database-name",
				string(PostgresUser):         "admin",
				string(PostgresPassword):     "password",
			},
			wantConfig: Config{
				Address: "127.0.0.1",
				Port:    9000,
				Database: DatabaseConfig{
					Host:         "localhost",
					DatabaseName: "database-name",
					Port:         5432,
					Username:     "admin",
					Password:     "password",
				},
				Worker: WorkerConfig{
					Counter: 3,
				},
				Job: JobConfig{
					Timeout:       60 * time.Second,
					RetryDelay:    30 * time.Second,
					MaxAttempts:   3,
					LeaseDuration: 30 * time.Minute,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//set variables
			for key, value := range tt.envVariables {
				t.Setenv(key, value)
			}

			config, err := LoadConfig()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantConfig, *config)
		})
	}

}
