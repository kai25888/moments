---
name: moments-project-memory
description: Maintain durable memory and daily logs specifically for the `moments` fork in this workspace. Use this skill whenever work in this repository changes branch strategy, tag strategy, GHCR image publishing, self-hosting deployment rules, branding direction, or high-priority hardening tasks. Also use it whenever the user asks to update `MEMORY.md`, save today's key progress, or preserve important decisions for the `moments-plus` line.
---

# Moments Project Memory

This skill is the project-specific version of the general memory/journal workflow.

It is for the `moments` fork whose long-term custom line is `moments-plus`.

## Default files for this repository

- Durable memory file: `MEMORY.md`
- Daily log directory: `daily-logs/`
- Daily log filename: `daily-logs/YYYY-MM-DD.md`

Always update these files at the repository root unless the user explicitly requests a different location.

## Project rules to preserve

These facts should be treated as current long-term defaults unless the user changes them:

- The formal long-term branch is `moments-plus/main`.
- Namespaced fork tags use `moments-plus/v*`.
- Do not replace or redefine upstream-style `v*` tags.
- `moments-plus/main` branch images publish to `ghcr.io/kai25888/moments-plus:latest`.
- `moments-plus/v*` namespaced tag images publish to `ghcr.io/kai25888/moments-plus:vX.Y.Z`.
- `tech-standards/main` branch images publish to `ghcr.io/kai25888/tech-standards:latest`.
- `tech-standards/*` namespaced tag images publish to `ghcr.io/kai25888/tech-standards:vX.Y.Z`.
- The fork is intended for self-hosted deployment, especially NAS-friendly Docker pulls.
- The customization strategy is light branding plus ongoing hardening, not a full rewrite.

## What belongs in MEMORY.md here

Persist only stable facts such as:

- branch and tag conventions
- image naming and publishing rules
- deployment preferences
- branding direction
- long-term security or architecture priorities
- user preferences about naming or release hygiene

If today's work changes one of these defaults, update `MEMORY.md`.

## What belongs in the daily log here

Use the daily log for:

- what changed today in the fork
- which workflow or release actions were triggered today
- network, git, CI, or publishing problems encountered today
- workarounds used today
- current project status at end of day
- immediate next steps

## Update workflow

1. Read the existing `MEMORY.md` and today's daily log if they exist.
2. Extract new information from the current session.
3. Merge new durable facts into `MEMORY.md` without duplicating existing rules.
4. Merge today's execution details into `daily-logs/YYYY-MM-DD.md`.
5. Replace stale project conventions instead of leaving conflicting notes side by side.

## Merge rules

- Do not duplicate the same branch, tag, or image convention in multiple bullets.
- If an old convention has been superseded, remove or rewrite the old entry.
- If today's log already mentions an event, add the result or consequence instead of repeating the event.
- Keep exact identifiers accurate: branch names, tag names, image tags, workflow intent, and paths.

## Suggested daily log sections

```md
# YYYY-MM-DD

## 今日目标
## 今日完成
## 关键决策
## 遇到的问题
## 当前状态
## 下一步建议
```

## Suggested durable memory sections

```md
# Moments Plus Long-Term Memory

## Project Positioning
## Branch And Release Strategy
## Image Publishing Rules
## Product And UX Direction
## Implemented Direction So Far
## Deployment Notes
## High-Priority Follow-Ups
## Working Preferences
```

## Response expectation

When asked to save memory or today's journal for this repository:

- update the files directly
- keep them concise and non-redundant
- mention which files were updated
- mention any newly established durable convention
