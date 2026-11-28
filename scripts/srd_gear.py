"""Shared comparison of adventuring gear records (issues #256-#262) with the
pinned SRD 5.2.1 PDF. Each verify_srd_items_<issue>.py script calls verify().

Built from the printed text, without a hand-kept value list: the Adventuring
Gear table (weights and prices), each entry's heading, page span, actions,
saving throws, ability checks, damage, pack contents and spellcasting classes,
and the variant tables (rows, weights, prices, ammunition amounts, storage and
bracketed row notes).

Hand-written and passed in as evidence: purposes of checks and damage, measure
names and trait keys. Each measure and trait must quote a phrase its own entry
prints, a measure's phrase must print its number and unit, every digit left in
an entry must belong to a measure, and every sentence must be covered.
"""
import hashlib
import json
from fractions import Fraction
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[1]
DATA = ROOT / "internal/rules/srd521/testdata"
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
TABLE_PAGE = 95
ENTRY_PAGES = (94, 96, 97, 98, 99, 100)
HALF_WIDTH = 306
PRICE = r"[\d,]+ (?:CP|SP|GP)"
WEIGHT = r"—|[\d½/]+ lb\.(?: \(full\))?"
TABLE_ROW = re.compile(rf"(?P<name>.+?)\s{{2,}}(?P<weight>Varies|{WEIGHT})\s+(?P<price>Varies|{PRICE})")
HEADING = re.compile(rf"(?P<name>[A-Z][^()]*?) \((?P<price>Varies|{PRICE})\)")
SCROLL_HEADING = ("Spell Scroll (Cantrip, 30 GP;", "Level 1, 50 GP)")
VARIANT_TABLES = {"Ammunition": "Ammunition", "Arcane Focuses": "Arcane Focus",
                  "Druidic Focuses": "Druidic Focus", "Holy Symbols": "Holy Symbol"}
AMMUNITION_ROW = re.compile(rf"(?P<name>.+?)\s{{2,}}(?P<amount>\d+)\s+(?P<storage>[A-Z][a-z]+)\s+(?P<weight>{WEIGHT})\s+(?P<price>{PRICE})")
FOCUS_ROW = re.compile(rf"(?P<name>.+?)\s{{2,}}(?P<weight>{WEIGHT})\s+(?P<price>{PRICE})")
ABILITY = "Strength|Dexterity|Constitution|Intelligence|Wisdom|Charisma"
SAVE_RULE = "plus-dexterity-modifier-and-proficiency-bonus"
CHECKS = (
    ("ability-check", re.compile(rf"DC (?P<dc>\d+) (?P<ability>{ABILITY}) \((?P<skill>[A-Z][A-Za-z ]+)\) check")),
    ("saving-throw", re.compile(rf"DC (?P<dc>\d+) (?P<ability>{ABILITY}) saving throw")),
    ("saving-throw", re.compile(rf"(?P<ability>{ABILITY}) saving throw \(DC (?P<dc>8) (?P<rule>plus your Dexterity modifier and Proficiency Bonus)\)")),
)
DAMAGE = re.compile(r"\b(?P<amount>\d+d\d+|\d+) (?P<type>Acid|Bludgeoning|Cold|Fire|Force|Piercing|Poison|Radiant|Slashing) damage")
ACTIONS = (
    ("attack-replacement", r"replace one of your attacks"),
    ("attack-action", r"you can attack with the"),
    ("utilize", r"Utilize action"),
    ("bonus-action", r"Bonus Action"),
    ("action", r"\b(?:as an action|use its action|take an action)\b"),
)
FOCUS = re.compile(r"An? (?P<classes>[A-Z][A-Za-z, ]+?) can use (?:such an item|such an object|a Holy Symbol) as a Spellcasting Focus")
CONTENTS = re.compile(r"contains the following items: (?P<list>[^.]+)\.")
INLINE_HEADINGS = ("Dousing a Creature or an Object.", "Dousing a Space.", "Fuel.")
# Words that a measure's phrase prints for its number or unit.
NUMBER_WORDS = {"one": 1, "two": 2, "three": 3, "four": 4, "ten": 10, "twice": 2, "half": Fraction(1, 2), "one-fifth": Fraction(1, 5)}
UNIT_WORDS = {
    "additional-dim-light": "", "arrows": "Arrows", "armor-class": r"\bAC\b", "bolts": "Bolts", "bonus": "bonus", "character": "character",
    "creature": "creature", "creatures": "sleeps", "cubic-foot": r"cubic foot", "cubic-feet": "cubic feet",
    "feet": r"foot|feet", "gallon": "gallon", "gallons": "gallons", "hit-points": r"HP|Hit Points", "hour": "hour",
    "hours": "hours", "minute": "minute", "minutes": "minutes", "ounce": "ounce", "ounces": "ounce", "pages": "pages",
    "pieces": "pieces", "pint": "pint", "pints": "pint", "pounds": "pounds", "rounds": "rounds", "seconds": "seconds",
    "sheets": "sheets", "speed": "Speed", "spikes": "bundles", "times-normal-lift": "times the weight you can normally lift",
    "times-size": r"twice their size", "uses": "uses", "weapon": "weapon", "words": "words",
}
# A sentence that prints a measured number a second time.
REPEATS = {"oil": ("fuel-burn-duration", "burned for a total of 6 hours")}


def key(text):
    return re.sub(r"[^a-z0-9]+", "-", text.replace("’", "").lower()).strip("-")


def page_lines(pdf, page):
    """One page's lines in reading order: the left, then the right column."""
    for left in (0, HALF_WIDTH):
        args = ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(left), "-y", "0",
                "-W", str(HALF_WIDTH), "-H", "800", "-layout", str(pdf), "-"]
        for line in subprocess.run(args, check=True, capture_output=True, text=True).stdout.splitlines():
            line = line.strip()
            if line and "System Reference Document 5.2.1" not in line:
                yield line


def gear_table(pdf):
    """The Adventuring Gear table: name to (weight, price), in print order."""
    rows, scrolls, started = {}, [], False
    for line in page_lines(pdf, TABLE_PAGE):
        if line == "Adventuring Gear":
            started = True
            continue
        if not started or line.split() == ["Item", "Weight", "Cost"]:
            continue
        match = TABLE_ROW.fullmatch(line)
        assert match, ("unparsed table line", line)
        name = match["name"]
        if name.startswith("Spell Scroll ("):
            scrolls.append((name, match["weight"], match["price"]))
            name = "Spell Scroll"
        else:
            assert name not in rows, name
        rows[name] = (match["weight"], match["price"])
    assert scrolls == [("Spell Scroll (Cantrip)", "—", "30 GP"), ("Spell Scroll (Level 1)", "—", "50 GP")], scrolls
    return rows


def join(first, second):
    """Join wrapped lines. A word the layout hyphenated is closed; the hyphen
    of a compound such as 5-foot-square stays."""
    if not first:
        return second
    if first.endswith("-"):
        word = first.split()[-1]
        if re.search(r"\d-$|-\w+-$", word) or not second[:1].islower():
            return first + second
        return first[:-1] + second
    return first + " " + second


def gear_entries(pdf, table):
    """Every gear entry's text and page span, and the variant tables."""
    entries, tables, current, rows, started = {}, {}, None, None, False
    for page in ENTRY_PAGES:
        lines = page_lines(pdf, page)
        for line in lines:
            if not started:
                started = line == "Adventuring Gear"
                continue
            if line == "Mounts and Vehicles":
                assert list(entries) == list(table), ("entries disagree with the table", set(table) ^ set(entries))
                assert set(tables) == set(VARIANT_TABLES.values()), list(tables)
                return entries, tables
            heading = HEADING.fullmatch(line)
            if line == SCROLL_HEADING[0]:
                assert next(lines) == SCROLL_HEADING[1]
                heading = {"name": "Spell Scroll", "price": "Varies"}
            if heading and heading["name"] in table:
                name = heading["name"]
                assert name not in entries and (name == "Spell Scroll" or heading["price"] == table[name][1]), (name, heading["price"])
                current = entries[name] = {"price": heading["price"], "pages": [page, page], "text": ""}
                rows = None
                continue
            if line in VARIANT_TABLES:
                current, rows = None, tables.setdefault(VARIANT_TABLES[line], [])
                entries[VARIANT_TABLES[line]]["pages"][1] = page
                continue
            if rows is not None:
                pattern = AMMUNITION_ROW if rows is tables.get("Ammunition") else FOCUS_ROW
                row = pattern.fullmatch(line)
                if row:
                    rows.append(row.groupdict())
                else:
                    assert re.fullmatch(r"(?:Type|Focus|Symbol)\s+.*Cost", line), ("unparsed variant table line", line)
                continue
            if current is None:
                assert not entries, ("text outside an entry", line)
                continue  # the section introduction
            current["pages"][1] = page
            current["text"] = join(current["text"], line)
    raise AssertionError("gear section end not found")


def weight(text):
    if text == "—":
        return {"listed": False, "pounds_numerator": 0, "pounds_denominator": 1, "condition": "none"}
    found = re.fullmatch(r"(\d*?)(½|1/2)? lb\.( \(full\))?", text)
    assert found, text
    pounds = Fraction(int(found[1] or 0)) + (Fraction(1, 2) if found[2] else 0)
    return {"listed": True, "pounds_numerator": pounds.numerator, "pounds_denominator": pounds.denominator,
            "condition": "full" if found[3] else "none"}


def price(text):
    amount, currency = text.split()
    return {"amount": int(amount.replace(",", "")), "currency": currency.lower()}


def numbers(phrase):
    """Every number a phrase prints, in digits or words. 11/2 is one and a half."""
    found = set()
    for word in re.findall(r"one-fifth|\d+/\d+|\d[\d,]*|[A-Za-z]+", phrase):
        mixed = re.fullmatch(r"(\d*)1/2", word)
        if mixed:
            found.add(int(mixed[1] or 0) + Fraction(1, 2))
        elif word[0].isdigit():
            found.add(Fraction(word.replace(",", "")))
        elif word.lower() in NUMBER_WORDS:
            found.add(Fraction(NUMBER_WORDS[word.lower()]))
    return found


def checks(text):
    found = []
    for kind, pattern in CHECKS:
        for m in pattern.finditer(text):
            parts = m.groupdict()
            found.append((m.span(), {"kind": kind, "ability": parts["ability"].lower(),
                                     "skill": key(parts.get("skill") or "none"), "dc": int(parts["dc"]),
                                     "dc_rule": SAVE_RULE if parts.get("rule") else "fixed"}))
    found.sort()
    assert len(re.findall(r"\bDC \d+", text)) == len(found), ("unparsed DC", text)
    return found


def damage(text):
    found = []
    for m in DAMAGE.finditer(text):
        dice, _, sides = m["amount"].partition("d")
        amount = {"dice": int(dice), "sides": int(sides), "flat": 0} if sides else {"dice": 0, "sides": 0, "flat": int(dice)}
        found.append((m.span(), dict(amount, type=m["type"].lower())))
    assert len(re.findall(r"\d+d\d+", text)) == sum(1 for _, d in found if d["dice"]), ("unparsed dice", text)
    return found


def actions(text):
    found = []
    for name, pattern in ACTIONS:
        spans = [m.span() for m in re.finditer(pattern, text)]
        if spans:
            found.append((spans[0][0], name, spans))
    found.sort()
    return [name for _, name, _ in found], [span for _, _, spans in found for span in spans]


def focus(text):
    m = FOCUS.search(text)
    if not m:
        assert "Spellcasting Focus" not in text, text
        return [], []
    return [c.lower() for c in re.split(r",\s*(?:or )?| or ", m["classes"]) if c], [m.span()]


def contents(text, table):
    m = CONTENTS.search(text)
    if not m:
        return [], []
    # The table prints "Lantern, Hooded"; a pack list prints "Hooded Lantern".
    printed = {" ".join(reversed(name.split(", "))): key(name) for name in table}
    out = []
    for part in re.split(r",\s*(?:and )?", m["list"]):
        found = re.fullmatch(r"(?:(\d+) )?(?:(flask|day|sheet)s? of )?(.+)", part)
        assert found, part
        name = next((n for n in (found[3], found[3].removesuffix("s"), found[3].removesuffix("es")) if n in printed), None)
        assert name, ("unknown pack item", part)
        out.append({"item": printed[name], "quantity": int(found[1] or 1), "unit": found[2] or "item"})
    return out, [m.span()]


def spans_of(text, phrase, name):
    spans = [m.span() for m in re.finditer(re.escape(phrase), text)]
    assert spans, (name, "phrase is not printed in the entry", phrase)
    return spans


def entry_payload(name, entry, table, evidence):
    """The payload fields every variant of one entry shares."""
    text, slug = entry["text"], key(name)
    action_names, covered = actions(text)
    classes, spans = focus(text)
    covered += spans
    pack, pack_spans = contents(text, table)
    covered += pack_spans
    numbered = list(pack_spans)
    payload = {"actions": action_names, "spellcasting_focus": classes, "contents": pack}
    for field, parsed in (("checks", checks(text)), ("damage", damage(text))):
        purposes = evidence.get(field, [])
        assert len(purposes) == len(parsed), (name, field, "printed", [p for _, p in parsed], purposes)
        payload[field] = [dict(values, purpose=purpose) for purpose, (_, values) in zip(purposes, parsed)]
        numbered += [span for span, _ in parsed]
    measures = []
    for measure, value, unit, phrase in evidence.get("measures", []):
        value = Fraction(value)
        assert value in numbers(phrase), (name, measure, "number is not in the phrase", phrase)
        assert re.search(UNIT_WORDS[unit], phrase), (name, measure, "unit is not in the phrase", unit, phrase)
        numbered += spans_of(text, phrase, name)
        measures.append({"name": measure, "numerator": value.numerator, "denominator": value.denominator, "unit": unit})
    if slug in REPEATS:
        repeated, phrase = REPEATS[slug]
        value = next(m["numerator"] for m in measures if m["name"] == repeated)
        assert Fraction(value) in numbers(phrase), (name, phrase)
        numbered += spans_of(text, phrase, name)
    traits = []
    for trait, phrase in evidence.get("traits", []):
        covered += spans_of(text, phrase, name)
        traits.append(trait)
    covered += numbered
    # Every digit belongs to a check, damage, pack list or measure.
    for m in re.finditer(r"\d+", text):
        assert any(a <= m.start() and m.end() <= b for a, b in numbered), (name, "uncatalogued number", text[max(0, m.start() - 40):m.end() + 20])
    # Every sentence has some evidence.
    masked = text
    for heading in INLINE_HEADINGS:
        masked = masked.replace(heading, " " * len(heading))
    start = 0
    for end in [m.end() for m in re.finditer(r"\.”?(?=\s+[A-Z“])", masked)] + [len(masked)]:
        if masked[start:end].strip():
            assert any(a < end and b > start for a, b in covered), (name, "uncatalogued sentence", masked[start:end].strip())
        start = end
    return payload, measures, traits


def forms(name, entry, table, tables):
    """(slug, name, weight, price, price page, measures, traits) per form."""
    row_weight, row_price = table[name]
    if name not in tables:
        return [("standard", name, weight(row_weight), price(row_price), entry["pages"][0], [], [])]
    assert row_weight == row_price == entry["price"] == "Varies", name
    out = []
    for row in tables[name]:
        label, _, note = row["name"].partition(" (")
        measures, traits = [], []
        if "amount" in row:
            measures = [{"name": "bundle-amount", "numerator": int(row["amount"]), "denominator": 1, "unit": "pieces"}]
            traits = [f"stored-in-{row['storage'].lower()}"]
        if note:
            words = [w for w in note.removesuffix(")").split() if w.lower() not in ("a", "an", "the")]
            traits.append(key(" ".join(words)))
        # A variant table is on the last page of its entry.
        out.append((key(label), label, weight(row["weight"]), price(row["price"]), entry["pages"][1], measures, traits))
    return out


def keys(p):
    """The mechanic keys that restate a payload."""
    out = ["item.gear"] + [f"action.{a}" for a in p["actions"]] + [f"spellcasting-focus.{c}" for c in p["spellcasting_focus"]]
    for c in p["checks"]:
        value = f"{c['kind']}.{c['purpose']}-{c['ability']}" + ("" if c["skill"] == "none" else f"-{c['skill']}") + f"-dc-{c['dc']}"
        out.append(value + ("" if c["dc_rule"] == "fixed" else f"-{c['dc_rule']}"))
    for d in p["damage"]:
        amount = f"{d['dice']}d{d['sides']}" if d["dice"] else str(d["flat"])
        out.append(f"damage.{d['purpose']}-{amount}-{d['type']}")
    for m in p["measures"]:
        amount = str(m["numerator"]) + ("" if m["denominator"] == 1 else f"-over-{m['denominator']}")
        out.append(f"measure.{m['name']}-{amount}-{m['unit']}")
    out += [f"trait.{t}" for t in p["traits"]]
    return out + [f"content.{c['item']}-{c['quantity']}-{c['unit']}" for c in p["contents"]]


def same(got, want):
    """Equality that does not treat False, 0 and 0.0 as the same value."""
    return json.dumps(got, sort_keys=True) == json.dumps(want, sort_keys=True)


def verify(argv, issue, batch, family_count, variant_count, evidence, joined):
    """Compare one batch. joined lists inventory IDs of the issue that another
    record owns; the batch must hold none of them."""
    if len(argv) != 2:
        raise SystemExit(f"usage: verify_srd_items_{issue}.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    table = gear_table(pdf)
    assert len(table) == 81, len(table)
    entries, tables = gear_entries(pdf, table)

    catalogue = json.loads((DATA / f"items_{issue}.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    gear = [r for r in inventory["records"] if r["family"] == "adventuring-gear"]
    # The planning inventory lists exactly the printed entries and their pages.
    assert sorted((r["name"], r["source_page"]) for r in gear) == sorted((n, e["pages"][0]) for n, e in entries.items())
    listed = [r for r in gear if r["catalogue_issue"] == issue]
    assert [r["id"] for r in listed if r["id"] in joined] == list(joined), joined
    expected = [(r["id"], r["name"], r["source_page"]) for r in listed if r["id"] not in joined]
    assert len(expected) == family_count, len(expected)
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == issue and catalogue["batch"] == batch
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert set(evidence) == {key(f["name"]) for f in families}

    total = 0
    for f in families:
        name, slug = f["name"], key(f["name"])
        entry = entries[name]
        pages = [entry["pages"][0]] * 2
        assert f["id"] == f"srd-5.2.1:content/items/{slug}" and f["kind"] == "mundane" and f["category"] == "adventuring-gear", name
        assert f["base"] == {"options": [], "excluded": []} and len(f["presentations"]) == 1, name
        assert f["presentations"][0] == {"id": f"srd-5.2.1:content/adventuring-gear/{slug}", "pages": pages}, name
        shared, measures, traits = entry_payload(name, entry, table, evidence[slug])
        printed = forms(name, entry, table, tables)
        assert len(f["variants"]) == len(printed), (name, len(printed))
        for v, (vslug, vname, vweight, vprice, price_page, extra_measures, extra_traits) in zip(f["variants"], printed):
            total += 1
            payload = dict(shared, measures=measures + extra_measures, traits=traits + extra_traits,
                           weight=vweight, source_pages=entry["pages"])
            use = "none" if not payload["actions"] else payload["actions"][0] if len(payload["actions"]) == 1 else "multiple"
            assert v["id"] == f["id"] + "/variant/" + vslug and v["name"] == vname and v["rarity"] == "none", (name, vname, vslug)
            assert v["attunement"] == {"mode": "none", "restrictions": []}, vname
            assert same(v["charges"], {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}), vname
            assert v["equip"] == {"mode": "none", "restrictions": []} and v["use"] == {"mode": use, "restrictions": []}, (vname, use)
            cost = dict(vprice, kind="listed", source_pages=[price_page, price_page])
            assert same(v["cost"], cost), (vname, v["cost"], cost)
            assert same(v["details"], {"gear": payload}), (vname, v["details"], payload)
            assert v["required_mechanic_keys"] == keys(payload), (vname, v["required_mechanic_keys"], keys(payload))
    assert total == variant_count, total
    return {"table": table, "entries": entries}


def owners(name):
    """Every item record file that holds a family of this name."""
    found = []
    for path in sorted(DATA.glob("items_*.json")) + [DATA / "item-families.json"]:
        for family in json.loads(path.read_text())["families"]:
            if family["name"] == name:
                found.append((path.name, family))
    return found


def check_potion_of_healing(context):
    """The p. 99 entry against the Healing contract fixture, its only owner."""
    entry, row = context["entries"]["Potion of Healing"], context["table"]["Potion of Healing"]
    [(path, family)] = owners("Potion of Healing")
    assert path == "item-families.json"
    assert {"id": "srd-5.2.1:content/adventuring-gear/potion-of-healing", "pages": entry["pages"]} in family["presentations"]
    basic = family["variants"][0]
    assert basic["name"] == "Potion of Healing" and entry["pages"] == [99, 99]
    assert same(basic["cost"], dict(price(entry["price"]), kind="listed", source_pages=entry["pages"])) and price(row[1]) == price(entry["price"])
    assert row[0] == "1/2 lb." and same(basic["weight"], weight(row[0]))
    assert all("weight" not in v for v in family["variants"][1:])  # the table weighs only this potion
    healed = re.search(r"regains (\d+)d(\d+) \+ (\d+) Hit Points", entry["text"])
    assert same(basic["healing"], {"dice": int(healed[1]), "sides": int(healed[2]), "bonus": int(healed[3])})
    assert "As a Bonus Action, you can drink it or administer it to another creature within 5 feet of yourself" in entry["text"]
    assert basic["use"] == {"mode": "drink-or-administer", "restrictions": ["bonus-action", "consume-one", "administer-within-5-feet"]}
    assert "This potion is a magic item" in entry["text"] and family["kind"] == "magic"


def check_spell_scroll(context):
    """The p. 99 entry against the #242 records, their only owner. They hold the
    gear listing and the two prices its heading prints; gear_table and
    gear_entries check the table rows and the heading against the PDF."""
    entry = context["entries"]["Spell Scroll"]
    text = entry["text"]
    [(path, family)] = owners("Spell Scroll")
    assert path == "items_242.json" and entry["pages"] == [99, 99]
    assert family["presentations"][1:] == [{"id": "srd-5.2.1:content/adventuring-gear/spell-scroll", "pages": entry["pages"]}]
    prices = [int(amount) for amount in re.findall(r"(\d+) GP", " ".join(SCROLL_HEADING))]
    found = re.search(r"the spell save DC is (\d+), and the attack bonus is \+(\d+)", text)
    assert (int(found[1]), int(found[2])) == (13, 5)
    for level, variant in enumerate(family["variants"][:2]):
        scroll = variant["details"]["scroll"]
        assert (scroll["spell_level"], scroll["save_dc"], scroll["attack_bonus"]) == (level, int(found[1]), int(found[2]))
        assert scroll["material_components_required"] is False and scroll["readable_if"] == "spell-on-your-spell-list"
        assert same(variant["cost"], {"kind": "listed", "amount": prices[level], "currency": "gp", "source_pages": entry["pages"]})
    assert all(v["cost"]["kind"] == "unspecified" for v in family["variants"][2:])
    for phrase in ("is a magic item that bears the words of a cantrip or level 1 spell", "If the spell is on your class’s spell list",
                   "without providing any Material components", "using its normal casting time"):
        assert phrase in text, phrase
