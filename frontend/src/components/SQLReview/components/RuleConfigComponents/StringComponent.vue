<template>
  <BBTextField
    :value="value"
    :disabled="disabled"
    :placeholder="`${config.payload.default}`"
    @update:value="$emit('update:value', $event)"
  />
</template>

<script lang="ts" setup>
import { BBTextField } from "@/bbkit";
import type { RuleConfigComponent } from "@/types";

defineProps<{
  config: RuleConfigComponent;
  value: string;
  disabled: boolean;
}>();

defineEmits<{
  (event: "update:value", value: string): void;
}>();
</script>
