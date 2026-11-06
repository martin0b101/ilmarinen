package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	Address EnvKey = "ILMARINEN_ADDRESS"
	Port    EnvKey = "ILMARINEN_PORT"

	PostgresDatabaseHost EnvKey = "POSTGRES_HOST"
	PostgresDB           EnvKey = "POSTGRES_DB"
	PostgresUser         EnvKey = "POSTGRES_USER"
	PostgresPassword     EnvKey = "POSTGRES_PWD"
	PostgresPort         EnvKey = "POSTGRES_PORT"

	WorkerCount   EnvKey = "WORKER_COUNT"
	JobTimeout    EnvKey = "JOB_TIMEOUT"
	MaxAttempts   EnvKey = "MAX_ATTEMPTS"
	RetryDelay    EnvKey = "RETRY_DELAY"
	LeaseDuration EnvKey = "LEASE_DURATION"

	DefaultAddress       string = "127.0.0.1"
	DefaultPort          string = "9000"
	DefaultPostgresPort  string = "5432"
	DefaultWorkerCount   string = "3"
	DefaultJobTimeout    string = "60s"
	DefaultMaxAttempts   string = "3"
	DefaultRetryDelay    string = "30s"
	DefaultLeaseDuration string = "30m"
)

type EnvKey string

func (key EnvKey) GetValue() string {
	return os.Getenv(string(key))
}

func (key EnvKey) GetValueOrDefault(defaultValue string) string {
	envValue := key.GetValue()
	if envValue != "" {
		return envValue
	}
	return defaultValue
}

type Config struct {
	Address  string
	Port     int
	Database DatabaseConfig
	Worker   WorkerConfig
	Job      JobConfig
}

func (c Config) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("server address is required")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}

	dbError := c.Database.validate()
	if dbError != nil {
		return dbError
	}

	workerError := c.Worker.validation()
	if workerError != nil {
		return workerError
	}

	jobError := c.Job.validate()
	if jobError != nil {
		return jobError
	}
	return nil
}

type DatabaseConfig struct {
	Host         string
	DatabaseName string
	Username     string
	Password     string
	Port         int
}

func (dc DatabaseConfig) validate() error {
	if dc.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if dc.DatabaseName == "" {
		return fmt.Errorf("database name is required")
	}
	if dc.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if dc.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if dc.Port < 1 || dc.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	return nil
}

type WorkerConfig struct {
	Counter int
}

func (wc WorkerConfig) validation() error {
	if wc.Counter <= 0 {
		return fmt.Errorf("invalid value for counter")
	}
	return nil
}

type JobConfig struct {
	Timeout       time.Duration
	MaxAttempts   int
	RetryDelay    time.Duration
	LeaseDuration time.Duration
}

func (jc JobConfig) validate() error {
	if jc.Timeout <= 0 {
		return fmt.Errorf("job timeout must be greater than zero")
	}

	if jc.MaxAttempts <= 0 {
		return fmt.Errorf("max attempts must be greater than zero")
	}

	if jc.RetryDelay <= 0 {
		return fmt.Errorf("retry delay must be greater than zero")
	}

	if jc.LeaseDuration <= 0 {
		return fmt.Errorf("lease duration must be greater than zero")
	}
	return nil
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	port, err := strconv.Atoi(Port.GetValueOrDefault(DefaultPort))
	if err != nil {
		return nil, fmt.Errorf("parse ILMARINEN_PORT: %w", err)
	}

	databasePort, err := strconv.Atoi(PostgresPort.GetValueOrDefault(DefaultPostgresPort))
	if err != nil {
		return nil, fmt.Errorf("parse POSTGRES_PORT: %w", err)
	}

	workerCounter, err := strconv.Atoi(WorkerCount.GetValueOrDefault(DefaultWorkerCount))
	if err != nil {
		return nil, fmt.Errorf("parse WORKER_COUNT: %w", err)
	}

	jobTimeout, err := time.ParseDuration(JobTimeout.GetValueOrDefault(DefaultJobTimeout))
	if err != nil {
		return nil, fmt.Errorf("parse JOB_TIMEOUT: %w", err)
	}

	maxAttempts, err := strconv.Atoi(MaxAttempts.GetValueOrDefault(DefaultMaxAttempts))
	if err != nil {
		return nil, fmt.Errorf("parse MAX_ATTEMPTS: %w", err)
	}

	retryDelay, err := time.ParseDuration(RetryDelay.GetValueOrDefault(DefaultRetryDelay))
	if err != nil {
		return nil, fmt.Errorf("parse RETRY_DELAY: %w", err)
	}

	leaseDuration, err := time.ParseDuration(LeaseDuration.GetValueOrDefault(DefaultLeaseDuration))
	if err != nil {
		return nil, fmt.Errorf("parse LEASE_DURATION: %w", err)
	}

	cfg := &Config{
		Address: Address.GetValueOrDefault(DefaultAddress),
		Port:    port,
		Database: DatabaseConfig{
			Host:         PostgresDatabaseHost.GetValue(),
			DatabaseName: PostgresDB.GetValue(),
			Port:         databasePort,
			Username:     PostgresUser.GetValue(),
			Password:     PostgresPassword.GetValue(),
		},
		Worker: WorkerConfig{
			Counter: workerCounter,
		},
		Job: JobConfig{
			Timeout:       jobTimeout,
			MaxAttempts:   maxAttempts,
			RetryDelay:    retryDelay,
			LeaseDuration: leaseDuration,
		},
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
