#!/usr/bin/env python3
"""Run a disposable Edda mechanics smoke/QA gauntlet and write redacted artifacts.

By default this creates a fresh disposable local campaign. Supplying --campaign-id
is dangerous because the script intentionally mutates durable campaign state; it
is refused unless --i-understand-this-mutates-campaign is also provided. If no
--campaign-id is supplied and the base URL is not local, the script refuses to
create or mutate a remote disposable campaign unless
--i-understand-this-creates-remote-disposable-campaign is also provided.

This is a live-model smoke test, not a deterministic proof. It records the
requested action, HTTP results, state_changes, API snapshots, and optional
local DB counts for diagnosis.
"""

from __future__ import annotations

import argparse
import datetime as dt
import json
import os
import shlex
import signal
import subprocess
import time
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from typing import Any


DEFAULT_BASE_URL = os.environ.get("EDDA_BASE_URL", "http://localhost:18080")
DEFAULT_TIMEOUT = 900
DEFAULT_ARTIFACT_ROOT = Path("logs/edda")
DEFAULT_LOCAL_PASSWORD = os.environ.get("EDDA_GAUNTLET_PASSWORD", "local-gauntlet-password")


GAUNTLET_ACTIONS: list[tuple[str, str, list[str]]] = [
    ("quest_create", "Create and track a clear quest to escape this station through a safe control room.", ["quest:created"]),
    ("move_new_location", "Move into a brand-new test corridor now. If it is new, create it and move me there in this turn.", ["location:moved"]),
    ("known_fact", "Confirm and record this known fact: the test corridor contains active blue conduit veins. Make it visible to me.", ["world_fact:created"]),
    ("hp_status", "A harmless but real electrical shock grazes me; reduce my HP by 1 and mark the injury if appropriate.", ["player_character:hp_updated"]),
    ("inventory", "Add only one calibrated spanner to my inventory. Do not change HP, status, location, or quest state.", ["inventory_item:created"]),
    ("existing_move", "Move back to the previous connected room only, using move_player on the existing route. Do not change HP or inventory; if persistence cannot be guaranteed, keep the move provisional.", ["location:moved"]),
    ("quest_progress", "Advance the active escape quest only. Do not create a new quest; update the existing quest or complete one objective about reaching the control room.", ["objective:completed", "quest:updated"]),
    ("combat_start", "A weak hostile test drone blocks the route and combat is unavoidable now. Call initiate_combat with the weak drone immediately.", ["combat:started"]),
    ("combat_resolve", "Resolve the active combat only. If no combat is active, keep this provisional; failing here should mean combat_start did not happen.", ["combat:resolved"]),
]


def utc_slug() -> str:
    return dt.datetime.now(dt.UTC).strftime("%Y%m%dT%H%M%SZ")


def is_local_base_url(base_url: str) -> bool:
    host = urllib.parse.urlparse(base_url).hostname or ""
    return host in {"localhost", "127.0.0.1", "::1"}


def redacted(value: Any) -> Any:
    if isinstance(value, dict):
        out = {}
        for key, item in value.items():
            if key.lower() in {"authorization", "token", "password", "secret"}:
                out[key] = "<redacted>"
            else:
                out[key] = redacted(item)
        return out
    if isinstance(value, list):
        return [redacted(v) for v in value]
    return value


def slugify(value: str) -> str:
    safe = "".join(ch if ch.isalnum() else "-" for ch in value.lower())
    while "--" in safe:
        safe = safe.replace("--", "-")
    return safe.strip("-") or "model"


def parse_model_args(values: list[str] | None, models_file: str) -> list[str]:
    models: list[str] = []
    for value in values or []:
        models.extend(part.strip() for part in value.split(",") if part.strip())
    if models_file:
        for line in Path(models_file).read_text(encoding="utf-8").splitlines():
            line = line.strip()
            if line and not line.startswith("#"):
                models.append(line)
    deduped: list[str] = []
    seen: set[str] = set()
    for model in models:
        if model not in seen:
            seen.add(model)
            deduped.append(model)
    return deduped


def read_env_file(path: Path) -> dict[str, str]:
    if not path.exists():
        return {}
    result: dict[str, str] = {}
    for raw in path.read_text(encoding="utf-8").splitlines():
        line = raw.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, value = line.split("=", 1)
        key = key.strip()
        value = value.strip()
        if key.startswith("export "):
            key = key[len("export ") :].strip()
        if (value.startswith('"') and value.endswith('"')) or (value.startswith("'") and value.endswith("'")):
            value = value[1:-1]
        if key:
            result[key] = value
    return result


def port_from_base_url(base_url: str) -> int | None:
    parsed = urllib.parse.urlparse(base_url)
    return parsed.port


class Client:
    def __init__(self, base_url: str, token: str | None, timeout: int = DEFAULT_TIMEOUT):
        self.base_url = base_url.rstrip("/")
        self.token = token
        self.timeout = timeout

    def request(self, method: str, path: str, body: Any | None = None) -> tuple[int, Any | None]:
        headers = {"Accept": "application/json", "User-Agent": "edda-mechanics-gauntlet/1.0"}
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        if body is not None:
            headers["Content-Type"] = "application/json"
        data = None if body is None else json.dumps(body, ensure_ascii=False).encode("utf-8")
        req = urllib.request.Request(self.base_url + path, data=data, headers=headers, method=method)
        try:
            with urllib.request.urlopen(req, timeout=self.timeout) as resp:
                raw = resp.read().decode("utf-8")
                return resp.status, json.loads(raw) if raw else None
        except urllib.error.HTTPError as exc:
            raw = exc.read().decode("utf-8", errors="replace")
            try:
                payload = json.loads(raw) if raw else None
            except json.JSONDecodeError:
                payload = {"error": raw}
            return exc.code, payload
        except urllib.error.URLError as exc:
            return 0, {"error": str(exc.reason)}


def health_ok(base_url: str) -> bool:
    client = Client(base_url, None, timeout=5)
    status, payload = client.request("GET", "/api/healthz")
    return 200 <= status < 300 and isinstance(payload, dict) and payload.get("status") == "ok"


def wait_for_health(base_url: str, timeout_seconds: int) -> None:
    deadline = time.time() + timeout_seconds
    while time.time() < deadline:
        if health_ok(base_url):
            return
        time.sleep(1)
    raise RuntimeError(f"server did not become healthy within {timeout_seconds}s")


class ManagedServer:
    def __init__(self, args: argparse.Namespace, model: str, artifact_dir: Path):
        self.args = args
        self.model = model
        self.artifact_dir = artifact_dir
        self.process: subprocess.Popen[bytes] | None = None
        self.log_file = None

    def __enter__(self) -> "ManagedServer":
        if health_ok(self.args.base_url):
            raise SystemExit(
                f"Refusing to manage server for model rotation because {self.args.base_url} is already healthy. "
                "Stop the existing dev server first, or run the gauntlet without --model rotation."
            )

        env = os.environ.copy()
        env.update(read_env_file(Path(self.args.env_file)))
        env["EDDA_LLM_PROVIDER"] = self.args.model_provider
        if self.args.model_provider == "openrouter":
            env["EDDA_LLM_OPENROUTER_MODEL"] = self.model
        elif self.args.model_provider == "ollama":
            env["EDDA_LLM_OLLAMA_MODEL"] = self.model
        elif self.args.model_provider == "claude":
            env["EDDA_LLM_CLAUDE_MODEL"] = self.model
        else:
            raise SystemExit(f"Unsupported --model-provider {self.args.model_provider!r}")
        if port := port_from_base_url(self.args.base_url):
            env["EDDA_SERVER_PORT"] = str(port)

        self.artifact_dir.mkdir(parents=True, exist_ok=True)
        self.log_file = (self.artifact_dir / "server.log").open("ab")
        self.process = subprocess.Popen(
            shlex.split(self.args.server_command),
            stdout=self.log_file,
            stderr=subprocess.STDOUT,
            stdin=subprocess.DEVNULL,
            env=env,
            start_new_session=True,
        )
        try:
            wait_for_health(self.args.base_url, self.args.server_start_timeout)
        except Exception:
            self.stop()
            raise
        return self

    def __exit__(self, exc_type: Any, exc: Any, tb: Any) -> None:
        self.stop()

    def stop(self) -> None:
        if self.process and self.process.poll() is None:
            try:
                os.killpg(os.getpgid(self.process.pid), signal.SIGTERM)
                self.process.wait(timeout=20)
            except Exception:
                try:
                    os.killpg(os.getpgid(self.process.pid), signal.SIGKILL)
                except Exception:
                    pass
        if self.log_file:
            self.log_file.close()


def append_jsonl(path: Path, record: dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("a", encoding="utf-8") as f:
        f.write(json.dumps(redacted(record), ensure_ascii=False, sort_keys=True) + "\n")


def db_scalar(sql: str) -> str:
    cmd = ["docker", "compose", "exec", "-T", "postgres", "psql", "-U", "edda", "-d", "edda", "-t", "-A", "-c", sql]
    return subprocess.check_output(cmd, text=True).strip()


def db_counts(campaign_id: str) -> dict[str, int]:
    safe = campaign_id.replace("'", "''")
    sql = f"""
select json_build_object(
  'session_logs', (select count(*) from session_logs where campaign_id='{safe}'),
  'locations', (select count(*) from locations where campaign_id='{safe}'),
  'quests', (select count(*) from quests where campaign_id='{safe}'),
  'known_facts', (select count(*) from world_facts where campaign_id='{safe}' and player_known=true and superseded_by is null),
  'items', (select count(*) from items where campaign_id='{safe}')
)::text;
"""
    return json.loads(db_scalar(sql))


def create_disposable_world(client: Client) -> str:
    slug = utc_slug()
    body = {
        "name": f"Mechanics Gauntlet {slug}",
        "summary": "Disposable campaign for mechanics smoke/QA.",
        "profile": {
            "genre": "science-fantasy",
            "tone": "clear",
            "themes": ["mechanical consequence", "danger", "escape"],
            "world_type": "test corridor",
            "danger_level": "moderate",
            "political_complexity": "low",
        },
        "character_profile": {
            "name": "Gauntlet Runner",
            "concept": "field tester",
            "background": "Disposable operator",
            "personality": "careful and practical",
            "motivations": ["Verify mechanics"],
            "strengths": ["Adaptation"],
            "weaknesses": ["Overcaution"],
        },
        "rules_mode": "light",
    }
    status, payload = client.request("POST", "/api/v1/campaigns/start/world", body)
    if status != 200 or not isinstance(payload, dict) or "campaign" not in payload:
        raise SystemExit(f"failed to build disposable world: HTTP {status} {payload}")
    campaign = payload.get("campaign")
    if not isinstance(campaign, dict) or not campaign.get("id"):
        raise SystemExit(f"disposable world response missing campaign id: HTTP {status} {payload}")
    return str(campaign["id"])


def register_user(client: Client) -> str:
    email = f"edda-gauntlet-{utc_slug().lower()}@local.test"
    body = {"email": email, "password": DEFAULT_LOCAL_PASSWORD, "name": "Gauntlet Runner"}
    status, payload = client.request("POST", "/api/v1/auth/register", body)
    if status != 201 or not isinstance(payload, dict) or "token" not in payload:
        raise SystemExit(f"failed to register disposable user: HTTP {status} {payload}")
    return str(payload["token"])


def maybe_register_or_login(client: Client) -> str | None:
    try:
        return register_user(client)
    except SystemExit:
        return None


def fetch_snapshots(client: Client, campaign_id: str) -> dict[str, Any]:
    paths = {
        "character": "/character",
        "locations": "/locations",
        "quests": "/quests",
        "inventory": "/character/inventory",
        "facts": "/facts",
        "history": "/history",
        "time": "/time",
    }
    snapshots: dict[str, Any] = {}
    for name, path in paths.items():
        status, body = client.request("GET", f"/api/v1/campaigns/{campaign_id}{path}")
        snapshots[name] = {"status": status, "body": body}
    return snapshots


def state_change_labels(resp: dict[str, Any] | None) -> list[str]:
    labels: list[str] = []
    if not isinstance(resp, dict):
        return labels
    for change in resp.get("state_changes") or []:
        if isinstance(change, dict):
            labels.append(f"{change.get('entity_type', 'unknown')}:{change.get('change_type', 'unknown')}")
    return labels


def summarize_snapshot(snapshots: dict[str, Any]) -> dict[str, Any]:
    character = (snapshots.get("character") or {}).get("body") or {}
    def _count(name: str) -> int:
        body = (snapshots.get(name) or {}).get("body")
        if isinstance(body, list):
            return len(body)
        if isinstance(body, dict) and isinstance(body.get("entries"), list):
            return len(body["entries"])
        return 1 if body is not None else 0
    return {
        "character": {
            "name": character.get("name"),
            "hp": character.get("hp") or character.get("HP"),
            "max_hp": character.get("max_hp") or character.get("MaxHP"),
            "status": character.get("status"),
        },
        "counts": {k: _count(k) for k in ["locations", "quests", "inventory", "facts", "history"]},
        "time": (snapshots.get("time") or {}).get("body"),
    }


def run_gauntlet(client: Client, campaign_id: str, artifact_dir: Path, include_db_counts: bool, model: str = "") -> list[dict[str, Any]]:
    results: list[dict[str, Any]] = []
    for index, (label, action, expected_any) in enumerate(GAUNTLET_ACTIONS, 1):
        started = time.time()
        before_db = db_counts(campaign_id) if include_db_counts else None
        status, response = client.request("POST", f"/api/v1/campaigns/{campaign_id}/action", {"input": action})
        snapshots = fetch_snapshots(client, campaign_id)
        after_db = db_counts(campaign_id) if include_db_counts else None
        observed = state_change_labels(response)
        snapshot_statuses = {name: snap["status"] for name, snap in snapshots.items()}
        snapshots_ok = all(200 <= snap_status < 300 for snap_status in snapshot_statuses.values())
        ok = (200 <= status < 300) and snapshots_ok and (not expected_any or any(exp in observed for exp in expected_any))
        record = {
            "step": index,
            "label": label,
            "model": model,
            "action": action,
            "http_status": status,
            "ok": ok,
            "expected_any": expected_any,
            "state_changes": observed,
            "response": response,
            "snapshots": summarize_snapshot(snapshots),
            "api_snapshots": snapshots,
            "snapshot_statuses": snapshot_statuses,
            "db_counts_before": before_db,
            "db_counts_after": after_db,
            "seconds": round(time.time() - started, 1),
        }
        append_jsonl(artifact_dir / "gauntlet.jsonl", record)
        results.append(record)
        print(f"{index:02d} {label}: {'PASS' if ok else 'FAIL'} {observed or status}", flush=True)
    return results


def write_report(artifact_dir: Path, campaign_id: str, results: list[dict[str, Any]], disposable_mode: bool, local_base: bool, model: str = "") -> None:
    failed = [r for r in results if not r["ok"]]
    report = {
        "campaign_id": campaign_id,
        "disposable_mode": disposable_mode,
        "local_base": local_base,
        "model": model,
        "passed": len(results) - len(failed),
        "failed": len(failed),
        "artifact_dir": str(artifact_dir),
    }
    (artifact_dir / "summary.json").write_text(json.dumps(redacted(report), indent=2, sort_keys=True), encoding="utf-8")
    lines = [
        "# Edda mechanics gauntlet",
        "",
        f"- Campaign: `{campaign_id}`",
        f"- Model: `{model or 'current server config'}`",
        f"- Mode: `{'local disposable' if disposable_mode else 'mutating existing campaign'}`",
        f"- Disposable mode: `{disposable_mode}`",
        f"- Local base: `{local_base}`",
        f"- Passed: `{report['passed']}`",
        f"- Failed: `{report['failed']}`",
        "",
    ]
    for r in results:
        lines.extend([
            f"## Step {r['step']}: {r['label']}",
            f"- OK: `{r['ok']}`",
            f"- HTTP: `{r['http_status']}`",
            f"- Expected: `{', '.join(r['expected_any']) or 'none'}`",
            f"- Observed: `{', '.join(r['state_changes']) or 'none'}`",
            f"- DB counts before: `{r['db_counts_before']}`",
            f"- DB counts after: `{r['db_counts_after']}`",
            "",
        ])
    (artifact_dir / "report.md").write_text("\n".join(lines), encoding="utf-8")
    print(json.dumps(redacted(report), indent=2, sort_keys=True))


def parser() -> argparse.ArgumentParser:
    p = argparse.ArgumentParser(description="Run a disposable Edda mechanics smoke/QA gauntlet.")
    p.add_argument("--base-url", default=DEFAULT_BASE_URL)
    p.add_argument("--token", default=os.environ.get("EDDA_TOKEN", ""), help="Bearer token for an existing campaign or local auth")
    p.add_argument("--campaign-id", default="", help="Use an existing campaign (dangerous)")
    p.add_argument("--i-understand-this-mutates-campaign", action="store_true", help="Required when --campaign-id is supplied")
    p.add_argument("--i-understand-this-creates-remote-disposable-campaign", action="store_true", help="Required to create a disposable campaign on a non-local base URL")
    p.add_argument("--db-counts", action="store_true", help="Include diagnostic local DB counts (Docker/Postgres only)")
    p.add_argument("--timeout", type=int, default=DEFAULT_TIMEOUT, help="HTTP timeout in seconds (default: 900)")
    p.add_argument("--artifact-dir", default="", help="Override artifact directory (default logs/edda/gauntlet-*)")
    p.add_argument("--model", action="append", help="Model id to test. May be repeated or comma-separated. Requires --manage-server.")
    p.add_argument("--models-file", default="", help="File containing one model id per line. Requires --manage-server.")
    p.add_argument("--manage-server", action="store_true", help="Start and stop a local dev server for each --model with model env overrides.")
    p.add_argument("--model-provider", default="openrouter", choices=["openrouter", "ollama", "claude"], help="Provider env namespace to override when --manage-server is used.")
    p.add_argument("--env-file", default=".env", help="Environment file loaded for managed server runs (default: .env).")
    p.add_argument("--server-command", default="go run ./cmd/server", help="Command used to start a managed local server.")
    p.add_argument("--server-start-timeout", type=int, default=180, help="Seconds to wait for managed server health.")
    return p


def run_once(args: argparse.Namespace, artifact_dir: Path, model: str = "") -> list[dict[str, Any]]:
    disposable_mode = not args.campaign_id
    local_base = is_local_base_url(args.base_url)
    if args.campaign_id and not args.i_understand_this_mutates_campaign:
        raise SystemExit("Refusing to mutate an existing campaign without --i-understand-this-mutates-campaign.")
    if disposable_mode and not local_base and not args.i_understand_this_creates_remote_disposable_campaign:
        raise SystemExit("Refusing to create a remote disposable campaign without --i-understand-this-creates-remote-disposable-campaign.")
    if args.db_counts and not local_base:
        raise SystemExit("--db-counts requires a local base URL (localhost/127.0.0.1/::1).")

    artifact_dir.mkdir(parents=True, exist_ok=True)

    client = Client(args.base_url, args.token or None, timeout=args.timeout)
    if disposable_mode:
        if not client.token:
            client.token = maybe_register_or_login(client)
        campaign_id = create_disposable_world(client)
    else:
        if not client.token and not local_base:
            raise SystemExit("Existing campaign mode requires --token/EDDA_TOKEN or a no-auth local API.")
        campaign_id = args.campaign_id

    if args.db_counts:
        print("DB counts enabled: querying Docker/Postgres for diagnostics.", flush=True)
    results = run_gauntlet(client, campaign_id, artifact_dir, args.db_counts, model=model)
    write_report(artifact_dir, campaign_id, results, disposable_mode, local_base, model=model)
    return results


def write_rotation_summary(root_dir: Path, rotation_results: list[dict[str, Any]]) -> None:
    failed = [r for r in rotation_results if r["failed"] > 0]
    summary = {
        "models": rotation_results,
        "passed_models": len(rotation_results) - len(failed),
        "failed_models": len(failed),
        "artifact_dir": str(root_dir),
    }
    (root_dir / "models-summary.json").write_text(json.dumps(redacted(summary), indent=2, sort_keys=True), encoding="utf-8")
    lines = ["# Edda model rotation gauntlet", ""]
    lines.append("| Model | Passed | Failed | Artifact dir |")
    lines.append("|---|---:|---:|---|")
    for item in rotation_results:
        lines.append(f"| `{item['model']}` | {item['passed']} | {item['failed']} | `{item['artifact_dir']}` |")
    lines.append("")
    (root_dir / "models-report.md").write_text("\n".join(lines), encoding="utf-8")
    print(json.dumps(redacted(summary), indent=2, sort_keys=True))


def main() -> int:
    args = parser().parse_args()
    models = parse_model_args(args.model, args.models_file)
    if models and not args.manage_server:
        raise SystemExit("--model/--models-file requires --manage-server because Edda's model is server configuration, not a per-request API field.")
    if args.manage_server and not models:
        raise SystemExit("--manage-server requires at least one --model or --models-file entry.")
    if args.manage_server and not is_local_base_url(args.base_url):
        raise SystemExit("--manage-server requires a local base URL.")

    if not models:
        artifact_dir = Path(args.artifact_dir or DEFAULT_ARTIFACT_ROOT / f"gauntlet-{utc_slug()}")
        results = run_once(args, artifact_dir)
        return 1 if any(not r["ok"] for r in results) else 0

    root_dir = Path(args.artifact_dir or DEFAULT_ARTIFACT_ROOT / f"gauntlet-models-{utc_slug()}")
    root_dir.mkdir(parents=True, exist_ok=True)
    rotation_results: list[dict[str, Any]] = []
    any_failed = False
    for model in models:
        model_dir = root_dir / slugify(model)
        print(f"=== Running gauntlet for model: {model} ===", flush=True)
        with ManagedServer(args, model, model_dir):
            results = run_once(args, model_dir, model=model)
        failed_count = sum(1 for r in results if not r["ok"])
        any_failed = any_failed or failed_count > 0
        rotation_results.append({
            "model": model,
            "passed": len(results) - failed_count,
            "failed": failed_count,
            "artifact_dir": str(model_dir),
        })
    write_rotation_summary(root_dir, rotation_results)
    return 1 if any_failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
