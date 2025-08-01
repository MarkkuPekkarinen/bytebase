<template>
  <div class="w-full max-h-[calc(100vh-10rem)]">
    <div v-if="isValidProjectName(project.name)">
      <NButton size="small" text @click="gotoWorkspace">
        <template #icon>
          <ChevronLeftIcon class="w-4 opacity-80" />
        </template>
        {{ $t("common.back-to-workspace") }}
      </NButton>
    </div>
    <NTabs
      :value="actualSelectedTab"
      type="line"
      @update:value="state.selectedTab = $event"
    >
      <template #suffix>
        <div class="flex flex-row justify-end items-center gap-x-2">
          <SearchBox
            v-model:value="state.searchText"
            :placeholder="$t('common.search-by-name')"
            :autofocus="false"
            class="!w-40"
            size="small"
          />
          <NTooltip v-if="allowToCreateProject" trigger="hover">
            <template #trigger>
              <NButton size="small" @click="$emit('on-create')">
                <template #icon>
                  <PlusIcon class="w-4 h-auto" />
                </template>
              </NButton>
            </template>
            {{ $t("quick-action.new-project") }}
          </NTooltip>
        </div>
      </template>
      <NTabPane
        v-for="tab in tabList"
        :key="tab.id"
        :name="tab.id"
        :tab="tab.title"
        :disabled="
          tab.id === 'recent' &&
          state.searchText.trim().length > 0 &&
          filteredRecentProjectList.length === 0
        "
      >
        <ProjectV1Table
          v-if="tab.id === 'recent'"
          :project-list="tab.list"
          :current-project="
            isValidProjectName(project.name) ? project : undefined
          "
          :keyword="state.searchText"
          @row-click="onProjectSelect"
        />
        <PagedProjectTable
          v-else
          class="mb-2"
          session-key="bb.project-table"
          :filter="filter"
          :loading="state.loading"
          @row-click="onProjectSelect"
        />
      </NTabPane>
    </NTabs>
  </div>
</template>

<script lang="ts" setup>
import { ChevronLeftIcon, PlusIcon } from "lucide-vue-next";
import { NButton, NTabPane, NTabs, NTooltip } from "naive-ui";
import { computed, reactive, onMounted, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import { useRecentProjects } from "@/components/Project/useRecentProjects";
import { SearchBox, ProjectV1Table, PagedProjectTable } from "@/components/v2";
import { PROJECT_V1_ROUTE_DETAIL } from "@/router/dashboard/projectV1";
import { WORKSPACE_ROUTE_LANDING } from "@/router/dashboard/workspaceRoutes";
import { useRecentVisit } from "@/router/useRecentVisit";
import { useCurrentProjectV1 } from "@/store";
import { getProjectName } from "@/store/modules/v1/common";
import type { ComposedProject } from "@/types";
import { isValidProjectName, DEFAULT_PROJECT_NAME } from "@/types";
import {
  filterProjectV1ListByKeyword,
  hasWorkspacePermissionV2,
} from "@/utils";

type LocalTabType = "recent" | "all";

interface LocalState {
  showPopover: boolean;
  searchText: string;
  loading: boolean;
  selectedTab: LocalTabType;
}

defineEmits<{
  (event: "on-create"): void;
}>();

const { t } = useI18n();
const state = reactive<LocalState>({
  showPopover: false,
  searchText: "",
  selectedTab: "all",
  loading: true,
});

const { recentViewProjects } = useRecentProjects();
const router = useRouter();
const { record } = useRecentVisit();

const { project } = useCurrentProjectV1();

const filter = computed(() => ({
  query: state.searchText,
  excludeDefault: true,
}));

onMounted(() => {
  state.selectedTab = recentViewProjects.value.length < 1 ? "all" : "recent";
});

const getFilteredProjectList = (
  projectList: ComposedProject[]
): ComposedProject[] => {
  const list = projectList.filter(
    (project) => project.name !== DEFAULT_PROJECT_NAME
  );
  return filterProjectV1ListByKeyword(list, state.searchText);
};

const filteredRecentProjectList = computed(() => {
  return getFilteredProjectList(recentViewProjects.value);
});

const allowToCreateProject = computed(() =>
  hasWorkspacePermissionV2("bb.projects.create")
);

const tabList = computed(
  (): { title: string; id: LocalTabType; list: ComposedProject[] }[] => [
    {
      title: t("common.recent"),
      id: "recent",
      list: filteredRecentProjectList.value,
    },
    {
      title: t("common.all"),
      id: "all",
      list: [],
    },
  ]
);

const actualSelectedTab = computed((): LocalState["selectedTab"] => {
  if (
    state.searchText.trim().length > 0 &&
    filteredRecentProjectList.value.length === 0
  ) {
    // Force to view 'ALL' tab when search by keyword but "Recent" is empty.
    return "all";
  }
  return state.selectedTab;
});

const onProjectSelect = (project: ComposedProject) => {
  const route = router.resolve({
    name: PROJECT_V1_ROUTE_DETAIL,
    params: {
      projectId: getProjectName(project.name),
    },
  });
  record(route.fullPath);
};

const gotoWorkspace = (e: MouseEvent) => {
  const route = router.resolve({
    name: WORKSPACE_ROUTE_LANDING,
  });
  record(route.fullPath);
  if (e.ctrlKey || e.metaKey) {
    window.open(route.fullPath, "_blank");
  } else {
    router.push(route.fullPath);
  }
};

// Close popover when current project changed.
watch(
  () => project.value.name,
  () => {
    state.showPopover = false;
  }
);
</script>
