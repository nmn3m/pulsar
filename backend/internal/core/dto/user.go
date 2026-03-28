package dto

type UpdateProfileRequest struct {
	FullName *string `json:"full_name,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Timezone *string `json:"timezone,omitempty"`
}
