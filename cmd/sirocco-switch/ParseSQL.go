

// QueryInfo holds extracted routing metadata
type QueryInfo struct {
	Type  string
	Table string
}

// ParseSQL analyzes a SQL query and extracts routing info
func ParseSQL(query string) (*QueryInfo, error) {

	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("invalid SQL: %v", err)
	}

	info := &QueryInfo{
		Type:  "UNKNOWN",
		Table: "",
	}

	switch v := stmt.(type) {

	// SELECT
	case *sqlparser.Select:
		info.Type = "SELECT"
		info.Table = extractTableFromSelect(v.From)

	// INSERT
	case *sqlparser.Insert:
		info.Type = "INSERT"
		info.Table = v.Table.Name.String()

	// UPDATE
	case *sqlparser.Update:
		info.Type = "UPDATE"
		info.Table = extractTableGeneric(v.TableExprs)

	// DELETE
	case *sqlparser.Delete:
		info.Type = "DELETE"
		info.Table = extractTableGeneric(v.TableExprs)

	default:
		info.Type = "UNKNOWN"
	}

	return info, nil
}

// Extract table from SELECT
func extractTableFromSelect(from sqlparser.TableExprs) string {
	for _, expr := range from {
		if aliased, ok := expr.(*sqlparser.AliasedTableExpr); ok {
			if tbl, ok := aliased.Expr.(sqlparser.TableName); ok {
				return tbl.Name.String()
			}
		}
	}
	return ""
}

// Extract table from UPDATE/DELETE
func extractTableGeneric(exprs sqlparser.TableExprs) string {
	for _, expr := range exprs {
		if aliased, ok := expr.(*sqlparser.AliasedTableExpr); ok {
			if tbl, ok := aliased.Expr.(sqlparser.TableName); ok {
				return tbl.Name.String()
			}
		}
	}
	return ""
}
