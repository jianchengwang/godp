import { App, defineComponent } from "vue";
import webTerminalrt from "./src/WebTerminal.vue";

export const WebTerminal = Object.assign(webTerminalrt, {
  install(app: App) {
    app.component(webTerminalrt.name, WebTerminal);
  }
});

export default {
  WebTerminal
};
