// +build integration

package base

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/voyagegroup/gin-boilerplate/model"

	"github.com/gin-gonic/gin"
)

func defaultServer() *Server {
	s := &Server{Engine: gin.New()}
	s.Init("dbconfig.yml", "test")
	return s
}

func tearDown() {
	s := defaultServer()
	tx := s.dbx.MustBegin()
	model.TodosDeleteAll(tx)
	tx.Commit()
}

func TestMain(m *testing.M) {
	r := m.Run()
	tearDown()
	os.Exit(r)
}

func TestHealthCheck(t *testing.T) {
	s := defaultServer()

	ts := httptest.NewServer(s.Engine)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Fatalf("ping failed %s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body error %s", err)
	}
	if string(b) != "pong" {
		t.Fatalf("response should be pong but actual %s", string(b))
	}
}

func TestPostGetDelete(t *testing.T) {
	s := defaultServer()
	ts := httptest.NewServer(s.Engine)
	defer ts.Close()

	// create new todos
	req, err := http.NewRequest("PUT", ts.URL+"/api/todos",
		strings.NewReader(`{"title": "test", "completed": false}`))
	if err != nil {
		t.Fatalf("create put request failed: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post failed: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create todo failed. got %d response.", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %s", err)
	}

	var todo model.Todo
	if err := json.Unmarshal(b, &todo); err != nil {
		t.Fatalf("unmarshal failed: %s", err)
	}
	if todo.Title != "test" {
		t.Fatal("response is different: want test, got %s", todo.Title)
	}

	// test get
	respGet, err := http.Get(ts.URL + "/api/todos")
	if err != nil {
		t.Fatalf("get todos failed: %s", err)
	}
	defer respGet.Body.Close()
	if respGet.StatusCode != http.StatusOK {
		t.Fatalf("get todos status code is unexpected: want %d got %d",
			http.StatusOK, respGet.StatusCode)
	}

	bb, err := ioutil.ReadAll(respGet.Body)
	if err != nil {
		t.Fatalf("read body failed: %s", err)
	}
	var gotTodos []model.Todo
	if err := json.Unmarshal(bb, &gotTodos); err != nil {
		t.Fatalf("unmarshal failed: %s", err)
	}
	if len(gotTodos) != 1 {
		t.Fatal("todo length want 1 got %d", len(gotTodos))
	}
	if got := gotTodos[0].Title; got != "test" {
		t.Fatal("response is different: want test, got %s", got)
	}

	bbb, err := json.Marshal(gotTodos[0])
	if err != nil {
		t.Fatalf("marshal posted todo failed: %s", err)
	}

	// and delete it
	reqDel, err := http.NewRequest("DELETE", ts.URL+"/api/todos",
		bytes.NewReader(bbb))
	if err != nil {
		t.Fatalf("create put request failed: %s", err)
	}
	reqDel.Header.Set("Content-Type", "application/json")

	respDel, err := http.DefaultClient.Do(reqDel)
	if err != nil {
		t.Fatalf("post failed: %s", err)
	}
	defer respDel.Body.Close()

	if respDel.StatusCode != 200 {
		t.Fatalf("delete todo failed: status want %d got %d", 200, respDel.StatusCode)
	}
}
