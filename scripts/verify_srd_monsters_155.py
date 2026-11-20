#!/usr/bin/env python3
"""Compare the issue #155 stat blocks with the pinned SRD 5.2.1 PDF.

Each record is printed back in the source's stat-block notation from its typed
fields and retained entry text, then compared with the PDF block word for word
(ignoring only line breaks and line-break hyphens). Entry headings are read
separately from the PDF's bold-italic runs, so a missed or merged entry fails.
Every typed value and every dependency key is tied to a phrase its own entry
prints, and the record list and pages must equal the #82 decision inventory.
"""
import hashlib
import html
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
# Record name -> (first page, last page, heading that follows its block in the PDF).
BLOCKS = {
    "Black Bear": (345, 345, "Blood Hawk"),
    "Black Dragon Wyrmling": (263, 264, "Young Black Dragon"),
    "Black Pudding": (265, 266, "Blink Dog"),
    "Blink Dog": (266, 266, "Blue Dragons"),
    "Blood Hawk": (345, 346, "Boar"),
    "Blue Dragon Wyrmling": (266, 266, "Young Blue Dragon"),
    "Boar": (346, 346, "Brown Bear"),
    "Bone Devil": (267, 268, "Brass Dragons"),
    "Brass Dragon Wyrmling": (268, 268, "Young Brass Dragon"),
    "Bronze Dragon Wyrmling": (269, 270, "Young Bronze Dragon"),
    "Brown Bear": (346, 346, "Camel"),
    "Bugbear Stalker": (271, 271, "Bugbear Warrior"),
}
SIZES = "Tiny|Small|Medium|Large|Huge|Gargantuan"
ABILITY = {"str": "Strength", "dex": "Dexterity", "con": "Constitution",
           "int": "Intelligence", "wis": "Wisdom", "cha": "Charisma"}
# Dependency keys whose phrase cannot be derived from the key's own name. Each
# phrase must be printed by the entry or effect that names the key.
PHRASES = {
    "save/advantage-against-magic": "Advantage on saving throws against spells and other magical effects",
    "effect/revive": "gains a new body instantly, reviving with all its Hit Points",
    "limit/recharge": "(Recharge",
    "save/half-damage": "Success: Half damage",
    "save/repeat": "repeats the save",
    "save/second-failure": "Second Failure:",
    "effect/end-on-damage-or-wake": "This effect ends for the target if it takes damage",
    "effect/push": "is pushed up to 30 feet straight away from the dragon",
    "effect/destroy-nonmagical-ammunition": "Nonmagical ammunition is destroyed",
    "effect/weapon-attack-penalty": "cumulative −1 penalty to attack rolls",
    "effect/destroy-weapon": "The weapon is destroyed if the penalty reaches −5",
    "effect/eat-through-material": "eat through 2 feet of nonmagical wood or metal",
    "effect/armor-ac-penalty": "−1 penalty to the AC it offers",
    "effect/destroy-armor": "The armor is destroyed if the penalty reduces its AC to 10",
    "reaction/trigger": "Trigger:",
    "requirement/size-large-or-medium": "Large or Medium",
    "requirement/minimum-hit-points": "10+ Hit Points",
    "effect/split": "splits into two new Black Puddings",
    "effect/smaller-size": "one size smaller than the original pudding",
    "effect/divide-hit-points": "divided evenly between the new puddings (round down)",
    "effect/teleport": "teleports up to 40 feet",
    "requirement/unoccupied-visible-space": "an unoccupied space it can see",
    "effect/no-hit-point-recovery": "can’t regain Hit Points",
    "attack/advantage": "has Advantage on",
    "requirement/ally-within-5-feet": "at least one of the hawk’s allies is within 5 feet of the creature",
    "requirement/ally-not-incapacitated": "the ally doesn’t have the Incapacitated condition",
    "requirement/bloodied": "While Bloodied",
    "effect/bloodied-piercing-damage": "Piercing damage if the target is Bloodied",
    "requirement/target-bloodied": "if the target is Bloodied",
    "requirement/target-size": "or smaller creature",
    "requirement/charge-distance": "moved 20+ feet straight toward it immediately before the hit",
    "effect/extra-piercing-damage": "an extra 3 (1d6) Piercing damage",
    "attack/advantage-against-grappled-target": "(with Advantage if the target is Grappled by the bugbear)",
    "target/creature": "one Medium or smaller creature the bugbear can see within 10 feet",
}
ATTACK_KEYS = ("action/melee-attack", "action/ranged-attack", "action/melee-or-ranged-attack")


def run(*args):
    return subprocess.run(args, check=True, capture_output=True, text=True).stdout


def column_lines(pdf, page):
    """Read columns independently so continuation lines do not interleave."""
    out = []
    for x in (0, 306):
        text = run("pdftotext", "-f", str(page), "-l", str(page), "-x", str(x), "-y", "0",
                   "-W", "306", "-H", "792", "-layout", str(pdf), "-")
        for line in text.splitlines():
            line = " ".join(line.replace("­", "").replace("\x0c", "").split())
            if line and not re.fullmatch(r"\d+ System Reference Document 5\.2\.1|(MOD SAVE ?)+", line):
                out.append((page, line))
    return out


def rejoin(lines):
    """Join wrapped lines; a line-break hyphen is dropped, a printed hyphen kept."""
    text = ""
    for line in lines:
        if text.endswith("-"):
            keep = re.search(r"\d-$|foot-$", text)  # a printed hyphen stays, the break adds no space
            text = (text if keep else text[:-1]) + line
        else:
            text += (" " if text else "") + line
    return text


def heading_runs(pdf, page):
    """Text runs in reading order with their bold-italic entry heading, if any."""
    runs = []
    for m in re.finditer(r'<text top="(\d+)" left="(\d+)"[^>]*>(.*?)</text>', run(
            "pdftohtml", "-xml", "-i", "-f", str(page), "-l", str(page), "-stdout", str(pdf))):
        top, left, body = int(m[1]), int(m[2]), m[3]
        head = re.match(r"<i><b>(.*?)</b>", body)
        plain = html.unescape(re.sub(r"<[^>]+>", "", body)).strip()
        runs.append((left >= 445, top, left, plain, html.unescape(head[1]).strip().rstrip(".") if head else None))
    return sorted(runs)


def signed(n):
    return ("−%d" % -n) if n < 0 else "+%d" % n


def dice(d):
    mod = " %s %d" % ("+" if d["modifier"] > 0 else "−", abs(d["modifier"])) if d["modifier"] else ""
    return "%d (%dd%d%s)" % (d["average"], d["dice"], d["sides"], mod)


def title(s):
    return " ".join(w.capitalize() for w in s.replace("-", " ").split())


def heading(e):
    uses = e.get("uses")
    if not uses:
        return e["name"]
    assert uses["lair_per_day"] == 0 and uses["per_day"] == 0, e["name"]  # this batch prints only Recharge
    return "%s (Recharge %d–6)" % (e["name"], uses["recharge"])


def render(m):
    """Print one record in source order from typed fields and entry text."""
    kind = title(m["type"]) + (" (%s)" % ", ".join(m["tags"]) if m.get("tags") else "")
    speed = "Speed %d ft." % m["speeds"]["walk"] + "".join(
        ", %s %d ft." % (title(k), v) for k, v in sorted(m["speeds"].items()) if k != "walk")
    out = [
        "%s %s, %s" % (" or ".join(title(s) for s in m["sizes"]), kind, m["alignment"]),
        "AC %d Initiative %s (%d)" % (m["ac"], signed(m["initiative"]["modifier"]), m["initiative"]["score"]),
        "HP " + dice(m["hp"]), speed,
    ]
    for key in ("str", "dex", "con", "int", "wis", "cha"):
        a = m["abilities"][key]
        out.append("%s %d %s %s" % (title(key), a["score"], signed(a["modifier"]), signed(a["save"])))
    if "skills" in m:
        out.append("Skills " + ", ".join("%s %s" % (title(k), signed(v)) for k, v in sorted(m["skills"].items())))
    if "resistances" in m:
        out.append("Resistances " + ", ".join(title(v) for v in m["resistances"]))
    immune = [", ".join(title(v) for v in m[k]) for k in ("immunities", "condition_immunities") if k in m]
    if immune:
        out.append("Immunities " + "; ".join(immune))
    if m["gear"]:
        out.append("Gear " + ", ".join(m["gear"]))
    notes = m["senses"].get("notes", {})
    special = ", ".join("%s %d ft.%s" % (title(k), v, " (%s)" % notes[k] if k in notes else "")
                        for k, v in sorted(m["senses"]["special"].items()))
    out.append("Senses " + (special + "; " if special else "") + "Passive Perception %d" % m["senses"]["passive_perception"])
    out.append("Languages " + (", ".join(m["languages"]) or "None") + ("; telepathy %d ft." % m["telepathy"] if m.get("telepathy") else ""))
    out.append("CR %s (XP %s; PB +%d)" % (m["cr"], "{:,}".format(m["xp"]), m["proficiency_bonus"]))
    for label, key in (("Traits", "traits"), ("Actions", "actions"), ("Bonus Actions", "bonus_actions"), ("Reactions", "reactions")):
        if m.get(key):
            out += [label] + [e["text"] for e in m[key]]
    return " ".join(out)


def sections(m):
    return m["actions"] + m.get("bonus_actions", []) + m.get("reactions", [])


def check_key(where, key, text):
    """A dependency key must be backed by a phrase the given source text prints."""
    kind, _, name = key.partition("/")
    if key in PHRASES:
        phrase = PHRASES[key]
    elif kind == "effect" and name.endswith("-damage"):
        phrase = title(name[:-len("-damage")]) + " damage"
    elif kind == "save" and name in ABILITY:
        phrase = ABILITY[name] + " Saving Throw"
    elif kind == "area":
        phrase = title(name)
    elif kind == "condition":
        phrase = title(name) + " condition"
    elif key == "limit/per-day":
        phrase = "/Day"
    else:
        raise AssertionError((where, key, "no source phrase reviewed for this key"))
    assert phrase.lower() in text.lower(), (where, key, phrase)


def check_entry(m, e, effects):
    """Each typed payload must be printed by the entry's own source text."""
    text, where = e["text"], (m["name"], e["name"])
    assert text.startswith(heading(e) + ". "), where
    a = e.get("attack")
    if e.get("key") == "action/melee-attack":
        want = r"Melee Attack Roll: \+%d(?: \([^)]*\))?, reach %d ft\. Hit: %s %s damage" % (
            a["bonus"], a["reach"], re.escape(dice(a["damage"])), title(a["damage_type"]))
        assert a["kind"] == "melee" and re.search(want, text), where
    else:
        assert a is None, where
        # The one untyped attack roll is the Javelin's, whose reach and ranges are both printed.
        assert ("Attack Roll" in text) == (e.get("key") == "action/melee-or-ranged-attack" and e["name"] != "Multiattack"), where
    if e.get("key") == "action/melee-or-ranged-attack" and e["name"] != "Multiattack":
        assert text == "Javelin. Melee or Ranged Attack Roll: +5, reach 10 ft. or range 30/120 ft. Hit: %s Piercing damage." % dice(effects[0]["dice"]), where
    if e.get("save"):
        assert e["key"] == "action/saving-throw", where
        assert "%s Saving Throw: DC %d, " % (ABILITY[e["save"]["ability"]], e["save"]["dc"]) in text, where
    if e.get("target"):
        t = e["target"]
        want = {
            "line": r"each creature in an? %d-foot-long, %d-foot-wide Line" % (t["size"], t["width"]),
            "cone": r"each creature in a %d-foot Cone" % t["size"],
            "creature": r"one [^.]*creature the \w+ can see within %d feet" % t["range"],
        }[t["shape"]]
        assert re.search(want, text), where
    printed = re.findall(r"\d+ \(\d+d\d+(?: [+−] \d+)?\)", text)
    if effects is None:  # a trait types its one dice expression itself
        assert printed == ([dice(e["dice"])] if e.get("dice") else []), where
        assert not e.get("more_dice"), where
        for key in e.get("requires", []):
            check_key(where, key, text)
        return
    # Every dice expression in an action is typed in order: the attack's own first, then each further effect.
    typed = [dice(a["damage"])] if a else []
    rest = effects[1:] if a else effects
    assert printed == typed + [dice(x["dice"]) for x in rest if x.get("dice")], where
    if a:
        assert effects[0]["key"] == "effect/%s-damage" % a["damage_type"] and effects[0]["dice"] == a["damage"], where
    for x in effects:
        assert x["text"].rstrip(".") in text, (where, x["key"])
        if x.get("dice"):
            assert dice(x["dice"]) in x["text"], (where, x["key"])
        if e["name"] == "Multiattack" and x["key"] in ATTACK_KEYS:
            continue  # checked against the named attacks below
        for key in [x["key"], *x.get("requires", [])]:
            check_key(where, key, x["text"])
    if e["name"] == "Multiattack":
        named = {s["name"] for s in m["actions"] if s["name"] != "Multiattack" and s["name"] in text}
        keyed = {s["name"] for s in m["actions"] if s["name"] != "Multiattack" and s["key"] in {x["key"] for x in effects}}
        recharge = {s["name"] for s in m["actions"] if s.get("uses")}
        assert named and named == keyed - recharge, (where, named, keyed)


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_monsters_155.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/monsters_155.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-82-world-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 155]
    records = catalogue["monsters"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["source"]["license"] == "CC-BY-4.0"
    assert len(records) == 12 and [(m["id"], m["name"], m["source_pages"][0]) for m in records] == expected
    assert [m["name"] for m in records] == list(BLOCKS)
    entries = 0
    for m in records:
        first, last, following = BLOCKS[m["name"]]
        assert m["source_pages"] == [first, last], m["name"]
        lines = column_lines(pdf, first) + column_lines(pdf, first + 1)
        start = next(i for i in range(len(lines) - 1)
                     if lines[i][1] == m["name"] and re.match(r"(%s) " % SIZES, lines[i + 1][1]))
        end = next(i for i in range(start + 1, len(lines)) if lines[i][1] == following)
        block = lines[start + 1:end]
        assert sorted({p for p, _ in block}) == list(range(first, last + 1)), (m["name"], "page range")
        source = rejoin([line for _, line in block])
        assert render(m) == source, (m["name"], "block text differs from PDF", render(m), source)
        # Entry headings, from font runs rather than the text layer.
        runs = [(page,) + r for page in (first, first + 1) for r in heading_runs(pdf, page)]
        begin = next(i for i in range(len(runs) - 1)
                     if runs[i][4] == m["name"] and re.match(r"(%s) " % SIZES, runs[i + 1][4]))
        stop = next(i for i in range(begin + 1, len(runs)) if runs[i][4] == following)
        heads = [r[5] for r in runs[begin:stop] if r[5]]
        ours = m["traits"] + sections(m)
        assert heads == [heading(e) for e in ours], (m["name"], heads)
        entries += len(ours)
        for t in m["traits"]:
            assert t["key"] == "trait/" + t["name"].lower().replace(" ", "-"), t["name"]
            check_entry(m, t, None)
        for a in sections(m):
            check_entry(m, a, a["effects"])
        used = {"monster/stat-block"}
        for e in m["traits"]:
            used |= {e["key"], *e.get("requires", [])}
        for a in sections(m):
            used.add(a["key"])
            for x in a["effects"]:
                used |= {x["key"], *x.get("requires", [])}
        header = {"monster/skills": "skills", "defense/damage-resistance": "resistances",
                  "defense/damage-immunity": "immunities", "defense/condition-immunity": "condition_immunities",
                  "action/bonus-action": "bonus_actions", "action/reaction": "reactions", "language/telepathy": "telepathy"}
        used |= {key for key, field in header.items() if m.get(field)}
        used |= {"movement/" + k for k in m["speeds"] if k != "walk"} | {"sense/" + k for k in m["senses"]["special"]}
        used |= {"sense/%s-%s" % (k, re.sub(r"[^a-z0-9]+", "-", v.lower()).strip("-")) for k, v in m["senses"].get("notes", {}).items()}
        assert used == set(m["required_mechanics"]), (m["name"], used ^ set(m["required_mechanics"]))
    print("#155 pinned PDF: twelve stat blocks, %d entries, headers, page ranges, typed values and keys match" % entries)


if __name__ == "__main__":
    main()
