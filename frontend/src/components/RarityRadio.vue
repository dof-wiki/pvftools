<script setup lang="ts">
const model = defineModel()
</script>

<template>
  <n-radio-group v-model="model">
    <n-radio :value="0">普通</n-radio>
    <n-radio :value="1">精良</n-radio>
    <n-radio :value="2">稀有</n-radio>
    <n-radio :value="3">神器</n-radio>
    <n-radio :value="4">史诗</n-radio>
    <n-radio :value="5">异界</n-radio>
  </n-radio-group>
</template>

<style scoped>

</style>
