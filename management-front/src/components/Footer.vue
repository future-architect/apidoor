<template>
  <ion-footer>
    <ion-toolbar> This is Footer </ion-toolbar>
  </ion-footer>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { IonToolbar, IonFooter } from "@ionic/vue";

export default defineComponent({
  name: "footer-page",
  components: {
    IonToolbar,
    IonFooter,
  },
});
</script>
