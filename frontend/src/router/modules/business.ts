const assetsHostRouter = {
  path: "/assetsHost",
  name: "assetsHost",
  redirect: "/assetsHost/index",
  meta: {
    title: "主机管理",
    icon: "monitor",
    rank: 1
  },
  children: [
    {
      path: "/assetsHost/index",
      name: "assetsHostIndex",
      meta: {
        title: "主机管理"
      }
    }
  ]
};

const batchRouter = {
  path: "/batch",
  name: "batch",
  redirect: "/batch/cmd/index",
  meta: {
    title: "批量执行",
    icon: "help",
    rank: 1
  },
  children: [
    {
      path: "/batch/cmd/index",
      name: "batchCmdIndex",
      meta: {
        title: "执行脚本"
      }
    },
    {
      path: "/batch/script/index",
      name: "batchScriptIndex",
      meta: {
        title: "脚本模板"
      }
    },
    {
      path: "/batch/file/index",
      name: "batchFileIndex",
      meta: {
        title: "文件分发"
      }
    }
  ]
};

export { assetsHostRouter, batchRouter };
