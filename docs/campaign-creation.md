# Campaign creation persistence

The direct `POST /campaigns` endpoint reports success only after a requested rules mode has been saved and the campaign has been read back. A failed rules write or read returns HTTP 500. An omitted rules mode uses the database default. A failed follow-up write can leave the initial campaign row; callers should inspect their campaign list before retrying this direct endpoint.

World creation requires a character name that is not blank. `world.Orchestrator.Run` refuses one before any model request or write, for both rules versions and for the TUI and the API; `POST /campaigns/start/world` answers HTTP 400. See [data-checks/blank-character-names.md](data-checks/blank-character-names.md) for characters saved before this rule.
