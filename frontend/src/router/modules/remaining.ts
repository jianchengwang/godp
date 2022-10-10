import { $t } from "/@/plugins/i18n";
import type { RouteConfigsTable } from "/#/index";
const Layout = () => import("/@/layout/index.vue");

const remainingRouter: Array<RouteConfigsTable> = [
  {
    path: "/login",
    name: "Login",
    component: () => import("/@/views/login/index.vue"),
    meta: {
      title: $t("menus.hslogin"),
      showLink: false,
      rank: 101
    }
  },
  {
    path: "/redirect",
    component: Layout,
    meta: {
      icon: "home-filled",
      title: $t("menus.hshome"),
      showLink: false,
      rank: 104
    },
    children: [
      {
        path: "/redirect/:path(.*)",
        name: "Redirect",
        component: () => import("/@/layout/redirect.vue")
      }
    ]
  },
  {
    path: "/ssh",
    name: "SSH",
    component: () => import("/@/views/ssh/index.vue"),
    meta: {
      title: "web终端",
      showLink: false,
      rank: 200
    }
  }
];

export default remainingRouter;
