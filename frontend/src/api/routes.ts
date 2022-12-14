import { http } from "../utils/http";

type Result = {
  code: number;
  data: Array<any>;
};

export const getAsyncRoutes = (params?: object) => {
  return http.request<Result>("get", "/getAsyncRoutes", { params });
};
