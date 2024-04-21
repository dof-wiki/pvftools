<script setup lang="ts">
import {useMessage} from "naive-ui";
import {ref} from "vue";
import {GetHellDrop, SetHellDrop} from "../../wailsjs/go/api/App";
import DropForm from "./DropForm.vue";

const message = useMessage()

const hallRate1 = ref<number[]>([0, 0, 0, 0, 0])
const hallRate2 = ref<number[]>([0, 0, 0, 0, 0])
const hallIsRead = ref(false)
const hallLoading = ref(false)

const loadHallRates = async () => {
  try {
    const data = await GetHellDrop()
    hallRate1.value = data.slice(0, 5)
    hallRate2.value = data.slice(5)
    hallIsRead.value = true
  } catch (e) {
    message.error(String(e))
  }
}

const onHallRate1Change = (rates: number[]) => {
  hallRate1.value = rates
}

const onHallRate2Change = (rates: number[]) => {
  hallRate2.value = rates
}

const saveRates = async () => {
  if (hallRate1.value.reduce((acc, cur) => acc + cur, 0) !== 100 ||
      hallRate2.value.reduce((acc, cur) => acc + cur, 0) !== 100) {
    message.warning('爆率总和必须为100')
    return
  }
  try {
    await SetHellDrop(hallRate1.value, hallRate2.value)
    message.success('保存成功')
  } catch (e) {
    message.error(String(e))
  }
}
</script>

<template>
  <n-card title="深渊">
    <div class="flex justify-around">
      <div class="p-2 text-center">
        <p>困难</p>
        <drop-form
            :rates="hallRate1"
            @change="onHallRate1Change"
        />
      </div>
      <div class="p-2 text-center">
        <p>非常困难</p>
        <drop-form
            :rates="hallRate2"
            @change="onHallRate2Change"
        />
      </div>
    </div>
    <template #footer>
      <n-space style="float:right;">
        <n-button
            type="warning"
            :loading="hallLoading"
            @click="loadHallRates"
        >
          {{ hallIsRead ? '重载' : '加载' }}
        </n-button>
        <n-button
            type="primary"
            :loading="hallLoading"
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