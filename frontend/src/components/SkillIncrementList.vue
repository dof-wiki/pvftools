<script setup lang="ts">
import {api, model, proto} from "../../wailsjs/go/models";
import SkillIncrement = proto.SkillIncrement;
import {computed, onMounted, reactive, ref} from "vue";
import storage from "../common/storage";
import SkillList from "./SkillList.vue";
import Skill = model.Skill;
import SkillDesc from "./SkillDesc.vue";
import SkillData = api.SkillData;
import {GetSkillDetail} from "../../wailsjs/go/api/App";

const props = defineProps<{
  data: Array<SkillIncrement>,
}>()

const emit = defineEmits(['update'])

const genExplain = (item: SkillIncrement) => {
  const jobStr = item.job.slice(1, item.job.length - 1)
  const jobName = storage.getJobNameByStr(jobStr)
  const skl = storage.getSkill(storage.getJobId(jobStr)!, item.skill_id)
  const skillName = skl ? skl.name : String(item.skill_id)
  let range = 'pk场'
  switch (item.type) {
    case '[dungeon type]':
      range = '地下城'
      break
    case '[all]':
      range = '全局'
      break
  }
  let dataType = item.data_type
  let useIdx = false
  switch (item.data_type) {
    case '[level]':
      dataType = '等级数据'
      useIdx = true
      break
    case '[static]':
      dataType = '静态数据'
      useIdx = true
      break
    case '[cooltime]':
      dataType = '冷却时间'
      break
    // TODO: 消耗
  }
  let c = `(${range})${jobName}-[${skillName}]${dataType}`
  if (useIdx) {
    c += `第${item.data_index}列`
  }
  if (item.value > 0) {
    c += ` +${item.value}`
  } else {
    c += ` ${item.value}`
  }
  if (item.increment_type === '%') {
    c += item.increment_type
  }
  return c
}

const editCtx = reactive({
  show: false,
  selectSkillShow: false,
  skillDetailShow: false,
  idx: -1,
  formData: {
    job: '[swordman]',
    skill_id: 0,
    type: '[dungeon type]',
    data_type: '[level]',
    data_index: 0,
    increment_type: '%',
    value: 0,
  },
  skill: undefined as Skill | undefined,
})

const addRow = () => {
  editCtx.idx = -1
  if (props.data.length > 0) {
    const row = props.data[0]
    const jobStr = row.job.slice(1, row.job.length - 1)
    editCtx.skill = storage.getSkill(storage.getJobId(jobStr)!, row.skill_id)
    editCtx.formData.skill_id = row.skill_id
    editCtx.formData.job = row.job
  }
  editCtx.show = true
}

const editRow = (idx: number) => {
  editCtx.idx = idx
  const row = props.data[idx]
  editCtx.formData.job = row.job
  editCtx.formData.skill_id = row.skill_id
  editCtx.formData.type = row.type
  editCtx.formData.data_type = row.data_type
  editCtx.formData.data_index = row.data_index
  editCtx.formData.increment_type = row.increment_type
  editCtx.formData.value = row.value
  editCtx.show = true
  const jobStr = row.job.slice(1, row.job.length - 1)
  editCtx.skill = storage.getSkill(storage.getJobId(jobStr)!, row.skill_id)
}

const deleteRow = (idx: number) => {
  const newData = props.data.filter((v, i) => i !== idx)
  emit('update', newData)
}

const confirmEdit = () => {
  if (editCtx.idx === -1) {
    if (!editCtx.skill) {
      return
    }
    const newData = props.data.concat([new SkillIncrement(editCtx.formData)])
    emit('update', newData)
  } else {
    const newData = props.data.map((item, i) => {
      if (i === editCtx.idx) {
        return new SkillIncrement(editCtx.formData)
      } else {
        return item
      }
    })
    emit('update', newData)
  }
  editCtx.show = false
}

const onSkillSelected = async (skl: Skill) => {
  editCtx.formData.job = `[${storage.getJobStr(skl.job)}]`
  editCtx.formData.skill_id = skl.code
  editCtx.skill = skl
  editCtx.selectSkillShow = false
}

const onSkillEdit = (idx?: number) => {
  if (idx === undefined) {
    return
  }
  const desc = skillDetail.value!.description[idx]
  editCtx.formData.data_type = desc.tp >= 0 ? '[static]' : '[level]'
  editCtx.formData.data_index = desc.idx
  editCtx.skillDetailShow = false
}
const skillDetail = ref<SkillData | null>(null)

const loadDetail = async () => {
  if (!editCtx.skill) {
    return
  }
  skillDetail.value = await GetSkillDetail(editCtx.skill.id)
  editCtx.skillDetailShow = true
}

onMounted(async () => {
  await Promise.all([
    storage.loadJobMap(),
    storage.loadSkillMap(),
  ])
})
</script>

<template>
  <n-list>
    <n-list-item v-for="(item, idx) in data" :key="idx">
      <div class="flex justify-between items-center">
        <p>
          {{ genExplain(item) }}
        </p>
        <div>
          <n-button-group size="small">
            <n-button @click="editRow(idx)">编辑</n-button>
            <n-button type="error" @click="deleteRow(idx)">删除</n-button>
          </n-button-group>
        </div>
      </div>
    </n-list-item>
    <n-list-item>
      <div class="flex justify-end">
        <n-button type="info" size="small" @click="addRow">增加</n-button>
      </div>
    </n-list-item>
  </n-list>
  <n-modal v-model:show="editCtx.show">
    <n-card style="width: 40%">
      <n-form label-placement="left" label-width="80px">
        <n-form-item label="技能">
          <span class="text-blue-300 mr-2">{{ editCtx.skill?.name }}</span>
          <n-button size="small" @click="editCtx.selectSkillShow = true">选择</n-button>
        </n-form-item>
        <n-form-item label="生效范围">
          <n-radio-group v-model:value="editCtx.formData.type">
            <n-radio value="[dungeon type]">地下城</n-radio>
            <n-radio value="[pvp type]">pk场</n-radio>
            <n-radio value="[all]">全局</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="数据类型">
          <n-radio-group v-model:value="editCtx.formData.data_type">
            <n-radio value="[level]">等级数据</n-radio>
            <n-radio value="[static]">静态数据</n-radio>
            <n-radio value="[cooltime]">冷却时间</n-radio>
            <n-radio value="[mp]">蓝耗</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="数据索引" v-if="editCtx.formData.data_type === '[level]' || editCtx.formData.data_type === '[static]'">
          <n-input-number v-model:value="editCtx.formData.data_index" :min="0" :show-button="false"></n-input-number>
          <n-button v-if="editCtx.skill" class="ml-2" @click="loadDetail">选择</n-button>
        </n-form-item>
        <n-form-item label="修正类型">
          <n-radio-group v-model:value="editCtx.formData.increment_type">
            <n-radio value="+">+</n-radio>
            <n-radio value="%">%</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="值">
          <n-input-number v-model:value="editCtx.formData.value" :show-button="false"></n-input-number>
        </n-form-item>
        <n-form-item>
          <div class="flex justify-end">
            <n-button type="primary" @click="confirmEdit">确定</n-button>
          </div>
        </n-form-item>
      </n-form>
    </n-card>
  </n-modal>
  <n-drawer v-model:show="editCtx.selectSkillShow" :width="300">
    <skill-list class="h-screen" :on-select="onSkillSelected"></skill-list>
  </n-drawer>
  <n-drawer v-model:show="editCtx.skillDetailShow" :width="500">
    <skill-desc class="p-5" v-if="skillDetail" :on-edit="onSkillEdit" :skill-detail="skillDetail"></skill-desc>
  </n-drawer>
</template>

<style scoped>

</style>
