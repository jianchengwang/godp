<script lang="ts">
export default {
  name: "WebTerminal"
};
</script>

<script lang="ts" setup>
import { ElMessage } from "element-plus";
import { loadEnv } from "@build/index";
import { ref, reactive, toRaw, computed, onMounted, onUnmounted } from "vue";

import { getAccessToken } from "/@/utils/auth";

import { Base64 } from "js-base64";
import { Dracula, Atom, Github, Material } from "xterm-theme";
import "xterm/css/xterm.css";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { AttachAddon } from "xterm-addon-attach";
import { WebLinksAddon } from "xterm-addon-web-links";

const { VITE_PROXY_DOMAIN_REAL } = loadEnv();

const props = defineProps({
  host: String,
  sessionId: String,
  cardLayout: Boolean
});
const emit = defineEmits(["clear-event"]);

let wsUrl = computed(() => {
  return (
    `ws://localhost:8081/api/v1/ws/ssh/ssh-${props.sessionId}?host=${props.host}&Authorization=` +
    getAccessToken()
  );
});

const DefaultTheme = {
  foreground: "#ffffff", // 字体
  background: "#1b212f", // 背景色
  cursor: "#ffffff", // 设置光标
  selection: "rgba(255, 255, 255, 0.3)",
  black: "#000000",
  brightBlack: "#808080",
  red: "#ce2f2b",
  brightRed: "#f44a47",
  green: "#00b976",
  brightGreen: "#05d289",
  yellow: "#e0d500",
  brightYellow: "#f4f628",
  magenta: "#bd37bc",
  brightMagenta: "#d86cd8",
  blue: "#1d6fca",
  brightBlue: "#358bed",
  cyan: "#00a8cf",
  brightCyan: "#19b8dd",
  white: "#e5e5e5",
  brightWhite: "#ffffff"
};

const themeOptions = [
  { label: "Default", theme: DefaultTheme },
  { label: "Dracula", theme: Dracula },
  { label: "Atom", theme: Atom },
  { label: "Github", theme: Github },
  { label: "Material", theme: Material }
];

const themeMap = new Map();
themeMap.set("Default", DefaultTheme);
themeMap.set("Dracula", Dracula);
themeMap.set("Atom", Atom);
themeMap.set("Github", Github);
themeMap.set("Material", Material);

const xtermConfig = reactive({
  term: null,
  ws: null,
  rows: 41,
  cols: 80,
  heartBeatTimer: null,
  cmdBuffer: ""
});

let xtermThemeLabel = ref("Default");
const changeTheme = function (xtermThemeLabel) {
  xtermConfig.term.setOption("theme", themeMap.get(xtermThemeLabel));
  console.info(xtermConfig.term);
};

const initTerm = function () {
  if (!xtermConfig.term) {
    xtermConfig.term = new Terminal({
      rendererType: "canvas", //渲染类型
      rows: xtermConfig.rows, //行数
      cols: xtermConfig.cols, // 不指定行数，自动回车后光标从下一行开始
      convertEol: true, //启用时，光标将设置为下一行的开头
      // scrollback: 50, //终端中的回滚量
      disableStdin: false, //是否应禁用输入
      // cursorStyle: "underline", //光标样式
      cursorBlink: true, //光标闪烁
      theme: DefaultTheme
    });
    xtermConfig.term.open(document.getElementById("ssh-" + props.sessionId));
  } else {
    xtermConfig.term.clear();
  }

  // 换行并输入起始符 $
  xtermConfig.term.prompt = () => {
    xtermConfig.term.write("\r\n\x1b[33m$\x1b[0m ");
  };
  xtermConfig.term.rn = () => {
    xtermConfig.term.write("\r\n");
  };

  const attachAddon = new AttachAddon(xtermConfig.ws);
  xtermConfig.term.loadAddon(attachAddon);
  const fitAddon = new FitAddon();
  xtermConfig.term.loadAddon(fitAddon);
  fitAddon.fit();
  xtermConfig.term.loadAddon(new WebLinksAddon());

  window.addEventListener("resize", resizeScreen);
  function resizeScreen() {
    try {
      // 窗口大小改变时，触发xterm的resize方法使自适应
      fitAddon.fit();
      xtermConfig.ws.send(
        JSON.stringify({
          type: "resize",
          rows: xtermConfig.rows,
          cols: xtermConfig.cols
        })
      );
    } catch (e) {
      console.log("e", e.message);
    }
  }
  xtermConfig.term.focus();
  runFakeTerminal();
};

const runFakeTerminal = function () {
  if (xtermConfig.term._initialized) return;
  // 初始化
  xtermConfig.term._initialized = true;
  xtermConfig.term.writeln("Welcome to \x1b[1;32mGODP\x1b[0m.");
  xtermConfig.term.writeln("connecting to " + props.host + "....");
  // 添加事件监听器，支持输入方法
  xtermConfig.term.onKey(e => {
    xtermConfig.ws.send(
      JSON.stringify({
        type: "cmd",
        cmd: Base64.encode(e.key), // encode data as base64 format
        rows: xtermConfig.rows,
        cols: xtermConfig.cols
      })
    );
  });
  xtermConfig.term.onData(key => {
    if (key.length > 1) {
      xtermConfig.term.write(key);
    }
  });
};

const state = reactive({
  fullscreen: false,
  teleport: true
});

function doFullscreen() {}

const doLink = function (ev, url) {
  if (ev.type === "click") {
    window.open(url);
  }
};

const doOpen = function () {
  initSocket();
  initTerm();
};

const doClose = function () {
  if (xtermConfig.heartBeatTimer != null) {
    clearInterval(xtermConfig.heartBeatTimer);
  }
  if (xtermConfig.ws) {
    xtermConfig.ws.send(
      JSON.stringify({
        type: "close",
        rows: xtermConfig.rows,
        cols: xtermConfig.cols
      })
    );
    xtermConfig.ws.close();
    xtermConfig.ws = null;
  }
  if (xtermConfig.term) {
    xtermConfig.term.dispose();
  }

  xtermConfig.term._initialized = false;
  xtermConfig.heartBeatTimer = null;
  xtermConfig.cmdBuffer = "";
};

const initSocket = function () {
  xtermConfig.ws = new WebSocket(wsUrl.value);
  xtermConfig.ws.onopen = () => {};
  xtermConfig.ws.onclose = () => {
    delete xtermConfig.ws;
    ElMessage({
      message: "console.web_socket_disconnect",
      type: "warning"
    });
  };
  xtermConfig.ws.onerror = () => {
    ElMessage({
      message: "ws onerror",
      type: "warning"
    });
  };
  xtermConfig.heartBeatTimer = setInterval(function () {
    xtermConfig.ws.send(
      JSON.stringify({
        type: "heartbeat",
        rows: xtermConfig.rows,
        cols: xtermConfig.cols
      })
    );
  }, 20 * 1000);
};

onMounted(() => {
  doOpen();
});
onUnmounted(() => {
  try {
    doClose();
  } catch (err) {
    console.log(err.message);
  }
});
</script>

<template>
  <el-card class="box-card" v-if="props.cardLayout">
    <template #header>
      <div class="card-header">
        <span>命令终端</span>
        <div style="display: flex; align-items: center">
          <el-select
            v-model="xtermThemeLabel"
            class="m-2"
            placeholder="主题选择"
            @change="changeTheme"
          >
            <el-option
              v-for="item in themeOptions"
              :key="item.label"
              :label="item.label"
              :value="item.label"
            />
          </el-select>
          <!-- <span @click="doOpen"
              ><IconifyIconOffline icon="refresh-right"
            /></span>
            <span @click="doOpen"
              ><IconifyIconOffline icon="full-screen"
            /></span> -->
          <el-button class="button" text>文件管理</el-button>
        </div>
      </div>
    </template>
    <div :id="'ssh-' + props.sessionId" />
  </el-card>
  <div v-else :id="'ssh-' + props.sessionId" />
</template>

<style scoped>
.icon {
  margin-left: 10px;
  cursor: pointer;
  color: var(--el-color-primary);
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
