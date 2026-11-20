#!/usr/bin/env python3
"""Compare issue #228 item identities, headers and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (215, 216)
RARITY = {"uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")


def text(pdf, page, x):
    args = ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x), "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"]
    return subprocess.run(args, check=True, capture_output=True, text=True).stdout


def clean(raw):
    return re.sub(r"(\w)- (\w)", r"\1\2", FOOTER.sub(" ", " ".join(raw.split())))


def header(family):
    """Rebuild the printed category and rarity line from the catalogued fields."""
    variant = family["variants"][0]
    assert family["category"] == "wondrous-item" and not any(family["base"].values())
    line = f"Wondrous Item, {RARITY[variant['rarity']]}"
    return line + (" (Requires Attunement)" if variant["attunement"]["mode"] == "required" else "")


def tied(values, phrases, source, label):
    """Require the catalogued keys to be exactly the mapped ones, each with its printed phrase."""
    assert list(values) == list(phrases), f"{label}: {values}"
    for value, phrase in phrases.items():
        assert phrase in source, f"{label} {value}: {phrase}"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_228.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order, so the Cloak
    # of Invisibility continues from the end of p. 215 to the start of p. 216.
    stream = " ".join(clean(text(pdf, p, x)) for p in PAGES for x in (0, 306))

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_228.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 228]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 228
    assert len(families) == 12 and all(len(f["variants"]) == 1 for f in families)
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    block = {}
    names = [f["name"] for f in families]
    for f in families:
        page = f["presentations"][0]["pages"]
        assert page[0] == page[1] and page[0] in PAGES, f["name"]
        printed = f"{f['name']} {header(f)} "
        assert stream.count(printed) == 1, printed
        # The section runs to the next heading of any family on the pages.
        body = stream.split(printed, 1)[1]
        end = [m.start() for m in re.finditer(r" (?:[A-Z][A-Za-z’']+ )+(?:Wondrous Item|Weapon|Armor|Ring|Rod|Staff|Wand|Potion|Scroll)[ ,(]", body)]
        block[f["name"]] = body[:end[0]] if end else body
        cost = f["variants"][0]["cost"]
        assert cost == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}
        assert " GP" not in block[f["name"]]
    # The printed order of the twelve headings follows the catalogue order within each page.
    for page in PAGES:
        order = [stream.index(f"{n} {header(f)} ") for n, f in zip(names, families) if f["presentations"][0]["pages"][0] == page]
        assert order == sorted(order)
    detail = {f["name"]: f["variants"][0]["details"] for f in families}
    charges = {f["name"]: f["variants"][0]["charges"] for f in families}
    none = {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}

    c = detail["Censer of Controlling Air Elementals"]["censer"]
    b = block["Censer of Controlling Air Elementals"]
    tied([c["condition"]], {"gently-swinging": "While gently swinging this censer"}, b, "censer condition")
    tied([c["action"]], {"magic": "you can take a Magic action to summon an Air Elemental"}, b, "censer action")
    tied([c["summons"]], {"air-elemental": "summon an Air Elemental"}, b, "censer summons")
    tied([c["location"]], {"unoccupied-space-nearest-censer": "appears in an unoccupied space as close to the censer as possible"}, b, "censer location")
    assert c["understands_owner_languages"] and "understands your languages" in b
    assert c["obeys_owner_commands"] and "obeys your commands" in b
    tied([c["initiative"]], {"immediately-after-summoner": "takes its turn immediately after you on your Initiative count"}, b, "censer initiative")
    assert (c["duration_seconds"], c["uses"], c["recharge"]) == (3600, 1, "next-dawn")
    tied(c["end_triggers"], {"duration-elapsed": "disappears after 1 hour", "dies": "when it dies", "dismissed-as-bonus-action": "when you dismiss it as a Bonus Action"}, b, "censer end")
    assert "can’t be used this way again until the next dawn" in b

    c = detail["Chime of Opening"]["chime"]
    b = block["Chime of Opening"]
    tied([c["appearance"]], {"hollow-metal-tube": "This hollow metal tube"}, b, "chime appearance")
    assert (c["length_feet"], c["weight_pounds"], c["audible_feet"], c["uses"]) == (1, 1, 300, 10)
    for snippet in ("measures about 1 foot long and weighs 1 pound", "audible out to 300 feet", "The chime can be used 10 times"):
        assert snippet in b, snippet
    tied([c["action"]], {"magic": "As a Magic action"}, b, "chime action")
    tied([c["spell"]], {"knock": "strike the chime to cast Knock"}, b, "chime spell")
    tied([c["sound_replaces"]], {"customary-knocking-sound": "customary knocking sound is replaced by the clear, ringing tone"}, b, "chime sound")
    tied([c["after_final_use"]], {"cracks-and-becomes-useless": "After the tenth time, it cracks and becomes useless"}, b, "chime end")
    assert charges["Chime of Opening"] == none and "charge" not in b.lower()

    c = detail["Circlet of Blasting"]["circlet"]
    b = block["Circlet of Blasting"]
    tied([c["condition"]], {"while-wearing": "While wearing this circlet"}, b, "circlet condition")
    tied([c["spell"]], {"scorching-ray": "you can cast Scorching Ray with it"}, b, "circlet spell")
    assert c["attack_bonus"] == 5 and "(+5 to hit)" in b
    assert (c["uses"], c["recharge"]) == (1, "next-dawn") and "can’t cast this spell again until the next dawn" in b

    c = detail["Cloak of Arachnida"]["arachnida"]
    b = block["Cloak of Arachnida"]
    tied([c["condition"]], {"while-wearing": "While wearing it, you gain the following benefits"}, b, "arachnida condition")
    tied(c["resistances"], {"poison": "You have Resistance to Poison damage"}, b, "arachnida resistance")
    tied([c["climb_speed"]], {"equal-to-speed": "You have a Climb Speed equal to your Speed"}, b, "arachnida climb")
    tied(c["climb_directions"], {"up": "move up", "down": "down", "across-vertical-surfaces": "and across vertical surfaces", "along-ceilings": "and along ceilings"}, b, "arachnida directions")
    assert c["climb_hands_free"] and "while leaving your hands free" in b
    tied([c["web_catch"]], {"cannot-be-caught-in-webs": "You can’t be caught in webs of any sort"}, b, "arachnida web catch")
    tied([c["web_movement"]], {"as-difficult-terrain": "move through webs as if they were Difficult Terrain"}, b, "arachnida web movement")
    tied([c["web_spell"]], {"web": "You can cast Web (save DC 13)"}, b, "arachnida spell")
    assert c["web_save_dc"] == 13 and c["web_area_multiplier"] == 2 and "fills twice its normal area" in b
    assert (c["web_uses"], c["web_recharge"]) == (1, "next-dawn") and "can’t be used again until the next dawn" in b

    c = detail["Cloak of Displacement"]["displacement"]
    b = block["Cloak of Displacement"]
    tied([c["condition"]], {"while-wearing": "While you wear this cloak"}, b, "displacement condition")
    tied([c["effect"]], {"illusion-appears-near-actual-location": "projects an illusion that makes you appear to be standing in a place near your actual location"}, b, "displacement effect")
    tied([c["attack_rolls_against_wearer"]], {"disadvantage": "any creature to have Disadvantage on attack rolls against you"}, b, "displacement attacks")
    tied([c["suspended_by"], c["resumes_at"]], {"taking-damage": "If you take damage", "start-of-next-turn": "until the start of your next turn"}, b, "displacement suspension")
    tied([c["suppressed_while"]], {"speed-0": "suppressed while your Speed is 0"}, b, "displacement suppression")

    c = detail["Cloak of Elvenkind"]["elvenkind"]
    b = block["Cloak of Elvenkind"]
    tied([c["condition"]], {"while-wearing": "While you wear this cloak"}, b, "elvenkind condition")
    assert (c["perceive_wearer_ability"], c["perceive_wearer_skill"], c["perceive_wearer_mode"]) == ("wisdom", "perception", "disadvantage")
    assert "Wisdom (Perception) checks made to perceive you have Disadvantage" in b
    assert (c["stealth_ability"], c["stealth_skill"], c["stealth_mode"]) == ("dexterity", "stealth", "advantage")
    assert "you have Advantage on Dexterity (Stealth) checks" in b

    c = detail["Cloak of Invisibility"]["invisibility"]
    b = block["Cloak of Invisibility"]
    assert "head and expend 1 charge" in b, "continuation across the p. 215/216 break"
    assert charges["Cloak of Invisibility"] == {"mode": "limited", "maximum": 3, "per_use": 1, "recharge": "regain-1d3-at-dawn"}
    assert c["max_charges"] == 3 and (c["regain_dice"], c["regain_timing"]) == ({"dice": 1, "sides": 3, "type": "charges"}, "dawn")
    assert "This cloak has 3 charges and regains 1d3 expended charges daily at dawn" in b
    tied([c["condition"]], {"while-wearing": "While wearing the cloak"}, b, "invisibility condition")
    tied([c["action"]], {"magic": "you can take a Magic action"}, b, "invisibility action")
    tied([c["gesture"]], {"pull-hood-over-head": "to pull its hood over your head"}, b, "invisibility gesture")
    assert c["charge_cost"] == 1 and "expend 1 charge" in b
    tied([c["grants_condition"]], {"invisible": "give yourself the Invisible condition"}, b, "invisibility condition granted")
    assert c["duration_seconds"] == 3600 and "for 1 hour" in b
    tied(c["ends_early"], {"hood-pulled-down": "if you pull the hood down", "cloak-no-longer-worn": "or cease wearing the cloak"}, b, "invisibility end")
    tied([c["hood_down_action"]], {"none": "(no action required)"}, b, "invisibility hood action")

    c = detail["Cloak of Protection"]["protection"]
    b = block["Cloak of Protection"]
    tied([c["condition"]], {"while-wearing": "while you wear this cloak"}, b, "protection condition")
    assert (c["ac_bonus"], c["saving_throw_bonus"]) == (1, 1) and "You gain a +1 bonus to Armor Class and saving throws" in b

    c = detail["Cloak of the Bat"]["bat"]
    b = block["Cloak of the Bat"]
    tied([c["condition"]], {"while-wearing": "While wearing this cloak"}, b, "bat condition")
    assert (c["stealth_ability"], c["stealth_skill"], c["stealth_mode"]) == ("dexterity", "stealth", "advantage")
    assert "you have Advantage on Dexterity (Stealth) checks" in b
    assert c["fly_speed_feet"] == 40 and "Fly Speed of 40 feet" in b
    tied(c["fly_light"], {"dim-light": "In an area of Dim Light", "darkness": "or Darkness"}, b, "bat fly light")
    tied([c["fly_grip"]], {"grip-edges-of-cloak": "you can grip the edges of the cloak"}, b, "bat grip")
    tied(c["fly_ends_on"], {"fail-to-grip-edges-while-flying": "If you ever fail to grip the cloak’s edges while flying in this way",
                            "no-longer-in-dim-light-or-darkness": "if you are no longer in Dim Light or Darkness"}, b, "bat fly end")
    tied([c["polymorph_spell"]], {"polymorph": "you can cast Polymorph on yourself"}, b, "bat spell")
    tied([c["polymorph_target"]], {"self": "Polymorph on yourself"}, b, "bat target")
    tied([c["polymorph_form"]], {"bat": "shape-shifting into a Bat"}, b, "bat form")
    tied(c["polymorph_light"], {"dim-light": "While wearing the cloak in an area of Dim Light", "darkness": "or Darkness, you can cast Polymorph"}, b, "bat polymorph light")
    tied(c["retained_scores"], {"intelligence": "you retain your Intelligence", "wisdom": "Wisdom", "charisma": "and Charisma scores"}, b, "bat scores")
    assert (c["polymorph_uses"], c["polymorph_recharge"]) == (1, "next-dawn") and "can’t be used this way again until the next dawn" in b

    c = detail["Cloak of the Manta Ray"]["manta_ray"]
    b = block["Cloak of the Manta Ray"]
    tied([c["condition"]], {"while-wearing": "While wearing this cloak"}, b, "manta ray condition")
    assert c["breathe_underwater"] and "you can breathe underwater" in b
    assert c["swim_speed_feet"] == 60 and "Swim Speed of 60 feet" in b

    c = detail["Crystal Ball"]["crystal_ball"]
    b = block["Crystal Ball"]
    tied([c["condition"]], {"while-touching": "While touching this crystal orb"}, b, "crystal ball condition")
    tied([c["spell"]], {"scrying": "you can cast Scrying (save DC 17) with it"}, b, "crystal ball spell")
    assert c["save_dc"] == 17 and not re.search(r"Detect|Telepathy|Truesight", b)

    c = detail["Crystal Ball of Mind Reading"]["mind_reading"]
    b = block["Crystal Ball of Mind Reading"]
    tied([c["condition"]], {"while-touching": "While touching this crystal orb"}, b, "mind reading condition")
    tied([c["spell"]], {"scrying": "you can cast Scrying (save DC 17) with it"}, b, "mind reading spell")
    tied([c["extra_spell"]], {"detect-thoughts": "you can cast Detect Thoughts (save DC 17)"}, b, "mind reading extra spell")
    assert (c["save_dc"], c["extra_save_dc"]) == (17, 17)
    tied([c["extra_targets"]], {"creatures-you-can-see-within-range-of-sensor": "targeting creatures you can see within 30 feet of the spell’s sensor"}, b, "mind reading targets")
    assert c["extra_range_feet"] == 30
    tied([c["extra_range_from"]], {"spell-sensor": "of the spell’s sensor"}, b, "mind reading range origin")
    assert not c["extra_concentration_required"] and "You don’t need to concentrate on this Detect Thoughts spell to maintain it" in b
    tied([c["extra_ends_with"]], {"scrying-ends": "but it ends if the Scrying spell ends"}, b, "mind reading end")

    for name in names:
        if name not in ("Cloak of Invisibility",):
            assert charges[name] == none, name
    print("#228 pinned PDF, twelve identities/pages/headers, one variant each, charges and source values pass")


if __name__ == "__main__":
    main()
