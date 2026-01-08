import type { User } from './user';

export type TeamRole = 'lead' | 'member';

export interface Team {
  id: string;
  organization_id: string;
  name: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface TeamMember {
  team_id: string;
  user_id: string;
  role: TeamRole;
  joined_at: string;
}

export interface UserWithTeamRole extends User {
  role: TeamRole;
  joined_at: string;
}

export interface TeamWithMembers extends Team {
  members: UserWithTeamRole[];
}

export interface CreateTeamRequest {
  name: string;
  description?: string;
}

export interface UpdateTeamRequest {
  name?: string;
  description?: string;
}

export interface AddTeamMemberRequest {
  user_id: string;
  role?: TeamRole;
}

export interface UpdateTeamMemberRoleRequest {
  role: TeamRole;
}
