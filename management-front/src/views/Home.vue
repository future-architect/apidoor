<template>
  <h1>This is API management page<br /></h1>
  <div class="container">
    <div v-for="product in APIs.products" :key="product.id">
      <APICard :info="product" />
    </div>
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

<style lang="scss">
.container {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  align-content: flex-start;
}
</style>
