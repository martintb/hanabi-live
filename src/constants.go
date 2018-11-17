package main

import (
	"time"
)

const (
	actionTypeClue = iota
	actionTypePlay
	actionTypeDiscard
	actionTypeDeckPlay
	actionTypeTimeLimitReached
	actionTypeIdleLimitReached
)

const (
	clueTypeNumber = iota
	clueTypeColor
)

const (
	replayActionTypeTurn = iota
	replayActionTypeArrow
	replayActionTypeLeaderTransfer
	replayActionTypeMorph
	replayActionTypeSound
)

const (
	// The amount of time that a game is inactive before it is killed by the server
	idleGameTimeout = time.Minute * 30

	// The amount of time that someone can be on the waiting list
	idleWaitingListTimeout = time.Hour * 8

	// The amount of time in between allowed @here Discord alerts
	discordAtHereTimeout = time.Hour * 2

	// Discord emotes
	pogChamp   = "<:PogChamp:254683883033853954>"
	bibleThump = "<:BibleThump:254683882601840641>"
)
