import type { User } from './user';

export type RotationType = 'daily' | 'weekly' | 'custom';

export interface Schedule {
	id: string;
	organization_id: string;
	team_id?: string;
	name: string;
	description?: string;
	timezone: string;
	created_at: string;
	updated_at: string;
}

export interface ScheduleRotation {
	id: string;
	schedule_id: string;
	name: string;
	rotation_type: RotationType;
	rotation_length: number;
	start_date: string;
	start_time: string;
	end_time?: string;
	handoff_day?: number;
	handoff_time: string;
	created_at: string;
	updated_at: string;
}

export interface ScheduleRotationParticipant {
	id: string;
	rotation_id: string;
	user_id: string;
	position: number;
	created_at: string;
}

export interface ParticipantWithUser extends ScheduleRotationParticipant {
	user: User;
}

export interface ScheduleOverride {
	id: string;
	schedule_id: string;
	user_id: string;
	start_time: string;
	end_time: string;
	note?: string;
	created_at: string;
	updated_at: string;
}

export interface OnCallUser {
	user_id: string;
	user?: User;
	schedule_id: string;
	start_time: string;
	end_time: string;
	is_override: boolean;
}

export interface ScheduleWithRotations extends Schedule {
	rotations: ScheduleRotation[];
}

export interface RotationWithParticipants extends ScheduleRotation {
	participants: ParticipantWithUser[];
}

// Request types

export interface CreateScheduleRequest {
	team_id?: string;
	name: string;
	description?: string;
	timezone?: string;
}

export interface UpdateScheduleRequest {
	team_id?: string;
	name?: string;
	description?: string;
	timezone?: string;
}

export interface CreateRotationRequest {
	name: string;
	rotation_type: RotationType;
	rotation_length: number;
	start_date: string;
	start_time?: string;
	end_time?: string;
	handoff_day?: number;
	handoff_time?: string;
}

export interface UpdateRotationRequest {
	name?: string;
	rotation_type?: RotationType;
	rotation_length?: number;
	start_date?: string;
	start_time?: string;
	end_time?: string;
	handoff_day?: number;
	handoff_time?: string;
}

export interface AddParticipantRequest {
	user_id: string;
	position: number;
}

export interface ReorderParticipantsRequest {
	user_ids: string[];
}

export interface CreateOverrideRequest {
	user_id: string;
	start_time: string;
	end_time: string;
	note?: string;
}

export interface UpdateOverrideRequest {
	user_id?: string;
	start_time?: string;
	end_time?: string;
	note?: string;
}

// Response types

export interface ListSchedulesResponse {
	schedules: Schedule[];
}

export interface ListRotationsResponse {
	rotations: ScheduleRotation[];
}

export interface ListParticipantsResponse {
	participants: ParticipantWithUser[];
}

export interface ListOverridesResponse {
	overrides: ScheduleOverride[];
}
