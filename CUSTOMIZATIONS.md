# OpenCode Customizations

This document describes the customizations made to OpenCode to avoid conflicts during future merges from the official repository.

## Summary of Changes

### 1. Full Width Chat Window
**File Modified:** `packages/tui/internal/tui/tui.go` (lines 461, 466)
- **Original**: Chat window limited to 104 characters width
- **Modified**: Chat window uses full terminal width
- **Reason**: Better utilization of wide screens

### 2. Custom Keybinding System
**Files Added:**
- `packages/tui/internal/commands/custom_commands.go` - Separate file for custom keybindings
- `packages/tui/internal/components/dialog/mode_select.go` - Mode selection popup dialog

**File Modified:**
- `packages/tui/internal/commands/command.go` (lines 337-338) - Integration point for custom commands
- `packages/tui/internal/tui/tui.go` (lines 813-815, 535-541) - Handlers for custom commands

**Changes:**
- **Tab key**: Changed from mode cycling to mode selection popup
- **Shift+Tab**: Now used for mode cycling (was Tab)

## File Organization for Merge Safety

### Safe to Customize (Low Merge Conflict Risk)
- `packages/tui/internal/commands/custom_commands.go` - **NEW FILE** - Contains all custom keybindings
- `packages/tui/internal/components/dialog/mode_select.go` - **NEW FILE** - Mode selection dialog
- `CUSTOMIZATIONS.md` - **NEW FILE** - This documentation

### Modified Core Files (Higher Merge Conflict Risk)
1. **`packages/tui/internal/tui/tui.go`**
   - Lines 461, 466: Chat width changes
   - Lines 813-815: Mode select command handler
   - Lines 535-541: Mode selected message handler

2. **`packages/tui/internal/commands/command.go`**
   - Lines 337-338: Integration with custom commands system

## How to Maintain During Updates

### Before Merging from Official Repository:
1. **Backup your customizations:**
   ```bash
   cp packages/tui/internal/commands/custom_commands.go /backup/
   cp packages/tui/internal/components/dialog/mode_select.go /backup/
   cp CUSTOMIZATIONS.md /backup/
   ```

2. **Note the line numbers** of modified sections in core files.

### After Merging from Official Repository:
1. **Restore new files:**
   ```bash
   cp /backup/custom_commands.go packages/tui/internal/commands/
   cp /backup/mode_select.go packages/tui/internal/components/dialog/
   cp /backup/CUSTOMIZATIONS.md ./
   ```

2. **Re-apply modifications to core files:**
   - Follow the modification patterns described in this document
   - Test the build after each change: `cd packages/tui && go build ./cmd/opencode`

### Conflict Resolution:
If conflicts occur during merge:
1. **Chat width changes**: Look for `container := min(a.width, 104)` and apply the full-width modifications
2. **Command integration**: Look for the `LoadFromConfig` function and add the custom commands integration
3. **TUI handlers**: Add the mode selection handlers in the appropriate switch statements

## Architecture Overview

The customization system is designed to be:
- **Modular**: Custom commands are in separate files
- **Non-intrusive**: Minimal changes to core files
- **Extensible**: Easy to add more custom commands using the same pattern

### Custom Commands System Flow:
1. `custom_commands.go` defines new commands and overrides
2. `command.go` integrates custom commands via `ApplyCustomizations()`
3. `tui.go` handles the new commands and their messages
4. `mode_select.go` provides the popup UI for mode selection

## Usage

### Current Keybindings:
- **Tab**: Open mode selection popup (shows all modes in a list)
- **Shift+Tab**: Cycle through modes sequentially
- **↑/↓ or k/j** (in mode popup): Navigate through modes
- **Enter** (in mode popup): Select highlighted mode
- **Esc** (in mode popup): Cancel mode selection

### Mode Selection Popup Features:
- **Visual highlighting**: Selected mode has colored background
- **Current mode indicator**: Active mode marked with ● symbol
- **Smart positioning**: Opens with current mode pre-selected
- **Keyboard navigation**: Standard vi-style (j/k) and arrow key navigation

### Adding More Custom Commands:
To add more custom keybindings, edit `packages/tui/internal/commands/custom_commands.go`:

```go
// Add to GetCustomCommands() function:
{
    Name:        "your_custom_command",
    Description: "your command description", 
    Keybindings: parseBindings("your_key"),
},
```

Then add the handler in `packages/tui/internal/tui/tui.go` in the `executeCommand` function.

## Testing

After making changes, always test:
```bash
cd packages/tui
go build ./cmd/opencode
```

The customizations should:
1. Allow tab to open mode selection popup
2. Allow shift+tab to cycle modes
3. Display full-width chat window
4. Build without errors

## Version Compatibility

These customizations were created for OpenCode commit: `[current commit hash]`

When updating, check if the customization patterns still apply to the new version's codebase structure.