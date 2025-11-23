package dto

import "Git_PR_plugin_avito_tech/internal/model"

func ToTeamInfo(team *model.Team, members []MemberDto) *TeamInfo {
	if team == nil {
		return nil
	}

	return &TeamInfo{
		TeamName: team.TeamName,
		Members:  members,
	}
}

func ToTeamResponse(team *model.Team, users []*model.User) *TeamResponse {
	if team == nil {
		return nil
	}

	members := make([]MemberDto, 0, len(users))
	for _, user := range users {
		if user == nil {
			continue
		}

		members = append(members, MemberDto{
			UserID:   user.UserID,
			Username: user.Username,
			IsActive: user.IsActive,
		})
	}

	return &TeamResponse{Team: *ToTeamInfo(team, members)}
}
