#!/usr/bin/env python3
"""Benchmark assembled Edda prompt fixtures across OpenRouter models.

This is a live-model prompt benchmark. It does not call the Edda API and does
not mutate campaign state. It reads the Markdown fixtures in
docs/playtest/assembled-prompt-fixtures.md, sends each message array to each
configured OpenRouter model, and records latency, usage, text output, and any
tool calls requested by the model.

Tool schemas are intentionally generic stubs derived from the fixture's
representative filtered tool names. The benchmark measures first-response tool
selection and narrative behavior only; use the live API gauntlet for full
database-backed turn evaluation.
"""

from __future__ import annotations

import argparse
import csv
import datetime as dt
import json
import os
import re
import time
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any


DEFAULT_FIXTURES = Path("docs/playtest/assembled-prompt-fixtures.md")
DEFAULT_ARTIFACT_ROOT = Path("logs/edda")
DEFAULT_OPENROUTER_URL = "https://openrouter.ai/api/v1/chat/completions"
DEFAULT_TIMEOUT = 180
DEFAULT_MODELS = [
    "qwen/qwen3-235b-a22b-2507",
    "openai/gpt-5-nano",
]


@dataclass
class Fixture:
    slug: str
    title: str
    purpose: str = ""
    expected: list[str] = field(default_factory=list)
    tool_names: list[str] = field(default_factory=list)
    messages: list[dict[str, Any]] = field(default_factory=list)


def utc_slug() -> str:
    return dt.datetime.now(dt.UTC).strftime("%Y%m%dT%H%M%SZ")


def slugify(value: str) -> str:
    value = value.lower()
    value = re.sub(r"[^a-z0-9]+", "-", value)
    return value.strip("-") or "fixture"


def split_fixture_sections(markdown: str) -> list[tuple[str, str]]:
    matches = list(re.finditer(r"^## Fixture\s+\d+\s+—\s+(.+)$", markdown, re.MULTILINE))
    sections: list[tuple[str, str]] = []
    for idx, match in enumerate(matches):
        start = match.end()
        end = matches[idx + 1].start() if idx + 1 < len(matches) else len(markdown)
        sections.append((match.group(1).strip(), markdown[start:end]))
    return sections


def extract_json_blocks(section: str) -> list[Any]:
    blocks: list[Any] = []
    for raw in re.findall(r"```json\n(.*?)\n```", section, flags=re.DOTALL):
        blocks.append(json.loads(raw))
    return blocks


def extract_purpose(section: str) -> str:
    match = re.search(r"^Purpose:\s*(.+)$", section, flags=re.MULTILINE)
    return match.group(1).strip() if match else ""


def extract_expected(section: str) -> list[str]:
    match = re.search(r"^Expected outcome:\n((?:- .+\n?)+)", section, flags=re.MULTILINE)
    if not match:
        return []
    return [line[2:].strip() for line in match.group(1).splitlines() if line.startswith("- ")]


def looks_like_messages(value: Any) -> bool:
    return (
        isinstance(value, list)
        and len(value) > 0
        and all(isinstance(item, dict) and "role" in item and "content" in item for item in value)
    )


def looks_like_tool_names(value: Any) -> bool:
    return isinstance(value, list) and all(isinstance(item, str) for item in value)


def load_fixtures(path: Path, *, expand_gamemaster: bool = True) -> list[Fixture]:
    markdown = path.read_text(encoding="utf-8")
    gamemaster = ""
    if expand_gamemaster:
        gm_path = Path("internal/prompt/gamemaster.txt")
        if gm_path.exists():
            gamemaster = gm_path.read_text(encoding="utf-8").rstrip()

    fixtures: list[Fixture] = []
    for title, section in split_fixture_sections(markdown):
        blocks = extract_json_blocks(section)
        tool_names: list[str] = []
        messages: list[dict[str, Any]] = []
        for block in blocks:
            if looks_like_tool_names(block):
                tool_names = list(block)
            elif looks_like_messages(block):
                messages = block
        if not messages:
            continue
        if gamemaster:
            messages = [dict(msg) for msg in messages]
            messages[0]["content"] = str(messages[0]["content"]).replace("{{GAMEMASTER_PROMPT}}", gamemaster)
        fixtures.append(
            Fixture(
                slug=slugify(title),
                title=title,
                purpose=extract_purpose(section),
                expected=extract_expected(section),
                tool_names=tool_names,
                messages=messages,
            )
        )
    return fixtures


def tool_stub(name: str) -> dict[str, Any]:
    return {
        "type": "function",
        "function": {
            "name": name,
            "description": f"Representative Edda tool stub for benchmark fixture: {name}.",
            "parameters": {
                "type": "object",
                "properties": {
                    "summary": {
                        "type": "string",
                        "description": "Short description of the intended state mutation or resolution.",
                    }
                },
                "additionalProperties": True,
            },
        },
    }


def openrouter_headers(api_key: str) -> dict[str, str]:
    headers = {
        "Accept": "application/json",
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
        "User-Agent": "edda-prompt-benchmark/1.0",
    }
    referer = os.environ.get("OPENROUTER_HTTP_REFERER", "")
    title = os.environ.get("OPENROUTER_X_TITLE", "Edda Prompt Benchmark")
    if referer:
        headers["HTTP-Referer"] = referer
    if title:
        headers["X-Title"] = title
    return headers


def request_openrouter(
    *,
    api_key: str,
    model: str,
    fixture: Fixture,
    timeout: int,
    temperature: float,
    max_tokens: int,
    send_tools: bool,
) -> dict[str, Any]:
    body: dict[str, Any] = {
        "model": model,
        "messages": fixture.messages,
        "temperature": temperature,
        "max_tokens": max_tokens,
    }
    if send_tools and fixture.tool_names:
        body["tools"] = [tool_stub(name) for name in fixture.tool_names]

    req = urllib.request.Request(
        DEFAULT_OPENROUTER_URL,
        data=json.dumps(body, ensure_ascii=False).encode("utf-8"),
        headers=openrouter_headers(api_key),
        method="POST",
    )
    started = time.monotonic()
    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            raw = resp.read().decode("utf-8")
            latency_ms = int((time.monotonic() - started) * 1000)
            payload = json.loads(raw) if raw else {}
            return normalize_success(payload, latency_ms)
    except urllib.error.HTTPError as exc:
        raw = exc.read().decode("utf-8", errors="replace")
        latency_ms = int((time.monotonic() - started) * 1000)
        return {"ok": False, "status": exc.code, "latency_ms": latency_ms, "error": raw[:4000]}
    except urllib.error.URLError as exc:
        latency_ms = int((time.monotonic() - started) * 1000)
        return {"ok": False, "status": 0, "latency_ms": latency_ms, "error": str(exc.reason)}


def normalize_success(payload: dict[str, Any], latency_ms: int) -> dict[str, Any]:
    choice = (payload.get("choices") or [{}])[0]
    message = choice.get("message") or {}
    tool_calls = []
    for call in message.get("tool_calls") or []:
        function = call.get("function") or {}
        tool_calls.append(
            {
                "id": call.get("id"),
                "name": function.get("name"),
                "arguments": function.get("arguments"),
            }
        )
    return {
        "ok": True,
        "status": 200,
        "latency_ms": latency_ms,
        "finish_reason": choice.get("finish_reason"),
        "content": message.get("content") or "",
        "tool_calls": tool_calls,
        "usage": payload.get("usage") or {},
        "provider": payload.get("provider"),
        "raw_model": payload.get("model"),
    }


def append_jsonl(path: Path, record: dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("a", encoding="utf-8") as f:
        f.write(json.dumps(record, ensure_ascii=False, sort_keys=True) + "\n")


def write_csv(path: Path, records: list[dict[str, Any]]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8", newline="") as f:
        writer = csv.DictWriter(
            f,
            fieldnames=[
                "model",
                "fixture",
                "ok",
                "status",
                "latency_ms",
                "prompt_tokens",
                "completion_tokens",
                "total_tokens",
                "tool_call_names",
                "content_chars",
                "error",
            ],
        )
        writer.writeheader()
        for record in records:
            result = record["result"]
            usage = result.get("usage") or {}
            writer.writerow(
                {
                    "model": record["model"],
                    "fixture": record["fixture"]["slug"],
                    "ok": result.get("ok"),
                    "status": result.get("status"),
                    "latency_ms": result.get("latency_ms"),
                    "prompt_tokens": usage.get("prompt_tokens"),
                    "completion_tokens": usage.get("completion_tokens"),
                    "total_tokens": usage.get("total_tokens"),
                    "tool_call_names": ";".join(tc.get("name") or "" for tc in result.get("tool_calls") or []),
                    "content_chars": len(result.get("content") or ""),
                    "error": (result.get("error") or "")[:300],
                }
            )


def write_markdown(path: Path, records: list[dict[str, Any]]) -> None:
    lines = ["# Edda prompt benchmark", ""]
    by_model: dict[str, list[dict[str, Any]]] = {}
    for record in records:
        by_model.setdefault(record["model"], []).append(record)
    for model, model_records in by_model.items():
        ok_count = sum(1 for r in model_records if r["result"].get("ok"))
        avg_latency = sum(r["result"].get("latency_ms") or 0 for r in model_records) / max(len(model_records), 1)
        lines.extend([f"## `{model}`", "", f"- Success: `{ok_count}/{len(model_records)}`", f"- Avg latency: `{avg_latency:.0f} ms`", ""])
        lines.append("| Fixture | OK | Latency | Tool calls | Finish | Notes |")
        lines.append("|---|---:|---:|---|---|---|")
        for record in model_records:
            result = record["result"]
            names = ", ".join(tc.get("name") or "?" for tc in result.get("tool_calls") or []) or "—"
            note = ""
            if not result.get("ok"):
                note = str(result.get("error") or "")[:120].replace("|", "\\|")
            elif not names or names == "—":
                note = "no tool call"
            lines.append(
                f"| {record['fixture']['slug']} | {result.get('ok')} | {result.get('latency_ms')} | {names} | {result.get('finish_reason') or '—'} | {note} |"
            )
        lines.append("")
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text("\n".join(lines), encoding="utf-8")


def model_list(args: argparse.Namespace) -> list[str]:
    models: list[str] = []
    if args.model:
        for value in args.model:
            models.extend(part.strip() for part in value.split(",") if part.strip())
    if args.models_file:
        for line in Path(args.models_file).read_text(encoding="utf-8").splitlines():
            line = line.strip()
            if line and not line.startswith("#"):
                models.append(line)
    return models or DEFAULT_MODELS


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--fixtures", type=Path, default=DEFAULT_FIXTURES)
    parser.add_argument("--model", action="append", help="OpenRouter model id. May be repeated or comma-separated.")
    parser.add_argument("--models-file", help="File containing one OpenRouter model id per line.")
    parser.add_argument("--fixture", action="append", help="Fixture slug/title substring to include. May be repeated.")
    parser.add_argument("--limit", type=int, default=0, help="Limit number of fixtures after filtering.")
    parser.add_argument("--repeat", type=int, default=1)
    parser.add_argument("--timeout", type=int, default=DEFAULT_TIMEOUT)
    parser.add_argument("--temperature", type=float, default=0.2)
    parser.add_argument("--max-tokens", type=int, default=900)
    parser.add_argument("--artifact-root", type=Path, default=DEFAULT_ARTIFACT_ROOT)
    parser.add_argument("--no-tools", action="store_true", help="Do not send representative tool stubs.")
    parser.add_argument("--dry-run", action="store_true", help="Parse fixtures and print planned runs without calling OpenRouter.")
    args = parser.parse_args()

    fixtures = load_fixtures(args.fixtures)
    if args.fixture:
        needles = [n.lower() for n in args.fixture]
        fixtures = [f for f in fixtures if any(n in f.slug.lower() or n in f.title.lower() for n in needles)]
    if args.limit > 0:
        fixtures = fixtures[: args.limit]
    models = model_list(args)

    if not fixtures:
        raise SystemExit("no fixtures matched")
    if args.dry_run:
        print(f"Would run {len(fixtures)} fixture(s) x {len(models)} model(s) x {args.repeat} repeat(s)")
        for fixture in fixtures:
            print(f"- {fixture.slug}: {len(fixture.messages)} messages, {len(fixture.tool_names)} tool stubs")
        for model in models:
            print(f"- model: {model}")
        return 0

    api_key = os.environ.get("OPENROUTER_API_KEY") or os.environ.get("EDDA_LLM_OPENROUTER_APIKEY")
    if not api_key:
        raise SystemExit("OPENROUTER_API_KEY or EDDA_LLM_OPENROUTER_APIKEY is required")

    run_dir = args.artifact_root / f"prompt-benchmark-{utc_slug()}"
    jsonl_path = run_dir / "results.jsonl"
    csv_path = run_dir / "summary.csv"
    md_path = run_dir / "summary.md"
    records: list[dict[str, Any]] = []

    for model in models:
        for fixture in fixtures:
            for iteration in range(1, args.repeat + 1):
                print(f"running model={model} fixture={fixture.slug} iteration={iteration}/{args.repeat}")
                result = request_openrouter(
                    api_key=api_key,
                    model=model,
                    fixture=fixture,
                    timeout=args.timeout,
                    temperature=args.temperature,
                    max_tokens=args.max_tokens,
                    send_tools=not args.no_tools,
                )
                record = {
                    "timestamp": dt.datetime.now(dt.UTC).replace(microsecond=0).isoformat().replace("+00:00", "Z"),
                    "model": model,
                    "iteration": iteration,
                    "fixture": {
                        "slug": fixture.slug,
                        "title": fixture.title,
                        "purpose": fixture.purpose,
                        "expected": fixture.expected,
                        "tool_names": fixture.tool_names,
                    },
                    "request": {
                        "temperature": args.temperature,
                        "max_tokens": args.max_tokens,
                        "sent_tools": not args.no_tools,
                        "message_count": len(fixture.messages),
                    },
                    "result": result,
                }
                append_jsonl(jsonl_path, record)
                records.append(record)

    write_csv(csv_path, records)
    write_markdown(md_path, records)
    print(f"wrote {jsonl_path}")
    print(f"wrote {csv_path}")
    print(f"wrote {md_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
