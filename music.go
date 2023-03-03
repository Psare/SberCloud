package main

import (
	"container/list"
	"path/filepath"
	"net/http"
	"os/signal"
	"sync"
	"time"
	"log"
	"os"
	"syscall"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gin-gonic/gin"
)

type MusicPlayer interface {
    Play() error
    Pause()
    Resume()
    Next()
    Prev()
    AddSong(path string)
    Delete()
}

type Song struct {
    path     string
    position *list.Element
}

type Player struct {
    playlist  *list.List
    current   *list.Element
    paused    bool
    pauseCond *sync.Cond
	duration float64
}

func NewPlayer() *Player {
    p := &Player{
        playlist:  list.New(),
        paused:    false,
        pauseCond: sync.NewCond(&sync.Mutex{}),
    }
    return p
}

func (p *Player) Play() error {
    if p.current == nil {
        return nil
    }
    f, err := os.Open(p.current.Value.(*Song).path)
    if err != nil {
        return err
    }
    s, format, err := mp3.Decode(f)
    if err != nil {
        return err
    }
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	p.duration = float64(s.Len()) / float64(format.SampleRate)
    speaker.Play(s)
	return err
}

func (p *Player) Pause() {
	if p.paused != false {
		return
	}
    p.pauseCond.L.Lock()
    defer p.pauseCond.L.Unlock()
    p.paused = true
    speaker.Lock()
}

func (p *Player) Resume() {
	if p.paused != true {
		return
	}
    p.pauseCond.L.Lock()
    defer p.pauseCond.L.Unlock()
    p.paused = false
    speaker.Unlock()
}

func (p *Player) AddSong(path string) {
	ext := filepath.Ext(path)
    if ext != ".mp3" {
        return
    }

    
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return
    }
    if err := os.Chmod(path, 0644); err != nil {
        return
    }
    song := &Song{path: path}
    song.position = p.playlist.PushBack(song)
}

func (p *Player) Delete() {
    if p.current != nil {
        speaker.Clear()
        p.playlist.Remove(p.current)
		p.Prev()
    }
}

func (p *Player) Next() {
    if p.current == nil {
        return
    }
    if p.current.Next() == nil {
        p.current = p.playlist.Front()
    } else {
		p.current = p.current.Next()
	}
    speaker.Clear()
    p.Play()
}

func (p *Player) Prev() {
    if p.current == nil {
        return
    }
    if p.current.Prev() == nil {
        
        p.current = p.playlist.Back()
    } else {
        p.current = p.current.Prev()
    }
    speaker.Clear()
    p.Play()
}

func (p *Player) PlayPlaylist() {
    p.current = p.playlist.Front()
    for p.current != nil {
        p.Play()
		for p.duration != 0{
			if p.paused != true{
				p.duration = p.duration - 1
			}
			time.Sleep(1 * time.Second)
		}
        p.Next()
    }
}

func main() {
	quit := make(chan os.Signal, 1)
    p := NewPlayer()

	

	r := gin.Default()
	r.LoadHTMLFiles("sait.html")

	
	r.GET("", func(ctx *gin.Context) {
		songs := make([]string, 0)
        for e := p.playlist.Front(); e != nil; e = e.Next() {
            songs = append(songs, e.Value.(*Song).path)
        }
		ctx.HTML(http.StatusOK, "sait.html", gin.H{
            "songs": songs,
        })
	})
	r.GET("/playplaylist", func(c *gin.Context) {
		go p.PlayPlaylist()
	})

	
	r.GET("/pause", func(c *gin.Context) {
		p.Pause()
	})

	
	r.GET("/resume", func(c *gin.Context) {
		p.Resume()
	})
	r.GET("/prev", func(c *gin.Context) {
		p.Prev()
	})
	r.GET("/next", func(c *gin.Context) {
		p.Next()
	})

	
	r.POST("/addSong", func(c *gin.Context) {
		path := c.PostForm("path")
		p.AddSong(path)
		c.Redirect(http.StatusSeeOther, "/")
	})

	
	r.GET("/delete", func(c *gin.Context) {
		p.Delete()
	})

	r.Run()
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("app Shutting down")
	os.Exit(0)
}