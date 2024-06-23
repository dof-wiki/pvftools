<script setup lang="ts">
import {computed, reactive, ref} from "vue";
import * as equAttrDict from "../resources/equ_attr.json"

type ChoiceItem = {
  label: string
  value: string
}

type TAttr = {
  name: string
  tmpl: string
  desc?: string[]
  choices?: Array<string | null>
}

type AttrItem = {
  name: string
  content: string
}

const cate = ref('trigger')
const triggerData = reactive({
  trigger: [] as AttrItem[],
  then: [] as AttrItem[],
})

const attrOptions = computed(() => {
  return (equAttrDict as any)[cate.value] as Array<TAttr>
})

</script>

<template>
  <div class="flex h-full w-full">
    <div class="border-r border-gray-800">
      <n-list clickable hoverable class="h-full">
        <n-list-item @click="cate = 'trigger'">触发器</n-list-item>
      </n-list>
    </div>
    <div class="border-r border-gray-800">
      <n-list clickable hoverable class="h-full">
        <n-list-item
            v-for="(item, idx) in attrOptions"
            @click="cate = 'trigger'"
        >{{ item.name }}</n-list-item>
      </n-list>
    </div>
    <div class="flex-grow ml-2 h-full overflow-auto">
    </div>
  </div>
</template>

<style scoped>

</style>
