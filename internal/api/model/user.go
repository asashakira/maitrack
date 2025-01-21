package model

import (
	database "github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	UserID          uuid.UUID        `json:"userID,omitempty"`
	Username        string           `json:"username,omitempty"`
	Password        string           `json:"password,omitempty"`
	SegaID          string           `json:"segaID,omitempty"`
	SegaPassword    string           `json:"segaPassword,omitempty"`
	GameName        string           `json:"gameName,omitempty"`
	TagLine         string           `json:"tagLine,omitempty"`
	Rating          int32            `json:"rating,omitempty"`
	SeasonPlayCount int32            `json:"seasonPlayCount,omitempty"`
	TotalPlayCount  int32            `json:"totalPlayCount,omitempty"`
	CreatedAt       pgtype.Timestamp `json:"createdAt,omitempty"`
	UpdatedAt       pgtype.Timestamp `json:"updatedAt,omitempty"`
}

func ConvertUsers(dbUsers []database.User, dbUserDatas []database.UserDatum) []User {
	users := []User{}
	for i := 0; i < len(dbUsers); i++ {
		users = append(users, ConvertUser(dbUsers[i], dbUserDatas[i]))
	}
	return users
}

func ConvertUser(dbUser database.User, dbUserData database.UserDatum) User {
	return User{
		UserID:          dbUser.UserID,
		Username:        dbUser.Username,
		Password:        dbUser.Password,
		SegaID:          dbUser.SegaID,
		SegaPassword:    dbUser.SegaPassword,
		GameName:        dbUser.GameName,
		TagLine:         dbUser.TagLine,
		Rating:          dbUserData.Rating,
		SeasonPlayCount: dbUserData.SeasonPlayCount,
		TotalPlayCount:  dbUserData.TotalPlayCount,
		CreatedAt:       dbUser.CreatedAt,
		UpdatedAt:       dbUser.UpdatedAt,
	}
}
