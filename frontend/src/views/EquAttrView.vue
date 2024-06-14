<script setup lang="ts">
import {computed, onMounted, reactive, ref} from "vue";
import {useDialog, useMessage} from 'naive-ui'
import * as equAttrDict from "../resources/equ_attr.json"
import {Trash, Save} from '@vicons/ionicons5'
import {model, proto} from "../../wailsjs/go/models";
import Skill = model.Skill;
import storage from "../common/storage";
import SkillList from "../components/SkillList.vue";
import SkillIncrementList from "../components/SkillIncrementList.vue";
import SkillIncrement = proto.SkillIncrement;
import {
  AddCustomAttrTmpl,
  DeleteCustomAttrTmpl,
  GetCustomAttrTmpls,
  UpdateCustomAttrTmpl
} from "../../wailsjs/go/api/App";
import CustomAttrTmpl = model.CustomAttrTmpl;

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
  id?: number
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
  {name: '技能', value: 'skill'},
  {name: '触发条件', value: 'trigger'},
  {name: '触发效果', value: 'then'},
  {name: '自定义模板', value: 'custom'},
]

const specialAttrs: Record<string, TAttr[]> = {
  skill: [
    {name: '技能等级', tmpl: 'level'},
    {name: '技能数据', tmpl: 'data'},
  ]
}

const customAttrs = ref<Array<TAttr>>([])

const filterText = ref('')

const attrOptions = computed<Array<TAttr>>(() => {
  let list: Array<TAttr> = []
  if (cate.value === 'custom') {
    list = customAttrs.value
  } else {
    list = specialAttrs[cate.value]
    if (!list) {
      list = (equAttrDict as any)[cate.value] as Array<TAttr>
    }
  }
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

const skillLevelCtx = reactive({
  showSelector: false,
  skills: [] as Skill[],
  level: 1,
})

const skillDataCtx = reactive({
  show: false,
  data: [] as Array<SkillIncrement>,
})

const addAttr = (item: TAttr) => {
  if (cate.value === 'custom') {
    customAttrOperationCtx.item = item
    customAttrOperationCtx.show = true
  } else {
    innerAddAttr(item)
  }
}

const innerAddAttr = (item: TAttr) => {
  const attr = attrOptions.value.find((it: any) => {
    if (item.id) {
      return it.id === item.id
    }
    return it.name === item.name
  })
  if (!attr) {
    return
  }
  switch (attr.tmpl) {
    case 'level':
      skillLevelCtx.level = 1
      skillLevelCtx.showSelector = true
      break
    case 'data':
      skillDataCtx.data = []
      skillDataCtx.show = true
      break
    default:
      attrCtx.attrTmpl = attr
      attrCtx.values = (attr.tmpl.match(/{\d+}/g) || []).map((v: any) => null)
      if (attrCtx.values.length === 0) {
        confirmAddAttr()
      } else {
        attrCtx.show = true
      }
  }
}

const onSelectSkills = (skls: Skill[]) => {
  skillLevelCtx.skills = skls
}

const onSkillIncrUpdate = async (newData: Array<SkillIncrement>) => {
  skillDataCtx.data = newData
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

const confirmAddSkillLevel = () => {
  if (skillLevelCtx.level <= 0 || skillLevelCtx.skills.length === 0) {
    return
  }
  const rows = skillLevelCtx.skills.map((skl) => {
    const job = `[${storage.getJobStr(skl.job)}]`
    return `\`${job}\`\t${skl.code}\t${skillLevelCtx.level}`
  })
  const item = attrs.value.find((v) => v.name === '技能等级')
  if (item) {
    item.content = item.content.replace('[/skill levelup]', '\t' + rows.join('\n\t') + '\n[/skill levelup]')
  } else {
    attrs.value.push({
      name: '技能等级',
      content: '[skill levelup]\n\t' + rows.join('\n\t') + '\n[/skill levelup]'
    })
  }
  skillLevelCtx.showSelector = false
}

const confirmAddSkillData = () => {
  skillDataCtx.show = false
  if (skillDataCtx.data.length === 0) {
    return
  }
  const rows = skillDataCtx.data.map((item: SkillIncrement) => {
    return `\`${item.job}\`\t${item.skill_id}\t\`${item.type}\`\t\`${item.data_type}\`\t${item.data_index}\t\`${item.increment_type}\`\t${item.value}`
  })
  const item = attrs.value.find((v) => v.name === '技能数据')
  if (item) {
    item.content = item.content.replace('[/skill data up]', '\t' + rows.join('\n\t') + '\n[/skill data up]')
  } else {
    attrs.value.push({
      name: '技能数据',
      content: '[skill data up]\n\t' + rows.join('\n\t') + '\n[/skill data up]'
    })
  }
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

const customAttrTmplCtx = reactive({
  show: false,
  formData: {
    name: '',
    tmpl: '',
    desc: [] as Array<string>,
    choices: [] as Array<string | null>,
  },
  editId: 0,
})

const addCustomAttrTmpl = () => {
  customAttrTmplCtx.formData = {
    name: '',
    tmpl: '',
    desc: [] as Array<string>,
    choices: [] as Array<string | null>,
  }
  customAttrTmplCtx.editId = 0
  customAttrTmplCtx.show = true
}

const addCustomAttrTmplValue = () => {
  customAttrTmplCtx.formData.choices.push(null)
  const idx = customAttrTmplCtx.formData.desc.length
  customAttrTmplCtx.formData.desc.push(`数值${idx}`)
  customAttrTmplCtx.formData.tmpl += `{${idx}}`
}

const rmCustomAttrTmplValue = () => {
  if (customAttrTmplCtx.formData.choices.length < 1) {
    return
  }
  customAttrTmplCtx.formData.choices.pop()
  customAttrTmplCtx.formData.desc.pop()
  const idx = customAttrTmplCtx.formData.desc.length
  customAttrTmplCtx.formData.tmpl = customAttrTmplCtx.formData.tmpl.replace(`{${idx}}`, '')
}

const confirmAddCustomTmpl = async () => {
  const tmpl = new CustomAttrTmpl({
    name: customAttrTmplCtx.formData.name,
    tmpl: customAttrTmplCtx.formData.tmpl,
    desc: JSON.stringify(customAttrTmplCtx.formData.desc),
    choices: JSON.stringify(customAttrTmplCtx.formData.choices),
    id: customAttrTmplCtx.editId || undefined,
  })
  if (customAttrTmplCtx.editId) {
    await UpdateCustomAttrTmpl(tmpl)
  } else {
    await AddCustomAttrTmpl(tmpl)
  }
  await fetchCustomAttrTmpl()
  customAttrTmplCtx.show = false
}

const tmplSelectors = computed(() => {
  return Object.keys(equAttrDict.choices as Record<string, any>).map((key) => {
    return {
      value: key,
      label: key
    }
  })
})

const fetchCustomAttrTmpl = async () => {
  const tmpls = await GetCustomAttrTmpls()
  customAttrs.value = tmpls.map((item: CustomAttrTmpl) => {
    return {
      id: item.id,
      name: item.name,
      tmpl: item.tmpl,
      desc: JSON.parse(item.desc),
      choices: JSON.parse(item.choices),
    }
  })
}

const handleTab = (e: KeyboardEvent) => {
  e.preventDefault()
  const el = e.target as HTMLTextAreaElement
  const start = el.selectionStart
  const end = el.selectionEnd
  const value = el.value
  el.value = value.substring(0, start) + '\t' + value.substring(end)
  el.selectionStart = el.selectionEnd = start + 1
}

const customAttrOperationCtx = reactive({
  show: false,
  item: null as TAttr | null,
})

const dialog = useDialog()

const handleCustomAttrSelect = (key: string) => {
  const item = customAttrOperationCtx.item!
  customAttrOperationCtx.show = false
  switch (key) {
    case 'insert':
      innerAddAttr(item)
      break
    case 'edit':
      customAttrTmplCtx.formData = {
        name: item.name,
        tmpl: item.tmpl,
        desc: item.desc || [],
        choices: item.choices || [],
      }
      customAttrTmplCtx.editId = item.id!
      customAttrTmplCtx.show = true
      break
    case 'delete':
      dialog.warning({
        title: '提示',
        content: `确认删除模板<${item.name}>?`,
        positiveText: '确认',
        negativeText: '取消',
        onPositiveClick(e) {
          DeleteCustomAttrTmpl(item.id!).then(() => fetchCustomAttrTmpl())
        }
      })
      break
  }
}

onMounted(async () => {
  await fetchCustomAttrTmpl()
})
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
          <n-input-group>
            <n-input v-model:value="filterText" size="small" clearable placeholder="搜索..."></n-input>
            <n-button v-if="cate === 'custom'" @click="addCustomAttrTmpl">+</n-button>
          </n-input-group>
        </n-list-item>
        <n-list-item v-for="item in attrOptions" :key="item.name" @click="addAttr(item)">
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
        <n-input type="textarea" v-model:value="item.content" @keydown.tab="handleTab"></n-input>
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

    <n-drawer v-model:show="skillLevelCtx.showSelector" :width="300">
      <n-form label-placement="left" class="mt-2 ml-2">
        <n-flex justify="end" class="mb-2">
          <n-button type="primary" @click="confirmAddSkillLevel">确定</n-button>
        </n-flex>
        <n-form-item label="等级">
          <n-input-number v-model:value="skillLevelCtx.level" :min="1"></n-input-number>
        </n-form-item>
        <n-form-item>
          <skill-list class="h-[85vh]" multi-select @multi-select="onSelectSkills"></skill-list>
        </n-form-item>
      </n-form>
    </n-drawer>

    <n-modal v-model:show="skillDataCtx.show">
      <n-card style="width: 50%">
        <template #header>
          <span>技能数据</span>
        </template>
        <skill-increment-list
            :data="skillDataCtx.data"
            @update="onSkillIncrUpdate"
        ></skill-increment-list>
        <template #footer>
          <n-flex justify="end">
            <n-button type="primary" @click="confirmAddSkillData">确定</n-button>
          </n-flex>
        </template>
      </n-card>
    </n-modal>

    <n-modal v-model:show="customAttrTmplCtx.show">
      <n-card style="width: 50%">
        <template #header>
          <span>增加模板</span>
        </template>
        <n-form>
          <n-form-item label="模板名">
            <n-input v-model:value="customAttrTmplCtx.formData.name"></n-input>
          </n-form-item>
          <n-form-item label="模板">
            <n-input
                v-model:value="customAttrTmplCtx.formData.tmpl"
                type="textarea"
                :rows="5"
                class="input-with-tab"
                @keydown.tab="handleTab"
            ></n-input>
          </n-form-item>
          <n-form-item label="模板变量">
            <div>
              <n-button-group size="small" class="mb-5">
                <n-button @click="addCustomAttrTmplValue">+</n-button>
                <n-button @click="rmCustomAttrTmplValue">-</n-button>
              </n-button-group>
              <n-form-item v-for="(item, idx) in customAttrTmplCtx.formData.desc" :label="`变量${idx}`">
                <n-form-item label="描述">
                  <n-input v-model:value="customAttrTmplCtx.formData.desc[idx]"></n-input>
                </n-form-item>
                <n-form-item label="选项(可选)">
                  <n-select
                      class="ml-2"
                      style="width: 150px;"
                      :options="tmplSelectors"
                      v-model:value="customAttrTmplCtx.formData.choices[idx]"
                      clearable
                  ></n-select>
                </n-form-item>
              </n-form-item>
            </div>
          </n-form-item>
        </n-form>
        <template #footer>
          <n-flex justify="end">
            <n-button @click="customAttrTmplCtx.show = false">取消</n-button>
            <n-button @click="confirmAddCustomTmpl" type="primary">保存</n-button>
          </n-flex>
        </template>
      </n-card>
    </n-modal>
    <n-modal v-model:show="customAttrOperationCtx.show">
      <n-card style="width: 15%;">
        <template #header>
          <span>{{ customAttrOperationCtx.item?.name }}</span>
        </template>
        <n-list hoverable clickable>
          <n-list-item class="text-center" @click="handleCustomAttrSelect('insert')">插入</n-list-item>
          <n-list-item class="text-center" @click="handleCustomAttrSelect('edit')">编辑</n-list-item>
          <n-list-item class="text-center" @click="handleCustomAttrSelect('delete')">删除</n-list-item>
        </n-list>
      </n-card>
    </n-modal>
  </div>
</template>

<style scoped>
.active {
  background: rgba(255, 255, 255, 0.1);
}
</style>
