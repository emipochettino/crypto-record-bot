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
				CoinName:      alertDAO.CoinName,
				IsGreaterThan: alertDAO.IsGreaterThan,
				Price:         alertDAO.Price,
				CreatedAt:     alertDAO.CreatedAt,
			})
	}

	return alerts, result.Error
}

func (s *AlertRepository) Create(alert model.Alert) error {
	result := s.db.Where(
		"chat_id = ? AND user_id = ? AND coin_name = ? AND is_greater_than = ?",
		alert.ChatId,
		alert.UserId,
		alert.CoinName,
		alert.IsGreaterThan,
	).FirstOrCreate(&AlertDAO{
		alert.ChatId,
		alert.UserId,
		alert.CoinName,
		alert.IsGreaterThan,
		alert.Price,
		alert.CreatedAt,
	})
	return result.Error
}

func (s *AlertRepository) Delete(alert model.Alert) (bool, error) {
	result := s.db.Where(
		"chat_id = ? AND user_id = ? AND coin_name = ?",
		alert.ChatId,
		alert.UserId,
		alert.CoinName,
	).Delete(&AlertDAO{})

	return result.RowsAffected > 0, result.Error
}

func (s *AlertRepository) FindCoinNames() ([]string, error) {
	var alertDAOS []AlertDAO
	result := s.db.Distinct("coin_name").Find(&alertDAOS)

	var currencies []string
	for _, alertDAO := range alertDAOS {
		currencies = append(currencies, alertDAO.CoinName)
	}

	return currencies, result.Error
}

func (s *AlertRepository) FindByCoinName(coinName string) ([]model.Alert, error) {
	var alertDAOS []AlertDAO
	result := s.db.Where(
		"coin_name = ?",
		coinName,
	).Find(&alertDAOS)

	var alerts []model.Alert

	for _, alertDAO := range alertDAOS {
		alerts = append(
			alerts,
			model.Alert{
				ChatId:        alertDAO.ChatId,
				UserId:        alertDAO.UserId,
				CoinName:      alertDAO.CoinName,
				IsGreaterThan: alertDAO.IsGreaterThan,
				Price:         alertDAO.Price,
				CreatedAt:     alertDAO.CreatedAt,
			})
	}

	return alerts, result.Error
}
