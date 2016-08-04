package base

import (
	"log"
	"net/http"

	"github.com/voyagegroup/gin-boilerplate/controller"
	"github.com/voyagegroup/gin-boilerplate/db"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Serverはベースアプリケーションのserverを示します
//
// TODO: dbxをstructから分離したほうが複数人数開発だと見通しがよいかもしれない
type Server struct {
	dbx    *sqlx.DB
	Engine *gin.Engine
}

func (s *Server) Close() error {
	return s.dbx.Close()
}

// InitはServerを初期化する
func (s *Server) Init(dbconf, env string) {
	cs, err := db.NewConfigsFromFile(dbconf)
	if err != nil {
		log.Fatalf("cannot open database configuration. exit. %s", err)
	}
	dbx, err := cs.Open(env)
	if err != nil {
		log.Fatalf("db initialization failed: %s", err)
	}
	s.dbx = dbx
	s.Route()
}

// Newはベースアプリケーションを初期化します
func New() *Server {
	r := gin.Default()
	return &Server{Engine: r}
}

func (s *Server) Run(addr ...string) {
	s.Engine.Run(addr...)
}

// Routeはベースアプリケーションのroutingを設定します
//
// TODO muxを返すようにしてroutingのテストをしやすくする
func (s *Server) Route() {
	// ヘルスチェック用
	// TODO opsチームとnginxからのhealthcheck用のendpointについてあわせておく
	s.Engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", "pong")
	})

	todo := &controller.Todo{DB: s.dbx}
	s.Engine.GET("/api/todos", todo.Get)
	s.Engine.PUT("/api/todos", todo.Put)
	s.Engine.POST("/api/todos", todo.Post)
	s.Engine.DELETE("/api/todos", todo.Delete)

	s.Engine.DELETE("/api/todos/multi", todo.DeleteMulti)

	s.Engine.POST("/api/todos/toggle", todo.Toggle)
	s.Engine.POST("/api/todos/toggleall", todo.ToggleAll)

	s.Engine.StaticFile("/", "public/index.html")
	s.Engine.Static("/static", "public")
}
