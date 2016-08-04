package model

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Todoは管理するタスク
type Todo struct {
	ID        int64      `db:"todo_id" json:"id"`
	Title     string     `json:"title"`
	Completed bool       `json:"completed"`
	Created   *time.Time `json:"created"`
	Updated   *time.Time `json:"updated"`
}

func TodosAll(dbx *sqlx.DB) (todos []Todo, err error) {
	if err := dbx.Select(&todos, "select * from todos"); err != nil {
		return nil, err
	}
	return todos, nil
}

func TodoOne(dbx *sqlx.DB, id int64) (*Todo, error) {
	var todo Todo
	if err := dbx.Get(&todo, `
	select * from todos where todo_id = ?
	`, id); err != nil {
		return nil, err
	}
	return &todo, nil
}

// TodosToggleAllは全部のtoggleのステータスをトグルします
func TodosToggleAll(tx *sqlx.Tx, checked bool) (sql.Result, error) {
	stmt, err := tx.Prepare(`
	update todos set completed = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(checked)
}

func (t *Todo) Update(tx *sqlx.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`
	update todos set title = ? where todo_id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(t.Title, t.ID)
}

func (t *Todo) Insert(tx *sqlx.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`
	insert into todos (title, completed)
	values(?, ?)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(t.Title, t.Completed)
}

// Toggle は指定されたタスクについて現在の状態と入れ替えます。
func (t *Todo) Toggle(tx *sqlx.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`
	update todos set completed=?
	where todo_id=?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(!t.Completed, t.ID)
}

func (t *Todo) Delete(tx *sqlx.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`delete from todos where todo_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(t.ID)
}

// TodosDeleteAllはすべてのタスクを消去します。
// テストのために使用されます。
func TodosDeleteAll(tx *sqlx.Tx) (sql.Result, error) {
	return tx.Exec(`truncate table todos`)
}
