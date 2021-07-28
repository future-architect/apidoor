<template>
  <ion-card>
    <ion-card-header>
      <ion-card-title>
        {{ info.name }}
      </ion-card-title>
      <ion-card-subtitle>
        {{ info.source }}
      </ion-card-subtitle>
    </ion-card-header>

    <ion-card-content>
      {{ info.description }}
    </ion-card-content>
  </ion-card>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import {
  IonCard,
  IonCardContent,
  IonCardSubtitle,
  IonCardTitle,
  IonCardHeader,
} from "@ionic/vue";
import { Product } from "@/models/product";

export default defineComponent({
  name: "APICard",
  components: {
    IonCard,
    IonCardTitle,
    IonCardSubtitle,
    IonCardContent,
    IonCardHeader,
  },
  props: {
    info: {
      type: Object as PropType<Product>,
      required: true,
    },
  },
});
</script>
