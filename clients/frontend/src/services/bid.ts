import axios from "axios";
import { API_URL } from "../constants";
import { NewBid } from "../types/bid";

export const createBid = async (newBid: NewBid) => {
  return axios.post(API_URL + "/bidding/", newBid, { withCredentials: true });
};
