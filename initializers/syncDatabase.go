package initializers

import "github.com/saffer4u/udharam/v2/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{}, &models.Ledger{}, &models.Transaction{})

}
