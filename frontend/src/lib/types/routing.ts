export interface RoutingCondition {
  field: string; // source, priority, tags, message, or custom field
  operator:
    | 'equals'
    | 'not_equals'
    | 'contains'
    | 'not_contains'
    | 'regex'
    | 'gte'
    | 'lte'
    | 'starts_with'
    | 'ends_with';
  value: string;
}

export interface RoutingConditions {
  match: 'all' | 'any';
  conditions: RoutingCondition[];
}

export interface RoutingActions {
  assign_team_id?: string;
  assign_user_id?: string;
  assign_escalation_policy_id?: string;
  set_priority?: string;
  add_tags?: string[];
  suppress?: boolean;
}

export interface AlertRoutingRule {
  id: string;
  organization_id: string;
  name: string;
  description?: string;
  priority: number;
  conditions: RoutingConditions;
  actions: RoutingActions;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateRoutingRuleRequest {
  name: string;
  description?: string;
  priority?: number;
  conditions: RoutingConditions;
  actions: RoutingActions;
  enabled?: boolean;
}

export interface UpdateRoutingRuleRequest {
  name?: string;
  description?: string;
  priority?: number;
  conditions?: RoutingConditions;
  actions?: RoutingActions;
  enabled?: boolean;
}

export interface ReorderRoutingRulesRequest {
  rule_ids: string[];
}

export const ROUTING_FIELDS = [
  { value: 'source', label: 'Source' },
  { value: 'priority', label: 'Priority' },
  { value: 'message', label: 'Message' },
  { value: 'tags', label: 'Tags' },
] as const;

export const ROUTING_OPERATORS = [
  { value: 'equals', label: 'Equals' },
  { value: 'not_equals', label: 'Not Equals' },
  { value: 'contains', label: 'Contains' },
  { value: 'not_contains', label: 'Does Not Contain' },
  { value: 'regex', label: 'Matches Regex' },
  { value: 'starts_with', label: 'Starts With' },
  { value: 'ends_with', label: 'Ends With' },
] as const;

export const PRIORITY_OPTIONS = [
  { value: 'P1', label: 'P1 - Critical' },
  { value: 'P2', label: 'P2 - High' },
  { value: 'P3', label: 'P3 - Medium' },
  { value: 'P4', label: 'P4 - Low' },
  { value: 'P5', label: 'P5 - Info' },
] as const;
