<script setup lang="ts">
import { reactive } from 'vue'
import client from '../api/client';
import { Item } from '../api/types';

const state = reactive<{
  files: {
    Name: string;
    Id: string;
    CreatedAt: string;
  }[],
  selected: Item | null
}>({
  files: [],
  selected: null,
})

client.get("/index").then((resp) => {
  state.files = resp.data.Logs
})


const props = defineProps<{
  setSelected: (item: Item) => void
}>()

const setSelected = (item: Item) => {
  props.setSelected(item)
  state.selected = item
}

</script>

<template>
  <div class="overflow-scroll">
    <h1 class="text-3xl my-4 font-bold underline">shellscribe</h1>
    <div class="flex flex-col gap-4 px-4">
      <template v-for="f in state.files" v-bind:key="f.Id">
        <div
          class="flex flex-row justify-between py-4 rounded-md px-4" 
          :class="f.Id === state?.selected?.Id ? 'dark:bg-slate-900' : ''" 
          @click="setSelected(f)"
        >
          <p class="max-w-3/4 text-ellipsis"> 
          {{ f.Name }}
          </p>
          <p>{{ new Date(f.CreatedAt).toLocaleString() }}</p>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
</style>
