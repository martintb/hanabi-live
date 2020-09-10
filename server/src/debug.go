package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	badGameIDs = make([]int, 0)
)

func debugPrint() {
	tablesMutex.RLock()
	defer tablesMutex.RUnlock()

	logger.Debug("---------------------------------------------------------------")
	logger.Debug("Current total tables:", len(tables))

	numUnstarted := 0
	numRunning := 0
	numReplays := 0

	for _, t := range tables { // This is a map[int]*Table
		if !t.Running {
			numUnstarted++
		}

		if t.Running && !t.Replay {
			numRunning++
		}

		if t.Replay {
			numReplays++
		}
	}

	logger.Debug("Current unstarted tables:", numUnstarted)
	logger.Debug("Current ongoing tables:", numRunning)
	logger.Debug("Current replays:", numReplays)

	logger.Debug("---------------------------------------------------------------")
	logger.Debug("Current table list:")
	logger.Debug("---------------------------------------------------------------")

	// Print out all of the current tables
	if len(tables) == 0 {
		logger.Debug("[no current tables]")
	}

	for tableID, t := range tables { // This is a map[int]*Table
		logger.Debug(strconv.FormatUint(tableID, 10) + " - " + t.Name)
		logger.Debug("\n")

		// Print out all of the fields
		// https://stackoverflow.com/questions/24512112/how-to-print-struct-variables-in-console
		logger.Debug("    All fields:")
		fieldsToIgnore := []string{
			"Players",
			"Spectators",
			"DisconSpectators",

			"Options",
		}
		s := reflect.ValueOf(t).Elem()
		maxChars := 0
		for i := 0; i < s.NumField(); i++ {
			fieldName := s.Type().Field(i).Name
			if stringInSlice(fieldName, fieldsToIgnore) {
				continue
			}
			if len(fieldName) > maxChars {
				maxChars = len(fieldName)
			}
		}
		for i := 0; i < s.NumField(); i++ {
			fieldName := s.Type().Field(i).Name
			if stringInSlice(fieldName, fieldsToIgnore) {
				continue
			}
			f := s.Field(i)
			line := "  "
			for i := len(fieldName); i < maxChars; i++ {
				line += " "
			}
			line += "%s = %v"
			line = fmt.Sprintf(line, fieldName, f.Interface())
			if strings.HasSuffix(line, " = ") {
				line += "[empty string]"
			}
			line += "\n"
			logger.Debug(line)
		}
		logger.Debug("\n")

		// Manually enumerate the slices and maps
		logger.Debug("    Options:")
		if t.Options == nil {
			logger.Debug("      [Options is nil; this should never happen]")
		} else {
			s2 := reflect.ValueOf(t.Options).Elem()
			maxChars2 := 0
			for i := 0; i < s2.NumField(); i++ {
				fieldName := s2.Type().Field(i).Name
				if len(fieldName) > maxChars2 {
					maxChars2 = len(fieldName)
				}
			}
			for i := 0; i < s2.NumField(); i++ {
				fieldName := s2.Type().Field(i).Name
				f := s2.Field(i)
				line := "    "
				for i := len(fieldName); i < maxChars2; i++ {
					line += " "
				}
				line += "%s = %v"
				line = fmt.Sprintf(line, fieldName, f.Interface())
				if strings.HasSuffix(line, " = ") {
					line += "[empty string]"
				}
				line += "\n"
				logger.Debug(line)
			}
		}
		logger.Debug("\n")

		logger.Debug("    Players (" + strconv.Itoa(len(t.Players)) + "):")
		if t.Players == nil {
			logger.Debug("      [Players is nil; this should never happen]")
		} else {
			for j, p := range t.Players { // This is a []*Player
				logger.Debug("        " + strconv.Itoa(j) + " - " +
					"User ID: " + strconv.Itoa(p.ID) + ", " +
					"Username: " + p.Name + ", " +
					"Present: " + strconv.FormatBool(p.Present))
			}
			if len(t.Players) == 0 {
				logger.Debug("        [no players]")
			}
		}
		logger.Debug("\n")

		logger.Debug("    Spectators (" + strconv.Itoa(len(t.Spectators)) + "):")
		if t.Spectators == nil {
			logger.Debug("      [Spectators is nil; this should never happen]")
		} else {
			for j, sp := range t.Spectators { // This is a []*Session
				logger.Debug("        " + strconv.Itoa(j) + " - " +
					"User ID: " + strconv.Itoa(sp.ID) + ", " +
					"Username: " + sp.Name)
			}
			if len(t.Spectators) == 0 {
				logger.Debug("        [no spectators]")
			}
		}
		logger.Debug("\n")

		logger.Debug("    DisconSpectators (" + strconv.Itoa(len(t.DisconSpectators)) + "):")
		if t.DisconSpectators == nil {
			logger.Debug("      [DisconSpectators is nil; this should never happen]")
		} else {
			for k := range t.DisconSpectators { // This is a map[int]struct{}
				logger.Debug("        User ID: " + strconv.Itoa(k))
			}
			if len(t.DisconSpectators) == 0 {
				logger.Debug("        [no disconnected spectators]")
			}
		}
		logger.Debug("\n")

		logger.Debug("    Chat (" + strconv.Itoa(len(t.Chat)) + "):")
		if t.Chat == nil {
			logger.Debug("      [Chat is nil; this should never happen]")
		} else {
			for j, m := range t.Chat { // This is a []*GameChatMessage
				logger.Debug("        " + strconv.Itoa(j) + " - " +
					"[" + strconv.Itoa(m.UserID) + "] <" + m.Username + "> " + m.Msg)
			}
			if len(t.Chat) == 0 {
				logger.Debug("        [no chat]")
			}
		}
		logger.Debug("\n")

		logger.Debug("---------------------------------------------------------------")
	}

	// Print out all of the current users
	logger.Debug("Current users (" + strconv.Itoa(len(sessions)) + "):")
	if len(sessions) == 0 {
		logger.Debug("    [no users]")
	}
	sessionsMutex.RLock()
	for i, s2 := range sessions { // This is a map[int]*Session
		logger.Debug("    User ID: " + strconv.Itoa(i) + ", " +
			"Username: " + s2.Username() + ", " +
			"Status: " + strconv.Itoa(s2.Status()))
	}
	sessionsMutex.RUnlock()
	logger.Debug("---------------------------------------------------------------")
}

func debugFunction() {
	logger.Debug("Executing debug function(s).")

	// updateAllSeedNumGames()
	// updateAllUserStats()
	// updateAllVariantStats()
	// updateUserStatsFromPast24Hours()
	// getBadGameIDs()

	updateUserStatsFromInterval("2 hours")

	logger.Debug("Debug function(s) complete.")
}

/*
func updateAllSeedNumGames() {
	if err := models.Seeds.UpdateAll(); err != nil {
		logger.Error("Failed to update the number of games for every seed:", err)
		return
	}
	logger.Info("Updated the number of games for every seed.")
}

func updateAllUserStats() {
	if err := models.UserStats.UpdateAll(variantGetHighestID()); err != nil {
		logger.Error("Failed to update the stats for every user:", err)
		return
	}
	logger.Info("Updated the stats for every user.")
}

func updateAllVariantStats() {
	highestID := variantGetHighestID()
	maxScores := make([]int, 0)
	for i := 0; i <= highestID; i++ {
		variantName := variantIDMap[i]
		variant := variants[variantName]
		maxScores = append(maxScores, variant.MaxScore)
	}

	if err := models.VariantStats.UpdateAll(highestID, maxScores); err != nil {
		logger.Error("Failed to update the stats for every variant:", err)
	} else {
		logger.Info("Updated the stats for every variant.")
	}
}
*/

func updateUserStatsFromInterval(interval string) {
	// Get the games played in the last X hours/days/whatever
	// Interval must mast a valid Postgres interval
	// https://popsql.com/learn-sql/postgresql/how-to-query-date-and-time-in-postgresql
	var gameIDs []int
	if v, err := models.Games.GetGameIDsSinceInterval(interval); err != nil {
		logger.Error("Failed to get the game IDs for the last \""+interval+"\":", err)
	} else {
		gameIDs = v
	}

	updateStatsFromGameIDs(gameIDs)
}

func updateStatsFromGameIDs(gameIDs []int) {
	// Get the games corresponding to these IDs
	var gameHistoryList []*GameHistory
	if v, err := models.Games.GetHistory(gameIDs); err != nil {
		logger.Error("Failed to get the games from the database:", err)
		return
	} else {
		gameHistoryList = v
	}

	for _, gameHistory := range gameHistoryList {
		updateStatsFromGameHistory(gameHistory)
	}
}

// updateStatsFromGameHistory is mostly copied from the "Game.WriteDatabaseStats()" function
// (the difference is that it works on a "GameHistory" instead of a "Game")
func updateStatsFromGameHistory(gameHistory *GameHistory) {
	logger.Debug("Updating stats for game: " + strconv.Itoa(gameHistory.ID))

	// Local variables
	variant := variants[gameHistory.Options.VariantName]
	// 2-player is at index 0, 3-player is at index 1, etc.
	bestScoreIndex := gameHistory.Options.NumPlayers - 2

	// Update the variant-specific stats for each player
	modifier := gameHistory.Options.GetModifier()
	for _, playerName := range gameHistory.PlayerNames {
		// Check to see if this username exists in the database
		var userID int
		if exists, v, err := models.Users.Get(playerName); err != nil {
			logger.Error("Failed to get user \""+playerName+"\":", err)
			return
		} else if !exists {
			logger.Error("User \"" + playerName + "\" does not exist in the database.")
			return
		} else {
			userID = v.ID
		}

		// Get their current best scores
		var userStats *UserStatsRow
		if v, err := models.UserStats.Get(userID, variant.ID); err != nil {
			logger.Error("Failed to get the stats for user "+playerName+":", err)
			return
		} else {
			userStats = v
		}

		thisScore := &BestScore{
			Score:    gameHistory.Score,
			Modifier: modifier,
		}
		bestScore := userStats.BestScores[bestScoreIndex]
		if thisScore.IsBetterThan(bestScore) {
			bestScore.Score = gameHistory.Score
			bestScore.Modifier = modifier
		}

		// Update their stats
		// (even if they did not get a new best score,
		// we still want to update their average score and strikeout rate)
		if err := models.UserStats.Update(userID, variant.ID, userStats); err != nil {
			logger.Error("Failed to update the stats for user "+playerName+":", err)
			return
		}
	}

	// Get the current stats for this variant
	var variantStats VariantStatsRow
	if v, err := models.VariantStats.Get(variant.ID); err != nil {
		logger.Error("Failed to get the stats for variant "+strconv.Itoa(variant.ID)+":", err)
		return
	} else {
		variantStats = v
	}

	// If the game was played with no modifiers, update the stats for this variant
	if modifier == 0 {
		bestScore := variantStats.BestScores[bestScoreIndex]
		if gameHistory.Score > bestScore.Score {
			bestScore.Score = gameHistory.Score
		}
	}

	// Write the updated stats to the database
	// (even if the game was played with modifiers,
	// we still need to update the number of games played)
	if err := models.VariantStats.Update(variant.ID, variant.MaxScore, variantStats); err != nil {
		logger.Error("Failed to update the stats for variant "+strconv.Itoa(variant.ID)+":", err)
		return
	}
}

/*
func variantGetHighestID() int {
	highestID := 0
	for k := range variantIDMap {
		if k > highestID {
			highestID = k
		}
	}
	return highestID
}

func getBadGameIDs() {
	// Get all game IDs
	var ids []int
	if v, err := models.Games.GetAllIDs(); err != nil {
		logger.Fatal("Failed to get all of the game IDs:", err)
		return
	} else {
		ids = v
	}

	for i, id := range ids {
		if i > 1000 {
			break
		}
		logger.Debug("ON GAME:", id)
		s := newFakeSession(1, "Server")
		commandReplayCreate(s, &CommandData{ // Manual invocation
			Source:     "id",
			GameID:     id,
			Visibility: "solo",
		})
		commandTableUnattend(s, &CommandData{ // Manual invocation
			TableID: tableIDCounter,
		})
	}

	logger.Debug("BAD GAME IDS:")
	logger.Debug(badGameIDs)
}
*/
