package custom

import (
	"strings"
	"testing"
)

// TestSplitSeedStatements 回归：seed SQL 按分号切分时要过滤注释与空块，
// 否则注释里的分号会产生空语句，导致 db.Exec("") 报错。
func TestSplitSeedStatements(t *testing.T) {
	sql := "-- comment only;\nCREATE TABLE t (id int);\n\nINSERT INTO t VALUES (1);\n-- tail comment\n"
	got := splitSeedStatements(sql)
	if len(got) != 2 {
		t.Fatalf("切分出 %d 条语句，期望 2 条：%#v", len(got), got)
	}
	if got[0] != "CREATE TABLE t (id int)" {
		t.Errorf("第 1 条语句 = %q", got[0])
	}
	if got[1] != "INSERT INTO t VALUES (1)" {
		t.Errorf("第 2 条语句 = %q", got[1])
	}
}

// TestBusinessSeedSQLStructure 回归：随代码分发的 seed 必须同时覆盖
// 菜单规整、角色授权、客服演示数据与版本标记，缺一会让组员拉代码后数据不一致。
func TestBusinessSeedSQLStructure(t *testing.T) {
	stmts := splitSeedStatements(businessSeedSQL)
	if len(stmts) < 20 {
		t.Fatalf("seed 语句数 = %d，期望不少于 20 条", len(stmts))
	}
	joined := strings.Join(stmts, "\n")
	for _, want := range []string{
		"sys_seed_versions", "sys_base_menus", "sys_authority_menus",
		"cs_faq_categories", "cs_faqs", "cs_sessions", "cs_tickets",
		"cs_ticket_logs", "cs_notifications", "business_data",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("seed 缺少关键内容 %q", want)
		}
	}
	if !strings.Contains(joined, "'business_data', 1") {
		t.Error("seed 末尾未写入与 seedVersion 一致的版本号")
	}
}
