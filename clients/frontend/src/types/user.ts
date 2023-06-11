export interface User {
  Id: string;
  Username: string;
  Email: string;
  Phone: string;
  Password?: string;
  Roles: Array<string>;
}

export interface UserDetails {
  Id: string;
  Username: string;
  Email: string;
  Phone: string;
  Roles: Array<string>;
}
