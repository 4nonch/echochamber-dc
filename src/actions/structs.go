package actions

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

type StreamFiles struct {
	Files []*discordgo.File
	Resps []io.ReadCloser
}

func (f *StreamFiles) Close() {
	for _, file := range f.Resps {
		file.Close()
	}
}
