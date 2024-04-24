<script setup lang="ts">
import {api} from "../../wailsjs/go/models";
import SkillData = api.SkillData;
import {computed, ref} from "vue";

const props = defineProps<{
  skillDetail: SkillData,
  onEdit: (idx: number) => void,
}>()

const showLevel = ref(1)

const maxLevel = computed(() => {
  if (!props.skillDetail.level_data) {
    return 0
  }
  return props.skillDetail.level_data.length
})

type DescItem = {
  tmpl: string,
  indexes: number[],
}

const skillDesc = computed(() => {
  let tmpl = props.skillDetail.description_template.replaceAll('%%', '%')
  console.log(tmpl)
  const ret: DescItem[] = []
  let dataCount = 0  // 累计数据数量
  tmpl.split('\n').forEach((t) => {
    const indexes: number[] = []  // 当前行匹配到的数据索引
    const newT = t.replace(/<int>|<float\d?>/g, match => {
      const desc = props.skillDetail.description[dataCount]
      indexes.push(dataCount)
      // 每匹配成功, dataCount + 1
      dataCount++
      let data = -1
      if (desc.tp >= 0) {
        data = props.skillDetail.static_data[desc.idx]
      } else {
        if (props.skillDetail.level_data) {
          const levelData = props.skillDetail.level_data[showLevel.value - 1]
          data = levelData[desc.idx]
        }
      }
      data *= desc.rate
      if (!Number.isInteger(desc.rate)) {
        return data.toFixed(String(desc.rate).length - 2)
      }
      return String(data)
    })
    ret.push({
      tmpl: newT,
      indexes,
    })
  })
  return ret
})
</script>

<template>
  <div>
    <div class="flex items-center mt-4">
      <span>lv</span>
      <n-input-number
          v-model:value="showLevel"
          :max="maxLevel"
          :min="1"
          class="w-10 ml-1"
          :show-button="false"
      ></n-input-number>
      <n-slider v-model:value="showLevel" :max="maxLevel" :min="1" class="ml-2"></n-slider>
    </div>
    <div class="mt-2 max-h-[40vh] overflow-auto w-full">
      <div v-for="item in skillDesc" :key="item.tmpl" class="mt-1 flex w-full justify-between">
        <p>{{ item.tmpl }}</p>
        <n-button-group>
          <n-button v-for="(idx, i) in item.indexes" :key="idx" @click="onEdit(idx)">数值{{ i + 1 }}</n-button>
        </n-button-group>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
