<script setup lang="ts">
import {computed, reactive, ref} from "vue";
import {useMessage} from 'naive-ui'
import * as equAttrDict from "../resources/equ_attr.json"
import {Trash, Save} from '@vicons/ionicons5'
import {model} from "../../wailsjs/go/models";
import Skill = model.Skill;

type ChoiceItem = {
  label: string
  value: string
}

type TInnerSelector = 'skill' | 'skill-increment'

type TAttr = {
  name: string
  tmpl: string
  desc?: string[]
  choices?: Array<string | null>
  inner_selector?: Array<TInnerSelector>
}

type AttrItem = {
  name: string
  content: string
}

const getChoices = (key: string): ChoiceItem[] => {
  return (equAttrDict.choices as any)[key]
}

const cate = ref('base')

const cateOptions = [
  {name: '基础属性', value: 'base'},
  {name: '抗性属性', value: 'resistance'},
  {name: '速度属性', value: 'speed'},
  {name: '武器元素属性', value: 'elemental_weapon'},
  {name: '触发条件', value: 'trigger'},
  {name: '触发效果', value: 'then'},
]

const filterText = ref('')

const attrOptions = computed<Array<TAttr>>(() => {
  const list = (equAttrDict as any)[cate.value] as Array<TAttr>
  if (!filterText.value) {
    return list
  }
  return list.filter((item: TAttr) => item.name.includes(filterText.value))
})

const attrs = ref<Array<AttrItem>>([])

const attrCtx = reactive({
  show: false,
  attrTmpl: null as TAttr | null,
  values: [] as Array<any>,
})

const addAttr = (name: string) => {
  const attr = attrOptions.value.find((item: any) => item.name == name)
  if (!attr) {
    return
  }
  attrCtx.attrTmpl = attr
  attrCtx.values = (attr.tmpl.match(/{\d+}/g) || []).map((v: any) => null)
  if (attrCtx.values.length === 0) {
    confirmAddAttr()
  } else {
    attrCtx.show = true
  }
}

const confirmAddAttr = () => {
  if (!attrCtx.attrTmpl) {
    return
  }
  const matches = attrCtx.attrTmpl.tmpl.match(/{\d+}/g)
  let attr;
  if (!matches) {
    attr = {
      name: attrCtx.attrTmpl.name,
      content: attrCtx.attrTmpl.tmpl,
    }
  } else {
    let tmpl = attrCtx.attrTmpl.tmpl
    matches.forEach((v) => {
      const idx = Number(v.slice(1, -1))
      tmpl = tmpl.replace(v, attrCtx.values[idx])
    })
    attr = {
      name: attrCtx.attrTmpl.name,
      content: tmpl
    }
  }
  if (triggerCtx.inTrigger) {
    if (cate.value === 'trigger') {
      triggerCtx.trigger.push(attr)
    } else if (cate.value === 'then') {
      triggerCtx.then.push(attr)
    }
  } else {
    attrs.value.push(attr)
  }
  attrCtx.show = false
}

const getDescList = (attr: TAttr) => {
  if (attr.desc) {
    return attr.desc
  }
  const matches = attr.tmpl.match(/{\d+}/g)
  if (!matches) {
    return []
  }
  return matches.map((v, idx) => `数值${idx + 1}`)
}

const deleteAttr = (idx: number) => {
  attrs.value.splice(idx, 1)
}

const message = useMessage()

const copyAll = () => {
  const text = attrs.value.map((item) => item.content).join('\n\n')
  navigator.clipboard.writeText(text)
  message.success('已复制到剪贴板')
}

const clearAll = () => {
  attrs.value = []
}

const triggerCtx = reactive({
  inTrigger: false,
  trigger: [] as Array<AttrItem>,
  then: [] as Array<AttrItem>
})

const triggerContent = computed(() => {
  if (triggerCtx.inTrigger) {
    return ('[if]\n\t' + triggerCtx.trigger.map((item) => item.content).join('\n\t') + '\n[/if]\n\n[then]\n\t' +
        triggerCtx.then.map(item => item.content).join('\n\t') + '\n[/then]')
  }
  return ''
})

const isTriggerCate = (c: string) => {
  return c === 'trigger' || c === 'then'
}

const selectCate = (c: string) => {
  if (triggerCtx.inTrigger && !isTriggerCate(c)) {
    if (triggerCtx.trigger.length > 0 || triggerCtx.then.length > 0) {
      message.warning('有修改尚未保存, 请先处理')
      return
    }
  }
  cate.value = c
  if (isTriggerCate(c)) {
    triggerCtx.inTrigger = true
  } else {
    triggerCtx.inTrigger = false
  }
}

const saveTrigger = () => {
  attrs.value.push({
    name: '触发效果',
    content: triggerContent.value
  })
  cate.value = 'base'
  triggerCtx.inTrigger = false
}

const deleteTrigger = () => {
  cate.value = 'base'
  triggerCtx.inTrigger = false
}
</script>

<template>
  <div class="flex h-full">
    <div class="border-r border-gray-800">
      <n-list clickable hoverable class="h-full">
        <n-list-item
            v-for="item in cateOptions"
            :key="item.value"
            :class="{'active': cate === item.value}"
            @click="selectCate(item.value)"
        >{{ item.name }}</n-list-item>
      </n-list>
    </div>
    <div class="border-r border-gray-800">
      <n-list clickable hoverable class="h-full overflow-auto">
        <n-list-item>
          <n-input v-model:value="filterText" size="small" clearable placeholder="搜索..."></n-input>
        </n-list-item>
        <n-list-item v-for="item in attrOptions" :key="item.name" @click="addAttr(item.name)">
          {{ item.name }}
        </n-list-item>
      </n-list>
    </div>
    <div class="flex-grow ml-2 h-full overflow-auto">
      <n-flex justify="end" class="mb-2">
        <n-button @click="clearAll">清空</n-button>
        <n-button type="primary" @click="copyAll">复制</n-button>
      </n-flex>
      <n-card v-for="(item, idx) in attrs" :key="idx" class="mb-2">
        <template #header>
          {{ item.name }}
        </template>
        <template #header-extra>
          <n-button text class="text-2xl" @click="deleteAttr(idx)">
            <n-icon>
              <trash/>
            </n-icon>
          </n-button>
        </template>
        <n-input type="textarea" v-model:value="item.content"></n-input>
      </n-card>
      <n-card v-if="triggerCtx.inTrigger">
        <template #header>
          触发
        </template>
        <template #header-extra>
          <n-flex>
            <n-button text class="text-2xl" @click="saveTrigger">
              <n-icon><save/></n-icon>
            </n-button>
            <n-button text class="text-2xl" @click="deleteTrigger">
              <n-icon><trash/></n-icon>
            </n-button>
          </n-flex>
        </template>
        <n-input type="textarea" :rows="10" :value="triggerContent" disabled></n-input>
      </n-card>
    </div>
    <n-modal v-model:show="attrCtx.show">
      <n-card style="width: 30%;" v-if="attrCtx.attrTmpl">
        <template #header>
          {{ attrCtx.attrTmpl.name }}
        </template>
        <n-form>
          <n-form-item v-for="(item, idx) in getDescList(attrCtx.attrTmpl)" :key="idx" :label="item">
            <n-select
                v-if="attrCtx.attrTmpl.choices && attrCtx.attrTmpl.choices[idx]"
                v-model:value="attrCtx.values[idx]"
                :options="getChoices(attrCtx.attrTmpl.choices[idx]!)"
                filterable
            ></n-select>
            <n-input v-else v-model:value="attrCtx.values[idx]"></n-input>
          </n-form-item>
        </n-form>
        <template #footer>
          <n-flex justify="end">
            <n-button @click="attrCtx.show = false">取消</n-button>
            <n-button type="primary" @click="confirmAddAttr">确定</n-button>
          </n-flex>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<style scoped>
.active {
  background: rgba(255, 255, 255, 0.1);
}
</style>
