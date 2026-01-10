export type EscalationTargetType = 'user' | 'team' | 'schedule';

export interface EscalationPolicy {
  id: string;
  organization_id: string;
  name: string;
  description?: string;
  repeat_enabled: boolean;
  repeat_count?: number;
  created_at: string;
  updated_at: string;
}

export interface EscalationRule {
  id: string;
  policy_id: string;
  position: number;
  escalation_delay: number;
  created_at: string;
  updated_at: string;
}

export interface EscalationTarget {
  id: string;
  rule_id: string;
  target_type: EscalationTargetType;
  target_id: string;
  notification_channels?: TargetNotificationConfig;
  created_at: string;
}

// Notification channel override configuration for escalation targets
export interface TargetNotificationConfig {
  channels: string[]; // e.g., ["email", "slack", "sms", "webhook"]
  urgent?: boolean; // If true, use urgent/high-priority notification
}

export interface EscalationRuleWithTargets extends EscalationRule {
  targets: EscalationTarget[];
}

export interface EscalationPolicyWithRules extends EscalationPolicy {
  rules: EscalationRuleWithTargets[];
}

// Request types

export interface CreateEscalationPolicyRequest {
  name: string;
  description?: string;
  repeat_enabled?: boolean;
  repeat_count?: number;
}

export interface UpdateEscalationPolicyRequest {
  name?: string;
  description?: string;
  repeat_enabled?: boolean;
  repeat_count?: number;
}

export interface CreateEscalationRuleRequest {
  position: number;
  escalation_delay: number;
}

export interface UpdateEscalationRuleRequest {
  position?: number;
  escalation_delay?: number;
}

export interface AddEscalationTargetRequest {
  target_type: EscalationTargetType;
  target_id: string;
  notification_channels?: TargetNotificationConfig;
}

// Response types

export interface ListEscalationPoliciesResponse {
  policies: EscalationPolicy[];
}

export interface ListEscalationRulesResponse {
  rules: EscalationRule[];
}

export interface ListEscalationTargetsResponse {
  targets: EscalationTarget[];
}
