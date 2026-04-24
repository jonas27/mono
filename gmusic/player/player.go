package player

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
	"github.com/gopxl/beep/wav"
)

type State int

const (
	StateStopped State = iota
	StatePlaying
	StatePaused
)

type Player struct {
	mu       sync.Mutex
	state    State
	ctrl     *beep.Ctrl
	vol      *effects.Volume
	streamer beep.StreamSeekCloser
	format   beep.Format
	volume   int // 0-100
	onDone   func()
	done     chan struct{}
}

func New() *Player {
	return &Player{
		volume: 100,
		state:  StateStopped,
	}
}

func (p *Player) Load(path string, onDone func()) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Stop existing playback
	if p.ctrl != nil {
		p.ctrl.Paused = true
	}
	if p.streamer != nil {
		p.streamer.Close()
		p.streamer = nil
	}
	if p.done != nil {
		close(p.done)
		p.done = nil
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	ext := filepath.Ext(path)
	var streamer beep.StreamSeekCloser
	var format beep.Format

	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".flac":
		streamer, format, err = flac.Decode(f)
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
	default:
		f.Close()
		return fmt.Errorf("unsupported format: %s", ext)
	}

	if err != nil {
		f.Close()
		return fmt.Errorf("decode: %w", err)
	}

	// Initialize speaker if needed
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	p.streamer = streamer
	p.format = format
	p.onDone = onDone
	p.done = make(chan struct{})
	p.state = StatePlaying

	resampled := beep.Resample(4, format.SampleRate, format.SampleRate, streamer)

	p.ctrl = &beep.Ctrl{Streamer: resampled, Paused: false}
	p.vol = &effects.Volume{
		Streamer: p.ctrl,
		Base:     2,
		Volume:   p.beepVolume(),
		Silent:   p.volume == 0,
	}

	done := p.done
	onDoneCb := onDone
	speaker.Play(beep.Seq(p.vol, beep.Callback(func() {
		select {
		case <-done:
		default:
			p.mu.Lock()
			p.state = StateStopped
			p.mu.Unlock()
			if onDoneCb != nil {
				onDoneCb()
			}
		}
	})))

	return nil
}

func (p *Player) PlayPause() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.ctrl == nil {
		return
	}

	speaker.Lock()
	if p.ctrl.Paused {
		p.ctrl.Paused = false
		p.state = StatePlaying
	} else {
		p.ctrl.Paused = true
		p.state = StatePaused
	}
	speaker.Unlock()
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.ctrl != nil {
		speaker.Lock()
		p.ctrl.Paused = true
		speaker.Unlock()
	}
	if p.streamer != nil {
		p.streamer.Close()
		p.streamer = nil
	}
	if p.done != nil {
		close(p.done)
		p.done = nil
	}
	p.ctrl = nil
	p.vol = nil
	p.state = StateStopped
}

func (p *Player) SetVolume(vol int) {
	if vol < 0 {
		vol = 0
	}
	if vol > 100 {
		vol = 100
	}

	p.mu.Lock()
	p.volume = vol
	bv := p.beepVolume()
	silent := vol == 0

	if p.vol != nil {
		speaker.Lock()
		p.vol.Volume = bv
		p.vol.Silent = silent
		speaker.Unlock()
	}
	p.mu.Unlock()
}

func (p *Player) beepVolume() float64 {
	return (float64(p.volume)/100.0)*10.0 - 10.0
}

func (p *Player) State() State {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.state
}

func (p *Player) Volume() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.volume
}

func (p *Player) Progress() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.streamer == nil {
		return 0
	}

	total := p.streamer.Len()
	if total == 0 {
		return 0
	}
	pos := p.streamer.Position()
	return float64(pos) / float64(total)
}

func (p *Player) Position() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.streamer == nil {
		return 0
	}
	return p.format.SampleRate.D(p.streamer.Position())
}

func (p *Player) Duration() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.streamer == nil {
		return 0
	}
	return p.format.SampleRate.D(p.streamer.Len())
}

func (p *Player) Seek(offset time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.streamer == nil {
		return
	}

	samples := p.format.SampleRate.N(offset)
	pos := p.streamer.Position() + samples
	if pos < 0 {
		pos = 0
	}
	if max := p.streamer.Len(); pos > max {
		pos = max
	}

	speaker.Lock()
	p.streamer.Seek(pos)
	speaker.Unlock()
}
