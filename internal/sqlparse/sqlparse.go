package sqlparse

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Operation string

const (
	OpInsert Operation = "insert"
	OpSelect Operation = "select"
	OpUpdate Operation = "update"
	OpDelete Operation = "delete"
	OpCount  Operation = "count"
)

type Query struct {
	Operation    Operation
	Table        string
	KeyColumn    string
	KeyValue     string
	Columns      map[string]string
	Updates      map[string]string
	SelectAll    bool
	IsCount      bool
	Raw          string
}

var (
	insertRe = regexp.MustCompile(`(?i)^\s*INSERT\s+INTO\s+([a-zA-Z_][\w]*)\s*\(([^)]*)\)\s*VALUES\s*\(([^)]*)\)\s*;?\s*$`)
	selectRe = regexp.MustCompile(`(?i)^\s*SELECT\s+(.+?)\s+FROM\s+([a-zA-Z_][\w]*)(?:\s+WHERE\s+(.+?))?\s*;?\s*$`)
	updateRe = regexp.MustCompile(`(?i)^\s*UPDATE\s+([a-zA-Z_][\w]*)\s+SET\s+(.+?)\s+WHERE\s+(.+?)\s*;?\s*$`)
	deleteRe = regexp.MustCompile(`(?i)^\s*DELETE\s+FROM\s+([a-zA-Z_][\w]*)\s+WHERE\s+(.+?)\s*;?\s*$`)
	eqRe     = regexp.MustCompile(`(?i)^\s*([a-zA-Z_][\w]*)\s*=\s*(.+?)\s*$`)
)

func Parse(raw string) (*Query, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimSuffix(raw, ";")

	if m := insertRe.FindStringSubmatch(raw); m != nil {
    cols := splitCSV(m[2])
    vals := splitCSV(m[3])

    if len(cols) != len(vals) {
        return nil, fmt.Errorf("insert columns and values length mismatch")
    }

    q := &Query{
        Operation: OpInsert,
        Table:     strings.ToLower(m[1]),
        Columns:   map[string]string{},
        Raw:       raw,
    }

    for i := range cols {
        key := strings.ToLower(strings.TrimSpace(cols[i]))
        val := unquote(strings.TrimSpace(vals[i]))
        q.Columns[key] = val
    }

    // ❌ REMOVE shard key requirement entirely
    q.KeyColumn = ""
    q.KeyValue = ""

    return q, nil
}

	if m := selectRe.FindStringSubmatch(raw); m != nil {
		fields := strings.TrimSpace(m[1])
		table := strings.ToLower(m[2])
		where := strings.TrimSpace(m[3])

		q := &Query{
			Operation: OpSelect,
			Table:     table,
			Raw:       raw,
			Columns:   map[string]string{},
		}

		if strings.EqualFold(fields, "count(*)") {
			q.Operation = OpCount
			q.IsCount = true
			return q, nil
		}
		if fields == "*" {
			q.SelectAll = true
		}

		if where != "" {
			col, val, err := parseEquality(where)
			if err != nil {
				return nil, err
			}
			q.KeyColumn = strings.ToLower(col)
			q.KeyValue = val
			if q.KeyColumn == "user_id" {
				return q, nil
			}
			return q, nil
		}

		return q, nil
	}

	if m := updateRe.FindStringSubmatch(raw); m != nil {
		table := strings.ToLower(m[1])
		setClause := m[2]
		where := m[3]

		keyCol, keyVal, err := parseEquality(where)
		if err != nil {
			return nil, err
		}
		updates, err := parseAssignments(setClause)
		if err != nil {
			return nil, err
		}
		q := &Query{
			Operation: OpUpdate,
			Table:     table,
			KeyColumn: strings.ToLower(keyCol),
			KeyValue:  keyVal,
			Updates:   updates,
			Raw:       raw,
		}
		if q.KeyColumn != "user_id" {
			return nil, errors.New("update routing requires user_id in where clause")
		}
		return q, nil
	}

	if m := deleteRe.FindStringSubmatch(raw); m != nil {
		table := strings.ToLower(m[1])
		where := m[2]
		keyCol, keyVal, err := parseEquality(where)
		if err != nil {
			return nil, err
		}
		q := &Query{
			Operation: OpDelete,
			Table:     table,
			KeyColumn: strings.ToLower(keyCol),
			KeyValue:  keyVal,
			Raw:       raw,
		}
		if q.KeyColumn != "user_id" {
			return nil, errors.New("delete routing requires user_id in where clause")
		}
		return q, nil
	}

	return nil, fmt.Errorf("unsupported SQL: %s", raw)
}

func parseEquality(expr string) (string, string, error) {
	m := eqRe.FindStringSubmatch(strings.TrimSpace(expr))
	if m == nil {
		return "", "", fmt.Errorf("unsupported condition: %s", expr)
	}
	return strings.TrimSpace(m[1]), unquote(strings.TrimSpace(m[2])), nil
}

func parseAssignments(expr string) (map[string]string, error) {
	pairs := splitCSV(expr)
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		col, val, err := parseEquality(p)
		if err != nil {
			return nil, err
		}
		out[strings.ToLower(col)] = val
	}
	return out, nil
}

func splitCSV(s string) []string {
	var out []string
	var cur strings.Builder
	inSingle := false
	inDouble := false
	escaped := false

	for _, r := range s {
		switch {
		case escaped:
			cur.WriteRune(r)
			escaped = false
		case r == '\\':
			escaped = true
			cur.WriteRune(r)
		case r == '\'' && !inDouble:
			inSingle = !inSingle
			cur.WriteRune(r)
		case r == '"' && !inSingle:
			inDouble = !inDouble
			cur.WriteRune(r)
		case r == ',' && !inSingle && !inDouble:
			out = append(out, strings.TrimSpace(cur.String()))
			cur.Reset()
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		out = append(out, strings.TrimSpace(cur.String()))
	}
	return out
}

func unquote(v string) string {
	v = strings.TrimSpace(v)
	if len(v) >= 2 {
		if (v[0] == '\'' && v[len(v)-1] == '\'') || (v[0] == '"' && v[len(v)-1] == '"') {
			v = v[1 : len(v)-1]
		}
	}
	v = strings.ReplaceAll(v, `\"`, `"`)
	v = strings.ReplaceAll(v, `\'`, `'`)
	return v
}
