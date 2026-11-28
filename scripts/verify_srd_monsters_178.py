#!/usr/bin/env python3
"""Compare the issue #178 stat blocks with the pinned SRD 5.2.1 PDF.

Each record is printed back in the source's stat-block notation from its typed
fields and retained entry text, then compared with the PDF block character by
character (ignoring whitespace and line-break hyphens). Entry headings are read
separately from the PDF's bold-italic runs, so a missed or merged entry fails.
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
# Record name -> heading that follows its block in the PDF.
NEXT = {
    "Young Copper Dragon": "Adult Copper Dragon", "Young Gold Dragon": "Adult Gold Dragon",
    "Young Green Dragon": "Adult Green Dragon", "Young Red Dragon": "Adult Red Dragon",
    "Young Silver Dragon": "Adult Silver Dragon", "Young White Dragon": "Adult White Dragon",
    "Zombie": "Ogre Zombie",
}
LAST_PAGE = 364
SIZES = "Tiny|Small|Medium|Large|Huge|Gargantuan"
ABILITY = {"str": "Strength", "dex": "Dexterity", "con": "Constitution",
           "int": "Intelligence", "wis": "Wisdom", "cha": "Charisma"}
SECTIONS = (("Traits", "traits"), ("Actions", "actions"), ("Bonus Actions", "bonus_actions"), ("Reactions", "reactions"),
            ("Legendary Actions", "legendary_actions"))
# No heading in this batch carries a form qualifier.
QUALIFIER = {}
# Records that print "(450 XP; PB +2)" instead of "(XP 450; PB +2)".
XP_BEFORE_LABEL = {"Young White Dragon"}
# Printed skill bonuses that no ability modifier plus PB or twice PB explains.
PRINTED_SKILLS = {}
# Saves printed unlike their typed value. The Young White Dragon's Intelligence save
# prints "2" beside a −2 modifier (p. 340); the record keeps the derived −2.
PRINTED_SAVES = {("Young White Dragon", "int"): "2"}
# Every melee-or-ranged entry in this batch prints "ft." after its reach.
NO_REACH_PERIOD = set()
DAMAGE_TYPES = "acid bludgeoning cold fire force lightning necrotic piercing poison psychic radiant slashing thunder".split()
PRINTED_DICE = r"\d+ \(\d+d\d+(?: [+−] \d+)?\)"


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


def heading_runs(pdf, page):
    """Text runs in reading order with their bold-italic entry heading, if any."""
    runs = []
    for m in re.finditer(r'<text top="(\d+)" left="(\d+)"[^>]*>(.*?)</text>', run(
            "pdftohtml", "-xml", "-i", "-f", str(page), "-l", str(page), "-stdout", str(pdf))):
        top, left, body = int(m[1]), int(m[2]), m[3]
        # A long heading can be set as two adjacent bold-italic spans (the Medusa's, p. 306).
        head = re.match(r"(?:<i><b>.*?</b></i> ?)+|<i><b>.*?</b>", body)
        bold = " ".join(part.strip() for part in re.findall(r"<i><b>(.*?)</b>", head[0])) if head else None
        plain = html.unescape(re.sub(r"<[^>]+>", "", body)).strip()
        runs.append((left >= 445, top, left, plain, html.unescape(bold).strip().rstrip(".") if head else None))
    return sorted(runs)


def squash(text):
    return re.sub(r"[\s\-­]", "", text)


def signed(n):
    return ("−%d" % -n) if n < 0 else "+%d" % n


def dice(d):
    """Rolled values print as "N (XdY ± Z)"; a fixed hit prints its number alone."""
    if d["dice"] == 0:
        assert d["sides"] == 0 and d["modifier"] == 0, d
        return str(d["average"])
    mod = " %s %d" % ("+" if d["modifier"] > 0 else "−", abs(d["modifier"])) if d["modifier"] else ""
    return "%d (%dd%d%s)" % (d["average"], d["dice"], d["sides"], mod)


def title(s):
    return " ".join(w.capitalize() for w in s.replace("-", " ").split())


def heading(record, e):
    uses, notes = e.get("uses"), []
    if uses:
        if uses["recharge"]:  # "Recharge 6" is printed alone, never as a range
            notes.append("Recharge 6" if uses["recharge"] == 6 else "Recharge %d–6" % uses["recharge"])
        elif uses["lair_per_day"]:
            notes.append("%d/Day, or %d/Day in Lair" % (uses["per_day"], uses["lair_per_day"]))
        else:
            notes.append("%d/Day" % uses["per_day"])
    if (record, e["name"]) in QUALIFIER:
        notes.append(QUALIFIER[record, e["name"]])
    return e["name"] + (" (%s)" % "; ".join(notes) if notes else "")


def render(m):
    """Print one record in source order from typed fields and entry text."""
    kind = title(m["type"])
    if "swarm_of" in m:
        kind = "Swarm of %s %s" % (title(m["swarm_of"]), {"beast": "Beasts", "undead": "Undead"}[m["type"]])
    kind += " (%s)" % ", ".join(m["tags"]) if m.get("tags") else ""
    speed = "Speed %d ft." % m["speeds"]["walk"] + "".join(
        ", %s %d ft.%s" % (title(k), v, " (hover)" if k == "fly" and m.get("hover") else "")
        for k, v in sorted(m["speeds"].items()) if k != "walk")
    for note in ("speed_note", "conditional_speed_note"):
        if note in m:
            speed += ", " + m[note]
    out = [
        "%s %s, %s" % (" or ".join(title(s) for s in m["sizes"]), kind, m["alignment"]),
        "AC %d Initiative %s (%d)" % (m["ac"], signed(m["initiative"]["modifier"]), m["initiative"]["score"]),
        "HP " + dice(m["hp"]), speed,
    ]
    for key in ("str", "dex", "con", "int", "wis", "cha"):
        a = m["abilities"][key]
        save = signed(a["save"])
        if (m["name"], key) in PRINTED_SAVES:
            assert a["save"] == a["modifier"], (m["name"], key, "the record must keep the derived save")
            save = PRINTED_SAVES[m["name"], key]
        out.append("%s %d %s %s" % (title(key), a["score"], signed(a["modifier"]), save))
    if "skills" in m:
        out.append("Skills " + ", ".join(
            "%s %s" % (title(k).replace(" Of ", " of "), signed(v)) for k, v in sorted(m["skills"].items())))
    for field, note, label in (("vulnerabilities", "vulnerability_note", "Vulnerabilities"),
                               ("resistances", "resistance_note", "Resistances")):
        if field in m:
            out.append(label + " " + ", ".join(title(v) for v in m[field]))
        if note in m:
            out.append(label + " " + m[note])
    immune = [", ".join(title(v) for v in m[k]) for k in ("immunities", "condition_immunities") if k in m]
    if "condition_immunity_note" in m:
        immune.append(m["condition_immunity_note"])
    if immune:
        out.append("Immunities " + "; ".join(immune))
    if m["gear"]:
        out.append("Gear " + ", ".join(m["gear"]))
    notes = m["senses"].get("notes", {})
    special = ", ".join("%s %d ft.%s" % (title(k), v, " (%s)" % notes[k] if k in notes else "")
                        for k, v in sorted(m["senses"]["special"].items()))
    out.append("Senses " + (special + "; " if special else "") + "Passive Perception %d" % m["senses"]["passive_perception"])
    languages = ", ".join(m["languages"]) or "None"
    if m.get("telepathy"):
        languages += "; telepathy %d ft." % m["telepathy"]
    if m.get("telepathy_note"):
        languages += " (%s)" % m["telepathy_note"]
    out.append("Languages " + languages)
    lair = ", or {:,} in lair".format(m["lair_xp"]) if m.get("lair_xp") else ""
    xp = ("{:,} XP" if m["name"] in XP_BEFORE_LABEL else "XP {:,}").format(m["xp"])
    out.append("CR %s (%s%s; PB +%d)" % (m["cr"], xp, lair, m["proficiency_bonus"]))
    for label, field in SECTIONS:
        if m.get(field):
            out.append(label)
            if field == "legendary_actions":
                out.append(m["legendary"]["text"])
            out += [e["text"] for e in m[field]]
    return " ".join(out)


def check_key(where, key, text):
    """A dependency key that names a damage type, condition, save or spell must be printed by its entry."""
    kind, _, name = key.partition("/")
    if kind == "effect" and name.endswith("-damage") and name[:-len("-damage")] in DAMAGE_TYPES:
        phrase = title(name[:-len("-damage")]) + " damage"
    elif kind == "condition":  # also printed in pairs, as in "the Blinded and Poisoned conditions"
        assert re.search(r"\b%s( and \w+)? conditions?\b" % title(name), text), (where, key)
        return
    elif kind == "save" and name in ABILITY:
        phrase = ABILITY[name] + " Saving Throw"
    elif kind == "spell":  # keys drop apostrophes, as in "Heroes' Feast"
        phrase, text = title(name), text.replace("’", "")
    elif key == "limit/recharge":
        phrase = "(Recharge"
    elif key == "save/half-damage":
        phrase = "Half damage"
    elif kind == "area" and name in ("cone", "emanation"):
        phrase = "-foot " + title(name)
    else:
        return  # structural keys are tied to their entry's own quoted text by the Go tests
    assert phrase.lower() in text.lower(), (where, key, phrase)


def check_entry(m, e, effects):
    """Each typed payload must be printed by the entry's own source text."""
    text, where = e["text"], (m["name"], e["name"])
    assert text.startswith(heading(m["name"], e) + ". "), where
    if (m["name"], e["name"]) in QUALIFIER:
        assert QUALIFIER[where] in [x["text"] for x in effects], where
    expected = []
    a = e.get("attack")
    both = e.get("key") == "action/melee-or-ranged-attack"
    # No untyped attack roll may hide in an entry except a melee-or-ranged attack.
    # The Roper's Tentacle (p. 320) rolls to hit but its hit prints no damage.
    harmless = e.get("key") == "action/melee-attack-without-damage"
    assert ("Attack Roll" in text) == (harmless or e.get("key") in ("action/melee-attack", "action/ranged-attack", "action/melee-or-ranged-attack")), where
    if harmless:
        assert a is None and re.search(r"Melee Attack Roll: \+\d+, reach \d+ ft\. Hit: The target has ", text), where
        assert effects[0]["key"] == "attack/roll-without-damage" and not any("dice" in x for x in effects), where
    if both:
        d = effects[0]["dice"]
        period = "" if where in NO_REACH_PERIOD else r"\."
        assert a is None and re.search(r"Melee or Ranged Attack Roll: \+\d+, reach \d+ ft%s or range \d+(/\d+)? ft\. Hit: %s %s damage" % (
            period, re.escape(dice(d)), title(effects[0]["key"][len("effect/"):-len("-damage")])), text), where
    if a:
        if a["kind"] == "melee":
            reach = "reach %d ft." % a["reach"]
        else:  # one printed range is retained as an equal long range
            reach = "range %d ft." % a["range"] if a["range"] == a["long_range"] else "range %d/%d ft." % (a["range"], a["long_range"])
        # A printed Advantage condition on the roll must be retained as its own effect.
        advantage = "".join(" " + x["text"] for x in effects if x["key"].startswith("attack/advantage-against-"))
        want = "%s Attack Roll: +%d%s, %s Hit: %s %s damage" % (
            title(a["kind"]), a["bonus"], advantage, reach, dice(a["damage"]), title(a["damage_type"]))
        assert want in text and (e["key"], effects[0]["key"]) == ("action/%s-attack" % a["kind"], "effect/%s-damage" % a["damage_type"]), where
        first = effects[0]
        if a["damage"]["dice"]:
            expected.append(dice(a["damage"]))
            assert first.get("dice") == a["damage"], where
        else:
            assert "dice" not in first and first.get("requires") == ["damage/fixed"], where
        effects = effects[1:]
        assert all(x["text"].rstrip(".,") in e["text"] for x in [first] + effects), where
    if e.get("save"):
        assert re.search(r"%s Saving Throw: DC %d[,.] " % (ABILITY[e["save"]["ability"]], e["save"]["dc"]), text), where
    if e.get("target"):
        t = e["target"]
        want = {
            "cone": r"each creature( that isn’t currently affected by this breath)? in a %d-foot Cone" % t["size"],
            "line": r"each creature in a %d-foot-long, %d-foot-wide Line" % (t["size"], t["width"]),
            "emanation": r"(?<!\d)%d-foot Emanation( originating from the \w+|(?= while ))" % t["size"],
            "sphere": r"%d-foot-radius Sphere centered on a point [^.]*within %d feet" % (t["size"], t["range"]),
            "creature": r"one [^.]*creature [^.]*within %d feet" % t["range"] if t["range"]
                        else r"one [^.]*creature (Grappled by the \w+|in the \w+’s space)",
        }[t["shape"]]
        assert re.search(want, text), where
    if e.get("spellcasting"):
        s = e["spellcasting"]
        bonus = ", +%d to hit with spell attacks" % s["attack_bonus"] if s["attack_bonus"] else ""
        assert bool(s["attack_bonus"]) == ("to hit with spell attacks" in text), where
        # A zero DC means the entry prints none (the Priest Acolyte, p. 316).
        assert bool(s["save_dc"]) == ("spell save DC" in text), where
        printed_dc = " (spell save DC %d%s)" % (s["save_dc"], bonus) if s["save_dc"] else ""
        want = "using %s as the spellcasting ability%s:" % (ABILITY[s["ability"]], printed_dc)
        if s["at_will"]:
            want += " At Will: " + ", ".join(s["at_will"])
        for group in s["per_day"]:
            want += " %d/Day%s: %s" % (group["uses"], " Each" if len(group["spells"]) > 1 else "", ", ".join(group["spells"]))
        assert text.endswith(want), (where, want)
        listed = [(x, "spellcasting/at-will") for x in s["at_will"]] + [(x, "limit/per-day") for g in s["per_day"] for x in g["spells"]]
        assert [(x["text"], x["requires"][0]) for x in effects if x["key"].startswith("spell/")] == listed, where
        others = [x["key"] for x in effects if not x["key"].startswith("spell/")]
        expected_others = [key for key, phrase in (("spellcasting/no-material-components", "requiring no Material components"),
                                                   ("spellcasting/no-somatic-or-material-components", "requiring no Somatic or Material components"),
                                                   ("spellcasting/no-components", "requiring no spell components"),
                                                   ("spellcasting/spell-attack-bonus", "to hit with spell attacks")) if phrase in text]
        assert others == expected_others, where
    printed = re.findall(PRINTED_DICE, text)
    if effects is None:  # a trait types its dice itself, in source order
        for key in e.get("requires", []):
            check_key(where, key, text)
        typed = ([e["dice"]] if e.get("dice") else []) + [d["dice"] for d in e.get("more_dice", [])]
        assert printed == [dice(d) for d in typed], where
        return
    # Every dice expression in an action is typed, in order, on the attack or an effect.
    assert printed == expected + [dice(x["dice"]) for x in effects if x.get("dice")], where
    for x in effects:
        assert x["text"].rstrip(".,") in e["text"], (where, x["key"])
        if x.get("dice"):
            assert dice(x["dice"]) in x["text"], (where, x["key"])
        if not e.get("spellcasting"):
            for key in [x["key"], *x.get("requires", [])]:
                check_key(where, key, text)
    if e["name"] == "Multiattack":
        named = [n for n in (a["name"] for a in m["actions"]) if n != "Multiattack" and n in text]
        assert named, where


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_monsters_178.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    data = ROOT / "internal/rules/srd521/testdata"
    catalogue = json.loads((data / "monsters_178.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-82-world-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 178]
    records = catalogue["monsters"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["source"]["license"] == "CC-BY-4.0"
    assert len(records) == 7 and [(m["id"], m["name"], m["source_pages"][0]) for m in records] == expected
    assert [m["name"] for m in records] == list(NEXT)
    entries = 0
    for m in records:
        for (record, skill), bonus in PRINTED_SKILLS.items():
            if record == m["name"]:
                assert m["skills"][skill] == bonus, record
        first, last = m["source_pages"]
        assert last in (first, first + 1), m["name"]
        pages = [p for p in (first, first + 1) if p <= LAST_PAGE]
        lines = [line for p in pages for line in column_lines(pdf, p)]
        start = next(i for i in range(len(lines) - 1)
                     if lines[i][1] == m["name"] and re.match(r"(%s) " % SIZES, lines[i + 1][1]))
        after = NEXT[m["name"]]
        end = len(lines) if after is None else next(i for i in range(start + 1, len(lines)) if lines[i][1] == after)
        block = lines[start + 1:end]
        assert sorted({p for p, _ in block}) == list(range(first, last + 1)), (m["name"], "page range")
        source = " ".join(line for _, line in block)
        assert squash(render(m)) == squash(source), (m["name"], "block text differs from PDF")
        # Entry headings, from font runs rather than the text layer.
        runs = [(page,) + r for page in pages for r in heading_runs(pdf, page)]
        begin = next(i for i in range(len(runs) - 1)
                     if runs[i][4] == m["name"] and re.match(r"(%s) " % SIZES, runs[i + 1][4]))
        stop = len(runs) if after is None else next(i for i in range(begin + 1, len(runs)) if runs[i][4] == after)
        heads = [r[5] for r in runs[begin:stop] if r[5]]
        ours = [e for _, field in SECTIONS for e in m.get(field, [])]
        assert heads == [heading(m["name"], e) for e in ours], (m["name"], heads)
        entries += len(ours)
        if m.get("legendary"):
            uses = m["legendary"]
            want = "Legendary Action Uses: %d%s. " % (uses["uses"], " (%d in Lair)" % uses["lair_uses"] if uses["lair_uses"] else "")
            assert uses["text"].startswith(want), m["name"]
        for t in m["traits"]:
            check_entry(m, t, None)
        for a in ours[len(m["traits"]):]:
            check_entry(m, a, a["effects"])
        used = {"monster/stat-block"}
        for e in m["traits"]:
            used |= {e["key"], *e.get("requires", []), *(d["key"] for d in e.get("more_dice", []))}
        for a in ours[len(m["traits"]):]:
            used.add(a["key"])
            for x in a["effects"]:
                used |= {x["key"], *x.get("requires", [])}
        header = {"monster/skills": "skills", "monster/swarm": "swarm_of", "movement/chosen-mode": "speed_note",
                  "defense/conditional-damage-vulnerability": "vulnerability_note",
                  "defense/conditional-condition-immunity": "condition_immunity_note",
                  "movement/conditional-speed": "conditional_speed_note",
                  "defense/damage-resistance": "resistances", "defense/damage-immunity": "immunities",
                  "defense/condition-immunity": "condition_immunities", "movement/hover": "hover", "language/telepathy": "telepathy",
                  "action/bonus-action": "bonus_actions", "action/reaction": "reactions",
                  "defense/damage-vulnerability": "vulnerabilities", "defense/chosen-damage-resistance": "resistance_note",
                  "language/telepathy-condition": "telepathy_note", "lair/xp": "lair_xp", "legendary/actions": "legendary"}
        used |= {key for key, field in header.items() if m.get(field)}
        if m.get("legendary", {}).get("lair_uses"):
            used.add("lair/uses")
        used |= {"movement/" + k for k in m["speeds"] if k != "walk"} | {"sense/" + k for k in m["senses"]["special"]}
        used |= {"sense/%s-%s" % (k, re.sub(r"[^a-z0-9]+", "-", v.lower()).strip("-")) for k, v in m["senses"].get("notes", {}).items()}
        used |= {"effect/%s-damage" % a["attack"]["damage_type"] for a in ours if a.get("attack")}
        assert used == set(m["required_mechanics"]), (m["name"], used ^ set(m["required_mechanics"]))
    print("#178 pinned PDF: seven stat blocks, %d entries, headers, page ranges and typed values match" % entries)


if __name__ == "__main__":
    main()
