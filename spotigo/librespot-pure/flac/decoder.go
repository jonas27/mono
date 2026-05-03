package flac

import (
	"fmt"
	"math"

	librespot "github.com/devgianlu/go-librespot"
	"github.com/mewkiz/flac"
)

// Decoder implements an FLAC decoder using a pure-Go FLAC library.
type Decoder struct {
	log librespot.Logger

	SampleRate int32
	Channels   int32

	gain   float32
	stream *flac.Stream
	buf    []float32
	pos    int64
}

func New(log librespot.Logger, r librespot.SizedReadAtSeeker, gain float32) (*Decoder, error) {
	s, err := flac.NewSeek(r)
	if err != nil {
		return nil, fmt.Errorf("flac: open stream: %w", err)
	}

	return &Decoder{
		log:        log,
		SampleRate: int32(s.Info.SampleRate),
		Channels:   int32(s.Info.NChannels),
		gain:       gain,
		stream:     s,
	}, nil
}

func (d *Decoder) Read(p []float32) (n int, err error) {
	norm := float32(math.Pow(2, float64(d.stream.Info.BitsPerSample-1)))
	for n < len(p) {
		if len(d.buf) > 0 {
			copied := copy(p[n:], d.buf)
			d.buf = d.buf[copied:]
			n += copied
			continue
		}
		f, ferr := d.stream.ParseNext()
		if ferr != nil {
			return n, ferr
		}
		for i := 0; i < int(f.BlockSize); i++ {
			for ch := 0; ch < int(d.Channels); ch++ {
				d.buf = append(d.buf, float32(f.Subframes[ch].Samples[i])/norm*d.gain)
			}
		}
		d.pos += int64(f.BlockSize)
	}
	return n, nil
}

func (d *Decoder) SetPositionMs(pos int64) error {
	sample := uint64(pos) * uint64(d.SampleRate) / 1000
	if _, err := d.stream.Seek(sample); err != nil {
		return fmt.Errorf("flac: seek: %w", err)
	}
	d.pos = int64(sample)
	d.buf = nil
	return nil
}

func (d *Decoder) PositionMs() int64 {
	return d.pos * 1000 / int64(d.SampleRate)
}

func (d *Decoder) Close() error {
	return d.stream.Close()
}
