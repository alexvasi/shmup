package main

import (
	"fmt"
	"path/filepath"

	"github.com/Jragonmiris/go-al/al"
	"github.com/Jragonmiris/go-al/alc"
	"github.com/Jragonmiris/go-al/decoder/wav"
)

var soundData struct {
	device  alc.Device
	context alc.Context

	storage []al.Data
	buffers map[string]al.Buffer
	sources []al.Source
}

func InitSound() {
	device, err := alc.OpenDefaultDevice()
	PanicOnError(err)

	context, err := alc.CreateDefaultContext(device)
	PanicOnError(err)

	soundData.device = device
	soundData.context = context
	soundData.buffers = make(map[string]al.Buffer)

	PanicOnError(context.MakeCurrent())
	PanicOnError(context.Process())
}

func TerminateSound() {
	for _, source := range soundData.sources {
		source.Delete()
	}
	soundData.sources = nil

	for _, buffer := range soundData.buffers {
		buffer.Delete()
	}
	soundData.buffers = nil

	soundData.storage = nil
	soundData.context.Destroy()
	soundData.device.Close()
}

func LoadSoundFile(label, path string) {
	rawSound, err := wav.LoadWavFile(path)
	PanicOnError(err)

	sound, err := wav.ToALData(rawSound)
	PanicOnError(err)
	soundData.storage = append(soundData.storage, sound)

	buffer, err := al.GenBuffer()
	PanicOnError(err)

	soundData.buffers[label] = buffer
	PanicOnError(buffer.BufferData(sound))
}

func LoadSoundAssets(globPattern string) {
	files, err := filepath.Glob(globPattern)
	PanicOnError(err)

	for _, path := range files {
		_, name := filepath.Split(path)
		ext := filepath.Ext(name)
		LoadSoundFile(name[:len(name)-len(ext)], path)
	}
}

func PlaySound(label string, gain, pitch float32) {
	buffer, found := soundData.buffers[label]
	if !found {
		panic(fmt.Sprint("Can't find sound buffer with label:", label))
	}

	var source al.Source
	found = false
	for _, s := range soundData.sources {
		state, err := s.GetState()
		PanicOnError(err)
		if state == al.Stopped {
			source = s
			found = true
			break
		}
	}
	if !found {
		var err error
		source, err = al.GenSource()
		PanicOnError(err)
		soundData.sources = append(soundData.sources, source)
	}

	PanicOnError(source.SetBuffer(buffer))
	source.SetGain(gain)
	source.SetPitch(pitch)
	PanicOnError(source.Play())
}
