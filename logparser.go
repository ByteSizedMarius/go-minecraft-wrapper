package wrapper

import (
	"fmt"
	"github.com/ByteSizedMarius/go-minecraft-wrapper/events"
	"regexp"
	"strconv"
	"strings"
)

// LogParser is an interface func to decode any server log line
// to its respective event type. The returned events must be either:
// - Cmd: event holds data to be returned to a user command.
// - State: event affects the state of the wrapper.
// - Game: event related to in-game events, like a player died...
// - Nil: event that hold no value and usually ignored/
type LogParser func(string, int) (events.Event, events.EventType)

type logLine struct {
	timestamp  string
	threadName string
	level      string
	output     string
}

var logRegex = regexp.MustCompile(`(\[[0-9:]*\]) \[([A-z(-| )#0-9]*)\/([A-z #]*)\]: (.*)`)

func parseToLogLine(line string) *logLine {
	matches := logRegex.FindAllStringSubmatch(line, 4)
	return &logLine{
		timestamp:  matches[0][1],
		threadName: matches[0][2],
		level:      matches[0][3],
		output:     matches[0][4],
	}
}

var stateEventToRegexp = map[string]*regexp.Regexp{
	events.Started:  regexp.MustCompile(`^Done (?s)(.*)! For help`),
	events.Starting: regexp.MustCompile(`^Starting Minecraft server on (.*)`),
	events.Stopping: regexp.MustCompile(`^Stopping (.*) server`),
	events.Saving:   regexp.MustCompile(`^Saving the game`),
	events.Saved:    regexp.MustCompile(`^Saved (?s)(.*)`),
}

var gameEventToRegex = map[string]*regexp.Regexp{
	events.Banned:           regexp.MustCompile(`^Banned (?s)(.*): (?s)(.*)`),
	events.BanList:          regexp.MustCompile(`^There are (no|\d+) bans(:|\z)`),
	events.BanListEntry:     regexp.MustCompile(`(?s)(.*) was banned by Server: (.*)`),
	events.DataGet:          regexp.MustCompile(`(?s)(.*) has the following (entity|block|storage) data: (.*)`),
	events.DataGetNoEntity:  regexp.MustCompile(`^No (entity|block|storage) was found`),
	events.DefaultGameMode:  regexp.MustCompile(`^The default game mode is now (Survival|Creative|Adventure|Spectator) Mode`),
	events.Difficulty:       regexp.MustCompile(`^The difficulty (?s)(.*)`),
	events.ExperienceAdd:    regexp.MustCompile(`^Gave ([0-9]+) experience (levels|points) to (?s)(.*)`),
	events.ExperienceQuery:  regexp.MustCompile(`(?s)(.*) has ([0-9]+) experience (levels|points)`),
	events.Give:             regexp.MustCompile(`^Gave ([0-9]+) \[(?s)(.*) (?s)(.*)\] to (?s)(.*)`),
	events.NoPlayerFound:    regexp.MustCompile(`^No player was found`),
	events.PlayerPos:        regexp.MustCompile(`(?s)(.*) has the following entity data: \[(?s)([0-9]*)\.[0-9]*d, (?s)([0-9]*)\.[0-9]*d, (?s)([0-9]*)\.[0-9]*d\]`),
	events.PlayerDied:       regexp.MustCompile(`(?s)(.*) (died|was squished too much|was squashed by|was poked|was killed|didn't want to live in the same|withered|froze to death|was fireballed by|was stung to death|starved to death|discovered the floor was lava|was shot|was pummeled|was pricked|experienced kinetic energy|drowned|blew up|was struck by lightning|was blown up|tried to swim in lava|went off with a bang|hit the ground|fell|was burnt to a crisp|was impaled|burned to death|was slain|walked into|was skewered|went up in flames|suffocated)(.*)`),
	events.PlayerJoined:     regexp.MustCompile(`(?s)(.*) joined the game`),
	events.PlayerLeft:       regexp.MustCompile(`(?s)(.*) left the game`),
	events.PlayerUUID:       regexp.MustCompile(`^UUID of player (?s)(.*) is (?s)(.*)`),
	events.Say:              regexp.MustCompile(`[<\[]((?s)([A-z0-9_]*))[>\]] ((?s)(.*))`),
	events.Kicked:           regexp.MustCompile(`^Kicked (?s)(.*): (.*)`),
	events.Seed:             regexp.MustCompile(`^Seed: (.*)`),
	events.ServerOverloaded: regexp.MustCompile(`^Can't keep up! Is the server overloaded\? Running ([0-9]+)ms or ([0-9]+) ticks behind`),
	events.TimeIs:           regexp.MustCompile(`^The time is (?s)(.*)`),
	events.UnknownItem:      regexp.MustCompile(`^Unknown item (.*)`),
	events.Version:          regexp.MustCompile(`^Starting minecraft server version (.*)`),
	events.WhisperTo:        regexp.MustCompile(`^You whisper to (?s)(.*): (.*)`),
	events.WhitelistList:    regexp.MustCompile(`There are [0-9]* whitelisted players: (?s)(.*)`),
	events.WhitelistAdd:     regexp.MustCompile(`Added (?s)(.*) to the whitelist|That player does not exist|Player is already whitelisted`),
	events.WhitelistRemove:  regexp.MustCompile(`Removed (?s)(.*) from the whitelist|That player does not exist|Player is not whitelisted`),
}

var activeGameEvents = map[string]*regexp.Regexp{
	events.PlayerDied:       gameEventToRegex[events.PlayerDied],
	events.PlayerJoined:     gameEventToRegex[events.PlayerJoined],
	events.PlayerLeft:       gameEventToRegex[events.PlayerLeft],
	events.PlayerUUID:       gameEventToRegex[events.PlayerUUID],
	events.Say:              gameEventToRegex[events.Say],
	events.ServerOverloaded: gameEventToRegex[events.ServerOverloaded],
	events.TimeIs:           gameEventToRegex[events.TimeIs],
	events.Version:          gameEventToRegex[events.Version],
	events.PlayerPos:        gameEventToRegex[events.PlayerPos],
	events.WhitelistAdd:     gameEventToRegex[events.WhitelistAdd],
	events.WhitelistRemove:  gameEventToRegex[events.WhitelistRemove],
	events.WhitelistList:    gameEventToRegex[events.WhitelistList],
}

func registerGameEvent(ev string) {
	_, ok := activeGameEvents[ev]
	if !ok {
		activeGameEvents[ev] = gameEventToRegex[ev]
	}
}

func logParserFunc(line string, tick int) (events.Event, events.EventType) {
	ll := parseToLogLine(line)
	if ll.output == "" {
		return events.NilEvent, events.TypeNil
	}

	for e, reg := range stateEventToRegexp {
		if reg.MatchString(ll.output) {
			return events.NewStateEvent(e), events.TypeState
		}
	}
	for e, reg := range activeGameEvents {
		matches := reg.FindStringSubmatch(ll.output)
		if matches == nil {
			continue
		}

		switch e {
		case events.BanList:
			return handleBanList(matches)
		case events.BanListEntry:
			return handleBanListEntry(matches)
		case events.Difficulty:
			return handleDifficulty(matches)
		case events.ExperienceQuery:
			return handleExperienceQuery(matches)
		case events.PlayerJoined:
			return handlePlayerJoined(matches, tick)
		case events.PlayerLeft:
			return handlePlayerLeft(matches, tick)
		case events.PlayerDied:
			return handlePlayerDied(matches, tick)
		case events.PlayerUUID:
			return handlePlayerUUIDEvent(matches, tick)
		case events.PlayerPos:
			return handlePlayerPosEvent(matches, tick)
		case events.Say:
			return handlePlayerSayEvent(matches, tick)
		case events.Version:
			return handleVersionEvent(matches)
		case events.TimeIs:
			return handleTimeEvent(matches)
		case events.DataGet:
			return handleDataGet(matches)
		case events.DataGetNoEntity:
			return handleDataGetNoEntity(matches)
		case events.Seed:
			return handleSeed(matches)
		case events.ServerOverloaded:
			return handleServerOverloaded(matches, tick)
		case events.DefaultGameMode:
			return handleDefaultGameMode(matches)
		case events.Banned:
			return handleBanned(matches)
		case events.WhitelistAdd:
			return handleWhitelistAdd(matches)
		case events.WhitelistRemove:
			return handleWhitelistRemove(matches)
		case events.WhitelistList:
			return handleWhitelistList(matches)
		case events.WhisperTo, events.ExperienceAdd, events.Give, events.NoPlayerFound,
			events.Kicked, events.UnknownItem:
			return events.NewGameEvent(e), events.TypeCmd
		default:
			gameEvent := events.NewGameEvent(e)
			gameEvent.Tick = tick
			return gameEvent, events.TypeGame
		}
	}
	return events.NilEvent, events.TypeNil
}

func handleWhitelistRemove(matches []string) (events.GameEvent, events.EventType) {
	waEvent := events.NewGameEvent(events.WhitelistRemove)

	if strings.Contains(matches[0], "not whitelisted") || strings.Contains(matches[0], "exist") {
		waEvent.Data = map[string]string{
			"error": matches[0],
		}
	} else if strings.Contains(matches[0], "Removed") {
		waEvent.Data = map[string]string{
			"removed_player_name": matches[1],
		}
	} else {
		// todo remove if its fine :)
		panic("????" + fmt.Sprint(matches))
	}

	return waEvent, events.TypeCmd
}

func handleWhitelistAdd(matches []string) (events.GameEvent, events.EventType) {
	waEvent := events.NewGameEvent(events.WhitelistAdd)

	if strings.Contains(matches[0], "already") || strings.Contains(matches[0], "exist") {
		waEvent.Data = map[string]string{
			"error": matches[0],
		}
	} else if strings.Contains(matches[0], "Added") {
		waEvent.Data = map[string]string{
			"added_player_name": matches[1],
		}
	} else {
		// todo remove if its fine :)
		panic("????" + fmt.Sprint(matches))
	}

	return waEvent, events.TypeCmd
}

func handleWhitelistList(matches []string) (events.GameEvent, events.EventType) {
	wlEvent := events.NewGameEvent(events.WhitelistList)
	wlEvent.Data = map[string]string{
		"players": matches[1],
	}
	return wlEvent, events.TypeCmd
}

func handleBanList(matches []string) (events.GameEvent, events.EventType) {
	blEvent := events.NewGameEvent(events.BanList)
	blEvent.Data = map[string]string{
		"entry_type": "header",
	}
	if matches[1] != "no" {
		// This indicates that there are entries to report back...
		blEvent.Data["entry_count"] = matches[1]
	}
	return blEvent, events.TypeCmd
}

func handleBanListEntry(matches []string) (events.GameEvent, events.EventType) {
	bleEvent := events.NewGameEvent(events.BanList)
	bleEvent.Data = map[string]string{
		"entry_type": "item",
		"entry_name": matches[1],
		"reason":     matches[2],
	}
	return bleEvent, events.TypeCmd
}

func handleDifficulty(matches []string) (events.GameEvent, events.EventType) {
	dfEvent := events.NewGameEvent(events.Difficulty)
	dfEvent.Data = map[string]string{}
	if strings.Contains(matches[1], "did not change") {
		dfEvent.Data["error_message"] = matches[0]
	}
	return dfEvent, events.TypeCmd
}

func handleExperienceQuery(matches []string) (events.GameEvent, events.EventType) {
	xqEvent := events.NewGameEvent(events.ExperienceQuery)
	xqEvent.Data = map[string]string{
		"amount": matches[2],
	}
	return xqEvent, events.TypeCmd
}

func handlePlayerJoined(matches []string, tick int) (events.GameEvent, events.EventType) {
	pjEvent := events.NewGameEvent(events.PlayerJoined)
	pjEvent.Tick = tick
	pjEvent.Data = map[string]string{
		"player_name": matches[1],
	}
	return pjEvent, events.TypeGame
}

func handlePlayerLeft(matches []string, tick int) (events.GameEvent, events.EventType) {
	plEvent := events.NewGameEvent(events.PlayerLeft)
	plEvent.Tick = tick
	plEvent.Data = map[string]string{
		"player_name": matches[1],
	}
	return plEvent, events.TypeGame
}

func handlePlayerDied(matches []string, tick int) (events.GameEvent, events.EventType) {
	pdEvent := events.NewGameEvent(events.PlayerDied)
	pdEvent.Tick = tick
	pdEvent.Data = map[string]string{
		"player_name":   matches[1],
		"death_by":      strings.ReplaceAll(matches[2], "\r", ""),
		"death_details": "",
	}
	if len(matches) >= 4 {
		pdEvent.Data["death_details"] = strings.ReplaceAll(matches[3], "\r", "")
	}
	return pdEvent, events.TypeGame
}

func handlePlayerUUIDEvent(matches []string, tick int) (events.GameEvent, events.EventType) {
	puEvent := events.NewGameEvent(events.PlayerUUID)
	puEvent.Tick = tick
	puEvent.Data = map[string]string{
		"player_name": matches[1],
		"player_uuid": matches[2],
	}
	return puEvent, events.TypeGame
}

func handlePlayerPosEvent(matches []string, tick int) (events.GameEvent, events.EventType) {
	puEvent := events.NewGameEvent(events.PlayerPos)
	puEvent.Tick = tick
	puEvent.Data = map[string]string{
		"player_name":  matches[1],
		"player_pos_x": matches[2],
		"player_pos_y": matches[3],
		"player_pos_z": matches[4],
	}
	return puEvent, events.TypeGame
}

func handlePlayerSayEvent(matches []string, tick int) (events.GameEvent, events.EventType) {
	psEvent := events.NewGameEvent(events.Say)
	psEvent.Tick = tick
	psEvent.Data = map[string]string{
		"name":    matches[1],
		"message": strings.ReplaceAll(matches[3], "\r", ""),
	}
	return psEvent, events.TypeGame
}

func handleVersionEvent(matches []string) (events.GameEvent, events.EventType) {
	versionEvent := events.NewGameEvent(events.Version)
	versionEvent.Data = map[string]string{
		"version": matches[1],
	}
	return versionEvent, events.TypeCmd
}

func handleTimeEvent(matches []string) (events.GameEvent, events.EventType) {
	tickStr := matches[1]
	tick, _ := strconv.Atoi(tickStr)
	timeEvent := events.NewGameEvent(events.TimeIs)
	timeEvent.Tick = tick
	return timeEvent, events.TypeCmd
}

func handleDataGet(matches []string) (events.GameEvent, events.EventType) {
	dgEvent := events.NewGameEvent(events.DataGet)
	dgEvent.Data = map[string]string{
		"player_name": matches[1],
		"data_type":   matches[2],
		"data_raw":    matches[3],
	}
	return dgEvent, events.TypeCmd
}

func handleDataGetNoEntity(matches []string) (events.GameEvent, events.EventType) {
	dgEvent := events.NewGameEvent(events.DataGet)
	dgEvent.Data = map[string]string{
		"error_message": matches[0],
	}
	return dgEvent, events.TypeCmd
}

func handleSeed(matches []string) (events.GameEvent, events.EventType) {
	sdEvent := events.NewGameEvent(events.Seed)
	sdEvent.Data = map[string]string{
		"data_raw": matches[1],
	}
	return sdEvent, events.TypeCmd
}

func handleServerOverloaded(matches []string, tick int) (events.GameEvent, events.EventType) {
	soEvent := events.NewGameEvent(events.ServerOverloaded)
	soEvent.Tick = tick
	soEvent.Data = map[string]string{
		"lag_time": matches[1],
		"lag_tick": matches[2],
	}
	return soEvent, events.TypeGame
}

func handleDefaultGameMode(matches []string) (events.GameEvent, events.EventType) {
	gmEvent := events.NewGameEvent(events.DefaultGameMode)
	gmEvent.Data = map[string]string{
		"default_game_mode": matches[1],
	}
	return gmEvent, events.TypeGame
}

func handleBanned(matches []string) (events.GameEvent, events.EventType) {
	bnEvent := events.NewGameEvent(events.Banned)
	bnEvent.Data = map[string]string{
		"player_name": matches[1],
		"reason":      matches[2],
	}
	return bnEvent, events.TypeGame
}
