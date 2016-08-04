// +build integration

package model

import (
	"os"
	"testing"

	"github.com/voyagegroup/gin-boilerplate/db"

	"github.com/jmoiron/sqlx"
)

func DefaultDB() *sqlx.DB {
	configs, err := db.NewConfigsFromFile("../dbconfig.yml")
	if err != nil {
		panic(err)
	}
	dbx, err := configs.Open("test")
	if err != nil {
		panic(err)
	}
	return dbx
}

func setUp() {}

func tearDown() {
	dbx := DefaultDB()
	tx := dbx.MustBegin()
	TodosDeleteAll(tx)
	tx.Commit()
}

func TestMain(m *testing.M) {
	setUp()
	r := m.Run()
	tearDown()
	os.Exit(r)
}

func TestTodoCRUD(t *testing.T) {
	dbx := DefaultDB()
	defer dbx.Close()

	todos, err := TodosAll(dbx)
	if err != nil {
		t.Fatalf("select todos failed: %s", err)
	}
	if len(todos) != 0 {
		t.Fatalf("len(todos) want 0 got %d", len(todos))
	}

	tx := dbx.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	todo := &Todo{Title: "homework"}
	if _, err := todo.Insert(tx); err != nil {
		t.Fatalf("insertion error: %s", err)
	}
	tx.Commit()

	afterTodos, err := TodosAll(dbx)
	if err != nil {
		t.Fatalf("select todos failed: %s", err)
	} else if len(afterTodos) != 1 {
		t.Fatalf("len(todos) want 1 got %d", len(afterTodos))
	}

	tx2 := dbx.MustBegin()
	if _, err := afterTodos[0].Delete(tx2); err != nil {
		t.Fatalf("delete todo failed: %s", err)
	}
	tx2.Commit()

	afterTodos, err = TodosAll(dbx)
	if err != nil {
		t.Fatalf("select todos failed: %s", err)
	} else if len(afterTodos) != 0 {
		t.Fatalf("len(todos) want 0 got %d", len(afterTodos))
	}
}

func TestTodoToggle(t *testing.T) {
	dbx := DefaultDB()
	defer dbx.Close()

	tx := dbx.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	todo := &Todo{Title: "homework", Completed: true}
	result, err := todo.Insert(tx)
	if err != nil {
		t.Fatalf("insertion error: %s", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get id failed: %s", err)
	}
	tx.Commit()

	tx2 := dbx.MustBegin()
	insertedTodo, err := TodoOne(dbx, id)
	if err != nil {
		t.Fatalf("get todo failed: %s", err)
	}
	if _, err := insertedTodo.Toggle(tx2); err != nil {
		t.Fatalf("toggle todo failed: %s", err)
	}
	tx2.Commit()

	afterTodo, err := TodoOne(dbx, insertedTodo.ID)
	if err != nil {
		t.Fatalf("get todo failed: %s", err)
	}
	if afterTodo.Completed {
		t.Fatalf("todo should be toggled, but not.: %v", afterTodo)
	}
}
