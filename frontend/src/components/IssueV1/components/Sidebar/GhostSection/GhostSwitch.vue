<template>
  <NTooltip :disabled="errors.length === 0">
    <template #trigger>
      <NSwitch
        :value="checked"
        :disabled="!allowChange"
        :loading="isUpdating"
        class="bb-ghost-switch"
        @update:value="toggleChecked"
      >
        <template #checked>
          <span style="font-size: 10px">{{ $t("common.on") }}</span>
        </template>
        <template #unchecked>
          <span style="font-size: 10px">{{ $t("common.off") }}</span>
        </template>
      </NSwitch>
    </template>
    <template #default>
      <ErrorList :errors="errors" />
    </template>
  </NTooltip>

  <InstanceAssignment
    v-if="showMissingInstanceLicense"
    :show="showInstanceAssignmentDrawer"
    @dismiss="showInstanceAssignmentDrawer = false"
  />
</template>

<script setup lang="tsx">
import { cloneDeep } from "lodash-es";
import { NSwitch, NTooltip } from "naive-ui";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";
import InstanceAssignment from "@/components/InstanceAssignment.vue";
import {
  databaseForTask,
  notifyNotEditableLegacyIssue,
  specForTask,
  useIssueContext,
} from "@/components/IssueV1/logic";
import type { ErrorItem } from "@/components/misc/ErrorList.vue";
import { default as ErrorList } from "@/components/misc/ErrorList.vue";
import { planServiceClient } from "@/grpcweb";
import { hasFeature, pushNotification, useIssueCommentStore } from "@/store";
import { Engine } from "@/types/proto/v1/common";
import { Plan_ChangeDatabaseConfig_Type } from "@/types/proto/v1/plan_service";
import { Task_Status } from "@/types/proto/v1/rollout_service";
import { engineNameV1, hasWorkspacePermissionV2 } from "@/utils";
import {
  allowGhostForTask,
  MIN_GHOST_SUPPORT_MARIADB_VERSION,
  MIN_GHOST_SUPPORT_MYSQL_VERSION,
  useIssueGhostContext,
} from "./common";

const { t } = useI18n();
const { isCreating, issue, events, selectedTask } = useIssueContext();
const { viewType, showFeatureModal, showMissingInstanceLicense } =
  useIssueGhostContext();
const isUpdating = ref(false);
const showInstanceAssignmentDrawer = ref(false);

const checked = computed(() => {
  return viewType.value === "ON";
});

const canManageSubscription = computed((): boolean => {
  return hasWorkspacePermissionV2("bb.settings.set");
});

const errors = computed(() => {
  const errors: ErrorItem[] = [];
  if (showMissingInstanceLicense.value && !canManageSubscription.value) {
    // Only show the tooltip when current user is not allowed to manage subscription
    // since we will show the InstanceAssignmentDrawer for high-privileged users
    // when clicking on the switch
    errors.push(
      t("subscription.instance-assignment.missing-license-attention")
    );
  }
  const database = databaseForTask(issue.value, selectedTask.value);
  // As we use the same database from backup to save temp tables in gh-ost, check if backup is available.
  if (!database.backupAvailable) {
    errors.push(
      t(
        "task.online-migration.error.not-applicable.needs-database-for-saving-temp-data",
        {
          // The same database name as backup.
          database: "bbdataarchive",
        }
      )
    );
  }
  if (!allowGhostForTask(issue.value, selectedTask.value)) {
    errors.push(
      t(
        "task.online-migration.error.not-applicable.task-doesnt-meet-ghost-requirement"
      )
    );
    errors.push({
      error: `${engineNameV1(Engine.MYSQL)} >= ${MIN_GHOST_SUPPORT_MYSQL_VERSION}, ${engineNameV1(Engine.MARIADB)} >= ${MIN_GHOST_SUPPORT_MARIADB_VERSION}`,
      indent: 1,
    });
  }
  return errors;
});

const allowChange = computed(() => {
  if (errors.value.length > 0) {
    return false;
  }
  // Always allow changing ghost status when creating.
  if (isCreating.value) {
    return true;
  }
  // Allow changing ghost status only when task is in one of the following states.
  if (
    [
      Task_Status.FAILED,
      Task_Status.CANCELED,
      Task_Status.NOT_STARTED,
    ].includes(selectedTask.value.status)
  ) {
    return true;
  }
  // Otherwise, disallow changing ghost status.
  return false;
});

const toggleChecked = async (on: boolean) => {
  if (!hasFeature("bb.feature.online-migration")) {
    showFeatureModal.value = true;
    return;
  }
  if (showMissingInstanceLicense.value) {
    if (canManageSubscription.value) {
      showInstanceAssignmentDrawer.value = true;
    }
    return;
  }
  if (errors.value.length > 0) {
    return;
  }

  if (isCreating.value) {
    const spec = specForTask(issue.value.planEntity, selectedTask.value);
    if (!spec || !spec.changeDatabaseConfig) return;
    spec.changeDatabaseConfig.type = on
      ? Plan_ChangeDatabaseConfig_Type.MIGRATE_GHOST
      : Plan_ChangeDatabaseConfig_Type.MIGRATE;
  } else {
    const planPatch = cloneDeep(issue.value.planEntity);
    const spec = specForTask(planPatch, selectedTask.value);
    if (!planPatch || !spec || !spec.changeDatabaseConfig) {
      notifyNotEditableLegacyIssue();
      return;
    }

    spec.changeDatabaseConfig.type = on
      ? Plan_ChangeDatabaseConfig_Type.MIGRATE_GHOST
      : Plan_ChangeDatabaseConfig_Type.MIGRATE;
    const updatedPlan = await planServiceClient.updatePlan({
      plan: planPatch,
      updateMask: ["steps"],
    });
    issue.value.planEntity = updatedPlan;

    const action = on ? "Enable" : "Disable";
    try {
      await useIssueCommentStore().createIssueComment({
        issueName: issue.value.name,
        comment: `${action} online migration for task [${selectedTask.value.target}].`,
      });
    } catch {
      // fail to comment won't be too bad.
    }

    events.emit("status-changed", { eager: true });
    pushNotification({
      module: "bytebase",
      style: "SUCCESS",
      title: t("common.updated"),
    });
  }
};
</script>

<style lang="postcss" scoped>
.bb-ghost-switch {
  --n-width: max(
    var(--n-rail-width),
    calc(var(--n-rail-width) + var(--n-button-width) - var(--n-rail-height))
  ) !important;
}
.bb-ghost-switch :deep(.n-switch__checked) {
  padding-right: calc(var(--n-rail-height) - var(--n-offset) + 1px);
}
.bb-ghost-switch :deep(.n-switch__unchecked) {
  padding-left: calc(var(--n-rail-height) - var(--n-offset) + 1px);
}
.bb-ghost-switch :deep(.n-switch__button-placeholder) {
  width: calc(1.25 * var(--n-rail-height));
}
</style>
