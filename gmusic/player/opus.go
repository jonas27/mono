package player

import (
	"errors"
	"io"
	"os"

	"github.com/gopxl/beep"
	"github.com/pion/opus"
	"github.com/pion/opus/pkg/oggreader"
)

const (
	opusSampleRate       = 48000
	opusMaxSamplesPerPkt = 5760 // 120 ms at 48 kHz per channel
)

type opusStreamer struct {
	f        *os.File
	samples  []float32 // interleaved PCM
	channels int
	pos      int // frame position (per channel)
}

func decodeOpus(f *os.File) (beep.StreamSeekCloser, beep.Format, error) {
	r, header, err := oggreader.NewWith(f)
	if err != nil {
		return nil, beep.Format{}, err
	}

	channels := int(header.Channels)
	if channels < 1 {
		channels = 1
	}
	if channels > 2 {
		channels = 2
	}

	dec, err := opus.NewDecoderWithOutput(opusSampleRate, channels)
	if err != nil {
		return nil, beep.Format{}, err
	}

	// Second Ogg packet is the comment header (OpusTags); skip it.
	if _, _, err = r.ParseNextPacket(); err != nil {
		return nil, beep.Format{}, err
	}

	outBuf := make([]float32, opusMaxSamplesPerPkt*channels)
	var samples []float32

	for {
		pkt, _, err := r.ParseNextPacket()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			break
		}
		n, err := dec.DecodeToFloat32(pkt, outBuf)
		if err != nil || n == 0 {
			continue
		}
		samples = append(samples, outBuf[:n*channels]...)
	}

	// Discard encoder pre-skip samples.
	skip := int(header.PreSkip) * channels
	if skip > len(samples) {
		skip = len(samples)
	}
	samples = samples[skip:]

	format := beep.Format{
		SampleRate:  beep.SampleRate(opusSampleRate),
		NumChannels: channels,
		Precision:   4,
	}

	return &opusStreamer{f: f, samples: samples, channels: channels}, format, nil
}

func (s *opusStreamer) Stream(out [][2]float64) (int, bool) {
	total := len(s.samples) / s.channels
	if s.pos >= total {
		return 0, false
	}
	for i := range out {
		if s.pos >= total {
			return i, i > 0
		}
		base := s.pos * s.channels
		if s.channels == 1 {
			v := float64(s.samples[base])
			out[i] = [2]float64{v, v}
		} else {
			out[i] = [2]float64{float64(s.samples[base]), float64(s.samples[base+1])}
		}
		s.pos++
	}
	return len(out), true
}

func (s *opusStreamer) Err() error { return nil }

func (s *opusStreamer) Len() int { return len(s.samples) / s.channels }

func (s *opusStreamer) Position() int { return s.pos }

func (s *opusStreamer) Seek(p int) error {
	if p < 0 || p > s.Len() {
		return errors.New("opus: seek position out of range")
	}
	s.pos = p
	return nil
}

func (s *opusStreamer) Close() error { return s.f.Close() }
