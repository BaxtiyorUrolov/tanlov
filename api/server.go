package api

import (
	"context"
	"it-tanlov/api/handler"
	"it-tanlov/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	httpHeaderTimeout = time.Second * 10
	shutdownTimeout   = time.Second * 5
)

type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	handler    handler.Handler
}

func New(
	cfg config.Config,
	h handler.Handler,
) *Server {
	r := gin.New()
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		cors(),
	)

	s := &Server{
		router:  r,
		handler: h,
	}
	s.endpoints()

	s.httpServer = &http.Server{
		Addr:              cfg.HTTPPort,
		Handler:           r,
		ReadHeaderTimeout: httpHeaderTimeout,
	}

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// routers
func (s *Server) endpoints() {
	s.router.POST("/partner", s.handler.CreatePartner)
	s.router.GET("/partners", s.handler.GetPartnerList)

	// Swagger documentation route
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, GET")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With") //nolint:lll
		c.Header("Access-Control-Max-Age", "3600")
		c.Header("Content-Type", "application/json")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
