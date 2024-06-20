<script setup lang="ts">
import {ref, h, onMounted, reactive} from 'vue'
import {DataTableColumns, NButton, NInput, NInputGroup, NInputNumber, NSelect, useMessage} from "naive-ui";
import {proto} from "../../wailsjs/go/models";
import UpgradeItem = proto.UpgradeItem;
import {GetItemsName, GetUpgradeData, SaveUpgradeData, SetUpgradeItem} from "../../wailsjs/go/api/App";
import ItemSelector from "../components/ItemSelector.vue";
import AnySearch from "../components/AnySearch.vue";

const loaded = ref(false)
const loading = ref(false)
const tp = ref<number>()
const items = ref<UpgradeItem[]>([])

const loadData = async () => {
  if (!tp.value) {
    return
  }
  loading.value = true
  try {
    items.value = await GetUpgradeData(tp.value)
    await getItemNames()
    loaded.value = true
  } finally {
    loading.value = false
  }
}

const message = useMessage()

const save = async () => {
  if (!tp.value) {
    return
  }
  loading.value = true
  try {
    await SetUpgradeItem(tp.value, items.value)
    await SaveUpgradeData(tp.value)
    message.success('保存成功')
  } finally {
    loading.value = false
  }
}

const itemsInfo = ref<Record<number, string>>({})

const getItemNames = async () => {
  itemsInfo.value = await GetItemsName(Array.from(new Set(items.value.map(item => item.cost_id))))
}

const columns: DataTableColumns<UpgradeItem> = [
  {
    title: '强化等级',
    key: 'level',
  },
  {
    title: '成功率',
    key: 'fail_rate',
    render(row, index) {
      return h(NInputNumber, {
        min: 0,
        max: 100,
        value: Number((100 - row.fail_rate / 1000).toFixed(2)),
        onUpdateValue(v) {
          if (v === null) {
            v = 0
          }
          items.value[index].fail_rate = Math.round((100 - v) * 1000)
        }
      })
    }
  },
  {
    title: '失败操作',
    key: 'fail_op',
    render(row, index) {
      return h(NSelect, {
        value: row.fail_op,
        options: [
          {label: '不处理', value: 1},
          {label: '降级', value: 2},
          {label: '损坏', value: 3},
        ],
        onUpdateValue(v) {
          items.value[index].fail_op = v
          if (v !== 2) {
            items.value[index].fail_op_value = 0
          }
        },
      })
    },
  },
  {
    title: '降级等级',
    key: 'fail_op_value',
    render(row, index) {
      return h(NInputNumber, {
        min: 0,
        max: 31,
        value: row.fail_op_value,
        onUpdateValue(v) {
          items.value[index].fail_op_value = v || 0
        },
        disabled: row.fail_op !== 2
      })
    }
  },
  {
    title: '消耗道具',
    key: 'cost_id',
    render(row, index) {
      return h(NInputGroup, {}, [h(NInputNumber, {
        min: 0,
        value: row.cost_id,
        showButton: false,
        onUpdateValue(v) {
          items.value[index].cost_id = v || 0
        },
      }), h(NInput, {
        value: itemsInfo.value[row.cost_id],
        disabled: false,
        readonly: true,
      })])
    }
  },
  {
    title: '消耗数量',
    key: 'cost_count',
    render(row, index) {
      return h(NInputNumber, {
        min: 0,
        value: row.cost_count,
        onUpdateValue(v) {
          items.value[index].cost_count = v || 0
        }
      })
    }
  }
]

const setCostCtx = reactive({
  show: false,
  costId: null as number | null,
  showSelector: false,
  itemIds: [] as number[],
})

const setCost = () => {
  setCostCtx.costId = null
  setCostCtx.show = true
}

const confirmSetCost = () => {
  if (setCostCtx.costId === null) {
    return
  }
  items.value.forEach(v => {
    v.cost_id = setCostCtx.costId!
  })
  setCostCtx.show = false
  getItemNames()
}

const onSelectItem = (v: number) => {
  setCostCtx.costId = v
}
</script>

<template>
  <div>
    <n-flex align="center">
      <n-radio-group v-model:value="tp" @change="loadData">
        <n-radio :value="1">强化</n-radio>
        <n-radio :value="2">增幅</n-radio>
      </n-radio-group>
      <n-button :disabled="!loaded" :loading="loading" @click="setCost">批量设置道具</n-button>
      <n-button :disabled="!loaded" type="warning" @click="save" :loading="loading">保存</n-button>
    </n-flex>
    <n-data-table :data="items" :loading="loading" :columns="columns" max-height="80vh" class="mt-5"></n-data-table>

    <n-modal v-model:show="setCostCtx.show">
      <n-card style="width: 20%">
        <template #header>
          <span>批量设置道具</span>
        </template>
        <n-input-group>
          <n-input-number class="w-full" v-model:value="setCostCtx.costId" :min="1" :show-button="false"></n-input-number>
          <n-button @click="setCostCtx.showSelector=true">选择</n-button>
        </n-input-group>
        <template #footer>
          <div class="flex justify-end">
            <n-button @click="confirmSetCost" type="primary">确定</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>

    <n-drawer v-model:show="setCostCtx.showSelector" class="p-5" width="500">
      <any-search :search-type="1" @select="onSelectItem"></any-search>
    </n-drawer>
  </div>
</template>

<style scoped>

</style>
