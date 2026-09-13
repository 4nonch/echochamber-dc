// Package actions provides helper functions to use discord / discordgo API
package actions

import (
	"io"
	"log"

	"github.com/bwmarrin/discordgo"
)

type StreamFiles struct {
	Files []*discordgo.File
	Resps []io.ReadCloser
}

func (f *StreamFiles) Close() {
	for _, file := range f.Resps {
		err := file.Close()
		if err != nil {
			log.Printf("Warning: Failed to close file. Error: %v", err)
		}
	}
}
