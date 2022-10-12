<script lang="ts">
export default {
  name: "SSH"
};
</script>

<script setup lang="ts">
import splitpane, { ContextProps } from "/@/components/ReSplitPane";
import { ref, reactive, computed, onMounted } from "vue";
import { buildUUID } from "/@/utils/uuid";
import { ElMessage, ElNotification } from "element-plus";
import { HostTree } from "/@/components/HostTree";
import { WebTerminal } from "/@/components/WebTerminal";

const settingLR: ContextProps = reactive({
  minPercent: 12,
  defaultPercent: 15,
  split: "vertical"
});

const clickEvent = function (id, host, hostName, hostRemark) {
  addSshTab(id, host, hostName, hostRemark);
};

let tabIndex = 1;
const sshTabsValue = ref("1");
const sshTabs = ref([]);

const addSshTab = (
  id: Number,
  host: String,
  hostName: String,
  hostRemark: String
) => {
  let sessionId = buildUUID();
  const newTabName = `ssh-${host}-${sessionId.slice(-4)}`;
  sshTabs.value.push({
    title: newTabName,
    name: newTabName,
    id: id,
    host: host,
    sessionId: sessionId,
    hostName: hostName,
    hostRemark: hostRemark
  });
  sshTabsValue.value = newTabName;
};

const removeSshTab = (targetName: string) => {
  const tabs = sshTabs.value;
  let activeName = sshTabsValue.value;
  if (activeName === targetName) {
    tabs.forEach((tab, index) => {
      if (tab.name === targetName) {
        const nextTab = tabs[index + 1] || tabs[index - 1];
        if (nextTab) {
          activeName = nextTab.name;
        }
      }
    });
  }

  sshTabsValue.value = activeName;
  sshTabs.value = tabs.filter(tab => tab.name !== targetName);
};

const activeName = ref("WebTerminal");
const handleClick = (tab, event: Event) => {
  console.log(tab, event);
};
</script>

<template>
  <div class="main">
    <div class="split-pane">
      <splitpane :splitSet="settingLR">
        <template #paneL>
          <div class="dv-a">
            <HostTree @click-event="clickEvent" />
          </div>
        </template>
        <template #paneR>
          <div class=".dv-b">
            <el-tabs
              v-model="sshTabsValue"
              type="card"
              closable
              @tab-remove="removeSshTab"
            >
              <el-tab-pane
                v-for="item in sshTabs"
                :key="item.name"
                :label="item.title"
                :name="item.name"
              >
                <div>
                  <WebTerminal
                    :host="item.host"
                    :sessionId="item.sessionId"
                    :cardLayout="true"
                    :hostName="item.hostName"
                    :hostRemark="item.hostRemark"
                  />
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </template>
      </splitpane>
    </div>
  </div>
</template>

<style lang="scss" scoped>
$W: 100%;
$H: 95vh;

.main {
  // background: #fff;
}

.split-pane {
  width: 100%;
  height: $H;
  font-size: 14px;
  color: #fff;
  border: 1px solid #e5e6eb;
  .dv-a,
  .dv-b {
    width: $W;
    height: $W;
    color: rgba($color: dodgerblue, $alpha: 0.8);
  }
}
</style>
