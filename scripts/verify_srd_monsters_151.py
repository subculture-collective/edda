#!/usr/bin/env python3
"""Compare the issue #151 stat blocks with the pinned SRD 5.2.1 PDF.

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
    "Aboleth": "Air Elemental", "Adult Black Dragon": "Ancient Black Dragon",
    "Adult Blue Dragon": "Ancient Blue Dragon", "Adult Brass Dragon": "Ancient Brass Dragon",
    "Adult Bronze Dragon": "Ancient Bronze Dragon", "Adult Copper Dragon": "Ancient Copper Dragon",
    "Adult Gold Dragon": "Ancient Gold Dragon", "Adult Green Dragon": "Ancient Green Dragon",
    "Adult Red Dragon": "Ancient Red Dragon", "Adult Silver Dragon": "Ancient Silver Dragon",
    "Adult White Dragon": "Ancient White Dragon", "Air Elemental": "Animated Objects",
}
SIZES = "Tiny|Small|Medium|Large|Huge|Gargantuan"
ABILITY = {"str": "Strength", "dex": "Dexterity", "con": "Constitution",
           "int": "Intelligence", "wis": "Wisdom", "cha": "Charisma"}


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
    mod = " + %d" % d["modifier"] if d["modifier"] else ""
    return "%d (%dd%d%s)" % (d["average"], d["dice"], d["sides"], mod)


def title(s):
    return " ".join(w.capitalize() for w in s.replace("-", " ").split())


def heading(e):
    uses = e.get("uses")
    if not uses:
        return e["name"]
    if uses["recharge"]:
        return "%s (Recharge %d–6)" % (e["name"], uses["recharge"])
    if uses["lair_per_day"]:
        return "%s (%d/Day, or %d/Day in Lair)" % (e["name"], uses["per_day"], uses["lair_per_day"])
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
    if immune:
        out.append("Immunities " + "; ".join(immune))
    special = ", ".join("%s %d ft." % (title(k), v) for k, v in sorted(m["senses"]["special"].items()))
    out.append("Senses " + (special + "; " if special else "") + "Passive Perception %d" % m["senses"]["passive_perception"])
    out.append("Languages " + ", ".join(m["languages"]) + ("; telepathy %d ft." % m["telepathy"] if m.get("telepathy") else ""))
    lair = ", or {:,} in lair".format(m["lair_xp"]) if m.get("lair_xp") else ""
    out.append("CR %s (XP %s%s; PB +%d)" % (m["cr"], "{:,}".format(m["xp"]), lair, m["proficiency_bonus"]))
    out.append("Traits")
    out += [t["text"] for t in m["traits"]]
    out.append("Actions")
    out += [a["text"] for a in m["actions"]]
    if m.get("legendary"):
        out += ["Legendary Actions", m["legendary"]["text"]] + [a["text"] for a in m["legendary_actions"]]
    return " ".join(out)


def check_entry(m, e, effects):
    """Each typed payload must be printed by the entry's own source text."""
    text, where = e["text"], (m["name"], e["name"])
    assert text.startswith(heading(e) + ". "), where
    if e.get("attack"):
        a = e["attack"]
        want = "Melee Attack Roll: +%d, reach %d ft. Hit: %s %s damage" % (
            a["bonus"], a["reach"], dice(a["damage"]), title(a["damage_type"]))
        assert a["kind"] == "melee" and want in text, where
    if e.get("save"):
        assert "%s Saving Throw: DC %d, " % (ABILITY[e["save"]["ability"]], e["save"]["dc"]) in text, where
    if e.get("target"):
        t = e["target"]
        want = {
            "line": r"each creature in an? %d-foot-long, %d-foot-wide Line" % (t["size"], t["width"]),
            "cone": r"each creature (that [^.]* )?in a %d-foot Cone" % t["size"],
            "sphere": r"each creature in a %d-foot-radius Sphere centered on a point [^.]* within %d feet" % (t["size"], t["range"]),
            "emanation": r"each creature in a %d-foot Emanation" % t["size"],
            "creature": r"one creature (the \w+ can see )?within %d feet" % t["range"] if t["range"]
                        else r"one [^.]*creature in the \w+’s space",
        }[t["shape"]]
        assert re.search(want, text), where
    if e.get("spellcasting"):
        s = e["spellcasting"]
        bonus = ", +%d to hit with spell attacks" % s["attack_bonus"] if s["attack_bonus"] else ""
        want = "using %s as the spellcasting ability (spell save DC %d%s): At Will: %s" % (
            ABILITY[s["ability"]], s["save_dc"], bonus, ", ".join(s["at_will"]))
        for group in s["per_day"]:
            want += " %d/Day%s: %s" % (group["uses"], " Each" if len(group["spells"]) > 1 else "", ", ".join(group["spells"]))
        assert text.endswith(want), where
    printed = re.findall(r"\d+ \(\d+d\d+(?: \+ \d+)?\)", text)
    if effects is None:  # a trait types its one dice expression itself
        assert printed == ([dice(e["dice"])] if e.get("dice") else []), where
        return
    # Every dice expression in an action is typed, in order, on one of its effects.
    assert printed == [dice(x["dice"]) for x in effects if x.get("dice")], where
    for x in effects:
        assert x["text"].rstrip(".") in text, (where, x["key"])
        if x.get("dice"):
            assert dice(x["dice"]) in x["text"], (where, x["key"])


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_monsters_151.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/monsters_151.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-82-world-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 151]
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
        source = re.sub(r"Con(\d)", r"Con \1", " ".join(line for _, line in block))  # p. 340 prints "Con22"
        assert squash(render(m)) == squash(source), (m["name"], "block text differs from PDF")
        # Entry headings, from font runs rather than the text layer.
        runs = [(page,) + r for page in (first, first + 1) for r in heading_runs(pdf, page)]
        starts = [i for i, r in enumerate(runs) if r[4] == m["name"]]
        begin = starts[-1] if m["name"] != "Aboleth" else starts[0]
        stop = next(i for i in range(begin + 1, len(runs)) if runs[i][4] == NEXT[m["name"]])
        heads = [r[5] for r in runs[begin:stop] if r[5]]
        ours = m["traits"] + m["actions"] + m.get("legendary_actions", [])
        assert heads == [heading(e) for e in ours], (m["name"], heads)
        entries += len(ours)
        for t in m["traits"]:
            check_entry(m, t, None)
        for a in m["actions"] + m.get("legendary_actions", []):
            check_entry(m, a, a["effects"])
        if m.get("legendary"):
            want = "Legendary Action Uses: %d (%d in Lair). " % (m["legendary"]["uses"], m["legendary"]["lair_uses"])
            assert m["legendary"]["text"].startswith(want), m["name"]
        used = {"monster/stat-block"}
        for e in m["traits"]:
            used |= {e["key"], *e.get("requires", [])}
        for a in m["actions"] + m.get("legendary_actions", []):
            used.add(a["key"])
            for x in a["effects"]:
                used |= {x["key"], *x.get("requires", [])}
        assert used <= set(m["required_mechanics"]), (m["name"], used - set(m["required_mechanics"]))
    print("#151 pinned PDF: twelve stat blocks, %d entries, headers, page ranges and typed values match" % entries)


if __name__ == "__main__":
    main()
