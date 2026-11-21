#!/usr/bin/env python3
"""Compare the issue #161 stat blocks with the pinned SRD 5.2.1 PDF.

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
    "Giant Crocodile": "Giant Eagle", "Giant Eagle": "Giant Elk", "Giant Elk": "Giant Fire Beetle",
    "Giant Fire Beetle": "Giant Frog", "Giant Frog": "Giant Goat", "Giant Goat": "Giant Hyena",
    "Giant Hyena": "Giant Lizard", "Giant Lizard": "Giant Octopus", "Giant Octopus": "Giant Owl",
    "Giant Owl": "Giant Rat", "Giant Rat": "Giant Scorpion", "Giant Scorpion": "Giant Seahorse",
}
# The Giant Rat's Bite is the one attack the source prints as "reach 5 feet."
FEET_REACH = {"Giant Rat"}
SIZES = "Tiny|Small|Medium|Large|Huge|Gargantuan"
ABILITY = {"str": "Strength", "dex": "Dexterity", "con": "Constitution",
           "int": "Intelligence", "wis": "Wisdom", "cha": "Charisma"}
SECTIONS = (("Traits", "traits"), ("Actions", "actions"), ("Bonus Actions", "bonus_actions"), ("Reactions", "reactions"))
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
        head = re.match(r"<i><b>(.*?)</b>", body)
        plain = html.unescape(re.sub(r"<[^>]+>", "", body)).strip()
        runs.append((left >= 445, top, left, plain, html.unescape(head[1]).strip().rstrip(".") if head else None))
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


def heading(e):
    uses = e.get("uses")
    if not uses:
        return e["name"]
    if uses["recharge"]:
        return "%s (Recharge %d–6)" % (e["name"], uses["recharge"])
    assert not uses["lair_per_day"], e["name"]
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
    for label, field in SECTIONS:
        if m.get(field):
            out.append(label)
            out += [e["text"] for e in m[field]]
    return " ".join(out)


def check_key(where, key, text):
    """A dependency key that names a damage type, condition, save or spell must be printed by its entry."""
    kind, _, name = key.partition("/")
    if kind == "effect" and name.endswith("-damage") and name[:-len("-damage")] in DAMAGE_TYPES:
        phrase = title(name[:-len("-damage")]) + " damage"
    elif kind == "condition":
        assert re.search(r"\b%s\b[^.]*\bconditions?\b" % title(name), text), (where, key)
        return
    elif kind == "save" and name in ABILITY:
        phrase = ABILITY[name] + " Saving Throw"
    elif kind == "spell":
        phrase = title(name)
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
    assert text.startswith(heading(e) + ". "), where
    expected = []
    a = e.get("attack")
    # No untyped attack roll may hide in an entry.
    assert ("Attack Roll" in text) == (e.get("key") in ("action/melee-attack", "action/ranged-attack")), where
    if a:
        if a["kind"] == "melee":
            reach = "reach %d %s" % (a["reach"], "feet." if m["name"] in FEET_REACH else "ft.")
        else:  # one printed range is retained as an equal long range
            reach = "range %d ft." % a["range"] if a["range"] == a["long_range"] else "range %d/%d ft." % (a["range"], a["long_range"])
        want = "%s Attack Roll: +%d, %s Hit: %s %s damage" % (title(a["kind"]), a["bonus"], reach, dice(a["damage"]), title(a["damage_type"]))
        assert want in text and (e["key"], effects[0]["key"]) == ("action/%s-attack" % a["kind"], "effect/%s-damage" % a["damage_type"]), where
        first = effects[0]
        if a["damage"]["dice"]:
            expected.append(dice(a["damage"]))
            assert first.get("dice") == a["damage"], where
        else:
            assert "dice" not in first and first.get("requires") == ["damage/fixed"], where
        effects = effects[1:]
        assert all(x["text"].rstrip(".,") in text for x in [first] + effects), where
    if e.get("save"):
        assert re.search(r"%s Saving Throw: DC %d[,.] " % (ABILITY[e["save"]["ability"]], e["save"]["dc"]), text), where
    if e.get("target"):
        t = e["target"]
        want = {
            "line": r"each creature in an? %d-foot-long, %d-foot-wide Line" % (t["size"], t["width"]),
            "cone": r"each creature in a %d-foot Cone" % t["size"],
            "emanation": r"each creature in a %d-foot Emanation originating from the \w+" % t["size"],
            "creature": r"one [^.]*creature (the \w+ can see within %d feet|Grappled by the \w+)" % t["range"] if t["range"]
                        else r"one [^.]*creature Grappled by the \w+",
        }[t["shape"]]
        assert re.search(want, text), where
    if e.get("spellcasting"):
        s = e["spellcasting"]
        bonus = ", +%d to hit with spell attacks" % s["attack_bonus"] if s["attack_bonus"] else ""
        assert bool(s["attack_bonus"]) == ("to hit with spell attacks" in text), where
        # A zero DC means the entry prints none (the Giant Owl).
        assert bool(s["save_dc"]) == ("spell save DC" in text), where
        printed_dc = " (spell save DC %d%s)" % (s["save_dc"], bonus) if s["save_dc"] else ""
        want = "using %s as the spellcasting ability%s: At Will: %s" % (
            ABILITY[s["ability"]], printed_dc, ", ".join(s["at_will"]))
        for group in s["per_day"]:
            want += " %d/Day%s: %s" % (group["uses"], " Each" if len(group["spells"]) > 1 else "", ", ".join(group["spells"]))
        assert text.endswith(want), where
        listed = [(x, "spellcasting/at-will") for x in s["at_will"]] + [(x, "limit/per-day") for g in s["per_day"] for x in g["spells"]]
        assert [(x["text"], x["requires"][0]) for x in effects if x["key"].startswith("spell/")] == listed, where
        others = [x["key"] for x in effects if not x["key"].startswith("spell/")]
        expected_others = [key for key, phrase in (("spellcasting/no-material-components", "requiring no Material components"),
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
        assert x["text"].rstrip(".,") in text, (where, x["key"])
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
        raise SystemExit("usage: verify_srd_monsters_161.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    data = ROOT / "internal/rules/srd521/testdata"
    catalogue = json.loads((data / "monsters_161.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-82-world-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 161]
    records = catalogue["monsters"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["source"]["license"] == "CC-BY-4.0"
    assert len(records) == 12 and [(m["id"], m["name"], m["source_pages"][0]) for m in records] == expected
    assert [m["name"] for m in records] == list(NEXT)
    entries = 0
    for m in records:
        first, last = m["source_pages"]
        assert last in (first, first + 1), m["name"]
        lines = column_lines(pdf, first) + column_lines(pdf, first + 1)
        start = next(i for i in range(len(lines) - 1)
                     if lines[i][1] == m["name"] and re.match(r"(%s) " % SIZES, lines[i + 1][1]))
        end = next(i for i in range(start + 1, len(lines)) if lines[i][1] == NEXT[m["name"]])
        block = lines[start + 1:end]
        assert sorted({p for p, _ in block}) == list(range(first, last + 1)), (m["name"], "page range")
        source = " ".join(line for _, line in block)
        assert squash(render(m)) == squash(source), (m["name"], "block text differs from PDF")
        # Entry headings, from font runs rather than the text layer.
        runs = [(page,) + r for page in (first, first + 1) for r in heading_runs(pdf, page)]
        begin = next(i for i in range(len(runs) - 1)
                     if runs[i][4] == m["name"] and re.match(r"(%s) " % SIZES, runs[i + 1][4]))
        stop = next(i for i in range(begin + 1, len(runs)) if runs[i][4] == NEXT[m["name"]])
        heads = [r[5] for r in runs[begin:stop] if r[5]]
        ours = [e for _, field in SECTIONS for e in m.get(field, [])]
        assert heads == [heading(e) for e in ours], (m["name"], heads)
        entries += len(ours)
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
        header = {"monster/skills": "skills", "defense/damage-resistance": "resistances", "defense/damage-immunity": "immunities",
                  "defense/condition-immunity": "condition_immunities", "movement/hover": "hover", "language/telepathy": "telepathy",
                  "action/bonus-action": "bonus_actions", "action/reaction": "reactions"}
        used |= {key for key, field in header.items() if m.get(field)}
        used |= {"movement/" + k for k in m["speeds"] if k != "walk"} | {"sense/" + k for k in m["senses"]["special"]}
        used |= {"sense/%s-%s" % (k, re.sub(r"[^a-z0-9]+", "-", v.lower()).strip("-")) for k, v in m["senses"].get("notes", {}).items()}
        used |= {"effect/%s-damage" % a["attack"]["damage_type"] for a in ours if a.get("attack")}
        assert used == set(m["required_mechanics"]), (m["name"], used ^ set(m["required_mechanics"]))
    print("#161 pinned PDF: twelve stat blocks, %d entries, headers, page ranges and typed values match" % entries)


if __name__ == "__main__":
    main()
