<template>
  <NSelect
    :value="index.expressions"
    :options="columnOptions"
    :disabled="readonly"
    :consistent-menu-width="true"
    :style="style"
    :multiple="true"
    :max-tag-count="'responsive'"
    :placeholder="$t('schema-editor.columns')"
    :filterable="true"
    :render-tag="renderTag"
    :show-arrow="!readonly"
    suffix-style="right: 3px"
    class="bb-schema-editor--index-columns-select"
    @focus="focused = true"
    @blur="focused = false"
    @update:value="$emit('update:expressions', $event as string[])"
  />
</template>

<script lang="ts" setup>
import type { SelectOption } from "naive-ui";
import { NSelect, NTag } from "naive-ui";
import type { CSSProperties } from "vue";
import { computed, h, ref } from "vue";
import type { ComposedDatabase } from "@/types";
import type {
  ColumnMetadata,
  DatabaseMetadata,
  IndexMetadata,
  SchemaMetadata,
  TableMetadata,
} from "@/types/proto-es/v1/database_service_pb";

type ColumnOption = SelectOption & {
  label: string;
  value: string;
  column: ColumnMetadata;
};

const props = defineProps<{
  readonly?: boolean;
  db: ComposedDatabase;
  database: DatabaseMetadata;
  schema: SchemaMetadata;
  table: TableMetadata;
  index: IndexMetadata;
}>();
defineEmits<{
  (event: "update:expressions", expressions: string[]): void;
}>();

const focused = ref(false);

const columnOptions = computed(() => {
  return props.table.columns.map<ColumnOption>((column) => ({
    label: column.name,
    value: column.name,
    column,
  }));
});

const style = computed(() => {
  const style: CSSProperties = {
    "--n-color": "transparent",
    "--n-color-disabled": "transparent",
    "--n-text-color-disabled": "rgb(var(--color-main))",
    cursor: "default",
  };
  const border = focused.value
    ? "1px solid rgb(var(--color-control-border))"
    : "none";
  style["--n-border"] = border;
  style["--n-border-disabled"] = border;

  return style;
});

const renderTag = (item: { option: SelectOption; handleClose: () => void }) => {
  const { option, handleClose } = item;
  const { column } = option as ColumnOption;
  return h(
    NTag,
    {
      closable: !props.readonly,
      onClose: handleClose,
    },
    { default: () => column.name }
  );
};
</script>

<style lang="postcss" scoped>
.bb-schema-editor--index-columns-select :deep(.n-base-selection) {
  --n-padding-multiple: 2px 16px 0px 2px !important;
  --n-color: transparent !important;
  --n-color-disabled: transparent !important;
  --n-border: none !important;
  --n-text-color-disabled: rgb(var(--color-main)) !important;
}
.bb-schema-editor--index-columns-select
  :deep(.n-base-selection .n-base-suffix) {
  right: 4px;
}
</style>
