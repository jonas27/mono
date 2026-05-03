//go:build !cgo && !android && !darwin && !js && !windows && !nintendosdk

package output

import "errors"

type alsaOutput struct{}

func newAlsaOutput(_ *NewOutputOptions) (*alsaOutput, error) {
	return nil, errors.New("alsa output not available (built without cgo)")
}

func (out *alsaOutput) Pause() error             { return nil }
func (out *alsaOutput) Resume() error            { return nil }
func (out *alsaOutput) Drop() error              { return nil }
func (out *alsaOutput) DelayMs() (int64, error)  { return 0, nil }
func (out *alsaOutput) SetVolume(_ float32)      {}
func (out *alsaOutput) Error() <-chan error       { return nil }
func (out *alsaOutput) Close() error             { return nil }
func (out *alsaOutput) setupMixer() error        { return nil }
func (out *alsaOutput) waitForMixerEvents()      {}
