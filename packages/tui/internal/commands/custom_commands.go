package commands

// Custom commands for user-specific modifications
// This file can be safely customized without conflicts during merges

const (
	// Custom command names for mode operations
	ModeSelectCommand CommandName = "mode_select" // Shows popup to select mode directly
)

// GetCustomCommands returns user-customized command definitions
// This function overrides default keybindings and adds new commands
func GetCustomCommands() []Command {
	return []Command{
		// Override the default tab behavior for mode switching
		{
			Name:        SwitchModeCommand,
			Description: "cycle through modes",
			Keybindings: parseBindings("shift+tab"), // Changed from "tab" to "shift+tab"
		},
		// New command: Tab opens mode selection popup
		{
			Name:        ModeSelectCommand,
			Description: "select mode from popup",
			Keybindings: parseBindings("tab"), // Tab now opens mode selection popup
		},
	}
}

// ApplyCustomizations modifies the default commands with custom keybindings
func ApplyCustomizations(commands []Command) []Command {
	customCommands := GetCustomCommands()
	
	// Create a map for efficient lookup
	customMap := make(map[CommandName]Command)
	for _, cmd := range customCommands {
		customMap[cmd.Name] = cmd
	}
	
	// Apply customizations to existing commands and add new ones
	result := make([]Command, 0, len(commands)+len(customCommands))
	
	// Process existing commands, applying customizations where applicable
	for _, cmd := range commands {
		if customCmd, exists := customMap[cmd.Name]; exists {
			// Override with custom definition
			result = append(result, customCmd)
			delete(customMap, cmd.Name) // Remove from map to track what's been processed
		} else {
			// Keep original command
			result = append(result, cmd)
		}
	}
	
	// Add any remaining new custom commands
	for _, cmd := range customMap {
		result = append(result, cmd)
	}
	
	return result
}