<template>
  <div>{{ info.name }}</div>
  <div>{{ info.source }}</div>
  <div>{{ info.description }}</div>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { APIInfo } from "@/models/api-info";

export default defineComponent({
  name: "APICard",
  props: {
    info: {
      type: Object as PropType<APIInfo>,
      required: true,
    },
  },
});
</script>
