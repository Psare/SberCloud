package main

import (
	"os"
	_"path/filepath"
	"testing"
	"time"
)

func TestNewPlayer(t *testing.T) {
	p := NewPlayer()
	if p == nil {
		t.Errorf("expected player to not be nil")
	}
	if p.playlist == nil {
		t.Errorf("expected playlist to not be nil")
	}
	if p.pauseCond == nil {
		t.Errorf("expected pause condition to not be nil")
	}
}

func TestAddSong(t *testing.T) {
	p := NewPlayer()

	p.AddSong("song1.mp3")

	if p.playlist.Len() != 1 {
		t.Errorf("Expected playlist length 1, got %d", p.playlist.Len())
	}

	p.AddSong("invalid.mp4")

	if p.playlist.Len() != 1 {
		t.Errorf("Expected playlist length 1, got %d", p.playlist.Len())
	}
}

func TestAddSong2(t *testing.T) {
	p := NewPlayer()
	p.AddSong("song1.mp3")
	if p.playlist.Len() != 1 {
		t.Errorf("expected playlist length to be 1")
	}
	p.AddSong("song2.mp3")
	p.AddSong("song3.mp3")
	if p.playlist.Len() != 3 {
		t.Errorf("expected playlist length to be 3")
	}
	p.AddSong("notasong.txt")
	if p.playlist.Len() != 3 {
		t.Errorf("expected playlist length to be 3")
	}
}

func TestDelete(t *testing.T) {
	p := NewPlayer()

	p.AddSong("song1.mp3")
	time.Sleep(1 * time.Second)
	
	p.Delete()


}

func TestPlay(t *testing.T) {
	p := NewPlayer()


	p.AddSong("song1.mp3")


	err := p.Play()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}


	time.Sleep(5 * time.Second)
}

func TestPauseResume(t *testing.T) {
	p := NewPlayer()

	p.AddSong("song1.mp3")

	err := p.Play()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}


	p.Pause()

	time.Sleep(2 * time.Second)
	
	p.Resume()

	time.Sleep(5 * time.Second)
}

func TestNextPrev(t *testing.T) {
	p := NewPlayer()

	p.AddSong("song1.mp3")
	p.AddSong("song2.mp3")
	p.AddSong("song3.mp3")

	err := p.Play()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	p.Next()

	time.Sleep(5 * time.Second)

	p.Prev()

	time.Sleep(5 * time.Second)
}

func TestPlayPlaylist(t *testing.T) {
	p := NewPlayer()

	p.AddSong("song1.mp3")
	p.AddSong("song2.mp3")
	p.AddSong("song3.mp3")

	go p.PlayPlaylist()

	time.Sleep(20 * time.Second)
}


func TestPrev(t *testing.T) {
	p := NewPlayer()
	p.AddSong("song1.mp3")
	p.AddSong("song2.mp3")
	p.AddSong("song3.mp3")
	p.Next()
	p.Prev()
	if p.current != nil {
		t.Errorf("expected current song to be the second song in the playlist")
	}
	p.Prev()
	p.Prev()
	p.Prev()
	if p.current != nil {
		t.Errorf("expected current song to be the last song in the playlist")
	}
}

func TestMain(m *testing.M) {

    code := m.Run()

    os.Exit(code)
}