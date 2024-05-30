<script setup lang="ts">
import {computed, onMounted, reactive, ref, watch, h} from "vue";
import {DataTableColumns, NButton, useDialog, useMessage} from 'naive-ui';
import {GetBreathData, ReloadBreath, SaveBreath, SetBreathBaseData, SetBreathSkills} from "../../wailsjs/go/api/App";
import {model, proto} from "../../wailsjs/go/models";
import {EquTypes} from "../common/consts";
import storage from "../common/storage";

const message = useMessage()

import Job = model.Job;
import BreathData = proto.BreathData;
import BreathCheckSkillItem = proto.BreathCheckSkillItem;
import SkillIncrement = proto.SkillIncrement;
import SkillIncrementList from "../components/SkillIncrementList.vue";
import BreathCheckItem = proto.BreathCheckItem;

const filters = reactive({
  breathId: 0,
  job: 0,
  subJob: 0,
  equType: 'weapon',
})

type JobOpt = {
  value: number,
  name: string,
  growTypes: Job[],
}

const jobOptions = ref<Array<JobOpt>>([])

const subJobOptions = computed<Array<Job>>(() => {
  const job = jobOptions.value.find((j) => j.value === filters.job)
  return job ? job.growTypes : []
})

const getGroupTypeOpts = (items: Job[]) => {
  return items.filter((item) => item.grow_type <= 4).map((item) => {
    if (item.grow_type === 0) {
      const ex = items.find((v) => v.grow_type === 5)
      if (ex) {
        item.grow_type_name = ex.grow_type_name
      }
    }
    return item
  })
}

const loadedData = reactive({
  loadedId: 0,
  data: null as BreathData | null,
})

const dialog = useDialog()

const reloadData = async () => {
  loadedData.data = await GetBreathData(filters.breathId)
  loadedData.loadedId = filters.breathId
}

const onDataChange = async () => {
  if (!loadedData.data) {
    return
  }
  loadedData.data.has_change = true
}

const calcCheckList = () => {
  if (!loadedData.data) {
    return []
  }
  const item = loadedData.data.check_list.find((it) => it.job === filters.job && it.sub_job === filters.subJob && it.equ_type === filters.equType)
  return item ? item.skill_list : []
}

const checkList = computed(() => {
  if (!loadedData.data) {
    return []
  }
  const item = loadedData.data.check_list.find((it) => it.job === filters.job && it.sub_job === filters.subJob && it.equ_type === filters.equType)
  return item ? item.skill_list : []
})

const onCheckListUpdate = async () => {
  await SetBreathSkills(filters.breathId, filters.job, filters.subJob, filters.equType, calcCheckList())
  loadedData.data!.has_change = true
}

const addCheckSkill = async () => {
  if (!loadedData.data) {
    return
  }
  const item = loadedData.data.check_list.find((it) => it.job === filters.job && it.sub_job === filters.subJob && it.equ_type === filters.equType)
  if (!item) {
    loadedData.data.check_list.push(new BreathCheckItem({
      job: filters.job,
      sub_job: filters.subJob,
      equ_type: filters.equType,
      skill_list: [new BreathCheckSkillItem({
        index: 0,
        increment_list: [],
        explain: ''
      })],
    }))
  } else {
    item.skill_list.push(new BreathCheckSkillItem({
      index: item.skill_list.length,
      increment_list: [],
      explain: ''
    }))
  }
  await onCheckListUpdate()
}

const rmCheckSkill = async (row: BreathCheckSkillItem) => {
  if (!loadedData.data) {
    return
  }
  const item = loadedData.data.check_list.find((it) => it.job === filters.job && it.sub_job === filters.subJob && it.equ_type === filters.equType)
  if (!item) {
    return
  }
  item.skill_list = item.skill_list.filter((it) => it.index !== row.index)
  item.skill_list = item.skill_list.map((it, idx) => {
    it.index = idx
    return it
  })
  await onCheckListUpdate()
}

const skillDetailCtx = reactive({
  row: null as BreathCheckSkillItem | null,
  show: false,
  explainShow: false,
})

const onSkillIncrUpdate = async (newData: Array<SkillIncrement>) => {
  skillDetailCtx.row!.increment_list = newData
  await onCheckListUpdate()
}

const checkTableColumns: DataTableColumns<BreathCheckSkillItem> = [
  {
    title: '索引',
    key: 'index',
  },
  {
    title: '描述',
    key: 'explain',
    render: (rowData) => {
      if (rowData.increment_list.length === 0) {
        return rowData.explain
      }
      const skl = storage.getSkill(filters.job, rowData.increment_list[0].skill_id)
      return `[${skl?.name}]${rowData.explain}`
    },
  },
  {
    title: '',
    key: 'actions',
    render: (row) => {
      return [
        h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              onClick: () => {
                skillDetailCtx.row = row
                skillDetailCtx.show = true
              }
            },
            { default: () => '数据' },
        ),
        h(
            NButton,
            {
              size: 'small',
              type: 'info',
              style: {
                marginLeft: '5px'
              },
              onClick: () => {
                skillDetailCtx.row = row
                skillDetailCtx.explainShow = true
              }
            },
            { default: () => '解释' },
        ),
        h(
            NButton,
            {
              size: 'small',
              type: 'error',
              style: {
                marginLeft: '5px'
              },
              onClick: () => {
                rmCheckSkill(row)
              }
            },
            { default: () => '删除' },
        )
      ]
    }
  }
]

watch(filters, async (v) => {
  if (v.breathId === loadedData.loadedId) {
    return
  }
  await reloadData()
})

const loadJobInfo = async () => {
  await storage.loadJobMap()
  jobOptions.value = []
  storage.getJobMap().forEach((v, k) => {
    jobOptions.value.push({
      value: k,
      name: v[0].job_name,
      growTypes: getGroupTypeOpts(v),
    })
  })
}

const save = async () => {
  await SetBreathBaseData(loadedData.loadedId, loadedData.data!.probability1, loadedData.data!.probability2)
  await SaveBreath(loadedData.loadedId)
  message.success('已保存')
  loadedData.data!.has_change = false
}

const reload = async () => {
  if (loadedData.data?.has_change) {
    dialog.warning({
      title: '警告',
      content: '当前有修改未保存, 重载将丢失未保存的数据, 是否继续?',
      positiveText: '忽略并继续',
      negativeText: '取消',
      onPositiveClick(e) {
        ReloadBreath(filters.breathId).then((resp) => {
          loadedData.data = resp
          loadedData.loadedId = filters.breathId
        })
      },
      onNegativeClick(e) {
        filters.breathId = loadedData.loadedId
      },
    })
  } else {
    loadedData.data = await ReloadBreath(filters.breathId)
  }
}

onMounted(async () => {
  await Promise.all([
      loadJobInfo(),
      storage.loadSkillMap(),
  ])
})
</script>

<template>
  <div>
    <n-flex>
      <n-button :disabled="!loadedData.data?.has_change" type="warning" @click="save">保存</n-button>
      <n-button :disabled="!loadedData.data" type="info" @click="reload">重载</n-button>
    </n-flex>
    <n-form label-placement="left" label-width="120" class="mt-5">
      <n-form-item label="气息">
        <n-radio-group v-model:value="filters.breathId">
          <n-radio :value="1254">红气息</n-radio>
          <n-radio :value="1255">蓝气息</n-radio>
          <n-radio :value="1256">绿气息</n-radio>
          <n-radio :value="1257">黄气息</n-radio>
        </n-radio-group>
      </n-form-item>
      <template v-if="loadedData.data">
        <n-form-item label="成功率">
          <n-input-number
              :show-button="false"
              :min="0"
              :max="100"
              v-model:value="loadedData.data.probability1"
              @update-value="onDataChange"
          >
            <template #prefix>首次</template>
            <template #suffix>%</template>
          </n-input-number>
          <n-input-number
              class="ml-2"
              :show-button="false"
              :min="0"
              :max="100"
              v-model:value="loadedData.data.probability2"
              @update-value="onDataChange"
          >
            <template #prefix>第2次</template>
            <template #suffix>%</template>
          </n-input-number>
        </n-form-item>
        <n-divider></n-divider>
        <n-form-item label="职业">
          <n-select
              v-model:value="filters.job"
              :options="jobOptions"
              label-field="name"
              style="width: 200px"
          ></n-select>
          <n-select
              v-model:value="filters.subJob"
              :options="subJobOptions"
              label-field="grow_type_name"
              value-field="grow_type"
              style="width: 200px"
          ></n-select>
        </n-form-item>
        <n-form-item label="装备类型">
          <n-radio-group v-model:value="filters.equType">
            <n-radio v-for="item in EquTypes" :value="item.value" :key="item.value">{{ item.label }}</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="数据">
          <div class="w-full">
            <n-flex class="mt-2">
              <n-button size="small" type="info" @click="addCheckSkill">增加</n-button>
            </n-flex>
            <n-data-table
                class="mt-2"
                size="small"
                striped
                :data="checkList"
                :columns="checkTableColumns"
                :row-key="(row: BreathCheckSkillItem) => row.index"
            ></n-data-table>
          </div>
        </n-form-item>
      </template>
    </n-form>

    <n-modal v-model:show="skillDetailCtx.show">
      <n-card style="width: 50%">
        <skill-increment-list
            v-if="skillDetailCtx.row"
            :data="skillDetailCtx.row.increment_list"
            @update="onSkillIncrUpdate"
        ></skill-increment-list>
      </n-card>
    </n-modal>

    <n-modal v-model:show="skillDetailCtx.explainShow">
      <n-card style="width: 40%">
        <n-input v-if="skillDetailCtx.row" v-model:value="skillDetailCtx.row.explain" @updateValue="onCheckListUpdate"></n-input>
      </n-card>
    </n-modal>
  </div>
</template>

<style scoped>

</style>
