import { http } from "../utils/http";

const BASE_API = "assetsHost";

export const assetsHostPage = (params?: object) => {
  return http.request("get", `${BASE_API}/page`, { params });
};

export const assetsHostGet = (id: number) => {
  return http.get(`${BASE_API}/${id}`);
};

export const assetsHostPost = (params?: object) => {
  return http.request("post", `${BASE_API}`, { data: params });
};

export const assetsHostPut = (id: number, params?: object) => {
  return http.request("put", `${BASE_API}/${id}`, { data: params });
};

export const assetsHostDelete = (id: number) => {
  return http.request("delete", `${BASE_API}/${id}`);
};