<script setup lang="ts">
import {CopyOutline, Dice, GiftOutline, GlobeOutline, Nuclear, SettingsOutline} from "@vicons/ionicons5";
import {MenuOption, NIcon, useMessage} from "naive-ui";
import {h, onMounted, ref} from "vue";
import {RouterLink} from "vue-router";
import {EventsOn} from "../../wailsjs/runtime";

const menuItems = [
  {
    name: '设置',
    icon: SettingsOutline,
    to: '/'
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
    to: '/equbatch',
  },
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

const message = useMessage()

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
})
</script>

<template>
  <n-layout class="h-screen" has-sider>
    <n-layout-sider
        :width="200"
        bordered
    >
      <n-menu :options="menuOptions" default-value="/"></n-menu>
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
</template>

<style scoped>

</style>