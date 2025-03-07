package clerk

type SessionReverificationLevel string

const (
	// SessionReverificationLevelFirstFactor uses first factor methods
	// e.g. email/phone code, password
	SessionReverificationLevelFirstFactor SessionReverificationLevel = "first_factor"

	// SessionReverificationLevelSecondFactor uses second factor methods, if available
	// e.g. authenticator app, backup code
	SessionReverificationLevelSecondFactor SessionReverificationLevel = "second_factor"

	// SessionReverificationLevelMultiFactor requires both first and second factor
	// methods, if available
	SessionReverificationLevelMultiFactor SessionReverificationLevel = "multi_factor"
)

type SessionReverificationPolicy struct {
	// AfterMinutes is the session age threshold before reverification
	AfterMinutes int64

	// Level specifies which verification factors are required
	Level SessionReverificationLevel
}

var (
	SessionReverificationStrictMFA = SessionReverificationPolicy{
		AfterMinutes: 10,
		Level:        SessionReverificationLevelMultiFactor,
	}

	SessionReverificationStrict = SessionReverificationPolicy{
		AfterMinutes: 10,
		Level:        SessionReverificationLevelSecondFactor,
	}

	SessionReverificationModerate = SessionReverificationPolicy{
		AfterMinutes: 60,
		Level:        SessionReverificationLevelSecondFactor,
	}

	SessionReverificationLax = SessionReverificationPolicy{
		AfterMinutes: 1_440,
		Level:        SessionReverificationLevelSecondFactor,
	}
)
