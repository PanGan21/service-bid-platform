export interface Request {
  Id: string;
  Title: string;
  CreatorId: string;
  Postcode: string;
  Info: string;
  Status: string;
  RejectionReason: string;
}

export interface NewRequest {
  Title: Request["Title"];
  Postcode: Request["Postcode"];
  Info: Request["Info"];
}
