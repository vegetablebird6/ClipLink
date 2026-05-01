package db

import "fmt"

// EnsureIndexes 创建或重建所有应用层索引。
// 索引不依赖 GORM struct tag（struct 字段声明顺序决定列顺序，不可控），
// 全部由本函数以 raw SQL 精确控制列顺序。
// 必须在 MigrateDB() 的 AutoMigrate 之后调用。
func EnsureIndexes() error {
	dialect := GetDB().Dialector.Name()

	type indexDef struct {
		name    string
		table   string
		columns string
		unique  bool
	}

	indexes := []indexDef{
		// clipboard_items
		{"idx_clip_created", "clipboard_items", "channel_id, created_at, id", false},
		{"idx_clip_type", "clipboard_items", "channel_id, type, created_at, id", false},
		{"idx_clip_fav", "clipboard_items", "channel_id, favorite, updated_at, id", false},
		{"idx_channel_content_hash", "clipboard_items", "channel_id, content_hash", false},
		// sync_events
		{"idx_sync_channel_created", "sync_events", "channel_id, created_at, id", false},
		// device_channels
		{"idx_device_channel", "device_channels", "device_id, channel_id", true},
		{"idx_channel_last_seen", "device_channels", "channel_id, last_seen_at", false},
	}

	for _, idx := range indexes {
		if err := createIndex(dialect, idx.name, idx.table, idx.columns, idx.unique); err != nil {
			return fmt.Errorf("创建索引 %s 失败: %w", idx.name, err)
		}
	}

	return nil
}

func createIndex(dialect, name, table, columns string, unique bool) error {
	db := GetDB()
	uniqueKeyword := ""
	if unique {
		uniqueKeyword = "UNIQUE "
	}

	switch dialect {
	case "sqlite":
		// SQLite 支持 DROP INDEX IF EXISTS，失败忽略；然后 CREATE INDEX。
		_ = db.Exec(fmt.Sprintf("DROP INDEX IF EXISTS %s", name)).Error
		return db.Exec(fmt.Sprintf("CREATE %sINDEX %s ON %s (%s)", uniqueKeyword, name, table, columns)).Error

	case "mysql":
		// 查询 information_schema 确认索引是否存在，再决定是否 DROP。
		// 兼容 MySQL 5.7+ / MariaDB 10.x。
		var count int
		err := db.Raw(
			"SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?",
			table, name,
		).Scan(&count).Error
		if err != nil {
			return fmt.Errorf("查询索引信息失败: %w", err)
		}
		if count > 0 {
			if dropErr := db.Exec(fmt.Sprintf("ALTER TABLE %s DROP INDEX %s", table, name)).Error; dropErr != nil {
				return fmt.Errorf("删除旧索引失败: %w", dropErr)
			}
		}
		return db.Exec(fmt.Sprintf("CREATE %sINDEX %s ON %s (%s)", uniqueKeyword, name, table, columns)).Error

	default:
		return fmt.Errorf("不支持的数据库方言: %s", dialect)
	}
}
