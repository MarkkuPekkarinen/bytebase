import { cloneDeep, pullAt, without } from "lodash-es";
import { computed } from "vue";
import { useAppFeature } from "@/store";
import { DatabaseChangeMode } from "@/types/proto-es/v1/setting_service_pb";

export type SemanticIssueStatus = "OPEN" | "CLOSED";

export const UIIssueFilterScopeIdList = [
  "approver",
  "approval",
  "releaser",
] as const;
type UIIssueFilterScopeId = (typeof UIIssueFilterScopeIdList)[number];

export const CommonFilterScopeIdList = ["environment", "instance"] as const;
type CommonFilterScopeId = (typeof CommonFilterScopeIdList)[number];

export const AllSearchScopeIdList = [
  // common search scopes.
  "project",
  "environment",
  "instance",
  "database",
  "creator",
  "created",
  "updated",
  "status",
  "state",
  // database related search scopes.
  "engine",
  "database-label",
  "drifted",
  "table",
  // issue related search scopes.
  "issue-label",
  "taskType",
  // auditLog related search scopes.
  "method",
  "level",
  "actor",
  // instance related search scopes.
  "host",
  "port",
  // rollout related search scopes.
  "stage",
] as const;
export const useSearchScopeIdList = () => {
  const databaseChangeMode = useAppFeature("bb.feature.database-change-mode");
  return computed(() => {
    if (databaseChangeMode.value === DatabaseChangeMode.PIPELINE) {
      return AllSearchScopeIdList;
    }
    return without(AllSearchScopeIdList, "taskType");
  });
};

export type SearchScopeId =
  | (typeof AllSearchScopeIdList)[number]
  | UIIssueFilterScopeId
  | CommonFilterScopeId;

export type SearchScope = {
  id: SearchScopeId;
  value: string;
  readonly?: boolean;
};

export interface SearchParams {
  query: string;
  scopes: SearchScope[];
}

export const isValidSearchScopeId = (id: string): id is SearchScopeId => {
  return (
    AllSearchScopeIdList.includes(id as any) ||
    UIIssueFilterScopeIdList.includes(id as any) ||
    CommonFilterScopeIdList.includes(id as any)
  );
};

export const buildSearchTextBySearchParams = (
  params: SearchParams | undefined
): string => {
  const parts: string[] = [];
  params?.scopes.forEach((scope) => {
    parts.push(`${scope.id}:${encodeURIComponent(scope.value.trim())}`);
  });
  const query = (params?.query ?? "").trim();
  if (params?.query) {
    parts.push(encodeURIComponent(query));
  }
  return parts.join(" ");
};

export const mergeSearchParams = (base: SearchParams, patch: SearchParams) => {
  for (const scope of patch.scopes) {
    if (!base.scopes.find((s) => s.id === scope.id)) {
      base.scopes.push(scope);
    }
  }
  if (!base.query && patch.query) {
    base.query = patch.query;
  }
  return base;
};

export const buildSearchParamsBySearchText = (text: string): SearchParams => {
  const params = emptySearchParams();
  const segments = text.split(/\s+/g);
  const querySegments: string[] = [];

  for (let i = 0; i < segments.length; i++) {
    const seg = segments[i];
    const parts = seg.split(":");
    if (parts.length === 2 && isValidSearchScopeId(parts[0]) && parts[1]) {
      params.scopes.push({
        id: parts[0],
        value: decodeURIComponent(parts[1]),
      });
    } else {
      querySegments.push(decodeURIComponent(seg));
    }
  }
  params.query = querySegments.join(" ");
  params.scopes = params.scopes.filter((scope) => scope.id && scope.value);

  return params;
};

export const getValueFromSearchParams = (
  params: SearchParams,
  scopeId: SearchScopeId,
  prefix: string = "",
  validValues: readonly string[] = []
): string => {
  const scope = params.scopes.find((s) => s.id === scopeId);
  if (!scope) {
    return "";
  }
  const value = scope.value;
  if (validValues.length !== 0) {
    if (!validValues.includes(value)) {
      return "";
    }
  }
  return `${prefix}${scope.value}`;
};

export const getTsRangeFromSearchParams = (
  params: SearchParams,
  scopeId: SearchScopeId
) => {
  const scope = params.scopes.find((s) => s.id === scopeId);
  if (!scope) return undefined;
  const parts = scope.value.split(",");
  if (parts.length !== 2) return undefined;
  const range = [parseInt(parts[0], 10), parseInt(parts[1], 10)];
  return range as [number, number];
};

/**
 * @param scope will remove `scope` from `params.scopes` if `scope.value` is empty.
 * @param mutate true to mutate `params`. false to create a deep cloned copy. Default to false.
 * @returns `params` itself or a deep-cloned copy.
 */
export const upsertScope = ({
  params,
  scopes,
  mutate = false,
  allowMultiple = false,
}: {
  params: SearchParams;
  scopes: SearchScope | SearchScope[];
  mutate?: boolean;
  allowMultiple?: boolean;
}) => {
  const target = mutate ? params : cloneDeep(params);
  if (!Array.isArray(scopes)) {
    scopes = [scopes];
  }
  scopes.forEach((scope) => {
    const index = target.scopes.findIndex((s) => s.id === scope.id);
    if (index >= 0) {
      if (!scope.value) {
        pullAt(target.scopes, index);
      } else {
        if (allowMultiple) {
          target.scopes.push(scope);
        } else {
          target.scopes[index].value = scope.value;
        }
      }
    } else {
      if (scope.value) {
        target.scopes.push(scope);
      }
    }
  });
  return target;
};

export const emptySearchParams = (): SearchParams => {
  return {
    query: "",
    scopes: [],
  };
};
