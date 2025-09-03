package handlers

import "github.com/bwmarrin/discordgo"

type handlersType []any

var handlers = handlersType{}

// Used to add event handler with signature (*discordgo.Session, some discodgo event struct)
func (h *handlersType) add(handler any) {
	*h = append(*h, handler)
}

func Register(s *discordgo.Session) {
	for _, h := range handlers {
		s.AddHandler(h)
	}
}
