import { http } from "../utils/http";
import { getBlob, uploadBlob } from "../utils/http/getBlob";

const BASE_API = "/ws/ssh";

export const wsSshDownloadFile = (
  sessionId: string,
  params: object,
  filename: string
) => {
  return getBlob(`${BASE_API}/${sessionId}/downloadFile`, params, filename);
};

export const wsSshUploadFile = (sessionId: string, formData: FormData) => {
  return uploadBlob(`${BASE_API}/${sessionId}/uploadFile`, formData);
};
