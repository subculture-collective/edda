#!/usr/bin/env python3
"""Compare issue #246 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (253,)
RARITY = {"common": "Common", "uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
CATEGORY = {"armor": "Armor (", "potion": "Potion, ", "ring": "Ring, ", "rod": "Rod, ", "scroll": "Scroll, ", "staff": "Staff, ", "weapon": "Weapon (", "wand": "Wand, ", "wondrous-item": "Wondrous Item, "}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# A category line under a heading. Each family block must hold only its own.
CATEGORY_LINE = re.compile(r"(?:Wondrous Item|Potion|Ring|Rod|Staff|Wand|Scroll|(?:Weapon|Armor) \([^)]*\)), (?:Common|Uncommon|Rare|Very Rare|Legendary|Artifact|Rarity Varies)")
NUMBER_WORDS = {"cantrip": 0, "a": 1, "an": 1, "once": 1, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10}
# The headings that follow the batch's last families in print.
FOLLOWING = []

# Each family's printed category line.
HEADERS = {
    'weapon-of-warning': 'Weapon (Any Simple or Martial), Uncommon (Requires Attunement)',
    'weapon-1-2-or-3': 'Weapon (Any Simple or Martial), Uncommon (+1), Rare (+2), or Very Rare (+3)',
    'well-of-many-worlds': 'Wondrous Item, Legendary',
    'wind-fan': 'Wondrous Item, Uncommon',
    'winged-boots': 'Wondrous Item, Uncommon (Requires Attunement)',
    'wings-of-flying': 'Wondrous Item, Rare (Requires Attunement)',
}

# Printed charge sentence of each family with a charge pool.
CHARGES = {
    'winged-boots': ({'mode': 'limited', 'maximum': 4, 'per_use': 1, 'recharge': 'regain-1d4-at-dawn'}, 'These boots have 4 charges and regain 1d4 expended charges daily at dawn'),
}

# Every payload leaf of every variant: (path, value, printed phrase). The
# catalogued payload must have exactly these leaves, and each phrase must be
# printed inside the family's own text. A "none" key or an empty list records
# that nothing is printed. source_pages is checked from the page layout.
VALUES = {
    'weapon-of-warning': {
        'standard': [
            ('condition', 'within-reach-while-attuned', 'As long as this weapon is within your reach and you are attuned to it'),
            ('radius_feet', 30, 'you and allies within 30 feet of you'),
            ('subjects', 'you-and-allies', 'you and allies within 30 feet of you gain the following benefits'),
            ('alarm.wakes', 'naturally-sleeping-subjects', 'magically awakens each subject who is sleeping naturally'),
            ('alarm.when', 'combat-begins', 'when combat begins'),
            ('alarm.wakes_from_magical_sleep', False, 'This benefit doesn’t wake a subject from magically induced sleep'),
            ('supernatural_readiness.roll', 'initiative', 'Each subject has Advantage on its Initiative rolls'),
            ('supernatural_readiness.mode', 'advantage', 'Each subject has Advantage on its Initiative rolls'),
        ],
    },
    'weapon-1-2-or-3': {
        'plus-1': [
            ('attack_bonus', 1, 'Uncommon (+1)'),
            ('damage_bonus', 1, 'Uncommon (+1)'),
            ('bonus_rule', 'determined-by-rarity', 'The bonus is determined by the weapon’s rarity'),
        ],
        'plus-2': [
            ('attack_bonus', 2, 'Rare (+2)'),
            ('damage_bonus', 2, 'Rare (+2)'),
            ('bonus_rule', 'determined-by-rarity', 'The bonus is determined by the weapon’s rarity'),
        ],
        'plus-3': [
            ('attack_bonus', 3, 'Very Rare (+3)'),
            ('damage_bonus', 3, 'Very Rare (+3)'),
            ('bonus_rule', 'determined-by-rarity', 'The bonus is determined by the weapon’s rarity'),
        ],
    },
    'well-of-many-worlds': {
        'standard': [
            ('material', 'fine-black-cloth', 'This fine black cloth, soft as silk'),
            ('folded_size', 'handkerchief', 'folded up to the dimensions of a handkerchief'),
            ('unfolded_diameter_feet', 6, 'a circular sheet 6 feet in diameter'),
            ('open.action', 'magic', 'You can take a Magic action to unfold the Well of Many Worlds'),
            ('open.placement', 'solid-surface', 'place it on a solid surface'),
            ('portal.direction', 'two-way', 'it forms a two-way'),
            ('portal.diameter_feet', 6, '6-foot-diameter, circular portal'),
            ('portal.destination', 'another-world-or-plane', 'portal to another world or plane of existence'),
            ('portal.chosen_by', 'gm', 'the GM decides where it leads'),
            ('close.action', 'magic', 'takes a Magic action to close it'),
            ('close.within_feet', 5, 'a creature within 5 feet of it'),
            ('close.method', 'fold-cloth-by-edges', 'by taking hold of the edges of the cloth and folding it up'),
            ('reuse_after', {'dice': 1, 'sides': 8, 'type': 'hours'}, 'it can’t do so again for 1d8 hours'),
        ],
    },
    'wind-fan': {
        'standard': [
            ('condition', 'while-holding', 'While holding this fan'),
            ('spell', 'gust-of-wind', 'you can cast Gust of Wind'),
            ('save_dc', 13, '(save DC 13)'),
            ('repeat_use.failure_percent_cumulative', 20, 'it has a cumulative 20 percent chance of not working'),
            ('repeat_use.applies_to', 'each-subsequent-use-before-next-dawn', 'Each subsequent time the fan is used before the next dawn'),
            ('repeat_use.on_failure', 'tears-into-nonmagical-tatters', 'if the fan fails to work, it tears into useless, nonmagical tatters'),
        ],
    },
    'winged-boots': {
        'standard': [
            ('max_charges', 4, 'These boots have 4 charges'),
            ('recharge_dice', {'dice': 1, 'sides': 4, 'type': 'charges'}, 'regain 1d4 expended charges daily at dawn'),
            ('action', 'magic', 'you can take a Magic action to expend 1 charge'),
            ('charge_cost', 1, 'a Magic action to expend 1 charge'),
            ('fly_speed_feet', 30, 'gaining a Fly Speed of 30 feet'),
            ('duration_hours', 1, 'for 1 hour'),
            ('on_expiry.condition', 'flying-when-duration-expires', 'If you are flying when the duration expires'),
            ('on_expiry.descent_feet_per_round', 30, 'you descend at a rate of 30 feet per round until you land'),
        ],
    },
    'wings-of-flying': {
        'standard': [
            ('form', 'cloak-becomes-wings-on-back', 'turn the cloak into a pair of wings on your back'),
            ('action', 'magic', 'you can take a Magic action to turn the cloak'),
            ('duration_hours', 1, 'The wings lasts for 1 hour'),
            ('early_end_action', 'magic', 'or until you end the effect early as a Magic action'),
            ('fly_speed_feet', 60, 'The wings give you a Fly Speed of 60 feet'),
            ('aloft_when_wings_disappear', 'fall', 'If you are aloft when the wings disappear, you fall'),
            ('reuse_after', {'dice': 1, 'sides': 12, 'type': 'hours'}, 'you can’t use them again for 1d12 hours'),
        ],
    },
}


def text(pdf, page, x):
    args = ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x), "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"]
    return subprocess.run(args, check=True, capture_output=True, text=True).stdout


def clean(raw):
    return FOOTER.sub(" ", " ".join(raw.split()))


def flat(value):
    """Ignore spaces and hyphens so a word split across a line end still matches."""
    return re.sub(r"[\s\-]", "", value)


def has(source, phrase):
    return flat(phrase) in flat(source)


def numbers(phrase):
    """The numbers a phrase prints, as digits or words. A printed 00 means 100."""
    digits = re.findall(r"\d+", re.sub(r"(?<=\d),(?=\d)", "", phrase))
    found = {int(d) for d in digits} | ({100} if "00" in digits else set())
    return found | {NUMBER_WORDS[w] for w in re.findall(r"[a-z]+", phrase.lower()) if w in NUMBER_WORDS}


def leaves(value, path=""):
    if isinstance(value, dict) and set(value) != {"dice", "sides", "type"}:
        for key, child in value.items():
            if key != "source_pages":
                yield from leaves(child, f"{path}.{key}" if path else key)
    elif isinstance(value, list) and value:
        for i, child in enumerate(value):
            yield from leaves(child, f"{path}[{i}]")
    else:
        yield path, value


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_246.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream, spans = "", {}
    for page in PAGES:
        start = len(stream)
        stream += " ".join(clean(text(pdf, page, x)).strip() for x in (0, 306))
        spans[page] = (start, len(stream))
        stream += " "

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_246.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 246]
    assert len(expected) == 6
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 246 and catalogue["batch"] == "data-magic-items-22"
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert [f["id"].rsplit("/", 1)[1] for f in families] == list(HEADERS) == list(VALUES)

    starts = {}
    for f in families:
        slug = f["id"].rsplit("/", 1)[1]
        printed = f"{f['name']} {HEADERS[slug]}"
        assert stream.count(printed) == 1, printed
        starts[slug] = stream.index(printed)
    bounds = sorted(starts.values())
    for heading in FOLLOWING:
        bounds.append(stream.index(heading, bounds[0]))
    count = 0
    for f in families:
        slug = f["id"].rsplit("/", 1)[1]
        start = starts[slug]
        end = min([b for b in bounds if b > start] + [len(stream)])
        block = stream[start:end].strip()
        assert len(CATEGORY_LINE.findall(block)) == 1, (slug, CATEGORY_LINE.findall(block))
        header = HEADERS[slug]
        assert header.startswith(CATEGORY[f["category"]]), (slug, header)
        pages = [p for p, (a, b) in spans.items() if a < start + len(block) and start < b]
        assert f["presentations"][0]["pages"] == [pages[0], pages[0]], (slug, pages)
        pool = CHARGES.get(slug)
        if pool is None:
            assert not re.search(r"(?:has|have) \d+ charges", block), slug
        else:
            assert has(block, pool[1]), (slug, pool[1])
            assert pool[0]["maximum"] in numbers(pool[1]), slug
        assert [v["id"].rsplit("/", 1)[1] for v in f["variants"]] == list(VALUES[slug]), slug
        for v in f["variants"]:
            vslug = v["id"].rsplit("/", 1)[1]
            label = f"{slug}/{vslug}"
            assert ("Requires Attunement" in header) == (v["attunement"]["mode"] == "required"), label
            assert "Rarity Varies" in header or RARITY[v["rarity"]] in header, label
            assert v["attunement"]["restrictions"] == [], label
            assert v["charges"] == (pool[0] if pool else {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}), label
            assert v["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": f["presentations"][0]["pages"]}, label
            assert len(v["details"]) == 1, label
            payload = next(iter(v["details"].values()))
            assert payload.get("rarity", v["rarity"]) == v["rarity"], label
            # A family whose text crosses a page records both pages.
            assert payload.get("source_pages") == (pages if len(pages) > 1 else None), (label, pages)
            want = VALUES[slug][vslug]
            got = list(leaves(payload))
            assert [p for p, _ in got] == [row[0] for row in want], (label, [p for p, _ in got])
            for (path, value), (_, printed_value, phrase) in zip(got, want):
                where = (label, path, value, phrase)
                assert type(value) is type(printed_value) and value == printed_value, where
                if phrase is None:
                    assert value == [], where
                    continue
                assert has(block, phrase), where
                if isinstance(value, dict):
                    assert f"{value['dice']}d{value['sides']}" in phrase, where
                elif type(value) is int:
                    assert value in numbers(phrase), where
                count += 1
    print(f"ok: {len(families)} families, {sum(len(f['variants']) for f in families)} variants, {count} printed values, p. 253")


main()
