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
        const logView = document.getElementById("log-content");
        logView?.scrollTo(0, logView.scrollHeight)
    };
    state.eventSource = source
})

onBeforeUnmount(() => {
    state.eventSource?.close()
})

</script>

<template>
    <div class="col-span-2 h-full max-h-screen p-4">
        <div class="flex flex-col h-full justify-center items-center">
            <div class="h-full w-full dark:bg-stone-800 rounded-lg p-4">
                <pre id="log-content" class="text-left break-words text-wrap h-full overflow-scroll">{{ state.logContent }}</pre>
            </div>
        </div>
    </div>
</template> 

<style scoped>

</style>