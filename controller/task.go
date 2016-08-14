package controller

import (
	"net/http"

	"github.com/voyagegroup/gin-boilerplate/model"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Todo はTodoへのリクエストに関する制御をします
type Todo struct {
	DB *sqlx.DB
}

// GetはDBからユーザを取得して結果を返します
func (t *Todo) Get(c *gin.Context) {
	todos, err := model.TodosAll(t.DB)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (t *Todo) Post(c *gin.Context) {
	var todo model.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}

	tx, err := t.DB.Beginx()
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	result, err := todo.Update(tx)
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	todo.ID = id
	c.JSON(http.StatusOK, todo)
}

// PutはタスクをDBに追加します
// todoをJSONとして受け取ることを想定しています。
func (t *Todo) Put(c *gin.Context) {
	var todo model.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}

	tx, err := t.DB.Beginx()
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	result, err := todo.Insert(tx)
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	todo.ID = id
	c.JSON(http.StatusCreated, todo)
	return
}

func (t *Todo) Delete(c *gin.Context) {
	var todo model.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}

	tx, err := t.DB.Beginx()
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if _, err := todo.Delete(tx); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (t *Todo) DeleteMulti(c *gin.Context) {
	var todos []model.Todo
	if err := c.BindJSON(&todos); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}

	tx, err := t.DB.Beginx()
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			return
		}
	}()

	for _, todo := range todos {
		if _, err := todo.Delete(tx); err != nil {
			c.JSON(500, gin.H{"err": err.Error()})
			return
		}
	}
	tx.Commit()
	c.Status(http.StatusOK)
}

func (t *Todo) Toggle(c *gin.Context) {
	var todo model.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}

	tx, err := t.DB.Beginx()
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	result, err := todo.Toggle(tx)
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	if n, err := result.RowsAffected(); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	} else if n != 1 {
		c.JSON(500, gin.H{"err": "update failed"})
		return
	}
	c.Status(http.StatusOK)
}

func (t *Todo) ToggleAll(c *gin.Context) {
	var req = struct {
		Checked bool `json:"checked"`
	}{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}

	tx, err := t.DB.Beginx()
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if _, err := model.TodosToggleAll(tx, req.Checked); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
