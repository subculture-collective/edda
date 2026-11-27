#!/usr/bin/env python3
"""Compare issue #239 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (238, 239)
RARITY = {"common": "Common", "uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
CATEGORY = {"armor": "Armor (", "potion": "Potion, ", "ring": "Ring, ", "rod": "Rod, ", "scroll": "Scroll, ", "staff": "Staff, ", "weapon": "Weapon (", "wondrous-item": "Wondrous Item, "}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# A category line under a heading. Each family block must hold only its own.
CATEGORY_LINE = re.compile(r"(?:Wondrous Item|Potion|Ring|Rod|Staff|Wand|Scroll|(?:Weapon|Armor) \([^)]*\)), (?:Common|Uncommon|Rare|Very Rare|Legendary|Artifact|Rarity Varies)")
NUMBER_WORDS = {"cantrip": 0, "a": 1, "an": 1, "once": 1, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10}
# The headings that follow the batch's last families in print.
FOLLOWING = ['Ring of the Ram Ring,']

# Each family's printed category line.
HEADERS = {
    'ring-of-free-action': 'Ring, Rare (Requires Attunement)',
    'ring-of-invisibility': 'Ring, Legendary (Requires Attunement)',
    'ring-of-jumping': 'Ring, Uncommon (Requires Attunement)',
    'ring-of-mind-shielding': 'Ring, Uncommon (Requires Attunement)',
    'ring-of-protection': 'Ring, Rare (Requires Attunement)',
    'ring-of-regeneration': 'Ring, Very Rare (Requires Attunement)',
    'ring-of-resistance': 'Ring, Rare',
    'ring-of-shooting-stars': 'Ring, Very Rare (Requires Attunement)',
    'ring-of-spell-storing': 'Ring, Rare (Requires Attunement)',
    'ring-of-spell-turning': 'Ring, Legendary (Requires Attunement)',
    'ring-of-swimming': 'Ring, Uncommon',
    'ring-of-telekinesis': 'Ring, Very Rare (Requires Attunement)',
}

# Printed charge sentence of each family with a charge pool.
CHARGES = {
    'ring-of-shooting-stars': ({'mode': 'limited', 'maximum': 6, 'per_use': 1, 'recharge': 'regain-1d6-at-dawn'}, 'The ring has 6 charges and regains 1d6 expended charges daily at dawn'),
}

# Every payload leaf of every variant: (path, value, printed phrase). The
# catalogued payload must have exactly these leaves, and each phrase must be
# printed inside the family's own text. A "none" key or an empty list records
# that nothing is printed. source_pages is checked from the page layout.
VALUES = {
    'ring-of-free-action': {
        'standard': [
            ('condition', 'while-wearing', 'While you wear this ring'),
            ('difficult_terrain_costs_extra_movement', False, 'Difficult Terrain doesn’t cost you extra movement'),
            ('magic_cannot[0]', 'reduce-any-speed', 'magic can neither reduce any of your Speeds'),
            ('magic_cannot[1]', 'cause-paralyzed-condition', 'nor cause you to have the Paralyzed or'),
            ('magic_cannot[2]', 'cause-restrained-condition', 'or Restrained condition'),
        ],
    },
    'ring-of-invisibility': {
        'standard': [
            ('condition', 'while-wearing', 'While wearing this ring'),
            ('action', 'magic', 'you can take a Magic action to give yourself the Invisible condition'),
            ('grants_condition', 'invisible', 'give yourself the Invisible condition'),
            ('ends_when[0]', 'ring-removed', 'until the ring is removed'),
            ('ends_when[1]', 'bonus-action-to-become-visible', 'or until you take a Bonus Action to become visible again'),
        ],
    },
    'ring-of-jumping': {
        'standard': [
            ('condition', 'while-wearing', 'While wearing this ring'),
            ('spell', 'jump', 'you can cast Jump from it'),
            ('target', 'yourself-only', 'can target only yourself'),
        ],
    },
    'ring-of-mind-shielding': {
        'standard': [
            ('condition', 'while-wearing', 'While wearing this ring'),
            ('immune_to_magic_that[0]', 'reads-thoughts', 'magic that allows other creatures to read your thoughts'),
            ('immune_to_magic_that[1]', 'detects-lies', 'determine whether you are lying'),
            ('immune_to_magic_that[2]', 'reveals-alignment', 'know your alignment'),
            ('immune_to_magic_that[3]', 'reveals-creature-type', 'or know your creature type'),
            ('telepathy_needs_permission', True, 'Creatures can telepathically communicate with you only if you allow it'),
            ('imperceptible.action', 'magic', 'You can take a Magic action to cause the ring to become imperceptible'),
            ('imperceptible.ends_when[0]', 'magic-action-to-make-perceptible', 'until you take another Magic action to make it perceptible'),
            ('imperceptible.ends_when[1]', 'ring-removed', 'until you remove the ring'),
            ('imperceptible.ends_when[2]', 'wearer-dies', 'or until you die'),
            ('soul.trigger', 'die-while-wearing', 'If you die while wearing the ring, your soul enters it'),
            ('soul.blocked_if', 'ring-already-houses-a-soul', 'unless it already houses a soul'),
            ('soul.choices[0]', 'remain-in-ring', 'You can remain in the ring'),
            ('soul.choices[1]', 'depart-for-afterlife', 'or depart for the afterlife'),
            ('soul.telepathy_while', 'soul-is-in-the-ring', 'As long as your soul is in the ring'),
            ('soul.telepathy_with', 'any-creature-wearing-the-ring', 'you can telepathically communicate with any creature wearing it'),
            ('soul.wearer_can_prevent', False, 'A wearer can’t prevent this telepathic communication'),
        ],
    },
    'ring-of-protection': {
        'standard': [
            ('condition', 'while-wearing', 'while wearing this ring'),
            ('armor_class_bonus', 1, '+1 bonus to Armor Class'),
            ('saving_throw_bonus', 1, '+1 bonus to Armor Class and saving throws'),
        ],
    },
    'ring-of-regeneration': {
        'standard': [
            ('condition', 'while-wearing', 'While wearing this ring'),
            ('healing', {'dice': 1, 'sides': 6, 'type': 'hit-points'}, 'you regain 1d6 Hit Points'),
            ('interval_minutes', 10, 'every 10 minutes'),
            ('minimum_hit_points', 1, 'if you have at least 1 Hit Point'),
            ('regrowth.trigger', 'lose-a-body-part', 'If you lose a body part'),
            ('regrowth.result', 'part-regrows-with-full-functionality', 'regrow and return to full functionality'),
            ('regrowth.delay', {'dice': 1, 'sides': 6, 'type': 'days'}, 'after 1d6 + 1 days'),
            ('regrowth.delay_bonus_days', 1, '1d6 + 1 days'),
            ('regrowth.requires_hit_points_whole_time', True, 'if you have at least 1 Hit Point the whole time'),
        ],
    },
    'ring-of-resistance': {
        'standard': [
            ('condition', 'while-wearing', 'while wearing this ring'),
            ('defense', 'resistance', 'You have Resistance to one damage type'),
            ('type_shown_by', 'gemstone', 'The gemstone in the ring indicates the type'),
            ('type_chosen_by', 'gm-chooses-or-rolls-randomly', 'which the GM chooses or determines randomly'),
            ('table_die_sides', 10, '1d10 Damage Type Gemstone'),
            ('damage_types[0].roll', 1, '1 Acid Pearl'),
            ('damage_types[0].damage_type', 'acid', '1 Acid Pearl'),
            ('damage_types[0].gemstone', 'pearl', '1 Acid Pearl'),
            ('damage_types[1].roll', 2, '2 Cold Tourmaline'),
            ('damage_types[1].damage_type', 'cold', '2 Cold Tourmaline'),
            ('damage_types[1].gemstone', 'tourmaline', '2 Cold Tourmaline'),
            ('damage_types[2].roll', 3, '3 Fire Garnet'),
            ('damage_types[2].damage_type', 'fire', '3 Fire Garnet'),
            ('damage_types[2].gemstone', 'garnet', '3 Fire Garnet'),
            ('damage_types[3].roll', 4, '4 Force Sapphire'),
            ('damage_types[3].damage_type', 'force', '4 Force Sapphire'),
            ('damage_types[3].gemstone', 'sapphire', '4 Force Sapphire'),
            ('damage_types[4].roll', 5, '5 Lightning Citrine'),
            ('damage_types[4].damage_type', 'lightning', '5 Lightning Citrine'),
            ('damage_types[4].gemstone', 'citrine', '5 Lightning Citrine'),
            ('damage_types[5].roll', 6, '6 Necrotic Jet'),
            ('damage_types[5].damage_type', 'necrotic', '6 Necrotic Jet'),
            ('damage_types[5].gemstone', 'jet', '6 Necrotic Jet'),
            ('damage_types[6].roll', 7, '7 Poison Amethyst'),
            ('damage_types[6].damage_type', 'poison', '7 Poison Amethyst'),
            ('damage_types[6].gemstone', 'amethyst', '7 Poison Amethyst'),
            ('damage_types[7].roll', 8, '8 Psychic Jade'),
            ('damage_types[7].damage_type', 'psychic', '8 Psychic Jade'),
            ('damage_types[7].gemstone', 'jade', '8 Psychic Jade'),
            ('damage_types[8].roll', 9, '9 Radiant Topaz'),
            ('damage_types[8].damage_type', 'radiant', '9 Radiant Topaz'),
            ('damage_types[8].gemstone', 'topaz', '9 Radiant Topaz'),
            ('damage_types[9].roll', 10, '10 Thunder Spinel'),
            ('damage_types[9].damage_type', 'thunder', '10 Thunder Spinel'),
            ('damage_types[9].gemstone', 'spinel', '10 Thunder Spinel'),
        ],
    },
    'ring-of-shooting-stars': {
        'standard': [
            ('free_spells[0]', 'dancing-lights', 'You can cast Dancing Lights or Light'),
            ('free_spells[1]', 'light', 'or Light from the ring'),
            ('max_charges', 6, 'The ring has 6 charges'),
            ('recharge', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 expended charges daily at dawn'),
            ('recharge_when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('faerie_fire.charge_cost', 1, 'You can expend 1 charge to cast Faerie Fire'),
            ('faerie_fire.spell', 'faerie-fire', 'to cast Faerie Fire from the ring'),
            ('lightning_spheres.charge_cost', 2, 'You can expend 2 charges as a Magic action'),
            ('lightning_spheres.action', 'magic', 'as a Magic action to create up to four'),
            ('lightning_spheres.max_spheres', 4, 'create up to four 3-foot-diameter spheres'),
            ('lightning_spheres.sphere_diameter_feet', 3, '3-foot-diameter spheres of lightning'),
            ('lightning_spheres.placement_range_feet', 120, 'an unoccupied space you can see within 120 feet of yourself'),
            ('lightning_spheres.concentration_required', True, 'The spheres last as long as you maintain Concentration'),
            ('lightning_spheres.max_duration_minutes', 1, 'up to 1 minute'),
            ('lightning_spheres.dim_light_radius_feet', 30, 'Each sphere sheds Dim Light in a 30-foot radius'),
            ('lightning_spheres.move_action', 'bonus', 'As a Bonus Action, you can move each sphere'),
            ('lightning_spheres.move_feet', 30, 'move each sphere up to 30 feet'),
            ('lightning_spheres.leash_feet', 120, 'but no farther than 120 feet away from yourself'),
            ('lightning_spheres.discharge_distance_feet', 5, 'comes within 5 feet of a creature other than you'),
            ('lightning_spheres.discharge_excludes[0]', 'you', 'a creature other than you'),
            ('lightning_spheres.discharge_excludes[1]', 'behind-total-cover', 'that isn’t behind Total Cover'),
            ('lightning_spheres.sphere_disappears_on_discharge', True, 'discharges lightning at that creature and disappears'),
            ('lightning_spheres.save_ability', 'dexterity', 'DC 15 Dexterity saving throw. On a failed save, the creature takes Lightning damage'),
            ('lightning_spheres.save_dc', 15, 'DC 15 Dexterity saving throw'),
            ('lightning_spheres.success_damage', 'half', 'On a successful save, the creature takes half as much damage'),
            ('lightning_spheres.damage_table[0].spheres', 1, '1 4d12'),
            ('lightning_spheres.damage_table[0].damage', {'dice': 4, 'sides': 12, 'type': 'lightning'}, '1 4d12'),
            ('lightning_spheres.damage_table[1].spheres', 2, '2 5d4'),
            ('lightning_spheres.damage_table[1].damage', {'dice': 5, 'sides': 4, 'type': 'lightning'}, '2 5d4'),
            ('lightning_spheres.damage_table[2].spheres', 3, '3 2d6'),
            ('lightning_spheres.damage_table[2].damage', {'dice': 2, 'sides': 6, 'type': 'lightning'}, '3 2d6'),
            ('lightning_spheres.damage_table[3].spheres', 4, '4 2d4'),
            ('lightning_spheres.damage_table[3].damage', {'dice': 2, 'sides': 4, 'type': 'lightning'}, '4 2d4'),
            ('shooting_stars.minimum_charges', 1, 'You can expend 1 to 3 charges as a Magic action'),
            ('shooting_stars.maximum_charges', 3, 'expend 1 to 3 charges'),
            ('shooting_stars.action', 'magic', '1 to 3 charges as a Magic action. For every charge'),
            ('shooting_stars.motes_per_charge', 1, 'For every charge you expend, you launch a glowing mote of light'),
            ('shooting_stars.range_feet', 60, 'at a point you can see within 60 feet of yourself'),
            ('shooting_stars.area_shape', 'cube', '15-foot Cube originating from that point'),
            ('shooting_stars.area_feet', 15, '15-foot Cube'),
            ('shooting_stars.save_ability', 'dexterity', 'is showered in sparks and makes a DC 15 Dexterity saving throw'),
            ('shooting_stars.save_dc', 15, 'showered in sparks and makes a DC 15'),
            ('shooting_stars.damage', {'dice': 5, 'sides': 4, 'type': 'radiant'}, 'taking 5d4 Radiant damage'),
            ('shooting_stars.success_damage', 'half', 'or half as much damage on a successful one'),
        ],
    },
    'ring-of-spell-storing': {
        'standard': [
            ('holds', 'spells-cast-into-it-until-attuned-wearer-uses-them', 'stores spells cast into it, holding them until the attuned wearer uses them'),
            ('capacity_spell_levels', 5, 'up to 5 levels worth of spells'),
            ('found_levels', {'dice': 1, 'sides': 6, 'type': 'stored-spell-levels'}, 'it contains 1d6 − 1 levels of stored spells'),
            ('found_levels_penalty', 1, '1d6 − 1 levels'),
            ('found_spells_chosen_by', 'gm', 'stored spells chosen by the GM'),
            ('storing.caster', 'any-creature', 'Any creature can cast a spell'),
            ('storing.minimum_spell_level', 1, 'a spell of level 1 through 5'),
            ('storing.maximum_spell_level', 5, 'level 1 through 5 into the ring'),
            ('storing.method', 'touch-ring-as-spell-is-cast', 'by touching the ring as the spell is cast'),
            ('storing.other_effect', 'none-besides-being-stored', 'The spell has no effect other than to be stored in the ring'),
            ('storing.if_it_cannot_fit', 'spell-expended-without-effect', 'the spell is expended without effect'),
            ('storing.space_used', 'slot-level-used', 'The level of the slot used to cast the spell determines how much space it uses'),
            ('casting.condition', 'while-wearing', 'While wearing this ring, you can cast any spell stored in it'),
            ('casting.uses_from_original_caster[0]', 'slot-level', 'The spell uses the slot level'),
            ('casting.uses_from_original_caster[1]', 'spell-save-dc', 'spell save DC'),
            ('casting.uses_from_original_caster[2]', 'spell-attack-bonus', 'spell attack bonus'),
            ('casting.uses_from_original_caster[3]', 'spellcasting-ability', 'and spellcasting ability of the original caster'),
            ('casting.otherwise', 'treated-as-if-you-cast-it', 'is otherwise treated as if you cast the spell'),
            ('casting.frees_space', True, 'no longer stored in it, freeing up space'),
        ],
    },
    'ring-of-spell-turning': {
        'standard': [
            ('condition', 'while-wearing', 'While wearing this ring'),
            ('save_mode', 'advantage', 'you have Advantage on saving throws against spells'),
            ('negate_max_spell_level', 7, 'If you succeed on the save for a spell of level 7 or lower'),
            ('negate_result', 'spell-has-no-effect-on-you', 'the spell has no effect on you'),
            ('deflect.requires[0]', 'spell-targeted-only-you', 'If that spell targeted only you'),
            ('deflect.requires[1]', 'spell-created-no-area-of-effect', 'and didn’t create an area of effect'),
            ('deflect.action', 'reaction', 'you can take a Reaction to deflect the spell'),
            ('deflect.target', 'original-caster', 'back at the spell’s caster'),
            ('deflect.caster_save_uses', 'casters-own-spell-save-dc', 'the caster must make a saving throw against the spell using their own spell save DC'),
        ],
    },
    'ring-of-swimming': {
        'standard': [
            ('condition', 'while-wearing', 'while wearing this ring'),
            ('swim_speed_feet', 40, 'You have a Swim Speed of 40 feet'),
        ],
    },
    'ring-of-telekinesis': {
        'standard': [
            ('condition', 'while-wearing', 'While wearing this ring'),
            ('spell', 'telekinesis', 'you can cast Telekinesis from it'),
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
        raise SystemExit("usage: verify_srd_items_239.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream, spans = "", {}
    for page in PAGES:
        start = len(stream)
        stream += " ".join(clean(text(pdf, page, x)).strip() for x in (0, 306))
        spans[page] = (start, len(stream))
        stream += " "

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_239.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 239]
    assert len(expected) == 12
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 239 and catalogue["batch"] == "data-magic-items-15"
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
    print(f"ok: {len(families)} families, {sum(len(f['variants']) for f in families)} variants, {count} printed values, pp. 238-239")


main()
