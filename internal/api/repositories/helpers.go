package repositories

import (
	"fmt"
	"reflect"
	"strings"
)

func addSorting(query *strings.Builder, sortings []string) {
	addedSort := false
	for _, param := range sortings {
		parts := strings.Split(param, ":")
		if len(parts) != 2 {
			continue
		}
		dbField, order := parts[0], parts[1]
		if !isValidSortField(dbField) || !isValidSortOrder(order) {
			continue
		}
		if !addedSort {
			query.WriteString(" ORDER BY")
			addedSort = true
		} else {
			query.WriteString(",")
		}
		query.WriteString(" " + dbField + " " + order)
	}
}

func isValidSortOrder(order string) bool {
	orderLowerCase := strings.ToLower(order)
	return orderLowerCase == "asc" || orderLowerCase == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

func generateInsertQuery(model any) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders []string
	for field := range modelType.Fields() {
		dbTag := field.Tag.Get("db")
		dbTag = strings.Split(dbTag, ",")[0]
		if dbTag != "" && dbTag != "id" {
			columns = append(columns, dbTag)
			placeholders = append(placeholders, "?")
		}
	}
	return fmt.Sprintf("INSERT INTO teachers (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(placeholders, ", "))
}

func getStructValues(model any) []any {
	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()
	values := []any{}
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		dbTag = strings.Split(dbTag, ",")[0]
		if dbTag != "" && dbTag != "id" {
			values = append(values, modelValue.Field(i).Interface())
		}
	}
	return values
}
