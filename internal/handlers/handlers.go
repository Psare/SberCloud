package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sber/internal/player"
)

type Handler struct {
	player *player.Player
}

func NewHandler(p *player.Player) *Handler {
	return &Handler{
		player: p,
	}
}

func (h *Handler) RootHandler(ctx *gin.Context) {
	songs := make([]string, 0)
	for e := h.player.Playlist.Front(); e != nil; e = e.Next() {
		songs = append(songs, e.Value.(*player.Song).Path)
	}
	ctx.HTML(http.StatusOK, "sait.html", gin.H{
		"songs": songs,
	})
}

func (h *Handler) PlayPlaylistHandler(c *gin.Context) {
	go h.player.PlayPlaylist()
}

func (h *Handler) PauseHandler(c *gin.Context) {
	h.player.Pause()
}

func (h *Handler) ResumeHandler(c *gin.Context) {
	h.player.Resume()
}

func (h *Handler) PrevHandler(c *gin.Context) {
	h.player.Prev()
}

func (h *Handler) NextHandler(c *gin.Context) {
	h.player.Next()
}

func (h *Handler) AddSongHandler(c *gin.Context) {
	path := c.PostForm("path")
	h.player.AddSong(path)
	c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) DeleteHandler(c *gin.Context) {
	h.player.Delete()
}