# Moments Plus Long-Term Memory

## Project Positioning
- This repository is a long-term fork of `moments` for light customization, security hardening, and self-hosted deployment.
- The fork should keep upstream defaults intact where possible and add an independent release lane for custom work.

## Branch And Release Strategy
- The formal long-term development branch is `moments-plus/main`.
- Avoid temporary-feeling branch names like `feature/*` for the primary custom line.
- Do not alter the upstream-style `v*` release tags.
- Use namespaced tags for fork releases: `moments-plus/v*`.

## Image Publishing Rules
- Prefer `GHCR` for fork images so they can be pulled directly on a NAS or other self-hosted Docker hosts.
- Map each long-term branch line to its own GHCR package instead of reusing the upstream `moments` package name.
- `moments-plus/main` branch images publish to:
  - `ghcr.io/kai25888/moments-plus:latest`
- `moments-plus/v*` namespaced tags publish to:
  - `ghcr.io/kai25888/moments-plus:vX.Y.Z`
- `moments-plus/main` and `moments-plus/v*` may share one workflow as long as GitHub Actions still shows distinct refs for branch and tag runs.
- `tech-standards/main` branch images publish to:
  - `ghcr.io/kai25888/tech-standards:latest`
- `tech-standards/*` namespaced tags publish to:
  - `ghcr.io/kai25888/tech-standards:vX.Y.Z`
- Keep fork image workflows separate from upstream `DockerHub` and upstream `v*` release automation.

## Product And UX Direction
- The customization direction is "light branding + stable long-term operation", not a full rewrite.
- Branding and homepage presentation are configurable from the admin settings page.
- The preferred naming and release style should feel formal and durable, not temporary.

## Implemented Direction So Far
- Added branding-related configuration fields and connected them to the homepage, layout, footer, and admin settings page.
- Hardened multiple security and stability risks in frontend and backend flows.
- Consolidated `moments-plus/main` and `moments-plus/v*` image publishing into one workflow while keeping refs distinguishable in Actions.
- 移动端 App 方案已迁移至独立仓库 `kai25888/moments-app`（Capacitor 壳化），原仓库已回退相关改动。

## Deployment Notes
- The fork is optimized for direct image pulls on a NAS.
- Keep the fork's deployment path simple: one product line can use one workflow for both branch builds and namespaced tag builds when the Actions history stays easy to read.
- Canonical NAS pull commands:
  - `docker pull ghcr.io/kai25888/moments-plus:latest`
  - `docker pull ghcr.io/kai25888/moments-plus:v0.0.1`
  - `docker pull ghcr.io/kai25888/tech-standards:latest`
  - `docker pull ghcr.io/kai25888/tech-standards:v0.0.1`

## Related Repositories
| 仓库 | 用途 | 路径 |
|------|------|------|
| moments | 原仓库，Web 应用主开发 | `/Users/andy/moments` |
| moments-app | 移动端 App（Capacitor 壳化） | `/Users/andy/moments-app` |

## High-Priority Follow-Ups
- Replace admin checks based on `id == 1` with an explicit role model.
- Add JWT expiration and stricter claims validation.
- Run the container as a non-root user.
- Add route-level auth middleware on the frontend for sensitive pages.

## Working Preferences
- Responses should be in Simplified Chinese.
- When persisting session knowledge, separate durable facts into `MEMORY.md` and date-specific facts into `daily-logs/YYYY-MM-DD.md`.
