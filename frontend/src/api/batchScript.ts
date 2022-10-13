import { http } from "../utils/http";

const BASE_API = "batchScript";

export const batchScriptPage = (params?: object) => {
  return http.request("get", `${BASE_API}/page`, { params });
};

export const batchScriptGet = (id: number) => {
  return http.get(`${BASE_API}/${id}`);
};

export const batchScriptPost = (params?: object) => {
  return http.request("post", `${BASE_API}`, { data: params });
};

export const batchScriptPut = (id: number, params?: object) => {
  return http.request("put", `${BASE_API}/${id}`, { data: params });
};

export const batchScriptDelete = (id: number) => {
  return http.request("delete", `${BASE_API}/${id}`);
};
