package constants

// Key constants for keyboard input
const (
	KeyQ       = "q"
	KeyCtrlC   = "ctrl+c"
	KeyEnter   = "enter"
	KeyEsc     = "esc"
	KeyUp      = "up"
	KeyDown    = "down"
	KeyAltUp   = "k"
	KeyAltDown = "j"
	KeyAltBack = "-"
	KeyTab     = "tab"
)

// Authentication method constants
const (
	// AWS authentication methods
	AWSProfileAuth = "profile"

	// Azure authentication methods (future)
	AzureCliAuth       = "cli"
	AzureConfigDirAuth = "config-dir"

	// GCP authentication methods (future)
	GCPServiceAccountAuth     = "service-account"
	GCPApplicationDefaultAuth = "adc"
)

// Configuration key constants
const (
	// AWS configuration keys
	AWSProfileKey = "profile"
	AWSRegionKey  = "region"

	// Azure configuration keys (future)
	AzureSubscriptionKey = "subscription"
	AzureLocationKey     = "location"
	AzureTenantKey       = "tenant"
	AzureConfigDirKey    = "config-dir"

	// GCP configuration keys (future)
	GCPProjectKey        = "project"
	GCPZoneKey           = "zone"
	GCPRegionKey         = "region"
	GCPServiceAccountKey = "service-account-path"
)
