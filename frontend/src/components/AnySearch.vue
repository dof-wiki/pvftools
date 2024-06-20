<script setup lang="ts">
import {onMounted, reactive, ref} from "vue";
import {DataTableColumns} from "naive-ui";
import {api} from "../../wailsjs/go/models";
import SearchResult = api.SearchResult;
import {SearchByName} from "../../wailsjs/go/api/App";

const props = defineProps<{
  searchType?: number
}>()

const searchOptions = [
  {label: '道具', value: 1},
  {label: '装备', value: 2},
  {label: '技能', value: 3},
  {label: 'NPC', value: 4},
  {label: '副本', value: 5},
]

const filters = reactive({
  q: '',
  searchType: 1,
})

const columns: DataTableColumns<SearchResult> = [
  {
    key: 'id',
    title: 'ID',
  },
  {
    key: 'name',
    title: '名称',
  },
]

const data = ref<SearchResult[]>([])
const loading = ref(false)

const search = async () => {
  loading.value = true
  try {
    data.value = await SearchByName(filters.q, filters.searchType)
  } finally {
    loading.value = false
  }
}

const rowKey = (row: SearchResult) => {
  return row.id
}

const emit = defineEmits(['select'])

const rowProps = (row: SearchResult) => {
  return {
    onClick: () => {
      emit('select', row.id, row.name)
    }
  }
}

onMounted(() => {
  if (props.searchType) {
    filters.searchType = props.searchType
  }
})

</script>

<template>
  <div>
    <n-form label-placement="left">
      <n-form-item>
        <n-radio-group v-model:value="filters.searchType" :disabled="searchType">
          <n-radio v-for="opt in searchOptions" :key="opt.value" :label="opt.label" :value="opt.value"></n-radio>
        </n-radio-group>
      </n-form-item>
      <n-form-item>
        <n-input-group>
          <n-input v-model:value="filters.q" placeholder="输入关键词" @keydown.enter="search"></n-input>
          <n-button :loading="loading" type="primary" @click="search">搜索</n-button>
        </n-input-group>
      </n-form-item>
    </n-form>
    <n-data-table
        :columns="columns"
        :data="data"
        :loading="loading"
        max-height="80vh"
        :row-key="rowKey"
        :row-props="rowProps"
        row-class-name="search-result-item"
        virtual-scroll
    ></n-data-table>
  </div>
</template>

<style>
.search-result-item {
  cursor: pointer;
}
</style>
