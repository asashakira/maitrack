package scraper

import (
	"context"
	"log"

	"github.com/asashakira/mai.gg/internal/database"
	"github.com/google/uuid"
)

func StartScraping(db *database.Queries) error {
	return nil
	users, err := db.GetAllUsers(context.Background())
	if err != nil {
		log.Println("GetAllUsers error: ", err)
		return err
	}

	for _, user := range users {
		segaID := user.SegaID
		password := user.Password
		userID := user.UserID

		// set up maimaiClient
		maimaiClient := New()
		err = maimaiClient.Login(segaID, password)
		if err != nil {
			log.Println("maimai Login error: ", err)
			return err
		}

		newUserData, err := scrapeUserData(maimaiClient)
		if err != nil {
			log.Println("Error getting user: ", err)
			return err
		}
		_, err = db.UpdateUser(context.Background(), database.UpdateUserParams{
			UserID:   userID,
			SegaID:   segaID,
			Password: password,
			GameName: newUserData.GameName,
			TagLine:  newUserData.TagLine,
		})
		if err != nil {
			log.Println("UpdateUser error while scraping: ", err)
			return err
		}

		err = db.InsertUserData(context.Background(), database.InsertUserDataParams{
			ID:              uuid.New(),
			UserID:          userID,
			GameName:        newUserData.GameName,
			TagLine:         newUserData.TagLine,
			Rating:          newUserData.Rating,
			SeasonPlayCount: newUserData.SeasonPlayCount,
			TotalPlayCount:  newUserData.TotalPlayCount,
		})
		if err != nil {
			log.Println("InsertUserData error while scraping: ", err)
			return err
		}
		log.Println("Done scraping ", user.GameName)
	}

	return nil
}
