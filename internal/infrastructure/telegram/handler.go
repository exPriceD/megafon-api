package telegram

import (
	"context"
	"fmt"
	"megafon-buisness-reports/internal/domain/ports"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"megafon-buisness-reports/internal/domain/entities"
	"megafon-buisness-reports/internal/interfaces"
	"megafon-buisness-reports/internal/usecase/services"
)

type handler struct {
	bot     *tgbotapi.BotAPI
	states  *stateStore
	rs      *services.ReportService
	cp      ports.CityPhoneRepository
	log     interfaces.Logger
	fromBuf map[int64]time.Time
	cityBuf map[int64]string
}

func newHandler(b *tgbotapi.BotAPI, rs *services.ReportService, cp ports.CityPhoneRepository, lg interfaces.Logger) *handler {
	return &handler{
		bot:     b,
		rs:      rs,
		cp:      cp,
		log:     lg,
		states:  newStateStore(),
		fromBuf: make(map[int64]time.Time),
		cityBuf: make(map[int64]string),
	}
}

func (h *handler) text(chat int64, txt string) {
	_, _ = h.bot.Send(tgbotapi.NewMessage(chat, txt))
}

func (h *handler) onMessage(m *tgbotapi.Message) {
	uid := m.Chat.ID

	switch h.states.get(uid) {
	case awaitFrom:
		from, err := parseRuDate(m.Text)
		fmt.Println(from)
		if err != nil {
			h.text(uid, err.Error())
			return
		}
		h.fromBuf[uid] = from
		h.states.set(uid, awaitTo)
		h.text(uid, "Введите дату окончания не включительно (дд.мм.гггг)")

	case awaitTo:
		to, err := parseRuDate(m.Text)
		if err != nil {
			h.text(uid, err.Error())
			return
		}
		fmt.Println(to)
		from := h.fromBuf[uid]
		h.states.set(uid, idle)
		filter := entities.CallFilter{
			Start:         &from,
			End:           &to,
			ProcessMissed: true,
		}
		h.makeReport(uid, filter, h.cityBuf[uid])

	default:
		// главное меню
		msg := tgbotapi.NewMessage(uid, "Главное меню")
		msg.ReplyMarkup = mainMenu()
		_, _ = h.bot.Send(msg)
	}
}

func (h *handler) onCallback(q *tgbotapi.CallbackQuery) {
	uid := q.Message.Chat.ID
	data := q.Data
	switch {
	// ───── 1. Главное меню ─────────────────────────────────────────────
	case data == "get_report":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cities, err := h.cp.GetCities(ctx)
		if err != nil {
			h.text(uid, fmt.Sprintf("Не удалось получить список городов: %v", err))
			break
		}

		h.states.set(uid, awaitCity)
		msg := tgbotapi.NewMessage(uid, "Выберите город")
		msg.ReplyMarkup = cityMenuDynamic(cities)
		_, _ = h.bot.Send(msg)

	// ───── 2. Пользователь выбрал город ────────────────────────────────
	case strings.HasPrefix(data, "city:"):
		city := strings.TrimPrefix(data, "city:")
		h.cityBuf[uid] = city

		h.states.set(uid, awaitPeriod)
		msg := tgbotapi.NewMessage(uid, "Выберите период")
		msg.ReplyMarkup = periodMenu()
		_, _ = h.bot.Send(msg)

	// ───── 3. Пользователь выбрал «Произвольный период» ────────────────
	case data == "period":
		if h.states.get(uid) != awaitPeriod {
			break
		}
		h.states.set(uid, awaitFrom)
		msg := tgbotapi.NewMessage(uid, "Введите дату начала (дд.мм.гггг)")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		_, _ = h.bot.Send(msg)

	// ───── 4. Остальные периоды («today», «yesterday» …) ───────────────
	default:
		if h.states.get(uid) != awaitPeriod {
			break
		}
		h.states.set(uid, idle)

		filter := entities.CallFilter{
			Period:        entities.Period(data),
			ProcessMissed: true,
		}
		h.makeReport(uid, filter, h.cityBuf[uid])
	}
	_, _ = h.bot.Request(tgbotapi.CallbackConfig{CallbackQueryID: q.ID})
}

func (h *handler) makeReport(chat int64, f entities.CallFilter, city string) {
	h.text(chat, "Генерирую отчёт…")
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	buf, err := h.rs.GenerateCallReport(ctx, f, city)
	if err != nil {
		h.text(chat, fmt.Sprintf("Ошибка: %v", err))
		return
	}
	doc := tgbotapi.FileBytes{Name: "report.xlsx", Bytes: buf.Bytes()}
	_, _ = h.bot.Send(tgbotapi.NewDocument(chat, doc))
}
