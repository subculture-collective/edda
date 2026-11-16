#!/usr/bin/env python3
"""Compare #300's exact identities, alternatives and benefits to pinned PDF p87."""

import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
KEY = "srd-5.2.1:mechanic/"


def column(pdf, x):
    result = subprocess.run(
        ["pdftotext", "-f", "87", "-l", "87", "-x", str(x), "-y", "0",
         "-W", "310", "-H", "792", "-layout", str(pdf), "-"],
        capture_output=True, text=True, check=True,
    )
    return " ".join(re.sub(r"-\s+", "", result.stdout).split())


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_general_feats.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    assert "A feat can be taken only once unless its description states otherwise" in column(pdf, 0)
    text = column(pdf, 310).split("General Feats ", 1)[1].split(" Fighting Style Feats", 1)[0]
    assert text == (
        "Ability Score Improvement General Feat (Prerequisite: Level 4+) "
        "Increase one ability score of your choice by 2, or increase two ability scores of your choice by 1. "
        "This feat can’t increase an ability score above 20. Repeatable. You can take this feat more than once. "
        "Grappler General Feat (Prerequisite: Level 4+, Strength or Dexterity 13+) "
        "You gain the following benefits. Ability Score Increase. Increase your Strength or Dexterity score by 1, to a maximum of 20. "
        "Punch and Grab. When you hit a creature with an Unarmed Strike as part of the Attack action on your turn, "
        "you can use both the Damage and the Grapple option. You can use this benefit only once per turn. "
        "Attack Advantage. You have Advantage on attack rolls against a creature Grappled by you. "
        "Fast Wrestler. You don’t have to spend extra movement to move a creature Grappled by you if the creature is your size or smaller."
    ), ("source text changed", text)
    c = json.loads((ROOT / "internal/rules/srd521/testdata/character-general-feats.json").read_text())
    source = c["source"]
    assert source["sha256"] == PIN and source["license"] == "CC-BY-4.0"
    assert source["attribution"] and source["changes"]
    records = c["records"]
    assert len(records) == len({r["id"] for r in records}) == 11
    inventory = json.loads((ROOT / "docs/investigations/issue-74-character-inventory.json").read_text())
    expected = {r["id"]: r["source_page"] for r in inventory["records"] if r.get("catalogue_issue") == 300}
    actual = {s["inventory_id"]: s["page"] for r in records if r["kind"] == "feat" for s in r["sources"]}
    assert len(expected) == 2 and actual == expected
    for r in records:
        assert r["sources"] == [{"inventory_id": "srd-5.2.1:planning/feats/" + r["id"].split("/")[2], "page": 87}]
        assert r["kind"] in ("feat", "choice", "prerequisite")
    a, level, pattern, one, cond_one, two, cond_two, g, g_level, ability, g_choice = records
    assert a["parameters"] == {KEY+"feats/ability-score-improvement/maximum-ability-score": 20}
    assert KEY+"feats/repeatability/unrestricted" in a["required_mechanics"]
    assert level["requirement"]["parameters"] == g_level["requirement"]["parameters"] == {"minimum-level": 4}
    assert ability["requirement"] == {"key": KEY+"prerequisite/strength-or-dexterity-at-least", "record_refs": [], "parameters": {"minimum-score": 13}}
    assert pattern["select"] == one["select"] == g_choice["select"] == 1 and two["select"] == 2
    base = KEY+"feats/ability-score-improvement/"
    assert pattern["option_keys"] == [base+"one-score-plus-two", base+"two-scores-plus-one"]
    assert pattern["options"] == [
        {"key": base+"one-score-plus-two", "quantities": {base+"increase-per-score": 2}},
        {"key": base+"two-scores-plus-one", "quantities": {base+"increase-per-score": 1}},
    ]
    assert one["option_keys"] == two["option_keys"] == [KEY+"abilities/"+n for n in ("strength", "dexterity", "constitution", "intelligence", "wisdom", "charisma")]
    for choice, condition, option in ((one, cond_one, "one-score-plus-two"), (two, cond_two, "two-scores-plus-one")):
        assert choice["parent_id"] == pattern["id"] and condition["parent_id"] == choice["id"]
        assert condition["requirement"] == {"key": KEY+"prerequisite/ability-score-improvement/"+option+"-selected", "record_refs": [pattern["id"]], "parameters": {}}
    assert g_choice["option_keys"] == [KEY+"abilities/strength", KEY+"abilities/dexterity"]
    assert g["parameters"] == {KEY+"feats/maximum-times-taken": 1, KEY+"feats/grappler/ability-score-increase": 1, KEY+"feats/grappler/maximum-ability-score": 20, KEY+"feats/grappler/punch-and-grab/uses-per-turn": 1}
    assert g["required_mechanics"] == [KEY+k for k in (
        "feats/category/general", "feats/maximum-times-taken", "feats/grappler/ability-score-increase", "feats/grappler/maximum-ability-score",
        "feats/grappler/punch-and-grab/optional-on-hit-creature", "feats/grappler/punch-and-grab/unarmed-strike",
        "feats/grappler/punch-and-grab/part-of-attack-action-on-your-turn", "feats/grappler/punch-and-grab/use-both-damage-and-grapple-options",
        "feats/grappler/punch-and-grab/uses-per-turn", "feats/grappler/attack-advantage/on-attack-rolls-against-creature-grappled-by-you",
        "feats/grappler/fast-wrestler/no-extra-movement-to-move-creature-grappled-by-you", "feats/grappler/fast-wrestler/only-if-creature-your-size-or-smaller")]
    print("SRD " + PIN + ": exact two General feat identities on PDF p87; complete source text, prerequisites, alternatives, nested choices, caps and conditional benefits agree")


if __name__ == "__main__":
    main()
