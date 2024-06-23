<script setup lang="ts">
import {computed, reactive, ref} from "vue";
import {DgnDifficulty} from "../common/consts";
import GrowTypeSelector from "../components/GrowTypeSelector.vue";
import {model} from "../../wailsjs/go/models";
import Job = model.Job;
import storage from "../common/storage";
import {NButton, useMessage} from "naive-ui";
import {SaveFile} from "../../wailsjs/go/api/App";
import AnySearch from "../components/AnySearch.vue";

/*
[monster reward item]
56762 203 -1 10100119 2 80 600
[/monster reward item]
怪物    地图    -1   任务附加怪物掉落   数量   几率    上限
 */

const questDifficultyOptions = ['S', 'A', 'B', 'C', 'M', 'Z', 'Y', 'X', 'W', 'O', 'E', 'F', 'G', 'D', 'U'].map((v) => ({
  label: v,
  value: '`' + v + '`'
}))

const gradeOptions = [
  {label: '每日任务', value: '[daily]'},
  {label: '重复任务', value: '[normaly repeat]'},
  {label: '普通任务', value: '[common unique]'},
  {label: '成就任务', value: '[achievement]'},
  {label: '主线任务', value: '[epic]'},
]

type IntChoiceType = {
  label: string
  value: number
}

type QuestTypeOpt = {
  value: number
  label: string
  type: string
  sub_type?: number
  int_data: string[]
  repeat?: boolean
  choices?: Array<IntChoiceType[] | null>
  comp?: Array<string | null>
}

const typeOptions: QuestTypeOpt[] = [
  {
    value: 1,
    label: '杀怪',
    type: '[hunt monster]',
    sub_type: -1,
    int_data: ['副本id', '难度', '怪物id', '数量'],
    choices: [null, DgnDifficulty, null, null],
    comp: ['select-5'],
    repeat: true
  },
  {
    value: 2,
    label: '通关房间',
    type: '[clear map]',
    sub_type: -1,
    int_data: ['map id']
  },
  {
    value: 3,
    label: '收集物品',
    type: '[seeking]',
    sub_type: -1,
    int_data: ['物品id', '数量'],
    comp: ['select-1'],
    repeat: true
  },
  {
    value: 4,
    label: '对话',
    type: '[meet npc]',
    sub_type: -1,
    int_data: ['npc id'],
    comp: ['select-4'],
  },
  {
    value: 5,
    label: '通关副本',
    type: '[condition under clear]',
    sub_type: 6,
    int_data: ['副本id', '难度'],
    choices: [null, DgnDifficulty],
    comp: ['select-5'],
  },
  {
    value: 6,
    label: '通关房间数量',
    type: '[condition under clear]',
    sub_type: 11,
    int_data: ['副本id', '难度', '房间数量'],
    choices: [null, DgnDifficulty, null],
    comp: ['select-5'],
  },
  {
    value: 7,
    label: '限时通关',
    type: '[condition under clear]',
    sub_type: 0,
    int_data: ['副本id', '难度', '时间/s'],
    choices: [null, DgnDifficulty, null],
    comp: ['select-5'],
  },
  {
    value: 8,
    label: '组队通关',
    type: '[condition under clear]',
    sub_type: 5,
    int_data: ['副本id', '难度', '人数'],
    choices: [null, DgnDifficulty, null],
    comp: ['select-5'],
  },
]

const rewardTypeOptions: Array<QuestTypeOpt> = [
  {
    value: 1,
    label: '道具',
    type: '[item]',
    int_data: ['道具id', '数量'],
    comp: ['select-1'],
    repeat: true,
  },
  {
    value: 2,
    label: '称号',
    type: '[title]',
    int_data: ['称号id', '数量'],
    comp: ['select-2'],
  },
  {
    value: 3,
    label: '转职',
    type: '[grow type]',
    int_data: ['转职id'],
    comp: ['grow_type']
  },
  {
    value: 4,
    label: '觉醒',
    type: '[awakening type]',
    int_data: ['觉醒次数']
  }
]

const formData = reactive({
  grade: '[daily]',
  qp: 0,
  name: '',
  difficulty: 'E',
  npc: 2,
  npc_name: '赛丽亚·克鲁敏',
  complete_npc: null as number | null,
  complete_npc_name: '',
  job: '',
  grow_type: -1,
  level: [1, 120],
  pre_quest: null as number | null,
  depend_message: '',
  condition_message: '',
  solve_message: '',
  tp: 1,
  int_data: [] as Array<number | undefined>,
  repeat: [0],
  reward_type: 1,
  reward_int_data: [] as Array<number | undefined>,
  reward_repeat: [0],
  reward_selection_int_data: [] as Array<number | undefined>,
  reward_selection_repeat: [] as number[],
})

const genResult = reactive({
  visible: false,
  content: '',
  path: '',
})

const generate = () => {
  const now = (new Date().getTime() / 1000).toFixed(0)
  genResult.path = `quest/genquest/${now}.stk`
  const type = typeOptions[formData.tp - 1]
  const rewardType = rewardTypeOptions[formData.reward_type - 1]
  const intData: string[] = []
  const rewardIntData: string[] = []
  for (let i = 0; i < type.int_data.length * formData.repeat.length; i++) {
    intData.push(String(formData.int_data[i] === undefined ? -1 : formData.int_data[i]))
  }
  for (let i = 0; i < rewardType.int_data.length * formData.reward_repeat.length; i++) {
    rewardIntData.push(String(formData.reward_int_data[i] === undefined ? -1 : formData.reward_int_data[i]))
  }
  let rewardSelectionIntDataStr = ''
  if (rewardType.value === 1) {
    const rewardSelectionIntData: string[] = []
    for (let i = 0; i < 2 * formData.reward_selection_repeat.length; i++) {
      rewardSelectionIntData.push(String(formData.reward_selection_int_data[i] === undefined ? -1 : formData.reward_selection_int_data[i]))
    }
    rewardSelectionIntDataStr = `
[reward selection int data]
\t${rewardSelectionIntData.join('\t')}
[/reward selection int data]\n`
  }

  let qp = ''
  if (formData.qp > 0) {
    qp = `
[quest point]
\t${formData.qp}
`
  }

  genResult.content = `#PVF_File

[grade]
\t\`${formData.grade}\`

[difficulty]
\t\`${formData.difficulty}\`

[npc index]
\t${formData.npc}

[complete npc index]
\t${formData.npc || -1}

[job]
\t\`[${formData.job || 'all'}]\`

[grow type]
\t${formData.grow_type}

[level]
\t${formData.level[0]}\t${formData.level[1]}
${qp}
[type]
\t\`${type.type}\`

[sub type]
\t${type.sub_type}

[int data]
\t${intData.join('\t')}
[/int data]

[reward type]
\t\`${rewardType.type}\`

[reward int data]
\t${rewardIntData.join('\t')}
[/reward int data]
${rewardSelectionIntDataStr}
[name]
\t\`${formData.name}\`

[depend message]
\t\`${formData.depend_message}\`

[condition message]
\t\`${formData.condition_message}\`

[solve message]
\t\`${formData.solve_message}\`
`
  genResult.visible = true
}

const onGrowTypeLimitChange = (job: Job | null, jobCode: number | null) => {
  if (!job && !jobCode) {
    formData.job = ''
    formData.grow_type = -1
  }
  if (!job) {
    formData.job = storage.getJobStr(jobCode!)
    formData.grow_type = -1
  } else {
    formData.job = job.job
    formData.grow_type = job.grow_type
  }
}

const intDataList = computed(() => {
  const opt = typeOptions[formData.tp - 1]
  return opt.int_data.map((v, i) => ({
    desc: v,
    choices: opt.choices ? opt.choices[i] : null,
    idx: i,
    comp: opt.comp ? opt.comp[i] : null,
  }))
})

const isRepeat = computed(() => {
  return typeOptions[formData.tp - 1].repeat
})

const addRepeat = () => {
  formData.repeat.push(formData.repeat.length)
}

const rmRepeat = () => {
  if (formData.repeat.length === 1) {
    return
  }
  formData.repeat.pop()
}

const rewardIntDataList = computed(() => {
  const opt = rewardTypeOptions[formData.reward_type - 1]
  return opt.int_data.map((v, i) => ({
    desc: v,
    choices: opt.choices ? opt.choices[i] : null,
    idx: i,
    comp: opt.comp ? opt.comp[i] : null,
  }))
})

const isRewardRepeat = computed(() => {
  return typeOptions[formData.reward_type - 1].repeat
})

const addRewardRepeat = () => {
  formData.reward_repeat.push(formData.reward_repeat.length)
}

const rmRewardRepeat = () => {
  if (formData.reward_repeat.length === 1) {
    return
  }
  formData.reward_repeat.pop()
}

const addRewardSelectionRepeat = () => {
  formData.reward_selection_repeat.push(formData.reward_selection_repeat.length)
}

const rmRewardSelectionRepeat = () => {
  if (formData.reward_selection_repeat.length === 1) {
    return
  }
  formData.reward_selection_repeat.pop()
}

const onRewardTypeChange = () => {
  formData.reward_int_data = []
  formData.reward_repeat = [0]
}

const onTypeChange = () => {
  formData.int_data = []
  formData.repeat = [0]
}

const rewardEnsureJob = computed(() => {
  if (formData.job === '') {
    return -1
  }
  return storage.getJobId(formData.job)
})

const onRewardGrowTypeChange = (idx: number, job: Job | null) => {
  formData.reward_int_data[idx] = job?.grow_type
}

const saveResult = async () => {
  await SaveFile(genResult.path, genResult.content)
}

const message = useMessage()

const copyResult = async () => {
  const clip = await navigator.clipboard
  await clip.writeText(genResult.content)
  message.success('已复制')
}

const selectorCtx = reactive({
  show: false,
  searchType: 1,

  field: '',
  idx: 0,
})

const showItemSelector = (comp: string, field: string, idx: number) => {
  selectorCtx.field = field
  selectorCtx.idx = idx
  selectorCtx.searchType = Number(comp.slice(7))
  selectorCtx.show = true
}

const onItemSelect = (id: number, name: string) => {
  if (selectorCtx.idx === -1) {
    (formData as any)[selectorCtx.field] = id;
    (formData as any)[selectorCtx.field + '_name'] = name;
  } else {
    ((formData as any)[selectorCtx.field] as any)[selectorCtx.idx] = id;
  }
  selectorCtx.show = false
}
</script>

<template>
  <n-split direction="horizontal" class="h-full">
    <template #1>
      <n-form class="pr-5 overflow-auto h-full">
        <n-form-item label="任务名称">
          <n-input v-model:value="formData.name"></n-input>
        </n-form-item>
        <n-form-item label="任务难度">
          <n-radio-group v-model:value="formData.difficulty">
            <n-radio v-for="item in questDifficultyOptions" :key="item.value" :label="item.label" :value="item.value"></n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="任务类型">
          <n-radio-group v-model:value="formData.grade">
            <n-radio v-for="item in gradeOptions" :key="item.value" :label="item.label" :value="item.value"></n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="发布npc">
          <n-input-group class="items-center">
            <n-input-number v-model:value="formData.npc" :show-button="false"></n-input-number>
            <n-button @click="showItemSelector('select-4', 'npc', -1)">搜索</n-button>
            <span class="pl-3">{{ formData.npc_name }}</span>
          </n-input-group>
        </n-form-item>
        <n-form-item label="任务引导npc">
          <n-input-group class="items-center">
            <n-input-number v-model:value="formData.complete_npc" :show-button="false"></n-input-number>
            <n-button @click="showItemSelector('select-4', 'complete_npc', -1)">搜索</n-button>
            <span class="pl-3">{{ formData.complete_npc_name }}</span>
          </n-input-group>
        </n-form-item>
        <n-form-item label="职业, 留空为(all)">
          <grow-type-selector @change="onGrowTypeLimitChange" with-job></grow-type-selector>
        </n-form-item>
        <n-form-item :label="`等级限制 ${formData.level[0]}-${formData.level[1]}`">
          <n-slider v-model:value="formData.level" range :step="1" :min="1" :max="120"></n-slider>
        </n-form-item>
        <n-form-item label="前置任务">
          <n-input-number :show-button="false" v-model:value="formData.pre_quest"></n-input-number>
        </n-form-item>
        <n-form-item label="接受对话">
          <n-input v-model:value="formData.depend_message" type="textarea" autosize></n-input>
        </n-form-item>
        <n-form-item label="任务提示">
          <n-input v-model:value="formData.condition_message" type="textarea" autosize></n-input>
        </n-form-item>
        <n-form-item label="完成对话">
          <n-input v-model:value="formData.solve_message" type="textarea" autosize></n-input>
        </n-form-item>
      </n-form>
    </template>

    <template #2>
      <n-form class="pl-5 overflow-auto h-full">
        <n-form-item style="display: flex; justify-content: flex-end;">
          <n-button type="primary" @click="generate">生成</n-button>
        </n-form-item>
        <n-form-item label="完成条件">
          <n-card>
            <n-select :options="typeOptions" v-model:value="formData.tp" @change="onTypeChange"></n-select>
            <n-divider></n-divider>
            <n-flex>
              <n-form size="small" label-placement="left" label-width="80" v-for="rt in formData.repeat" :key="rt">
                <n-form-item v-for="item in intDataList" :key="item.desc" :label="item.desc">
                  <n-input-group>
                    <n-select v-if="item.choices" :options="item.choices" v-model:value="formData.int_data[rt * intDataList.length + item.idx]"></n-select>
                    <n-input-number v-else v-model:value="formData.int_data[rt * intDataList.length + item.idx]" :show-button="false"></n-input-number>

                    <n-button v-if="item.comp?.startsWith('select-')" @click="showItemSelector(item.comp, 'int_data', rt * intDataList.length + item.idx)">搜索</n-button>
                  </n-input-group>
                </n-form-item>
              </n-form>
            </n-flex>
            <n-button-group v-if="isRepeat">
              <n-button @click="addRepeat" size="small">加一组</n-button>
              <n-button @click="rmRepeat" size="small">减一组</n-button>
            </n-button-group>
          </n-card>
        </n-form-item>
        <n-form-item label="奖励">
          <n-card>
            <n-select :options="rewardTypeOptions" v-model:value="formData.reward_type" @change="onRewardTypeChange"></n-select>
            <n-divider></n-divider>
            <n-flex>
              <n-form size="small" label-placement="left" label-width="80" v-for="rt in formData.reward_repeat" :key="rt">
                <n-form-item v-for="item in rewardIntDataList" :key="item.desc" :label="item.desc">
                  <n-input-group>
                    <n-select v-if="item.choices" :options="item.choices" v-model:value="formData.reward_int_data[rt * rewardIntDataList.length + item.idx]"></n-select>
                    <n-input-number v-else v-model:value="formData.reward_int_data[rt * rewardIntDataList.length + item.idx]" :show-button="false"></n-input-number>

                    <n-button v-if="item.comp?.startsWith('select-')" @click="showItemSelector(item.comp, 'reward_int_data', rt * rewardIntDataList.length + item.idx)">搜索</n-button>
                    <grow-type-selector v-else-if="item.comp === 'grow_type'" :ensure-job="rewardEnsureJob" @change="onRewardGrowTypeChange(rt * rewardIntDataList.length + item.idx, $event)"></grow-type-selector>
                  </n-input-group>
                </n-form-item>
              </n-form>
            </n-flex>
            <n-button-group v-if="isRewardRepeat">
              <n-button @click="addRewardRepeat" size="small">加一组</n-button>
              <n-button @click="rmRewardRepeat" size="small">减一组</n-button>
            </n-button-group>
            <template v-if="formData.reward_type === 1">
              <n-divider></n-divider>
              <p class="mb-5">可选奖励</p>
              <n-flex>
                <n-form size="small" label-placement="left" label-width="80" v-for="rt in formData.reward_selection_repeat" :key="rt">
                  <n-form-item label="道具Id">
                    <n-input-group>
                      <n-input-number v-model:value="formData.reward_selection_int_data[rt * 2]" :min="1" :show-button="false"></n-input-number>
                      <n-button @click="showItemSelector('select-1', 'reward_selection_int_data', rt * 2)">搜索</n-button>
                    </n-input-group>
                  </n-form-item>
                  <n-form-item label="数量">
                    <n-input-number v-model:value="formData.reward_selection_int_data[rt * 2 + 1]" :min="1" :show-button="false"></n-input-number>
                  </n-form-item>
                </n-form>
              </n-flex>
              <n-button-group>
                <n-button @click="addRewardSelectionRepeat" size="small">加一组</n-button>
                <n-button @click="rmRewardSelectionRepeat" size="small">减一组</n-button>
              </n-button-group>
            </template>
            <n-divider></n-divider>
            <n-form label-placement="left" label-width="80" size="small">
              <n-form-item label="QP">
                <n-input-number v-model:value="formData.qp"></n-input-number>
              </n-form-item>
            </n-form>
          </n-card>
        </n-form-item>
      </n-form>

      <n-modal v-model:show="genResult.visible">
        <n-card style="width: 80%" title="生成预览">
          <template #header-extra>
            <n-input v-model:value="genResult.path" style="width: 300px"></n-input>
            <n-button type="primary" @click="saveResult">写入pvf</n-button>
            <n-button @click="copyResult">复制</n-button>
          </template>
          <n-input v-model:value="genResult.content" type="textarea" :rows="30"></n-input>
        </n-card>
      </n-modal>

      <n-drawer v-model:show="selectorCtx.show" width="500" class="p-5">
        <any-search :search-type="selectorCtx.searchType" @select="onItemSelect"></any-search>
      </n-drawer>
    </template>
  </n-split>
</template>

<style scoped>

</style>
