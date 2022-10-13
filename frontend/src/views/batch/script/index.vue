<script lang="ts">
export default {
  name: "batchScriptIndex"
};
</script>

<script setup lang="ts">
import { ref, reactive, toRaw } from "vue";
import { useRouter } from "vue-router";

import {
  VXETable,
  VxeGridInstance,
  VxeGridProps,
  VxeGridListeners
} from "vxe-table";

import { ElMessage, ElForm } from "element-plus";
import { CodemirrorCard } from "/@/components/CodemirrorCard";

import {
  batchScriptPage,
  batchScriptGet,
  batchScriptDelete,
  batchScriptPost,
  batchScriptPut
} from "/@/api/batchScript";

const xGrid = ref({} as VxeGridInstance);

const buttons = [
  {
    name: "创建脚本",
    status: "perfect",
    icon: "vxe-icon--plus",
    code: "showForm"
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
        return batchScriptPage(queryParams).then(({ data }) => {
          return data;
        });
      }
    }
  },
  columns: [
    { type: "seq", title: "序号", width: 60 },
    {
      field: "name",
      title: "文件名"
    },
    {
      field: "remark",
      title: "说明"
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
      case "showForm": {
        showFormEvent(0);
        break;
      }
    }
  },
  toolbarToolClick({ code }) {
    const $grid = xGrid.value;
  }
} as VxeGridListeners);

const deleteEvent = async (id: number) => {
  const type = await VXETable.modal.confirm("您确认删除此脚本吗？");
  if (type === "confirm") {
    batchScriptDelete(id).then(() => {
      VXETable.modal.message({ content: "已删除", status: "success" });
      xGrid.value.commitProxy("reload");
    });
  }
};

const formRef = ref<InstanceType<typeof ElForm>>();
const showForm = ref(false);
const formId = ref(0);
const emptyForm = {
  id: 0,
  name: "",
  content: "",
  args: "",
  remark: ""
};
const formObj = reactive(Object.assign({}, emptyForm));
const showFormEvent = (rowId: number) => {
  formId.value = rowId;
  if (formId.value > 0) {
    batchScriptGet(formId.value).then(({ data }) => {
      Object.assign(formObj, data);
    });
  } else {
    Object.assign(formObj, emptyForm);
  }
  showForm.value = true;
};

const cancelFormEvent = () => {
  showForm.value = false;
  formId.value = 0;
  Object.assign(formObj, emptyForm);
};

const confirmFormEvent = () => {
  if (formId.value) {
    batchScriptPut(formId.value, toRaw(formObj)).then(() => {
      ElMessage.success("保存成功");
      xGrid.value.commitProxy("reload");
      cancelFormEvent();
      xGrid.value.commitProxy("reload");
    });
  } else {
    batchScriptPost(toRaw(formObj)).then(() => {
      ElMessage.success("保存成功");
      cancelFormEvent();
      xGrid.value.commitProxy("reload");
    });
  }
};

const propFilename = ref("");
const propMode = ref("");
const propUpdateString = ref("");
const codemirrorCardVisible = ref(false);

const editScript = () => {
  if (formObj.name.endsWith(".sh") || formObj.name.endsWith(".py")) {
    propFilename.value = formObj.name;
    propUpdateString.value = formObj.content;
    codemirrorCardVisible.value = true;
  } else {
    ElMessage.error("文件名称要以.sh或者.py结尾");
  }
};

const updateText = (_filename, val) => {
  formObj.content = val;
  ElMessage.error("已保存");
};

const codemirrorOnCancel = () => {
  propFilename.value = "";
  propMode.value = "";
  propUpdateString.value = "";
  codemirrorCardVisible.value = false;
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
          placeholder="输入文件名检索"
        />
      </template>
      <template #submit_item>
        <vxe-button type="submit" status="primary" content="查询" />
      </template>
      <template #reset_item>
        <vxe-button type="reset" content="重置" />
      </template>

      <template #opt_field="{ row }">
        <vxe-button
          status="primary"
          type="text"
          content="编辑"
          @click="showFormEvent(row.id)"
        />
        <vxe-button
          status="danger"
          type="text"
          content="删除"
          @click="deleteEvent(row.id)"
        />
      </template>
    </vxe-grid>
    <el-drawer v-model="showForm" direction="rtl">
      <template #title>
        <h4>{{ formId > 0 ? "编辑脚本" : "新建脚本" }}</h4>
      </template>
      <template #default>
        <el-form ref="formRef" :model="formObj" label-width="120px">
          <el-form-item label="文件名">
            <el-input v-model="formObj.name" clearable />
          </el-form-item>
          <el-form-item label="帮助说明">
            <el-input type="textarea" v-model="formObj.remark" clearable />
          </el-form-item>
          <el-form-item label="脚本内容">
            <el-button @click="editScript">编辑</el-button>
          </el-form-item>
        </el-form>
      </template>
      <template #footer>
        <div style="flex: auto">
          <el-button @click="cancelFormEvent">取消</el-button>
          <el-button type="primary" @click="confirmFormEvent">确定</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
  <el-dialog
    :showHeader="false"
    v-if="codemirrorCardVisible"
    v-model="codemirrorCardVisible"
    width="80%"
    show-zoom
    resize
    remember
    fullscreen
    :show-close="false"
  >
    <CodemirrorCard
      :filename="propFilename"
      :updateText="propUpdateString"
      remark="备注：确认无误后请点击保存按钮"
      @save-event="updateText"
      @cancel-event="codemirrorOnCancel()"
      style="margin-top: -40px"
    />
  </el-dialog>
</template>

<style scoped></style>
