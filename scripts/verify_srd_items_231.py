#!/usr/bin/env python3
"""Compare issue #231 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
# pp. 220, 221 and 223 name the families; the Eversmoking Bottle continues on p. 221
# and the Figurine of Wondrous Power on p. 222.
PAGES = (220, 221, 222, 223)
RARITY = {"uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# Printed category line of each family, before the rarity.
CATEGORY = {
    "energy-bow": "Weapon (Longbow or Shortbow)",
    "flame-tongue": "Weapon (Any Melee Weapon)",
    "frost-brand": "Weapon (Glaive, Greatsword, Longsword, Rapier, Scimitar, or Shortsword)",
}
# The next heading after the last family on p. 223.
FOLLOWING = "Gem of Seeing Wondrous Item"


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


def tied(values, phrases, source, label):
    """Require the catalogued keys to be exactly the mapped ones, each with its printed phrase."""
    assert list(values) == list(phrases), f"{label}: {values}"
    for value, phrase in phrases.items():
        assert has(source, phrase), f"{label} {value}: {phrase}"


def header(family):
    variants = family["variants"]
    attune = variants[0]["attunement"]["mode"] == "required"
    suffix = " (Requires Attunement)" if attune else ""
    if family["name"] in ("Feather Token", "Figurine of Wondrous Power"):
        return "Wondrous Item, Rarity Varies"
    category = CATEGORY.get(family["id"].rsplit("/", 1)[1], "Wondrous Item")
    assert (family["category"] == "weapon") == category.startswith("Weapon")
    return f"{category}, {RARITY[variants[0]['rarity']]}{suffix}"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_231.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream = " ".join(clean(text(pdf, p, x)) for p in PAGES for x in (0, 306))

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_231.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 231]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 231 and catalogue["batch"] == "data-magic-items-07"
    assert len(families) == 12 and sum(len(f["variants"]) for f in families) == 25
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert [f["presentations"][0]["pages"][0] for f in families] == [220, 220, 221, 221, 221, 221, 221, 223, 223, 223, 223, 223]
    names = [f["name"] for f in families]
    starts = []
    for f in families:
        printed = f"{f['name']} {header(f)}"
        assert stream.count(printed) == 1, printed
        starts.append(stream.index(printed))
    assert starts == sorted(starts), "printed order differs from the catalogue order"
    ends = starts[1:] + [stream.index(FOLLOWING)]
    block = {f["name"]: stream[s:e] for f, s, e in zip(families, starts, ends)}
    detail = {f["name"]: [v["details"] for v in f["variants"]] for f in families}
    charges = {f["name"]: f["variants"][0]["charges"] for f in families}
    none = {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}
    for f in families:
        page = f["presentations"][0]["pages"]
        for v in f["variants"]:
            assert v["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}, f["name"]
    assert " GP" not in " ".join(b for n, b in block.items() if n != "Gem of Brightness")

    d = detail["Energy Bow"][0]["energy_bow"]
    b = block["Energy Bow"]
    assert (d["attack_bonus"], d["damage_bonus"]) == (1, 1) and has(b, "You gain a +1 bonus to attack rolls and damage rolls made with this magic weapon")
    assert d["has_string"] is False and has(b, "which has no string")
    tied([d["arrow_appears"]], {"when-arm-pulled-back-in-firing-motion": "Each time you pull your arm back in a firing motion"}, b, "bow appears")
    assert (d["arrow_damage_type"], d["replaces_damage_type"]) == ("force", "piercing") and has(b, "deals Force damage instead of Piercing damage on a hit")
    tied([d["arrow_ends"]], {"disappears-after-hit-or-miss": "it disappears after it hits or misses its target"}, b, "bow end")
    assert (d["bright_light_feet"], d["dim_light_additional_feet"]) == (20, 20) and has(b, "Bright Light in a 20-foot radius and Dim Light for an additional 20 feet")
    r = d["arrow_of_restraint"]
    tied([r["trigger"]], {"ranged-attack-instead-of-damage": "make a ranged attack against a creature, you can try to restrain the target instead of dealing damage"}, b, "restraint trigger")
    assert (r["save_ability"], r["save_dc"], r["condition"], r["duration_seconds"]) == ("strength", 15, "restrained", 60)
    assert has(b, "succeed on a DC 15 Strength saving throw or have the Restrained condition for 1 minute")
    assert r["escape_action"] == "action" and r["escape_check"] == {"ability": "strength", "skill": "athletics", "dc": 20}
    assert has(b, "As an action, a creature Restrained by an arrow can make a DC 20 Strength (Athletics) check")
    tied([r["escape_ends"]], {"effect-on-itself": "ending the effect on itself on a successful check"}, b, "restraint escape")
    t = d["arrow_of_transport"]
    assert t["action"] == "magic" and (t["range_feet"], t["object_max_cube_feet"], t["destination_range_feet"]) == (60, 5, 10)
    assert has(b, "As a Magic action, you can fire one energy arrow from this weapon at a target you can see within 60 feet of yourself")
    tied(t["targets"], {"willing-medium-or-smaller-creature": "a willing Medium or smaller creature", "unworn-uncarried-object": "an object that isn’t being worn or carried"}, b, "transport targets")
    assert has(b, "small enough to fit inside a 5-foot Cube")
    tied([t["destination"]], {"unoccupied-space-you-can-see": "teleports the target to an unoccupied space you can see within 10 feet of you"}, b, "transport destination")
    ladder = d["energy_ladder"]
    assert ladder["action"] == "magic" and (ladder["wall_range_feet"], ladder["ladder_length_feet"], ladder["duration_seconds"]) == (60, 60, 60)
    assert has(b, "As a Magic action, you can loose a flurry of energy arrows from this weapon at a wall up to 60 feet away from yourself")
    assert has(b, "a magical ladder up to 60 feet long on the wall") and has(b, "This ladder lasts for 1 minute before disappearing")
    tied([ladder["ends_with"]], {"disappears": "before disappearing"}, b, "ladder end")
    assert charges["Energy Bow"] == none

    d = detail["Eversmoking Bottle"][0]["eversmoking_bottle"]
    b = block["Eversmoking Bottle"]
    assert d["source_pages"] == [220, 221] and "Closing the bottle causes the cloud" in b
    assert d["action"] == "magic" and has(b, "As a Magic action, you can open or close this bottle")
    tied(d["operations"], {"open": "open or close", "close": "or close this bottle"}, b, "bottle operations")
    assert (d["emanation_feet"], d["growth_feet_per_minute"], d["max_emanation_feet"]) == (60, 10, 120)
    assert has(b, "fills a 60-foot Emanation originating from the bottle") and has(b, "the size of the Emanation increases by 10 feet until it reaches its maximum size of 120 feet")
    tied([d["obscurement"]], {"heavily-obscured": "The area within the smoke is Heavily Obscured"}, b, "bottle obscurement")
    tied([d["closed_state"]], {"cloud-fixed-in-place": "the cloud to become fixed in place"}, b, "bottle closed")
    assert (d["closed_disperses_after_seconds"], d["strong_wind_disperses_after_seconds"]) == (600, 60)
    assert has(b, "disperses after 10 minutes") and has(b, "disperses the cloud after 1 minute")
    tied([d["strong_wind_example"]], {"gust-of-wind": "such as that created by the Gust of Wind spell"}, b, "bottle wind")
    assert charges["Eversmoking Bottle"] == none

    d = detail["Eyes of Charming"][0]["eyes_of_charming"]
    b = block["Eyes of Charming"]
    tied([d["appearance"]], {"crystal-lenses-over-eyes": "These crystal lenses fit over the eyes"}, b, "charming appearance")
    tied([d["condition"]], {"while-wearing": "While wearing them"}, b, "charming condition")
    assert charges["Eyes of Charming"] == {"mode": "limited", "maximum": 3, "per_use": 1, "recharge": "regain-all-at-dawn"} and d["max_charges"] == 3
    assert has(b, "They have 3 charges") and has(b, "you can expend 1 or more charges")
    tied([d["spell"]], {"charm-person": "cast Charm Person (save DC 13)"}, b, "charming spell")
    assert d["save_dc"] == 13 and (d["base_spell_level"], d["base_charges"], d["levels_per_additional_charge"]) == (1, 1, 1)
    assert has(b, "For 1 charge, you cast the level 1 version of the spell") and has(b, "increase the spell’s level by one for each additional charge you expend")
    tied([d["recharge"]], {"all-expended-charges-daily-at-dawn": "The lenses regain all expended charges daily at dawn"}, b, "charming recharge")

    d = detail["Eyes of Minute Seeing"][0]["eyes_of_minute_seeing"]
    b = block["Eyes of Minute Seeing"]
    tied([d["appearance"]], {"crystal-lenses-over-eyes": "These crystal lenses fit over the eyes"}, b, "minute appearance")
    tied([d["condition"]], {"while-wearing": "While wearing them"}, b, "minute condition")
    assert d["range_feet"] == 1 and has(b, "out to a range of 1 foot")
    assert d["darkvision_within_range"] is True and has(b, "granting you Darkvision within that range")
    assert (d["check_ability"], d["check_skill"], d["check_mode"], d["check_range_feet"]) == ("intelligence", "investigation", "advantage", 1)
    assert has(b, "Advantage on Intelligence (Investigation) checks made to examine something within that range")
    tied([d["check_subject"]], {"examine-something-within-range": "to examine something within that range"}, b, "minute subject")
    assert charges["Eyes of Minute Seeing"] == none

    d = detail["Eyes of the Eagle"][0]["eyes_of_the_eagle"]
    b = block["Eyes of the Eagle"]
    tied([d["appearance"]], {"crystal-lenses-over-eyes": "These crystal lenses fit over the eyes"}, b, "eagle appearance")
    tied([d["condition"]], {"while-wearing": "While wearing them"}, b, "eagle condition")
    assert (d["check_ability"], d["check_skill"], d["check_mode"]) == ("wisdom", "perception", "advantage")
    assert has(b, "Advantage on Wisdom (Perception) checks that rely on sight")
    tied([d["check_requirement"]], {"relies-on-sight": "that rely on sight"}, b, "eagle requirement")
    tied([d["detail_condition"]], {"clear-visibility": "In conditions of clear visibility"}, b, "eagle detail condition")
    assert d["detail_size_feet"] == 2 and has(b, "objects as small as 2 feet across")
    tied([d["detail_subjects"]], {"extremely-distant-creatures-and-objects": "details of even extremely distant creatures and objects"}, b, "eagle subjects")

    # Feather Token: one variant per printed kind, rarity printed after its name, and the 1d100 table.
    b = block["Feather Token"]
    assert has(b, "The type of token determines its rarity") and has(b, "rolling on the Feather Tokens table")
    table = re.search(r"Feather Tokens 1d100 Token Rarity (.*)$", b).group(1)
    rows = re.findall(r"(\d\d)–(\d\d) ([A-Za-z ]+?) (Uncommon|Rare)(?= \d\d–|\s*$)", table)
    assert len(rows) == 6, rows
    kinds = detail["Feather Token"]
    variants = families[5]["variants"]
    assert len(variants) == 6
    for (lo, hi, name, rarity), v, k in zip(rows, variants, kinds):
        t = k["feather_token"]
        hi = 100 if hi == "00" else int(hi)
        assert (t["table_die_sides"], t["table_min"], t["table_max"]) == (100, int(lo), hi), name
        assert name.lower().replace(" ", "-") == t["token"] == v["id"].rsplit("/", 1)[1], name
        assert v["name"].lower() == name.lower() and RARITY[v["rarity"]] == rarity, name
        assert f"{v['name']} ({rarity})." in b, name
        assert t["use_limit"] == "single-use" and t["action"] == "magic"
    assert [t["feather_token"]["table_max"] for t in kinds] == [20, 35, 50, 65, 90, 100] and kinds[-1]["feather_token"]["table_min"] == 91 and "91–00 Whip" in b
    for t in kinds:
        facts = [f["key"] for f in t["feather_token"]["facts"]]
        assert len(facts) == len(set(facts))
    sect = {m.group(1): m.group(0) for m in re.finditer(r"(Anchor|Bird|Fan|Swan Boat|Tree|Whip) \((?:Uncommon|Rare)\)\..*?(?=(?: (?:Bird|Fan|Swan Boat|Tree|Whip) \((?:Uncommon|Rare)\)\.)|\s*Feather Tokens 1d100)", b)}
    assert list(sect) == ["Anchor", "Bird", "Fan", "Swan Boat", "Tree", "Whip"], list(sect)

    def fact(kind, key, amount, unit, phrase):
        facts = {f["key"]: f for f in detail["Feather Token"][kind]["feather_token"]["facts"]}
        assert (facts[key]["amount"], facts[key]["unit"]) == (amount, unit), (kind, key)
        assert has(sect[names_of[kind]], phrase), (kind, key, phrase)

    names_of = ["Anchor", "Bird", "Fan", "Swan Boat", "Tree", "Whip"]
    t = kinds[0]["feather_token"]
    tied([t["placement"]], {"touch-token-to-boat-or-ship": "touch the token to a boat or ship"}, sect["Anchor"], "anchor placement")
    tied([t["effect"]], {"vessel-cannot-be-moved-by-any-means": "the vessel can’t be moved by any means"}, sect["Anchor"], "anchor effect")
    fact(0, "duration", 24, "hours", "For the next 24 hours")
    tied(t["ends_on"], {"duration-elapsed": "For the next 24 hours", "token-touched-to-vessel-again": "Touching the token to the vessel again ends the effect"}, sect["Anchor"], "anchor end")
    assert t["traits"] == ["touching-vessel-again-ends-effect"] and t["dismiss_action"] == "none" and "dismiss" not in sect["Anchor"]
    t = kinds[1]["feather_token"]
    tied([t["placement"]], {"toss-token-into-air": "toss the token 5 feet into the air"}, sect["Bird"], "bird placement")
    tied([t["effect"]], {"enormous-multicolored-bird-replaces-token": "an enormous, multicolored bird takes its place"}, sect["Bird"], "bird effect")
    for key, amount, unit, phrase in [("toss-height", 5, "feet", "toss the token 5 feet"), ("carry-weight-at-maximum-speed", 500, "pounds", "carry up to 500 pounds while flying at its maximum speed"),
                                      ("carry-weight-at-half-speed", 1000, "pounds", "or 1,000 pounds at half that speed"), ("maximum-flying-speed", 16, "miles-per-hour", "16 miles per hour"),
                                      ("maximum-distance-per-day", 144, "miles", "a maximum of 144 miles per day"), ("flying-hours-per-rest", 3, "hours", "for every 3 hours of flying"), ("rest-hours", 1, "hours", "with a 1-hour rest")]:
        fact(1, key, amount, unit, phrase)
    tied(t["traits"], {"statistics-of-roc": "The bird has the statistics of a Roc", "cannot-attack": "but it can’t attack", "obeys-simple-commands": "It obeys your simple commands"}, sect["Bird"], "bird traits")
    tied(t["ends_on"], {"maximum-distance-flown-for-a-day": "after flying its maximum distance for a day", "drops-to-0-hit-points": "if it drops to 0 Hit Points"}, sect["Bird"], "bird end")
    assert t["dismiss_action"] == "magic" and has(sect["Bird"], "You can dismiss the bird as a Magic action")
    t = kinds[2]["feather_token"]
    tied([t["placement"]], {"toss-token-into-air-while-on-boat-or-ship": "If you are on a boat or ship, you can take a Magic action to toss the token up to 10 feet in the air"}, sect["Fan"], "fan placement")
    tied([t["effect"]], {"giant-flapping-fan-replaces-token": "a giant flapping fan takes its place"}, sect["Fan"], "fan effect")
    for key, amount, unit, phrase in [("toss-height", 10, "feet", "up to 10 feet in the air"), ("ships-with-sails-filled", 1, "ship", "fill the sails of one ship"),
                                      ("speed-increase", 5, "miles-per-hour", "increasing its speed by 5 miles per hour"), ("duration", 8, "hours", "for 8 hours")]:
        fact(2, key, amount, unit, phrase)
    tied(t["traits"], {"floats": "The fan floats", "creates-strong-wind": "creates a strong wind"}, sect["Fan"], "fan traits")
    assert t["ends_on"] == ["duration-elapsed"] and t["dismiss_action"] == "magic" and has(sect["Fan"], "You can dismiss the fan as a Magic action")
    t = kinds[3]["feather_token"]
    tied([t["placement"]], {"touch-token-to-body-of-water": "touch the token to a body of water"}, sect["Swan Boat"], "swan placement")
    tied([t["effect"]], {"swan-shaped-boat-replaces-token": "boat shaped like a swan takes its place"}, sect["Swan Boat"], "swan effect")
    for key, amount, unit, phrase in [("minimum-water-diameter", 60, "feet", "at least 60 feet in diameter"), ("boat-length", 50, "feet", "a 50-foot-long"), ("boat-width", 20, "feet", "20-foot-wide boat"),
                                      ("speed", 6, "miles-per-hour", "at a speed of 6 miles per hour"), ("maximum-turn", 90, "degrees", "turn up to 90 degrees"), ("duration", 24, "hours", "The boat remains for 24 hours")]:
        fact(3, key, amount, unit, phrase)
    tied(t["traits"], {"self-propelled": "The boat is self-propelled", "moves-across-water": "moves across water", "commanded-with-magic-action-while-aboard": "You can take a Magic action while on the boat to command it to move or to turn"}, sect["Swan Boat"], "swan traits")
    assert t["ends_on"] == ["duration-elapsed"] and t["dismiss_action"] == "magic" and has(sect["Swan Boat"], "You can dismiss the boat as a Magic action")
    t = kinds[4]["feather_token"]
    tied([t["placement"]], {"touch-token-to-unoccupied-ground-space": "touch it to an unoccupied space on the ground"}, sect["Tree"], "tree placement")
    tied([t["effect"]], {"nonmagical-oak-tree-replaces-token": "a nonmagical oak tree springs into existence"}, sect["Tree"], "tree effect")
    for key, amount, unit, phrase in [("tree-height", 60, "feet", "The tree is 60 feet tall"), ("trunk-diameter", 5, "feet", "a 5-foot-diameter trunk"), ("branch-spread-radius", 20, "feet", "spread out in a 20-foot radius")]:
        fact(4, key, amount, unit, phrase)
    tied(t["traits"], {"must-be-outdoors": "You must be outdoors to use this token", "nonmagical-oak-tree": "a nonmagical oak tree"}, sect["Tree"], "tree traits")
    assert t["ends_on"] == [] and t["dismiss_action"] == "none" and "dismiss" not in sect["Tree"]
    t = kinds[5]["feather_token"]
    tied([t["placement"]], {"throw-token-to-point-within-10-feet": "throw the token to a point within 10 feet of yourself"}, sect["Whip"], "whip placement")
    tied([t["effect"]], {"floating-whip-replaces-token": "a floating whip takes its place"}, sect["Whip"], "whip effect")
    for key, amount, unit, phrase in [("placement-range", 10, "feet", "to a point within 10 feet of yourself"), ("attack-bonus", 9, "plus", "with an attack bonus of +9"),
                                      ("damage-dice", 1, "d6", "takes 1d6 + 5 Force damage"), ("damage-bonus", 5, "force-damage", "1d6 + 5 Force damage"),
                                      ("attack-range", 10, "feet", "within 10 feet of the whip"), ("flight-distance", 20, "feet", "fly up to 20 feet"), ("duration", 1, "hours", "disappears after 1 hour")]:
        fact(5, key, amount, unit, phrase)
    tied(t["traits"], {"melee-spell-attack": "make a melee spell attack against a creature", "attack-as-bonus-action": "You can then take a Bonus Action to make", "fly-and-repeat-attack-as-bonus-action": "As a Bonus Action, you can direct the whip to fly up to 20 feet and repeat the attack"}, sect["Whip"], "whip traits")
    tied(t["ends_on"], {"duration-elapsed": "disappears after 1 hour", "dismissed-with-magic-action": "when you take a Magic action to dismiss it", "you-die": "when you die", "you-are-incapacitated": "have the Incapacitated condition"}, sect["Whip"], "whip end")
    assert t["dismiss_action"] == "magic"
    assert charges["Feather Token"] == none

    # Figurine of Wondrous Power: shared rules, nine printed figurines, three goats.
    b = block["Figurine of Wondrous Power"]
    figs = detail["Figurine of Wondrous Power"]
    fvars = families[6]["variants"]
    assert [v["id"].rsplit("/", 1)[1] for v in fvars] == ["bronze-griffon", "ebony-fly", "golden-lions", "ivory-goats", "marble-elephant", "obsidian-steed", "onyx-dog", "serpentine-owl", "silver-raven"]
    for v in fvars:
        assert f"{v['name']} ({RARITY[v['rarity']]})." in b, v["name"]
    common = {k: figs[0]["figurine"][k] for k in figs[0]["figurine"] if k not in ("forms", "figurines_in_set")}
    for f in figs:
        assert {k: f["figurine"][k] for k in common} == common, "shared rules differ between figurines"
    c = common
    assert c["source_pages"] == [221, 222] and "Wondrous Item, Rarity Varies" in b and "lions are always created in pairs" in b
    assert c["action"] == "magic" and c["throw_range_feet"] == 60 and has(b, "If you take a Magic action to throw the figurine to a point on the ground within 60 feet of yourself")
    tied([c["throw_target"]], {"point-on-ground": "to a point on the ground"}, b, "figurine target")
    tied([c["no_space_result"]], {"figurine-does-not-become-creature": "the figurine doesn’t become a creature"}, b, "figurine no space")
    tied(c["friendly_to"], {"you": "The creature is Friendly to you", "your-allies": "to you and your allies"}, b, "figurine friendly")
    assert c["understands_owner_languages"] and has(b, "It understands your languages") and c["obeys_owner_commands"] and has(b, "obeys your commands")
    tied([c["initiative"]], {"immediately-after-owner": "takes its turn immediately after you on your Initiative count"}, b, "figurine initiative")
    tied([c["no_command_behavior"]], {"defends-itself-and-takes-no-other-actions": "If you issue no commands, the creature defends itself but takes no other actions"}, b, "figurine no command")
    tied(c["revert_triggers"], {"duration-elapsed": "At the end of the duration, the creature reverts to its figurine form", "creature-form-at-0-hit-points": "if its creature form drops to 0 Hit Points",
                                "magic-action-while-touching-creature": "if you take a Magic action while touching the creature to make it revert to figurine form"}, b, "figurine revert")
    assert c["revert_action"] == "magic"
    forms = {f["key"]: f for fig in figs for f in fig["figurine"]["forms"]}
    assert [f["figurine"]["figurines_in_set"] for f in figs] == [1, 1, 2, 3, 1, 1, 1, 1, 1] and len(forms) == 11
    assert has(b, "always created in sets of three") and has(b, "always created in pairs")

    def para(start, end):
        i = b.index(start)
        return b[i:] if end is None else b[i:b.index(end, i)]

    def form_check(key, creature, hours, days, text_, hours_phrase, days_phrase):
        f = forms[key]
        assert (f["creature"], f["duration_hours"], f["cooldown_days"]) == (creature, hours, days), key
        assert has(text_, hours_phrase) and has(text_, days_phrase), key
        return f

    p = para("Bronze Griffon (Rare)", "Ebony Fly (Rare)")
    form_check("bronze-griffon", "griffon", 6, 5, p, "can become a Griffon for up to 6 hours", "until 5 days have passed")
    p = para("Ebony Fly (Rare)", "Golden Lions (Rare)")
    f = form_check("ebony-fly", "giant-fly", 12, 2, p, "can become a Giant Fly (see the accompanying stat block) for up to 12 hours", "until 2 days have passed")
    tied(f["traits"], {"can-be-ridden-as-mount": "can be ridden as a mount", "stat-block-accompanies-item": "see the accompanying stat block", "large-beast": "Large Beast", "unaligned": "Unaligned", "languages-none": "Languages None"}, p, "fly traits")
    stat = {x["key"]: (x["amount"], x["unit"]) for x in f["facts"]}
    fly = para("Giant Fly Large Beast", "Golden Lions (Rare)")
    assert has(fly, "AC 11 Initiative +1 (11) HP 19 (3d10 + 3) Speed 30 ft., Fly 60 ft.")
    assert has(fly, "Str 14 +2 +2 Dex 13 +1 +1 Con 13 +1 +1 Int 2 −4 −4 Wis 10 +0 +0 Cha 3 −4 −4")
    assert has(fly, "Senses Darkvision 60 ft., Passive Perception 10 Languages None CR 0 (XP 0; PB +2)")
    assert [k for k in stat] == ["armor-class", "initiative-bonus", "initiative-score", "hit-points", "hit-point-dice", "hit-point-bonus", "speed", "fly-speed",
                                 "strength", "strength-modifier", "strength-save", "dexterity", "dexterity-modifier", "dexterity-save", "constitution", "constitution-modifier", "constitution-save",
                                 "intelligence", "intelligence-modifier", "intelligence-save", "wisdom", "wisdom-modifier", "wisdom-save", "charisma", "charisma-modifier", "charisma-save",
                                 "darkvision", "passive-perception", "challenge-rating", "experience-points", "proficiency-bonus"]
    assert (stat["armor-class"][0], stat["initiative-bonus"][0], stat["initiative-score"][0], stat["hit-points"][0], stat["hit-point-dice"], stat["hit-point-bonus"][0]) == (11, 1, 11, 19, (3, "d10"), 3)
    assert (stat["speed"][0], stat["fly-speed"][0], stat["darkvision"][0], stat["passive-perception"][0], stat["challenge-rating"][0], stat["experience-points"][0], stat["proficiency-bonus"][0]) == (30, 60, 60, 10, 0, 0, 2)
    for ability, score, mod, sign in [("strength", 14, 2, "plus"), ("dexterity", 13, 1, "plus"), ("constitution", 13, 1, "plus"), ("intelligence", 2, 4, "minus"), ("wisdom", 10, 0, "plus"), ("charisma", 3, 4, "minus")]:
        assert stat[ability] == (score, "score") and stat[ability + "-modifier"] == (mod, sign) and stat[ability + "-save"] == (mod, sign), ability
    p = para("Golden Lions (Rare)", "Ivory Goats (Rare)")
    f = form_check("golden-lions", "lion", 1, 7, p, "Each can become a Lion for up to 1 hour", "until 7 days have passed")
    assert f["facts"] == [{"key": "figurines-usable-simultaneously", "amount": 2, "unit": "figurines"}] and has(p, "You can use one figurine or both simultaneously")
    assert f["traits"] == ["cooldown-per-lion"] and has(p, "Once a lion has been used")
    p = para("Goat of Terror.", "Goat of Traveling.")
    f = form_check("goat-of-terror", "giant-goat", 3, 15, p, "can become a Giant Goat for up to 3 hours", "until 15 days have passed")
    tied(f["traits"], {"cannot-attack": "The goat can’t attack", "horns-removable-as-weapons-with-magic-action": "Removing a horn requires a Magic action",
                       "horn-weapons-disappear-when-goat-reverts": "the weapons disappear and the horns return when the goat reverts to figurine form",
                       "frightens-hostile-creatures-starting-turn-in-emanation-while-ridden": "While you ride the goat, any Hostile creature that starts its turn within a 30-foot Emanation originating from the goat",
                       "wisdom-save": "DC 15 Wisdom saving throw", "frightened-ends-when-no-longer-riding-or-goat-reverts": "until you are no longer riding the goat, or until the goat reverts to figurine form",
                       "repeats-save-at-end-of-each-turn": "repeats the save at the end of each of its turns"}, p, "terror traits")
    facts = {x["key"]: (x["amount"], x["unit"]) for x in f["facts"]}
    assert facts == {"horn-weapons": (2, "horns"), "lance-bonus": (1, "plus"), "longsword-bonus": (2, "plus"), "fear-emanation": (30, "feet"), "save-dc": (15, "dc"), "frightened-duration": (60, "seconds"), "immunity-after-success": (24, "hours")}
    for phrase in ("One horn becomes a +1 Lance, and the other becomes a +2 Longsword", "a 30-foot Emanation", "DC 15 Wisdom saving throw", "Frightened condition for 1 minute", "immune to this effect for the next 24 hours"):
        assert has(p, phrase), phrase
    p = para("Goat of Traveling.", "Goat of Travail.")
    f = form_check("goat-of-traveling", "large-goat", 0, 7, p, "Large goat with the same statistics as a Riding Horse", "until 7 days have passed")
    assert f["facts"] == [{"key": "charges", "amount": 24, "unit": "charges"}, {"key": "charge-cost-per-hour-or-portion", "amount": 1, "unit": "charges"}]
    assert has(p, "It has 24 charges, and each hour or portion thereof it spends in goat form costs 1 charge") and "hours" not in p.split("Large goat")[1].split("It has 24")[0]
    tied(f["traits"], {"statistics-of-riding-horse": "the same statistics as a Riding Horse", "usable-as-often-as-wished-while-charged": "While it has charges, you can use it as often as you wish",
                       "reverts-when-out-of-charges": "When it runs out of charges, it reverts to a figurine", "regains-all-charges-after-cooldown": "when it regains all expended charges"}, p, "traveling traits")
    p = para("Goat of Travail.", "Marble Elephant (Rare)")
    f = form_check("goat-of-travail", "giant-goat", 3, 30, p, "can become a Giant Goat for up to 3 hours", "until 30 days have passed")
    assert f["facts"] == [] and f["traits"] == []
    p = para("Marble Elephant (Rare)", "Obsidian Steed (Very Rare)")
    f = form_check("marble-elephant", "elephant", 24, 7, p, "can become an Elephant for up to 24 hours", "until 7 days have passed")
    assert f["facts"] == [] and f["traits"] == []
    p = para("Obsidian Steed (Very Rare)", "Onyx Dog (Rare)")
    f = form_check("obsidian-steed", "nightmare", 24, 5, p, "can become a Nightmare for up to 24 hours", "until 5 days have passed")
    assert f["facts"] == [{"key": "ignore-orders-chance", "amount": 10, "unit": "percent"}] and has(p, "a 10 percent chance each time you use it to ignore your orders")
    tied(f["traits"], {"fights-only-to-defend-itself": "The nightmare fights only to defend itself", "may-ignore-revert-command": "including a command to revert to figurine form",
                       "mounted-while-ignoring-orders-transports-to-random-location-on-plane-of-hades": "If you mount the nightmare while it is ignoring your orders, you and the nightmare are instantly transported to a random location on the plane of Hades",
                       "reverts-to-figurine-on-plane-of-hades": "where the nightmare reverts to figurine form"}, p, "steed traits")
    p = para("Onyx Dog (Rare)", "Serpentine Owl (Rare)")
    f = form_check("onyx-dog", "mastiff", 6, 7, p, "can become a Mastiff for up to 6 hours", "until 7 days have passed")
    assert f["facts"] == [{"key": "intelligence", "amount": 8, "unit": "score"}, {"key": "blindsight", "amount": 60, "unit": "feet"}] and f["traits"] == ["can-speak-common"]
    assert has(p, "has an Intelligence of 8 and can speak Common") and has(p, "Blindsight with a range of 60 feet")
    p = para("Serpentine Owl (Rare)", "Silver Raven (Uncommon)")
    f = form_check("serpentine-owl", "giant-owl", 8, 2, p, "can become a Giant Owl for up to 8 hours", "until 2 days have passed")
    assert f["facts"] == [] and f["traits"] == ["telepathy-with-owner-at-any-range-on-same-plane"]
    assert has(p, "communicate telepathically with you at any range if you and it are on the same plane of existence")
    p = para("Silver Raven (Uncommon)", None)
    f = form_check("silver-raven", "raven", 12, 2, p, "can become a Raven for up to 12 hours", "until 2 days have passed")
    assert f["facts"] == [] and f["traits"] == ["raven-form-grants-owner-animal-messenger-cast-on-it"]
    assert has(p, "While in raven form, the figurine grants you the ability to cast Animal Messenger on it")
    assert charges["Figurine of Wondrous Power"] == none and "charges" in forms["goat-of-traveling"]["facts"][0]["key"]

    d = detail["Flame Tongue"][0]["flame_tongue"]
    b = block["Flame Tongue"]
    tied([d["condition"]], {"while-holding": "While holding this magic weapon"}, b, "tongue condition")
    assert d["action"] == "bonus" and has(b, "you can take a Bonus Action and use a command word")
    tied([d["activation"]], {"command-word": "use a command word"}, b, "tongue activation")
    tied([d["effect"]], {"flames-engulf-damage-dealing-part": "to cause flames to engulf the damage-dealing part of the weapon"}, b, "tongue effect")
    assert (d["bright_light_feet"], d["dim_light_additional_feet"]) == (40, 40) and has(b, "Bright Light in a 40-foot radius and Dim Light for an additional 40 feet")
    assert d["extra_damage"] == {"dice": 2, "sides": 6, "type": "fire"} and has(b, "While the weapon is ablaze, it deals an extra 2d6 Fire damage on a hit")
    tied([d["extra_damage_condition"]], {"on-hit-while-ablaze": "While the weapon is ablaze, it deals an extra 2d6 Fire damage on a hit"}, b, "tongue damage condition")
    tied(d["ends_on"], {"bonus-action-command-again": "until you take a Bonus Action to issue the command again", "drop-weapon": "you drop", "stow-weapon": "stow", "sheathe-weapon": "or sheathe the weapon"}, b, "tongue end")
    assert charges["Flame Tongue"] == none

    d = detail["Folding Boat"][0]["folding_boat"]
    b = block["Folding Boat"]
    assert (d["length_inches"], d["width_inches"], d["depth_inches"], d["weight_pounds"]) == (12, 6, 6, 4)
    assert has(b, "a wooden box that measures 12 inches long, 6 inches wide, and 6 inches deep") and has(b, "It weighs 4 pounds and floats")
    assert d["floats"] is True and d["stores_items_when_open"] is True and has(b, "It can be opened to store items inside")
    assert d["command_action"] == "magic" and has(b, "three command words, each requiring a Magic action to use")
    cmds = d["command_words"]
    assert [c["ordinal"] for c in cmds] == [1, 2, 3] and [c["result"] for c in cmds] == ["rowboat", "keelboat", "folds-back-into-box"]
    assert has(b, "First Command Word. The box unfolds into a Rowboat") and has(b, "Second Command Word. The box unfolds into a Keelboat")
    assert has(b, "Third Command Word. The Folding Boat folds back into a box if no creatures are aboard") and cmds[:2] == [{"ordinal": 1, "result": "rowboat", "requires": []}, {"ordinal": 2, "result": "keelboat", "requires": []}]
    assert cmds[2]["requires"] == ["no-creatures-aboard"]
    tied([d["objects_that_do_not_fit"]], {"remain-outside-box": "Any objects in the vessel that can’t fit inside the box remain outside the box as it folds"}, b, "boat outside")
    tied([d["objects_that_fit"]], {"stored-inside-box": "Any objects in the vessel that can fit inside the box do so"}, b, "boat inside")
    tied([d["vessel_weight"]], {"normal-weight-of-vessel-its-size": "its weight becomes that of a normal vessel its size"}, b, "boat weight")
    tied([d["stored_items_when_vessel"]], {"remain-in-boat": "anything that was stored in the box remains in the boat"}, b, "boat stored")
    tied(d["vessel_statistics_in"], {"equipment": "Statistics for the Rowboat and Keelboat appear in “Equipment.”"}, b, "boat statistics")
    tied([d["destroyed_when"]], {"vessel-reduced-to-0-hit-points": "If either vessel is reduced to 0 Hit Points, the Folding Boat is destroyed"}, b, "boat destroyed")
    assert charges["Folding Boat"] == none

    d = detail["Frost Brand"][0]["frost_brand"]
    b = block["Frost Brand"]
    assert d["extra_damage"] == {"dice": 1, "sides": 6, "type": "cold"} and has(b, "the target takes an extra 1d6 Cold damage")
    tied([d["extra_damage_trigger"]], {"hit-with-attack-roll-using-weapon": "When you hit with an attack roll using this magic weapon"}, b, "brand trigger")
    tied([d["resistance_condition"]], {"while-holding": "while you hold the weapon"}, b, "brand resistance condition")
    assert d["resistance"] == "fire" and has(b, "you have Resistance to Fire damage")
    tied([d["light_condition"]], {"freezing-temperatures": "In freezing temperatures"}, b, "brand light condition")
    assert (d["bright_light_feet"], d["dim_light_additional_feet"]) == (10, 10) and has(b, "Bright Light in a 10-foot radius and Dim Light for an additional 10 feet")
    tied([d["extinguish_trigger"]], {"draw-weapon": "When you draw this weapon"}, b, "brand extinguish trigger")
    tied([d["extinguish_targets"]], {"nonmagical-flames": "extinguish all nonmagical flames"}, b, "brand extinguish targets")
    assert (d["extinguish_radius_feet"], d["extinguish_uses"], d["extinguish_cooldown_hours"]) == (30, 1, 1)
    assert has(b, "within 30 feet of yourself") and has(b, "Once used, this property can’t be used again for 1 hour")
    assert charges["Frost Brand"] == none

    d = detail["Gauntlets of Ogre Power"][0]["gauntlets_of_ogre_power"]
    b = block["Gauntlets of Ogre Power"]
    tied([d["condition"]], {"while-wearing": "while you wear these gauntlets"}, b, "gauntlets condition")
    assert d["strength_score"] == 19 and has(b, "Your Strength is 19")
    assert d["no_effect_at_or_above_score"] == 19 and has(b, "no effect on you if your Strength is 19 or higher without them")
    assert charges["Gauntlets of Ogre Power"] == none

    d = detail["Gem of Brightness"][0]["gem_of_brightness"]
    b = block["Gem of Brightness"]
    assert charges["Gem of Brightness"] == {"mode": "limited", "maximum": 50, "per_use": 1, "recharge": "none"} and d["max_charges"] == 50 and has(b, "This prism has 50 charges")
    tied([d["condition"]], {"while-holding": "While you are holding it"}, b, "gem condition")
    assert d["action"] == "magic" and d["command_words"] == 3 and has(b, "you can take a Magic action and use one of three command words")
    first = d["first_command_word"]
    assert (first["charge_cost"], first["bright_light_feet"], first["dim_light_additional_feet"]) == (0, 30, 30)
    assert has(b, "First Command Word. The gem sheds Bright Light in a 30-foot radius and Dim Light for an additional 30 feet") and has(b, "This effect doesn’t expend a charge")
    tied(first["ends_on"], {"bonus-action-repeat-command-word": "until you take a Bonus Action to repeat the command word", "another-function-used": "or until you use another function of the gem"}, b, "gem first end")
    second = d["second_command_word"]
    assert second["charge_cost"] == 1 and has(b, "Second Command Word. You expend 1 charge")
    tied([second["target"]], {"one-creature-you-can-see": "one creature you can see within 60 feet of yourself"}, b, "gem target")
    assert second["range_feet"] == 60 and (second["save_ability"], second["save_dc"], second["condition"], second["duration_seconds"]) == ("constitution", 15, "blinded", 60)
    assert has(b, "fire a brilliant beam of light") and has(b, "succeed on a DC 15 Constitution saving throw or have the Blinded condition for 1 minute")
    tied([second["repeat_save"]], {"end-of-each-of-its-turns": "repeats the save at the end of each of its turns"}, b, "gem repeat")
    tied([second["success_ends"]], {"effect-on-itself": "ending the effect on itself on a success"}, b, "gem success")
    third = d["third_command_word"]
    assert (third["charge_cost"], third["shape"], third["size_feet"]) == (5, "cone", 30) and has(b, "Third Command Word. You expend 5 charges") and has(b, "intense light in a 30-foot Cone")
    tied([third["targets"]], {"each-creature-in-cone": "Each creature in the Cone makes a saving throw"}, b, "gem cone targets")
    tied([third["save"]], {"as-struck-by-second-command-word-beam": "as if struck by the beam created with the second command word"}, b, "gem cone save")
    tied([d["when_depleted"]], {"becomes-nonmagical-jewel": "the gem becomes a nonmagical jewel"}, b, "gem depleted")
    assert d["depleted_value_gp"] == 50 and has(b, "nonmagical jewel worth 50 GP")

    for name in names:
        if name not in ("Eyes of Charming", "Gem of Brightness"):
            assert charges[name] == none, name
    print("#231 pinned PDF, twelve identities/pages/headers, 25 variants, the 1d100 token table, charges and source values pass")


if __name__ == "__main__":
    main()
