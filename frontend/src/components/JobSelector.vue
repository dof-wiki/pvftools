<script setup lang="ts">
import {computed, ref} from "vue";
import storage from "../common/storage";

const model = defineModel()

const jobOptions = computed(() => {
  const opts: Array<{value: number, name: string}> = []
  storage.getJobMap().forEach((v, k) => {
    opts.push({
      value: k,
      name: v[0].job_name
    })
  })
  return opts
})

</script>

<template>
  <n-select
      v-model="model"
      :options="jobOptions"
      label-field="name"
      style="width: 200px"
  ></n-select>
</template>

<style scoped>

</style>
