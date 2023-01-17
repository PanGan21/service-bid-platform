import axios from "axios";
import { API_URL } from "../constants";
import { FormattedRequest, NewRequest, Request } from "../types/request";

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

export const createRequest = async (newRequest: NewRequest) => {
  return axios.post(API_URL + "/request/", newRequest, {
    withCredentials: true,
  });
};

export const formatRequests = (requests: Request[]): FormattedRequest[] => {
  return requests.map((r) => {
    return {
      ...r,
      Deadline: new Date().toLocaleDateString("en-US"),
    };
  });
};
