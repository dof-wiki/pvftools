<script setup lang="ts">
import {GetClearRewardDrop, SetClearRewardDrop} from "../../wailsjs/go/api/App";
import {ref} from "vue";
import {useMessage} from "naive-ui";
import DropForm from "./DropForm.vue";

const message = useMessage()

const rate = ref<number[]>([0, 0, 0, 0, 0])
const rateIsRead = ref(false)
const isLoading = ref(false)

const loadRates = async () => {
  try {
    rate.value = await GetClearRewardDrop()
    rateIsRead.value = true
  } catch (e) {
    message.error(String(e))
  }
}

const onRateChange = (rates: number[]) => {
  rate.value = rates
}

const saveRates = async () => {
  if (rate.value.reduce((acc, cur) => acc + cur, 0) !== 100) {
    message.warning('爆率总和必须为100')
    return
  }
  try {
    await SetClearRewardDrop(rate.value)
    message.success('保存成功')
  } catch (e) {
    message.error(String(e))
  }
}
</script>

<template>
  <n-card title="翻牌">
    <div class="flex justify-around">
      <div class="p-2 text-center">
        <p>困难</p>
        <drop-form
            :rates="rate"
            @change="onRateChange"
        />
      </div>
    </div>
    <template #footer>
      <n-space style="float:right;">
        <n-button
            type="warning"
            :loading="isLoading"
            @click="loadRates"
        >
          {{ rateIsRead ? '重载' : '加载' }}
        </n-button>
        <n-button
            type="primary"
            :loading="isLoading"
            @click="saveRates"
        >
          保存
        </n-button>
      </n-space>
    </template>
  </n-card>
</template>

<style scoped>

</style>