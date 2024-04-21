<script setup lang="ts">
import {useMessage} from "naive-ui";
import {computed, reactive, ref, watch} from "vue";
import {world_drop} from "../../wailsjs/go/models";
import WorldDropItem = world_drop.WorldDropItem;
import {
  AddWorldDropItem,
  DelWorldDropItem,
  GetItemsName,
  GetWorldDropData,
  SaveWorldDropData
} from "../../wailsjs/go/api/App";

const message = useMessage()

const columns = [
  { title: 'id', key: 'id' },
  { title: '掉落物', key: 'name' },
  { title: '爆率', key: 'rate' },
]

const loading = ref(false)
const show0 = ref(false)

const sdata = ref<WorldDropItem[]>([])
const itemsInfo = ref<Record<number, string>>({})
const data = computed(() => {
  return sdata.value.map((item: WorldDropItem) => {
    return {
      id: item.id,
      name: itemsInfo.value[item.id],
      rate: item.rate,
    }
  }).filter((item) => item.rate > 0 || show0.value).sort((a, b) => a.id - b.id)
})

const selectedLevel = ref(1)
const rowKey = (row: WorldDropItem) => row.id

const getItemNames = async () => {
  itemsInfo.value = await GetItemsName(Array.from(new Set(sdata.value.map(item => item.id))))
}

const getData = async () => {
  loading.value = true
  try {
    sdata.value = await GetWorldDropData(selectedLevel.value)
    await getItemNames()
  } finally {
    loading.value = false
  }
}

const formLoading = ref(false)
const formDialogOpened = ref(false)
const isAdd = ref(false)

const dropItem = reactive({
  levelRange: [1, 200],
  code: 0,
  rate: 500,
})

const onAddDropItem = (add: boolean) => {
  dropItem.levelRange = [1, 200]
  dropItem.rate = 500
  isAdd.value = add
  formDialogOpened.value = true
}

const delRate0 = async () => {
}

const addDropItem = async () => {
  formLoading.value = true
  try {
    if (isAdd.value) {
      await AddWorldDropItem(dropItem.levelRange[0], dropItem.levelRange[1], dropItem.code, dropItem.rate)
      message.success('添加成功')
    } else {
      await DelWorldDropItem(dropItem.levelRange[0], dropItem.levelRange[1], [dropItem.code])
      message.success('删除成功')
    }
    formDialogOpened.value = false
    await getData()
  } finally {
    formLoading.value = false
  }
}

const saveLoading = ref(false)

const saveData = async () => {
  saveLoading.value = true
  try {
    await SaveWorldDropData()
    message.success('保存成功')
  } finally {
    saveLoading.value = false
  }
}

watch(selectedLevel, getData)
</script>

<template>
  <div>
    <n-space class="flex items-center">
      <n-input-number
          v-model:value="selectedLevel"
          class="w-[120px] mx-2 text-center"
          button-placement="both"
          :max="200"
          :min="1"
      />
      <n-button
          type="primary"
          :loading="loading"
          @click="getData"
      >
        加载掉落数据
      </n-button>
      <n-switch v-model:value="show0">
        <template #checked>
          显示掉落率为0的物品
        </template>
        <template #unchecked>
          隐藏掉落率为0的物品
        </template>
      </n-switch>
      <n-button @click="onAddDropItem(true)">
        增加掉落项
      </n-button>
      <n-button @click="onAddDropItem(false)">
        删除掉落项
      </n-button>
<!--      <n-button @click="delRate0">删除掉落率为0的项</n-button>-->
      <n-button
          :loading="saveLoading"
          type="warning"
          @click="saveData"
      >
        保存数据
      </n-button>
    </n-space>
    <div class="mt-[20px] overflow-auto h-[80vh]">
      <n-data-table
          ref="table"
          striped
          size="small"
          remote
          :columns="columns"
          :data="data"
          :loading="loading"
          :row-key="rowKey"
      />
    </div>
    <n-modal v-model:show="formDialogOpened">
      <n-card
          title="增加掉落项"
          style="width: 60%"
      >
        <n-form>
          <n-form-item :label="`掉落等级: ${dropItem.levelRange[0]} - ${dropItem.levelRange[1]}`">
            <n-slider
                v-model:value="dropItem.levelRange"
                range
                :step="1"
                :min="1"
                :max="200"
            />
          </n-form-item>
          <n-form-item label="掉落物品代码">
            <n-input-number
                v-model:value="dropItem.code"
                :min="1"
                :show-button="false"
            />
          </n-form-item>
          <n-form-item
              v-if="isAdd"
              label="掉落几率"
          >
            <n-input-number
                v-model:value="dropItem.rate"
                :min="0"
                :show-button="false"
            />
          </n-form-item>
        </n-form>
        <template #footer>
          <div>
            <n-button
                type="primary"
                :loading="formLoading"
                @click="addDropItem"
            >
              确定
            </n-button>
          </div>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<style scoped>

</style>