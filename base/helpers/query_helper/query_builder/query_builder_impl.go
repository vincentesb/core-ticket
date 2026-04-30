package query_builder

import (
	"core-ticket/base/helpers/query_helper"
	"core-ticket/base/helpers/query_helper/query_sanitizer"
	"core-ticket/base/helpers/struct_helper"
	"core-ticket/base/utility/nullable"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/guregu/null.v4"

	"github.com/jmoiron/sqlx"
)

var (
	EmptySqlxInstance    = errors.New("[Query Builder]: sqlx DB or Tx instance not found")
	TableNameCannotEmpty = errors.New("[Query Builder]: table name cannot be empty")
)

/*
sqlQueryBuilder represents a struct that holds various components for building SQL queries, such as
select clause, where condition, order condition, group by condition, limit condition, conditions map, insert configuration, update configuration, parameters, where parameters, existence flag, database connection, transaction, and table name.

Fields:
- selectClause: a slice of strings representing the select clause of the SQL query
- whereCondition: a string representing the where condition of the SQL query
- orderCondition: a string representing the order condition of the SQL query
- groupByCondition: a string representing the group by condition of the SQL query
- limitCondition: a string representing the limit condition of the SQL query
- condition: a map[string]query_helper.WhereCondition representing the conditions for the SQL query
- insertConfig: a query_helper.InsertConfig representing the insert configuration for the SQL query
- updateConfig: a query_helper.UpdateConfig representing the update configuration for the SQL query
- params: a slice of interface{} representing the parameters for the SQL query
- whereParams: a slice of interface{} representing the where parameters for the SQL query
- isExist: a boolean flag indicating the existence of a record in the SQL query
- db: a pointer to sqlx.DB representing the database connection
- tx: a pointer to sqlx.Tx representing the transaction
- tableName: a pointer to string representing the table name for the SQL query
*/
type sqlQueryBuilder struct {
	selectClause     []string
	whereCondition   string
	orderCondition   string
	groupByCondition string
	havingCondition  string
	limitCondition   string
	condition        map[string]query_helper.WhereCondition
	insertConfig     query_helper.InsertConfig
	updateConfig     query_helper.UpdateConfig
	params           []interface{}
	whereParams      []interface{}
	havingParams     []interface{}
	isExist          bool
	db               *sqlx.DB
	tx               *sqlx.Tx
	tableName        *string
	fromParams       []any
	joinClause       []string
	joinParams       []any
}

/*
Join adds a join clause to the SQL query being built by the QueryBuilder instance.
It specifies the type of join (e.g., LEFT JOIN, RIGHT JOIN), the table to join, the ON condition for the join, and any additional parameters for the join.

Parameters:
- joinType (JoinType): The type of join to perform (LeftJoin, RightJoin, InnerJoin, CrossJoin).
- tableName (string): The name of the table to join.
- on (string): The ON condition for the join.
- params (any): Additional parameters for the join condition.

Returns:
- QueryBuilder: The updated QueryBuilder instance with the join clause added.

Note:
- The joinType must be a valid JoinType value (LeftJoin, RightJoin, InnerJoin, CrossJoin) for the join clause to be added.
- The tableName and ON condition must be provided for the join to be valid and added to the query.
- Additional parameters can be passed to customize the join condition further.
*/
func (qb *sqlQueryBuilder) Join(
	joinType JoinType,
	tableName string,
	on string,
	params ...any,
) QueryBuilder {
	if !joinType.Valid() {
		return qb
	}

	if tableName == "" {
		return qb
	}

	if on == "" {
		return qb
	}

	pattern := `\w+\.\w+(\s\w+)?$`
	re := regexp.MustCompile(pattern)

	if re.MatchString(tableName) {
		tb := strings.Split(tableName, ".")
		if len(tb) == 2 {
			tb[0] = strings.Trim(tb[0], " ")
			tb[1] = strings.Trim(tb[1], " ")

			var alias string
			tbb := strings.Split(tb[1], " ")
			if len(tbb) == 2 {
				alias = fmt.Sprintf(" `%s`", tbb[1])
				tb[1] = strings.Trim(tbb[0], " ")
			}

			tableName = fmt.Sprintf("`%s`.`%s`%s", tb[0], tb[1], alias)
		}
	}

	qb.joinClause = append(
		qb.joinClause,
		fmt.Sprintf("%s %s ON %s", joinType.String(), tableName, on),
	)
	qb.joinParams = append(qb.joinParams, params...)

	return qb
}

/*
LeftJoin adds a LEFT JOIN clause to the SQL query being built.
It specifies the table to join, the ON condition for the join, and any additional parameters for the join.

Parameters:
- tableName (string): The name of the table to perform the LEFT JOIN on.
- on (string): The ON condition for the LEFT JOIN.
- params (any): Additional parameters for the LEFT JOIN.

Returns:
- QueryBuilder: The updated QueryBuilder instance with the LEFT JOIN clause added.

Example:
qb.LeftJoin("orders", "users.id = orders.user_id")

Note:
The LEFT JOIN clause is used to retrieve all records from the left table (table1), and the matched records from the right table (table2). The result is NULL from the right side if there is no match.
*/
func (qb *sqlQueryBuilder) LeftJoin(tableName string, on string, params ...any) QueryBuilder {
	return qb.Join(LeftJoin, tableName, on, params...)
}

/*
RightJoin adds a RIGHT JOIN clause to the SQL query being built.
It specifies the table to join, the ON condition for the join, and any additional parameters for the join.

Parameters:
- tableName (string): The name of the table to perform the RIGHT JOIN on.
- on (string): The ON condition for the RIGHT JOIN.
- params (any): Additional parameters for the RIGHT JOIN.

Returns:
- QueryBuilder: The updated QueryBuilder instance with the RIGHT JOIN added.

Example:
qb.RightJoin("orders", "users.id = orders.user_id")

Note:
The RightJoin method is used to join the current query with another table using a RIGHT JOIN operation.
*/
func (qb *sqlQueryBuilder) RightJoin(tableName string, on string, params ...any) QueryBuilder {
	return qb.Join(RightJoin, tableName, on, params...)
}

/*
InnerJoin adds an INNER JOIN clause to the SQL query being built.
It specifies the table to join, the ON condition for the join, and any additional parameters for the join.

Parameters:
- tableName (string): The name of the table to perform the inner join on.
- on (string): The ON condition for the inner join.
- params (any): Additional parameters for the inner join.

Returns:
- QueryBuilder: The updated QueryBuilder instance with the INNER JOIN clause added.

Example:
qb.InnerJoin("orders", "users.id = orders.user_id")

Note:
The INNER JOIN clause is used to combine rows from two or more tables based on a related column between them.
*/
func (qb *sqlQueryBuilder) InnerJoin(tableName string, on string, params ...any) QueryBuilder {
	return qb.Join(InnerJoin, tableName, on, params...)
}

/*
CrossJoin performs a CROSS JOIN operation in the SQL query being built by the QueryBuilder instance.
It joins the specified table with the current query using the CROSS JOIN type and the provided ON condition.
Additional parameters can be passed to customize the join operation further.

Parameters:
- tableName (string): The name of the table to perform the CROSS JOIN with.
- on (string): The ON condition specifying the join criteria.
- params (...any): Additional parameters that can be used in the join operation.

Returns:
- QueryBuilder: The updated QueryBuilder instance after applying the CROSS JOIN operation.

Example:
qb.CrossJoin("orders", "users.id = orders.user_id")

Note:
The CROSS JOIN operation combines each row from the first table with each row from the second table, resulting in a Cartesian product of the two tables.
*/
func (qb *sqlQueryBuilder) CrossJoin(tableName string, on string, params ...any) QueryBuilder {
	return qb.Join(CrossJoin, tableName, on, params...)
}

/*
AddWhereParams appends the provided parameters to the list of whereParams in the QueryBuilder instance.
It takes a variadic number of parameters and adds them to the existing whereParams slice.
The method then returns a reference to the updated QueryBuilder instance for method chaining.

Parameters:
- params (any): A variadic number of parameters to be added to the whereParams slice.

Returns:
- QueryBuilder: A reference to the QueryBuilder instance with the new parameters added.

Example:
qb.AddWhereParams(10, "John", true)

This will add the values 10, "John", and true to the whereParams slice in the QueryBuilder instance.

Note:
Make sure to aware that params ordering and location is must match with the query
*/
func (qb *sqlQueryBuilder) AddWhereParams(params ...any) QueryBuilder {
	qb.whereParams = append(qb.whereParams, params...)
	return qb
}

/*
Copy creates a new instance of the QueryBuilder with the same configuration as the current instance.
It copies the select clause, where condition, order condition, group by condition, limit condition,
condition map, insert configuration, update configuration, parameters, where parameters, existence flag,
database connection, transaction, and table name from the current QueryBuilder instance.

Returns:
- QueryBuilder: a new QueryBuilder instance with the same configuration as the current instance.

Note:
This method is useful for creating a copy of the QueryBuilder to make modifications without affecting the original instance.
*/
func (qb *sqlQueryBuilder) Copy() QueryBuilder {
	return &sqlQueryBuilder{
		selectClause:     qb.selectClause,
		whereCondition:   qb.whereCondition,
		orderCondition:   qb.orderCondition,
		groupByCondition: qb.groupByCondition,
		havingCondition:  qb.havingCondition,
		limitCondition:   qb.limitCondition,
		condition:        qb.condition,
		insertConfig:     qb.insertConfig,
		updateConfig:     qb.updateConfig,
		params:           qb.params,
		whereParams:      qb.whereParams,
		havingParams:     qb.havingParams,
		isExist:          qb.isExist,
		db:               qb.db,
		tx:               qb.tx,
		tableName:        qb.tableName,
		fromParams:       qb.fromParams,
	}
}

/*
CopyAll The difference with the usual Copy() method is that this method copy everything, including join clauses and join parameters.
*/
func (qb *sqlQueryBuilder) CopyAll() QueryBuilder {
	return &sqlQueryBuilder{
		selectClause:     qb.selectClause,
		whereCondition:   qb.whereCondition,
		orderCondition:   qb.orderCondition,
		groupByCondition: qb.groupByCondition,
		havingCondition:  qb.havingCondition,
		limitCondition:   qb.limitCondition,
		condition:        qb.condition,
		insertConfig:     qb.insertConfig,
		updateConfig:     qb.updateConfig,
		params:           qb.params,
		whereParams:      qb.whereParams,
		havingParams:     qb.havingParams,
		isExist:          qb.isExist,
		db:               qb.db,
		tx:               qb.tx,
		tableName:        qb.tableName,
		fromParams:       qb.fromParams,
		joinClause:       qb.joinClause,
		joinParams:       qb.joinParams,
	}
}

/*
From will set the table name for the SQL query being built.

Parameters:
- tableName (string): The name of the table to set for the query.
- params (any): Additional parameters for the LEFT JOIN.

Returns:
- QueryBuilder: The updated QueryBuilder instance with the table name set.

Note:
The 'From' method must be called before executing methods like Exist, Count, One, All, Create, Update, or Delete to specify the table to operate on.
*/
func (qb *sqlQueryBuilder) From(tableName string, params ...any) QueryBuilder {
	qb.tableName = &tableName
	qb.fromParams = params
	return qb
}

/*
ConfigUpdate sets the update configuration for the SQL query builder.
It assigns the provided update configuration to the sqlQueryBuilder instance for building update queries.

Parameters:
- updateConfig (query_helper.UpdateConfig): The update configuration containing the model and saved attributes.

Returns:
- QueryBuilder: The updated QueryBuilder instance with the new update configuration set.

Example:
qb.ConfigUpdate(updateConfig)

Note:
This method must be called before executing the Update method to ensure the update query is properly configured.
*/
func (qb *sqlQueryBuilder) ConfigUpdate(updateConfig query_helper.UpdateConfig) QueryBuilder {
	qb.updateConfig = updateConfig
	return qb
}

/*
ConfigInsert sets the configuration for inserting data into the database using the provided InsertConfig struct.

Parameters:
- insertConfig (query_helper.InsertConfig): The configuration settings for the insert operation, including the model and saved attributes.

Returns:
- QueryBuilder: The updated QueryBuilder instance with the insert configuration set.

Example:
qb.ConfigInsert(query_helper.InsertConfig{Model: User{}, SavedAttribute: []string{"name", "email"}})

Note:
This method must be called before executing the insert operation using the Create method.
*/
func (qb *sqlQueryBuilder) ConfigInsert(insertConfig query_helper.InsertConfig) QueryBuilder {
	qb.insertConfig = insertConfig
	return qb
}

/*
ConfigWhereCondition sets the conditions to be used in the WHERE clause of the SQL query builder.

Deprecated: This method is deprecated and will be removed soon.

Parameters:
- condition: a map where the key is the column name or expression, and the value is a WhereCondition struct containing the comparison operator and value for that column.

Returns:
- QueryBuilder: a reference to the QueryBuilder instance for method chaining.

Example:

	condition := map[string]query_helper.WhereCondition{
		"age": {Operator: ">", Value: 18},
		"name": {Operator: "LIKE", Value: "John%"},
	}

qb.ConfigWhereCondition(condition)

Note:
This method allows setting multiple conditions to be used in the WHERE clause of the SQL query. The conditions are applied using the specified comparison operators and values for each column.
*/
func (qb *sqlQueryBuilder) ConfigWhereCondition(
	condition map[string]query_helper.WhereCondition,
) QueryBuilder {
	qb.condition = condition
	return qb
}

/*
AndWhereUsingCondition checks if the specified column exists in the condition map.
If the column exists, it applies the AndWhere method with the operator and value from the condition map.

Deprecated: This method is deprecated and will be removed soon.

Parameters:
- column: the name of the column to check in the condition map

Returns:
- QueryBuilder: the updated query builder instance
*/
func (qb *sqlQueryBuilder) AndWhereUsingCondition(column string) QueryBuilder {
	if _, exists := qb.condition[column]; !exists {
		return qb
	}
	return qb.AndWhere(qb.condition[column].Operator, column, qb.condition[column].Value)
}

/*
OrWhereUsingCondition checks if the specified column exists in the condition map.
If the column exists, it applies the OrWhere method with the operator and value from the condition map.

Deprecated: This method is deprecated and will be removed soon.

Parameters:
- column: the name of the column to check in the condition map

Returns:
- QueryBuilder: the updated query builder instance
*/
func (qb *sqlQueryBuilder) OrWhereUsingCondition(column string) QueryBuilder {
	if _, exists := qb.condition[column]; !exists {
		return qb
	}
	return qb.OrWhere(qb.condition[column].Operator, column, qb.condition[column].Value)
}

/*
AndFilterWhereUsingCondition checks if the specified column exists in the condition map.
If the column exists, it applies the AndFilterWhere method with the operator and value from the condition map.

Deprecated: This method is deprecated and will be removed soon.

Parameters:
- column: the name of the column to check in the condition map

Returns:
- QueryBuilder: the updated query builder instance
*/
func (qb *sqlQueryBuilder) AndFilterWhereUsingCondition(column string) QueryBuilder {
	if _, exists := qb.condition[column]; !exists {
		return qb
	}
	return qb.AndFilterWhere(qb.condition[column].Operator, column, qb.condition[column].Value)
}

/*
OrFilterWhereUsingCondition checks if the specified column exists in the condition map.
If the column exists, it applies the OrFilterWhere method with the operator and value from the condition map.

Deprecated: This method is deprecated and will be removed soon.

Parameters:
- column: the name of the column to check in the condition map

Returns:
- QueryBuilder: the updated query builder instance
*/
func (qb *sqlQueryBuilder) OrFilterWhereUsingCondition(column string) QueryBuilder {
	if _, exists := qb.condition[column]; !exists {
		return qb
	}
	return qb.OrFilterWhere(qb.condition[column].Operator, column, qb.condition[column].Value)
}

/*
AndFilterWhere adds an AND filter condition to the query builder based on the provided operator, column, and value.
If the value is empty (nil, zero, or null), the filter condition will not be added to the query.
This method then delegates to the AndWhere method to add the filter condition.

Parameters:
- operator: The operator to use in the filter condition (e.g., "=", ">", "<").
- column: The column on which the filter condition is applied.
- value: The value to compare against in the filter condition.

Returns:
- QueryBuilder: The updated query builder instance with the filter condition added.

Note: This method is used to add an AND filter condition to the query builder.
*/
func (qb *sqlQueryBuilder) AndFilterWhere(
	operator string,
	column interface{},
	value interface{},
) QueryBuilder {
	if qb.checkForEmptyValue(value) {
		return qb
	}
	return qb.AndWhere(operator, column, value)
}

/*
OrFilterWhere adds an OR filter condition to the query builder based on the provided operator, column, and value.
If the value is empty (nil, zero, or null), the filter condition will not be added to the query.
This method then delegates to the OrWhere method to add the filter condition.

Parameters:
- operator: The comparison operator for the filter condition.
- column: The column name or expression to filter on.
- value: The value to compare against in the filter condition.

Returns:
- QueryBuilder: The query builder instance with the OR filter condition added, if the value is not empty.

Example:
qb.OrFilterWhere("=", "age", 30)

Note:
This method is used to add OR filter conditions to the query builder.
*/
func (qb *sqlQueryBuilder) OrFilterWhere(
	operator string,
	column interface{},
	value interface{},
) QueryBuilder {
	if qb.checkForEmptyValue(value) {
		return qb
	}
	return qb.OrWhere(operator, column, value)
}

func (qb *sqlQueryBuilder) checkForEmptyValue(value interface{}) bool {
	rfVal := reflect.ValueOf(value)
	if !rfVal.IsValid() {
		return true
	}

	if strings.Contains(rfVal.Type().PkgPath(), "null") {
		return !rfVal.FieldByName("Valid").Bool()
	}
	switch rfVal.Kind() {
	case reflect.Ptr:
		return rfVal.IsNil()
	default:
		return rfVal.IsZero()
	}
}

/*
Clear resets all the fields of the sqlQueryBuilder struct to their zero values or empty states.
This method is useful for reusing the sqlQueryBuilder instance for a new query without any residual data from previous queries.
*/
func (qb *sqlQueryBuilder) Clear() {
	qb.selectClause = nil
	qb.whereCondition = ""
	qb.orderCondition = ""
	qb.havingCondition = ""
	qb.limitCondition = ""
	qb.condition = nil
	qb.insertConfig = query_helper.InsertConfig{}
	qb.updateConfig = query_helper.UpdateConfig{}
	qb.whereParams = nil
	qb.params = nil
	qb.isExist = false
	qb.groupByCondition = ""
}

/*
ClearAll The difference with the usual Clear() method is that this method removes everything, including join clauses and join parameters.
*/
func (qb *sqlQueryBuilder) ClearAll() {
	qb.selectClause = nil
	qb.whereCondition = ""
	qb.orderCondition = ""
	qb.havingCondition = ""
	qb.limitCondition = ""
	qb.condition = nil
	qb.insertConfig = query_helper.InsertConfig{}
	qb.updateConfig = query_helper.UpdateConfig{}
	qb.whereParams = nil
	qb.params = nil
	qb.isExist = false
	qb.groupByCondition = ""
	qb.joinClause = []string{}
	qb.joinParams = []any{}
	qb.tableName = nil
	qb.fromParams = nil
}

/*
Limit sets the limit and offset for the SQL query.
It takes two parameters: take (number of rows to retrieve) and skip (number of rows to skip).
It constructs the LIMIT and OFFSET clauses for the query based on the provided parameters.
It returns the QueryBuilder instance to allow for method chaining.

Parameters:
- take: number of rows to retrieve
- skip: number of rows to skip

Returns:
- QueryBuilder: instance of the QueryBuilder interface

Example:
qb.Limit(10, 5)

This will set the limit to retrieve 10 rows starting from the 6th row.
*/
func (qb *sqlQueryBuilder) Limit(take int, skip int) QueryBuilder {
	qb.limitCondition = fmt.Sprintf("LIMIT %d OFFSET %d", take, skip)
	return qb
}

/*
setWhereParams sets the parameters for the WHERE clause in the SQL query based on the provided column, operator, and value.

Parameters:
- column: The column to apply the condition on. It can be a string or an Expression.
- operator: The comparison operator to use in the condition (e.g., "=", "LIKE", "IN", "BETWEEN", "IS").
- value: The value to compare the column with.

Returns:
- QueryBuilder: The updated QueryBuilder instance with the WHERE parameters set.

Note:
- This method dynamically handles different types of conditions LIKE, IN, BETWEEN, IS NULL, etc.
- It constructs the WHERE condition and appends the corresponding parameters to the whereParams slice in the QueryBuilder instance.
*/
func (qb *sqlQueryBuilder) setWhereParams(
	column interface{},
	operator string,
	value interface{},
) QueryBuilder {
	var col string
	switch v := column.(type) {
	case string:
		if split := strings.Split(v, "."); len(split) == 2 {
			col = fmt.Sprintf("`%s`.`%s`", split[0], split[1])
		} else {
			col = fmt.Sprintf("`%s`", v)
		}
	case Expression:
		col = v.String()
	default:
		return qb
	}

	operator = strings.ToUpper(operator)
	switch {
	case strings.Contains(operator, "LIKE"):
		rfVal := reflect.ValueOf(value)
		if strings.Contains(rfVal.Type().PkgPath(), "null") {
			if slices.Contains([]reflect.Type{
				reflect.TypeOf(null.String{}),
				reflect.TypeOf(nullable.String{}),
			}, rfVal.Type()) {
				field := rfVal.FieldByName("String")
				value = field.String()
			} else {
				break
			}
		}
		if rfVal.Type().Kind() == reflect.Ptr {
			value = rfVal.Elem().Interface()
		}
		qb.whereParams = append(qb.whereParams, fmt.Sprintf("%%%s%%", value))
		qb.whereCondition += fmt.Sprintf("%s %s ?", col, operator)
	case strings.Contains(operator, "IN"):
		rfVal := reflect.ValueOf(value)
		if rfVal.Kind() != reflect.Slice {
			return qb
		}
		if rfVal.IsNil() {
			qb.whereCondition += col + " IS NULL"
			return qb
		}
		str := " ("
		for i := 0; i < rfVal.Len(); i++ {
			str += "?"
			if i != rfVal.Len()-1 {
				str += ","
			}
			qb.whereParams = append(qb.whereParams, rfVal.Index(i).Interface())
		}
		str += ")"
		qb.whereCondition += fmt.Sprintf("(%s %s%s)", col, operator, str)
	case strings.Contains(operator, "BETWEEN"):
		rfVal := reflect.ValueOf(value)
		if rfVal.Kind() != reflect.Slice && rfVal.Kind() != reflect.Array {
			return qb
		}

		if rfVal.Len() != 2 {
			return qb
		}

		qb.whereParams = append(
			qb.whereParams,
			rfVal.Index(0).Interface(),
			rfVal.Index(1).Interface(),
		)
		qb.whereCondition += fmt.Sprintf("(%s %s ? AND ?)", col, operator)
	case strings.Contains(operator, "IS"):
		if value != nil {
			return qb
		}
		qb.whereCondition += fmt.Sprintf("(%s %s NULL)", col, operator)

	default:
		qb.whereParams = append(qb.whereParams, value)
		qb.whereCondition += fmt.Sprintf("(%s %s ?)", col, operator)
	}
	return qb
}

func (qb *sqlQueryBuilder) whereRaw(rawSql string, params ...any) QueryBuilder {
	qb.whereCondition += rawSql
	qb.whereParams = append(qb.whereParams, params...)
	return qb
}

func (qb *sqlQueryBuilder) havingRaw(rawSql string, params ...any) QueryBuilder {
	qb.havingCondition += rawSql
	qb.havingParams = append(qb.havingParams, params...)
	return qb
}

/*
AndWhereRaw appends an 'AND' operator to the existing WHERE condition and adds a raw SQL string with parameters to the query builder.

Warning: This method does not check for SQL injection vulnerabilities.

Parameters:
- rawSql: a string representing the raw SQL condition to be added.
- params: variadic parameters to be used in the raw SQL condition.

Returns:
- QueryBuilder: the query builder instance for method chaining.
*/
func (qb *sqlQueryBuilder) AndWhereRaw(rawSql string, params ...any) QueryBuilder {
	if qb.whereCondition != "" {
		qb.whereCondition += " AND "
	}
	return qb.whereRaw(rawSql, params...)
}

/*
OrWhereRaw appends the given raw SQL string with OR operator to the existing WHERE condition in the query builder.

Warning: This function does not check for SQL injection vulnerabilities.

Parameters:
- rawSql: the raw SQL string to be appended with OR operator.
- params: optional parameters to be used in the raw SQL string.

Returns:
- QueryBuilder: the updated query builder instance.

Example:
qb.OrWhereRaw("age > ?", 18)
*/
func (qb *sqlQueryBuilder) OrWhereRaw(rawSql string, params ...any) QueryBuilder {
	if qb.whereCondition != "" {
		qb.whereCondition += " OR "
	}
	return qb.whereRaw(rawSql, params...)
}

/*
AndWhere adds an 'AND' condition to the WHERE clause of the SQL query being built.

Parameters:
- operator: The comparison operator for the condition (e.g., '=', '>', '<', 'LIKE', 'IN', 'BETWEEN', 'IS').
- column: The column name or expression to apply the condition on.
- value: The value to compare the column with.

Returns:
- QueryBuilder: The updated QueryBuilder object with the new 'AND' condition added to the WHERE clause.

Note:
This method appends the 'AND' keyword to the existing WHERE clause if it is not empty, and then sets the parameters for the specified condition based on the provided operator, column, and value.
*/
func (qb *sqlQueryBuilder) AndWhere(
	operator string,
	column interface{},
	value interface{},
) QueryBuilder {
	if qb.whereCondition != "" {
		qb.whereCondition += " AND "
	}
	return qb.setWhereParams(column, operator, value)
}

/*
OrWhere adds an OR condition to the WHERE clause of the SQL query being built.

Parameters:
- operator: The comparison operator for the condition (e.g., "=", "LIKE", "IN").
- column: The column name or expression to apply the condition on.
- value: The value to compare against.

Returns:
- QueryBuilder: The updated QueryBuilder instance with the OR condition added.

Note:
- If the WHERE clause already has conditions, " OR " will be appended before adding the new condition.
*/
func (qb *sqlQueryBuilder) OrWhere(
	operator string,
	column interface{},
	value interface{},
) QueryBuilder {
	if qb.whereCondition != "" {
		qb.whereCondition += " OR "
	}
	return qb.setWhereParams(column, operator, value)
}

/*
OrderBy sets the order condition for the SQL query.

Parameters:
- order (string): The order condition to be set.

Returns:
- QueryBuilder: The QueryBuilder instance with the order condition set.

Example:
qb.OrderBy("created_at DESC")

Note:
This method is used to specify the order in which the query results should be returned.
*/
func (qb *sqlQueryBuilder) OrderBy(order string) QueryBuilder {
	qb.orderCondition = query_sanitizer.SanitizeOrderBy(order)
	return qb
}

/*
OrderBy sets the order condition for the SQL query.

Warning: This method does not check for SQL injection vulnerabilities. You can do it manually by using query_sanitizer.SanitizeOrderBy

Parameters:
- order (string): The order condition to be set.

Returns:
- QueryBuilder: The QueryBuilder instance with the order condition set.

Example:
qb.OrderBy("created_at DESC")

Note:
This method is used to specify the order in which the query results should be returned.
*/
func (qb *sqlQueryBuilder) OrderByRaw(order string) QueryBuilder {
	qb.orderCondition = order
	return qb
}

/*
GroupBy sets the GROUP BY condition for the SQL query.

Parameters:
- groupBy (string): The column or columns to group the results by.

Returns:
QueryBuilder: The QueryBuilder instance to allow for method chaining.
*/
func (qb *sqlQueryBuilder) GroupBy(groupBy string) QueryBuilder {
	qb.groupByCondition = groupBy
	return qb
}

/*
setHavingParams sets the parameters for the HAVING clause in the SQL query based on the provided column, operator, and value.

Parameters:
- column: The column to apply the condition on. It can be a string or an Expression.
- operator: The comparison operator to use in the condition (e.g., "=", ">", "<").
- value: The value to compare the column with.

Returns:
- QueryBuilder: The updated QueryBuilder instance with the HAVING parameters set.

Note:
- It constructs the HAVING condition and appends the corresponding parameters to the havingParams slice in the QueryBuilder instance.
*/
func (qb *sqlQueryBuilder) setHavingParams(
	column interface{},
	operator string,
	value interface{},
) QueryBuilder {
	var col string
	switch v := column.(type) {
	case string:
		if split := strings.Split(v, "."); len(split) == 2 {
			col = fmt.Sprintf("`%s`.`%s`", split[0], split[1])
		} else {
			col = fmt.Sprintf("`%s`", v)
		}
	case Expression:
		col = v.String()
	default:
		return qb
	}

	operator = strings.ToUpper(operator)
	switch {
	case strings.Contains(operator, "LIKE"):
		rfVal := reflect.ValueOf(value)
		if strings.Contains(rfVal.Type().PkgPath(), "null") {
			if slices.Contains([]reflect.Type{
				reflect.TypeOf(null.String{}),
				reflect.TypeOf(nullable.String{}),
			}, rfVal.Type()) {
				field := rfVal.FieldByName("String")
				value = field.String()
			} else {
				break
			}
		}
		if rfVal.Type().Kind() == reflect.Ptr {
			value = rfVal.Elem().Interface()
		}
		qb.havingParams = append(qb.havingParams, fmt.Sprintf("%%%s%%", value))
		qb.havingCondition += fmt.Sprintf("(%s %s ?)", col, operator)
	case strings.Contains(operator, "IN"):
		rfVal := reflect.ValueOf(value)
		if rfVal.Kind() != reflect.Slice {
			return qb
		}
		if rfVal.IsNil() {
			qb.havingCondition += col + " IS NULL"
			return qb
		}
		str := " ("
		for i := 0; i < rfVal.Len(); i++ {
			str += "?"
			if i != rfVal.Len()-1 {
				str += ","
			}
			qb.havingParams = append(qb.havingParams, rfVal.Index(i).Interface())
		}
		str += ")"
		qb.havingCondition += fmt.Sprintf("(%s %s%s)", col, operator, str)
	case strings.Contains(operator, "BETWEEN"):
		rfVal := reflect.ValueOf(value)
		if rfVal.Kind() != reflect.Slice && rfVal.Kind() != reflect.Array {
			return qb
		}

		if rfVal.Len() != 2 {
			return qb
		}

		qb.havingParams = append(
			qb.havingParams,
			rfVal.Index(0).Interface(),
			rfVal.Index(1).Interface(),
		)
		qb.havingCondition += fmt.Sprintf("(%s %s ? AND ?)", col, operator)
	case strings.Contains(operator, "IS"):
		if value != nil {
			return qb
		}
		qb.havingCondition += fmt.Sprintf("(%s %s NULL)", col, operator)

	default:
		qb.havingParams = append(qb.havingParams, value)
		qb.havingCondition += fmt.Sprintf("(%s %s ?)", col, operator)
	}
	return qb
}

/*
AndHaving adds an 'AND' condition to the HAVING clause of the SQL query being built.

Parameters:
- operator: The comparison operator for the condition (e.g., '=', '>', '<').
- column: The column name or expression to apply the condition on.
- value: The value to compare the column with.

Returns:
- QueryBuilder: The updated QueryBuilder object with the new 'AND' condition added to the HAVING clause.

Note:
This method appends the 'AND' keyword to the existing HAVING clause if it is not empty, and then sets the parameters for the specified condition based on the provided operator, column, and value.
*/
func (qb *sqlQueryBuilder) AndHaving(
	operator string,
	column interface{},
	value interface{},
) QueryBuilder {
	if qb.havingCondition != "" {
		qb.havingCondition += " AND "
	}
	return qb.setHavingParams(column, operator, value)
}

/*
AndHavingRaw adds an 'AND' condition to the HAVING clause and adds a raw SQL string with parameters to the query builder.

Warning: This method does not check for SQL injection vulnerabilities.

Parameters:
- rawSql: a string representing the raw SQL condition to be added.
- params: variadic parameters to be used in the raw SQL condition.

Returns:
- QueryBuilder: The updated QueryBuilder object with the new 'AND' condition added to the HAVING clause.

Note:
This method appends the 'AND' keyword to the existing HAVING clause if it is not empty, and then sets the parameters for the specified condition based on the provided operator, column, and value.
*/
func (qb *sqlQueryBuilder) AndHavingRaw(
	rawSql string,
	params ...any,
) QueryBuilder {
	if qb.havingCondition != "" {
		qb.havingCondition += " AND "
	}
	return qb.havingRaw(rawSql, params...)
}

/*
OrHaving adds an OR condition to the HAVING clause of the SQL query being built.

Parameters:
- operator: The comparison operator for the condition (e.g., "=", "LIKE", "IN").
- column: The column name or expression to apply the condition on.
- value: The value to compare against.

Returns:
- QueryBuilder: The updated QueryBuilder instance with the OR having condition added.

Note:
- If the HAVING clause already has conditions, " OR " will be appended before adding the new condition.
*/
func (qb *sqlQueryBuilder) OrHaving(
	operator string,
	column interface{},
	value interface{},
) QueryBuilder {
	if qb.havingCondition != "" {
		qb.havingCondition += " OR "
	}
	return qb.setHavingParams(column, operator, value)
}

/*
AndFilterHaving adds an AND filter having condition to the query builder based on the provided operator, column, and value.
If the value is empty (nil, zero, or null), the filter having condition will not be added to the query.
This method then delegates to the AndHaving method to add the filter condition.

Parameters:
- operator: The operator to use in the filter condition (e.g., "=", ">", "<").
- column: The column on which the filter condition is applied.
- value: The value to compare against in the filter condition.

Returns:
- QueryBuilder: The updated query builder instance with the filter having condition added.

Note: This method is used to add an AND filter having condition to the query builder.
*/
func (qb *sqlQueryBuilder) AndFilterHaving(
	operator string,
	column interface{},
	value interface{},
) QueryBuilder {
	if qb.checkForEmptyValue(value) {
		return qb
	}
	return qb.AndHaving(operator, column, value)
}

/*
OrFilterHaving adds an 'OR' condition to the HAVING clause of the SQL query being built based on the provided operator, column, and value.
If the value is empty (nil, zero, or null), the filter condition will not be added to the query.
This method then delegates to the OrHaving method to add the filter condition.

Parameters:
- operator: The comparison operator for the condition (e.g., '=', '>', '<').
- column: The column name or expression to apply the condition on.
- value: The value to compare the column with.

Returns:
  - QueryBuilder: The updated QueryBuilder object with the new 'OR' condition added to the HAVING clause, if the value is not empty.

Note:
This method appends the 'OR' keyword to the existing HAVING clause if it is not empty, and then sets the parameters for the specified condition based on the provided operator, column, and value.
*/
func (qb *sqlQueryBuilder) OrFilterHaving(
	operator string,
	column interface{},
	value interface{},
) QueryBuilder {
	if qb.checkForEmptyValue(value) {
		return qb
	}
	return qb.OrHaving(operator, column, value)
}

/*
Select sets the columns to be selected in the SQL query.

Parameters:
- columns: a variadic parameter of strings representing the columns to be selected.

Returns:
- QueryBuilder: a reference to the QueryBuilder instance for method chaining.

Example:
qb.Select("id", "name", "email")

Note:
This method must be called before building the SQL query using BuildSQL or Build method.
*/
func (qb *sqlQueryBuilder) Select(columns ...string) QueryBuilder {
	qb.selectClause = columns
	return qb
}

/*
buildSelect constructs the SELECT query based on the selectClause field of the sqlQueryBuilder struct.
If selectClause is not empty, it joins the elements with commas. Otherwise, it selects all columns with '*'.
Returns the constructed SELECT query as a string.
*/
func (qb *sqlQueryBuilder) buildSelect() string {
	query := "SELECT "
	if len(qb.selectClause) > 0 {
		query += strings.Join(qb.selectClause, ",")
	} else {
		query += "*"
	}
	return query
}

/*
buildWhere builds the WHERE clause of the SQL query based on the stored whereCondition in the sqlQueryBuilder struct.
If whereCondition is not empty, it appends "WHERE" followed by the whereCondition to the query string and adds whereParams to the params slice.
Returns the constructed WHERE clause as a string.
*/
func (qb *sqlQueryBuilder) buildWhere() string {
	var query string
	if qb.whereCondition != "" {
		query += fmt.Sprintf("WHERE %s", qb.whereCondition)
		qb.params = append(qb.params, qb.whereParams...)
	}
	return query
}

/*
buildLimit returns the LIMIT clause of the SQL query being built by the sqlQueryBuilder instance.

Returns:
- string: The LIMIT clause of the SQL query, or an empty string if no limit is set.
*/
func (qb *sqlQueryBuilder) buildLimit() string {
	var query string
	if qb.limitCondition != "" {
		query += fmt.Sprintf("%s", qb.limitCondition)
	}
	return query
}

/*
buildCreate generates the SQL query for inserting a new record into the database based on the provided InsertConfig.
It iterates over the SavedAttribute list in the InsertConfig to build the column names and placeholders for values.
If an attribute is not found in the model, it will panic with an error message.
The method returns the constructed SQL query string.

Returns:

	string: The SQL query for inserting a new record.

Panic:

	If an attribute specified in SavedAttribute is not found in the model.
*/
func (qb *sqlQueryBuilder) buildCreate() string {
	var sqlVarPlaceholder string

	query := "("
	sqlVarPlaceholder += "("
	for j, attr := range qb.insertConfig.SavedAttribute {
		//if reflect.ValueOf(qb.insertConfig.Model).FieldByName(attr).IsValid() {
		//	continue
		//}
		if j != 0 {
			query += ","
			sqlVarPlaceholder += ","
		}
		field, exist := reflect.TypeOf(qb.insertConfig.Model).FieldByName(attr)
		if !exist {
			panic(fmt.Sprintf(fmt.Sprintf("attribute %s is not exist on the model", attr)))
		}
		query += "`" + field.Tag.Get("db") + "`"
		sqlVarPlaceholder += "?"
		if j == len(qb.insertConfig.SavedAttribute)-1 {
			query += ") VALUE "
			sqlVarPlaceholder += ")"
		}
		qb.params = append(
			qb.params,
			reflect.ValueOf(qb.insertConfig.Model).FieldByName(attr).Interface(),
		)
	}
	query += sqlVarPlaceholder
	return query
}

/*
buildBatchCreate generates a SQL batch insert query based on the insertConfig provided in the sqlQueryBuilder instance.
It constructs the query by iterating over the saved attributes of the model and building the VALUES clause for each entry in the model slice.

Returns:
- string: The constructed SQL batch insert query.

Panics:
- If the model is empty or not a slice or array of structs.
- If an attribute specified in the insertConfig is not found in the model struct.
- If a db tag is missing for an attribute in the model struct.
*/
func (qb *sqlQueryBuilder) buildBatchCreate() string {
	rtMod := reflect.ValueOf(qb.insertConfig.Model)
	if !slices.Contains([]reflect.Kind{
		reflect.Slice,
		reflect.Array,
	}, rtMod.Kind()) ||
		rtMod.Len() == 0 {
		panic(fmt.Errorf("model is empty or not slice or array of struct"))
	}

	if rtMod.Index(0).Kind() != reflect.Struct {
		panic(fmt.Errorf("model is not slice or array of struct"))
	}

	var querySB strings.Builder
	querySB.WriteRune('(')
	for i, attr := range qb.insertConfig.SavedAttribute {
		f, e := rtMod.Index(0).Type().FieldByName(attr)
		if !e {
			panic(fmt.Errorf("attr %s is not exist on the model", attr))
		}

		dbTag := f.Tag.Get("db")
		if dbTag == "" {
			panic(fmt.Errorf("there is no db tag on the attr %s", attr))
		}

		if i != 0 {
			querySB.WriteRune(',')
		}

		querySB.WriteString(fmt.Sprintf("`%s`", dbTag))
	}
	querySB.WriteString(") VALUES ")

	for i := 0; i < rtMod.Len(); i++ {
		if i != 0 {
			querySB.WriteRune(',')
		}

		querySB.WriteRune('(')
		for j, attr := range qb.insertConfig.SavedAttribute {
			if j != 0 {
				querySB.WriteRune(',')
			}
			querySB.WriteRune('?')
			qb.params = append(qb.params, rtMod.Index(i).FieldByName(attr).Interface())
		}
		querySB.WriteRune(')')
	}

	return querySB.String()
}

/*
buildUpdate generates the SQL update query string based on the UpdateConfig provided in the sqlQueryBuilder instance.

Returns:
- string: The generated SQL update query string.

This method iterates over the SavedAttribute list in the UpdateConfig, retrieves the corresponding field values from the Model, and constructs the SET clause of the SQL update query. The field names are obtained from the 'db' tag of the struct fields.
*/
func (qb *sqlQueryBuilder) buildUpdate() string {
	query := " SET "
	for j, attr := range qb.updateConfig.SavedAttribute {
		if j != 0 {
			query += ","
		}
		var field reflect.StructField
		var value interface{}
		if reflect.TypeOf(qb.updateConfig.Model).Kind() == reflect.Ptr {
			field, _ = reflect.TypeOf(qb.updateConfig.Model).Elem().FieldByName(attr)
			value = reflect.ValueOf(qb.updateConfig.Model).Elem().FieldByName(attr).Interface()
		} else {
			field, _ = reflect.TypeOf(qb.updateConfig.Model).FieldByName(attr)
			value = reflect.ValueOf(qb.updateConfig.Model).FieldByName(attr).Interface()
		}
		qb.params = append(qb.params, value)
		query += "`" + field.Tag.Get("db") + "`" + " = ?"
	}

	return query
}

/*
buildOrder generates the ORDER BY clause for the SQL query based on the orderCondition field of the sqlQueryBuilder struct.

Returns:
- string: The ORDER BY clause for the SQL query. If orderCondition is empty, an empty string is returned.
*/
func (qb *sqlQueryBuilder) buildOrder() string {
	var query string
	if qb.orderCondition != "" {
		query += fmt.Sprintf("ORDER BY %s", qb.orderCondition)
	}
	return query
}

/*
buildGroup generates the GROUP BY clause for the SQL query based on the groupByCondition field in the sqlQueryBuilder struct.

Returns:
- string: The generated GROUP BY clause.
*/
func (qb *sqlQueryBuilder) buildGroup() string {
	var query string
	if qb.groupByCondition != "" {
		query += fmt.Sprintf("GROUP BY %s", qb.groupByCondition)
	}
	return query
}

/*
buildHaving builds the HAVING clause of the SQL query based on the stored havingCondition in the sqlQueryBuilder struct.
If havingCondition is not empty, it appends "HAVING" followed by the havingCondition to the query string and adds havingParams to the params slice.
Returns the constructed HAVING clause as a string.
*/
func (qb *sqlQueryBuilder) buildHaving() string {
	var query string
	if qb.havingCondition != "" {
		query += fmt.Sprintf("HAVING %s", qb.havingCondition)
		qb.params = append(qb.params, qb.havingParams...)
	}
	return query
}

func (qb *sqlQueryBuilder) buildJoin() string {
	var s strings.Builder

	for i, join := range qb.joinClause {
		if i != 0 && i < len(join)-1 {
			s.WriteRune(' ')
		}
		s.WriteString(join)
	}

	qb.params = append(qb.params, qb.joinParams...)

	return s.String()
}

/*
BuildSQL generates a SQL query based on the provided table name and build type.

Parameters:
- tableName (string): The name of the table to perform the SQL operation on.
- buildType (BuildType): The type of SQL operation to perform (Select, Create, Update, Delete).

Returns:
- string: The generated SQL query.
- []interface{}: The parameters to be used in the SQL query.

Example:

	query, params := qb.BuildSQL("users", Select)
	fmt.Println(query) // Output: "SELECT * FROM users"
	fmt.Println(params) // Output: []

Note:
- If a custom table name is set in the sqlQueryBuilder instance, it will override the tableName parameter.
- The returned SQL query may vary based on the buildType provided.
*/
func (qb *sqlQueryBuilder) BuildSQL(tableName string, buildType BuildType) (string, []interface{}) {
	qb.params = nil
	if qb.tableName != nil {
		tableName = *qb.tableName
		if qb.fromParams != nil {
			qb.params = append(qb.params, qb.fromParams...)
		}
	}
	var query string
	switch buildType {
	case Select:
		query = fmt.Sprintf(
			"%s %s %s %s %s %s %s %s",
			qb.buildSelect(),
			fmt.Sprintf("FROM %s", tableName),
			qb.buildJoin(),
			qb.buildWhere(),
			qb.buildGroup(),
			qb.buildHaving(),
			qb.buildOrder(),
			qb.buildLimit(),
		)
		if qb.isExist {
			query = fmt.Sprintf("SELECT EXISTS(%s)", query)
		}
	case Create:
		query = fmt.Sprintf(
			"INSERT INTO %s %s",
			tableName,
			qb.buildCreate(),
		)
	case Update:
		if qb.whereCondition == "" {
			panic("update query must have condition")
		}
		query = fmt.Sprintf(
			"UPDATE %s %s %s",
			tableName,
			qb.buildUpdate(),
			qb.buildWhere(),
		)
	case Delete:
		if qb.whereCondition == "" {
			panic("delete query must have condition")
		}
		query = fmt.Sprintf(
			"DELETE FROM %s %s",
			tableName,
			qb.buildWhere(),
		)
	case BatchCreate:
		query = fmt.Sprintf(
			"INSERT INTO %s %s",
			tableName,
			qb.buildBatchCreate(),
		)
	}
	return query, qb.params
}

/*
Build generates a SQL query based on the provided table name and build type.

Parameters:
- tableName: the name of the table to perform the query on
- buildType: the type of query to build (Select, Create, Update, Delete)

Returns:
- query: a slice of strings representing the generated SQL query
- params: a slice of interface{} containing the parameters for the query

The method handles different build types by calling specific methods to build different parts of the SQL query, such as SELECT, INSERT, UPDATE, or DELETE clauses. It then constructs the final query by combining these parts based on the provided build type.
*/
func (qb *sqlQueryBuilder) Build(tableName string, buildType BuildType) ([]string, []interface{}) {
	qb.params = nil
	if qb.tableName != nil {
		tableName = *qb.tableName
	}
	var query []string
	switch buildType {
	case Select:
		query = append(
			query,
			qb.buildSelect(),
			fmt.Sprintf("FROM %s", tableName),
			qb.buildJoin(),
			qb.buildWhere(),
			qb.buildHaving(),
			qb.buildLimit(),
		)
	case Create:
		query = append(
			query,
			fmt.Sprintf("INSERT INTO %s", tableName),
			qb.buildCreate(),
		)
	case Update:
		query = append(
			query,
			fmt.Sprintf("UPDATE %s", tableName),
			qb.buildUpdate(),
			qb.buildWhere(),
		)
	case Delete:
		query = append(
			query,
			fmt.Sprintf("DELETE FROM %s", tableName),
			qb.buildWhere(),
		)
	case BatchCreate:
		query = append(
			query,
			fmt.Sprintf("INSERT INTO %s", tableName),
			qb.buildBatchCreate(),
		)
	}
	return query, qb.params
}

/*
checkForSQLxInstance checks if the sqlx DB or Tx instance is present in the sqlQueryBuilder struct.
It returns an error if the instance is not found or if the table name is empty.

Returns:
- EmptySqlxInstance error if both db and tx are nil.
- TableNameCannotEmpty error if the table name is nil.

Returns nil if both conditions are met.
*/
func (qb *sqlQueryBuilder) checkForSQLxInstance() error {
	if qb.db == nil && qb.tx == nil {
		return EmptySqlxInstance
	}

	if qb.tableName == nil {
		return TableNameCannotEmpty
	}

	return nil
}

/*
buildForSQLx checks for the SQLx instance and then builds the SQL query based on the provided BuildType.
It calls the checkForSQLxInstance method to ensure that the SQLx instance is valid.
If the check fails, it returns an error.
It then calls the BuildSQL method to construct the SQL query based on the provided buildType.
The method returns the constructed query string and the associated parameters.

Parameters:
- buildType: The type of build operation to perform (Select, Create, Update, Delete).

Returns:
- string: The constructed SQL query.
- []interface{}: The parameters associated with the query.
- error: An error if the check for SQLx instance fails.
*/
func (qb *sqlQueryBuilder) buildForSQLx(buildType BuildType) (string, []interface{}, error) {
	if err := qb.checkForSQLxInstance(); err != nil {
		return "", nil, err
	}
	query, params := qb.BuildSQL("", buildType)
	return query, params, nil
}

/*
execDB executes the SQL query built by the sqlQueryBuilder for the specified BuildType.

Parameters:
- buildType (BuildType): The type of query to build.

Returns:
- sql.Result: The result of executing the SQL query.
- error: An error if the execution of the SQL query fails.
*/
func (qb *sqlQueryBuilder) execDB(buildType BuildType) (sql.Result, error) {
	query, params, err := qb.buildForSQLx(buildType)
	if err != nil {
		return nil, err
	}

	if qb.db != nil {
		return qb.db.Exec(query, params...)
	} else {
		return qb.tx.Exec(query, params...)
	}
}

/*
All retrieves all records from the database based on the provided destination struct.
If the select clause is not set and the isExist flag is false, it automatically selects columns based on the 'db' tag of the destination struct fields.
It then builds the SQL query using the Select build type and executes the query using either the database connection or transaction based on availability.

Warning: The From function must be called before calling this method.

Parameters:
- dest: The destination struct where the retrieved records will be mapped.

Returns:
- error: An error if the operation fails, nil otherwise.
*/
func (qb *sqlQueryBuilder) All(dest interface{}) error {
	if qb.selectClause == nil && !qb.isExist {
		columns := struct_helper.GetDBTagFromStruct(dest)
		qb.Select(columns...)
	}
	query, params, err := qb.buildForSQLx(Select)
	if err != nil {
		return err
	}

	if qb.db != nil {
		return qb.db.Select(dest, query, params...)
	} else {
		return qb.tx.Select(dest, query, params...)
	}
}

/*
One executes the query and retrieves a single row from the database.
If selectClause is nil and isExist is false, it fetches the columns from the destination struct using GetDBTagFromStruct method and sets them as the selectClause.
Then it builds the query for SQLx using buildForSQLx method with BuildType Select.
Finally, it executes the query using the database connection (db or tx) and retrieves the result into the provided destination interface.

Warning: The From function must be called before calling this method.

Parameters:
- dest: A pointer to the destination struct where the result will be scanned.

Returns:
- error: An error if any occurred during the execution of the query.
*/
func (qb *sqlQueryBuilder) One(dest interface{}) error {
	if qb.selectClause == nil && !qb.isExist {
		columns := struct_helper.GetDBTagFromStruct(dest)
		qb.Select(columns...)
	}
	query, params, err := qb.buildForSQLx(Select)
	if err != nil {
		return err
	}

	if qb.db != nil {
		return qb.db.Get(dest, query, params...)
	} else {
		return qb.tx.Get(dest, query, params...)
	}
}

/*
Create executes a SQL query to insert a new record into the database using the sqlQueryBuilder instance.
It calls the execDB method with the Create BuildType to perform the database operation.
If an error occurs during execution, it returns 0, nil, and the error.
Otherwise, it retrieves the last inserted ID from the result and returns it along with the result and nil error.

Warning: The From function must be called before calling this method.

Returns:
- int: The last inserted ID as an integer.
- sql.Result: The result of the SQL operation.
- error: Any error that occurred during execution.
*/
func (qb *sqlQueryBuilder) Create() (int, sql.Result, error) {
	res, err := qb.execDB(Create)
	if err != nil {
		return 0, nil, err
	}

	insertedID, _ := res.LastInsertId()

	return int(insertedID), res, nil
}

/*
Update executes an update query using the sqlQueryBuilder instance.
It calls the execDB method with the BuildType Update to build and execute the update query.

Warning: The From function must be called before calling this method.

Returns:
- sql.Result: The result of the update query execution.
- error: An error if the update query execution fails.
*/
func (qb *sqlQueryBuilder) Update() (sql.Result, error) {
	if qb.whereCondition == "" {
		return nil, fmt.Errorf("update must have condition")
	}
	return qb.execDB(Update)
}

/*
Delete executes a delete operation in the database using the provided SQL query builder. It calls the execDB method with the BuildType set to Delete, which generates the appropriate SQL query for deletion.

Warning: The From function must be called before calling this method.

Returns:
- sql.Result: The result of the delete operation.
- error: An error if the delete operation fails.
*/
func (qb *sqlQueryBuilder) Delete() (sql.Result, error) {
	if qb.whereCondition == "" {
		return nil, fmt.Errorf("delete must have condition")
	}
	return qb.execDB(Delete)
}

/*
Exist checks if a record exists based on the current query builder configuration.
It sets the isExist flag to true and calls the One method to retrieve a single result.
If an error occurs during the retrieval process, it returns false and the error.
If the record exists, it returns true and nil.

Warning: The From function must be called before calling this method.
*/
func (qb *sqlQueryBuilder) Exist() (bool, error) {
	qb.isExist = true

	var exist bool
	if err := qb.One(&exist); err != nil {
		return false, err
	}
	return exist, nil
}

/*
Count returns the total count of records based on the current query builder configuration.

It sets the select clause to "COUNT(*)" and executes the query to fetch the count.
The count value is stored in the provided integer variable.

Warning: The From function must be called before calling this method.

Returns:
- int: The total count of records.
- error: An error if the query execution fails.
*/
func (qb *sqlQueryBuilder) Count() (int, error) {
	if qb.selectClause == nil {
		qb.Select("COUNT(*)")
	}

	var count int
	if err := qb.One(&count); err != nil {
		return 0, err
	}
	return count, nil
}

/*
BatchCreate executes a SQL query to insert many new record into the database using the sqlQueryBuilder instance.
It calls the execDB method with the BatchCreate BuildType to perform the database operation.
If an error occurs during execution, it returns nil and the error.
Otherwise, it retrieves the result and nil error.

Warning: The From function must be called before calling this method.

Returns:
- sql.Result: The result of the batch create operation.
- error: An error if the batch create operation fails.
*/
func (qb *sqlQueryBuilder) BatchCreate() (sql.Result, error) {
	res, err := qb.execDB(BatchCreate)
	if err != nil {
		return nil, err
	}

	return res, nil
}

/*
BuildRawSQL generates a raw SQL query string by building the SQL query based on the specified
build type and replacing placeholders with formatted parameter values. It handles various
data types, including nullable types, to ensure correct SQL syntax. The method does not
affect the existing query flow by working on a dereferenced copy of the QueryBuilder.

Deprecated: READ THE WARNING LABEL.

Warning: DO NOT USE the query returned from this function as query to be run on the DB. BEWARE OF SQL Injection. FOR DEBUGGING PURPOSE ONLY.

Parameters:
- tableName (string): The name of the table to perform the SQL operation on or nil if already set the table name using From function.
- buildType (BuildType): The type of SQL operation to perform (Select, Create, Update, Delete).

Returns:
- string: The raw query.
*/
func (qb *sqlQueryBuilder) BuildRawSQL(tableName string, buildType BuildType) string {
	q := *qb // dereference Query Builder, so it doesn't affect existing query flow.
	query, params := q.BuildSQL(tableName, buildType)

	for _, param := range params {
		var formattedValue string
		rv := reflect.ValueOf(param)
		if rv.Type().Kind() == reflect.Pointer {
			rv = rv.Elem()
		}

		switch rv.Type().Kind() {
		case reflect.String:
			formattedValue = fmt.Sprintf("'%s'", rv.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			formattedValue = fmt.Sprintf("%d", rv.Int())
		case reflect.Float32, reflect.Float64:
			formattedValue = fmt.Sprintf("%f", rv.Float())
		case reflect.Bool:
			formattedValue = fmt.Sprintf("%t", rv.Bool())
		case reflect.Struct:
			switch rv.Type() {
			case reflect.TypeOf(null.String{}), reflect.TypeOf(nullable.String{}):
				formattedValue = fmt.Sprintf("'%s'", rv.FieldByName("String").String())
			case reflect.TypeOf(null.Int{}), reflect.TypeOf(nullable.Int{}):
				formattedValue = fmt.Sprintf("%d", rv.FieldByName("Int64").Int())
			case reflect.TypeOf(null.Float{}), reflect.TypeOf(nullable.Float{}):
				formattedValue = fmt.Sprintf("%f", rv.FieldByName("Float64").Float())
			case reflect.TypeOf(null.Bool{}), reflect.TypeOf(nullable.Bool{}):
				formattedValue = fmt.Sprintf("%t", rv.FieldByName("Bool").Bool())
			default:
				formattedValue = fmt.Sprintf("'%v'", rv.Interface())
			}
		default:
			formattedValue = fmt.Sprintf("'%v'", rv.Interface())
		}

		// Replace the first occurrence of '?' with the formatted value
		query = strings.Replace(query, "?", formattedValue, 1)
	}

	return query
}

/*
internalNewQBInstance creates a new instance of a QueryBuilder with the provided *sqlx.DB and *sqlx.Tx.
It initializes a sqlQueryBuilder struct, sets the db and tx fields, and returns the QueryBuilder interface.

Parameters:
- db: *sqlx.DB - The database connection to be used by the QueryBuilder.
- tx: *sqlx.Tx - The database transaction to be used by the QueryBuilder.

Returns:
- QueryBuilder - The newly created QueryBuilder instance.

Example:
qb := internalNewQBInstance(db, tx)
*/
func internalNewQBInstance(db *sqlx.DB, tx *sqlx.Tx) QueryBuilder {
	qb := &sqlQueryBuilder{}
	qb.db = db
	qb.tx = tx
	return qb
}

/*
New returns a new instance of a QueryBuilder without database and transaction values.
*/
func New() QueryBuilder {
	return internalNewQBInstance(nil, nil)
}

/*
NewWithDB creates a new instance of QueryBuilder with the provided sqlx.DB connection.
It initializes the QueryBuilder with the given database connection and a nil transaction.
The returned QueryBuilder instance can be used to construct SQL queries and perform various database operations.

Parameters:
- db: The sqlx.DB connection to be used by the QueryBuilder.

Returns:
- QueryBuilder: A new instance of QueryBuilder initialized with the provided database connection.

Example:

	qb := NewWithDB(db)


	qb.Select("name", "age").From("users").Where("age", ">", 18).OrderBy("name").Limit(10, 0).All(&users)
*/
func NewWithDB(db *sqlx.DB) QueryBuilder {
	return internalNewQBInstance(db, nil)
}

/*
NewWithTx creates a new instance of QueryBuilder that operates within the provided transaction.
It initializes the QueryBuilder with a nil database connection and the given transaction.

Parameters:
- tx: The transaction object to be used by the QueryBuilder.

Returns:
- QueryBuilder: A new instance of QueryBuilder configured to use the provided transaction.
*/
func NewWithTx(tx *sqlx.Tx) QueryBuilder {
	return internalNewQBInstance(nil, tx)
}
