export interface Request {
  Id: string;
  Title: string;
  CreatorId: string;
  Postcode: string;
  Info: string;
  Deadline: number;
  Status: string;
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

export interface NewRequest {
  Title: Request["Title"];
  Postcode: Request["Postcode"];
  Info: Request["Info"];
  Deadline: Request["Deadline"];
}
