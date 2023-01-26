package bid

import "github.com/PanGan21/pkg/entity"

type BidEvents interface {
	PublishBidCreated(topic string, bid *entity.Bid) error
}
