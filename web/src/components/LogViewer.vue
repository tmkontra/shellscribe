<script setup lang="ts">
import { onBeforeUnmount, reactive, watchEffect } from 'vue';
import { Item } from '../api/types';


const props = defineProps<{
    selected: Item | null
}>()

const state = reactive<{
    logContent: string;
    eventSource: EventSource | null;
    currentSource: Item | null;
}>({
    logContent: "",
    eventSource: null,
    currentSource: null,
})

watchEffect(() => {
    if (!props.selected) {
        return 
    }
    if (props.selected?.Id == state.currentSource?.Id) {
        return
    }
    if (state.eventSource !== null) {
        state.eventSource.close()
    }
    state.logContent = ''
    const source = new EventSource("/tail/" + props.selected?.Id);
    source.onmessage = (event) => {
        state.logContent = state.logContent + event.data + "\n"
    };
    state.eventSource = source
})

onBeforeUnmount(() => {
    state.eventSource?.close()
})

</script>

<template>
    <div class="text-left">
        <pre class="break-words text-wrap overflow-wrap">{{ state.logContent }}</pre>
    </div>
</template> 

<style scoped>

</style>