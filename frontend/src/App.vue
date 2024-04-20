<script lang="ts" setup>
import {darkTheme, MenuOption, NIcon} from 'naive-ui'
import {GlobeOutline, SettingsOutline} from "@vicons/ionicons5";
import {h, onMounted, ref} from "vue";
import {RouterLink} from "vue-router";
import {EventsOn} from "../wailsjs/runtime";

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
})
</script>

<template>
  <n-config-provider :theme="darkTheme">
    <n-message-provider>
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
    </n-message-provider>
  </n-config-provider>
</template>

<style>
</style>
