<script setup lang="ts">
import {CopyOutline, Search, Dice, GiftOutline, GlobeOutline, Nuclear, SettingsOutline, Aperture, Shirt, ArrowUp, HelpOutline} from "@vicons/ionicons5";
import {MenuOption, NIcon, useMessage} from "naive-ui";
import {h, onMounted, reactive, ref} from "vue";
import {RouterLink} from "vue-router";
import {EventsOn} from "../../wailsjs/runtime";
import storage from "../common/storage";
import {DownloadNewVersion} from "../../wailsjs/go/api/App";

const menuItems = [
  {
    name: '设置',
    icon: SettingsOutline,
    to: '/'
  },
  {
    name: '搜索',
    icon: Search,
    to: '/search'
  },
  {
    name: '世界掉落',
    icon: GlobeOutline,
    to: '/worlddrop'
  },
  {
    name: '其他掉落',
    icon: Dice,
    to: '/drop'
  },
  {
    name: '礼盒生成器',
    icon: GiftOutline,
    to: '/box'
  },
  {
    name: '技能修改',
    icon: Nuclear,
    to: '/skill'
  },
  {
    name: '装备批处理',
    icon: CopyOutline,
    to: '/batch',
  },
  {
    name: '异次元气息',
    icon: Aperture,
    to: '/breath',
  },
  {
    name: '装备属性生成',
    icon: Shirt,
    to: '/equ_attr',
  },
  {
    name: '强化/增幅',
    icon: ArrowUp,
    to: '/upgrade',
  },
  {
    name: '任务生成器',
    icon: HelpOutline,
    to: '/quest',
  }
]

const menuOptions: MenuOption[] = menuItems.map((item) => {
  return {
    label: () => h(RouterLink, {to: item.to}, {default: () => item.name}),
    key: item.to,
    icon: () => h(NIcon, null, {default: () => h(item.icon)})
  }
})

type ProgressController = {
  name: string
  total: number
  current: number
}

const progressList = ref<ProgressController[]>([])

const collapsed = ref(false)

const message = useMessage()

const selfUpdateCtx = reactive({
  show: false,
  version: '',
  message: '',
})

const downloadNewVersion = async () => {
  selfUpdateCtx.show = false
  await DownloadNewVersion()
}

onMounted(() => {
  EventsOn('progress-start', async (name: string, total: number) => {
    progressList.value.push({name, total, current: 0})
  })
  EventsOn('progress-update', async (name: string, current: number) => {
    console.log('progress-update', name, current)
    const index = progressList.value.findIndex(item => item.name === name)
    if (index !== -1) {
      progressList.value[index].current = current
    }
  })
  EventsOn('progress-end', async (name: string) => {
    progressList.value = progressList.value.filter(item => item.name !== name)
  })

  EventsOn('error', async (msg: string) => {
    message.error(msg)
  })

  EventsOn('job-update', async () => {
    await storage.loadJobMap()
  })

  EventsOn('skill-update', async () => {
    await storage.loadSkillMap()
  })

  EventsOn('has-self-update', async (version: string, msg: string) => {
    selfUpdateCtx.version = version
    selfUpdateCtx.message = msg
    selfUpdateCtx.show = true
  })

  storage.loadJobMap()
})
</script>

<template>
  <n-layout class="h-screen" has-sider>
    <n-layout-sider
        :width="200"
        bordered
        collapse-mode="width"
        :collapsed-width="64"
        :collapsed="collapsed"
        @collapse="collapsed = true"
        @expand="collapsed = false"
        show-trigger
    >
      <n-menu
          :options="menuOptions"
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          default-value="/"
      ></n-menu>
      <div class="p-4" v-if="progressList.length > 0">
        <div class="font-bold text-yellow-500">数据加载进度</div>
        <div v-for="item in progressList" :key="item.name" class="flex items-center">
          <span class="w-10">{{ item.name}}</span>
          <n-progress
              :percentage="(item.current * 100 / item.total).toFixed(2)"
              :indicator-placement="'inside'"
              processing
          ></n-progress>
        </div>
      </div>
    </n-layout-sider>
    <n-layout-content class="p-[30px] h-full">
      <router-view></router-view>
    </n-layout-content>
  </n-layout>

  <n-modal v-model:show="selfUpdateCtx.show">
    <n-card style="width: 50%;">
      <template #header>
        <div class="text-lg">
          有新版本可用: {{ selfUpdateCtx.version }}
        </div>
      </template>
      <n-input readonly type="textarea" autosize v-model:value="selfUpdateCtx.message"></n-input>
      <template #footer>
        <n-flex justify="end">
          <n-button @click="selfUpdateCtx.show = false">暂不更新</n-button>
          <n-button type="primary" @click="downloadNewVersion">立即更新</n-button>
        </n-flex>
      </template>
    </n-card>
  </n-modal>
</template>

<style scoped>

</style>
