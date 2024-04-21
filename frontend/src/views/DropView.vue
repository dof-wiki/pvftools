<script setup lang="ts">
import {onMounted, ref} from "vue";
import {GetExtraMonsterDrop, GetMonsterDrop, SetMonsterDrop, SetMonsterExtraDrop} from "../../wailsjs/go/api/App";
import HellDrop from "../components/HellDrop.vue";
import ClearRewardDrop from "../components/ClearRewardDrop.vue";
import MonsterDrop from "../components/MonsterDrop.vue";

const width = ref(0)

onMounted(() => {
  width.value = window.document.body.clientWidth
  window.onresize = () => {
    width.value = window.document.body.clientWidth
  }
})
</script>

<template>
  <div>
    <n-grid
        :cols="width < 1200 ? 1 : 2"
    >
      <n-gi class="drop-item">
        <hell-drop />
      </n-gi>
      <n-gi class="drop-item">
        <clear-reward-drop />
      </n-gi>
      <n-gi class="drop-item">
        <monster-drop
            title="怪物"
            :get-rate="GetMonsterDrop"
            :save-rate="SetMonsterDrop"
        />
      </n-gi>
      <n-gi class="drop-item">
        <monster-drop
            title="绿名怪物"
            :get-rate="GetExtraMonsterDrop"
            :save-rate="SetMonsterExtraDrop"
        />
      </n-gi>
    </n-grid>
  </div>
</template>

<style scoped>

</style>