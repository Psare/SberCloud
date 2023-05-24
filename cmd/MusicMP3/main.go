package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"sber/internal/handlers"
	"sber/internal/player"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	p := player.NewPlayer()
	h := handlers.NewHandler(p)

	r := gin.Default()
	r.LoadHTMLFiles("/Users/macbook/Desktop/SberCloud/web/sait.html")

	r.GET("", h.RootHandler)
	r.GET("/playplaylist", h.PlayPlaylistHandler)
	r.GET("/pause", h.PauseHandler)
	r.GET("/resume", h.ResumeHandler)
	r.GET("/prev", h.PrevHandler)
	r.GET("/next", h.NextHandler)
	r.POST("/addSong", h.AddSongHandler)
	r.GET("/delete", h.DeleteHandler)

	go func() {
		if err := r.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	<-quit
	log.Println("App shutting down")
	os.Exit(0)
}