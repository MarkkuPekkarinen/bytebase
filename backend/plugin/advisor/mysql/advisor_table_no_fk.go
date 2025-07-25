package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/pkg/errors"

	mysql "github.com/bytebase/mysql-parser"

	"github.com/bytebase/bytebase/backend/common"
	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	"github.com/bytebase/bytebase/backend/plugin/advisor"
	mysqlparser "github.com/bytebase/bytebase/backend/plugin/parser/mysql"
)

var (
	_ advisor.Advisor = (*TableNoFKAdvisor)(nil)
)

func init() {
	advisor.Register(storepb.Engine_MYSQL, advisor.MySQLTableNoFK, &TableNoFKAdvisor{})
	advisor.Register(storepb.Engine_MARIADB, advisor.MySQLTableNoFK, &TableNoFKAdvisor{})
	advisor.Register(storepb.Engine_OCEANBASE, advisor.MySQLTableNoFK, &TableNoFKAdvisor{})
}

// TableNoFKAdvisor is the advisor checking table disallow foreign key.
type TableNoFKAdvisor struct {
}

// Check checks table disallow foreign key.
func (*TableNoFKAdvisor) Check(_ context.Context, checkCtx advisor.Context) ([]*storepb.Advice, error) {
	root, ok := checkCtx.AST.([]*mysqlparser.ParseResult)
	if !ok {
		return nil, errors.Errorf("failed to convert to mysql parse result")
	}
	level, err := advisor.NewStatusBySQLReviewRuleLevel(checkCtx.Rule.Level)
	if err != nil {
		return nil, err
	}

	checker := &tableNoFKChecker{
		level: level,
		title: string(checkCtx.Rule.Type),
	}
	for _, stmtNode := range root {
		checker.baseLine = stmtNode.BaseLine
		antlr.ParseTreeWalkerDefault.Walk(checker, stmtNode.Tree)
	}

	return checker.adviceList, nil
}

type tableNoFKChecker struct {
	*mysql.BaseMySQLParserListener

	baseLine   int
	adviceList []*storepb.Advice
	level      storepb.Advice_Status
	title      string
}

// EnterCreateTable is called when production createTable is entered.
func (checker *tableNoFKChecker) EnterCreateTable(ctx *mysql.CreateTableContext) {
	if !mysqlparser.IsTopMySQLRule(&ctx.BaseParserRuleContext) {
		return
	}
	if ctx.TableName() == nil || ctx.TableElementList() == nil {
		return
	}

	_, tableName := mysqlparser.NormalizeMySQLTableName(ctx.TableName())
	for _, tableElement := range ctx.TableElementList().AllTableElement() {
		if tableElement.TableConstraintDef() == nil {
			continue
		}
		checker.handleTableConstraintDef(tableName, tableElement.TableConstraintDef())
	}
}

// EnterAlterTable is called when production alterTable is entered.
func (checker *tableNoFKChecker) EnterAlterTable(ctx *mysql.AlterTableContext) {
	if !mysqlparser.IsTopMySQLRule(&ctx.BaseParserRuleContext) {
		return
	}
	if ctx.TableRef() == nil {
		return
	}
	_, tableName := mysqlparser.NormalizeMySQLTableRef(ctx.TableRef())

	if ctx.AlterTableActions() == nil || ctx.AlterTableActions().AlterCommandList() == nil {
		return
	}
	if ctx.AlterTableActions().AlterCommandList().AlterList() == nil {
		return
	}
	for _, option := range ctx.AlterTableActions().AlterCommandList().AlterList().AllAlterListItem() {
		switch {
		// ADD CONSTRANIT
		case option.ADD_SYMBOL() != nil && option.TableConstraintDef() != nil:
			checker.handleTableConstraintDef(tableName, option.TableConstraintDef())
		default:
			continue
		}
	}
}

func (checker *tableNoFKChecker) handleTableConstraintDef(tableName string, ctx mysql.ITableConstraintDefContext) {
	if ctx.GetType_() != nil {
		switch strings.ToUpper(ctx.GetType_().GetText()) {
		case "FOREIGN":
			checker.adviceList = append(checker.adviceList, &storepb.Advice{
				Status:        checker.level,
				Code:          advisor.TableHasFK.Int32(),
				Title:         checker.title,
				Content:       fmt.Sprintf("Foreign key is not allowed in the table `%s`", tableName),
				StartPosition: common.ConvertANTLRLineToPosition(checker.baseLine + ctx.GetStart().GetLine()),
			})
		default:
		}
	}
}
