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
    let id = node.key;
    let host = node.label;
    let name = node.parent.name;
    emit("click-event", id, host, name);
  }
};

onMounted(() => {
  assetsHostPage({ curPage: 1, limit: 10000 }).then(({ data }) => {
    var menusArr = [];
    const cloudVendors = ["None", "Aliyun", "Tencent"];
    for (let i = 0; i < cloudVendors.length; i++) {
      let menuObj = {
        id: cloudVendors[i],
        label: cloudVendors[i],
        children: []
      };
      let filterCloudVendors = data.list.filter(
        item => item.cloudVendor == cloudVendors[i]
      );
      for (let i = 0; i < filterCloudVendors.length; i++) {
        menuObj.children.push({
          id: data.list[i].id,
          label: data.list[i].ip,
          name: data.list[i].name
        });
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
          <span v-if="data.name" style="color: #e6a23c">({{ data.name }})</span>
        </span>
      </template>
    </el-tree-v2>
  </div>
</template>

<style scoped></style>
