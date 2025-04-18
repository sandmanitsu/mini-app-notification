package event

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mini-app-notifications/internal/domain"
	sl "mini-app-notifications/internal/logger"
)

const (
	eventCreateItem string = "create_item"
)

type TelegramBot interface {
	SendMessage(users []domain.User, text string)
}

type UserRepository interface {
	GetUsers() ([]domain.User, error)
}

type EventService struct {
	logger   *slog.Logger
	bot      TelegramBot
	userRepo UserRepository
}

func NewEventService(logger *slog.Logger, tgbot TelegramBot, ur UserRepository) *EventService {
	return &EventService{
		logger:   logger,
		bot:      tgbot,
		userRepo: ur,
	}
}

func (e *EventService) Process(event domain.Event) error {
	if event.EventType == eventCreateItem {
		err := e.itemCreateNotification(event)
		if err != nil {
			return err
		}
	}

	return nil
}

type createItemPayload struct {
	ItemName  string `json:"item_name"`
	BrandName string `json:"brand_name"`
	ItemId    int    `json:"item_id"`
	Price     int    `json:"price"`
}

func (e *EventService) itemCreateNotification(event domain.Event) error {
	const op = "service.event.itemCreateNotification"

	var itemPayload createItemPayload
	if err := json.Unmarshal(event.Value, &itemPayload); err != nil {
		e.logger.Error(fmt.Sprintf("%s : failed unmashal payload", op), sl.Err(err))
		return err
	}

	users, err := e.userRepo.GetUsers()
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Новый товар! От %s %s за %d", itemPayload.BrandName, itemPayload.ItemName, itemPayload.Price)
	e.bot.SendMessage(users, msg)

	return nil
}
