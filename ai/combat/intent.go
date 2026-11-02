package combat

// CombatIntent is a semantic response, not a keypress. A controller must
// validate the active encounter and movement constraints before acting.
type CombatIntent string

const (
	IntentNoAction   CombatIntent = "no-action"
	IntentDodge       CombatIntent = "dodge"
	IntentStayStill   CombatIntent = "stay-still"
	IntentKeepMoving  CombatIntent = "keep-moving"
	IntentCollect     CombatIntent = "collect-or-interact"
	IntentShield      CombatIntent = "shield"
	IntentShoot       CombatIntent = "shoot"
	IntentWarning     CombatIntent = "warning-only"
)

// AttackIntent converts documented contact semantics into a safe high-level
// intent. Unknown colors never become dodge commands by default.
func AttackIntent(mode SoulMode, class AttackClass) CombatIntent {
	if mode == SoulGreen && class == AttackWhite {
		return IntentShield
	}
	switch ClassRule(class) {
	case RuleDodge:
		return IntentDodge
	case RuleStayStill:
		return IntentStayStill
	case RuleKeepMoving:
		return IntentKeepMoving
	case RuleCollect:
		return IntentCollect
	case RuleShield:
		return IntentShield
	case RuleShoot:
		if mode == SoulYellow {
			return IntentShoot
		}
		return IntentWarning
	case RuleWarning:
		return IntentWarning
	case RuleHarmless:
		return IntentNoAction
	default:
		return IntentNoAction
	}
}

// ModeCanMove reports whether free directional movement is available in the
// current SOUL mode. Blue movement is gravity-constrained; green is shield-led;
// yellow and purple modes require their own encounter-specific controls.
func ModeCanMove(mode SoulMode) bool {
	switch mode {
	case SoulRed:
		return true
	case SoulBlue, SoulGreen, SoulYellow, SoulPurple:
		return false
	default:
		return false
	}
}
