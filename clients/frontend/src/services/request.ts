import axios from "axios";
import { API_URL } from "../constants";

export const getMyRequests = async () => {
  return axios.get(API_URL + "/request/own", { withCredentials: true });
};
