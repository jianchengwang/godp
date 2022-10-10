<script lang="ts">
export default {
  name: "HostTree"
};
</script>

<script lang="ts" setup>
import type { ElTreeV2 } from "element-plus";
import type { TreeNode } from "element-plus/es/components/tree-v2/src/types";
import { ref, reactive, computed, onMounted } from "vue";

import { assetsHostPage } from "/@/api/assetsHost";

const props = defineProps({});

const emit = defineEmits(["click-event"]);

const query = ref("");
let dataProps = ref({
  value: "id",
  label: "label",
  children: "children"
});
const treeRef = ref<InstanceType<typeof ElTreeV2>>();

const menusData = ref([]);

const onQueryChanged = (query: string) => {
  // @ts-expect-error
  treeRef.value!.filter(query);
};

const filterMethod = (query: string, node: TreeNode) => {
  return node.label!.includes(query);
};

const nodeClickEvent = function (data, node, e) {
  if (node.parent) {
    let type = node.data.type;
    let projectId = node.parent.key;
    let host = node.key;
    let projectNmame = node.parent.label;
    let projectApp = node.parent.data.code;
    emit("click-event", type, projectId, host, projectNmame, projectApp);
  }
};

onMounted(() => {
  assetsHostPage({ curPage: 1, limit: 10000 }).then(({ data }) => {
    var menusArr = [];
    for (let i = 0; i < data.list.length; i++) {
      let hostObj = data.list[i];
      let menuObj = {
        id: hostObj.id,
        label: hostObj.name,
        children: []
      };
      let projectConfig = hostObj.projectConfig;
      // ci
      menuObj.children.push({
        id: projectConfig.ci.ip,
        label: projectConfig.ci.ip,
        type: "ci"
      });
      for (let j = 0; j < projectConfig.ip_address_arr.length; j++) {
        if (projectConfig.ip_address_arr[j].ip) {
          menuObj.children.push({
            id: projectConfig.ip_address_arr[j].ip,
            label: projectConfig.ip_address_arr[j].ip,
            type: j == 0 ? "server" : "client"
          });
        }
      }
      menusArr.push(menuObj);
    }
    menusData.value = menusArr;
  });
});
</script>

<template>
  <div
    style="margin-top: 5px; margin-left: 5px; marigin-right: 5px; height: 100%"
  >
    <el-input
      class="mb-4"
      v-model="query"
      placeholder="请输入关键字查找"
      clearable
      @input="onQueryChanged"
    />
    <el-tree-v2
      ref="treeRef"
      :data="menusData"
      :props="dataProps"
      :filter-method="filterMethod"
      @node-click="nodeClickEvent"
      :height="650"
    >
      <template #default="{ data }">
        <span
          >{{ data.label }}
          <span v-if="data.type === 'ci'" style="color: #e6a23c">(ci)</span>
          <span v-else-if="data.type === 'server'" style="color: #67c23a"
            >(server)</span
          >
        </span>
      </template>
    </el-tree-v2>
  </div>
</template>

<style scoped></style>
