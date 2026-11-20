#!/usr/bin/env python3
"""Compare the issue #153 stat blocks with the pinned SRD 5.2.1 PDF.

Each record is printed back in the source's stat-block notation from its typed
fields and retained entry text, then compared with the PDF block word for word
(ignoring only line breaks and line-break hyphens). Entry headings are read
separately from the PDF's bold-italic runs, so a missed or merged entry fails.
Every typed value and every dependency key is tied to a phrase its own entry
prints.
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
    "Animated Armor": (259, 259, "Animated Flying Sword"),
    "Animated Flying Sword": (259, 259, "Animated Rug of Smothering"),
    "Animated Rug of Smothering": (259, 259, "Ankheg"),
    "Ankheg": (259, 259, "Assassin"),
    "Ankylosaurus": (344, 344, "Ape"),
    "Ape": (344, 344, "Archelon"),
    "Archelon": (344, 345, "Baboon"),
    "Archmage": (305, 305, "Magmin"),
    "Assassin": (260, 260, "Awakened Plants"),
    "Awakened Shrub": (260, 260, "Awakened Tree"),
    "Awakened Tree": (260, 260, "Axe Beak"),
    "Axe Beak": (260, 260, "Azer"),
}
SIZES = "Tiny|Small|Medium|Large|Huge|Gargantuan"
ABILITY = {"str": "Strength", "dex": "Dexterity", "con": "Constitution",
           "int": "Intelligence", "wis": "Wisdom", "cha": "Charisma"}
# Dependency key -> phrase printed by the entry (or effect) that names the key.
# Damage, spell and trait keys are derived from their own names below.
PHRASES = {
    "requirement/target-size": "or smaller creature",
    "condition/grappled": "Grappled condition (escape DC 13)",
    "effect/condition-instead-of-damage": "instead of dealing damage",
    "condition/blinded": "the target has the Blinded and Restrained conditions",
    "condition/restrained": "the target has the Blinded and Restrained conditions",
    "effect/suffocating": "is suffocating",
    "effect/start-of-turn-damage": "Bludgeoning damage at the start of each of its turns",
    "limit/one-grappled-creature": "The rug can smother only one creature at a time",
    "requirement/not-grappling": "While grappling the target, the rug can’t take this action",
    "effect/halve-damage-taken": "the rug halves the damage it takes (round down)",
    "effect/share-damage-taken": "the target takes the same amount of damage",
    "attack/advantage-against-grappled-target": "(with Advantage if the target is Grappled by the ankheg)",
    "save/dex": "Dexterity saving throw",
    "area/line": "each creature in a 30-foot-long, 5-foot-wide Line",
    "limit/recharge": "(Recharge 6)",
    "save/half-damage": "Success: Half damage",
    "movement/burrow-through-rock": "can burrow through solid rock",
    "movement/half-burrow-speed": "at half its Burrow Speed",
    "effect/leave-tunnel": "leaves a 10-foot-diameter tunnel in its wake",
    "condition/prone": "it has the Prone condition",
    "save/advantage-against-magic": "Advantage on saving throws against spells and other magical effects",
    "multiattack/any-combination": "in any combination",
    "spellcasting/higher-level-version": "version)",
    "spellcasting/included-in-ac": "(included in AC)",
    "spellcasting/cast-before-combat": "(cast before combat)",
    "defense/immunity-with-mind-blank": "Mind Blank",
    "action/spellcasting": "using the same spellcasting ability as Spellcasting",
    "reaction/spell-trigger": "in response to the spell’s trigger",
    "save/no-damage-on-success": "takes no damage if it succeeds on the save",
    "save/half-damage-on-failure": "only half damage if it fails",
    "requirement/not-incapacitated": "It can’t use this trait if it has the Incapacitated condition",
    "condition/poisoned": "the target has the Poisoned condition until the start of the assassin’s next turn",
    "action/dash": "Dash",
    "action/disengage": "Disengage",
    "action/hide": "Hide action",
}
ATTACK_KEYS = ("action/melee-attack", "action/ranged-attack", "action/melee-or-ranged-attack",
               "action/fixed-damage-melee-attack")


def run(*args):
    return subprocess.run(args, check=True, capture_output=True, text=True).stdout


def column_lines(pdf, page):
    """Read columns independently so continuation lines do not interleave."""
    out = []
    for x in (0, 306):
        text = run("pdftotext", "-f", str(page), "-l", str(page), "-x", str(x), "-y", "0",
                   "-W", "306", "-H", "792", "-layout", str(pdf), "-")
        for line in text.splitlines():
            line = " ".join(line.replace("\x0c", "").split())
            if line and not re.fullmatch(r"\d+ System Reference Document 5\.2\.1|(MOD SAVE ?)+", line):
                out.append((page, line))
    return out


def rejoin(lines):
    """Join wrapped lines; a line-break hyphen is dropped, a printed hyphen kept."""
    text = ""
    for line in lines:
        if text.endswith("-") and not re.search(r"\d-$|foot-$", text):
            text = text[:-1] + line
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
    mod = " + %d" % d["modifier"] if d["modifier"] else ""
    return "%d (%dd%d%s)" % (d["average"], d["dice"], d["sides"], mod)


def title(s):
    return " ".join(w.capitalize() for w in s.replace("-", " ").split())


def heading(e):
    uses = e.get("uses")
    if not uses:
        return e["name"]
    assert uses["lair_per_day"] == 0, e["name"]
    if uses["recharge"]:
        assert uses["recharge"] == 6 and uses["per_day"] == 0, e["name"]  # this batch prints only "Recharge 6"
        return "%s (Recharge 6)" % e["name"]
    return "%s (%d/Day)" % (e["name"], uses["per_day"])


def render(m):
    """Print one record in source order from typed fields and entry text."""
    kind = title(m["type"]) + (" (%s)" % ", ".join(m["tags"]) if m.get("tags") else "")
    speed = "Speed %d ft." % m["speeds"]["walk"] + "".join(
        ", %s %d ft." % (title(k), v) for k, v in sorted(m["speeds"].items()) if k != "walk")
    out = [
        "%s %s, %s" % (" or ".join(title(s) for s in m["sizes"]), kind, m["alignment"]),
        "AC %d Initiative %s (%d)" % (m["ac"], signed(m["initiative"]["modifier"]), m["initiative"]["score"]),
        "HP " + dice(m["hp"]), speed + (" (hover)" if m.get("hover") else ""),
    ]
    for key in ("str", "dex", "con", "int", "wis", "cha"):
        a = m["abilities"][key]
        out.append("%s %d %s %s" % (title(key), a["score"], signed(a["modifier"]), signed(a["save"])))
    if "skills" in m:
        out.append("Skills " + ", ".join("%s %s" % (title(k), signed(v)) for k, v in sorted(m["skills"].items())))
    for label in ("vulnerabilities", "resistances"):
        if label in m:
            out.append(title(label) + " " + ", ".join(title(v) for v in m[label]))
    immune = [", ".join(title(v) for v in m[k]) for k in ("immunities", "condition_immunities") if k in m]
    if "defense/immunity-with-mind-blank" in m["required_mechanics"]:
        immune[-1] += " (with Mind Blank)"
    if immune:
        out.append("Immunities " + "; ".join(immune))
    if m["gear"]:
        out.append("Gear " + ", ".join(m["gear"]))
    special = ", ".join("%s %d ft." % (title(k), v) for k, v in sorted(m["senses"]["special"].items()))
    out.append("Senses " + (special + "; " if special else "") + "Passive Perception %d" % m["senses"]["passive_perception"])
    out.append("Languages " + (", ".join(m["languages"]) or "None"))
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
    elif kind == "spell":
        phrase = title(name).replace(" Of ", " of ")
    elif key == "limit/per-day":
        phrase = "/Day"
    elif key == "spellcasting/at-will":
        return  # checked against the printed At Will list
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
    elif e.get("key") == "action/ranged-attack":
        want = "Ranged Attack Roll: +%d, range %d/%d ft. Hit: %s %s damage" % (
            a["bonus"], a["range"], a["long_range"], dice(a["damage"]), title(a["damage_type"]))
        assert a["kind"] == "ranged" and want in text, where
    else:
        assert a is None, where
        # No untyped attack roll may hide in an entry except the two reviewed below.
        attack_roll = "Attack Roll" in text
        assert attack_roll == (e.get("key") in ("action/fixed-damage-melee-attack", "action/melee-or-ranged-attack")), where
    if e.get("key") == "action/fixed-damage-melee-attack":
        assert text == "Rake. Melee Attack Roll: +1, reach 5 ft. Hit: 1 Slashing damage." and not re.search(r"\dd\d", text), where
    if e.get("key") == "action/melee-or-ranged-attack" and effects is not None and e["name"] != "Multiattack":
        assert text == "Arcane Burst. Melee or Ranged Attack Roll: +9, reach 5 ft. or range 150 ft. Hit: %s Force damage." % dice(effects[0]["dice"]), where
    if e.get("save"):
        assert e["key"] == "action/saving-throw", where
        assert "%s Saving Throw: DC %d, " % (ABILITY[e["save"]["ability"]], e["save"]["dc"]) in text, where
    if e.get("target"):
        t = e["target"]
        assert t["shape"] == "line" and t["range"] == 0, where
        assert "each creature in a %d-foot-long, %d-foot-wide Line" % (t["size"], t["width"]) in text, where
    if e.get("spellcasting"):
        s = e["spellcasting"]
        assert s["attack_bonus"] == 0 and "to hit with spell attacks" not in text, where
        want = "using %s as the spellcasting ability (spell save DC %d): At Will: %s" % (
            ABILITY[s["ability"]], s["save_dc"], ", ".join(s["at_will"]))
        for group in s["per_day"]:
            want += " %d/Day%s: %s" % (group["uses"], " Each" if len(group["spells"]) > 1 else "", ", ".join(group["spells"]))
        assert text.endswith(want), where
        listed = [(x, "spellcasting/at-will") for x in s["at_will"]] + [(x, "limit/per-day") for g in s["per_day"] for x in g["spells"]]
        assert [(x["text"], x["requires"][0]) for x in effects] == listed, where
    printed = re.findall(r"\d+ \(\d+d\d+(?: \+ \d+)?\)", text)
    if effects is None:  # a trait types its one dice expression itself
        assert printed == ([dice(e["dice"])] if e.get("dice") else []), where
        for key in e.get("requires", []):
            check_key(where, key, text)
        return
    # Every dice expression in an action is typed, in order, on one of its effects.
    assert printed == [dice(x["dice"]) for x in effects if x.get("dice")], where
    for x in effects:
        assert x["text"].rstrip(".") in text, (where, x["key"])
        if x.get("dice"):
            assert dice(x["dice"]) in x["text"], (where, x["key"])
        if e["name"] == "Multiattack" and x["key"] in ATTACK_KEYS:
            continue  # checked against the named attacks below
        # A listed spell's first requirement is its At Will or N/Day group, checked above.
        for key in [x["key"], *x.get("requires", [])[1 if e.get("spellcasting") else 0:]]:
            check_key(where, key, x["text"])
    if e["name"] == "Multiattack":
        named = {a["name"] for a in m["actions"] if a["name"] != "Multiattack" and a["name"] in text}
        keyed = {a["name"] for a in m["actions"] if a["name"] != "Multiattack" and a["key"] in {x["key"] for x in effects}}
        recharge = {a["name"] for a in m["actions"] if a.get("uses")}
        assert named and named == keyed - recharge, (where, named, keyed)


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_monsters_153.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/monsters_153.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-82-world-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 153]
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
        begin = [i for i, r in enumerate(runs) if r[4] == m["name"]][-1]
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
        if "defense/immunity-with-mind-blank" in m["required_mechanics"]:
            assert "Immunities Psychic; Charmed (with Mind Blank)" in source, m["name"]
        used = {"monster/stat-block"}
        for e in m["traits"]:
            used |= {e["key"], *e.get("requires", [])}
        for a in sections(m):
            used.add(a["key"])
            for x in a["effects"]:
                used |= {x["key"], *x.get("requires", [])}
        header = {"monster/skills": "skills", "defense/damage-vulnerability": "vulnerabilities",
                  "defense/damage-resistance": "resistances", "defense/damage-immunity": "immunities",
                  "defense/condition-immunity": "condition_immunities", "movement/hover": "hover",
                  "action/bonus-action": "bonus_actions", "action/reaction": "reactions"}
        used |= {key for key, field in header.items() if m.get(field)}
        used |= {"movement/" + k for k in m["speeds"] if k != "walk"} | {"sense/" + k for k in m["senses"]["special"]}
        assert used == set(m["required_mechanics"]), (m["name"], used ^ set(m["required_mechanics"]))
    # The Archmage's printed XP differs from the CR table on p. 256; both are source text.
    archmage = next(m for m in records if m["name"] == "Archmage")
    table = " ".join(run("pdftotext", "-f", "256", "-l", "256", "-layout", str(pdf), "-").split())
    assert (archmage["cr"], archmage["xp"]) == ("12", 8000) and "12 8,400" in table
    print("#153 pinned PDF: twelve stat blocks, %d entries, headers, page ranges, typed values and keys match" % entries)


if __name__ == "__main__":
    main()
