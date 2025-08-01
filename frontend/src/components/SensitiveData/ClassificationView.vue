<template>
  <div class="w-full space-y-4">
    <div class="text-sm text-control-light">
      {{ $t("database.classification.description") }}
      <LearnMoreLink
        url="https://docs.bytebase.com/security/data-masking/data-classification?source=console"
        class="ml-1"
      />
    </div>
    <div>
      <div class="flex items-center space-x-2">
        <NSwitch
          :value="!state.classification.classificationFromConfig"
          :disabled="
            !allowEdit || !hasClassificationFeature || !hasClassificationConfig
          "
          @update:value="onClassificationConfigChange"
        />
        <div class="font-medium leading-7 text-main">
          {{ $t("database.classification.sync-from-comment") }}
        </div>
      </div>
      <i18n-t
        class="textinfolabel mt-1"
        tag="div"
        keypath="database.classification.sync-from-comment-tip"
      >
        <template #format>
          <span class="font-semibold">{classification id}-{comment}</span>
        </template>
      </i18n-t>
    </div>

    <NDivider class="my-2" />

    <div class="flex items-center justify-between">
      <div class="textinfolabel">
        {{ $t("settings.sensitive-data.classification.upload-label") }}
        <span
          class="normal-link cursor-pointer hover:underline"
          @click="state.showExampleModal = true"
        >
          {{ $t("settings.sensitive-data.view-example") }}
        </span>
      </div>

      <div class="flex items-center justify-end gap-2">
        <NButton
          type="primary"
          :disabled="!allowEdit || !hasClassificationFeature"
          @click="onUpload"
        >
          <template #icon>
            <UploadIcon class="h-4 w-4" />
          </template>
          {{ $t("settings.sensitive-data.classification.upload") }}
        </NButton>
        <input
          ref="uploader"
          type="file"
          accept=".json"
          class="sr-only hidden"
          :disabled="!allowEdit || !hasClassificationFeature"
          @input="onFileChange"
        />
      </div>
    </div>

    <div
      v-if="emptyConfig"
      class="flex justify-center border-2 border-gray-300 border-dashed rounded-md relative h-72"
    >
      <SingleFileSelector
        class="space-y-1 text-center flex flex-col justify-center items-center absolute top-0 bottom-0 left-0 right-0"
        :support-file-extensions="['.json']"
        :max-file-size-in-mi-b="maxFileSizeInMiB"
        :disabled="!allowEdit || !hasClassificationFeature"
        @on-select="onFileSelect"
      >
      </SingleFileSelector>
    </div>
    <div v-else class="h-full">
      <ClassificationTree :classification-config="state.classification" />
    </div>
  </div>

  <DataExampleModal
    v-if="state.showExampleModal"
    :example="JSON.stringify(example, null, 2)"
    @dismiss="state.showExampleModal = false"
  />
</template>

<script lang="ts" setup>
import { create } from "@bufbuild/protobuf";
import { head, isEqual, isEmpty } from "lodash-es";
import { UploadIcon } from "lucide-vue-next";
import { NSwitch, useDialog, NDivider, NButton } from "naive-ui";
import { v4 as uuidv4 } from "uuid";
import { computed, reactive, ref, watchEffect } from "vue";
import { useI18n } from "vue-i18n";
import { featureToRef, useSettingV1Store, pushNotification } from "@/store";
import type {
  DataClassificationSetting_DataClassificationConfig_Level as ClassificationLevel,
  DataClassificationSetting_DataClassificationConfig_DataClassification as DataClassification,
  DataClassificationSetting_DataClassificationConfig,
} from "@/types/proto-es/v1/setting_service_pb";
import {
  DataClassificationSetting_DataClassificationConfigSchema,
  DataClassificationSetting_DataClassificationConfig_LevelSchema,
  DataClassificationSetting_DataClassificationConfig_DataClassificationSchema,
  DataClassificationSettingSchema,
  Setting_SettingName,
  ValueSchema as SettingValueSchema,
} from "@/types/proto-es/v1/setting_service_pb";
import { PlanFeature } from "@/types/proto-es/v1/subscription_service_pb";
import { hasWorkspacePermissionV2 } from "@/utils";
import LearnMoreLink from "../LearnMoreLink.vue";
import ClassificationTree from "../SchemaTemplate/ClassificationTree.vue";
import SingleFileSelector from "../SingleFileSelector.vue";
import DataExampleModal from "./components/DataExampleModal.vue";

const uploader = ref<HTMLInputElement | null>(null);
const maxFileSizeInMiB = 10;

interface UploadClassificationConfig {
  title: string;
  levels: ClassificationLevel[];
  classification: { [key: string]: DataClassification };
}

interface LocalState {
  classification: DataClassificationSetting_DataClassificationConfig;
  showExampleModal: boolean;
}

const { t } = useI18n();
const $dialog = useDialog();
const settingStore = useSettingV1Store();

const formerConfig = computed(() => {
  const classification = head(settingStore.classification);
  return create(DataClassificationSetting_DataClassificationConfigSchema, {
    id: uuidv4(),
    title: classification?.title || "",
    levels: classification?.levels || [],
    classification: classification?.classification || {},
    classificationFromConfig: classification?.classificationFromConfig || false,
  });
});

const hasClassificationConfig = computed(
  () => settingStore.classification.length > 0
);

const state = reactive<LocalState>({
  showExampleModal: false,
  classification: create(
    DataClassificationSetting_DataClassificationConfigSchema,
    {
      id: uuidv4(),
      title: "",
      levels: [],
      classification: {},
      classificationFromConfig: false,
    }
  ),
});

// Initialize state with formerConfig
watchEffect(() => {
  const config = formerConfig.value;
  Object.assign(state.classification, {
    id: config.id,
    title: config.title,
    levels: config.levels,
    classification: config.classification,
    classificationFromConfig: config.classificationFromConfig,
  });
});

const emptyConfig = computed(
  () => Object.keys(state.classification.classification).length === 0
);

const allowSave = computed(() => {
  return (
    allowEdit.value &&
    hasClassificationFeature.value &&
    !isEqual(formerConfig.value, state.classification)
  );
});

const onClassificationConfigChange = (fromComment: boolean) => {
  $dialog.warning({
    title: t("common.warning"),
    content: fromComment
      ? t("database.classification.sync-from-comment-enable-warning")
      : t("database.classification.sync-from-comment-disable-warning"),
    style: "z-index: 100000",
    negativeText: t("common.cancel"),
    positiveText: t("common.confirm"),
    onPositiveClick: async () => {
      state.classification.classificationFromConfig = !fromComment;
      await upsertSetting();
    },
  });
};

const saveChanges = async () => {
  if (Object.keys(formerConfig.value.classification).length !== 0) {
    $dialog.warning({
      title: t("settings.sensitive-data.classification.override-title"),
      content: t("settings.sensitive-data.classification.override-desc"),
      style: "z-index: 100000",
      negativeText: t("common.cancel"),
      positiveText: t(
        "settings.sensitive-data.classification.override-confirm"
      ),
      onPositiveClick: async () => {
        await upsertSetting();
      },
    });
    return;
  }
  await upsertSetting();
};

const upsertSetting = async () => {
  await settingStore.upsertSetting({
    name: Setting_SettingName.DATA_CLASSIFICATION,
    value: create(SettingValueSchema, {
      value: {
        case: "dataClassificationSettingValue",
        value: create(DataClassificationSettingSchema, {
          configs: [state.classification],
        }),
      },
    }),
  });
  pushNotification({
    module: "bytebase",
    style: "SUCCESS",
    title: t("common.updated"),
  });
};

const allowEdit = computed(() => {
  return hasWorkspacePermissionV2("bb.settings.set");
});

const hasClassificationFeature = featureToRef(
  PlanFeature.FEATURE_DATA_CLASSIFICATION
);

const onUpload = () => {
  uploader.value?.click();
};

const onFileChange = () => {
  const files: File[] = (uploader.value as any).files;
  if (files.length !== 1) {
    return;
  }
  const file = files[0];
  if (file.size > maxFileSizeInMiB * 1024 * 1024) {
    pushNotification({
      module: "bytebase",
      style: "CRITICAL",
      title: t("common.file-selector.size-limit", {
        size: maxFileSizeInMiB,
      }),
    });
    return;
  }
  onFileSelect(file);
};

const onFileSelect = (file: File) => {
  const fr = new FileReader();
  fr.onload = () => {
    if (!fr.result) {
      return;
    }
    const data: UploadClassificationConfig = JSON.parse(fr.result as string);
    if (isEmpty(data.classification) || Array.isArray(data.classification)) {
      return pushNotification({
        module: "bytebase",
        style: "CRITICAL",
        title: "Data format error",
        description: `Should has the "classification" field. Please check the example.`,
      });
    }
    if (Object.keys(data.classification).length === 0) {
      return pushNotification({
        module: "bytebase",
        style: "CRITICAL",
        title: "Data format error",
        description: `"classification" field is empty. Please check the example.`,
      });
    }
    if (!Array.isArray(data.levels) || data.levels.length === 0) {
      return pushNotification({
        module: "bytebase",
        style: "CRITICAL",
        title: "Data format error",
        description: `Should has the "levels" field. Please check the example.`,
      });
    }
    Object.assign(state.classification, {
      title: data.title || state.classification.title || "",
      levels: data.levels.map((level) =>
        create(
          DataClassificationSetting_DataClassificationConfig_LevelSchema,
          level
        )
      ),
      classification: Object.values(data.classification).reduce(
        (map, data) => {
          map[data.id] = create(
            DataClassificationSetting_DataClassificationConfig_DataClassificationSchema,
            data
          );
          return map;
        },
        {} as { [key: string]: DataClassification }
      ),
    });
    if (isEqual(formerConfig.value, state.classification)) {
      return pushNotification({
        module: "bytebase",
        style: "INFO",
        title: "Nothing changed",
      });
    }
    if (allowSave.value) {
      saveChanges();
    }
  };
  fr.onerror = () => {
    pushNotification({
      module: "bytebase",
      style: "CRITICAL",
      title: "Read file error",
      description: String(fr.error),
    });
    return;
  };
  fr.readAsText(file);
};

const example = {
  title: "Classification Example",
  levels: [
    create(DataClassificationSetting_DataClassificationConfig_LevelSchema, {
      id: "1",
      title: "Level 1",
      description: "",
    }),
    create(DataClassificationSetting_DataClassificationConfig_LevelSchema, {
      id: "2",
      title: "Level 2",
      description: "",
    }),
    create(DataClassificationSetting_DataClassificationConfig_LevelSchema, {
      id: "3",
      title: "Level 3",
      description: "",
    }),
    create(DataClassificationSetting_DataClassificationConfig_LevelSchema, {
      id: "4",
      title: "Level 4",
      description: "",
    }),
  ],
  classification: {
    "1": create(
      DataClassificationSetting_DataClassificationConfig_DataClassificationSchema,
      {
        id: "1",
        title: "Basic",
        description: "",
      }
    ),
    "1-1": create(
      DataClassificationSetting_DataClassificationConfig_DataClassificationSchema,
      {
        id: "1-1",
        title: "Basic",
        description: "",
        levelId: "1",
      }
    ),
    "1-2": create(
      DataClassificationSetting_DataClassificationConfig_DataClassificationSchema,
      {
        id: "1-2",
        title: "Contact",
        description: "",
        levelId: "2",
      }
    ),
    "1-3": create(
      DataClassificationSetting_DataClassificationConfig_DataClassificationSchema,
      {
        id: "1-3",
        title: "Health",
        description: "",
        levelId: "4",
      }
    ),
    "2": create(
      DataClassificationSetting_DataClassificationConfig_DataClassificationSchema,
      {
        id: "2",
        title: "Relationship",
        description: "",
      }
    ),
    "2-1": create(
      DataClassificationSetting_DataClassificationConfig_DataClassificationSchema,
      {
        id: "2-1",
        title: "Social",
        description: "",
        levelId: "1",
      }
    ),
    "2-2": create(
      DataClassificationSetting_DataClassificationConfig_DataClassificationSchema,
      {
        id: "2-2",
        title: "Business",
        description: "",
        levelId: "3",
      }
    ),
  },
} satisfies UploadClassificationConfig;
</script>
