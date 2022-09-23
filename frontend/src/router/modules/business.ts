const projectRouter = {
  path: "/project",
  name: "project",
  redirect: "/project/index",
  meta: {
    title: "项目管理",
    icon: "histogram",
    showLink: true,
    rank: 1
  },
  children: [
    {
      path: "/project/index",
      name: "projectIndex",
      meta: {
        title: "项目管理",
        showLink: true
      }
    },
    {
      path: "/project/edit",
      name: "projectEdit",
      meta: {
        title: "",
        showLink: false,
        dynamicLevel: 3,
        refreshRedirect: "/project/index"
      }
    },
    {
      path: "/project/config",
      name: "projectConfig",
      meta: {
        title: "",
        showLink: false,
        dynamicLevel: 3,
        refreshRedirect: "/project/index"
      }
    }
  ]
};

export { projectRouter };
