import {createRouter, createWebHashHistory, RouteRecordRaw} from "vue-router";
import IndexView from "./views/IndexView.vue";
import WorldDropView from "./views/WorldDropView.vue";
import DropView from "./views/DropView.vue";

const routes: RouteRecordRaw[] = [
    { path: '/', component: IndexView },
    { path: '/worlddrop', component: WorldDropView },
    { path: '/drop', component: DropView },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router