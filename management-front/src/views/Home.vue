<template>
  <div>This is API management page<br /></div>
  <div v-for="product in APIs.products" :key="product.id">
    <APICard :info="product" />
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";
import APICard from "@/components/APICard.vue";
import { getProducts } from "@/api/get-products";
import { Products } from "@/models/products";

export default defineComponent({
  name: "Home",
  components: {
    APICard,
  },
  setup() {
    const APIs = reactive<Products>({ products: [] });
    const getAPIs = async () => {
      const data = await getProducts();
      APIs.products = data.products;
    };

    return {
      APIs,
      getAPIs,
    };
  },
  mounted() {
    this.getAPIs();
  },
});
</script>
