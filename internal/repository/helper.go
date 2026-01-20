package repository

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LockStrength string

const (
	LockStrengthUpdate LockStrength = "UPDATE"
)

type Query struct {
	Page           int `json:"page"`
	Limit          int `json:"limit"`
	Query          string
	QueryOr        string
	ValueOr        []interface{}
	Values         []interface{}
	Sort           string
	Order          string
	SelectQuery    interface{}
	SelectArgs     []interface{}
	Joins          []string
	Preload        []string
	DistinctColumn string
	StartAt        time.Time `json:"start_at"`
	EndAt          time.Time `json:"end_at"`
	IsLive         *bool
	Select         string
	Group          string
	QuerySlice     []string
	DateFilterType string
	Table          string
	RemoveSort     bool
	LockStrength   LockStrength
}

func QueryHelperDB(q *gorm.DB, query Query) *gorm.DB {
	offset := (query.Page - 1) * query.Limit
	if query.Page > 0 {
		q = q.Offset(offset)
	}

	if query.Limit > 0 {
		q = q.Limit(query.Limit)
	}

	if query.Sort == "" {
		query.Sort = "created_at"
	}

	if query.Order == "" {
		query.Order = "DESC"
	}

	if query.Select != "" {
		q = q.Select(query.Select)
	}
	if query.Group != "" {
		q = q.Group(query.Group)
	}

	if query.SelectQuery != nil {
		q = q.Select(query.SelectQuery, query.SelectArgs...)
	}

	if !query.StartAt.IsZero() {
		if len(query.DateFilterType) > 0 {
			q = q.Where(fmt.Sprintf("%s >= ?", query.DateFilterType), query.StartAt.Format("2006-01-02 15:04:05 -07:00"))
		} else {
			// for cases where we join to a table with duplicate column name
			if query.Table != "" {
				startAtQuery := fmt.Sprintf("%s.created_at >= ?", query.Table)
				q = q.Where(startAtQuery, query.StartAt.Format("2006-01-02 15:04:05 -07:00"))

			} else {
				q = q.Where("created_at >= ?", query.StartAt.Format("2006-01-02 15:04:05 -07:00"))
			}
		}
	}

	if !query.EndAt.IsZero() {
		if len(query.DateFilterType) > 0 {
			q = q.Where(fmt.Sprintf("%s <= ?", query.DateFilterType), query.EndAt.Format("2006-01-02 15:04:05 -07:00"))
		} else {
			// for cases where we join to a table with duplicate column name
			if query.Table != "" {

				endAtQuery := fmt.Sprintf("%s.created_at <= ?", query.Table)
				q = q.Where(endAtQuery, query.EndAt.Format("2006-01-02 15:04:05 -07:00"))

			} else {
				q = q.Where("created_at <= ?", query.EndAt.Format("2006-01-02 15:04:05 -07:00"))
			}
		}
	}

	if query.IsLive != nil {
		q = q.Where("is_live = ?", query.IsLive)
	}

	for _, join := range query.Joins {
		q = q.Joins(join)
	}

	for _, preload := range query.Preload {
		q = q.Preload(preload)
	}

	if query.DistinctColumn != "" {
		q = q.Distinct(query.DistinctColumn)
	}

	if len(query.QuerySlice) > 0 {
		for i := range query.QuerySlice {
			if strings.Contains(query.QuerySlice[i], "BETWEEN") {
				v := query.Values[i].([]time.Time)
				q = q.Where(query.QuerySlice[i], v[0], v[1])
			} else if !strings.Contains(query.QuerySlice[i], "?") {
				q.Where(query.QuerySlice[i])
			} else {
				q = q.Where(query.QuerySlice[i], query.Values[i])
			}
		}
	} else if query.Query != "" {
		q = q.Where(query.Query, query.Values...)

		if query.QueryOr != "" {
			q = q.Where(
				q.Where(query.Query, query.Values...).Or(query.QueryOr, query.ValueOr...),
			)
		}
	}

	// added logic for Sum query where we do not want it to be sorted since it will cause an error
	if !query.RemoveSort {
		q = q.Order(fmt.Sprintf("%s %s", query.Sort, query.Order))
	}

	if query.Query == "" {
		return q
	}

	if query.LockStrength != LockStrength("") {
		q = q.Clauses(clause.Locking{Strength: string(query.LockStrength)})
	}

	return q
}

func WhereOnlyQueryHelperDB(q *gorm.DB, query Query) *gorm.DB {
	if len(query.QuerySlice) > 0 {
		for i := range query.QuerySlice {
			if strings.Contains(query.QuerySlice[i], "BETWEEN") {
				v := query.Values[i].([]time.Time)
				q = q.Where(query.QuerySlice[i], v[0], v[1])
			} else {
				q = q.Where(query.QuerySlice[i], query.Values[i])
			}
		}
	} else if query.Query != "" {
		if query.QueryOr != "" {
			q = q.Where(
				q.Where(query.Query, query.Values...).Or(query.QueryOr, query.ValueOr...),
			)
		} else {
			q = q.Where(query.Query, query.Values...)
		}
	}

	return q
}

func (q *Query) AddTermCondition(term string, value interface{}) {
	if len(q.Query) > 0 {
		term = " AND" + term
	}

	q.Query += term
	q.Values = append(q.Values, value)
}

func (q *Query) AddBetweenTermUsingDay(columnName string, start time.Time, days int) {
	if len(q.Query) > 0 {
		q.Query += " AND"
	}

	q.Query += " " + columnName + " BETWEEN ? AND ?"
	q.Values = append(q.Values, start, start.AddDate(0, 0, days))
}
