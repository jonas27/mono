package vorbis

import (
	librespot "github.com/devgianlu/go-librespot"
	"github.com/jfreymuth/oggvorbis"
)

// Decoder implements an OggVorbis decoder using a pure-Go library.
type Decoder struct {
	log librespot.Logger

	SampleRate int32
	Channels   int32

	gain   float32
	reader *oggvorbis.Reader
}

func New(log librespot.Logger, r librespot.SizedReadAtSeeker, _ *MetadataPage, gain float32) (*Decoder, error) {
	reader, err := oggvorbis.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &Decoder{
		log:        log,
		SampleRate: int32(reader.SampleRate()),
		Channels:   int32(reader.Channels()),
		gain:       gain,
		reader:     reader,
	}, nil
}

func (d *Decoder) Read(p []float32) (int, error) {
	n, err := d.reader.Read(p)
	if d.gain != 1.0 {
		for i := 0; i < n; i++ {
			p[i] *= d.gain
		}
	}
	return n, err
}

func (d *Decoder) SetPositionMs(pos int64) error {
	return d.reader.SetPosition(pos * int64(d.SampleRate) / 1000)
}

func (d *Decoder) PositionMs() int64 {
	return d.reader.Position() * 1000 / int64(d.SampleRate)
}

func (d *Decoder) Close() {}
