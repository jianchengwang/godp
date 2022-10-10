<template>
  <el-drawer direction="rtl">
    <template #title>
      <h4>{{ props.title }}</h4>
    </template>
    <template #default>
      <div>hello world</div>
    </template>
    <template #footer>
      <div style="flex: auto">
        <el-button @click="cancelClick">取消</el-button>
        <el-button type="primary" @click="confirmClick">确定</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessageBox } from "element-plus";

const props = defineProps({
  title: String,
  formObj: Object
});

const emit = defineEmits(["editFormCancel", "editFormConfirm"]);

const cancelClick = () => {
  emit("editFormCancel");
};

const confirmClick = () => {
  emit("editFormConfirm");
};
</script>

<style scoped></style>
