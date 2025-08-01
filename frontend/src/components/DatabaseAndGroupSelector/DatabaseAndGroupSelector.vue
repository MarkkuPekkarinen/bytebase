<template>
  <NTabs
    v-model:value="databaseSelectState.changeSource"
    type="card"
    size="small"
  >
    <NTabPane name="DATABASE" :tab="$t('common.databases')">
      <AdvancedSearch
        v-model:params="searchParams"
        class="w-full mb-2"
        :autofocus="false"
        :placeholder="$t('database.filter-database')"
        :scope-options="scopeOptions"
        :override-route-query="false"
      />
      <PagedDatabaseTable
        mode="PROJECT_SHORT"
        :show-selection="true"
        :custom-click="true"
        :parent="project.name"
        :filter="filter"
        :size="'small'"
        v-model:selected-database-names="
          databaseSelectState.selectedDatabaseNameList
        "
      />
    </NTabPane>
    <NTabPane name="GROUP" :tab="$t('common.database-group')">
      <DatabaseGroupDataTable
        :database-group-list="dbGroupList"
        :show-selection="true"
        :single-selection="true"
        :show-external-link="true"
        :selected-database-group-names="
          databaseSelectState.selectedDatabaseGroup
            ? [databaseSelectState.selectedDatabaseGroup]
            : []
        "
        @update:selected-database-group-names="
          databaseSelectState.selectedDatabaseGroup = head($event)
        "
      />
    </NTabPane>
  </NTabs>
</template>

<script lang="ts" setup>
import { head } from "lodash-es";
import { NTabs, NTabPane } from "naive-ui";
import { computed, reactive, ref, watch } from "vue";
import AdvancedSearch from "@/components/AdvancedSearch";
import DatabaseGroupDataTable from "@/components/DatabaseGroup/DatabaseGroupDataTable.vue";
import { PagedDatabaseTable } from "@/components/v2/Model/DatabaseV1Table";
import { useDBGroupListByProject } from "@/store";
import {
  instanceNamePrefix,
  environmentNamePrefix,
} from "@/store/modules/v1/common";
import { type ComposedProject } from "@/types";
import {
  CommonFilterScopeIdList,
  extractProjectResourceName,
  type SearchParams,
  type SearchScope,
} from "@/utils";
import { useCommonSearchScopeOptions } from "../AdvancedSearch/useCommonSearchScopeOptions";
import type { DatabaseSelectState } from "./types";

const props = defineProps<{
  project: ComposedProject;
  value?: DatabaseSelectState;
}>();

const emit = defineEmits<{
  (event: "close"): void;
  (event: "update:value", state: DatabaseSelectState): void;
}>();

const readonlyScopes = computed((): SearchScope[] => [
  {
    id: "project",
    value: extractProjectResourceName(props.project.name),
    readonly: true,
  },
]);

const searchParams = ref<SearchParams>({
  query: "",
  scopes: [...readonlyScopes.value],
});

const databaseSelectState = reactive<DatabaseSelectState>(
  props.value || {
    changeSource: "DATABASE",
    selectedDatabaseNameList: [],
  }
);

const { dbGroupList } = useDBGroupListByProject(props.project.name);

const selectedInstance = computed(() => {
  const instanceId = searchParams.value.scopes.find(
    (scope) => scope.id === "instance"
  )?.value;
  if (!instanceId) {
    return;
  }
  return `${instanceNamePrefix}${instanceId}`;
});

const selectedEnvironment = computed(() => {
  const environmentId = searchParams.value.scopes.find(
    (scope) => scope.id === "environment"
  )?.value;
  if (!environmentId) {
    return;
  }
  return `${environmentNamePrefix}${environmentId}`;
});

const filter = computed(() => ({
  instance: selectedInstance.value,
  environment: selectedEnvironment.value,
  query: searchParams.value.query,
}));

const scopeOptions = useCommonSearchScopeOptions([...CommonFilterScopeIdList]);

watch(
  () => databaseSelectState,
  () => {
    emit("update:value", databaseSelectState);
  },
  { deep: true }
);
</script>
