<script setup lang="ts">
import {computed, h, onMounted, reactive, ref} from "vue";
import {box} from "../../wailsjs/go/models";
import BoxItemGroup = box.BoxItemGroup;
import {NButton, useMessage, useModal } from "naive-ui";
import {GenBox, SaveFile} from "../../wailsjs/go/api/App";
import BoxParams = box.BoxParams;
import RarityRadio from "../components/RarityRadio.vue";
import InputModal from "../components/InputModal.vue";
import ItemSelector from "../components/ItemSelector.vue";

const formData = reactive({
  box_type: 0,
  name: '',
  explain: '',
  rarity: 0,
  grade: 1,
  items: [] as BoxItemGroup[],
  tp: 'etc',
  category_names: ['请选择'] as string[],
})

const isSelection = computed(() => formData.box_type === 1)

type TableData = {
  id: number
  count: number
  rate: number
  idx: number
}

const message = useMessage()

const selectedIdx = ref(0)

const tableData = computed(() => {
  if (formData.items.length === 0) {
    return []
  }
  return formData.items[selectedIdx.value].items.map((item, idx) => {
    return {
      id: item.id,
      count: item.count,
      rate: item.rate,
      idx,
    }
  })
})

const selectItemGroup = (idx: number) => {
  selectedIdx.value = idx
}

const addItemForm = reactive({
  visible: false,
  id: 0,
  count: 1,
  rate: 1000,
})

const addItem = () => {
  addItemForm.visible = true
}

const addGroup = () => {
  formData.items.push(new BoxItemGroup({
    items: [] as box.BoxItem[],
  }))
}

const deleteGroup = () => {
  if (formData.items.length === 1) {
    message.warning('至少保留一个分组')
    return
  }
  formData.items = formData.items.filter((_, idx) => idx !== selectedIdx.value)
  if (selectedIdx.value >= formData.items.length) {
    selectedIdx.value = 0
  }
}

const modal = useModal()

const changeGroupName = async () => {
  const modalInt = modal.create({
    title: '修改分组名称',
    preset: 'dialog',
    content: () => {
      return h(
          InputModal,
          {
            onConfirm: async (val: string) => {
              formData.category_names[selectedIdx.value + 1] = val
              modalInt.destroy()
            },
          }
      )
    },
  })
}

const doAddItem = () => {
  formData.items[selectedIdx.value].items.push({
    id: addItemForm.id,
    count: addItemForm.count,
    rate: addItemForm.rate,
  })
}

const deleteItem = (idx: number) => {
  formData.items[selectedIdx.value].items = formData.items[selectedIdx.value].items.filter((_, i) => i !== idx)
}

const tableColumns = [
  { title: '物品id', key: 'id' },
  { title: '数量', key: 'count' },
  { title: '几率', key: 'rate' },
  {
    title: '',
    key: 'action',
    render(row: TableData) {
      return h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: 'small',
            onClick(e) {
              deleteItem(row.idx)
            },
          },
          { default: () => '删除' }
      )
    }
  },
]

const genResult = reactive({
  visible: false,
  content: '',
  path: '',
})

const genContent = async () => {
  genResult.content = await GenBox(new BoxParams(formData))
  const now = (new Date().getTime() / 1000).toFixed(0)
  genResult.path = `stackable/genbox/${now}.stk`
  genResult.visible = true
}

const saveResult = async () => {
  await SaveFile(genResult.path, genResult.content)
}

const copyResult = async () => {
  const clip = await navigator.clipboard
  await clip.writeText(genResult.content)
  message.success('已复制')
}

const simpleAddItemForm = reactive({
  visible: false,
  data: '',
})

const simpleAddItem = () => {
  simpleAddItemForm.data = ''
  simpleAddItemForm.visible = true
}

const doSimpleAddItem = () => {
  const items = simpleAddItemForm.data.split('\t').map(id => {
    return {
      id: parseInt(id),
      count: 1,
      rate: 1000,
    }
  })
  formData.items[selectedIdx.value].items = formData.items[selectedIdx.value].items.concat(items)
  simpleAddItemForm.visible = false
}

const getCateName = (idx: number) => {
  return formData.category_names[idx+1] || `etc${idx+1}`
}

onMounted(() => {
  if (formData.items.length === 0) {
    addGroup()
  }
})
</script>

<template>
  <n-form>
    <n-form-item label="类型">
      <n-radio-group v-model:value="formData.box_type">
        <n-radio :value="0">随机礼盒</n-radio>
        <n-radio :value="1">自选礼盒</n-radio>
      </n-radio-group>
    </n-form-item>
    <n-form-item label="名称">
      <n-input v-model:value="formData.name"></n-input>
    </n-form-item>
    <n-form-item label="描述">
      <n-input type="textarea" v-model:value="formData.explain"></n-input>
    </n-form-item>
    <n-form-item label="掉落等级">
      <n-input-number v-model:value="formData.grade"></n-input-number>
    </n-form-item>
    <n-form-item label="稀有度">
      <rarity-radio v-model:value="formData.rarity"></rarity-radio>
    </n-form-item>
    <n-form-item label="物品类型">
      <n-radio-group v-model:value="formData.tp">
        <n-radio value="equipment">装备</n-radio>
        <n-radio value="etc">道具</n-radio>
      </n-radio-group>
    </n-form-item>
    <n-form-item label="提示" v-if="formData.box_type === 1">
      <n-input v-model:value="formData.category_names[0]"></n-input>
    </n-form-item>
    <n-form-item label="礼盒内物品">
      <div class="flex flex-col">
        <n-flex>
          <n-button @click="addGroup">添加组</n-button>
          <n-button @click="deleteGroup">删除选中组</n-button>
          <n-button v-if="isSelection" @click="changeGroupName">修改组名</n-button>
          <n-button @click="addItem">添加物品</n-button>
          <n-button @click="simpleAddItem">批量添加</n-button>
        </n-flex>
        <div class="flex mt-3 h-[300px]">
          <n-list class="w-[200px]" bordered>
            <n-list-item
                v-for="(items, idx) in formData.items"
                :key="idx"
                :class="{'active-group': selectedIdx === idx}"
                style="cursor: pointer"
                @click="selectItemGroup(idx)"
            >
              <span>{{ getCateName(idx) }}</span>
            </n-list-item>
          </n-list>
          <n-data-table
              class="w-[500px]"
              striped
              size="small"
              :columns="tableColumns"
              :data="tableData"
              max-height="300px"
          ></n-data-table>
        </div>
      </div>
    </n-form-item>
    <n-form-item>
      <n-button type="primary" @click="genContent">生成</n-button>
    </n-form-item>
  </n-form>

  <n-modal v-model:show="addItemForm.visible">
    <n-card title="增加物品" style="width: 60%">
      <n-form>
        <n-form-item label="物品">
          <n-input-number v-model:value="addItemForm.id" :show-button="false"></n-input-number>
          <item-selector v-model:value="addItemForm.id"></item-selector>
        </n-form-item>
        <n-form-item label="数量">
          <n-input-number v-model:value="addItemForm.count" :min="1"></n-input-number>
        </n-form-item>
        <n-form-item label="几率">
          <n-input-number v-model:value="addItemForm.rate" :min="1" :show-button="false"></n-input-number>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-flex justify="end">
          <n-button @click="addItemForm.visible = false">取消</n-button>
          <n-button type="primary" @click="doAddItem">确定</n-button>
        </n-flex>
      </template>
    </n-card>
  </n-modal>

  <n-modal v-model:show="simpleAddItemForm.visible">
    <n-card title="批量添加, 相等几率" style="width: 60%">
      <n-input v-model:value="simpleAddItemForm.data" type="textarea"></n-input>
      <template #footer>
        <n-flex justify="end">
          <n-button @click="simpleAddItemForm.visible = false">取消</n-button>
          <n-button type="primary" @click="doSimpleAddItem">确定</n-button>
        </n-flex>
      </template>
    </n-card>
  </n-modal>

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
</template>

<style scoped>
.active-group {
  background-color: #777;
}
</style>
