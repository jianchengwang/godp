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

import { CodemirrorCard } from "/@/components/CodemirrorCard";

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
  remoteFileList: []
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
    cmdListFile();
  };
  sftpConfig.ws.onmessage = ({ data }) => {
    if (data) {
      let jsonData = JSON.parse(data);
      if (jsonData.cmd) {
        if (jsonData.success) {
          switch (jsonData.cmd) {
            case "list":
              {
                sftpConfig.remoteFileList = jsonData.data;
              }
              break;
            case "fetch":
              {
                propUpdateString.value = jsonData.data;
                codemirrorCardVisible.value = true;
              }
              break;
            case "update": {
              ElMessage.success("保存成功");
              propUpdateString.value = "";
              codemirrorCardVisible.value = false;
              break;
            }
            default:
              break;
          }
        } else {
          ElMessage.error(
            jsonData.data ? jsonData.data : "webstock run error."
          );
        }
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

const cmdListFile = () => {
  sftpConfig.ws.send(
    JSON.stringify({
      type: "cmd",
      cmd: "list",
      remoteFilePath: Base64.encode(sftpConfig.remoteFileDir)
    })
  );
  tmpRemoteFileDir.value = sftpConfig.remoteFileDir;
};

const cmdFetchFileText = (fileName: String) => {
  sftpConfig.ws.send(
    JSON.stringify({
      type: "cmd",
      cmd: "fetch",
      remoteFilePath: Base64.encode(sftpConfig.remoteFileDir + "/" + fileName)
    })
  );
};

const cmdUpdateText = (fileName: string, val: string) => {
  sftpConfig.ws.send(
    JSON.stringify({
      type: "cmd",
      cmd: "update",
      remoteFilePath: Base64.encode(sftpConfig.remoteFileDir + "/" + fileName),
      updateString: Base64.encode(val)
    })
  );
};

const showEditRemoteFileDir = ref(false);
const tmpRemoteFileDir = ref("");
const changeRemoteFileDir = () => {
  if (tmpRemoteFileDir.value.startsWith("/")) {
    sftpConfig.remoteFileDir = tmpRemoteFileDir.value;
    cmdListFile();
  } else {
    alert("please input dir");
  }
  showEditRemoteFileDir.value = false;
  tmpRemoteFileDir.value = sftpConfig.remoteFileDir;
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
const dbClickEvent = row => {
  let fileName = row.name;
  if (fileName.endsWith("/")) {
    sftpConfig.remoteFileDir =
      sftpConfig.remoteFileDir +
      "/" +
      fileName.substring(0, fileName.length - 1);
    cmdListFile();
  } else {
    // if (isAssetTypeText(fileName.split(".").pop().toLowerCase())) {
    //   cmdFetchFileText(fileName);
    // }
    propFilename.value = fileName;
    cmdFetchFileText(fileName);
  }
};

const propFilename = ref("");
const propUpdateString = ref("");
const codemirrorCardVisible = ref(false);
const codemirrorOnCancel = () => {
  propFilename.value = "";
  propUpdateString.value = "";
  codemirrorCardVisible.value = false;
};

const filterType = ref("");
const remoteFileList = computed(() => {
  if (sftpConfig.remoteFileList) {
    if (filterType.value == "file") {
      return sftpConfig.remoteFileList.filter(item => !item.name.endsWith("/"));
    } else if (filterType.value == "dir") {
      return sftpConfig.remoteFileList.filter(item => item.name.endsWith("/"));
    } else if (filterType.value == "notShowHiddenFile") {
      return sftpConfig.remoteFileList.filter(
        item => !item.name.startsWith(".")
      );
    }
    return sftpConfig.remoteFileList;
  }
  return [];
});
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
          style="width: 400px"
        />
        <el-link v-else @click="showEditRemoteFileDir = true">{{
          sftpConfig.remoteFileDir
        }}</el-link>
      </div>

      <div>
        <el-select
          v-model="filterType"
          class="m-2"
          placeholder="Select"
          size="small"
        >
          <el-option key="0" label="全部" value="" />
          <el-option key="1" label="文件" value="file" />
          <el-option key="2" label="目录" value="dir" />
          <el-option key="3" label="隐藏.文件" value="notShowHiddenFile" />
        </el-select>
        <!-- <el-switch
          v-model="sftpConfig.showHiddenFile"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          :width="80"
        /> -->
      </div>
    </div>
    <el-divider />
    <el-table
      :data="remoteFileList"
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

  <vxe-modal
    :showHeader="false"
    v-if="codemirrorCardVisible"
    v-model="codemirrorCardVisible"
    width="80%"
    show-zoom
    resize
    remember
    fullscreen
  >
    <CodemirrorCard
      :filename="propFilename"
      :updateText="propUpdateString"
      remark="备注: 保存会直接覆盖服务器的文件"
      @save-event="cmdUpdateText"
      @cancel-event="codemirrorOnCancel()"
    />
  </vxe-modal>
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
