import { App, defineComponent } from "vue";
import webSftp from "./src/WebSftp.vue";

export const WebSftp = Object.assign(webSftp, {
  install(app: App) {
    app.component(webSftp.name, WebSftp);
  }
});

export default {
  WebSftp
};