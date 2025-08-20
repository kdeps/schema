# AGENTS.md - Kdeps Schema Development Guide

## Project Overview
This repository contains PKL schema definitions for the Kdeps AI Agent framework. It defines types for workflows, resources, LLMs, memory, tools, and other core components used by Kdeps agents.

## Build & Test Commands
- **Generate Go code**: `make generate` or `./gradlew makePackages`
- **Build documentation**: `./gradlew pkldoc`
- **Run all Go tests**: `go test ./...`
- **Test specific package**: `go test ./gen/resource`
- **Validate PKL modules**: `pkl project resolve --root-dir . --module-path deps/pkl deps/pkl`

## Architecture
- `deps/pkl/`: Core PKL schema definitions (Resource.pkl, Workflow.pkl, LLM.pkl, etc.)
- `gen/`: Auto-generated Go code from PKL schemas (DO NOT EDIT)
- `assets/`: Static assets and helper Go code
- Multi-language support: PKL generates Go bindings using pkl-go

## Code Style Guidelines
- **Indent**: 2 spaces (per .editorconfig)
- **Line length**: 100 characters max
- **PKL modules**: Use proper module documentation with `@ModuleInfo { minPklVersion = "0.28.2" }`
- **Go packages**: Use `@go.Package { name = "..." }` annotations for code generation
- **Naming**: Use ActionID for resource identifiers, follow camelCase for PKL properties
- **Validation**: Use regex patterns for ActionID validation (alphanumeric or @package/action:version format)
- **Imports**: Import standard PKL modules (pkl:json, pkl:yaml, etc.) and local modules consistently
- **Generated code**: Never edit files in gen/ directory - they're auto-generated from PKL sources
