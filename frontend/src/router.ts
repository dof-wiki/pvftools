import {createRouter, createWebHashHistory, RouteRecordRaw} from "vue-router";
import IndexView from "./views/IndexView.vue";

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        component: IndexView
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router