package onemedical

type authRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	GrantType string `json:"grant_type"`
}

type authResponse struct {
	AccessToken string `json:"access_token"`
}

type appointmentRequest struct {
	AppointmentTypeID int64   `json:"appointment_type_id"`
	ServiceAreaID     int64   `json:"service_area_id"`
	OfficeIDs         []int64 `json:"office_ids"`
	ProviderID        *int64  `json:"provider_id"`
	OnsiteOnly        bool    `json:"onsite_only"`

	StartDate string `json:"date_start"`
	EndDate   string `json:"date_end"`
}

type appointmentResponse struct {
	ProviderCount               int64  `json:"provider_count"`
	InventoryCount              int64  `json:"inventory_count"`
	FirstAvailableInventoryDate string `json:"first_available_inventory_date"`

	Query appointmentRequest `json:"query"`
}
