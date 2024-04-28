<script setup lang="ts">
import {ref} from "vue";
import {SearchByName} from "../../wailsjs/go/api/App";
import {SelectOption} from 'naive-ui'

const model = defineModel()

const options = ref<SelectOption[]>([])
const loading = ref(false)

const handleSearch = async (query: string) => {
  loading.value = true
  try {
    const result = await SearchByName(query)
    options.value = result.map((item) => {
      return {
        label: item.name,
        value: item.id,
      }
    })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <n-select
      v-model="model"
      filterable
      :options="options"
      :loading="loading"
      clearable
      remote
      tag
      virtual-scroll
      @search="handleSearch"
  ></n-select>
</template>

<style scoped>
</style>
