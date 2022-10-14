const assetsHostRouter = {
  path: "/assetsHost",
  name: "assetsHost",
  meta: {
    title: "主机管理",
    icon: "monitor",
    rank: 1
  }
};

const batchRouter = {
  path: "/batch",
  name: "batch",
  redirect: "/batch/script",
  meta: {
    title: "批量执行",
    icon: "help",
    rank: 2
  },
  children: [
    {
      path: "/batch/script",
      name: "batchScriptIndex",
      meta: {
        title: "脚本模板"
      }
    },
    {
      path: "/batch/cmd",
      name: "batchCmdIndex",
      meta: {
        title: "执行脚本"
      }
    },
    {
      path: "/batch/file",
      name: "batchFileIndex",
      meta: {
        title: "文件分发"
      }
    }
  ]
};

export { assetsHostRouter, batchRouter };
