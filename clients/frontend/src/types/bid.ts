export interface Bid {
  Id: number;
  Amount: number;
  CreatorId: string;
  RequestId: string;
}

export interface NewBid {
  Amount: Bid["Amount"];
  RequestId: Bid["RequestId"];
}
