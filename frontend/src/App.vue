<script lang="ts" setup>
import {darkTheme} from 'naive-ui'
import Layout from "./views/Layout.vue";
</script>

<template>
  <n-config-provider :theme="darkTheme">
    <n-message-provider>
      <n-dialog-provider>
        <n-modal-provider>
          <Layout></Layout>
        </n-modal-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<style>
</style>
