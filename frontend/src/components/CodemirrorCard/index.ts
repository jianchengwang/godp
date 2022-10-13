import { App, defineComponent } from "vue";
import codemirrorCard from "./src/CodemirrorCard.vue";

export const CodemirrorCard = Object.assign(codemirrorCard, {
  install(app: App) {
    app.component(codemirrorCard.name, CodemirrorCard);
  }
});

export default {
  CodemirrorCard
};