#!/usr/bin/env python3
"""Compare issue #253 tool records with the tool entries in the pinned PDF.

The expected records are built from the printed entries, not from a list kept
in this script. Utilize and craft keys are the printed words in lower case,
without articles and apostrophes.
"""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
ISSUE = 253
BATCH = "data-tools-01"
COUNT = 12
VARIANT_COUNT = 15
PAGES = (93, 94)
HALF_WIDTH = 306
HEADER = re.compile(r"(?P<name>[A-Z][A-Za-z’]*(?: [A-Z][A-Za-z’]*)*) \((?P<price>\d+ GP|Varies)\)")
FIELD = re.compile(r"(Ability|Utilize|Craft|Variants):\s*(.*)")
GROUPS = {"Artisan’s Tools": "artisans-tools", "Other Tools": "other-tools"}


def column_lines(pdf):
    """(page, line) pairs in reading order: each page's left, then right column."""
    for page in PAGES:
        for left in (0, HALF_WIDTH):
            args = ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(left), "-y", "0",
                    "-W", str(HALF_WIDTH), "-H", "800", "-layout", str(pdf), "-"]
            text = subprocess.run(args, check=True, capture_output=True, text=True).stdout
            for line in text.splitlines():
                line = line.strip()
                if line and not line.endswith("System Reference Document 5.2.1"):
                    yield page, line


def join(first, second):
    """Join a wrapped line, closing a word the layout hyphenated."""
    if not first:
        return second
    if re.search(r"[a-z]-$", first) and second[:1].islower():
        return first[:-1] + second
    return first + " " + second


def entries(pdf):
    """Every tool entry from the Artisan's Tools heading to Adventuring Gear."""
    found, group, entry, field = [], None, None, None
    for page, line in column_lines(pdf):
        if line in GROUPS:
            group, entry = GROUPS[line], None
            continue
        if line == "Adventuring Gear":
            return found
        if group is None:
            continue
        header = HEADER.fullmatch(line)
        if header:
            entry = {"name": header["name"], "price": header["price"], "group": group, "pages": [page, page], "fields": {}}
            found.append(entry)
            field = None
            continue
        if entry is None:
            continue  # the Artisan's Tools introduction
        entry["pages"][1] = page
        start = FIELD.fullmatch(line)
        if start:
            field = start[1]
            assert field not in entry["fields"], (entry["name"], field)
            entry["fields"][field] = start[2]
            continue
        assert field and field != "Ability", ("unparsed tool line", entry["name"], line)
        entry["fields"][field] = join(entry["fields"][field], line)
    raise AssertionError("tool section end not found")


def key(text):
    words = [w for w in text.replace("’", "").lower().split() if w not in ("a", "an", "the")]
    return re.sub(r"[^a-z0-9]+", "-", " ".join(words)).strip("-")


def top_level(text):
    """Split a list on commas that are not inside parentheses."""
    return [part.strip() for part in re.split(r",\s*(?![^()]*\))", text)]


def utilize(text):
    tasks, end = [], 0
    for match in re.finditer(r"\(DC (\d+)\)", text):
        phrase = re.sub(r"^(?:, or |, )", "", text[end:match.start()]).strip()
        assert phrase, text
        tasks.append({"task": key(phrase), "dc": int(match[1])})
        end = match.end()
    assert tasks and end == len(text), text
    return tasks


def craft(text):
    if text is None:
        return []
    out = []
    for part in top_level(text):
        found = re.fullmatch(r"(.*?)(?: \(except (.*)\))?", part)
        excepted = []
        if found[2]:
            excepted = [key(re.sub(r"^(?:and )", "", x)) for x in re.split(r",\s*| (?=and )", found[2]) if x]
        out.append({"item": key(found[1]), "except": excepted})
    return out


def records(entry):
    """(variant slug, name, cost, keys, payload) for each form of one entry."""
    fields = entry["fields"]
    found = re.fullmatch(r"(\w+)\s+Weight: (—|Varies|\d+ lb\.)", fields["Ability"])
    assert found, (entry["name"], fields["Ability"])
    ability, weight = found[1].lower(), found[2]
    tasks, crafted = utilize(fields["Utilize"]), craft(fields.get("Craft"))
    forms = []
    if "Variants" in fields:
        assert entry["price"] == "Varies", entry["name"]
        for part in top_level(fields["Variants"]):
            variant = re.fullmatch(r"(.+?) \((\d+) (CP|SP|GP)(?:, (\d+) lb\.)?\)", part)
            assert variant, (entry["name"], part)
            # A variant prints its own weight only where the entry says Varies.
            assert (weight == "Varies") == (variant[4] is not None) and weight in ("Varies", "—"), (entry["name"], part)
            name = variant[1][0].upper() + variant[1][1:]
            forms.append((key(variant[1]), name, int(variant[2]), variant[3].lower(), int(variant[4] or 0)))
    else:
        assert entry["price"] != "Varies" and weight.endswith(" lb."), entry["name"]
        forms.append(("standard", entry["name"], int(entry["price"].removesuffix(" GP")), "gp", int(weight.removesuffix(" lb."))))
    out = []
    for slug, name, amount, currency, pounds in forms:
        payload = {"tool_group": entry["group"], "ability": ability, "utilize": tasks, "craft": crafted,
                   "weight": {"listed": pounds > 0, "pounds": pounds}, "source_pages": entry["pages"]}
        keys = ["tool.utilize", "tool.proficiency", f"tool.{entry['group']}", f"ability.{ability}"]
        keys += [f"utilize.{t['task']}-dc-{t['dc']}" for t in tasks]
        keys += [f"craft.{c['item']}" + ("-except-" + "-".join(c["except"]) if c["except"] else "") for c in crafted]
        out.append((slug, name, {"amount": amount, "currency": currency}, keys, payload))
    return out


def same(got, want):
    """Equality that does not treat False, 0 and 0.0 as the same value."""
    return json.dumps(got, sort_keys=True) == json.dumps(want, sort_keys=True)


def main():
    if len(sys.argv) != 2:
        raise SystemExit(f"usage: verify_srd_items_{ISSUE}.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    rows = entries(pdf)
    assert len(rows) == 25, [r["name"] for r in rows]
    printed = {r["name"]: r for r in rows}
    assert len(printed) == 25

    catalogue = json.loads((ROOT / f"internal/rules/srd521/testdata/items_{ISSUE}.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    tools = [r for r in inventory["records"] if r["family"] == "tools"]
    # The planning inventory lists exactly the printed entries and their pages.
    assert {(r["name"], r["source_page"]) for r in tools} == {(r["name"], r["pages"][0]) for r in rows}
    expected = [(r["id"], r["name"], r["source_page"]) for r in tools if r.get("catalogue_issue") == ISSUE]
    assert len(expected) == COUNT
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == ISSUE and catalogue["batch"] == BATCH
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected

    total = 0
    for f in families:
        name = f["name"]
        slug = key(name).replace(" ", "-")
        entry = printed[name]
        pages = [entry["pages"][0]] * 2
        assert f["id"] == f"srd-5.2.1:content/items/{slug}" and f["kind"] == "mundane" and f["category"] == "tool", name
        assert f["base"] == {"options": [], "excluded": []} and len(f["presentations"]) == 1, name
        assert f["presentations"][0] == {"id": f"srd-5.2.1:content/tools/{slug}", "pages": pages}, name
        forms = records(entry)
        assert len(f["variants"]) == len(forms), (name, len(forms))
        for v, (vslug, vname, price, keys, payload) in zip(f["variants"], forms):
            total += 1
            assert v["id"] == f["id"] + "/variant/" + vslug and v["name"] == vname and v["rarity"] == "none", (name, vname)
            assert v["attunement"] == {"mode": "none", "restrictions": []}, vname
            assert same(v["charges"], {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}), vname
            assert v["equip"] == {"mode": "none", "restrictions": []} and v["use"] == {"mode": "utilize", "restrictions": []}, vname
            cost = dict(price, kind="listed", source_pages=pages)
            assert same(v["cost"], cost), (vname, v["cost"], cost)
            assert v["required_mechanic_keys"] == keys, (vname, v["required_mechanic_keys"], keys)
            assert same(v["details"], {"tool": payload}), (vname, v["details"], payload)
    assert total == VARIANT_COUNT, total
    print(f"ok: {len(families)} tool families and {total} variants match their printed entries")


main()
