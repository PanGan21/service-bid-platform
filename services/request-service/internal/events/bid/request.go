package request

import (
	"github.com/PanGan21/pkg/messaging"
)

type requestEvents struct {
	pub messaging.Publisher
}

const ()

func NewRequestEvents(pub messaging.Publisher) *requestEvents {
	return &requestEvents{pub: pub}
}
