<script setup lang="ts">
import {model} from "../../wailsjs/go/models";
import Skill = model.Skill;
import {computed, onMounted, ref} from "vue";
import {DataTableColumns, useMessage} from "naive-ui";
import {GetJobInfo, GetJobSkill} from "../../wailsjs/go/api/App";
import Job = model.Job;

const message = useMessage()

const props = withDefaults(defineProps<{
  onSelect: (skill: Skill) => Promise<void>,
  multiSelect: boolean,
}>(), {
  multiSelect: false,
})

const emit = defineEmits(['multiSelect'])

type Option = {
  value: number,
  name: string,
}

const loading = ref(false)
const jobOptions = ref<Option[]>([])
const jobMap = ref<Map<number, string>>(new Map())
const selectedJob = ref<number | null>(null)
const keyword = ref('')
const skills = ref<Skill[]>([])
const isAll = ref(false)

const skillColumns = computed(() => {
  const cols: DataTableColumns<Skill> = []
  if (props.multiSelect) {
    cols.push({
      type: 'selection'
    })
  }
  cols.push(...[
    { title: '代码', key: 'code' },
    { title: '名称', key: 'name' },
  ])
  if (isAll.value) {
    cols.push({
      title: '职业',
      key: 'job',
      render: (skl: Skill) => jobMap.value.get(skl.job)
    })
  }
  return cols
})

const onSelectChange = async (newVal: number) => {
  selectedJob.value = newVal
  isAll.value = false
  loading.value = true
  try {
    skills.value = (await GetJobSkill(newVal)).sort((a, b) => a.code - b.code)
  } catch (e) {
    message.error('加载技能失败')
  } finally {
    loading.value = false
  }
}

const loadAll = async () => {
  loading.value = true
  isAll.value = true
  try {
    skills.value = (await GetJobSkill(-1)).sort((a, b) => {
      if (a.job !== b.job) {
        return a.job - b.job
      }
      return a.code - b.code
    })
  } catch (e) {
    message.error('加载技能失败')
  } finally {
    loading.value = false
  }
}

const listContainer = ref<HTMLDivElement | null>(null)
const tableMaxHeight = ref(0)

const filteredSkills = computed(() => {
  if (!keyword.value) {
    return skills.value
  }
  return skills.value.filter(skl => {
    return skl.name.includes(keyword.value)
  })
})

const rowProps = (row: Skill) => {
  return {
    onClick: async () => {
      await props.onSelect(row)
    }
  }
}

const rowKey = (row: Skill) => {
  return row.id
}

const setTableMaxHeight = () => {
  const height = listContainer.value!.clientHeight
  tableMaxHeight.value = height - 34 - 34 - 58
}

const handleCheck = (rowKeys: number[]) => {
  emit('multiSelect', skills.value.filter(skl => rowKeys.includes(skl.id)))
}

onMounted(async () => {
  window.onresize = setTableMaxHeight
  setTableMaxHeight()
  const jobs  = await GetJobInfo()
  jobs.forEach((item) => {
    if (jobMap.value.has(item.code)) {
      return
    }
    jobOptions.value.push({
      value: item.code,
      name: item.job_name,
    })
    jobMap.value.set(item.code, item.job_name)
  })
})
</script>

<template>
  <div ref="listContainer" class="pr-3">
    <div class="flex">
      <n-button @click="loadAll">加载全职业</n-button>
      <n-select
          v-model:value="selectedJob"
          :options="jobOptions"
          label-field="name"
          clearable
          :on-update:value="onSelectChange"
      ></n-select>
    </div>
    <n-input class="mt-1" v-model:value="keyword"></n-input>
    <n-data-table
        class="mt-1"
        size="small"
        :data="filteredSkills"
        :loading="loading"
        :columns="skillColumns"
        :max-height="tableMaxHeight"
        :min-height="tableMaxHeight"
        virtual-scroll
        :row-props="rowProps"
        :row-key="rowKey"
        row-class-name="skill-item"
        @update:checked-row-keys="handleCheck"
    ></n-data-table>
  </div>
</template>

<style scoped>

</style>
