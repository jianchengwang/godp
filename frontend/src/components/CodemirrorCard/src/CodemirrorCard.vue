<script lang="ts">
export default {
  name: "CodemirrorCard"
};
</script>

<script lang="ts" setup>
import { ref, reactive, computed, onMounted } from "vue";

const props = defineProps({
  filename: String,
  mode: String,
  updateText: String,
  remark: String
});

const emit = defineEmits(["save-event", "cancel-event"]);

import Codemirror from "codemirror-editor-vue3";
import "codemirror/mode/nginx/nginx.js";
import "codemirror/mode/shell/shell.js";
import "codemirror/mode/yaml/yaml.js";
import "codemirror/mode/toml/toml.js";
import "codemirror/mode/python/python.js";
import "codemirror/mode/go/go.js";
import "codemirror/mode/javascript/javascript.js";
import "codemirror/mode/css/css.js";
import "codemirror/mode/sql/sql.js";
import "codemirror/mode/markdown/markdown.js";
import "codemirror/theme/dracula.css";
const cmOptions = reactive({
  mode: "text/javascript",
  theme: "dracula", // Theme
  lineNumbers: true, // Show line number
  smartIndent: true, // Smart indent
  indentUnit: 2, // The smart indent unit is 2 spaces in length
  foldGutter: true, // Code folding
  styleActiveLine: true // Display the style of the selected row
});

const getModeByFileName = fileName => {
  if (props.mode) {
    return props.mode;
  }
  if (fileName.endsWith("conf")) {
    return "nginx";
  } else if (fileName.endsWith("sh")) {
    return "shell";
  } else if (fileName.endsWith(".yaml") || fileName.endsWith(".yml")) {
    return "yaml";
  } else if (fileName.endsWith(".toml")) {
    return "toml";
  } else if (fileName.endsWith(".py")) {
    return "python";
  } else if (fileName.endsWith(".go")) {
    return "go";
  } else if (fileName.endsWith(".js")) {
    return "javascript";
  } else if (fileName.endsWith(".css")) {
    return "css";
  } else if (fileName.endsWith(".sql")) {
    return "sql";
  } else if (fileName.endsWith(".md")) {
    return "markdown";
  }
  return "javascript";
};

const updateText = ref("");
onMounted(() => {
  cmOptions.mode = getModeByFileName(props.filename);
  updateText.value = props.updateText;
});
</script>

<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <div>
          {{ props.filename }}
        </div>
        <div class="card-header">
          <el-alert
            style="font-size: 12px; margin-left: 10px; margin-right: 10px"
            :title="props.remark"
            type="warning"
            :closable="false"
          />
          <el-button
            type="primary"
            @click="emit('save-event', filename, updateText)"
            >保存</el-button
          >
          <el-button @click="emit('cancel-event')">返回</el-button>
        </div>
      </div>
    </template>
    <Codemirror
      v-model:value="updateText"
      :options="cmOptions"
      border
      placeholder=""
      :height="750"
    />
  </el-card>
</template>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
