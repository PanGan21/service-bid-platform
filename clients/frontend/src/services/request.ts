import axios from "axios";
import { API_URL } from "../constants";
import {
  ExtendedFormattedRequest,
  ExtendedRequest,
  FormattedRequest,
  NewRequest,
  Request,
} from "../types/request";

export const getAllRequests = async (limit: number, page: number) => {
  return axios.get(API_URL + `/request/?limit=${limit}&page=${page}&asc=true`, {
    withCredentials: true,
  });
};

export const countAllRequests = async () => {
  return axios.get(API_URL + "/request/count", { withCredentials: true });
};

export const getMyRequests = async (limit: number, page: number) => {
  return axios.get(
    API_URL + `/request/own?limit=${limit}&page=${page}&asc=true`,
    {
      withCredentials: true,
    }
  );
};

export const countMyRequests = async () => {
  return axios.get(API_URL + "/request/count/own", { withCredentials: true });
};

export const getOpenPastDeadlineRequests = async (
  limit: number,
  page: number
) => {
  return axios.get(
    API_URL +
      `/request/open/past-deadline?limit=${limit}&page=${page}&asc=true`,
    {
      withCredentials: true,
    }
  );
};

export const countOpenPastDeadlineRequests = async () => {
  return axios.get(API_URL + `/request/open/past-deadline/count`, {
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
      `/request/status?status=${status}&limit=${limit}&page=${page}&asc=true`,
    {
      withCredentials: true,
    }
  );
};

export const countRequestsByStatus = async (status: string) => {
  return axios.get(API_URL + `/request/status/count?status=${status}`, {
    withCredentials: true,
  });
};

export const updateWinner = async (requestId: string) => {
  return axios.post(
    API_URL + `/request/update/winner?requestId=${requestId}`,
    {},
    { withCredentials: true }
  );
};

export const updateRequestStatus = async (
  requestId: string,
  status: string
) => {
  return axios.post(
    API_URL + `/request/update/status?requestId=${requestId}`,
    {
      status,
    },
    { withCredentials: true }
  );
};

export const createRequest = async (newRequest: NewRequest) => {
  return axios.post(API_URL + "/request/", newRequest, {
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

export const formatExtendedRequests = (
  requests: ExtendedRequest[]
): ExtendedFormattedRequest[] => {
  return requests.map((r) => {
    return {
      ...r,
      Deadline: new Date(r.Deadline).toLocaleDateString("en-US"),
    };
  });
};
