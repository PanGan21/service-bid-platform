import { Bid } from "./bid";

export interface Auction {
  Id: string;
  Title: string;
  CreatorId: string;
  Postcode: string;
  Info: string;
  Deadline: number;
  Status: string;
}

export interface ExtendedAuction extends Auction {
  BidsCount: number;
}

export interface FormattedAuction {
  Id: Auction["Id"];
  Title: Auction["Title"];
  CreatorId: Auction["CreatorId"];
  Postcode: Auction["Postcode"];
  Info: Auction["Info"];
  Deadline: string;
  Status: Auction["Status"];
}

export interface ExtendedFormattedAuction extends FormattedAuction {
  BidsCount: number;
}

export interface Assignment {
  Id: Auction["Id"];
  Title: Auction["Title"];
  CreatorId: Auction["CreatorId"];
  Postcode: Auction["Postcode"];
  Info: Auction["Info"];
  Status: Auction["Status"];
  BidId: Bid["Id"];
  BidAmount: Bid["Amount"];
}

export interface NewAuction {
  Title: Auction["Title"];
  Postcode: Auction["Postcode"];
  Info: Auction["Info"];
  Deadline: Auction["Deadline"];
}
