<script lang="ts">
export default {
  name: "WebSftp"
};
</script>

<script lang="ts" setup>
import { ElMessage } from "element-plus";
import { loadEnv } from "@build/index";
import { ref, reactive, toRaw, computed, onMounted, onUnmounted } from "vue";

import { getAccessToken } from "/@/utils/auth";

import { Base64 } from "js-base64";

const { VITE_PROXY_DOMAIN_REAL } = loadEnv();

const props = defineProps({
  host: String,
  workDir: String,
  sessionId: String
});
const emit = defineEmits(["hide-event", "clear-event"]);

let wsUrl = computed(() => {
  return (
    `ws://localhost:8081/api/v1/ws/ssh/sftp-${props.sessionId}?host=${props.host}&Authorization=` +
    getAccessToken()
  );
});

const sftpConfig = reactive({
  ws: null,
  heartBeatTimer: null,
  showHiddenFile: true,
  remoteFileDir: props.workDir,
  remoteFileList: [],
  updateString: ""
});

const doClose = function () {
  if (sftpConfig.heartBeatTimer != null) {
    clearInterval(sftpConfig.heartBeatTimer);
  }
  if (sftpConfig.ws) {
    sftpConfig.ws.send(
      JSON.stringify({
        type: "close"
      })
    );
    sftpConfig.ws.close();
    sftpConfig.ws = null;
  }
  sftpConfig.heartBeatTimer = null;
};

const initSocket = function () {
  sftpConfig.ws = new WebSocket(wsUrl.value);
  sftpConfig.ws.onopen = () => {
    listFile();
  };
  sftpConfig.ws.onmessage = ({ data }) => {
    if (data) {
      let jsonData = JSON.parse(data);
      if (jsonData.cmd && jsonData.cmd === "list") {
        sftpConfig.remoteFileList = jsonData.data;
      }
      if (jsonData.cmd && jsonData.cmd === "get") {
        sftpConfig.updateString = jsonData.data;
        alert(sftpConfig.updateString);
      }
    }
  };
  sftpConfig.ws.onclose = () => {
    delete sftpConfig.ws;
    ElMessage({
      message: "console.web_socket_disconnect",
      type: "warning"
    });
  };
  sftpConfig.ws.onerror = () => {
    ElMessage({
      message: "ws onerror",
      type: "warning"
    });
  };
  sftpConfig.heartBeatTimer = setInterval(function () {
    sftpConfig.ws.send(
      JSON.stringify({
        type: "heartbeat"
      })
    );
  }, 20 * 1000);
};

onMounted(() => {
  initSocket();
});
onUnmounted(() => {
  try {
    doClose();
  } catch (err) {
    console.log(err.message);
  }
});

const listFile = () => {
  sftpConfig.ws.send(
    JSON.stringify({
      type: "cmd",
      cmd: "list",
      remoteFilePath: Base64.encode(sftpConfig.remoteFileDir)
    })
  );
  tmpRemoteFileDir.value = sftpConfig.remoteFileDir;
};

const getFileText = (fileName: String) => {
  sftpConfig.ws.send(
    JSON.stringify({
      type: "cmd",
      cmd: "get",
      remoteFilePath: Base64.encode(sftpConfig.remoteFileDir + "/" + fileName)
    })
  );
};

const isAssetTypeText = ext => {
  return (
    [
      ".txt",
      ".conf",
      ".yaml",
      ".yml",
      ".json",
      ".sh",
      ".go",
      ".py",
      ".java",
      ".md"
    ].indexOf(ext.toLowerCase()) !== -1
  );
};

const showEditRemoteFileDir = ref(false);
const tmpRemoteFileDir = ref("");
const changeRemoteFileDir = () => {
  if (tmpRemoteFileDir.value.startsWith("/")) {
    sftpConfig.remoteFileDir = tmpRemoteFileDir.value;
    listFile();
  } else {
    alert("please input dir");
  }
  showEditRemoteFileDir.value = false;
  tmpRemoteFileDir.value = sftpConfig.remoteFileDir;
};

const dbClickEvent = row => {
  let fileName = row.name;
  if (fileName.endsWith("/")) {
    sftpConfig.remoteFileDir =
      sftpConfig.remoteFileDir +
      "/" +
      fileName.substring(0, fileName.length - 1);
    listFile();
  } else {
    // if (isAssetTypeText(fileName.split(".").pop().toLowerCase())) {
    //   getFileText(fileName);
    // }
    getFileText(fileName);
  }
};

const downloadFile = (fileName: String) => {};
</script>

<template>
  <div :id="'sftp-' + props.sessionId">
    <div class="card-header">
      <div style="display: flex; align-items: center">
        <span class="el-icon-home"
          ><IconifyIconOffline icon="home-filled"
        /></span>
        <el-input
          v-if="showEditRemoteFileDir"
          v-model="tmpRemoteFileDir"
          placeholder=""
          @change="changeRemoteFileDir"
        />
        <el-link v-else @click="showEditRemoteFileDir = true">{{
          sftpConfig.remoteFileDir
        }}</el-link>
      </div>

      <div>
        <el-switch
          v-model="sftpConfig.showHiddenFile"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          :width="80"
        />
      </div>
    </div>
    <el-divider />
    <el-table
      :data="sftpConfig.remoteFileList"
      stripe
      style="width: 100%"
      max-height="720"
      @cell-dblclick="dbClickEvent"
    >
      <el-table-column prop="name" label="名称" sortable />
      <el-table-column prop="size" label="大小" width="180" sortable />
      <el-table-column prop="modTime" label="修改时间" width="180" sortable />
    </el-table>
  </div>
</template>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-icon-home {
  height: 48px;
  width: 38px;
  padding: 12px;
  display: flex;
  cursor: pointer;
  align-items: center;
  color: gray;
}
</style>
