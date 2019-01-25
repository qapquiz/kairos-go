package client

// Hub maintains the set of active remotes and broadcasts messages to the
// remotes
type Hub struct {
	// Registered remotes.
	remotes map[*Remote]bool

	// Inbound messages from the remotes.
	broadcast chan []byte

	// Register requests from the remotes.
	register chan *Remote

	// Unregister requests from remotes
	unregister chan *Remote
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Remote),
		unregister: make(chan *Remote),
		remotes:    make(map[*Remote]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case remote := <-h.register:
			h.remotes[remote] = true
		case remote := <-h.unregister:
			if _, ok := h.remotes[remote]; ok {
				delete(h.remotes, remote)
				close(remote.send)
			}
		case message := <-h.broadcast:
			for remote := range h.remotes {
				select {
				case remote.send <- message:
				default:
					close(remote.send)
					delete(h.remotes, remote)
				}
			}
		}
	}
}
