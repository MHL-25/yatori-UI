package dao

import (
	"errors"
	"time"
	"yatori-UI/internal/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var GlobalDB *gorm.DB

func InitDB(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&entity.UserPO{}, &entity.SettingPO{})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	GlobalDB = db
	return nil
}

func GetSetting(key string) string {
	if GlobalDB == nil {
		return ""
	}
	var s entity.SettingPO
	if err := GlobalDB.Where("key = ?", key).First(&s).Error; err != nil {
		return ""
	}
	return s.Value
}

func SetSetting(key, value string) error {
	if GlobalDB == nil {
		return errors.New("database not initialized")
	}
	var s entity.SettingPO
	if err := GlobalDB.Where("key = ?", key).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GlobalDB.Create(&entity.SettingPO{Key: key, Value: value}).Error
		}
		return err
	}
	return GlobalDB.Model(&entity.SettingPO{}).Where("key = ?", key).Update("value", value).Error
}

func InsertUser(user *entity.UserPO) error {
	if GlobalDB == nil {
		return errors.New("database not initialized")
	}
	if err := GlobalDB.Create(&user).Error; err != nil {
		return errors.New("插入数据失败: " + err.Error())
	}
	return nil
}

func DeleteUser(uid string) error {
	if GlobalDB == nil {
		return errors.New("database not initialized")
	}
	if err := GlobalDB.Where("uid = ?", uid).Delete(&entity.UserPO{}).Error; err != nil {
		return errors.New("删除用户失败: " + err.Error())
	}
	return nil
}

func QueryAllUsers() ([]entity.UserPO, error) {
	if GlobalDB == nil {
		return nil, errors.New("database not initialized")
	}
	var users []entity.UserPO
	if err := GlobalDB.Order("uid ASC").Find(&users).Error; err != nil {
		return nil, errors.New("查询用户失败: " + err.Error())
	}
	return users, nil
}

func QueryUserByUid(uid string) (*entity.UserPO, error) {
	if GlobalDB == nil {
		return nil, errors.New("database not initialized")
	}
	var user entity.UserPO
	if err := GlobalDB.Where("uid = ?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("查询用户失败: " + err.Error())
	}
	return &user, nil
}

func QueryUserByAccount(accountType, account string) (*entity.UserPO, error) {
	if GlobalDB == nil {
		return nil, errors.New("database not initialized")
	}
	var user entity.UserPO
	if err := GlobalDB.Where("account_type = ? AND account = ?", accountType, account).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("查询用户失败: " + err.Error())
	}
	return &user, nil
}

func UpdateUser(uid string, updateData map[string]interface{}) error {
	if GlobalDB == nil {
		return errors.New("database not initialized")
	}
	if uid == "" {
		return errors.New("UID 不能为空")
	}
	if err := GlobalDB.Model(&entity.UserPO{}).Where("uid = ?", uid).Updates(updateData).Error; err != nil {
		return errors.New("更新用户失败: " + err.Error())
	}
	return nil
}
