#!/usr/bin/env bun

import { $ } from "bun"

// Get version from git tag, fallback to dev-local if no tag found
const version = await $`git describe --tags --abbrev=0`
  .text()
  .then((x) => x.substring(1).trim()) // Remove 'v' prefix
  .catch(() => {
    console.warn("No git tag found, using dev-local")
    return "dev-local"
  })

console.log(`🎯 You are now in version ${version}`)
console.log(`Building local version: ${version}`)

// Get current platform info
const platform = process.platform === "darwin" ? "darwin" : process.platform === "win32" ? "windows" : "linux"

const arch = process.arch === "arm64" ? "arm64" : process.arch === "x64" ? "x64" : "x64"

console.log(`Building for ${platform}-${arch}`)

// Set up build directory structure
const buildDir = "./packages/opencode/dist/local"
const binDir = `${buildDir}/bin`
await $`rm -rf ${buildDir}`
await $`mkdir -p ${binDir}`

// Build the Go TUI binary
console.log("Building Go TUI binary...")
const GOARCH = arch === "arm64" ? "arm64" : "amd64"
const tui_name = platform === "windows" ? "tui.exe" : "tui"

await $`cd packages/tui && CGO_ENABLED=0 GOOS=${platform} GOARCH=${GOARCH} go build -ldflags="-s -w -X main.Version=${version}" -o ../${buildDir.replace("./packages/", "")}/${tui_name} ./cmd/opencode/main.go`

// Build the main CLI binary with embedded TUI
console.log("Building main CLI binary...")
const cli_name = platform === "windows" ? "opencode.exe" : "opencode"
const bun_target = platform === "windows" ? `bun-${platform}-${arch}` : `bun-${platform}-${arch}`

await $`cd packages/opencode && bun build --define OPENCODE_VERSION="'${version}'" --compile --minify --target=${bun_target} --outfile=./${buildDir.replace("./packages/opencode/", "")}/bin/${cli_name} ./src/index.ts ./${buildDir.replace("./packages/opencode/", "")}/${tui_name}`

// Clean up the separate TUI binary since it's now embedded
await $`rm -f ${buildDir}/${tui_name}`

console.log(`\n✅ Build complete!`)
console.log(`📁 Binary location: ${binDir}/${cli_name}`)
console.log(`🚀 Version: ${version}`)
