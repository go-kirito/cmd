package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/blastrain/vitess-sqlparser/tidbparser/ast"
	"github.com/blastrain/vitess-sqlparser/tidbparser/dependency/mysql"
	"github.com/blastrain/vitess-sqlparser/tidbparser/dependency/types"
	"github.com/blastrain/vitess-sqlparser/tidbparser/parser"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type ModelCodes struct {
	Package    string
	ImportPath []string
	StructCode tmplData
}

func ParseSql(sql string, options ...Option) ([]ModelCodes, error) {
	opt := parseOption(options)
	stmts, err := parser.New().Parse(sql, opt.Charset, opt.Collation)
	if err != nil {
		return nil, err
	}

	var codes []ModelCodes

	for _, stmt := range stmts {
		if ct, ok := stmt.(*ast.CreateTableStmt); ok {
			importPath := make(map[string]struct{})
			s, ipt, err := makeCode(ct, opt)
			if err != nil {
				return nil, err
			}
			for _, s := range ipt {
				importPath[s] = struct{}{}
			}
			importPathArr := make([]string, 0, len(importPath))
			for s := range importPath {
				importPathArr = append(importPathArr, s)
			}
			sort.Strings(importPathArr)
			modelCode := ModelCodes{
				Package:    opt.Package,
				ImportPath: importPathArr,
				StructCode: s,
			}
			codes = append(codes, modelCode)
		}
	}
	return codes, nil
}

type tmplData struct {
	TableName    string
	NameFunc     bool
	RawTableName string
	Fields       []tmplField
	Comment      string
}

type tmplField struct {
	Name    string
	GoType  string
	Tag     string
	Comment string
}

func makeCode(stmt *ast.CreateTableStmt, opt options) (tmplData, []string, error) {
	importPath := make([]string, 0, 1)
	data := tmplData{
		TableName:    stmt.Table.Name.String(),
		RawTableName: stmt.Table.Name.String(),
		Fields:       make([]tmplField, 0, 1),
	}
	tablePrefix := opt.TablePrefix
	if tablePrefix != "" && strings.HasPrefix(data.TableName, tablePrefix) {
		data.NameFunc = true
		data.TableName = data.TableName[len(tablePrefix):]
	}
	if opt.ForceTableName || data.RawTableName != inflection.Plural(data.RawTableName) {
		data.NameFunc = true
	}

	data.TableName = strcase.ToCamel(data.TableName)

	// find table comment
	for _, opt := range stmt.Options {
		if opt.Tp == ast.TableOptionComment {
			data.Comment = opt.StrValue
			break
		}
	}

	isPrimaryKey := make(map[string]bool)
	isIndex := make(map[string][]string)
	isUniq := make(map[string][]string)
	for _, con := range stmt.Constraints {
		switch con.Tp {
		case ast.ConstraintPrimaryKey:
			isPrimaryKey[con.Keys[0].Column.String()] = true
		case ast.ConstraintIndex:
			for _, item := range con.Keys {
				index, ok := isIndex[item.Column.String()]
				if !ok {
					index = make([]string, 0)
				}
				index = append(index, con.Name)
				isIndex[item.Column.String()] = index
			}
		case ast.ConstraintUniq:
			for _, item := range con.Keys {
				index, ok := isUniq[item.Column.String()]
				if !ok {
					index = make([]string, 0)
				}
				index = append(index, con.Name)
				isUniq[item.Column.String()] = index
			}
		default:
			fmt.Println("未处理类型 con.tp:", con.Tp, "column:", con.Keys[0].Column, con.Name)
		}
	}

	columnPrefix := opt.ColumnPrefix
	for _, col := range stmt.Cols {
		colName := col.Name.Name.String()
		goFieldName := colName
		if columnPrefix != "" && strings.HasPrefix(goFieldName, columnPrefix) {
			goFieldName = goFieldName[len(columnPrefix):]
		}

		field := tmplField{
			Name: strcase.ToCamel(goFieldName),
		}

		tags := make([]string, 0, 4)
		// make GORM's tag
		gormTag := strings.Builder{}
		gormTag.WriteString("column:")
		gormTag.WriteString(colName)
		if opt.GormType {
			gormTag.WriteString(";type:")
			gormTag.WriteString(col.Tp.InfoSchemaStr())
		}
		if isPrimaryKey[colName] {
			gormTag.WriteString(";primary_key")
		}

		if indexs, ok := isIndex[colName]; ok {
			for _, indexName := range indexs {
				gormTag.WriteString(";index:")
				gormTag.WriteString(indexName)
			}
		}

		if uniqName, ok := isUniq[colName]; ok {
			for _, indexName := range uniqName {
				gormTag.WriteString(";uniqueIndex:")
				gormTag.WriteString(indexName)
			}
		}

		isNotNull := false
		canNull := false
		for _, o := range col.Options {
			switch o.Tp {
			case ast.ColumnOptionPrimaryKey:
				if !isPrimaryKey[colName] {
					gormTag.WriteString(";primary_key")
					isPrimaryKey[colName] = true
				}
			case ast.ColumnOptionNotNull:
				isNotNull = true
			case ast.ColumnOptionAutoIncrement:
				gormTag.WriteString(";AUTO_INCREMENT")
			case ast.ColumnOptionDefaultValue:
				if value := getDefaultValue(o.Expr); value != "" {
					gormTag.WriteString(";default:")
					gormTag.WriteString(value)
				}
			case ast.ColumnOptionUniqKey:
				gormTag.WriteString(";unique")
			case ast.ColumnOptionNull:
				//gormTag.WriteString(";NULL")
				canNull = true
			case ast.ColumnOptionOnUpdate: // For Timestamp and Datetime only.
			case ast.ColumnOptionFulltext:
			case ast.ColumnOptionComment:
				field.Comment = o.Expr.GetDatum().GetString()
			default:
				//return "", nil, errors.Errorf(" unsupport option %d\n", o.Tp)
			}
		}
		if !isPrimaryKey[colName] && isNotNull {
			gormTag.WriteString(";NOT NULL")
		}

		tags = append(tags, "gorm", gormTag.String())

		if opt.JsonTag {
			tags = append(tags, "json", colName)
		}

		field.Tag = makeTagStr(tags)

		// get type in golang
		nullStyle := opt.NullStyle
		if !canNull {
			nullStyle = NullDisable
		}

		goType, pkg := mysqlToGoType(col.Tp, nullStyle)
		if pkg != "" {
			importPath = append(importPath, pkg)
		}
		field.GoType = goType

		data.Fields = append(data.Fields, field)
	}
	return data, importPath, nil
}

func mysqlToGoType(colTp *types.FieldType, style NullStyle) (name string, path string) {
	if style == NullInSql {
		path = "database/sql"
		switch colTp.Tp {
		case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24, mysql.TypeLong:
			name = "sql.NullInt32"
		case mysql.TypeLonglong:
			name = "sql.NullInt64"
		case mysql.TypeFloat, mysql.TypeDouble:
			name = "sql.NullFloat64"
		case mysql.TypeString, mysql.TypeVarchar, mysql.TypeVarString,
			mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob:
			name = "sql.NullString"
		case mysql.TypeTimestamp, mysql.TypeDatetime, mysql.TypeDate:
			name = "sql.NullTime"
		case mysql.TypeDecimal, mysql.TypeNewDecimal:
			name = "sql.NullString"
		case mysql.TypeJSON:
			name = "sql.NullString"
		default:
			return "UnSupport", ""
		}
	} else {
		switch colTp.Tp {
		case mysql.TypeTiny:
			if mysql.HasUnsignedFlag(colTp.Flag) {
				name = "uint8"
			} else {
				name = "int8"
			}
		case mysql.TypeShort:
			if mysql.HasUnsignedFlag(colTp.Flag) {
				name = "uint16"
			} else {
				name = "int16"
			}
		case mysql.TypeInt24,
			mysql.TypeLong:
			if mysql.HasUnsignedFlag(colTp.Flag) {
				name = "uint32"
			} else {
				name = "int32"
			}
		case mysql.TypeLonglong:
			if mysql.HasUnsignedFlag(colTp.Flag) {
				name = "uint64"
			} else {
				name = "int64"
			}
		case mysql.TypeFloat, mysql.TypeDouble, mysql.TypeDecimal, mysql.TypeNewDecimal:
			name = "float64"
		case mysql.TypeString, mysql.TypeVarchar, mysql.TypeVarString,
			mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob:
			name = "string"
		case mysql.TypeTimestamp, mysql.TypeDatetime, mysql.TypeDate:
			path = "time"
			name = "time.Time"
		case mysql.TypeJSON:
			name = "string"
		default:
			return "UnSupport", ""
		}
		if style == NullInPointer {
			name = "*" + name
		}
	}
	return
}

func makeTagStr(tags []string) string {
	builder := strings.Builder{}
	for i := 0; i < len(tags)/2; i++ {
		builder.WriteString(tags[i*2])
		builder.WriteString(`:"`)
		builder.WriteString(tags[i*2+1])
		builder.WriteString(`" `)
	}
	if builder.Len() > 0 {
		return builder.String()[:builder.Len()-1]
	}
	return builder.String()
}

func getDefaultValue(expr ast.ExprNode) (value string) {
	if expr.GetDatum().Kind() != types.KindNull {
		value = fmt.Sprintf("%v", expr.GetDatum().GetValue())
	} else if expr.GetFlag() != ast.FlagConstant {
		if expr.GetFlag() == ast.FlagHasFunc {
			if funcExpr, ok := expr.(*ast.FuncCallExpr); ok {
				value = funcExpr.FnName.O
			}
		}
	}
	return
}
