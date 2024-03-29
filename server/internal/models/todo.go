// Code generated by SQLBoiler 4.14.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Todo is an object representing the database table.
type Todo struct {
	ID          int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Subject     string      `boil:"subject" json:"subject" toml:"subject" yaml:"subject"`
	Description null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	AppUserID   int64       `boil:"app_user_id" json:"-" toml:"-" yaml:"-"`

	R *todoR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L todoL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var TodoColumns = struct {
	ID          string
	Subject     string
	Description string
	AppUserID   string
}{
	ID:          "id",
	Subject:     "subject",
	Description: "description",
	AppUserID:   "app_user_id",
}

var TodoTableColumns = struct {
	ID          string
	Subject     string
	Description string
	AppUserID   string
}{
	ID:          "todo.id",
	Subject:     "todo.subject",
	Description: "todo.description",
	AppUserID:   "todo.app_user_id",
}

// Generated where

var TodoWhere = struct {
	ID          whereHelperint64
	Subject     whereHelperstring
	Description whereHelpernull_String
	AppUserID   whereHelperint64
}{
	ID:          whereHelperint64{field: "`todo`.`id`"},
	Subject:     whereHelperstring{field: "`todo`.`subject`"},
	Description: whereHelpernull_String{field: "`todo`.`description`"},
	AppUserID:   whereHelperint64{field: "`todo`.`app_user_id`"},
}

// TodoRels is where relationship names are stored.
var TodoRels = struct {
	AppUser string
}{
	AppUser: "AppUser",
}

// todoR is where relationships are stored.
type todoR struct {
	AppUser *AppUser `boil:"AppUser" json:"AppUser" toml:"AppUser" yaml:"AppUser"`
}

// NewStruct creates a new relationship struct
func (*todoR) NewStruct() *todoR {
	return &todoR{}
}

func (r *todoR) GetAppUser() *AppUser {
	if r == nil {
		return nil
	}
	return r.AppUser
}

// todoL is where Load methods for each relationship are stored.
type todoL struct{}

var (
	todoAllColumns            = []string{"id", "subject", "description", "app_user_id"}
	todoColumnsWithoutDefault = []string{"subject", "description", "app_user_id"}
	todoColumnsWithDefault    = []string{"id"}
	todoPrimaryKeyColumns     = []string{"id"}
	todoGeneratedColumns      = []string{}
)

type (
	// TodoSlice is an alias for a slice of pointers to Todo.
	// This should almost always be used instead of []Todo.
	TodoSlice []*Todo

	todoQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	todoType                 = reflect.TypeOf(&Todo{})
	todoMapping              = queries.MakeStructMapping(todoType)
	todoPrimaryKeyMapping, _ = queries.BindMapping(todoType, todoMapping, todoPrimaryKeyColumns)
	todoInsertCacheMut       sync.RWMutex
	todoInsertCache          = make(map[string]insertCache)
	todoUpdateCacheMut       sync.RWMutex
	todoUpdateCache          = make(map[string]updateCache)
	todoUpsertCacheMut       sync.RWMutex
	todoUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single todo record from the query.
func (q todoQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Todo, error) {
	o := &Todo{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for todo")
	}

	return o, nil
}

// All returns all Todo records from the query.
func (q todoQuery) All(ctx context.Context, exec boil.ContextExecutor) (TodoSlice, error) {
	var o []*Todo

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Todo slice")
	}

	return o, nil
}

// Count returns the count of all Todo records in the query.
func (q todoQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count todo rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q todoQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if todo exists")
	}

	return count > 0, nil
}

// AppUser pointed to by the foreign key.
func (o *Todo) AppUser(mods ...qm.QueryMod) appUserQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`id` = ?", o.AppUserID),
	}

	queryMods = append(queryMods, mods...)

	return AppUsers(queryMods...)
}

// LoadAppUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (todoL) LoadAppUser(ctx context.Context, e boil.ContextExecutor, singular bool, maybeTodo interface{}, mods queries.Applicator) error {
	var slice []*Todo
	var object *Todo

	if singular {
		var ok bool
		object, ok = maybeTodo.(*Todo)
		if !ok {
			object = new(Todo)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeTodo)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeTodo))
			}
		}
	} else {
		s, ok := maybeTodo.(*[]*Todo)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeTodo)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeTodo))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &todoR{}
		}
		args = append(args, object.AppUserID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &todoR{}
			}

			for _, a := range args {
				if a == obj.AppUserID {
					continue Outer
				}
			}

			args = append(args, obj.AppUserID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`app_user`),
		qm.WhereIn(`app_user.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load AppUser")
	}

	var resultSlice []*AppUser
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice AppUser")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for app_user")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for app_user")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.AppUser = foreign
		if foreign.R == nil {
			foreign.R = &appUserR{}
		}
		foreign.R.Todos = append(foreign.R.Todos, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.AppUserID == foreign.ID {
				local.R.AppUser = foreign
				if foreign.R == nil {
					foreign.R = &appUserR{}
				}
				foreign.R.Todos = append(foreign.R.Todos, local)
				break
			}
		}
	}

	return nil
}

// SetAppUser of the todo to the related item.
// Sets o.R.AppUser to related.
// Adds o to related.R.Todos.
func (o *Todo) SetAppUser(ctx context.Context, exec boil.ContextExecutor, insert bool, related *AppUser) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `todo` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"app_user_id"}),
		strmangle.WhereClause("`", "`", 0, todoPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.AppUserID = related.ID
	if o.R == nil {
		o.R = &todoR{
			AppUser: related,
		}
	} else {
		o.R.AppUser = related
	}

	if related.R == nil {
		related.R = &appUserR{
			Todos: TodoSlice{o},
		}
	} else {
		related.R.Todos = append(related.R.Todos, o)
	}

	return nil
}

// Todos retrieves all the records using an executor.
func Todos(mods ...qm.QueryMod) todoQuery {
	mods = append(mods, qm.From("`todo`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`todo`.*"})
	}

	return todoQuery{q}
}

// FindTodo retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindTodo(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*Todo, error) {
	todoObj := &Todo{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `todo` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, todoObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from todo")
	}

	return todoObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Todo) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no todo provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(todoColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	todoInsertCacheMut.RLock()
	cache, cached := todoInsertCache[key]
	todoInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			todoAllColumns,
			todoColumnsWithDefault,
			todoColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(todoType, todoMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(todoType, todoMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `todo` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `todo` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `todo` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, todoPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into todo")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == todoMapping["id"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for todo")
	}

CacheNoHooks:
	if !cached {
		todoInsertCacheMut.Lock()
		todoInsertCache[key] = cache
		todoInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Todo.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Todo) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	var err error
	key := makeCacheKey(columns, nil)
	todoUpdateCacheMut.RLock()
	cache, cached := todoUpdateCache[key]
	todoUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			todoAllColumns,
			todoPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return errors.New("models: unable to update todo, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `todo` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, todoPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(todoType, todoMapping, append(wl, todoPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	_, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update todo row")
	}

	if !cached {
		todoUpdateCacheMut.Lock()
		todoUpdateCache[key] = cache
		todoUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q todoQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for todo")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o TodoSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), todoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `todo` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, todoPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in todo slice")
	}

	return nil
}

var mySQLTodoUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Todo) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no todo provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(todoColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLTodoUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	todoUpsertCacheMut.RLock()
	cache, cached := todoUpsertCache[key]
	todoUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			todoAllColumns,
			todoColumnsWithDefault,
			todoColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			todoAllColumns,
			todoPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("models: unable to upsert todo, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`todo`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `todo` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(todoType, todoMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(todoType, todoMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for todo")
	}

	var lastID int64
	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == todoMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(todoType, todoMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for todo")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for todo")
	}

CacheNoHooks:
	if !cached {
		todoUpsertCacheMut.Lock()
		todoUpsertCache[key] = cache
		todoUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Todo record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Todo) Delete(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil {
		return errors.New("models: no Todo provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), todoPrimaryKeyMapping)
	sql := "DELETE FROM `todo` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from todo")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q todoQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if q.Query == nil {
		return errors.New("models: no todoQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from todo")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o TodoSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), todoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `todo` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, todoPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from todo slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Todo) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindTodo(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *TodoSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := TodoSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), todoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `todo`.* FROM `todo` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, todoPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in TodoSlice")
	}

	*o = slice

	return nil
}

// TodoExists checks if the Todo row exists.
func TodoExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `todo` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if todo exists")
	}

	return exists, nil
}

// Exists checks if the Todo row exists.
func (o *Todo) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return TodoExists(ctx, exec, o.ID)
}
