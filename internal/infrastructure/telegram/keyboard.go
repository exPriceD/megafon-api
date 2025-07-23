package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func mainMenu() tgbotapi.InlineKeyboardMarkup {
	btn := tgbotapi.NewInlineKeyboardButtonData("📊 Получить отчёт", "get_report")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))
}

func cityMenuDynamic(cities []string) tgbotapi.InlineKeyboardMarkup {
	const perRow = 2

	var rows [][]tgbotapi.InlineKeyboardButton

	for i, c := range cities {
		btn := tgbotapi.NewInlineKeyboardButtonData(c, "city:"+c)

		if i%perRow == 0 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
		} else {
			rows[len(rows)-1] = append(rows[len(rows)-1], btn)
		}
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func periodMenu() tgbotapi.InlineKeyboardMarkup {
	row := func(btns ...tgbotapi.InlineKeyboardButton) []tgbotapi.InlineKeyboardButton { return btns }
	return tgbotapi.NewInlineKeyboardMarkup(
		row(tgbotapi.NewInlineKeyboardButtonData("Сегодня", "today"),
			tgbotapi.NewInlineKeyboardButtonData("Вчера", "yesterday")),
		row(tgbotapi.NewInlineKeyboardButtonData("Текущая неделя", "this_week"),
			tgbotapi.NewInlineKeyboardButtonData("Прошлая неделя", "last_week")),
		row(tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "this_month"),
			tgbotapi.NewInlineKeyboardButtonData("Прошлый месяц", "last_month")),
		row(tgbotapi.NewInlineKeyboardButtonData("Произвольный период", "period")),
	)
}
