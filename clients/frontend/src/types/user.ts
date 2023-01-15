export interface User {
  Id: string;
  Username: string;
  Password?: string;
  Roles: Array<string>;
}
