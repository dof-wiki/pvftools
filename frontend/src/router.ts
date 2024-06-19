import {createRouter, createWebHashHistory, RouteRecordRaw} from "vue-router";
import IndexView from "./views/IndexView.vue";
import WorldDropView from "./views/WorldDropView.vue";
import DropView from "./views/DropView.vue";
import BatchHandleView from "./views/BatchHandleView.vue";
import SkillView from "./views/SkillView.vue";
import BoxView from "./views/BoxView.vue";
import BreathView from "./views/BreathView.vue";
import EquAttrView from "./views/EquAttrView.vue";
import UpgradeView from "./views/UpgradeView.vue";
import QuestView from "./views/QuestView.vue";

const routes: RouteRecordRaw[] = [
    { path: '/', component: IndexView },
    { path: '/worlddrop', component: WorldDropView },
    { path: '/drop', component: DropView },
    { path: '/batch', component: BatchHandleView },
    { path: '/skill', component: SkillView },
    { path: '/box', component: BoxView },
    { path: '/breath', component: BreathView },
    { path: '/equ_attr', component: EquAttrView },
    { path: '/upgrade', component: UpgradeView },
    { path: '/quest', component: QuestView },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router
