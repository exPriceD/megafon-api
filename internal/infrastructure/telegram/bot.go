package telegram

import (
	"context"
	"megafon-buisness-reports/internal/domain/ports"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"megafon-buisness-reports/internal/interfaces"
	"megafon-buisness-reports/internal/usecase/services"
)

type Bot struct {
	api *tgbotapi.BotAPI
	h   *handler
}

func New(api *tgbotapi.BotAPI, rs *services.ReportService, cp ports.CityPhoneRepository, lg interfaces.Logger) *Bot {
	return &Bot{api: api, h: newHandler(api, rs, cp, lg)}
}

func (b *Bot) Run(ctx context.Context) error {
	offset := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			u := tgbotapi.NewUpdate(offset)
			u.Timeout = 30

			updates, err := b.api.GetUpdates(u)
			if err != nil {
				b.h.log.Warn("Ошибка получения обновлений: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, upd := range updates {
				offset = upd.UpdateID + 1

				if upd.CallbackQuery != nil {
					b.h.onCallback(upd.CallbackQuery)
				} else if upd.Message != nil {
					b.h.onMessage(upd.Message)
				}
			}
		}
	}
}
