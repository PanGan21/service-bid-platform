import { Bid } from "./bid";

export interface Auction {
  Id: string;
  Title: string;
  CreatorId: string;
  Postcode: string;
  Info: string;
  Deadline: number;
  Status: string;
  RejectionReason: string;
  WinnerId: string;
  WinningAmount: string;
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
  RejectionReason: Auction["RejectionReason"];
  WinnerId: Auction["WinnerId"];
  WinningAmount: Auction["WinningAmount"];
}

export interface ExtendedFormattedAuction extends FormattedAuction {
  BidsCount: number;
}

export interface NewAuction {
  Title: Auction["Title"];
  Postcode: Auction["Postcode"];
  Info: Auction["Info"];
  Deadline: Auction["Deadline"];
}
