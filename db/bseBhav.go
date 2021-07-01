package db

import (
	models "github.com/g33kzone/go-stox-ingest/models/entity"
)

// UploadBSEBhavCopy - Insert Daily Bhavcopy in DB
func (db *Database) UploadBSEBhavCopy(bse_bhav []models.BSE_BHAV) error {

	err := db.DbConn.Model(&models.BSE_BHAV{}).Create(bse_bhav).Error
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) FetchData() ([]models.Order, error) {
	var orders []models.Order

	err := db.DbConn.Debug().Table("order").Where("date > ?", "2021-01-01 00:00:00+05:30").Find(&orders).Error
	if err != nil {
		return orders, err
	}
	return orders, nil
}
