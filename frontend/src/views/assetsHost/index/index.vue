<script lang="ts">
export default {
  name: "assetsHostIndex"
};
</script>

<script setup lang="ts">
import { loadEnv } from "@build/index";
import { ref, reactive } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  VXETable,
  VxeGridInstance,
  VxeGridProps,
  VxeGridListeners
} from "vxe-table";
import { useMultiTagsStoreHook } from "/@/store/modules/multiTags";

import {
  assetsHostPage,
  assetsHostGet,
  assetsHostDelete
} from "/@/api/assetsHost";

const route = useRoute();
const router = useRouter();
const xGrid = ref({} as VxeGridInstance);

let searchParam = reactive({});
const buttons = [
  {
    name: "创建主机",
    status: "perfect",
    icon: "vxe-icon--plus",
    code: "toEdit"
  }
];
const gridOptions = reactive({
  stripe: true,
  border: true,
  resizable: true,
  keepSource: true,
  height: 678,
  printConfig: {},
  importConfig: {},
  exportConfig: {},
  pagerConfig: {
    perfect: true,
    pageSize: 10
  },
  formConfig: {
    data: {
      q: "",
      internal: ""
    },
    items: [
      { field: "q", title: "", itemRender: {}, slots: { default: "q_item" } },
      {
        field: "internal",
        title: "",
        itemRender: {},
        slots: { default: "internal_item" }
      },
      { itemRender: {}, slots: { default: "submit_item" } },
      { itemRender: {}, slots: { default: "reset_item" } }
    ]
  },
  toolbarConfig: {
    buttons: buttons,
    perfect: true,
    refresh: {
      icon: "fa fa-refresh",
      iconLoading: "fa fa-spinner fa-spin"
    },
    zoom: {
      iconIn: "fa fa-arrows-alt",
      iconOut: "fa fa-expand"
    },
    custom: {
      icon: "fa fa-cog"
    }
  },
  sortConfig: {
    trigger: "cell",
    remote: true
  },
  proxyConfig: {
    sort: true,
    form: true,
    props: {
      result: "list",
      total: "totalCount"
    },
    ajax: {
      // 接收 Promise
      query: ({ page, sorts, form }) => {
        searchParam = Object.assign({}, form);
        const queryParams: any = Object.assign({}, form);

        // 处理分页数据
        queryParams.curPage = page.currentPage;
        queryParams.limit = page.pageSize;
        // 处理排序
        const firstSort = sorts[0];
        if (firstSort) {
          queryParams.sidx = firstSort.property;
          queryParams.order = firstSort.order;
        }
        return assetsHostPage(queryParams).then(({ data }) => {
          return data;
        });
      }
    }
  },
  columns: [
    { type: "seq", title: "序号", width: 60 },
    {
      field: "name",
      title: "名称"
    },
    {
      field: "ip",
      title: "ip地址"
    },
    {
      field: "verified",
      title: "是否验证"
    },
    {
      field: "cloudVendor",
      title: "云厂商"
    },
    {
      field: "remark",
      title: "备注"
    },
    {
      field: "createdAt",
      title: "创建时间"
    },
    {
      field: "opt",
      title: "操作",
      width: 280,
      slots: { default: "opt_field" }
    }
  ]
} as VxeGridProps);

const gridEvents: VxeGridListeners = reactive({
  toolbarButtonClick({ code }) {
    const $grid = xGrid.value;
    switch (code) {
      case "toEdit": {
        toEditEvent(0);
        break;
      }
    }
  },
  toolbarToolClick({ code }) {
    const $grid = xGrid.value;
  }
} as VxeGridListeners);

const toEditEvent = (rowId: number) => {
  // router.push("/project/edit?id=" + rowId);
  let title = "创建项目";
  if (rowId > 0) {
    title = `项目基础信息 - ID.${rowId}`;
  }
  useMultiTagsStoreHook().handleTags("push", {
    path: `/project/edit`,
    parentPath: route.matched[0].path,
    name: "projectEdit",
    query: { id: String(rowId) },
    meta: {
      title: title,
      showLink: false,
      dynamicLevel: 3
    }
  });
  router.push({ name: "projectEdit", query: { id: String(rowId) } });
};

const deleteEvent = async (id: number) => {
  const type = await VXETable.modal.confirm("您确认删除此主机吗？");
  if (type === "confirm") {
    assetsHostDelete(id).then(() => {
      VXETable.modal.message({ content: "已删除", status: "success" });
      xGrid.value.commitProxy("reload");
    });
  }
};
</script>

<template>
  <div class="main">
    <vxe-grid ref="xGrid" v-bind="gridOptions" v-on="gridEvents">
      <template #q_item="{ data }">
        <vxe-input
          style="margin-left: 5px; width: 305px"
          v-model="data.q"
          type="text"
          placeholder="输入名称/IP检索"
        ></vxe-input>
      </template>
      <template #verified_item="{ data }">
        <vxe-select v-model="data.verified" placeholder="是否验证" clearable>
          <vxe-option value="1" label="是"></vxe-option>
          <vxe-option value="0" label="否"></vxe-option>
        </vxe-select>
      </template>
      <template #submit_item>
        <vxe-button type="submit" status="primary" content="查询"></vxe-button>
      </template>
      <template #reset_item>
        <vxe-button type="reset" content="重置"></vxe-button>
      </template>

      <template #opt_field="{ row }">
        <vxe-button
          status="primary"
          type="text"
          content="编辑"
          @click="toEditEvent(row.id)"
        ></vxe-button>
        <vxe-button
          v-if="row.id === 0"
          status="danger"
          type="text"
          content="删除"
          @click="deleteEvent(row.id)"
        ></vxe-button>
      </template>
    </vxe-grid>
  </div>
</template>

<style scoped>
</style>