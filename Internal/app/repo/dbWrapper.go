package repo

import (
	"github.com/Pave1kor/Optics-characteristics/internal/app/models"
)

// DBWrapper оборачивает DBManager и добавляет методы
type DBWrapper struct {
	*models.DBManager
}

// NewDBWrapper - конструктор
func NewDBWrapper(manager *models.DBManager) *DBWrapper {
	return &DBWrapper{manager}
}
