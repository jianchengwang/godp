import axios from "axios";
import qs from "qs";
import { loadEnv } from "@build/index";
import { getAccessToken } from "/@/utils/auth";

const { VITE_PROXY_DOMAIN, VITE_PROXY_DOMAIN_REAL } = loadEnv();
const baseURL =
  process.env.NODE_ENV === "production"
    ? VITE_PROXY_DOMAIN_REAL
    : VITE_PROXY_DOMAIN;
console.info(baseURL);

const getBlob = async (fetchUrl, param, filename) => {
  const accessToken = getAccessToken();
  const { data, headers } = await axios({
    method: "get",
    responseType: "blob",
    url: baseURL + fetchUrl,
    headers: {
      Authorization: "Bearer " + accessToken
    },
    params: param,
    paramsSerializer: params => qs.stringify(params, { indices: false })
  });
  if (!filename) {
    filename = headers["filename"];
  }
  const blob = new Blob([data]);
  const reader = new FileReader();
  reader.readAsDataURL(blob);
  reader.onload = e => {
    const a = document.createElement("a");
    a.download = filename;
    a.href = e.target.result;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
  };
};

const uploadBlob = async (fetchUrl, formData: FormData) => {
  const accessToken = getAccessToken();
  return await axios({
    method: "post",
    responseType: "json",
    url: baseURL + fetchUrl,
    headers: {
      Authorization: "Bearer " + accessToken,
      "Content-Type": "multipart/form-data"
    },
    data: formData
  });
};

export { getBlob, uploadBlob };
