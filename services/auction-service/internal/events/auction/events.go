package auction

import "github.com/PanGan21/pkg/entity"

type AuctionEvents interface {
	PublishAuctionCreated(auction *entity.Auction) error
	PublishAuctionUpdated(auction *entity.Auction) error
}
