package dto

type TeamResponse struct {
	Team TeamInfo `json:"team"`
}

type TeamInfo struct {
	TeamName string      `json:"team_name"`
	Members  []MemberDto `json:"members"`
}
