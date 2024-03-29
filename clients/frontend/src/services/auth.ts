import axios from "axios";
import { API_URL } from "../constants";
import { User } from "../types/user";

export const login = async (username: string, password: string) => {
  return axios.post(
    API_URL + "/user/login",
    {
      username,
      password,
    },
    { withCredentials: true }
  );
};

export const getLoggedInUserDetails = async () => {
  return axios.get(API_URL + "/user/", { withCredentials: true });
};

export const logout = async () => {
  localStorage.removeItem("user");
  clearCookies();
  await axios.post(API_URL + "/user/logout", {}, { withCredentials: true });
};

export const getCurrentUser = (): User | undefined => {
  const userStr = localStorage.getItem("user");
  if (userStr) {
    return JSON.parse(userStr);
  }
};

export const getUserDetailsById = async (id: string) => {
  return axios.get(API_URL + "/user/details?userId=" + id, {
    withCredentials: true,
  });
};

const clearCookies = () => {
  document.cookie.split(";").forEach((c) => {
    document.cookie = c
      .replace(/^ +/, "")
      .replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/");
  });
};
