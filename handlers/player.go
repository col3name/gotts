package handlers

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"time"
)

type BeepPlayer struct {
	Volume float64
	Speed  float64
	volume *effects.Volume
}

func (p *BeepPlayer) SetVolume(volume float64) {
	if volume >= 0 && volume <= 2.0 {
		p.Volume = volume
	}
}

func (p *BeepPlayer) SetSpeed(speed float64) {
	if speed >= 0.3 && speed <= 1.3 {
		p.Speed = speed
	}
}

func (p *BeepPlayer) Play(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctrl := &beep.Ctrl{Streamer: beep.Loop(1, streamer), Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	speedy := beep.ResampleRatio(4, 1, volume)
	done := make(chan bool)
	speaker.Play(beep.Seq(speedy, streamer, beep.Callback(func() {
		done <- true
	})))
	<-done

	return nil
}
