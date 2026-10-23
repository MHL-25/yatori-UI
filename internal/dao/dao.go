package dao

import (
	"errors"
	"time"
	"yatori-UI/internal/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GlobalDB *gorm.DB

func InitDB(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath+"?_journal_mode=WAL&_busy_timeout=5000"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&entity.UserPO{}, &entity.SettingPO{}, &entity.CourseLockPO{}, &entity.AccountCheckPO{})
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

// ReplaceCourseLocks 覆盖某账号某门课程的待处理课程锁记录
func ReplaceCourseLocks(uid, courseID string, locks []entity.CourseLockPO) error {
	if GlobalDB == nil {
		return errors.New("database not initialized")
	}
	return GlobalDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("uid = ? AND course_id = ? AND handled = ?", uid, courseID, false).Delete(&entity.CourseLockPO{}).Error; err != nil {
			return err
		}
		if len(locks) == 0 {
			return nil
		}
		return tx.Create(&locks).Error
	})
}

// QueryDueCourseLocks 查询已到解锁时间且未处理的课程锁
func QueryDueCourseLocks(now time.Time) ([]entity.CourseLockPO, error) {
	if GlobalDB == nil {
		return nil, errors.New("database not initialized")
	}
	var locks []entity.CourseLockPO
	if err := GlobalDB.Where("handled = ? AND unlock_at <= ?", false, now).Order("unlock_at ASC").Find(&locks).Error; err != nil {
		return nil, err
	}
	return locks, nil
}

// QueryPendingCourseLocks 查询某账号未处理的课程锁
func QueryPendingCourseLocks(uid string) ([]entity.CourseLockPO, error) {
	if GlobalDB == nil {
		return nil, errors.New("database not initialized")
	}
	var locks []entity.CourseLockPO
	if err := GlobalDB.Where("uid = ? AND handled = ?", uid, false).Order("unlock_at ASC").Find(&locks).Error; err != nil {
		return nil, err
	}
	return locks, nil
}

// MarkCourseLockHandled 标记课程锁已处理
func MarkCourseLockHandled(id uint) error {
	if GlobalDB == nil {
		return errors.New("database not initialized")
	}
	return GlobalDB.Model(&entity.CourseLockPO{}).Where("id = ?", id).Update("handled", true).Error
}

// SaveCheckReport 保存账号最近一次质检报告
func SaveCheckReport(report *entity.CheckReport, summary string) error {
	if GlobalDB == nil {
		return errors.New("database not initialized")
	}
	po := entity.AccountCheckPO{Uid: report.Uid, CheckedAt: report.CheckedAt, Summary: summary}
	return GlobalDB.Save(&po).Error
}

// GetCheckReport 读取账号最近一次质检报告 JSON
func GetCheckReport(uid string) (string, error) {
	if GlobalDB == nil {
		return "", errors.New("database not initialized")
	}
	var po entity.AccountCheckPO
	if err := GlobalDB.Where("uid = ?", uid).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return po.Summary, nil
}

// GetAllCheckReports 读取所有账号最近一次质检报告
func GetAllCheckReports() ([]entity.AccountCheckPO, error) {
	if GlobalDB == nil {
		return nil, errors.New("database not initialized")
	}
	var list []entity.AccountCheckPO
	if err := GlobalDB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
