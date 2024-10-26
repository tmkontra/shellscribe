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
  <div class="flex flex-col w-1/3 h-full">
    <h1 class="text-3xl font-bold underline">shellscribe</h1>
    <template v-for="f in state.files" v-bind:key="f.Id">
      <div
        class="flex flex-row py-2" 
        :class="f.Id === state?.selected?.Id ? 'bg-black' : ''" 
        @click="setSelected(f)"
      >
        <p> 
        {{ f.Name }}
        </p>
        <p>{{ new Date(f.CreatedAt).toLocaleString() }}</p>
      </div>
    </template>
  </div>
</template>

<style scoped>
</style>
