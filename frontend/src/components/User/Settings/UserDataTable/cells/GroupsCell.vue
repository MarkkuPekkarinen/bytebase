<template>
  <div class="flex flex-row items-center flex-wrap gap-2">
    <NTag
      v-for="group in groups"
      :key="group.name"
      class="!cursor-pointer hover:bg-gray-200"
      @click="$emit('select-group', group)"
    >
      {{ group.title }}
    </NTag>
    <span v-if="groups.length === 0">-</span>
  </div>
</template>

<script lang="ts" setup>
import { NTag } from "naive-ui";
import { computed } from "vue";
import { useGroupList } from "@/store";
import { extractUserId } from "@/store/modules/v1/common";
import type { Group } from "@/types/proto-es/v1/group_service_pb";
import { type User } from "@/types/proto-es/v1/user_service_pb";

const props = defineProps<{
  user: User;
}>();

defineEmits<{
  (event: "select-group", group: Group): void;
}>();

const groupList = useGroupList();

const groups = computed(() => {
  const groups = [];
  for (const group of groupList.value) {
    for (const member of group.members) {
      if (extractUserId(member.member) === props.user.email) {
        groups.push(group);
        break;
      }
    }
  }
  return groups;
});
</script>
