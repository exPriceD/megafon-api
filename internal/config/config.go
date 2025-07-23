package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"megafon-buisness-reports/internal/infrastructure/db"
	"os"
	"strconv"
	"sync"
	"time"
)

// Config — корневой объект конфигурации.
type Config struct {
	Logging         LoggingConfig         `yaml:"logging"          validate:"required"`
	Telegram        TelegramConfig        `yaml:"telegram"         validate:"required"`
	MegafonBuisness MegafonBuisnessConfig `yaml:"megafon_buisness" validate:"required"`
	Postgres        PostgresConfig        `yaml:"postgres"         validate:"required"`
}

type LoggingConfig struct {
	Level      string `yaml:"level"       validate:"required,oneof=debug info warn error dpanic panic fatal"`
	Format     string `yaml:"format"      validate:"required,oneof=json console"`
	OutputPath string `yaml:"output_path" validate:"required"`
	Service    string `yaml:"service,omitempty"`
	Env        string `yaml:"env,omitempty"`
}

type TelegramConfig struct {
	Token string `validate:"required"`
}

type MegafonBuisnessConfig struct {
	BaseURL string `validate:"required"`
	APIKey  string `validate:"required"`
}

type PostgresConfig struct {
	DSN             string        `yaml:"dsn"               validate:"required"`
	MaxOpenConns    int           `yaml:"max_open_conns"    validate:"gte=1"`
	MaxIdleConns    int           `yaml:"max_idle_conns"    validate:"gte=0"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" validate:"gt=0"`
	MigrationsDir   string        `yaml:"migrations_dir"`
}

func (p PostgresConfig) ToDBConfig() db.Config {
	return db.Config{
		DSN:             p.DSN,
		MaxOpenConns:    p.MaxOpenConns,
		MaxIdleConns:    p.MaxIdleConns,
		ConnMaxLifetime: p.ConnMaxLifetime,
		MigrationsDir:   p.MigrationsDir,
	}
}

var (
	cfg  *Config
	once sync.Once
)

// Load читает .env, YAML, подставляет секреты, валидирует и
// возвращает singleton-конфиг. Повторные вызовы отдадут уже готовый объект.
func Load(yamlPath, envPath string) (*Config, error) {
	var err error
	once.Do(func() {
		if e := godotenv.Load(envPath); e != nil && !os.IsNotExist(e) {
			err = fmt.Errorf("load env: %w", e)
			return
		}
		cfg, err = readYAML(yamlPath)
		if err != nil {
			return
		}

		overrideFromEnv(&cfg.Telegram.Token, "TELEGRAM_TOKEN")
		overrideFromEnv(&cfg.MegafonBuisness.APIKey, "MEGAFON_API_KEY")
		overrideFromEnv(&cfg.MegafonBuisness.BaseURL, "MEGAFON_BASE_URL")
		overrideFromEnv(&cfg.Postgres.DSN, "POSTGRES_DSN")
		overrideIntFromEnv(&cfg.Postgres.MaxOpenConns, "POSTGRES_MAX_OPEN")
		overrideIntFromEnv(&cfg.Postgres.MaxIdleConns, "POSTGRES_MAX_IDLE")
		overrideDurationFromEnv(&cfg.Postgres.ConnMaxLifetime, "POSTGRES_CONN_LIFETIME")
		overrideFromEnv(&cfg.Postgres.MigrationsDir, "POSTGRES_MIGRATIONS")

		err = validate(cfg)
	})
	return cfg, err
}

// MustLoad — обёртка для main(): паникует при ошибке.
func MustLoad(yamlPath, envPath string) *Config {
	c, err := Load(yamlPath, envPath)
	if err != nil {
		panic(err)
	}
	return c
}

func readYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %q: %w", path, err)
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}
	return &c, nil
}

func validate(c *Config) error {
	return validator.New().Struct(c)
}

func overrideFromEnv(field *string, key string) {
	if v := os.Getenv(key); v != "" {
		*field = v
	}
}

func overrideIntFromEnv(field *int, key string) {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			*field = n
		}
	}
}

func overrideDurationFromEnv(field *time.Duration, key string) {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			*field = d
		}
	}
}
