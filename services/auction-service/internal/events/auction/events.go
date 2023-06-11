package auction

import "github.com/PanGan21/pkg/entity"

type AuctionEvents interface {
	PublishAuctionUpdated(auction *entity.Auction) error
}
