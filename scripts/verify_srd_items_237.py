#!/usr/bin/env python3
"""Compare issue #237 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (235, 236)
RARITY = {"common": "Common", "uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# Heading that follows the last family on p. 235 (the Healing potions open p. 236)
# and the one that follows the last family on p. 236.
HEALING = "Potions of Healing Potion, Rarity Varies"
FOLLOWING = "Potion of Poison Potion, Uncommon"


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


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_237.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream = " ".join(clean(text(pdf, p, x)) for p in PAGES for x in (0, 306))

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_237.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 237]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 237 and catalogue["batch"] == "data-magic-items-13"
    assert len(families) == 12 and sum(len(f["variants"]) for f in families) == 16
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert [f["presentations"][0]["pages"][0] for f in families] == [235] * 7 + [236] * 5
    starts = []
    for f in families:
        variants = f["variants"]
        rarity = "Rarity Varies" if f["name"] == "Potion of Giant Strength" else RARITY[variants[0]["rarity"]]
        printed = f"{f['name']} Potion, {rarity}"
        assert stream.count(printed) == 1, printed
        assert f["category"] == "potion" and f["base"] == {"options": [], "excluded": []}
        starts.append(stream.index(printed))
        page = f["presentations"][0]["pages"]
        for v in variants:
            assert v["attunement"]["mode"] == "none" and v["charges"] == {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}, f["name"]
            assert v["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}, f["name"]
    assert starts == sorted(starts), "printed order differs from the catalogue order"
    ends = starts[1:] + [stream.index(FOLLOWING)]
    ends[6] = stream.index(HEALING)  # Potion of Growth ends the p. 235 right column
    block = {f["name"]: stream[s:e] for f, s, e in zip(families, starts, ends)}
    one = {f["name"]: f["variants"][0]["details"] for f in families}
    assert " GP" not in stream[starts[0]:ends[-1]], "a price is printed"

    d = one["Potion of Clairvoyance"]["clairvoyance"]
    b = block["Potion of Clairvoyance"]
    tied([d["spell"]], {"clairvoyance": "you gain the effect of the Clairvoyance spell"}, b, "clairvoyance spell")
    assert d["concentration_required"] is False and has(b, "(no Concentration required)")
    tied([d["appearance"]], {"yellowish-liquid-with-bobbing-eyeball-that-vanishes-when-opened": "An eyeball bobs in this potion’s yellowish liquid but vanishes when the potion is opened"}, b, "clairvoyance appearance")

    d = one["Potion of Climbing"]["climbing"]
    b = block["Potion of Climbing"]
    assert (d["speed_kind"], d["duration_hours"]) == ("climb", 1)
    tied([d["speed_equals"]], {"your-speed": "you gain a Climb Speed equal to your Speed for 1 hour"}, b, "climbing speed")
    assert (d["check_ability"], d["check_skill"], d["check_mode"]) == ("strength", "athletics", "advantage")
    tied([d["check_subject"]], {"climbing": "Advantage on Strength (Athletics) checks to climb"}, b, "climbing check")
    tied([d["appearance"]], {"brown-silver-gray-layers-like-bands-of-stone-that-shaking-fails-to-mix": "separated into brown, silver, and gray layers resembling bands of stone. Shaking the bottle fails to mix the colors"}, b, "climbing appearance")

    d = one["Potion of Diminution"]["diminution"]
    b = block["Potion of Diminution"]
    assert (d["spell"], d["spell_effect"]) == ("enlarge-reduce", "reduce") and has(b, "you gain the “reduce” effect of the Enlarge/Reduce spell")
    assert (d["duration_dice"], d["duration_die_sides"], d["duration_unit"]) == (1, 4, "hours") and has(b, "for 1d4 hours (no Concentration required)")
    assert d["concentration_required"] is False
    tied([d["appearance"]], {"red-contracts-to-tiny-bead-then-expands-to-color-clear-liquid-shaking-fails-to-interrupt": "The red in the potion’s liquid continuously contracts to a tiny bead and then expands to color the clear liquid around it. Shaking the bottle fails to interrupt this process"}, b, "diminution appearance")

    d = one["Potion of Flying"]["flying"]
    b = block["Potion of Flying"]
    assert (d["speed_kind"], d["duration_hours"], d["can_hover"]) == ("fly", 1, True)
    tied([d["speed_equals"]], {"your-speed": "you gain a Fly Speed equal to your Speed for 1 hour and can hover"}, b, "flying speed")
    tied([d["falls_when_it_ends"]], {"when-in-the-air": "If you’re in the air when the potion wears off, you fall"}, b, "flying fall")
    tied([d["fall_exception"]], {"other-means-of-staying-aloft": "unless you have some other means of staying aloft"}, b, "flying exception")
    tied([d["appearance"]], {"clear-liquid-floats-at-top-with-cloudy-white-impurities-drifting": "clear liquid floats at the top of its container and has cloudy white impurities drifting in it"}, b, "flying appearance")

    d = one["Potion of Gaseous Form"]["gaseous_form"]
    b = block["Potion of Gaseous Form"]
    assert d["spell"] == "gaseous-form" and d["duration_hours"] == 1 and d["concentration_required"] is False
    assert has(b, "you gain the effect of the Gaseous Form spell for 1 hour (no Concentration required)")
    tied([d["ends_early_action"]], {"bonus": "or until you end the effect as a Bonus Action"}, b, "gaseous end")
    tied([d["appearance"]], {"container-seems-to-hold-fog-that-moves-and-pours-like-water": "container seems to hold fog that moves and pours like water"}, b, "gaseous appearance")

    # Potion of Giant Strength: one variant per row of the table, rarity printed in each row.
    b = block["Potion of Giant Strength"]
    f = families[5]
    assert has(b, "your Strength score changes for 1 hour") and has(b, "The type of giant determines the score")
    assert has(b, "The potion has no effect on you if your Strength is equal to or greater than that score")
    assert has(b, "Potion Str. Rarity")
    rows = [
        ("hill", ["hill"], 21, "Uncommon", "Potion of Giant Strength (hill) 21 Uncommon"),
        ("frost-or-stone", ["frost", "stone"], 23, "Rare", "Potion of Giant Strength (frost or 23 Rare stone)"),
        ("fire", ["fire"], 25, "Rare", "Potion of Giant Strength (fire) 25 Rare"),
        ("cloud", ["cloud"], 27, "Very Rare", "Potion of Giant Strength (cloud) 27 Very Rare"),
        ("storm", ["storm"], 29, "Legendary", "Potion of Giant Strength (storm) 29 Legendary"),
    ]
    assert len(f["variants"]) == len(rows)
    for v, (slug, giants, score, rarity, printed) in zip(f["variants"], rows):
        g = v["details"]["giant_strength"]
        assert v["id"].endswith("/variant/" + slug) and has(b, printed), slug
        assert (g["giants"], g["strength_score"], g["no_effect_at_or_above_score"], g["duration_hours"]) == (giants, score, score, 1), slug
        assert v["name"] == "Potion of Giant Strength (" + " or ".join(giants) + ")" and RARITY[v["rarity"]] == rarity and has(b, printed), slug
        tied([g["score_change"]], {"strength-score-becomes-giant-score": "The type of giant determines the score"}, b, "giant score")
        tied([g["appearance"]], {"transparent-liquid-with-floating-sliver-of-light-resembling-giants-fingernail": "transparent liquid has floating in it a sliver of light resembling a giant’s fingernail"}, b, "giant appearance")
    assert RARITY[f["variants"][1]["rarity"]] == "Rare" and has(b, "(frost or 23 Rare stone)")

    d = one["Potion of Growth"]["growth"]
    b = block["Potion of Growth"]
    assert (d["spell"], d["spell_effect"], d["duration_minutes"], d["concentration_required"]) == ("enlarge-reduce", "enlarge", 10, False)
    assert has(b, "you gain the “enlarge” effect of the Enlarge/Reduce spell for 10 minutes (no Concentration required)")
    tied([d["appearance"]], {"red-expands-from-tiny-bead-to-color-clear-liquid-then-contracts-shaking-fails-to-interrupt": "The red in the potion’s liquid continuously expands from a tiny bead to color the clear liquid around it and then contracts. Shaking the bottle fails to interrupt this process"}, b, "growth appearance")

    d = one["Potion of Heroism"]["heroism"]
    b = block["Potion of Heroism"]
    assert (d["temporary_hit_points"], d["duration_hours"]) == (10, 1) and has(b, "you gain 10 Temporary Hit Points that last for 1 hour")
    assert d["spell"] == "bless" and d["concentration_required"] is False
    tied([d["spell_duration"]], {"same-as-temporary-hit-points": "For the same duration, you are under the effect of the Bless spell (no Concentration required)"}, b, "heroism spell")
    tied([d["appearance"]], {"blue-liquid-bubbles-and-steams-as-if-boiling": "blue liquid bubbles and steams as if boiling"}, b, "heroism appearance")

    d = one["Potion of Invisibility"]["invisibility"]
    b = block["Potion of Invisibility"]
    assert d["condition"] == "invisible" and d["duration_hours"] == 1 and has(b, "you have the Invisible condition for 1 hour")
    tied(d["ends_early_on"], {"make-attack-roll": "The effect ends early if you make an attack roll", "deal-damage": "deal damage", "cast-spell": "or cast a spell"}, b, "invisibility ends")
    tied([d["appearance"]], {"container-looks-empty-but-feels-as-though-it-holds-liquid": "container looks empty but feels as though it holds liquid"}, b, "invisibility appearance")

    d = one["Potion of Invulnerability"]["invulnerability"]
    b = block["Potion of Invulnerability"]
    assert d["duration_minutes"] == 1 and has(b, "For 1 minute after you drink this potion")
    tied([d["resistance"]], {"all-damage": "you have Resistance to all damage"}, b, "invulnerability resistance")
    tied([d["appearance"]], {"syrupy-liquid-looks-like-liquefied-iron": "syrupy liquid looks like liquefied iron"}, b, "invulnerability appearance")

    d = one["Potion of Longevity"]["longevity"]
    b = block["Potion of Longevity"]
    roll = {"dice": 1, "sides": 6, "bonus": 6}
    assert d["age_reduction_years"] == roll and d["minimum_age_years"] == 13
    assert has(b, "your physical age is reduced by 1d6 + 6 years, to a minimum of 13 years")
    assert d["aging_chance_percent"] == 10 and d["aging_chance_cumulative"] is True and d["aging_instead_years"] == roll
    assert has(b, "there is 10 percent cumulative chance that you instead age by 1d6 + 6 years")
    tied([d["aging_chance_trigger"]], {"each-subsequent-potion-of-longevity": "Each time you subsequently drink a Potion of Longevity"}, b, "longevity trigger")
    tied([d["appearance"]], {"amber-liquid-with-suspended-tiny-beating-heart-that-vanishes-when-opened": "Suspended in this amber liquid is a tiny heart that, against all reason, is still beating. These ingredients vanish when the potion is opened"}, b, "longevity appearance")

    d = one["Potion of Mind Reading"]["mind_reading"]
    b = block["Potion of Mind Reading"]
    assert (d["spell"], d["save_dc"], d["duration_minutes"], d["concentration_required"]) == ("detect-thoughts", 13, 10, False)
    assert has(b, "you gain the effect of the Detect Thoughts spell (save DC 13) for 10 minutes (no Concentration required)")
    tied([d["appearance"]], {"dense-purple-liquid-with-floating-ovoid-pink-cloud": "dense, purple liquid has an ovoid cloud of pink floating in it"}, b, "mind appearance")

    print("#237 pinned PDF, twelve identities/pages/headers, 16 variants, the giant table and source values pass")


if __name__ == "__main__":
    main()
