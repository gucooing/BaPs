package db

import (
	"errors"

	"gorm.io/gorm"
)

type BlackDevice struct {
	DeviceId string `gorm:"unique"`
}

// GetBlackDeviceByYostarUid 使用DeviceId拉取数据
func GetBlackDeviceByYostarUid(deviceId string) *BlackDevice {
	var data *BlackDevice
	err := SQL.Model(&BlackDevice{}).Where("device_id = ?", deviceId).First(&data).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil
	}
	return data
}
