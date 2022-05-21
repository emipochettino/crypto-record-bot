package persistance

import (
	"CryptoRecordBot/internal/domain/model"
	"CryptoRecordBot/internal/domain/ports"
	"gorm.io/gorm"
)

type AlertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) ports.AlertRepository {
	return &AlertRepository{
		db: db,
	}
}

func (s *AlertRepository) FindByChatIDAndUserID(chatID int64, userID int64) ([]model.Alert, error) {
	var alertDAOS []AlertDAO
	result := s.db.Where("chat_id = ? AND user_id = ?",
		chatID,
		userID,
	).Find(&alertDAOS)

	var alerts []model.Alert

	for _, alertDAO := range alertDAOS {
		alerts = append(
			alerts,
			model.Alert{
				ChatId:        alertDAO.ChatId,
				UserId:        alertDAO.UserId,
				Currency:      alertDAO.Currency,
				IsGreaterThan: alertDAO.IsGreaterThan,
				Price:         alertDAO.Price,
				CreatedAt:     alertDAO.CreatedAt,
			})
	}

	return alerts, result.Error
}

func (s *AlertRepository) Create(alert model.Alert) error {
	result := s.db.Where(
		"chat_id = ? AND user_id = ? AND currency = ? AND is_greater_than = ?",
		alert.ChatId,
		alert.UserId,
		alert.Currency,
		alert.IsGreaterThan,
	).FirstOrCreate(&AlertDAO{
		alert.ChatId,
		alert.UserId,
		alert.Currency,
		alert.IsGreaterThan,
		alert.Price,
		alert.CreatedAt,
	})
	return result.Error
}

func (s *AlertRepository) Delete(alert model.Alert) (bool, error) {
	result := s.db.Where(
		"chat_id = ? AND user_id = ? AND currency = ?",
		alert.ChatId,
		alert.UserId,
		alert.Currency,
	).Delete(&AlertDAO{})

	return result.RowsAffected > 0, result.Error
}

func (s *AlertRepository) FindCurrencies() ([]string, error) {
	var alertDAOS []AlertDAO
	result := s.db.Distinct("currency").Find(&alertDAOS)

	var currencies []string
	for _, alertDAO := range alertDAOS {
		currencies = append(currencies, alertDAO.Currency)
	}

	return currencies, result.Error
}

func (s *AlertRepository) FindByCurrency(currency string) ([]model.Alert, error) {
	var alertDAOS []AlertDAO
	result := s.db.Where(
		"currency = ?",
		currency,
	).Find(&alertDAOS)

	var alerts []model.Alert

	for _, alertDAO := range alertDAOS {
		alerts = append(
			alerts,
			model.Alert{
				ChatId:        alertDAO.ChatId,
				UserId:        alertDAO.UserId,
				Currency:      alertDAO.Currency,
				IsGreaterThan: alertDAO.IsGreaterThan,
				Price:         alertDAO.Price,
				CreatedAt:     alertDAO.CreatedAt,
			})
	}

	return alerts, result.Error
}
