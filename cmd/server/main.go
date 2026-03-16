package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	// DefaultPort はデフォルトのサーバーポートです。
	DefaultPort = "8080"
)

func main() {
	e := setupServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	// グレースフルシャットダウン用チャネル
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
			log.Printf("サーバー停止: %v", err)
		}
	}()

	log.Printf("サーバー起動: http://localhost:%s", port)
	<-quit
	log.Println("サーバーをシャットダウンしています...")
}

// setupServer はEchoサーバーを構成して返します。
func setupServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// ミドルウェア設定
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Printf("%s %s %d", v.Method, v.URI, v.Status)
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// 静的ファイル配信
	e.Static("/static", "static")

	// ヘルスチェックエンドポイント
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	return e
}
