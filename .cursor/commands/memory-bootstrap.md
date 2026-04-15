# Load moments memory and recent progress

Use the repository-specific `moments-project-memory` skill first, then read the repository memory files to produce a startup brief before working.

Complete this workflow:

1. Read `.cursor/skills/moments-project-memory/SKILL.md`.
2. Read `MEMORY.md` if it exists.
3. Read today's `daily-logs/YYYY-MM-DD.md` if it exists.
4. If today's log does not exist, read the most recent daily log instead.
5. Summarize:
   - active long-term rules
   - current project status
   - blockers or risks
   - likely next steps
   - rules that should be followed in this session
6. Mention exactly which files were used to build the brief.

Important:

- This is a read-only startup flow. Do not modify `MEMORY.md`, daily logs, project skills, or other files in this command.
- Treat `moments-project-memory` as the highest-priority source for repository-specific rules.
- In this bootstrap command, use `moments-project-memory` only as a read-only source of facts, paths, and conventions. Do not execute its update or merge workflow.
- Preserve exact identifiers such as branch names, tags, image names, pull commands, workflow conventions, and paths.
- Keep the startup brief concise and directly useful for continuing work in this repository.
