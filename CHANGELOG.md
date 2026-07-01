# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added
- Session token (JWT) authentication support via Cookie header
- Automatic token type detection (API key vs JWT)
- Smart parentDocumentId fallback for create operations
- Comprehensive authentication guide (AUTHENTICATION.md)
- Helper functions for finding parent documents in manifest

### Fixed
- Session tokens can now create documents (via parentDocumentId fallback)
- HTTP headers match browser behavior (x-api-version, x-editor-version)
- Authorization error handling with clear user messages
- **documents.info response parsing** - Fixed nested response structure wrapper
- **Clone command with session token** - Now works end-to-end (was failing with empty collectionID)
- GetCollection graceful fallback when session token lacks collections.info permission

### Changed
- JWT tokens sent via Cookie header instead of Authorization Bearer
- Create document strategy: try collectionId first, fallback to parentDocumentId
- Error messages now distinguish between different auth failure scenarios
- GetDocument now correctly parses `{data: {document: {...}}}` response structure

## [0.1.0] - 2026-07-01

### Added
- Initial release
- Commands: init, list, clone, status, add, push
- Support for 7 URL formats in clone command
- Automatic frontmatter management
- 3-way conflict detection (local, remote, manifest)
- Cross-platform binaries (Linux, macOS, Windows)

[Unreleased]: https://github.com/rayzalzero/outline-cli/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/rayzalzero/outline-cli/releases/tag/v0.1.0
