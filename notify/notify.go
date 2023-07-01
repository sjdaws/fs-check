package notify

import (
	"log"
	"strings"

	"github.com/containrrr/shoutrrr"
)

type Notify struct {
	urls []string
}

func New(urls []string) *Notify {
	return &Notify{
		urls: urls,
	}
}

func (n *Notify) Message(text string) {
	for _, url := range n.urls {
		url = strings.TrimSpace(url)

		if url == "" {
			continue
		}

		err := shoutrrr.Send(url, text)
		if err != nil {
			log.Printf("notify: %v", err)
		}
	}
}
