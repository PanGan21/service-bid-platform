export interface Bid {
  Id: number;
  Amount: number;
  CreatorId: string;
  AuctionId: string;
}

export interface NewBid {
  Amount: Bid["Amount"];
  AuctionId: Bid["AuctionId"];
}
