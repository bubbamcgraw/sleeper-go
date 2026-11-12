package sleeper

type Sport string

const (
	NFL Sport = "nfl"
)

type User struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"avatar"`
}

type League struct {
	LeagueID         string         `json:"league_id"`
	Name             string         `json:"name"`
	Season           string         `json:"season"`
	SeasonType       string         `json:"season_type"`
	Sport            string         `json:"sport"`
	Status           string         `json:"status"`
	TotalRosters     int            `json:"total_rosters"`
	RosterPositions  []string       `json:"roster_positions"`
	Settings         map[string]any `json:"settings"`
	ScoringSettings  map[string]any `json:"scoring_settings"`
	DraftID          string         `json:"draft_id"`
	Avatar           string         `json:"avatar"`
	PreviousLeagueID string         `json:"previous_league_id"`
}

// Roster is one franchise's set of players within a League.
type Roster struct {
	RosterID int            `json:"roster_id"`
	OwnerID  string         `json:"owner_id"`
	LeagueID string         `json:"league_id"`
	Players  []string       `json:"players"`
	Starters []string       `json:"starters"`
	Reserve  []string       `json:"reserve"`
	Settings map[string]any `json:"settings"`
}

// LeagueUser is a member of a League — richer than User in league context
// (adds ownership + a per-league team_name via Metadata) but without the
// account-level fields (Username) that only /user/<id> returns.
type LeagueUser struct {
	UserID      string         `json:"user_id"`
	DisplayName string         `json:"display_name"`
	Avatar      string         `json:"avatar"`
	IsOwner     bool           `json:"is_owner"`
	Metadata    map[string]any `json:"metadata"`
}

// Player is a single entry from Sleeper's full player pool (GetAllPlayers).
// Sleeper's dump has dozens of fields; only the identity/roster-relevant
// ones are modeled here — notably EspnID, which bridges a Sleeper player to
// the richer metadata available from espn-go's GetPlayer.
type Player struct {
	PlayerID         string   `json:"player_id"`
	EspnID           int      `json:"espn_id"`
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	FullName         string   `json:"full_name"`
	Position         string   `json:"position"`
	FantasyPositions []string `json:"fantasy_positions"`
	Team             string   `json:"team"`
	Status           string   `json:"status"`
	Active           bool     `json:"active"`
}

// Draft is a single Sleeper draft, from GET /draft/<draft_id>. Picks aren't
// included here — Sleeper serves them from a separate endpoint, see
// GetDraftPicks.
type Draft struct {
	DraftID        string         `json:"draft_id"`
	LeagueID       string         `json:"league_id"`
	Type           string         `json:"type"`
	Status         string         `json:"status"`
	Sport          string         `json:"sport"`
	Season         string         `json:"season"`
	SeasonType     string         `json:"season_type"`
	StartTime      int64          `json:"start_time"`
	Created        int64          `json:"created"`
	Settings       map[string]any `json:"settings"`
	Metadata       map[string]any `json:"metadata"`
	DraftOrder     map[string]any `json:"draft_order"`
	SlotToRosterID map[string]any `json:"slot_to_roster_id"`
}

// DraftPick is a single selection within a Draft, from
// GET /draft/<draft_id>/picks.
type DraftPick struct {
	DraftID   string         `json:"draft_id"`
	PlayerID  string         `json:"player_id"`
	PickedBy  string         `json:"picked_by"`
	RosterID  int            `json:"roster_id"`
	Round     int            `json:"round"`
	DraftSlot int            `json:"draft_slot"`
	PickNo    int            `json:"pick_no"`
	IsKeeper  *bool          `json:"is_keeper"`
	Metadata  map[string]any `json:"metadata"`
}

// Matchup is one roster's performance in a single week of a League, from
// GET /league/<league_id>/matchups/<week>. Rosters sharing the same
// MatchupID played each other that week (a bye week roster gets no
// MatchupID partner).
type Matchup struct {
	RosterID       int                `json:"roster_id"`
	MatchupID      int                `json:"matchup_id"`
	Points         float64            `json:"points"`
	Starters       []string           `json:"starters"`
	StartersPoints []float64          `json:"starters_points"`
	Players        []string           `json:"players"`
	PlayersPoints  map[string]float64 `json:"players_points"`
}

// BracketSource says a bracket slot is filled by the winner or loser of an
// earlier matchup (identified by its BracketMatchup.MatchupID) rather than
// a fixed roster — exactly one of Winner/Loser is set, matching whichever
// key ("w" or "l") Sleeper includes.
type BracketSource struct {
	Winner *int `json:"w"`
	Loser  *int `json:"l"`
}

// BracketMatchup is a single game within a Sleeper playoff bracket, from
// GET /league/<league_id>/winners_bracket or /losers_bracket. Team1/Team2
// are nil until the bracket determines who fills that slot (see
// Team1From/Team2From); Winner/Loser are nil until the game's been played.
type BracketMatchup struct {
	Round     int            `json:"r"`
	MatchupID int            `json:"m"`
	Team1     *int           `json:"t1"`
	Team2     *int           `json:"t2"`
	Team1From *BracketSource `json:"t1_from"`
	Team2From *BracketSource `json:"t2_from"`
	Winner    *int           `json:"w"`
	Loser     *int           `json:"l"`
	// Position, if set, marks this matchup as determining a specific
	// standings placement (e.g. 3rd place) rather than advancing its
	// winner further into the bracket.
	Position *int `json:"p"`
}

// State is Sleeper's current season/week for a sport, from GET /state/<sport>.
type State struct {
	Week               int    `json:"week"`
	SeasonType         string `json:"season_type"`
	SeasonStartDate    string `json:"season_start_date"`
	Season             string `json:"season"`
	PreviousSeason     string `json:"previous_season"`
	Leg                int    `json:"leg"`
	LeagueSeason       string `json:"league_season"`
	LeagueCreateSeason string `json:"league_create_season"`
	DisplayWeek        int    `json:"display_week"`
}
