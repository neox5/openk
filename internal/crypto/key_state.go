package crypto

// KeyState represents the lifecycle state of a key
type KeyState int

const (
	// KeyStateActive indicates key is enabled and available for use
	KeyStateActive KeyState = iota
	// KeyStatePendingRotation indicates key is scheduled for rotation
	KeyStatePendingRotation
	// KeyStateInactive indicates key is disabled but can be re-enabled
	KeyStateInactive
	// KeyStateDestroyed indicates key is permanently deactivated
	KeyStateDestroyed
)

// String returns a string representation of the KeyState
func (s KeyState) String() string {
	switch s {
	case KeyStateActive:
		return "Active"
	case KeyStatePendingRotation:
		return "PendingRotation"
	case KeyStateInactive:
		return "Inactive"
	case KeyStateDestroyed:
		return "Destroyed"
	default:
		return "Unknown"
	}
}
