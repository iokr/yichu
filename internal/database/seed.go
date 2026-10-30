package database

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/iokr/yichu/internal/config"
	"github.com/iokr/yichu/internal/model"
)

type seedCategory struct {
	Name     string
	Children []string
}

// 设计文档 3.2 分类种子数据。
var seedCategories = []seedCategory{
	{Name: "衣服", Children: []string{"上衣", "T恤", "衬衫", "卫衣", "毛衣", "外套", "羽绒服", "连衣裙", "半身裙", "裤子", "牛仔裤", "内衣", "睡衣"}},
	{Name: "鞋子", Children: []string{"运动鞋", "板鞋", "帆布鞋", "高跟鞋", "平底鞋", "靴子", "凉鞋", "拖鞋"}},
	{Name: "帽子", Children: []string{"棒球帽", "渔夫帽", "贝雷帽", "毛线帽", "遮阳帽"}},
}

// Seed 幂等写入分类、schema_version 与首个管理员账号。
func Seed(db *gorm.DB, cfg *config.Config) error {
	if err := seedCategoryTree(db); err != nil {
		return err
	}
	if err := db.Where("`key` = ?", model.SettingSchemaVersion).
		FirstOrCreate(&model.Setting{Key: model.SettingSchemaVersion, Value: model.SchemaVersion}).Error; err != nil {
		return err
	}
	return ensureAdmin(db, cfg)
}

func seedCategoryTree(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for i, top := range seedCategories {
			parent := model.Category{Name: top.Name, ParentID: 0}
			if err := tx.Where("parent_id = 0 AND name = ?", top.Name).
				Attrs(model.Category{Sort: i, IsSystem: true}).
				FirstOrCreate(&parent).Error; err != nil {
				return err
			}
			for j, child := range top.Children {
				sub := model.Category{Name: child, ParentID: parent.ID}
				if err := tx.Where("parent_id = ? AND name = ?", parent.ID, child).
					Attrs(model.Category{Sort: j, IsSystem: true}).
					FirstOrCreate(&sub).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func ensureAdmin(db *gorm.DB, cfg *config.Config) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Auth.InitPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := model.User{
		Username:     cfg.Auth.InitUsername,
		PasswordHash: string(hash),
		Nickname:     cfg.App.Name + "主人",
	}
	if err := db.Create(&user).Error; err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		return err
	}
	log.Printf("👤 已创建初始账号: %s", cfg.Auth.InitUsername)
	return nil
}
