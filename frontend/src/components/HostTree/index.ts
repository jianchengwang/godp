import { App, defineComponent } from "vue";
import hostTree from "./src/HostTree.vue";

export const HostTree = Object.assign(hostTree, {
  install(app: App) {
    app.component(hostTree.name, HostTree);
  }
});

export default {
  HostTree
};
