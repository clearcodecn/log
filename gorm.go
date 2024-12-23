package xlogger

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// gorm log.
const KeyLoggerStatusBegin = "LoggerStatusBegin"

type statusBegin struct {
	begin time.Time
	op    string
}

func NewStatusBegin(begin time.Time, op string) *statusBegin {
	return &statusBegin{begin: begin, op: op}
}

type loggerPlugin struct {
}

func NewLoggerPlugin() gorm.Plugin {
	return &loggerPlugin{}
}

func (p *loggerPlugin) Name() string {
	return "logger"
}

// Initialize registers all needed callbacks
func (p *loggerPlugin) Initialize(db *gorm.DB) (err error) {
	_ = db.Callback().Create().Before("gorm:create").Register("logger:before_create", p.before("insert"))
	_ = db.Callback().Create().After("gorm:create").Register("logger:after_create", p.after("insert"))
	_ = db.Callback().Update().Before("gorm:update").Register("logger:before_update", p.before("update"))
	_ = db.Callback().Update().After("gorm:update").Register("logger:after_update", p.after("update"))
	_ = db.Callback().Query().Before("gorm:query").Register("logger:before_query", p.before("select"))
	_ = db.Callback().Query().After("gorm:query").Register("logger:after_query", p.after("select"))
	_ = db.Callback().Delete().Before("gorm:delete").Register("logger:before_delete", p.before("delete"))
	_ = db.Callback().Delete().After("gorm:delete").Register("logger:after_delete", p.after("delete"))
	_ = db.Callback().Row().Before("gorm:row").Register("logger:before_row", p.before("row"))
	_ = db.Callback().Row().After("gorm:row").Register("logger:after_row", p.after("row"))
	_ = db.Callback().Raw().Before("gorm:raw").Register("logger:before_raw", p.before("raw"))
	_ = db.Callback().Raw().After("gorm:raw").Register("logger:after_raw", p.after("raw"))
	return
}

func (p *loggerPlugin) before(op string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db == nil || db.Statement == nil || db.Statement.Context == nil {
			return
		}
		before := NewStatusBegin(time.Now(), op)
		db.InstanceSet(KeyLoggerStatusBegin, before)
	}
}

func (p *loggerPlugin) after(op string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db == nil || db.Statement == nil || db.Statement.Context == nil {
			return
		}
		v, ok := db.InstanceGet(KeyLoggerStatusBegin)
		if !ok || v == nil {
			return
		}
		before, ok := v.(*statusBegin)
		if !ok || before == nil || before.op != op {
			return
		}
		latency := time.Since(before.begin)
		if db.Error != nil && !p.isErrorIgnorable(db.Error) {
			sql := formatSQL(db.Statement.SQL.String(), db.Statement.Vars)
			WithContext(db.Statement.Context).WithFields(Field{
				"scene":    "mysql",
				"table":    db.Statement.Table,
				"op":       op,
				"duration": latency,
				"sql":      sql,
				"error":    db.Error,
				"args":     db.Statement.Vars,
			}).Error("mysql exec failed")
		} else {
			sql := formatSQL(db.Statement.SQL.String(), db.Statement.Vars)
			WithContext(db.Statement.Context).WithFields(Field{
				"scene":    "mysql",
				"table":    db.Statement.Table,
				"op":       op,
				"duration": latency,
				"sql":      sql,
				"rows":     db.RowsAffected,
			}).Info("mysql exec success")
		}
	}
}

func (p *loggerPlugin) isErrorIgnorable(err error) bool {
	if err == gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func formatSQL(sql string, args []interface{}) string {
	for _, arg := range args {
		idx := strings.Index(sql, "?")
		if idx == -1 {
			break
		}
		sql = sql[:idx] + fmt.Sprintf(`%s`, argString(arg)) + sql[idx+1:]
	}
	return sql
}

func argString(v interface{}) string {
	switch x := v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", x)
	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(x), 'g', -1, 32)
	default:
		return fmt.Sprintf("'%v'", x)
	}
}
