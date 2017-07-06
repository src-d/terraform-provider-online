package online

type RPNv2Type string

const (
	Standard RPNv2Type = "STANDARD"
	QinQ     RPNv2Type = "QINQ"
	Demo     RPNv2Type = "DEMO"
)

type RPNv2 struct {
	ID                 int       `json:"id,omitempty"`
	Name               string    `json:"description"`
	Status             string    `json:"status,omitempty"`
	Type               RPNv2Type `json:"type"`
	CompatibilityRPNv1 bool      `json:"compatibility_rpn_v1"`
	Members            []*Member `json:"member,omitempty"`
}

func (r *RPNv2) MemberByServerID(id int) *Member {
	for _, m := range r.Members {
		if m.Linked.ID == id {
			return m
		}
	}

	return nil
}

type Member struct {
	ID     int `json:"id"`
	Linked struct {
		ID   int    `json:"id"`
		IP   string `json:"ip"`
		Type string `json:"type"`
		Ref  string `json:"$ref"`
	} `json:"linked"`
	Status string `json:"status"`
	VLAN   int    `json:"vlan"`
}
