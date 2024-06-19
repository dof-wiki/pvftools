<script setup lang="ts">
import {computed, ref} from "vue";
import storage from "../common/storage";
import {CascaderOption} from "naive-ui";

const props = defineProps<{
  withJob?: boolean
  ensureJob?: number
}>()

const emit = defineEmits(['change'])

const selected = ref<string>()

const options = computed(() => {
  const opts: Array<CascaderOption> = []
  storage.getJobMap().forEach((v, k) => {
    if (props.ensureJob !== undefined && props.ensureJob !== k) {
      return
    }
    opts.push({
      value: String(k),
      label: v[0].job_name,
      children: v.map((j) => ({
        value: `${k}-${j.grow_type}`,
        label: j.grow_type_name,
      }))
    })
  })
  return opts
})

const onChange = () => {
  if (!selected.value) {
    emit('change', null, null)
    return
  }
  const args = selected.value.split('-')
  if (args.length === 1) {
    emit('change', null, Number(args[0]))
  } else {
    const job = storage.getJobMap().get(Number(args[0]))?.find((j) => j.grow_type === Number(args[1]))
    if (job) {
      emit('change', job, Number(args[0]))
    }
  }
}
</script>

<template>
  <n-cascader
      v-model:value="selected"
      :options="options"
      :leaf-only="!withJob"
      filterable
      clearable
      @change="onChange"
  ></n-cascader>
</template>

<style scoped>

</style>
