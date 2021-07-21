<template>
  <div>{{ info.name }}</div>
  <div>{{ info.source }}</div>
  <div>{{ info.description }}</div>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { Product } from "@/models/product";

export default defineComponent({
  name: "APICard",
  props: {
    info: {
      type: Object as PropType<Product>,
      required: true,
    },
  },
});
</script>
