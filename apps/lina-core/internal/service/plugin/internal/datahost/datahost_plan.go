// This file applies typed plugindb query plans to governed data service
// requests.

package datahost

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/internal/service/plugin/internal/catalog"
	"lina-core/pkg/pluginbridge"
	"lina-core/pkg/plugindb"
)

// decodeDataListPlan restores a typed plugindb list plan or synthesizes one
// from the legacy list request fields.
func decodeDataListPlan(table string, request *pluginbridge.HostServiceDataListRequest) (*plugindb.DataQueryPlan, error) {
	var (
		requestPlan *plugindb.DataQueryPlan
		err         error
	)
	if request != nil && len(request.PlanJSON) > 0 {
		requestPlan, err = plugindb.UnmarshalQueryPlanJSON(request.PlanJSON)
		if err != nil {
			return nil, err
		}
	}
	if requestPlan == nil {
		request = normalizeDataListRequest(request)
		requestPlan = &plugindb.DataQueryPlan{
			Table:  strings.TrimSpace(table),
			Action: plugindb.DataPlanActionList,
			Page: &plugindb.DataPagination{
				PageNum:  request.PageNum,
				PageSize: request.PageSize,
			},
		}
		for field, value := range request.Filters {
			if strings.TrimSpace(value) == "" {
				continue
			}
			valueJSON, marshalErr := json.Marshal(value)
			if marshalErr != nil {
				return nil, marshalErr
			}
			filter := &plugindb.DataFilter{
				Field:     field,
				Operator:  plugindb.DataFilterOperatorEQ,
				ValueJSON: valueJSON,
			}
			requestPlan.Filters = append(requestPlan.Filters, filter)
		}
		return requestPlan, nil
	}
	if strings.TrimSpace(requestPlan.Table) == "" {
		requestPlan.Table = strings.TrimSpace(table)
	}
	if strings.TrimSpace(requestPlan.Table) != strings.TrimSpace(table) {
		return nil, gerror.Newf("plugindb query plan table mismatch: %s != %s", requestPlan.Table, table)
	}
	if requestPlan.Action == "" {
		requestPlan.Action = plugindb.DataPlanActionList
	}
	if requestPlan.Action != plugindb.DataPlanActionList && requestPlan.Action != plugindb.DataPlanActionCount {
		return nil, gerror.Newf("plugindb list request action is invalid: %s", requestPlan.Action)
	}
	if requestPlan.Action == plugindb.DataPlanActionList {
		if requestPlan.Page == nil {
			requestPlan.Page = &plugindb.DataPagination{PageNum: defaultDataListPageNum, PageSize: defaultDataListPageSize}
		}
		if requestPlan.Page.PageNum <= 0 {
			requestPlan.Page.PageNum = defaultDataListPageNum
		}
		if requestPlan.Page.PageSize <= 0 {
			requestPlan.Page.PageSize = defaultDataListPageSize
		}
		if requestPlan.Page.PageSize > maxDataListPageSize {
			requestPlan.Page.PageSize = maxDataListPageSize
		}
	}
	return requestPlan, plugindb.ValidateDataQueryPlan(requestPlan)
}

// decodeDataGetPlan restores a typed plugindb get plan or synthesizes one from
// the legacy get request key.
func decodeDataGetPlan(table string, request *pluginbridge.HostServiceDataGetRequest) (*plugindb.DataQueryPlan, error) {
	var (
		requestPlan *plugindb.DataQueryPlan
		err         error
	)
	if request != nil && len(request.PlanJSON) > 0 {
		requestPlan, err = plugindb.UnmarshalQueryPlanJSON(request.PlanJSON)
		if err != nil {
			return nil, err
		}
	}
	if requestPlan == nil {
		requestPlan = &plugindb.DataQueryPlan{Table: strings.TrimSpace(table), Action: plugindb.DataPlanActionGet}
	}
	if strings.TrimSpace(requestPlan.Table) == "" {
		requestPlan.Table = strings.TrimSpace(table)
	}
	if strings.TrimSpace(requestPlan.Table) != strings.TrimSpace(table) {
		return nil, gerror.Newf("plugindb get request table mismatch: %s != %s", requestPlan.Table, table)
	}
	if requestPlan.Action == "" {
		requestPlan.Action = plugindb.DataPlanActionGet
	}
	if requestPlan.Action != plugindb.DataPlanActionGet {
		return nil, gerror.Newf("plugindb get request action is invalid: %s", requestPlan.Action)
	}
	if request != nil && len(requestPlan.KeyJSON) == 0 {
		requestPlan.KeyJSON = append([]byte(nil), request.KeyJSON...)
	}
	if len(requestPlan.KeyJSON) == 0 {
		return nil, gerror.New("data key cannot be empty")
	}
	return requestPlan, plugindb.ValidateDataQueryPlan(requestPlan)
}

// applyPlanFilters applies typed plugindb filters against authorized resource fields.
func applyPlanFilters(model *gdb.Model, resource *catalog.ResourceSpec, filters []*plugindb.DataFilter) (*gdb.Model, error) {
	if model == nil || resource == nil || len(filters) == 0 {
		return model, nil
	}
	for _, filter := range filters {
		if err := plugindb.ValidateDataFilter(filter); err != nil {
			return nil, err
		}
		column := resolveResourceFieldColumn(resource, filter.Field)
		if column == "" {
			return nil, gerror.Newf("plugindb filter field is not authorized: %s", filter.Field)
		}
		switch filter.Operator {
		case plugindb.DataFilterOperatorEQ:
			value, err := plugindb.UnmarshalValueJSON(filter.ValueJSON)
			if err != nil {
				return nil, err
			}
			model = model.Where(column, value)
		case plugindb.DataFilterOperatorIN:
			values, err := plugindb.UnmarshalValuesJSON(filter.ValuesJSON)
			if err != nil {
				return nil, err
			}
			if len(values) == 0 {
				return nil, gerror.Newf("plugindb filter %s requires at least one value", filter.Operator)
			}
			model = model.WhereIn(column, values)
		case plugindb.DataFilterOperatorLike:
			value, err := plugindb.UnmarshalValueJSON(filter.ValueJSON)
			if err != nil {
				return nil, err
			}
			model = model.WhereLike(column, "%"+fmt.Sprint(value)+"%")
		default:
			return nil, gerror.Newf("plugindb filter operator is not supported: %s", filter.Operator)
		}
	}
	return model, nil
}

// buildPlanFieldArgs builds select expressions for the requested field subset.
func buildPlanFieldArgs(resource *catalog.ResourceSpec, selected []string) ([]any, error) {
	if len(selected) == 0 {
		return buildResourceFieldArgs(resource), nil
	}
	fields := make([]any, 0, len(selected))
	seen := make(map[string]struct{}, len(selected))
	for _, fieldName := range selected {
		normalizedField := strings.TrimSpace(fieldName)
		if normalizedField == "" {
			return nil, gerror.New("plugindb selected field cannot be empty")
		}
		if _, ok := seen[normalizedField]; ok {
			continue
		}
		seen[normalizedField] = struct{}{}
		column := resolveResourceFieldColumn(resource, normalizedField)
		if column == "" {
			return nil, gerror.Newf("plugindb selected field is not authorized: %s", normalizedField)
		}
		fields = append(fields, fmt.Sprintf("%s AS %s", column, quoteResourceAlias(normalizedField)))
	}
	return fields, nil
}

// buildPlanOrderBy builds the ORDER BY clause for the typed query plan.
func buildPlanOrderBy(resource *catalog.ResourceSpec, orders []*plugindb.DataOrder) (string, error) {
	if len(orders) == 0 {
		return buildResourceOrderBy(resource), nil
	}
	parts := make([]string, 0, len(orders))
	for _, order := range orders {
		if err := plugindb.ValidateDataOrder(order); err != nil {
			return "", err
		}
		column := resolveResourceFieldColumn(resource, order.Field)
		if column == "" {
			return "", gerror.Newf("plugindb order field is not authorized: %s", order.Field)
		}
		direction := "ASC"
		if order.Direction == plugindb.DataOrderDirectionDESC {
			direction = "DESC"
		}
		parts = append(parts, column+" "+direction)
	}
	return strings.Join(parts, ", "), nil
}

// buildResourceRecordWithSelection projects only the selected logical fields from one row.
func buildResourceRecordWithSelection(recordMap map[string]interface{}, resource *catalog.ResourceSpec, selected []string) map[string]interface{} {
	if len(selected) == 0 {
		return buildResourceRecord(recordMap, resource)
	}
	row := make(map[string]interface{}, len(selected))
	seen := make(map[string]struct{}, len(selected))
	for _, fieldName := range selected {
		normalizedField := strings.TrimSpace(fieldName)
		if normalizedField == "" {
			continue
		}
		if _, ok := seen[normalizedField]; ok {
			continue
		}
		seen[normalizedField] = struct{}{}
		field := findResourceField(resource, normalizedField)
		if field == nil {
			continue
		}
		row[normalizedField] = normalizeResourceValue(resolveResourceRecordValue(recordMap, field))
	}
	return row
}

// findResourceField returns the declared resource field for one logical name.
func findResourceField(resource *catalog.ResourceSpec, fieldName string) *catalog.ResourceField {
	if resource == nil {
		return nil
	}
	targetFieldName := strings.TrimSpace(fieldName)
	for _, field := range resource.Fields {
		if field != nil && field.Name == targetFieldName {
			return field
		}
	}
	return nil
}
