import type { User } from './user';

export interface APIKey {
  id: string;
  organization_id: string;
  user_id: string;
  user?: User;
  name: string;
  key_prefix: string;
  scopes: string[];
  is_active: boolean;
  last_used_at?: string;
  expires_at?: string;
  created_at: string;
  updated_at: string;
}

export interface APIKeyResponse extends APIKey {
  key: string; // Full key only returned on creation
}

export interface CreateAPIKeyRequest {
  name: string;
  scopes: string[];
  expires_at?: string; // RFC3339 format
}

export interface UpdateAPIKeyRequest {
  name?: string;
  is_active?: boolean;
}

export interface ListAPIKeysResponse {
  api_keys: APIKey[];
}

export interface APIKeyScope {
  name: string;
  description: string;
}
