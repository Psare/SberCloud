package player

import (
	"container/list"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
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
	Path     string
	Position *list.Element
}

type Player struct {
	Playlist  *list.List
	Current   *list.Element
	Paused    bool
	PauseCond *sync.Cond
	Duration  float64
}

func NewPlayer() *Player {
	p := &Player{
		Playlist:  list.New(),
		Paused:    false,
		PauseCond: sync.NewCond(&sync.Mutex{}),
	}
	return p
}

func (p *Player) Play() error {
	if p.Current == nil {
		return nil
	}
	f, err := os.Open(p.Current.Value.(*Song).Path)
	if err != nil {
		return err
	}
	s, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	p.Duration = float64(s.Len()) / float64(format.SampleRate)
	speaker.Play(s)
	return nil
}

func (p *Player) Pause() {
	if p.Paused != false {
		return
	}
	p.PauseCond.L.Lock()
	defer p.PauseCond.L.Unlock()
	p.Paused = true
	speaker.Lock()
}

func (p *Player) Resume() {
	if p.Paused != true {
		return
	}
	p.PauseCond.L.Lock()
	defer p.PauseCond.L.Unlock()
	p.Paused = false
	speaker.Unlock()
}

func (p *Player) AddSong(path string) {
	path = "/Users/macbook/Desktop/SberCloud/music/" + path
	ext := filepath.Ext(path)
	if ext != ".mp3" {
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
		return
	}
	if err := os.Chmod(path, 0644); err != nil {
		fmt.Println(err)
		return
	}
	song := &Song{Path: path}
	song.Position = p.Playlist.PushBack(song)
}

func (p *Player) Delete() {
	if p.Current != nil {
		speaker.Clear()
		p.Playlist.Remove(p.Current)
		p.Prev()
	}
}

func (p *Player) Next() {
	if p.Current == nil {
		return
	}
	if p.Current.Next() == nil {
		p.Current = p.Playlist.Front()
	} else {
		p.Current = p.Current.Next()
	}
	speaker.Clear()
	p.Play()
}

func (p *Player) Prev() {
	if p.Current == nil {
		return
	}
	if p.Current.Prev() == nil {
		p.Current = p.Playlist.Back()
	} else {
		p.Current = p.Current.Prev()
	}
	speaker.Clear()
	p.Play()
}

func (p *Player) PlayPlaylist() {
	p.Current = p.Playlist.Front()
	for p.Current != nil {
		p.Play()
		for p.Duration != 0 {
			if p.Paused != true {
				p.Duration = p.Duration - 1
			}
			time.Sleep(1 * time.Second)
		}
		p.Next()
	}
}
