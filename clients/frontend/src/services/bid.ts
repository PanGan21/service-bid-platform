import axios from "axios";
import { API_URL } from "../constants";
import { NewBid } from "../types/bid";

export const createBid = async (newBid: NewBid) => {
  return axios.post(API_URL + "/bidding/", newBid, { withCredentials: true });
};

export const getMyBids = async (limit: number, page: number) => {
  return axios.get(
    API_URL + `/bidding/own?limit=${limit}&page=${page}&asc=true`,
    {
      withCredentials: true,
    }
  );
};

export const countMyBids = async () => {
  return axios.get(API_URL + "/bidding/count/own", { withCredentials: true });
};
