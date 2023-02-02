package bid

import "github.com/PanGan21/pkg/entity"

type BidEvents interface {
	PublishBidCreated(bid *entity.Bid) error
}
