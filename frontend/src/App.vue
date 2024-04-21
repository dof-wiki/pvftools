<script lang="ts" setup>
import {darkTheme} from 'naive-ui'
import Layout from "./views/Layout.vue";
</script>

<template>
  <n-config-provider :theme="darkTheme">
    <n-message-provider>
      <Layout></Layout>
    </n-message-provider>
  </n-config-provider>
</template>

<style>
</style>
