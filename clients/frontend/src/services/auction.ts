import axios from "axios";
import { API_URL } from "../constants";
import {
  ExtendedFormattedAuction,
  ExtendedAuction,
  FormattedAuction,
  NewAuction,
  Auction,
} from "../types/auction";

export const getMyAuctions = async (limit: number, page: number) => {
  return axios.get(
    API_URL + `/auction/own?limit=${limit}&page=${page}&asc=false`,
    {
      withCredentials: true,
    }
  );
};

export const countMyAuctions = async () => {
  return axios.get(API_URL + "/auction/count/own", { withCredentials: true });
};

export const getOpenPastDeadlineAuctions = async (
  limit: number,
  page: number
) => {
  return axios.get(
    API_URL +
      `/auction/open/past-deadline?limit=${limit}&page=${page}&asc=false`,
    {
      withCredentials: true,
    }
  );
};

export const countOpenPastDeadlineAuctions = async () => {
  return axios.get(API_URL + `/auction/open/past-deadline/count`, {
    withCredentials: true,
  });
};

export const getAuctionsByStatus = async (
  status: string,
  limit: number,
  page: number
) => {
  return axios.get(
    API_URL +
      `/auction/status?status=${status}&limit=${limit}&page=${page}&asc=false`,
    {
      withCredentials: true,
    }
  );
};

export const countAuctionsByStatus = async (status: string) => {
  return axios.get(API_URL + `/auction/status/count?status=${status}`, {
    withCredentials: true,
  });
};

export const updateWinner = async (auctionId: string) => {
  return axios.post(
    API_URL + `/auction/update/winner?auctionId=${auctionId}`,
    {},
    { withCredentials: true }
  );
};

export const updateAuctionStatus = async (
  auctionId: string,
  status: string
) => {
  return axios.post(
    API_URL + `/auction/update/status?auctionId=${auctionId}`,
    {
      status,
    },
    { withCredentials: true }
  );
};

export const createAuction = async (newAuction: NewAuction) => {
  return axios.post(API_URL + "/auction/", newAuction, {
    withCredentials: true,
  });
};

export const getOwnAssignments = async (limit: number, page: number) => {
  return axios.get(
    API_URL +
      `/auction/own/assigned-bids?limit=${limit}&page=${page}&asc=false`,
    { withCredentials: true }
  );
};

export const countOwnAssignments = async () => {
  return axios.get(API_URL + "/auction/own/assigned-bids/count", {
    withCredentials: true,
  });
};

export const rejectAuction = async (
  auctionId: string,
  rejectionReason: string
) => {
  return axios.post(
    API_URL + `/auction/update/reject?auctionId=${auctionId}`,
    {
      rejectionReason,
    },
    { withCredentials: true }
  );
};

export const formatAuctions = (auctions: Auction[]): FormattedAuction[] => {
  return auctions.map((r) => {
    return {
      ...r,
      Deadline: new Date(r.Deadline).toLocaleDateString("en-US"),
    };
  });
};

export const formatExtendedAuctions = (
  auctions: ExtendedAuction[]
): ExtendedFormattedAuction[] => {
  return auctions.map((r) => {
    return {
      ...r,
      Deadline: new Date(r.Deadline).toLocaleDateString("en-US"),
    };
  });
};
