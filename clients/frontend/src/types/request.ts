import { Bid } from "./bid";

export interface Request {
  Id: string;
  Title: string;
  CreatorId: string;
  Postcode: string;
  Info: string;
  Deadline: number;
  Status: string;
}

export interface ExtendedRequest extends Request {
  BidsCount: number;
}

export interface FormattedRequest {
  Id: Request["Id"];
  Title: Request["Title"];
  CreatorId: Request["CreatorId"];
  Postcode: Request["Postcode"];
  Info: Request["Info"];
  Deadline: string;
  Status: Request["Status"];
}

export interface ExtendedFormattedRequest extends FormattedRequest {
  BidsCount: number;
}

export interface Assignment {
  Id: Request["Id"];
  Title: Request["Title"];
  CreatorId: Request["CreatorId"];
  Postcode: Request["Postcode"];
  Info: Request["Info"];
  Status: Request["Status"];
  BidId: Bid["Id"];
  BidAmount: Bid["Amount"];
}

export interface NewRequest {
  Title: Request["Title"];
  Postcode: Request["Postcode"];
  Info: Request["Info"];
  Deadline: Request["Deadline"];
}
