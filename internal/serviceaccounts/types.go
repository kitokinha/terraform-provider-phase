package serviceaccounts

// ServiceAccount is the Phase API representation used by both the summary and
// detail endpoints. Apps is populated by the detail endpoint.
type ServiceAccount struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      Role   `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Apps      []App  `json:"apps"`

	InitialToken *InitialToken `json:"initialToken,omitempty"`
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type App struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Environments []Environment `json:"environments"`
}

type Environment struct {
	ID string `json:"id"`
}

type InitialToken struct {
	ID          string `json:"id"`
	Token       string `json:"token"`
	BearerToken string `json:"bearerToken"`
}

type CreateRequest struct {
	Name      string `json:"name"`
	RoleID    string `json:"role_id"`
	TokenName string `json:"token_name,omitempty"`
	TeamID    string `json:"team_id,omitempty"`
}

type UpdateRequest struct {
	Name   *string `json:"name,omitempty"`
	RoleID *string `json:"role_id,omitempty"`
}

type AccessRequest struct {
	Apps []AccessApp `json:"apps"`
}

type AccessApp struct {
	ID           string   `json:"id"`
	Environments []string `json:"environments"`
}
