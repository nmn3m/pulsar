// Do Not Disturb (DND) Types

export interface DNDSettings {
  id?: string;
  user_id: string;
  enabled: boolean;
  schedule: DNDSchedule;
  overrides: DNDOverride[];
  allow_p1_override: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface DNDSchedule {
  weekly: DNDTimeSlot[];
  timezone: string;
}

export interface DNDTimeSlot {
  day: DayOfWeek;
  start: string; // HH:MM format (24-hour)
  end: string; // HH:MM format (24-hour)
}

export type DayOfWeek =
  | 'sunday'
  | 'monday'
  | 'tuesday'
  | 'wednesday'
  | 'thursday'
  | 'friday'
  | 'saturday';

export interface DNDOverride {
  start: string; // ISO date string
  end: string; // ISO date string
  reason?: string;
}

export interface UpdateDNDSettingsRequest {
  enabled?: boolean;
  schedule?: DNDSchedule;
  overrides?: DNDOverride[];
  allow_p1_override?: boolean;
}

export interface AddDNDOverrideRequest {
  start: string; // ISO date string
  end: string; // ISO date string
  reason?: string;
}

export interface DNDStatusResponse {
  in_dnd_mode: boolean;
  priority: string;
}

// Helper constants
export const DAYS_OF_WEEK: DayOfWeek[] = [
  'sunday',
  'monday',
  'tuesday',
  'wednesday',
  'thursday',
  'friday',
  'saturday',
];

export const DAY_LABELS: Record<DayOfWeek, string> = {
  sunday: 'Sunday',
  monday: 'Monday',
  tuesday: 'Tuesday',
  wednesday: 'Wednesday',
  thursday: 'Thursday',
  friday: 'Friday',
  saturday: 'Saturday',
};

// Common timezone options
export const COMMON_TIMEZONES = [
  'UTC',
  'America/New_York',
  'America/Chicago',
  'America/Denver',
  'America/Los_Angeles',
  'America/Toronto',
  'America/Vancouver',
  'Europe/London',
  'Europe/Paris',
  'Europe/Berlin',
  'Europe/Amsterdam',
  'Asia/Tokyo',
  'Asia/Shanghai',
  'Asia/Singapore',
  'Asia/Dubai',
  'Asia/Kolkata',
  'Australia/Sydney',
  'Australia/Melbourne',
  'Pacific/Auckland',
];
