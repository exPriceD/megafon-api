package app

import (
	"context"
	"database/sql"
	"log"
	"megafon-buisness-reports/internal/config"
	"megafon-buisness-reports/internal/infrastructure/db"
	zaplogger "megafon-buisness-reports/internal/infrastructure/logger/zap"
	"megafon-buisness-reports/internal/infrastructure/megafon"
	megsvc "megafon-buisness-reports/internal/infrastructure/megafon/services"
	"megafon-buisness-reports/internal/infrastructure/reporting"
	"megafon-buisness-reports/internal/infrastructure/repository"
	tg "megafon-buisness-reports/internal/infrastructure/telegram"
	"megafon-buisness-reports/internal/usecase/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Run(ctx context.Context) error {
	lg, err := zaplogger.New(a.cfg.Logging)
	if err != nil {
		return err
	}
	lg.Info(a.cfg.Postgres.DSN)
	pgPool, err := db.New(context.Background(), a.cfg.Postgres.ToDBConfig())
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	defer func(pgPool *sql.DB) {
		_ = pgPool.Close()
	}(pgPool)

	botAPI, err := tgbotapi.NewBotAPI(a.cfg.Telegram.Token)
	if err != nil {
		return err
	}

	mClient, err := megafon.NewClient(a.cfg.MegafonBuisness, lg)
	if err != nil {
		return err
	}

	callServiceRepo := megsvc.NewCallService(mClient, lg)
	cityPhoneRepo := repository.NewCityPhoneRepository(pgPool)
	builder := reporting.NewExcelBuilder()
	uc := services.NewReportService(callServiceRepo, cityPhoneRepo, builder)
	bot := tg.New(botAPI, uc, cityPhoneRepo, lg)

	return bot.Run(ctx)
}
