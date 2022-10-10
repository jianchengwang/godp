const assetsHostRouter = {
  path: "/assetsHost",
  name: "assetsHost",
  redirect: "/assetsHost/index",
  meta: {
    title: "主机管理",
    icon: "histogram",
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

export { assetsHostRouter };
