<template>
    <div class="flex flex-col w-2/3">
        <div>
            <h1>Logs: {{  props.selected?.Name }}</h1>
        </div>
        <pre class="break-words text-wrap overflow-wrap">
            {{ state.logContent }}
        </pre>
    </div>
</template> 

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
    console.log(props.selected)
    if (!props.selected) {
        return 
    }
    if (props.selected?.Id == state.currentSource?.Id) {
        return
    }
    if (state.eventSource !== null) {
        console.log("closing")
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

<style scoped>

</style>