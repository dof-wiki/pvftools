<script setup lang="ts">
import {ref} from "vue";

const props = withDefaults(defineProps<{
  message: string,
  onConfirm: (value: string) => Promise<void>
  onCancel?: () => Promise<void>
}>(), {
  message: '请输入',
})

const inputValue = ref('')
</script>

<template>
  <n-form style="margin-top: 20px;">
    <n-form-item>
      <n-input v-model:value="inputValue" :placeholder="message"></n-input>
    </n-form-item>
    <n-flex justify="end">
      <n-button v-if="onCancel" @click="onCancel">取消</n-button>
      <n-button type="primary" @click="onConfirm(inputValue)">确定</n-button>
    </n-flex>
  </n-form>
</template>

<style scoped>

</style>
