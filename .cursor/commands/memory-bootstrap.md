# Load project memory and recent progress

Use the same read-only startup workflow as the global `memory-bootstrap` skill: load repository-specific memory first, then `MEMORY.md`, then today’s daily log (or the most recent log). In this repository, the project-specific memory skill lives at `.cursor/skills/moments-project-memory/SKILL.md`.

Complete this workflow:

1. Resolve the current repository or workspace root.
2. Read `.cursor/skills/moments-project-memory/SKILL.md` (this repo’s project memory skill). Treat it only as a read-only source of facts, paths, and conventions for this step; do not execute its update or merge workflow.
3. Read `MEMORY.md` if it exists.
4. Read today's `daily-logs/YYYY-MM-DD.md` if it exists.
5. If today's daily log does not exist, read the most recent daily log instead.
6. Summarize the current operating context as a startup brief that includes:
   - active long-term rules
   - current status
   - blockers or risks
   - likely next steps
   - work rules the agent should follow in this session
7. Mention exactly which files were used to produce the brief.

Important:

- This is a read-only bootstrap flow. Do not modify `MEMORY.md` or daily logs in this command.
- Do not modify repository-specific memory skills or other project files in this command.
- These read-only limits apply to this bootstrap command itself. If the user later asks to update files, switch to the correct write workflow (for example `project-memory-journal` or this repo’s `moments-project-memory` update workflow) instead of extending this command.
- If `MEMORY.md` or today's daily log is missing, say so clearly instead of guessing.
- Prefer repository-specific memory rules from `moments-project-memory` over any generic default when both exist.
- Keep the brief concise and directly useful for starting work in a new window.
