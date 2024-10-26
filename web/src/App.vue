<script setup lang="ts">
import Sidebar from "./components/Sidebar.vue";
import LogViewer from "./components/LogViewer.vue";
import { reactive } from "vue";
import { Item } from "./api/types";

const state = reactive<{
  selected: Item | null;
}>({
  selected: null,
});

const setSelected = (item: Item) => {
  state.selected = item;
};
</script>

<template>
  <div class="flex flex-row h-full">
    <Sidebar :set-selected="setSelected" />
    <LogViewer :selected="state.selected" />
  </div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
