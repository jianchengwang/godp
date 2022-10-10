import { App, defineComponent } from "vue";
import projectHostTree from "./src/HostTree.vue";

export const ProjectHostTree = Object.assign(projectHostTree, {
  install(app: App) {
    app.component(projectHostTree.name, ProjectHostTree);
  }
});

export default {
  ProjectHostTree
};