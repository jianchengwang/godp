<script lang="ts">
export default {
  name: "assetsHostIndex"
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

import { ElMessage } from "element-plus";
import type { ElForm } from "element-plus";

import {
  assetsHostPage,
  assetsHostGet,
  assetsHostDelete,
  assetsHostPost,
  assetsHostPut
} from "/@/api/assetsHost";

const xGrid = ref({} as VxeGridInstance);

const buttons = [
  {
    name: "创建主机",
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
      {
        field: "verified",
        title: "",
        itemRender: {},
        slots: { default: "verified_item" }
      },
      { itemRender: {}, slots: { default: "submit_item" } },
      { itemRender: {}, slots: { default: "reset_item" } },
      { itemRender: {}, slots: { default: "webterminal_item" } }
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
  const type = await VXETable.modal.confirm("您确认删除此主机吗？");
  if (type === "confirm") {
    assetsHostDelete(id).then(() => {
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
  cloudVendor: "None",
  ip: "127.0.0.1",
  port: 22,
  user: "root",
  password: "",
  remark: "",
  workDir: "/root/godp"
};
const formObj = reactive(Object.assign({}, emptyForm));
const showFormEvent = (rowId: number) => {
  formId.value = rowId;
  if (formId.value > 0) {
    assetsHostGet(formId.value).then(({ data }) => {
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
    assetsHostPut(formId.value, toRaw(formObj)).then(() => {
      ElMessage.success("保存成功");
      xGrid.value.commitProxy("reload");
      cancelFormEvent();
      xGrid.value.commitProxy("reload");
    });
  } else {
    assetsHostPost(toRaw(formObj)).then(() => {
      ElMessage.success("保存成功");
      cancelFormEvent();
      xGrid.value.commitProxy("reload");
    });
  }
};

const router = useRouter();
const sshHref = router.resolve({
  name: "ssh",
  path: "/ssh"
});
const openWebTerminal = () => {
  window.open(sshHref.href, "_blank");
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
        />
      </template>
      <template #verified_item="{ data }">
        <vxe-select v-model="data.verified" placeholder="是否验证" clearable>
          <vxe-option value="1" label="是" />
          <vxe-option value="2" label="否" />
        </vxe-select>
      </template>
      <template #submit_item>
        <vxe-button type="submit" status="primary" content="查询" />
      </template>
      <template #reset_item>
        <vxe-button type="reset" content="重置" />
      </template>
      <template #webterminal_item>
        <vxe-button
          @click="openWebTerminal"
          status="primary"
          content="Web终端"
        />
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
        <h4>{{ formId > 0 ? "编辑主机" : "新建主机" }}</h4>
      </template>
      <template #default>
        <el-form ref="formRef" :model="formObj" label-width="120px">
          <el-form-item label="名称">
            <el-input v-model="formObj.name" clearable />
          </el-form-item>
          <el-form-item label="云厂商">
            <el-select
              v-model="formObj.cloudVendor"
              placeholder="请选择云厂商"
              clearable
            >
              <el-option label="None" value="None" />
              <el-option label="Aliyun" value="Aliyun" />
              <el-option label="Tencent" value="Tencent" />
            </el-select>
          </el-form-item>
          <el-form-item label="ip地址">
            <el-input v-model="formObj.ip" clearable />
          </el-form-item>
          <el-form-item label="端口">
            <el-input v-model.number="formObj.port" clearable />
          </el-form-item>
          <el-form-item label="用户">
            <el-input v-model="formObj.user" clearable />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="formObj.password" show-password clearable />
          </el-form-item>
          <el-form-item label="备注">
            <el-input type="textarea" v-model="formObj.remark" clearable />
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
</template>

<style scoped></style>
