export interface User {
  Id: string;
  Username: string;
  Email: string;
  Phone: string;
  Password?: string;
  Roles: Array<string>;
}
