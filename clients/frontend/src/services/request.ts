import axios from "axios";
import { API_URL } from "../constants";
import { FormattedRequest, NewRequest, Request } from "../types/request";

export const createRequest = async (newRequest: NewRequest) => {
  return axios.post(API_URL + "/request/", newRequest, {
    withCredentials: true,
  });
};

export const rejectRequest = async (
  requestId: string,
  rejectionReason: string
) => {
  return axios.post(
    API_URL + `/request/reject?requestId=${requestId}`,
    {
      rejectionReason,
    },
    { withCredentials: true }
  );
};

export const approveRequest = async (requestId: string) => {
  return axios.post(
    API_URL + `/request/approve?requestId=${requestId}`,
    {},
    { withCredentials: true }
  );
};

export const getMyRequests = async (limit: number, page: number) => {
  return axios.get(
    API_URL +
      `/request/status/own?status=new&limit=${limit}&page=${page}&asc=false`,
    {
      withCredentials: true,
    }
  );
};

export const countMyRequests = async () => {
  return axios.get(API_URL + "/request/status/own/count?status=new", {
    withCredentials: true,
  });
};

export const getOwnRejectedRequests = async (limit: number, page: number) => {
  return axios.get(
    API_URL +
      `/request/status/own?status=rejected&limit=${limit}&page=${page}&asc=false`,
    { withCredentials: true }
  );
};

export const countRequestsByStatus = async (status: string) => {
  return axios.get(API_URL + `/request/status/count?status=${status}`, {
    withCredentials: true,
  });
};

export const getRequestsByStatus = async (
  status: string,
  limit: number,
  page: number
) => {
  return axios.get(
    API_URL +
      `/request/status?status=${status}&limit=${limit}&page=${page}&asc=false`,
    { withCredentials: true }
  );
};

export const countOwnRejectedRequests = async () => {
  return axios.get(API_URL + "/request/status/own/count?status=rejected", {
    withCredentials: true,
  });
};

export const formatRequests = (requests: Request[]): FormattedRequest[] => {
  return requests.map((r) => {
    return {
      ...r,
      Deadline: new Date(r.Deadline).toLocaleDateString("en-US"),
    };
  });
};
