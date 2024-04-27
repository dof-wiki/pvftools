<script setup lang="ts">
import {computed, h, onMounted, reactive, ref} from "vue";
import {api, model, script} from "../../wailsjs/go/models";
import Equipment = model.Equipment;
import {EquipmentBatchHandle, ExportFiles, GetFileContent, QueryEquipments} from "../../wailsjs/go/api/App";
import EquQueryCond = api.EquQueryCond;
import {TableColumn} from "naive-ui/es/data-table/src/interface";
import {DataTableRowKey, NButton, useMessage} from "naive-ui";
import {eq, uniqBy} from "lodash";
import Action = script.Action;

const actionOptions = [
  { label: '设置字段', value: 'set' },
  { label: '增加字段', value: 'add' },
  { label: '删除字段', value: 'del' },
]

const filters = reactive({
  rarity: [] as number[],
  types: [] as string[],
  partSetType: 0,
})

const filters2 = reactive({
  name: '',
  minLevel: [0, 100],
})

const choiceAllRarity = () => {
  if (filters.rarity.length === 8) {
    filters.rarity = []
    return
  }
  filters.rarity = [0, 1, 2, 3, 4, 5, 6, 7]
}

const EquTypes = [
  {
    label: '武器',
    value: 'weapon',
  },
  {
    label: '上衣',
    value: 'coat',
  },
  {
    label: '下装',
    value: 'pants',
  },
  {
    label: '肩甲',
    value: 'shoulder',
  },
  {
    label: '腰带',
    value: 'waist',
  },
  {
    label: '鞋子',
    value: 'shoes',
  },
  {
    label: '项链',
    value: 'amulet',
  },
  {
    label: '手镯',
    value: 'wrist',
  },
  {
    label: '戒指',
    value: 'ring',
  },
  {
    label: '辅助装备',
    value: 'support',
  },
  {
    label: '魔法石',
    value: 'magic stone',
  },
]

const choiceAllTypes = () => {
  if (filters.types.length === EquTypes.length) {
    filters.types = []
    return
  }
  filters.types = EquTypes.map(item => item.value)
}

const equipments = ref<Equipment[]>([])
const loading = ref(false)

const loadEqu = async () => {
  const cond: EquQueryCond = {
    rarity: filters.rarity,
    types: filters.types,
    part_set_type: filters.partSetType,
  }
  loading.value = true
  try {
    equipments.value = await QueryEquipments(cond)
  } finally {
    loading.value = false
  }
}

const filteredEqu = computed(() => {
  return equipments.value.filter(item => {
    if (filters2.name !== '' && item.name.indexOf(filters2.name) === -1) {
      return false
    }
    return item.mini_level >= filters2.minLevel[0] && item.mini_level <= filters2.minLevel[1]
  })
})

const equContent = reactive({
  show: false,
  content: '',
})

const showContent = async (path: string) => {
  equContent.content = await GetFileContent(path)
  equContent.show = true
}

const listContainer = ref<HTMLDivElement | null>(null)
const tableMaxHeight = ref(0)

const setTableMaxHeight = () => {
  const height = listContainer.value!.clientHeight
  tableMaxHeight.value = height - 34 - 34 - 220
}

const selectedCodes = ref<number[]>([])
const selectedCodes2 = ref<number[]>([])
const selectedEqus = ref<Equipment[]>([])

const onChecked = (items: DataTableRowKey[]) => {
  selectedCodes.value = items as number[]
}

const onChecked2 = (items: DataTableRowKey[]) => {
  selectedCodes2.value = items as number[]
}

const removeSelected = () => {
  equipments.value = equipments.value.filter(item => !selectedCodes.value.includes(item.code))
}

const transfer = () => {
  let equs = equipments.value.filter(item => selectedCodes.value.includes(item.code))
  equs = equs.concat(selectedEqus.value)
  selectedEqus.value = uniqBy(equs, (v: Equipment) => v.code)
}

const removeSelected2 = () => {
  selectedEqus.value = selectedEqus.value.filter(item => !selectedCodes2.value.includes(item.code))
}

const message = useMessage()

const exportPaths = async () => {
  loading.value = true
  try {
    await ExportFiles(selectedEqus.value.map(item => item.path))
    message.success('导出成功')
  } finally {
    loading.value = false
  }
}

type IAction = {
  type: string
  key: string
  isClosed: boolean
  args: string
}

const batchCtx = reactive({
  show: false,
  actions: [] as Array<IAction>
})

const startBatch = () => {
  batchCtx.show = true
  batchCtx.actions = []
  addAction(0)
}

const rmAction = (idx: number) => {
  batchCtx.actions.splice(idx, 1)
}

const addAction = (idx: number) => {
  batchCtx.actions.splice(idx + 1, 0, {type: 'set', key: 'name', args: '', isClosed: false})
}

const submitBatch = async () => {
  const actions: Action[] = batchCtx.actions.map((item) => {
    const args: string[] = [item.key, item.isClosed ? 'true' : 'false']
    if (item.args && item.args !== '') {
      args.push(item.args)
    }
    return {
      type: item.type,
      args: args,
    }
  })
  await EquipmentBatchHandle(selectedEqus.value.map(item => item.code), actions)
  message.success('提交成功')
  batchCtx.show = false
}

const rarityClassNames = ['white', 'blue', 'purple', 'pink', 'epic', 'hero', 'legend', 'myth']

const rowClassName = (row: Equipment) => {
  return rarityClassNames[row.rarity]
}

const quickAttachType = ref('free')
const quickAttachTypeOpts = [
  { label: '免费', value: 'free' },
  { label: '封装', value: 'sealing' },
  { label: '不可交易', value: 'trade' },
  { label: '帐号绑定', value: 'account' },
  { label: '不可交易、删除', value: 'trade delete' },
  { label: '封装且不可交易', value: 'sealing trade' },
]

const quickSetAttachType = () => {
  batchCtx.actions.push({
    type: 'set',
    key: 'attach type',
    args: `\`[${quickAttachType.value}]\``,
    isClosed: false,
  })
}

const transferRow = (code: number) => {
  const equ = equipments.value.find(item => item.code === code)
  if (equ) {
    selectedEqus.value.push(equ)
    selectedEqus.value = uniqBy(selectedEqus.value, (v: Equipment) => v.code)
  }
}

const rmSelectedRow = (code: number) => {
  selectedEqus.value = selectedEqus.value.filter(item => item.code !== code)
}

const quickSetAction = (key: string) => {
  let akey = key
  let value = ''
  let isClosed = false
  switch (key) {
    case 'creation rate': {
      value = '200'
      break
    }
    case 'random option': {
      value ='1'
      break
    }
    case 'breath': {
      batchCtx.actions.push({
        type: 'breath',
        key: '',
        args: '',
        isClosed: false,
      })
      return
    }
  }
  batchCtx.actions.push({
    type: 'set',
    key: akey,
    args: value,
    isClosed: false,
  })
}

const equColumns: TableColumn[] = [
  { type: 'selection' },
  { title: '代码', key: 'code', sorter: 'default', width: '120px' },
  {
    title: '名称',
    key: 'name',
    sorter: 'default',
  },
  { title: '使用等级', key: 'mini_level', sorter: 'default', width: '80px' },
  {
    title: '',
    key: 'actions',
    render: (row: any) => {
      return [
          h(
            NButton,
            {
              size: 'small',
              onClick: () => showContent(row.path)
            },
            { default: () => '查看' }
          ),
          h(
              NButton,
              {
                size: 'small',
                onClick: () => transferRow(row.code)
              },
              { default: () => '移入' }
          )
      ]
    }
  }
]

const equColumns2: TableColumn[] = [
  { type: 'selection' },
  { title: '代码', key: 'code', sorter: 'default', width: '100px' },
  {
    title: '名称',
    key: 'name',
    sorter: 'default',
  },
  {
    title: '',
    key: 'actions',
    width: '80px',
    render: (row: any) => {
      return h(
          NButton,
          {
            size: 'small',
            onClick: () => rmSelectedRow(row.code)
          },
          { default: () => '移出' }
      )
    }
  }
]

onMounted(() => {
  setTableMaxHeight()
  window.addEventListener('resize', setTableMaxHeight)
  choiceAllTypes()
})
</script>

<template>
  <div class="flex w-full h-full">
    <div ref="listContainer" class="pr-3">
      <div>
        <n-checkbox-group v-model:value="filters.rarity">
          <n-space item-style="display: flex; align-items: center;">
            <n-button @click="choiceAllRarity" size="small">全选</n-button>
            <n-checkbox :value="0" label="白装" />
            <n-checkbox :value="1" label="蓝装" />
            <n-checkbox :value="2" label="紫装" />
            <n-checkbox :value="3" label="粉" />
            <n-checkbox :value="4" label="史诗" />
            <n-checkbox :value="5" label="勇者" />
            <n-checkbox :value="6" label="传说" />
            <n-checkbox :value="7" label="神话" />
          </n-space>
        </n-checkbox-group>
        <n-checkbox-group v-model:value="filters.types">
          <n-space item-style="display: flex; align-itemss: center;">
            <n-button @click="choiceAllTypes" size="small">全选</n-button>
            <n-checkbox v-for="item in EquTypes" :key="item.value" :value="item.value">{{ item.label }}</n-checkbox>
          </n-space>
        </n-checkbox-group>
        <n-radio-group v-model:value="filters.partSetType" class="mt-2">
          <n-space item-style="display: flex; align-items: center;">
            <n-radio :value="0">全部</n-radio>
            <n-radio :value="1">套装</n-radio>
            <n-radio :value="2">散件</n-radio>
            <span>注：此筛选只适用于新版套装</span>
          </n-space>
        </n-radio-group>
        <n-space class="mt-2" item-style="display: flex; align-items: center;">
          <n-button @click="loadEqu" type="primary" :loading="loading">加载</n-button>
          <span>共 {{ filteredEqu.length }} 件</span>
        </n-space>
        <div class="flex w-[30vw] mt-2 items-center" v-show="equipments.length > 0">
          <n-input v-model:value="filters2.name" placeholder="输入关键字筛选"></n-input>
          <span class="mx-2">Lv</span>
          <n-slider v-model:value="filters2.minLevel" range></n-slider>
        </div>
      </div>
      <div class="flex items-center mt-2">
        <div>
          <n-button @click="removeSelected">移除选中</n-button>
          <n-data-table
              class="mt-2"
              size="small"
              :data="filteredEqu"
              :loading="loading"
              :columns="equColumns"
              :max-height="tableMaxHeight"
              :min-height="tableMaxHeight"
              virtual-scroll
              :row-class-name="rowClassName"
              :row-key="(row: Equipment) => row.code"
              @update:checked-row-keys="onChecked"
          ></n-data-table>
        </div>
        <n-space class="mx-2">
          <n-button type="primary" @click="transfer">-></n-button>
          <n-button @click="removeSelected2"><-</n-button>
        </n-space>
        <div>
          <n-space>
            <n-button type="primary" @click="exportPaths" :loading="loading">导出路径</n-button>
            <n-button type="primary" @click="startBatch" :loading="loading">批处理</n-button>
          </n-space>
          <n-data-table
              class="mt-2"
              style="width: 420px"
              size="small"
              :data="selectedEqus"
              :loading="loading"
              :columns="equColumns2"
              :max-height="tableMaxHeight"
              :min-height="tableMaxHeight"
              virtual-scroll
              :row-class-name="rowClassName"
              :row-key="(row: Equipment) => row.code"
              @update:checked-row-keys="onChecked2"
          ></n-data-table>
        </div>
      </div>
    </div>
    <n-modal v-model:show="equContent.show">
      <n-card style="width: 80%">
        <n-input type="textarea" v-model:value="equContent.content" :rows="35"></n-input>
      </n-card>
    </n-modal>
    <n-modal v-model:show="batchCtx.show" :mask-closable="false" :close-on-esc="false">
      <n-card style="width: 80%">
        <template #header>
          <span>已选择 {{ selectedEqus.length }} 件装备</span>
        </template>
        <div>
          <p class="text-lg mb-1">脚本</p>
          <div v-for="(act, idx) in batchCtx.actions" :key="idx">
            <n-form inline>
              <n-form-item label="操作">
                <n-select style="width: 120px;" v-model:value="act.type" :options="actionOptions"></n-select>
              </n-form-item>
              <n-form-item label="键 (不要带中括号)">
                <n-input v-model:value="act.key"></n-input>
              </n-form-item>
              <n-form-item label="值">
                <n-input v-model:value="act.args"></n-input>
              </n-form-item>
              <n-form-item>
                <n-button-group>
                  <n-button @click="rmAction(idx)">-</n-button>
                  <n-button @click="addAction(idx)">+</n-button>
                </n-button-group>
              </n-form-item>
            </n-form>
          </div>
        </div>
        <div>
          <p class="mb-1 text-lg">预设</p>
          <n-flex align="center">
            <n-button size="small" @click="quickSetAction('creation rate')">掉率</n-button>
            <n-button size="small" @click="quickSetAction('random option')">设置魔法封印</n-button>
            <div class="flex items-center">
              <n-select size="small" v-model:value="quickAttachType" :options="quickAttachTypeOpts" style="width: 150px;"></n-select>
              <n-button size="small" @click="quickSetAttachType()">交易类型</n-button>
            </div>
            <n-button size="small" @click="quickSetAction('breath')">异次元气息</n-button>
          </n-flex>
        </div>
        <template #footer>
          <n-space>
            <n-button @click="submitBatch" type="primary">提交</n-button>
            <n-button @click="batchCtx.show = false">取消</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<style scoped>
:deep(.white td) {
  color: white !important;
}

:deep(.blue td) {
  color: blue !important;
}

:deep(.purple td) {
  color: rebeccapurple !important;
}

:deep(.pink td) {
  color: #ff00ff !important;
}

:deep(.epic td) {
  color: gold !important;
}

:deep(.hero td) {
  color: indianred !important;
}

:deep(.legend td) {
  color: #ff8400 !important;
}

:deep(.myth td) {
  color: #4d941e !important;
}
</style>