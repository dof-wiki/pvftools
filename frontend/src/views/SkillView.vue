<script setup lang="ts">
import {computed, reactive, ref} from "vue";
import {api, model} from "../../wailsjs/go/models";
import Skill = model.Skill;
import SkillData = api.SkillData;
import {GetSkillDetail, SetSkillData} from "../../wailsjs/go/api/App";
import {DataTableColumns, useMessage} from "naive-ui";
import SkillList from "../components/SkillList.vue";
import SkillDesc from "../components/SkillDesc.vue";

const message = useMessage()

const selectedSkill = ref<Skill | null>(null)
const skillDetail = ref<SkillData | null>(null)
const detailLoading = ref(false)

const loadDetail = async () => {
  if (!selectedSkill.value) return
  try {
    detailLoading.value = true
    skillDetail.value = await GetSkillDetail(selectedSkill.value!.id)
  } finally {
    detailLoading.value = false
  }
}

const onSkillSelected = async (skl: Skill) => {
  selectedSkill.value = skl
  await loadDetail()
}

const editSkillValueCtx = reactive({
  show: false,
  changeType: 0,  // 0: *倍率, 1: +/-数值, 2: =固定值
  num: 0,
  fixted: false,  // 是否固定以下两列
  tp: 0,  // 0: 静态数据, 1: 等级数据
  col: 0,  // 数据列
})

const onEditSkillValue = (idx?: number) => {
  editSkillValueCtx.num = 0
  editSkillValueCtx.changeType = 0
  if (idx === undefined) {
    editSkillValueCtx.tp = 0
    editSkillValueCtx.col = 0
    editSkillValueCtx.fixted = false
  } else {
    const desc = skillDetail.value!.description[idx]
    editSkillValueCtx.tp = desc.tp >= 0 ? 0 : 1
    editSkillValueCtx.col = desc.idx
    editSkillValueCtx.fixted = true
  }
  editSkillValueCtx.show = true
}

const operateAdd = (val: number) => {
  return val + editSkillValueCtx.num
}

const operateMul = (val: number) => {
  return val * editSkillValueCtx.num
}

const operateSet = (val: number) => {
  return editSkillValueCtx.num
}

const confirmChange = () => {
  if (!skillDetail.value) {
    return
  }
  const op = [operateMul, operateAdd, operateSet][editSkillValueCtx.changeType]
  if (editSkillValueCtx.tp === 0) {
    if (editSkillValueCtx.col >= skillDetail.value!.static_data.length) {
      message.warning('数据列超出范围')
      return
    }
    skillDetail.value.static_data[editSkillValueCtx.col] = op(skillDetail.value.static_data[editSkillValueCtx.col])
  } else {
    if (editSkillValueCtx.col >= skillDetail.value.level_data_num) {
      message.warning('数据列超出范围')
      return
    }
    for (let i = 0; i < skillDetail.value.level_data!.length; i++) {
      skillDetail.value.level_data[i][editSkillValueCtx.col] = op(skillDetail.value.level_data[i][editSkillValueCtx.col])
    }
  }
  editSkillValueCtx.show = false
}

const cancelChange = () => {
  editSkillValueCtx.show = false
}

const saveChange = async (tp: boolean) => {
  if (!skillDetail.value || !selectedSkill.value) {
    return
  }
  detailLoading.value = true
  try {
    await SetSkillData(selectedSkill.value.id, skillDetail.value.static_data, skillDetail.value.level_data, tp)
    message.success('保存成功')
  } finally {
    detailLoading.value = false
  }
}

const staticColumns = computed(() => {
  return skillDetail.value?.static_data.map((_, idx) => {
    return {
      title: String(idx),
      key: `col${idx}`
    }
  })
})

const staticTable = computed(() => {
  const row: Record<string, number> = {}
  if (!skillDetail.value) {
    return []
  }
  for (let i = 0; i < skillDetail.value!.static_data.length; i++) {
    row[`col${i}`] = skillDetail.value!.static_data[i]
  }
  return [row]
})

const levelColumns = computed(() => {
  const cols: DataTableColumns = [{ title: '等级', key: 'level', fixed: 'left'}]
  skillDetail.value?.level_data[0].forEach((_, idx) => {
    cols.push({
      title: String(idx),
      key: `col${idx}`,
    })
  })
  return cols
})

const levelTable = computed(() => {
  const rows: Array<Record<string, number>> = []
  if (!skillDetail.value) {
    return rows
  }
  for (let i = 0; i < skillDetail.value!.level_data.length; i++) {
    const levelData = skillDetail.value!.level_data[i]
    const row: Record<string, number> = {level: i + 1}
    for (let j = 0; j < levelData.length; j++) {
      row[`col${j}`] = levelData[j]
    }
    rows.push(row)
  }
  return rows
})
</script>

<template>
  <div class="flex w-full h-full">
    <skill-list :on-select="onSkillSelected" class="w-1/4 h-full"></skill-list>
    <div class="flex-grow p-3 w-3/4" :loading="detailLoading">
      <div v-if="skillDetail">
        <p class="text-xl font-bold">{{ skillDetail.name }}</p>
        <n-flex class="mt-3">
          <n-button type="primary" @click="saveChange(false)" :loading="detailLoading">保存</n-button>
          <n-popconfirm
              @positive-click="saveChange(true)"
              positive-text="继续"
              negative-text="不了不了"
          >
            <template #trigger>
              <n-button type="primary" :loading="detailLoading">保存并同时修改TP</n-button>
            </template>
            将同时修改对应TP技能(如果能找到的话), 需人工确认TP技能和SP技能数据列一致, 慎用!
          </n-popconfirm>
          <n-button @click="loadDetail" :loading="detailLoading">重新加载</n-button>
          <n-button type="warning" @click="onEditSkillValue(undefined)" :loading="detailLoading">修改数据</n-button>
        </n-flex>

        <skill-desc v-if="skillDetail" :on-edit="onEditSkillValue" :skill-detail="skillDetail"></skill-desc>
        <div v-if="skillDetail && skillDetail.static_data" class="w-full mt-2">
          <span class="font-bold text-lg">静态数据</span>
          <n-data-table
              size="small"
              :data="staticTable"
              :columns="staticColumns"
              :scroll-x="skillDetail.static_data.length > 10 ? 1800 : 0"
          ></n-data-table>
        </div>
        <span class="mt-2 font-bold text-lg" v-if="skillDetail && skillDetail.level_data">等级数据</span>
        <n-data-table
            v-if="skillDetail && skillDetail.level_data"
            :size="'small'"
            :data="levelTable"
            :columns="levelColumns"
            :max-height="skillDetail.description.length >= 10 ? '20vh' : '30vh'"
            :scroll-x="skillDetail.level_data_num > 10 ? 18000 : 0"
            virtual-scroll
        ></n-data-table>
      </div>
    </div>

    <n-modal v-model:show="editSkillValueCtx.show">
      <n-card style="width: 500px">
        <div>
          <n-radio-group v-model:value="editSkillValueCtx.tp" :disabled="editSkillValueCtx.fixted">
            <n-radio :value="0">静态数据</n-radio>
            <n-radio :value="1">等级数据</n-radio>
          </n-radio-group>
          <n-input-number
              :show-button="false"
              v-model:value="editSkillValueCtx.col"
              class="mt-3"
              :disabled="editSkillValueCtx.fixted"
          >
            <template #prefix>
              <span>第</span>
            </template>
            <template #suffix>
              列
            </template>
          </n-input-number>
          <n-radio-group v-model:value="editSkillValueCtx.changeType" class="mt-3">
            <n-radio :value="0">*倍率</n-radio>
            <n-radio :value="1">+/-数值</n-radio>
            <n-radio :value="2">=固定值</n-radio>
          </n-radio-group>
          <n-input-number :show-button="false" v-model:value="editSkillValueCtx.num" class="mt-3"></n-input-number>
        </div>
        <template #footer>
          <n-flex class="float-right">
            <n-button @click="cancelChange">取消</n-button>
            <n-button @click="confirmChange" type="primary">确定</n-button>
          </n-flex>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<style scoped>
</style>
