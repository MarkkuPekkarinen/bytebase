<template>
  <NPopover
    v-if="show"
    placement="bottom-end"
    trigger="click"
    :disabled="disabled"
  >
    <template #trigger>
      <NButton
        :disabled="disabled"
        type="primary"
        size="small"
        style="--n-padding: 0 0.25rem"
      >
        <template #icon>
          <ChevronDown />
        </template>
      </NButton>
    </template>
    <template #default>
      <div class="flex flex-col gap-1">
        <div>
          <p class="mb-1 textinfolabel">
            {{ $t("data-source.select-query-data-source") }}
          </p>
          <NRadioGroup
            class="max-w-44"
            :value="selectedDataSourceId"
            @update:value="onDataSourceSelected"
          >
            <NTooltip
              v-for="ds in dataSources"
              :key="ds.id"
              :disabled="!Boolean(dataSourceUnaccessibleReason(ds))"
            >
              <template #trigger>
                <NRadio
                  class="w-full"
                  :value="ds.id"
                  :disabled="Boolean(dataSourceUnaccessibleReason(ds))"
                >
                  <div
                    class="max-w-36 flex flex-row justify-start items-center"
                  >
                    <span class="text-xs opacity-60 shrink-0">{{
                      readableDataSourceType(ds.type)
                    }}</span>
                    <span class="ml-1 truncate">{{ ds.username }}</span>
                  </div>
                </NRadio>
              </template>
              <p class="text-nowrap">
                {{ dataSourceUnaccessibleReason(ds) }}
              </p>
            </NTooltip>
          </NRadioGroup>
        </div>
        <NTooltip
          v-if="showRedisConfig"
          :disabled="
            selectedDataSource?.redisType === DataSource_RedisType.CLUSTER
          "
        >
          <template #trigger>
            <div class="border-t pt-1" style="">
              <p class="mb-1 textinfolabel">
                {{ $t("sql-editor.redis-command.self") }}
              </p>
              <NRadioGroup
                :disabled="
                  selectedDataSource?.redisType !== DataSource_RedisType.CLUSTER
                "
                v-model:value="redisCommandOption"
                class="max-w-44"
              >
                <NRadio :value="QueryOption_RedisRunCommandsOn.SINGLE_NODE">
                  {{ $t("sql-editor.redis-command.single-node") }}
                </NRadio>
                <NRadio :value="QueryOption_RedisRunCommandsOn.ALL_NODES">
                  {{ $t("sql-editor.redis-command.all-nodes") }}
                </NRadio>
              </NRadioGroup>
            </div>
          </template>
          {{ $t("sql-editor.redis-command.only-for-cluster") }}
        </NTooltip>
        <div class="border-t pt-1">
          <ResultLimitSelect
            placement="right-start"
            trigger="hover"
            :maximum="maximumResultRows"
          >
            <template
              #default="{ resultRowsLimit }: { resultRowsLimit: number }"
            >
              <NButton
                icon-placement="right"
                quaternary
                style="justify-content: start; --n-padding: 0 8px; width: 100%"
              >
                {{ $t("sql-editor.result-limit.self") }}
                {{
                  $t("sql-editor.result-limit.n-rows", { n: resultRowsLimit })
                }}
                <template #icon>
                  <ChevronRight />
                </template>
              </NButton>
            </template>
          </ResultLimitSelect>
        </div>
      </div>
    </template>
  </NPopover>
</template>

<script lang="ts" setup>
import { orderBy } from "lodash-es";
import { ChevronDown, ChevronRight } from "lucide-vue-next";
import { NButton, NPopover, NRadioGroup, NRadio, NTooltip } from "naive-ui";
import { storeToRefs } from "pinia";
import { computed, watch, watchEffect } from "vue";
import { useI18n } from "vue-i18n";
import {
  useConnectionOfCurrentSQLEditorTab,
  usePolicyV1Store,
  useSQLEditorTabStore,
  useSQLEditorStore,
  useSettingV1Store,
} from "@/store";
import { isValidDatabaseName } from "@/types";
import { Engine } from "@/types/proto-es/v1/common_pb";
import type { DataSource } from "@/types/proto-es/v1/instance_service_pb";
import {
  DataSourceType,
  DataSource_RedisType,
} from "@/types/proto-es/v1/instance_service_pb";
import {
  DataSourceQueryPolicy_Restriction,
  PolicyType,
} from "@/types/proto-es/v1/org_policy_service_pb";
import { Setting_SettingName } from "@/types/proto-es/v1/setting_service_pb";
import { QueryOption_RedisRunCommandsOn } from "@/types/proto-es/v1/sql_service_pb";
import { getValidDataSourceByPolicy, readableDataSourceType } from "@/utils";
import { getAdminDataSourceRestrictionOfDatabase } from "@/utils";
import ResultLimitSelect from "./ResultLimitSelect.vue";

defineProps<{
  disabled?: boolean;
}>();

const { t } = useI18n();
const tabStore = useSQLEditorTabStore();
const { connection, database } = useConnectionOfCurrentSQLEditorTab();
const policyStore = usePolicyV1Store();
const settingV1Store = useSettingV1Store();

const { redisCommandOption, resultRowsLimit } =
  storeToRefs(useSQLEditorStore());

const show = computed(() => {
  return tabStore.currentTab?.mode !== "ADMIN";
});

const showRedisConfig = computed(() => {
  if (!database.value) {
    return false;
  }
  const instance = database.value.instanceResource;
  return instance.engine === Engine.REDIS;
});

const adminDataSourceRestriction = computed(() => {
  if (!database.value) {
    return {
      environmentPolicy:
        DataSourceQueryPolicy_Restriction.RESTRICTION_UNSPECIFIED,
      projectPolicy: DataSourceQueryPolicy_Restriction.RESTRICTION_UNSPECIFIED,
    };
  }
  return getAdminDataSourceRestrictionOfDatabase(database.value);
});

const selectedDataSourceId = computed(() => {
  return connection.value.dataSourceId;
});

const selectedDataSource = computed(() => {
  const instance = database.value.instanceResource;
  return instance.dataSources.find(
    (ds) => ds.id === selectedDataSourceId.value
  );
});

const dataSources = computed(() => {
  return orderBy(database.value.instanceResource.dataSources, "type");
});

const dataSourceUnaccessibleReason = (
  dataSource: DataSource
): string | undefined => {
  if (dataSource.type === DataSourceType.ADMIN) {
    if (
      adminDataSourceRestriction.value.environmentPolicy ===
        DataSourceQueryPolicy_Restriction.DISALLOW ||
      adminDataSourceRestriction.value.projectPolicy ===
        DataSourceQueryPolicy_Restriction.DISALLOW
    ) {
      return t(
        "sql-editor.query-context.admin-data-source-is-disallowed-to-query"
      );
    }
    const readOnlyDataSources = dataSources.value.filter(
      (ds) => ds.type === DataSourceType.READ_ONLY
    );
    if (
      readOnlyDataSources.length > 0 &&
      (adminDataSourceRestriction.value.environmentPolicy ===
        DataSourceQueryPolicy_Restriction.FALLBACK ||
        adminDataSourceRestriction.value.projectPolicy ===
          DataSourceQueryPolicy_Restriction.FALLBACK)
    ) {
      return t(
        "sql-editor.query-context.admin-data-source-is-disallowed-to-query-when-read-only-data-source-is-available"
      );
    }
  }

  return undefined;
};

const onDataSourceSelected = (dataSourceId?: string) => {
  tabStore.updateCurrentTab({
    connection: {
      ...connection.value,
      dataSourceId: dataSourceId,
    },
  });
};

watch(
  [() => selectedDataSourceId.value, () => database.value],
  ([current, database]) => {
    if (!isValidDatabaseName(database.name)) return;
    if (!current) {
      const fixed = getValidDataSourceByPolicy(database);
      onDataSourceSelected(fixed);
    }
  },
  {
    immediate: true,
  }
);

watch(
  () => database.value.name,
  async () => {
    if (!isValidDatabaseName(database.value.name)) {
      return;
    }
    await policyStore.getOrFetchPolicyByParentAndType({
      parentPath: database.value.effectiveEnvironment,
      policyType: PolicyType.DATA_SOURCE_QUERY,
    });
    await policyStore.getOrFetchPolicyByParentAndType({
      parentPath: database.value.project,
      policyType: PolicyType.DATA_SOURCE_QUERY,
    });
  },
  {
    immediate: true,
  }
);

const maximumResultRows = computed(() => {
  const setting = settingV1Store.getSettingByName(
    Setting_SettingName.SQL_RESULT_SIZE_LIMIT
  );
  if (setting?.value?.value?.case === "sqlQueryRestrictionSetting") {
    const limit = setting.value.value.value.maximumResultRows ?? -1;
    if (limit <= 0) {
      return Number.MAX_VALUE;
    }
    return limit;
  }
  return Number.MAX_VALUE;
});

watchEffect(() => {
  if (resultRowsLimit.value > maximumResultRows.value) {
    resultRowsLimit.value = maximumResultRows.value;
  }
});
</script>
