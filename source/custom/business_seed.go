package custom

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"gorm.io/gorm"
)

// businessSeedSQL 随代码分发，保证小组成员拉到代码后启动即可得到同一套菜单与演示数据。
//
//go:embed sql/business_seed.sql
var businessSeedSQL string

const (
	seedName          = "business_data"
	seedVersion       = 1
	initOrderBusiness = system.InitOrderExternal + 1
)

type businessSeed struct{}

func init() {
	system.RegisterInit(initOrderBusiness, &businessSeed{})
}

func (s *businessSeed) InitializerName() string { return "business_seed" }

func (s *businessSeed) MigrateTable(ctx context.Context) (context.Context, error) {
	db := dbFromContext(ctx)
	if db == nil {
		return ctx, fmt.Errorf("business seed: 缺少 ctx 中的 db")
	}
	return ctx, db.Exec(`CREATE TABLE IF NOT EXISTS sys_seed_versions (
		name        varchar(64) PRIMARY KEY,
		version     integer NOT NULL,
		applied_at  timestamptz NOT NULL DEFAULT now()
	)`).Error
}

func (s *businessSeed) TableCreated(ctx context.Context) bool {
	db := dbFromContext(ctx)
	return db != nil && db.Migrator().HasTable("sys_seed_versions")
}

// DataInserted 用版本表做幂等判断：同一版本只跑一次，避免重启覆盖手工修改的数据。
func (s *businessSeed) DataInserted(ctx context.Context) bool {
	db := dbFromContext(ctx)
	if db == nil {
		return false
	}
	var version int
	if err := db.Raw("SELECT version FROM sys_seed_versions WHERE name = ?", seedName).Scan(&version).Error; err != nil {
		return false
	}
	return version >= seedVersion
}

// InitializeData 在单个事务内逐条执行 SQL，任何一条失败都整体回滚，下次启动会重试。
func (s *businessSeed) InitializeData(ctx context.Context) (context.Context, error) {
	db := dbFromContext(ctx)
	if db == nil {
		return ctx, fmt.Errorf("business seed: 缺少 ctx 中的 db")
	}
	statements := splitSeedStatements(businessSeedSQL)
	err := db.Transaction(func(tx *gorm.DB) error {
		for i, stmt := range statements {
			if err := tx.Exec(stmt).Error; err != nil {
				return fmt.Errorf("第 %d 条语句失败: %w", i+1, err)
			}
		}
		return nil
	})
	if err != nil {
		return ctx, fmt.Errorf("business seed 执行失败: %w", err)
	}
	return ctx, nil
}

func dbFromContext(ctx context.Context) *gorm.DB {
	db, _ := ctx.Value("db").(*gorm.DB)
	return db
}

// splitSeedStatements 按分号切分并剔除注释行与空语句。
func splitSeedStatements(sql string) []string {
	var out []string
	for _, raw := range strings.Split(sql, ";") {
		var lines []string
		for _, line := range strings.Split(raw, "\n") {
			t := strings.TrimSpace(line)
			if t == "" || strings.HasPrefix(t, "--") {
				continue
			}
			lines = append(lines, line)
		}
		stmt := strings.TrimSpace(strings.Join(lines, "\n"))
		if stmt != "" {
			out = append(out, stmt)
		}
	}
	return out
}
