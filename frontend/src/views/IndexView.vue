<script setup lang="ts">
import {CheckDataUpdate, GetUserSettings, LoadData, SetUserSettings} from '../../wailsjs/go/api/App'
import {onMounted, ref} from "vue";
import {useMessage} from 'naive-ui'
import {data_loader, setting} from "../../wailsjs/go/models";
import SettingStruct = setting.SettingStruct;
import CheckUpdateResult = data_loader.CheckUpdateResult;
import {EventsOn} from "../../wailsjs/runtime";

const conf = ref<SettingStruct>({
  mode: 'pvfutility',
  target: 'http://127.0.0.1:27000',
  concurrent_count: 10,
})

const onModeChange = () => {
  if (conf.value.mode === 'pvfutility' && !conf.value.target.startsWith('http')) {
    conf.value.target = 'http://127.0.0.1:27000'
  } else if (conf.value.mode === 'file' && conf.value.target.startsWith('http')) {
    conf.value.target = ''
  }
}

const message = useMessage()

const save = async () => {
  try {
    await SetUserSettings(conf.value)
    message.success('保存成功')
  } catch (e) {
    message.error(String(e))
  }
}

const updateResult = ref<Array<CheckUpdateResult>>([])
const loading = ref(false)

const checkUpdate = async () => {
  loading.value = true
  try {
    updateResult.value = await CheckDataUpdate()
  } finally {
    loading.value = false
  }
}

const updateItems = ref<string[]>([])

const reloadData = async () => {
  await LoadData(updateItems.value)
}

const selectAll = () => {
  if (updateItems.value.length !== updateResult.value.length) {
    updateItems.value = updateResult.value.map(it => it.key)
  } else {
    updateItems.value = []
  }
}

const selectUpdate = () => {
  updateItems.value = updateResult.value.filter(it => it.need_update).map(it => it.key)
}

onMounted(async () => {
  try {
    conf.value = await GetUserSettings()
    await checkUpdate()
  } catch (e) {
    message.error(String(e))
  }
  EventsOn('progress-end', async (name: string) => {
    checkUpdate()
  })
})
</script>

<template>
  <div>
    <n-card>
      <template #header>
        <div class="text-lg">
          用户设置
        </div>
      </template>
      <n-form label-placement="left" class="mt-4" label-width="120">
        <n-form-item label="数据源">
          <n-radio-group v-model:value="conf.mode" @change="onModeChange">
            <n-space>
              <n-radio value="pvfutility">PvfUtility</n-radio>
              <n-radio value="file">文件系统</n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="PvfUtility地址" v-if="conf.mode === 'pvfutility'">
          <n-input v-model:value="conf.target"></n-input>
        </n-form-item>
        <n-form-item label="文件夹" v-if="conf.mode === 'file'">
          <n-input-group>
            <n-button @click="">选择</n-button>
            <n-input v-model:value="conf.target"></n-input>
          </n-input-group>
        </n-form-item>
        <n-form-item label="并发数">
          <n-input-number v-model:value="conf.concurrent_count" :min="1" :max="100"></n-input-number>
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="float-right">
          <n-button type="primary" @click="save">保存</n-button>
        </div>
      </template>
    </n-card>

    <n-spin :show="loading">
      <n-card class="mt-4" :loading="loading">
        <template #header>
          <div class="text-lg">
            数据更新
          </div>
          <n-form>
            <n-form-item>
              <n-checkbox-group v-model:value="updateItems">
                <n-space>
                  <n-badge v-for="item in updateResult" :key="item.key" :dot="item.need_update">
                    <n-checkbox :value="item.key" :label="item.name"></n-checkbox>
                  </n-badge>
                </n-space>
              </n-checkbox-group>
            </n-form-item>
          </n-form>
        </template>
        <template #footer>
          <n-space>
            <n-button @click="selectAll">全选</n-button>
            <n-button @click="selectUpdate">有更新的</n-button>
            <n-button type="primary" @click="reloadData">重载数据</n-button>
          </n-space>
        </template>
      </n-card>
    </n-spin>
  </div>
</template>

<style scoped>

</style>