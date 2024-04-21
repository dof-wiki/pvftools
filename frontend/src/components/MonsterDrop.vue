<script setup lang="ts">
import {useMessage} from "naive-ui";
import {ref} from "vue";
import DropForm from "./DropForm.vue";

const props = defineProps<{
  title: string,
  getRate:() => Promise<Array<any>>,
  saveRate: (rate: Array<any>) => Promise<void>,
}>()

const message = useMessage()

const monsterRate1 = ref<number[]>([0, 0, 0, 0, 0])
const monsterRate2 = ref<number[]>([0, 0, 0, 0, 0])
const monsterRate3 = ref<number[]>([0, 0, 0, 0, 0])
const monsterRate4 = ref<number[]>([0, 0, 0, 0, 0])
const monsterIsRead = ref(false)
const monsterLoading = ref(false)

const loadRates = async () => {
  try {
    const data = await props.getRate()
    monsterRate1.value = data[0]
    monsterRate2.value = data[1]
    monsterRate3.value = data[2]
    monsterRate4.value = data[3]
    monsterIsRead.value = true
  } catch (e) {
    message.error(String(e))
  }
}

const onRate1Change = (rates: number[]) => {
  monsterRate1.value = rates
}

const onRate2Change = (rates: number[]) => {
  monsterRate2.value = rates
}

const onRate3Change = (rates: number[]) => {
  monsterRate3.value = rates
}

const onRate4Change = (rates: number[]) => {
  monsterRate4.value = rates
}

const saveRates = async () => {
  if (monsterRate1.value.reduce((acc, cur) => acc + cur, 0) !== 100 ||
      monsterRate2.value.reduce((acc, cur) => acc + cur, 0) !== 100 ||
      monsterRate3.value.reduce((acc, cur) => acc + cur, 0) !== 100 ||
      monsterRate4.value.reduce((acc, cur) => acc + cur, 0) !== 100
  ) {
    message.warning('爆率总和必须为100')
    return
  }
  try {
    await props.saveRate([monsterRate1.value, monsterRate2.value, monsterRate3.value, monsterRate4.value])
    message.success('保存成功')
  } catch (e) {
    message.error(String(e))
  }
}
</script>

<template>
  <n-card :title="title">
    <div class="flex justify-around">
      <div class="drop-form-item">
        <p>普通</p>
        <drop-form
            :rates="monsterRate1"
            @change="onRate1Change"
        />
      </div>
      <div class="drop-form-item">
        <p>冒险</p>
        <drop-form
            :rates="monsterRate2"
            :label="false"
            @change="onRate2Change"
        />
      </div>
      <div class="drop-form-item">
        <p>王者</p>
        <drop-form
            :rates="monsterRate3"
            :label="false"
            @change="onRate3Change"
        />
      </div>
      <div class="drop-form-item">
        <p>地狱</p>
        <drop-form
            :rates="monsterRate4"
            :label="false"
            @change="onRate4Change"
        />
      </div>
    </div>
    <template #footer>
      <n-space style="float:right;">
        <n-button
            type="warning"
            :loading="monsterLoading"
            @click="loadRates"
        >
          {{ monsterIsRead ? '重载' : '加载' }}
        </n-button>
        <n-button
            type="primary"
            :loading="monsterLoading"
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