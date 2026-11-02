package combat

// RouteMode selects the long-term moral route the controller should preserve.
type RouteMode string

const (
	RouteTruePacifist RouteMode = "true-pacifist"
)

// EncounterKind is a normalized encounter identifier. Unknown encounters are
// handled conservatively; names are not inferred from a partially recognized frame.
type EncounterKind string

const (
	EncounterOrdinary    EncounterKind = "ordinary"
	EncounterToriel      EncounterKind = "toriel"
	EncounterPapyrus     EncounterKind = "papyrus"
	EncounterUndyne      EncounterKind = "undyne"
	EncounterMettatonEX  EncounterKind = "mettaton-ex"
	EncounterAsgore      EncounterKind = "asgore"
	EncounterOmegaFlowey EncounterKind = "omega-flowey"
	EncounterAsriel      EncounterKind = "asriel"
)

// BattleStrategy describes the preferred non-lethal way to finish an encounter.
// It is guidance for a future controller, not an instruction to press buttons.
type BattleStrategy string

const (
	StrategyActThenSpare BattleStrategy = "act-until-spare-then-spare"
	StrategySpare        BattleStrategy = "spare-when-available"
	StrategyFlee         BattleStrategy = "flee"
	StrategyRatings      BattleStrategy = "raise-ratings"
	StrategyBossProgress BattleStrategy = "boss-progress-no-exp"
	StrategySaveSouls    BattleStrategy = "save-lost-souls"
	StrategyNoAction     BattleStrategy = "no-action-unrecognized"
)

// RouteDecision contains conservative route guidance. RequiresHumanReview is
// true whenever the current CV state is insufficient to choose a safe action.
type RouteDecision struct {
	Strategy           BattleStrategy
	AllowFight         bool
	RequiresHumanReview bool
	Reason             string
}

// TruePacifistDecision returns route-safe guidance for a recognized encounter.
// Ordinary monsters are never attacked by default. Asgore and Flowey are
// progression fights on the required Neutral completion; Asriel is the final
// True Pacifist boss. These exceptions do not award EXP in the original game.
func TruePacifistDecision(encounter EncounterKind) RouteDecision {
	switch encounter {
	case EncounterOrdinary:
		return RouteDecision{
			Strategy: StrategyActThenSpare,
			Reason: "Use encounter-specific ACT actions until MERCY is available; never choose a lethal outcome.",
		}
	case EncounterToriel, EncounterPapyrus, EncounterUndyne:
		return RouteDecision{
			Strategy: StrategyActThenSpare,
			Reason: "Keep the character alive and resolve the encounter through its non-lethal dialogue or mercy sequence.",
		}
	case EncounterMettatonEX:
		return RouteDecision{
			Strategy: StrategyRatings,
			Reason: "Reach the peaceful ratings threshold; use the yellow-heart phone mechanic only when its phase is recognized.",
		}
	case EncounterAsgore, EncounterOmegaFlowey:
		return RouteDecision{
			Strategy: StrategyBossProgress,
			AllowFight: true,
			Reason: "Required to complete the prerequisite Neutral ending; these boss encounters do not grant EXP.",
		}
	case EncounterAsriel:
		return RouteDecision{
			Strategy: StrategySaveSouls,
			Reason: "Use SAVE on the Lost Souls and follow the final encounter's non-lethal progression.",
		}
	default:
		return RouteDecision{
			Strategy: StrategyNoAction,
			RequiresHumanReview: true,
			Reason: "Encounter identity is unknown; suspend combat input rather than risk a kill.",
		}
	}
}

// CanUseFightForRoute prevents an ordinary FIGHT command from being selected
// by the True Pacifist policy. Boss progression exceptions must still be
// gated by a verified encounter state in the caller.
func CanUseFightForRoute(mode RouteMode, encounter EncounterKind) bool {
	if mode != RouteTruePacifist {
		return false
	}
	decision := TruePacifistDecision(encounter)
	return decision.AllowFight
}
