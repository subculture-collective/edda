#!/usr/bin/env python3
"""Compare issue #233 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
# pp. 225, 226 and 227 name the families; the Helm of Brilliance continues from
# p. 225 to p. 226 and the Ioun Stone from p. 227 to p. 228.
PAGES = (225, 226, 227, 228)
RARITY = {"uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# Printed category and rarity line of each family, in catalogue order.
HEADER = {
    "Helm of Brilliance": "Wondrous Item, Very Rare (Requires Attunement)",
    "Helm of Comprehending Languages": "Wondrous Item, Uncommon",
    "Helm of Telepathy": "Wondrous Item, Uncommon (Requires Attunement)",
    "Helm of Teleportation": "Wondrous Item, Rare (Requires Attunement)",
    "Holy Avenger": "Weapon (Any Simple or Martial), Legendary (Requires Attunement by a Paladin)",
    "Horn of Blasting": "Wondrous Item, Rare",
    "Horn of Valhalla": "Wondrous Item, Rare (Silver or Brass), Very Rare (Bronze), or Legendary (Iron)",
    "Horseshoes of Speed": "Wondrous Item, Rare",
    "Horseshoes of a Zephyr": "Wondrous Item, Very Rare",
    "Immovable Rod": "Rod, Uncommon",
    "Instant Fortress": "Wondrous Item, Rare (Requires Attunement)",
    "Ioun Stone": "Wondrous Item, Rarity Varies (Requires Attunement)",
}
# The next heading after the last family on p. 228's Ioun Stone text.
FOLLOWING = "Iron Bands Wondrous Item"


def text(pdf, page, x):
    args = ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x), "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"]
    return subprocess.run(args, check=True, capture_output=True, text=True).stdout


def clean(raw):
    return FOOTER.sub(" ", " ".join(raw.split()))


def flat(value):
    """Ignore case, spaces and hyphens so a word split across a line end still matches."""
    return re.sub(r"[\s\-]", "", value).lower()


def has(source, phrase):
    return flat(phrase) in flat(source)


def tied(values, phrases, source, label):
    """Require the catalogued keys to be exactly the mapped ones, each with its printed phrase."""
    assert list(values) == list(phrases), f"{label}: {values}"
    for value, phrase in phrases.items():
        assert has(source, phrase), f"{label} {value}: {phrase}"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_233.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    pieces = [(p, clean(text(pdf, p, x))) for p in PAGES for x in (0, 306)]
    stream = " ".join(t for _, t in pieces)
    page_starts, at = [], 0
    for p, t in pieces:
        page_starts.append((at, p))
        at += len(t) + 1

    def page_of(index):
        return [p for start, p in page_starts if start <= index][-1]

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_233.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 233]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 233 and catalogue["batch"] == "data-magic-items-09"
    assert len(families) == 12 and sum(len(f["variants"]) for f in families) == 28
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert [f["presentations"][0]["pages"][0] for f in families] == [225, 226, 226, 226, 226, 226, 226, 227, 226, 227, 227, 227]
    assert [f["name"] for f in families] == list(HEADER)
    starts = {}
    for f in families:
        printed = f"{f['name']} {HEADER[f['name']]}"
        assert stream.count(printed) == 1, printed
        starts[f["name"]] = stream.index(printed)
    # Printed order differs from the catalogue order: Horseshoes of a Zephyr is
    # printed before Horseshoes of Speed. Blocks end at the next printed heading.
    order = sorted(starts.values()) + [stream.index(FOLLOWING)]
    block = {n: stream[s:order[order.index(s) + 1]] for n, s in starts.items()}
    assert starts["Horseshoes of a Zephyr"] < starts["Horseshoes of Speed"]
    assert [page_of(starts[f["name"]]) for f in families] == [f["presentations"][0]["pages"][0] for f in families], "presentation page disagrees with print"
    detail = {f["name"]: [v["details"] for v in f["variants"]] for f in families}
    charges = {f["name"]: f["variants"][0]["charges"] for f in families}
    none = {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}
    for f in families:
        page = f["presentations"][0]["pages"]
        for v in f["variants"]:
            assert v["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}, f["name"]
            assert v["attunement"] == {"mode": "required" if "Requires Attunement" in HEADER[f["name"]] else "none", "restrictions": []}, v["name"]
        if len(f["variants"]) == 1:
            assert RARITY[f["variants"][0]["rarity"]] in HEADER[f["name"]], f["name"]
    assert " GP" not in " ".join(block[n] for n in ("Helm of Brilliance", "Helm of Telepathy", "Helm of Teleportation", "Holy Avenger", "Horn of Blasting", "Horn of Valhalla", "Instant Fortress", "Ioun Stone"))

    # Helm of Brilliance (pp. 225-226).
    d = detail["Helm of Brilliance"][0]["helm_of_brilliance"]
    b = block["Helm of Brilliance"]
    assert d["source_pages"] == [225, 226] == sorted({page_of(starts["Helm of Brilliance"]), page_of(starts["Helm of Brilliance"] + len(b) - 1)})
    assert has(b, "This helm is set with 1d10 diamonds, 2d10 rubies, 3d10 fire opals, and 4d10 opals")
    assert d["gems"] == [{"gem": "diamond", "dice": 1, "sides": 10}, {"gem": "ruby", "dice": 2, "sides": 10}, {"gem": "fire-opal", "dice": 3, "sides": 10}, {"gem": "opal", "dice": 4, "sides": 10}]
    assert d["gems_crumble_when_pried"] is True and has(b, "Any gem pried from the helm crumbles to dust")
    assert d["loses_magic_when_all_gems_gone"] is True and has(b, "When all the gems are removed or destroyed, the helm loses its magic")
    tied([d["condition"]], {"while-wearing": "You gain the following benefits while wearing the helm"}, b, "helm condition")
    dl = d["diamond_light"]
    tied([dl["requires"]], {"at-least-one-diamond": "As long as it has at least one diamond"}, b, "diamond requires")
    assert dl["emanation_feet"] == 30 and has(b, "the helm emits a 30-foot Emanation")
    tied([dl["dim_light_when"]], {"undead-within-emanation": "When at least one Undead is within that area, the Emanation is filled with Dim Light"}, b, "diamond dim")
    tied([dl["damage_trigger"]], {"undead-starts-turn-in-emanation": "Any Undead that starts its turn in that area takes 1d6 Radiant damage"}, b, "diamond trigger")
    assert dl["damage"] == {"dice": 1, "sides": 6, "type": "radiant"}
    fo = d["fire_opal_flames"]
    tied([fo["requires"]], {"at-least-one-fire-opal": "As long as the helm has at least one fire opal"}, b, "flames requires")
    assert fo["action"] == "magic" and has(b, "you can take a Magic action to cause one weapon you are holding to burst into flames")
    tied([fo["target"]], {"one-weapon-you-are-holding": "one weapon you are holding to burst into flames"}, b, "flames target")
    assert (fo["bright_light_feet"], fo["dim_light_additional_feet"]) == (10, 10) and has(b, "Bright Light in a 10-foot radius and Dim Light for an additional 10 feet")
    tied(fo["harmless_to"], {"you": "The flames are harmless to you", "the-weapon": "harmless to you and the weapon"}, b, "flames harmless")
    assert fo["extra_damage"] == {"dice": 1, "sides": 6, "type": "fire"}
    tied([fo["extra_damage_trigger"]], {"hit-with-attack-using-blazing-weapon": "When you hit with an attack using the blazing weapon, the target takes an extra 1d6 Fire damage"}, b, "flames trigger")
    tied(fo["ends_on"], {"bonus-action-extinguish": "until you take a Bonus Action to extinguish them", "drop-weapon": "or until you drop", "stow-weapon": "drop or stow the weapon"}, b, "flames end")
    rr = d["ruby_resistance"]
    tied([rr["requires"]], {"at-least-one-ruby": "As long as the helm has at least one ruby"}, b, "ruby requires")
    assert rr["resistance"] == "fire" and has(b, "you have Resistance to Fire damage")
    sp = d["spells"]
    assert sp["save_dc"] == 18 and has(b, "You can cast one of the following spells (save DC 18)")
    tied([sp["component"]], {"one-helm-gem-of-specified-type": "using one of the helm’s gems of the specified type as a component"}, b, "helm component")
    assert has(b, "Daylight (opal), Fireball (fire opal), Prismatic Spray (diamond), or Wall of Fire (ruby)")
    assert sp["spells"] == [{"spell": "daylight", "gem": "opal"}, {"spell": "fireball", "gem": "fire-opal"}, {"spell": "prismatic-spray", "gem": "diamond"}, {"spell": "wall-of-fire", "gem": "ruby"}]
    assert sp["gem_destroyed_when_cast"] is True and sp["gem_disappears_from_helm"] is True and has(b, "The gem is destroyed when the spell is cast and disappears from the helm")
    fd = d["taking_fire_damage"]
    tied([fd["trigger"]], {"wearing-and-fire-damage-from-failed-save-against-spell": "Roll 1d20 if you are wearing the helm and take Fire damage as a result of failing a saving throw against a spell"}, b, "mishap trigger")
    assert (fd["roll_die_sides"], fd["destroyed_on_roll"]) == (20, 1) and has(b, "On a roll of 1, the helm emits beams of light from its remaining gems and is then destroyed")
    tied([fd["effect"]], {"beams-of-light-from-remaining-gems": "the helm emits beams of light from its remaining gems"}, b, "mishap effect")
    assert fd["helm_destroyed"] is True
    assert fd["emanation_feet"] == 60 and fd["origin"] == "you" and has(b, "Each creature within a 60-foot Emanation originating from you")
    assert (fd["save_ability"], fd["save_dc"]) == ("dexterity", 17) and has(b, "must succeed on a DC 17 Dexterity saving throw or be struck by a beam")
    assert fd["damage_type"] == "radiant"
    tied([fd["damage"]], {"number-of-gems-in-helm": "taking Radiant damage equal to the number of gems in the helm"}, b, "mishap damage")
    assert charges["Helm of Brilliance"] == none

    d = detail["Helm of Comprehending Languages"][0]["helm_of_comprehending_languages"]
    b = block["Helm of Comprehending Languages"]
    tied([d["condition"]], {"while-wearing": "While wearing this helm"}, b, "languages condition")
    tied([d["spell"]], {"comprehend-languages": "you can cast Comprehend Languages from it"}, b, "languages spell")
    assert charges["Helm of Comprehending Languages"] == none

    d = detail["Helm of Telepathy"][0]["helm_of_telepathy"]
    b = block["Helm of Telepathy"]
    tied([d["condition"]], {"while-wearing": "While wearing this helm"}, b, "telepathy condition")
    assert d["telepathy_range_feet"] == 30 and has(b, "you have telepathy with a range of 30 feet")
    tied(d["spells"], {"detect-thoughts": "you can cast Detect Thoughts or Suggestion", "suggestion": "Detect Thoughts or Suggestion (save DC 13)"}, b, "telepathy spells")
    assert (d["save_dc_spell"], d["save_dc"]) == ("suggestion", 13) and has(b, "Suggestion (save DC 13) from the helm")
    tied([d["spell_reuse"]], {"that-spell-cannot-be-cast-again": "that spell can’t be cast from it again"}, b, "telepathy reuse")
    tied([d["spell_reuse_window"]], {"until-next-dawn": "until the next dawn"}, b, "telepathy window")
    assert has(b, "Once either spell is cast from the helm") and charges["Helm of Telepathy"] == none

    d = detail["Helm of Teleportation"][0]["helm_of_teleportation"]
    b = block["Helm of Teleportation"]
    assert charges["Helm of Teleportation"] == {"mode": "limited", "maximum": 3, "per_use": 1, "recharge": "regain-1d3-at-dawn"} and d["max_charges"] == 3 and has(b, "This helm has 3 charges")
    tied([d["condition"]], {"while-wearing": "While wearing it"}, b, "teleportation condition")
    assert d["charge_cost"] == 1 and has(b, "you can expend 1 charge to cast Teleport from it")
    tied([d["spell"]], {"teleport": "cast Teleport from it"}, b, "teleportation spell")
    assert d["regain_dice"] == {"dice": 1, "sides": 3, "type": "charges"} and has(b, "The helm regains 1d3 expended charges daily at dawn")
    tied([d["regain_timing"]], {"dawn": "daily at dawn"}, b, "teleportation timing")

    d = detail["Holy Avenger"][0]["holy_avenger"]
    b = block["Holy Avenger"]
    tied(d["attunement_by_any_of"], {"paladin": "Requires Attunement by a Paladin"}, b, "avenger attunement")
    assert (d["attack_bonus"], d["damage_bonus"]) == (3, 3) and has(b, "You gain a +3 bonus to attack rolls and damage rolls made with this magic weapon")
    assert d["extra_damage"] == {"dice": 2, "sides": 10, "type": "radiant"} and has(b, "that creature takes an extra 2d10 Radiant damage")
    tied(d["extra_damage_against"], {"fiend": "When you hit a Fiend", "undead": "a Fiend or an Undead with it"}, b, "avenger targets")
    tied([d["extra_damage_trigger"]], {"hit-with-weapon": "When you hit a Fiend or an Undead with it"}, b, "avenger trigger")
    a = d["aura"]
    tied([a["condition"]], {"while-holding-drawn-weapon": "While you hold the drawn weapon"}, b, "avenger condition")
    assert a["emanation_feet"] == 10 and a["origin"] == "you" and has(b, "it creates a 10-foot Emanation originating from you")
    tied(a["affected"], {"you": "You and all creatures Friendly to you in the Emanation", "creatures-friendly-to-you-in-emanation": "all creatures Friendly to you in the Emanation"}, b, "avenger affected")
    tied([a["benefit"]], {"advantage-on-saving-throws": "have Advantage on saving throws"}, b, "avenger benefit")
    tied(a["against"], {"spells": "saving throws against spells", "other-magical-effects": "against spells and other magical effects"}, b, "avenger against")
    assert (a["larger_emanation_min_paladin_levels"], a["larger_emanation_feet"]) == (17, 30)
    assert has(b, "If you have 17 or more levels in the Paladin class, the size of the Emanation increases to 30 feet")
    assert charges["Holy Avenger"] == none

    d = detail["Horn of Blasting"][0]["horn_of_blasting"]
    b = block["Horn of Blasting"]
    assert d["action"] == "magic" and has(b, "You can take a Magic action to blow the horn")
    assert (d["cone_feet"], d["audible_feet"]) == (30, 600) and has(b, "thunderous blast in a 30-foot Cone that is audible out to 600 feet")
    assert (d["save_ability"], d["save_dc"]) == ("constitution", 15) and has(b, "Each creature in the Cone makes a DC 15 Constitution saving throw")
    assert d["damage"] == {"dice": 5, "sides": 8, "type": "thunder"} and has(b, "On a failed save, a creature takes 5d8 Thunder damage")
    assert d["condition"] == "deafened" and d["condition_duration_seconds"] == 60 and has(b, "has the Deafened condition for 1 minute")
    tied([d["success_result"]], {"half-damage-only": "On a successful save, a creature takes half as much damage only"}, b, "blasting success")
    tied([d["object_targets"]], {"glass-or-crystal-in-cone-not-worn-or-carried": "Glass or crystal objects in the Cone that aren’t being worn or carried"}, b, "blasting objects")
    assert d["object_damage"] == {"dice": 10, "sides": 8, "type": "thunder"} and has(b, "take 10d8 Thunder damage")
    assert d["explosion_chance_percent"] == 20 and has(b, "Each use of the horn’s magic has a 20 percent chance of causing the horn to explode")
    assert d["explosion_damage"] == {"dice": 10, "sides": 6, "type": "force"} and has(b, "The explosion deals 10d6 Force damage to the user")
    tied([d["explosion_target"]], {"user": "damage to the user"}, b, "blasting target")
    assert d["explosion_destroys_horn"] is True and has(b, "and destroys the horn")
    assert charges["Horn of Blasting"] == none

    # Horn of Valhalla: one variant per printed metal, rarity printed in the heading.
    b = block["Horn of Valhalla"]
    f = families[6]
    assert [(v["id"].rsplit("/", 1)[1], v["name"], v["rarity"]) for v in f["variants"]] == [("silver", "Silver", "rare"), ("brass", "Brass", "rare"), ("bronze", "Bronze", "very-rare"), ("iron", "Iron", "legendary")]
    assert has(b, "Rare (Silver or Brass), Very Rare (Bronze), or Legendary (Iron)") and has(b, "Four types of Horn of Valhalla are known to exist, each made of a different metal")
    table = re.search(r"1d100 Horn Type Spirits Requirement (.*?)\s*$", b).group(1)
    rows = re.findall(r"(\d\d)–(\d\d) (Silver|Brass|Bronze|Iron) (\d) (.*?)(?= \d\d–\d\d |$)", table)
    assert len(rows) == 4, rows
    for (lo, hi, metal, spirits, requirement), t in zip(rows, detail["Horn of Valhalla"]):
        t = t["horn_of_valhalla"]
        hi = 100 if hi == "00" else int(hi)
        assert (t["table_die_sides"], t["table_min"], t["table_max"], t["spirits"]) == (100, int(lo), hi, int(spirits)), metal
        assert t["type"] == metal.lower()
        assert flat(t["requirement"]) == flat(requirement), (metal, requirement)
        assert has(b, f"{metal}"), metal
    assert [t["horn_of_valhalla"]["table_max"] for t in detail["Horn of Valhalla"]] == [40, 75, 90, 100] and "91–00 Iron" in b
    for t in detail["Horn of Valhalla"]:
        t = t["horn_of_valhalla"]
        assert t["action"] == "magic" and has(b, "You can take a Magic action to blow this horn")
        assert t["spawn_range_feet"] == 60 and has(b, "appear in unoccupied spaces within 60 feet of you")
        tied([t["spirits_from_plane"]], {"ysgard": "warrior spirits from the plane of Ysgard"}, b, "valhalla plane")
        tied([t["spirit_stat_block"]], {"berserker": "Each spirit uses the Berserker stat block"}, b, "valhalla stat block")
        tied([t["spirit_appearance"]], {"living-breathing-warriors": "The spirits look like living, breathing warriors"}, b, "valhalla appearance")
        assert t["duration_hours"] == 1 and t["cooldown_days"] == 7 and t["types_known"] == 4
        tied(t["return_triggers"], {"duration-elapsed": "returns to Ysgard after 1 hour", "drops-to-0-hit-points": "or when it drops to 0 Hit Points"}, b, "valhalla returns")
        tied(t["spirit_immunities"], {"charmed": "Immunity to the Charmed", "frightened": "Charmed and Frightened conditions"}, b, "valhalla immunities")
        assert has(b, "it can’t be used again until 7 days have passed")
        tied([t["unmet_requirement_result"]], {"spirits-attack-you": "If you blow the horn without meeting its requirement, the summoned spirits attack you"}, b, "valhalla unmet")
        tied(t["met_requirement_result"], {"friendly-to-you-and-your-allies": "they are Friendly to you and your allies", "follow-your-commands": "and follow your commands"}, b, "valhalla met")
    assert charges["Horn of Valhalla"] == none

    for name in ("Horseshoes of Speed", "Horseshoes of a Zephyr"):
        t = detail[name][0]["horseshoes_of_speed" if name.endswith("Speed") else "horseshoes_of_a_zephyr"]
        b = block[name]
        assert t["set_size"] == 4 and has(b, "These horseshoes come in a set of four")
        assert t["affix_action"] == "magic" and has(b, "As a Magic action, you can touch one of the horseshoes to the hoof")
        tied([t["affix_target"]], {"hoof-of-horse-or-similar-creature": "to the hoof of a horse or similar creature"}, b, name + " target")
        assert t["remove_action"] == "magic" and has(b, "Removing a horseshoe also takes a Magic action")
        assert charges[name] == none
    d = detail["Horseshoes of Speed"][0]["horseshoes_of_speed"]
    b = block["Horseshoes of Speed"]
    assert d["requires_all_affixed_to_same_creature"] is True and has(b, "While all four horseshoes are attached to the same creature")
    assert d["speed_bonus_feet"] == 30 and has(b, "its Speed is increased by 30 feet")
    d = detail["Horseshoes of a Zephyr"][0]["horseshoes_of_a_zephyr"]
    b = block["Horseshoes of a Zephyr"]
    assert d["requires_all_affixed"] is True and has(b, "While all four shoes are affixed to the hooves of a horse or similar creature")
    assert d["float_inches"] == 4 and d["moves_normally"] is True and has(b, "allow the creature to move normally while floating 4 inches above a surface")
    tied([d["cross_surfaces"]], {"nonsolid-or-unstable": "cross or stand above nonsolid or unstable surfaces"}, b, "zephyr surfaces")
    tied(d["surface_examples"], {"water": "such as water or lava", "lava": "water or lava"}, b, "zephyr examples")
    assert d["leaves_tracks"] is False and has(b, "The creature leaves no tracks")
    assert d["ignores_difficult_terrain"] is True and has(b, "ignores Difficult Terrain")
    assert d["travel_hours_per_day"] == 12 and has(b, "can travel for up to 12 hours a day")
    tied([d["no_exhaustion_from"]], {"extended-travel": "without gaining Exhaustion levels from extended travel"}, b, "zephyr exhaustion")

    d = detail["Immovable Rod"][0]["immovable_rod"]
    b = block["Immovable Rod"]
    tied([d["material"]], {"iron": "This iron rod"}, b, "rod material")
    tied([d["button_on"]], {"one-end": "has a button on one end"}, b, "rod button")
    assert d["button_action"] == "utilize" and has(b, "You can take a Utilize action to press the button")
    tied([d["fixed_effect"]], {"magically-fixed-in-place": "causes the rod to become magically fixed in place"}, b, "rod fixed")
    assert d["holds_against_gravity"] is True and has(b, "the rod doesn’t move, even if it defies gravity")
    tied(d["released_by"], {"you": "Until you or another creature takes a Utilize action to push the button again", "another-creature": "you or another creature takes a Utilize action"}, b, "rod release")
    assert d["max_weight_pounds"] == 8000 and has(b, "The rod can hold up to 8,000 pounds of weight")
    tied([d["overload_result"]], {"deactivates-and-falls": "More weight causes the rod to deactivate and fall"}, b, "rod overload")
    assert d["move_action"] == "utilize" and d["move_check"] == {"ability": "strength", "skill": "athletics", "dc": 30}
    assert has(b, "A creature can take a Utilize action to make a DC 30 Strength (Athletics) check")
    assert d["move_feet"] == 10 and has(b, "moving the fixed rod up to 10 feet on a successful check")
    tied([d["move_on"]], {"success": "on a successful check"}, b, "rod success")
    assert charges["Immovable Rod"] == none

    d = detail["Instant Fortress"][0]["instant_fortress"]
    b = block["Instant Fortress"]
    assert d["statuette_inches"] == 1 and has(b, "place this 1-inch adamantine statuette")
    tied([d["statuette_material"]], {"adamantine": "1-inch adamantine statuette"}, b, "fortress material")
    assert d["action"] == "magic" and has(b, "As a Magic action, you can place this 1-inch adamantine statuette on the ground")
    tied([d["placed_on"]], {"ground": "statuette on the ground"}, b, "fortress placed")
    tied([d["activation"]], {"command-word": "using a command word, cause it to grow rapidly into a square adamantine tower"}, b, "fortress activation")
    tied([d["revert_activation"]], {"repeat-command-word": "Repeating the command word causes the tower to revert to statuette form"}, b, "fortress revert")
    tied([d["revert_condition"]], {"tower-is-empty": "which works only if the tower is empty"}, b, "fortress condition")
    tied([d["displaced_creatures"]], {"pushed-to-unoccupied-space-outside-next-to-tower": "is pushed to an unoccupied space outside but next to the tower"}, b, "fortress creatures")
    tied([d["displaced_objects"]], {"unworn-uncarried-objects-pushed-clear": "Objects in the area that aren’t being worn or carried are also pushed clear of the tower"}, b, "fortress objects")
    t = d["tower"]
    tied([t["shape"]], {"square": "square adamantine tower"}, b, "tower shape")
    assert (t["side_feet"], t["height_feet"]) == (20, 30) and has(b, "The tower is 20 feet on a side and 30 feet high")
    tied(t["features"], {"arrow-slits-on-all-sides": "with arrow slits on all sides", "battlement-atop": "and a battlement atop it"}, b, "tower features")
    assert t["floors"] == 2 and has(b, "Its interior is divided into two floors")
    tied(t["floor_connection_your_choice"], {"ladder": "a ladder, staircase, or ramp (your choice) connecting them", "staircase": "ladder, staircase, or ramp", "ramp": "or ramp (your choice)"}, b, "tower connection")
    tied([t["roof_access"]], {"trapdoor": "ends at a trapdoor leading to the roof"}, b, "tower roof")
    assert t["doors"] == 1 and has(b, "the tower has a single door at ground level on the side facing you")
    tied([t["door_location"]], {"ground-level-facing-you": "a single door at ground level on the side facing you"}, b, "tower door")
    tied([t["door_opens_by"]], {"your-command": "The door opens only at your command"}, b, "tower door opens")
    assert t["door_open_action"] == "bonus" and has(b, "which you can issue as a Bonus Action")
    tied([t["door_immune_to"]], {"knock-and-similar-magic": "It is immune to the Knock spell and similar magic"}, b, "tower door immune")
    assert t["cannot_be_tipped"] is True and has(b, "Magic prevents the tower from being tipped over")
    tied(t["durable_parts"], {"roof": "The roof, the door, and the walls each have AC 20", "door": "the door, and the walls", "walls": "and the walls each have"}, b, "tower parts")
    assert (t["armor_class"], t["hit_points"]) == (20, 100) and has(b, "each have AC 20; HP 100")
    tied(t["immune_to"], {"bludgeoning": "Immunity to Bludgeoning, Piercing, and Slashing damage", "piercing": "Bludgeoning, Piercing, and Slashing", "slashing": "and Slashing damage"}, b, "tower immunities")
    tied([t["immunity_exception"]], {"damage-dealt-by-siege-equipment": "except that which is dealt by siege equipment"}, b, "tower exception")
    tied([t["resistant_to"]], {"all-other-damage": "Resistance to all other damage"}, b, "tower resistance")
    r = d["damage_repair"]
    assert r["shrinking_repairs_damage"] is False and has(b, "Shrinking the tower back down to statuette form doesn’t repair damage to the tower")
    tied([r["only_spell"]], {"wish": "Only a Wish spell can repair the tower"}, b, "fortress repair")
    assert r["counts_as_spell_level_at_most"] == 8 and has(b, "this use of the spell counts as replicating a spell of level 8 or lower")
    tied([r["each_casting_restores"]], {"all-hit-points": "Each casting of Wish causes the tower to regain all its Hit Points"}, b, "fortress restore")
    assert charges["Instant Fortress"] == none

    # Ioun Stone: shared rules, then one variant per printed stone type with its rarity.
    b = block["Ioun Stone"]
    f = families[11]
    shared = f["variants"][0]["details"]["ioun_stone"]["shared_rules"]
    assert all(v["details"]["ioun_stone"]["shared_rules"] == shared for v in f["variants"])
    tied([shared["size"]], {"roughly-marble": "Roughly marble sized"}, b, "ioun size")
    assert shared["activation_action"] == "magic" and has(b, "When you take a Magic action to toss an Ioun Stone into the air")
    tied([shared["method"]], {"toss-into-the-air": "to toss an Ioun Stone into the air"}, b, "ioun method")
    assert (shared["orbit_distance_dice"], shared["orbit_distance_sides"]) == (1, 3) and has(b, "the stone orbits your head at a distance of 1d3 feet")
    assert shared["max_orbiting"] == 3 and has(b, "You can have up to three Ioun Stones orbiting your head at the same time")
    tied([shared["orbiting_stone_counts_as"]], {"object-you-are-wearing": "is considered to be an object you are wearing"}, b, "ioun worn")
    tied(shared["orbit_traits"], {"avoids-contact-with-creatures-and-objects": "The orbiting stone avoids contact with other creatures and objects",
                                  "adjusts-orbit-to-avoid-collisions": "adjusting its orbit to avoid collisions", "thwarts-attacks-and-snatching": "thwarting all attempts by other creatures to attack or snatch it"}, b, "ioun orbit")
    assert shared["stow_action"] == "utilize" and has(b, "As a Utilize action, you can seize and stow any number of Ioun Stones orbiting your head")
    tied([shared["stow_count"]], {"any-number-orbiting": "any number of Ioun Stones orbiting your head"}, b, "ioun stow")
    tied([shared["attunement_ends_result"]], {"stone-falls-as-if-dropped": "the stone falls as though you had dropped it"}, b, "ioun attunement")
    assert has(b, "The type of stone determines its rarity and effects")
    names = [v["name"] for v in f["variants"]]
    assert names == ["Absorption", "Agility", "Awareness", "Fortitude", "Greater Absorption", "Insight", "Intellect", "Leadership", "Mastery", "Protection", "Regeneration", "Reserve", "Strength", "Sustenance"]
    sect, cursor = {}, 0
    for v in f["variants"]:
        heading = f"{v['name']} ({RARITY[v['rarity']]})."
        at = b.index(heading, cursor)
        assert b.count(heading) == 1, heading
        sect[v["name"]] = (at, None)
        cursor = at + 1
    keys = list(sect)
    for i, name in enumerate(keys):
        end = sect[keys[i + 1]][0] if i + 1 < len(keys) else len(b)
        sect[name] = (sect[name][0], end)
    texts = {n: b[s:e] for n, (s, e) in sect.items()}
    base = starts["Ioun Stone"]
    for v in f["variants"]:
        t = v["details"]["ioun_stone"]
        s = texts[v["name"]]
        assert t["source_pages"] == [page_of(base + sect[v["name"]][0])], v["name"]
        assert has(s, f"this {t['color']} {t['shape']}"), v["name"]
        facts = {x["key"]: (x["amount"], x["unit"]) for x in t["facts"]}
        assert len(facts) == len(t["facts"]) and len(set(t["traits"])) == len(t["traits"])
        assert has(s, "orbits your head") or v["name"] == "Reserve"

    def stone(name, benefit, action, ability, expect_facts, phrases):
        t = next(v["details"]["ioun_stone"] for v in f["variants"] if v["name"] == name)
        assert (t["benefit"], t["action"], t["ability"]) == (benefit, action, ability), name
        assert {k: (a, u) for k, (a, u) in ((x["key"], (x["amount"], x["unit"])) for x in t["facts"])} == expect_facts, name
        for phrase in phrases:
            assert has(texts[name], phrase), (name, phrase)
        return t

    for name, level in (("Absorption", 4), ("Greater Absorption", 8)):
        t = stone(name, "cancel-spell", "reaction", "none", {"max-spell-level": (level, "level"), "burns-out-after-levels": (20, "levels")},
                  [f"you can take a Reaction to cancel a spell of level {level} or lower cast by a creature you can see", "Once the stone has canceled 20 levels of spells, it burns out, turns dull gray, and loses its magic"])
        tied(t["traits"], {"target-spell-cast-by-creature-you-can-see": "cast by a creature you can see", "canceled-spell-has-no-effect": "A canceled spell has no effect",
                           "resources-used-to-cast-are-wasted": "any resources used to cast it are wasted", "burns-out-turns-dull-gray-loses-magic": "it burns out, turns dull gray, and loses its magic"}, texts[name], name + " traits")
    for name, ability in (("Agility", "dexterity"), ("Fortitude", "constitution"), ("Insight", "wisdom"), ("Intellect", "intelligence"), ("Leadership", "charisma"), ("Strength", "strength")):
        t = stone(name, "ability-score-increase", "none", ability, {"increase": (2, "points"), "maximum-score": (20, "score")},
                  [f"Your {ability.capitalize()} increases by 2, to a maximum of 20, while this"])
        assert t["traits"] == ["while-orbiting-head"]
    t = stone("Awareness", "advantage-on-rolls", "none", "none", {}, ["you have Advantage on Initiative rolls and Wisdom (Perception) checks"])
    tied(t["traits"], {"advantage-on-initiative-rolls": "Advantage on Initiative rolls", "advantage-on-wisdom-perception-checks": "and Wisdom (Perception) checks", "while-orbiting-head": "While this dark-blue rhomboid orbits your head"}, texts["Awareness"], "awareness traits")
    t = stone("Mastery", "proficiency-bonus-increase", "none", "none", {"increase": (1, "points")}, ["Your Proficiency Bonus increases by 1"])
    assert t["traits"] == ["while-orbiting-head"]
    t = stone("Protection", "armor-class-bonus", "none", "none", {"bonus": (1, "points")}, ["You gain a +1 bonus to Armor Class"])
    assert t["traits"] == ["while-orbiting-head"]
    t = stone("Regeneration", "regain-hit-points", "none", "none", {"regain": (15, "hit-points"), "minimum-hit-points": (1, "hit-points")},
              ["You regain 15 Hit Points at the end of each hour this pearly white spindle orbits your head if you have at least 1 Hit Point"])
    tied(t["traits"], {"end-of-each-hour": "at the end of each hour", "while-orbiting-head": "this pearly white spindle orbits your head"}, texts["Regeneration"], "regeneration traits")
    t = stone("Reserve", "store-spells", "none", "none",
              {"capacity": (4, "levels"), "initial-levels-dice": (1, "dice"), "initial-levels-sides": (4, "sides"), "cast-into-min-level": (1, "level"), "cast-into-max-level": (4, "level")},
              ["The stone can store up to 4 levels of spells at a time", "it contains 1d4 levels of stored spells chosen by the GM", "Any creature can cast a spell of level 1 through 4 into the stone by touching it as the spell is cast"])
    tied(t["traits"], {"initial-spells-chosen-by-gm": "chosen by the GM", "any-creature-can-cast-into-stone-by-touching-it": "by touching it as the spell is cast",
                       "stored-spell-has-no-other-effect": "The spell has no effect, other than to be stored in the stone",
                       "spell-expended-without-effect-if-stone-cannot-hold-it": "If the stone can’t hold the spell, the spell is expended without effect",
                       "slot-level-determines-space-used": "The level of the slot used to cast the spell determines how much space it uses",
                       "you-can-cast-stored-spell-while-orbiting": "While this stone orbits your head, you can cast any spell stored in it",
                       "uses-original-caster-slot-level-save-dc-attack-bonus-and-spellcasting-ability": "The spell uses the slot level, spell save DC, spell attack bonus, and spellcasting ability of the original caster",
                       "cast-spell-is-no-longer-stored": "The spell cast from the stone is no longer stored in it, freeing up space"}, texts["Reserve"], "reserve traits")
    t = stone("Sustenance", "no-need-to-eat-or-drink", "none", "none", {}, ["You don’t need to eat or drink while this clear spindle orbits your head"])
    assert t["traits"] == ["while-orbiting-head"]
    assert charges["Ioun Stone"] == none and all(v["charges"] == none for v in f["variants"])
    # Every other family lists no charge pool.
    for name in HEADER:
        if name != "Helm of Teleportation":
            assert charges[name] == none, name
    print("#233 pinned PDF, twelve identities/pages/headers, 28 variants, the 1d100 horn table, charges and source values pass")


if __name__ == "__main__":
    main()
