<script setup lang="ts">
import {computed, onMounted, ref, watch} from "vue";

const props = withDefaults(defineProps<{
  rates: number[],
  label: boolean,
}>(), { label: true })

const emits = defineEmits(['change'])

const rateList = ref<number[]>([0, 0, 0, 0, 0])
const labels = ['白云', '蓝天', '稀有', '神器', '史诗']

const sum = computed(() => {
  return rateList.value.reduce((acc, cur) => acc + cur, 0)
})

const onRateChange = async () => {
  emits('change', rateList.value)
}

watch(props, (val) => {
  rateList.value = val.rates
})

onMounted(() => { rateList.value = props.rates })
</script>

<template>
<n-form label-placement="left">
  <n-form-item
      v-for="(rate, idx) in rateList"
      :key="idx"
      :label="label ? labels[idx] : ''"
  >
    <n-input-number
        v-model:value="rateList[idx]"
        :show-button="false"
        @update:value="onRateChange"
    >
      <template #suffix>
        %
      </template>
    </n-input-number>
  </n-form-item>
  <span>总掉率: <span :class="{'red': sum !== 100}">{{ sum }}</span>%</span>
</n-form>
</template>

<style scoped>

</style>