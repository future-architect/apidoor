<template>
  <ion-header>
    <ion-toolbar>
      <ion-title>This is Header</ion-title>
    </ion-toolbar>
  </ion-header>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { IonHeader, IonToolbar, IonTitle } from "@ionic/vue";

export default defineComponent({
  name: "header-page",
  components: {
    IonHeader,
    IonTitle,
    IonToolbar,
  },
});
</script>
