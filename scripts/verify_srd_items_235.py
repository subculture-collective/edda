#!/usr/bin/env python3
"""Compare issue #235 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
# The twelve families are on pp. 230-233; the Mysterious Deck continues from p. 231
# onto p. 232, where its Avatar of Death stat block and Necklace of Adaptation follow.
PAGES = (230, 231, 232, 233)
RARITY = {"uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# Printed header (category, rarity and attunement) after each family name.
HEADER = {
    "manual-of-quickness-of-action": "Wondrous Item, Very Rare",
    "marvelous-pigments": "Wondrous Item, Very Rare",
    "medallion-of-thoughts": "Wondrous Item, Uncommon (Requires Attunement)",
    "mirror-of-life-trapping": "Wondrous Item, Very Rare",
    "mithral-armor": "Armor (Any Medium or Heavy, Except Hide Armor), Uncommon",
    "mysterious-deck": "Wondrous Item, Legendary",
    "necklace-of-adaptation": "Wondrous Item, Uncommon (Requires Attunement)",
    "necklace-of-fireballs": "Wondrous Item, Rare",
    "necklace-of-prayer-beads": "Wondrous Item, Rare (Requires Attunement by a Cleric, Druid, or Paladin)",
    "nine-lives-stealer": "Weapon (Any Simple or Martial), Very Rare (Requires Attunement)",
    "oathbow": "Weapon (Longbow or Shortbow), Very Rare (Requires Attunement)",
    "oil-of-etherealness": "Potion, Rare",
}
# The next heading after the last family on p. 233.
FOLLOWING = "Oil of Sharpness Potion"


def text(pdf, page, x):
    args = ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x), "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"]
    return subprocess.run(args, check=True, capture_output=True, text=True).stdout


def clean(raw):
    return FOOTER.sub(" ", " ".join(raw.split()))


def flat(value):
    """Ignore spaces and hyphens so a word split across a line end still matches."""
    return re.sub(r"[\s\-]", "", value)


def has(source, *phrases):
    return all(flat(p) in flat(source) for p in phrases)


def need(source, label, *phrases):
    for p in phrases:
        assert has(source, p), f"{label}: {p}"


def tied(values, phrases, source, label):
    """Require the catalogued keys to be exactly the mapped ones, each with its printed phrase."""
    assert list(values) == list(phrases), f"{label}: {values}"
    for value, phrase in phrases.items():
        assert has(source, phrase), f"{label} {value}: {phrase}"


def roll(r):
    return (r["dice"], r["sides"], r["bonus"])


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_235.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream = " ".join(clean(text(pdf, p, x)) for p in PAGES for x in (0, 306))

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_235.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 235]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 235 and catalogue["batch"] == "data-magic-items-11"
    assert len(families) == 12 and all(len(f["variants"]) == 1 for f in families)
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert [f["presentations"][0]["pages"][0] for f in families] == [230, 230, 230, 230, 231, 231, 232, 233, 233, 233, 233, 233]
    names = [f["name"] for f in families]
    assert names == ["Manual of Quickness of Action", "Marvelous Pigments", "Medallion of Thoughts", "Mirror of Life Trapping", "Mithral Armor", "Mysterious Deck",
                     "Necklace of Adaptation", "Necklace of Fireballs", "Necklace of Prayer Beads", "Nine Lives Stealer", "Oathbow", "Oil of Etherealness"]
    starts = []
    for f in families:
        slug = f["id"].rsplit("/", 1)[1]
        v = f["variants"][0]
        header = HEADER[slug]
        assert header.endswith(("Uncommon", "Rare", "Legendary")) or "(Requires Attunement" in header
        assert RARITY[v["rarity"]] in header and (("Requires Attunement" in header) == (v["attunement"]["mode"] == "required")), f["name"]
        assert (f["category"] == "weapon") == header.startswith("Weapon") and (f["category"] == "armor") == header.startswith("Armor") and (f["category"] == "potion") == header.startswith("Potion")
        printed = f"{f['name']} {header}"
        assert stream.count(printed) == 1, printed
        starts.append(stream.index(printed))
    assert starts == sorted(starts), "printed order differs from the catalogue order"
    ends = starts[1:] + [stream.index(FOLLOWING)]
    block = {f["name"]: stream[s:e] for f, s, e in zip(families, starts, ends)}
    detail = {f["name"]: f["variants"][0]["details"] for f in families}
    charges = {f["name"]: f["variants"][0]["charges"] for f in families}
    none = {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}
    for f in families:
        page = f["presentations"][0]["pages"]
        assert f["variants"][0]["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}, f["name"]
        assert f["variants"][0]["attunement"]["restrictions"] == [], f["name"]
    # Base items printed in the category lines.
    base = {f["name"]: (f["base"]["options"], f["base"]["excluded"]) for f in families}
    assert base["Mithral Armor"] == (["any-medium-armor", "any-heavy-armor"], ["hide-armor"])
    assert base["Nine Lives Stealer"] == (["any-simple-weapon", "any-martial-weapon"], [])
    assert base["Oathbow"] == (["longbow", "shortbow"], [])
    assert all(base[n] == ([], []) for n in names if n not in ("Mithral Armor", "Nine Lives Stealer", "Oathbow"))

    d = detail["Manual of Quickness of Action"]["manual_of_quickness_of_action"]
    b = block["Manual of Quickness of Action"]
    assert (d["study_hours"], d["study_within_days"]) == (48, 6) and has(b, "If you spend 48 hours over a period of 6 days or fewer studying the book’s contents")
    assert d["ability"] == "dexterity" and (d["increase"], d["max_score"]) == (2, 30) and has(b, "your Dexterity increases by 2, to a maximum of 30")
    assert d["loses_magic_after_use"] is True and d["regains_magic_after_years"] == 100 and has(b, "The manual then loses its magic but regains it in a century")
    assert charges["Manual of Quickness of Action"] == none

    d = detail["Marvelous Pigments"]["marvelous_pigments"]
    b = block["Marvelous Pigments"]
    tied([d["container"]], {"fine-wooden-box": "This fine wooden box"}, b, "pigments box")
    assert roll(d["pots"]) == (1, 4, 0) and has(b, "contains 1d4 pots of pigment and a brush")
    assert d["pots_and_brush_weight_pounds"] == 1 and has(b, "(weighing 1 pound in total)")
    assert d["pots_expended_per_painting"] == 1 and has(b, "Using the brush and expending 1 pot of pigment")
    assert d["cube_feet"] == 20 and has(b, "confined to a 20-foot Cube")
    assert d["painting_minutes"] == 10 and has(b, "The effort takes 10 minutes")
    assert d["requires_concentration"] is True and has(b, "and requires Concentration")
    assert d["must_remain_in_cube"] is True and has(b, "you must remain in the Cube")
    tied(d["failure_triggers"], {"concentration-broken": "If your Concentration is broken", "leave-cube-before-work-done": "or you leave the Cube before the work is done"}, b, "pigments failure")
    tied([d["failure_result"]], {"all-painted-elements-vanish-and-pot-is-wasted": "all the painted elements vanish, and the pot of pigment is wasted"}, b, "pigments failure result")
    tied([d["completion_result"]], {"painted-objects-and-terrain-features-become-real": "all the painted objects and terrain features become real"}, b, "pigments completion")
    assert (d["object_max_value_gp"], d["total_max_value_per_pot_gp"]) == (25, 500)
    need(b, "pigments value", "can have a value greater than 25 GP", "can’t exceed 500 GP")
    tied([d["over_value_result"]], {"look-authentic-but-are-paste-cookies-or-other-worthless-material": "they look authentic, but close inspection reveals they’re made from paste, cookies, or some other worthless material"}, b, "pigments over value")
    tied([d["energy_result"]], {"energy-dissipates-when-painting-completes-doing-no-harm": "the energy dissipates as soon as you complete the painting, doing no harm"}, b, "pigments energy")
    assert charges["Marvelous Pigments"] == none and " GP" in b

    d = detail["Medallion of Thoughts"]["medallion_of_thoughts"]
    b = block["Medallion of Thoughts"]
    assert charges["Medallion of Thoughts"] == {"mode": "limited", "maximum": 5, "per_use": 1, "recharge": "regain-1d4-at-dawn"} and d["max_charges"] == 5 and has(b, "The medallion has 5 charges")
    tied([d["condition"]], {"while-wearing": "While wearing it"}, b, "medallion condition")
    assert d["charges_per_cast"] == 1 and d["spell"] == "detect-thoughts" and d["save_dc"] == 13 and has(b, "you can expend 1 charge to cast Detect Thoughts (save DC 13) from it")
    assert roll(d["recharge_amount"]) == (1, 4, 0) and d["recharge_time"] == "daily-at-dawn" and has(b, "The medallion regains 1d4 expended charges daily at dawn")

    d = detail["Mirror of Life Trapping"]["mirror_of_life_trapping"]
    b = block["Mirror of Life Trapping"]
    assert (d["height_feet"], d["width_feet"], d["weight_pounds"]) == (4, 2, 50) and has(b, "this 4-foot-tall, 2-foot-wide mirror", "The mirror weighs 50 pounds")
    assert (d["ac"], d["hp"]) == (11, 10) and has(b, "it has AC 11, HP 10")
    tied(d["immunities"], {"poison": "Immunity to Poison", "psychic": "and Psychic damage"}, b, "mirror immunities")
    tied(d["vulnerabilities"], {"bludgeoning": "Vulnerability to Bludgeoning damage"}, b, "mirror vulnerability")
    tied([d["destroyed_when"]], {"reduced-to-0-hit-points": "It shatters and is destroyed when reduced to 0 Hit Points"}, b, "mirror destroyed")
    assert d["activate_action"] == "magic" and d["activate_range_feet"] == 5 and d["deactivate_action"] == "magic"
    need(b, "mirror activate", "If the mirror is hanging on a vertical surface and you are within 5 feet of it, you can take a Magic action and use a command word to activate it", "until you take a Magic action and repeat the command word to deactivate it")
    tied(d["activate_requires"], {"command-word": "use a command word to activate it", "hanging-on-vertical-surface": "If the mirror is hanging on a vertical surface"}, b, "mirror requires")
    assert d["reflection_range_feet"] == 30 and has(b, "sees its reflection in the activated mirror while within 30 feet of the mirror")
    assert (d["save_ability"], d["save_dc"]) == ("charisma", 15) and has(b, "must succeed on a DC 15 Charisma saving throw")
    tied([d["save_failure"]], {"trapped-with-worn-and-carried-items-in-extradimensional-cell": "trapped, along with anything it is wearing or carrying, in one of the mirror’s twelve extradimensional cells"}, b, "mirror failure")
    tied([d["save_advantage_when"]], {"knows-the-mirrors-nature": "A creature that knows the mirror’s nature makes the save with Advantage"}, b, "mirror advantage")
    tied(d["save_auto_succeed"], {"constructs": "Constructs succeed on the save automatically"}, b, "mirror auto")
    assert d["cells"] == 12
    tied([d["cell_description"]], {"infinite-expanse-of-thick-fog": "An extradimensional cell is an infinite expanse filled with thick fog"}, b, "mirror cell")
    assert d["cell_visibility_feet"] == 10 and has(b, "reduces visibility to 10 feet")
    assert d["trapped_creatures_age"] is False and d["trapped_creatures_need_food_drink_or_sleep"] is False and has(b, "Creatures trapped in the mirror’s cells don’t age, and they don’t need to eat, drink, or sleep")
    tied([d["escape_method"]], {"magic-that-permits-planar-travel": "can escape using magic that permits planar travel"}, b, "mirror escape")
    tied([d["otherwise_confined"]], {"until-freed": "the creature is confined to the cell until freed"}, b, "mirror confined")
    tied([d["full_cells_result"]], {"frees-one-trapped-creature-at-random": "the mirror frees one trapped creature at random to accommodate the new prisoner"}, b, "mirror full")
    tied([d["freed_creature_appears"]], {"unoccupied-space-within-sight-of-mirror-facing-away": "A freed creature appears in an unoccupied space within sight of the mirror but facing away from it"}, b, "mirror freed")
    tied([d["shatter_result"]], {"all-creatures-freed-in-unoccupied-spaces-near-mirror": "If the mirror is shattered, all creatures it contains are freed and appear in unoccupied spaces near it"}, b, "mirror shatter")
    assert d["name_action"] == "magic" and d["name_range_feet"] == 5
    need(b, "mirror name", "While within 5 feet of the mirror, you can take a Magic action to name one creature trapped in it or call out a particular cell by number")
    tied([d["name_result"]], {"named-creature-or-cell-occupant-appears-as-image-and-can-communicate": "appears as an image on the mirror’s surface. You and the creature can then communicate"}, b, "mirror name result")
    assert d["free_action"] == "magic" and d["free_command_word_ordinal"] == 2 and has(b, "use a second command word to free one creature trapped in the mirror")
    tied([d["free_result"]], {"creature-appears-with-possessions-in-nearest-unoccupied-space-facing-away": "The freed creature appears, along with its possessions, in the unoccupied space nearest to the mirror and facing away from it"}, b, "mirror free")
    tied(d["extradimensional_items_destroyed"], {"bag-of-holding": "created by a Bag of Holding", "portable-hole": "Portable Hole", "similar-item": "or similar item"}, b, "mirror items")
    assert has(b, "instantly destroys both items")
    assert d["gate_to"] == "astral-plane" and d["gate_radius_feet"] == 10 and has(b, "opens a gate to the Astral Plane", "within 10 feet of the gate")
    tied([d["gate_affected"]], {"creatures-not-behind-total-cover-sucked-to-random-location": "not behind Total Cover is sucked through it to a random location on the Astral Plane"}, b, "mirror gate affected")
    assert d["gate_after"] == "closes" and d["gate_one_way"] is True and d["gate_can_reopen"] is False and has(b, "The gate then closes. The gate is one-way only and can’t be reopened")
    assert charges["Mirror of Life Trapping"] == none

    d = detail["Mithral Armor"]["mithral_armor"]
    b = block["Mithral Armor"]
    tied([d["wearable"]], {"under-normal-clothes": "Armor made of this substance can be worn under normal clothes"}, b, "mithral wearable")
    tied([d["stealth_disadvantage"]], {"not-imposed": "If the armor normally imposes Disadvantage on Dexterity (Stealth) checks"}, b, "mithral stealth")
    tied([d["strength_requirement"]], {"does-not-apply": "or has a Strength requirement, the mithral version of the armor doesn’t"}, b, "mithral strength")
    assert charges["Mithral Armor"] == none

    d = detail["Mysterious Deck"]["mysterious_deck"]
    b = block["Mysterious Deck"]
    assert d["source_pages"] == [231, 232] and "Avatar of Death Medium Undead" in b
    tied(d["card_materials"], {"ivory": "cards made of ivory", "vellum": "or vellum"}, b, "deck materials")
    assert d["deck_sizes"] == [13, 22] and d["thirteen_card_percent"] == 75 and d["table_die_sides"] == 100
    need(b, "deck sizes", "Most (75 percent) of these decks have thirteen cards, but some have twenty-two", "Use the appropriate column of the Mysterious Deck table")
    need(b, "deck draws", "you must declare how many cards you intend to draw", "Any cards drawn in excess of this number have no effect")
    assert d["declare_draw_count_before_drawing"] is True and d["excess_draws_have_no_effect"] is True
    assert d["max_hours_between_draws"] == 1 and has(b, "You must draw each card no more than 1 hour after the previous draw")
    tied([d["undrawn_cards_result"]], {"remaining-cards-fly-from-deck-and-take-effect-at-once": "the remaining number of cards fly from the deck on their own and take effect all at once"}, b, "deck undrawn")
    assert d["drawn_card_disappears"] is True and has(b, "Once a card is drawn, it disappears")
    tied(d["does_not_reappear"], {"fool": "Unless the card is the Fool", "jester": "or Jester, the card reappears in the deck"}, b, "deck reappear")
    tied(d["reroll_when_gone"], {"fool": "Once the Fool", "jester": "or Jester has left the deck, re-roll on the table if that card comes up again"}, b, "deck reroll")
    # The two columns of the table, parsed from the printed rows (00 is 100).
    table = re.search(r"\(22-Card Card Deck\) Deck\) (.*?) 1d100 1d100 \(13-Card \(22-Card Card Deck\) Deck\) (.*?) Each card’s effect", b)
    rows = re.findall(r"(—|\d\d–\d\d) (—|\d\d–\d\d) ([A-Z][a-z]+)", table.group(1) + " " + table.group(2))
    assert len(rows) == 22, rows
    span = lambda s: None if s == "—" else tuple(100 if x == "00" else int(x) for x in s.split("–"))
    want13 = [{"card": c.lower(), "min": span(a)[0], "max": span(a)[1]} for a, _, c in rows if a != "—"]
    want22 = [{"card": c.lower(), "min": span(z)[0], "max": span(z)[1]} for _, z, c in rows]
    assert d["deck_13_cards"] == want13 and len(want13) == 13, "13-card column"
    assert d["deck_22_cards"] == want22, "22-card column"
    assert [c["card"] for c in d["cards"]] == [r["card"] for r in want22] == sorted(r["card"] for r in want22)
    # Each card has its own paragraph, in table order.
    marks = [b.index(f" {c['card'].capitalize()}. ", b.index("Each card’s effect")) for c in d["cards"]]
    assert marks == sorted(marks), "card paragraphs out of order"
    section = {c["card"]: b[m:e] for c, m, e in zip(d["cards"], marks, marks[1:] + [b.index("Avatar of Death Medium")])}
    card = {c["card"]: c for c in d["cards"]}

    def fact(name, key, amount, unit, *phrases):
        facts = {f["key"]: f for f in card[name]["facts"]}
        assert (facts[key]["amount"], facts[key]["unit"]) == (amount, unit), (name, key)
        need(section[name], f"card {name} {key}", *phrases)

    fact("balance", "increase", 2, "points", "increase one of your ability scores by 2")
    fact("balance", "maximum-score", 22, "score", "to a maximum of 22")
    fact("balance", "decrease", 2, "points", "decrease another one of your ability scores by 2")
    fact("balance", "cannot-decrease-at-or-below-score", 5, "score", "an ability that has a score of 5 or lower")
    fact("comet", "advantage-duration", 1, "year", "Advantage on Death Saving Throws for 1 year")
    fact("euryale", "saving-throw-penalty", 2, "minus", "You take a −2 penalty to saving throws")
    fact("fool", "duration", 72, "hours", "Disadvantage on D20 Tests for the next 72 hours")
    fact("fool", "additional-draws", 1, "cards", "Draw another card; this draw doesn’t count as one of your declared draws")
    fact("gem", "jewelry-pieces", 25, "pieces", "Twenty-five pieces of jewelry")
    fact("gem", "jewelry-value-each", 2000, "gp", "worth 2,000 GP each")
    fact("gem", "gems", 50, "pieces", "fifty gems")
    fact("gem", "gem-value-each", 1000, "gp", "worth 1,000 GP each appear at your feet")
    fact("jester", "duration", 72, "hours", "Advantage on D20 Tests for the next 72 hours")
    fact("jester", "additional-draws", 2, "cards", "draw two additional cards beyond your declared draws")
    fact("knight", "appears-within", 30, "feet", "within 30 feet of yourself")
    fact("moon", "wish-casts", 1, "d3", "cast Wish 1d3 times")
    fact("puzzle", "reduction", 1, "d4", "Permanently reduce your Intelligence or Wisdom by 1d4 + 1")
    fact("puzzle", "reduction-bonus", 1, "points", "by 1d4 + 1")
    fact("puzzle", "minimum-score", 1, "score", "(to a minimum score of 1)")
    fact("puzzle", "additional-draws", 1, "cards", "You can draw one additional card beyond your declared draws")
    fact("sage", "window", 1, "year", "within one year of drawing this card")
    fact("star", "increase", 2, "points", "Increase one of your ability scores by 2")
    fact("star", "maximum-score", 24, "score", "to a maximum of 24")
    fact("sun", "temporary-hit-points", 10, "hit-points", "you gain 10 Temporary Hit Points daily at dawn until you die")
    assert [c["card"] for c in d["cards"] if c["facts"]] == ["balance", "comet", "euryale", "fool", "gem", "jester", "knight", "moon", "puzzle", "sage", "star", "sun"]
    need(section["balance"], "balance", "you can choose not to adjust your ability scores, in which case this card has no effect")
    need(section["comet"], "comet", "you can select one of them as your foe when you roll Initiative", "If someone else reduces your chosen foe to 0 Hit Points or you don’t choose a foe, this card has no effect")
    need(section["donjon"], "donjon", "entombed in a state of suspended animation in an extradimensional sphere", "except for Artifacts, which stay behind in the space you occupied", "You remain imprisoned until you are found and removed from the sphere", "You can’t be located by any Divination magic, but a Wish spell can reveal the location of your prison", "You draw no more cards")
    need(section["euryale"], "euryale", "medusa-like visage curses you", "Only a god or the magic of the Fates card can end this curse")
    need(section["fates"], "fates", "avoid or erase one event as if it never happened", "as soon as you draw the card or at any other time before you die")
    need(section["flames"], "flames", "A powerful devil becomes your enemy", "torments you, savoring your suffering before attempting to slay you", "until either you or the devil dies")
    need(section["key"], "key", "A Rare or rarer magic weapon with which you are proficient appears on your person", "The GM chooses the weapon")
    need(section["knight"], "knight", "You gain the service of a Knight", "The knight has the same alignment as you and serves you loyally until death", "Work with your GM to create a name and backstory", "The GM can use a different stat block")
    need(section["rogue"], "rogue", "An NPC of the GM’s choice becomes Hostile toward you", "You don’t know the identity of this NPC until they or someone else reveals it", "Nothing less than a Wish spell or divine intervention can end the NPC’s hostility")
    need(section["ruin"], "ruin", "All forms of wealth that you carry or own, other than magic items, are lost to you", "Portable property vanishes", "Businesses, buildings, and land you own are lost in a way that alters reality the least", "Any documentation that proves you should own something lost to this card also disappears")
    need(section["sage"], "sage", "ask a question in meditation and mentally receive a truthful answer")
    need(section["skull"], "skull", "An Avatar of Death (see the accompanying stat block) appears in an unoccupied space as close to you as possible", "The avatar targets only you with its attacks", "a ghostly skeleton clad in a tattered black robe and carrying a spectral scythe", "The avatar disappears when it drops to 0 Hit Points or you die", "If an ally of yours deals damage to the avatar, that ally summons another Avatar of Death", "targets only that ally with its attacks", "You and your allies can each summon only one avatar as a consequence of this draw", "A creature slain by an avatar can’t be restored to life")
    need(section["sun"], "sun", "A magic item (chosen by the GM) appears on your person")
    need(section["talons"], "talons", "Every magic item you wear or carry disintegrates", "Artifacts in your possession vanish instead")
    need(section["throne"], "throne", "proficiency and Expertise in your choice of History, Insight, Intimidation, or Persuasion", "rightful ownership of a small keep", "the keep is currently home to one or more monsters, which must be cleared out before you can claim the keep as yours")
    need(section["void"], "void", "Your soul is drawn from your body and contained in an object in a place of the GM’s choice", "One or more powerful beings guard the place", "your body is inert, ceases aging, and requires no food, air, or water", "A Wish spell can’t return your soul to your body, but the spell reveals the location of the object that holds your soul", "You draw no more cards")
    a = d["avatar_of_death"]
    need(b, "avatar", "Avatar of Death Medium Undead, Neutral evil AC 20 Initiative +3 (13) HP Half the HP maximum of its summoner Speed 60 ft., Fly 60 ft. (hover)")
    assert (a["size"], a["creature_type"], a["alignment"], a["ac"], a["initiative_modifier"], a["initiative_score"]) == ("medium", "undead", "neutral-evil", 20, 3, 13)
    assert a["hit_points"] == "half-the-hp-maximum-of-its-summoner" and (a["speed_feet"], a["fly_speed_feet"], a["hover"]) == (60, 60, True)
    assert a["abilities"] == [{"ability": n, "score": 16, "mod": 3, "save": 3} for n in ("strength", "dexterity", "constitution", "intelligence", "wisdom", "charisma")]
    need(b, "avatar abilities", "Str 16 +3 +3 Dex 16 +3 +3 Con 16 +3 +3 Int 16 +3 +3 Wis 16 +3 +3 Cha 16 +3 +3")
    assert a["damage_immunities"] == ["necrotic", "poison"] and a["condition_immunities"] == ["charmed", "exhaustion", "frightened", "paralyzed", "petrified", "poisoned", "unconscious"]
    need(b, "avatar immunities", "Immunities Necrotic, Poison; Charmed, Exhaustion, Frightened, Paralyzed, Petrified, Poisoned, Unconscious")
    assert (a["truesight_feet"], a["passive_perception"]) == (60, 13) and has(b, "Senses Truesight 60 ft., Passive Perception 13")
    assert a["languages"] == "all-languages-known-to-its-summoner" and has(b, "Languages All languages known to its summoner")
    assert (a["challenge_rating"], a["xp"], a["proficiency_bonus"]) == ("none", 0, "equals-its-summoners") and has(b, "CR None (XP 0; PB equals its summoner’s)")
    assert a["incorporeal_movement"] == "move-through-creatures-and-objects-as-difficult-terrain" and has(b, "The avatar can move through other creatures and objects as if they were Difficult Terrain")
    assert a["incorporeal_damage_ending_turn_inside_object"] == {"average": 5, "dice": 1, "sides": 10, "bonus": 0, "type": "force"} and has(b, "It takes 5 (1d10) Force damage if it ends its turn inside an object")
    assert a["multiattack_count"] == "half-the-summoners-proficiency-bonus-rounded-up" and a["multiattack_attack"] == "reaping-scythe"
    need(b, "avatar multiattack", "The avatar makes a number of Reaping Scythe attacks equal to half the summoner’s Proficiency Bonus (rounded up)")
    assert a["reaping_scythe_attack"] == "melee-attack-roll" and a["reaping_scythe_automatic_hit"] is True and a["reaping_scythe_reach_feet"] == 5
    need(b, "avatar scythe", "Reaping Scythe. Melee Attack Roll: Automatic hit, reach 5 ft.")
    assert a["reaping_scythe_damage"] == {"average": 7, "dice": 1, "sides": 8, "bonus": 3, "type": "slashing"} and a["reaping_scythe_extra_damage"] == {"average": 4, "dice": 1, "sides": 8, "bonus": 0, "type": "necrotic"}
    need(b, "avatar damage", "Hit: 7 (1d8 + 3) Slashing damage plus 4 (1d8) Necrotic damage")
    assert charges["Mysterious Deck"] == none

    d = detail["Necklace of Adaptation"]["necklace_of_adaptation"]
    b = block["Necklace of Adaptation"]
    tied([d["condition"]], {"while-wearing": "While wearing this necklace"}, b, "adaptation condition")
    tied([d["breathe_normally_in"]], {"any-environment": "you can breathe normally in any environment"}, b, "adaptation breathe")
    tied(d["advantage_on_saves_to"], {"avoid": "saving throws made to avoid", "end": "or end the Poisoned condition"}, b, "adaptation saves")
    assert d["advantage_against_condition"] == "poisoned" and has(b, "you have Advantage on saving throws made to avoid or end the Poisoned condition")
    assert charges["Necklace of Adaptation"] == none

    d = detail["Necklace of Fireballs"]["necklace_of_fireballs"]
    b = block["Necklace of Fireballs"]
    assert roll(d["beads"]) == (1, 6, 3) and has(b, "This necklace has 1d6 + 3 beads hanging from it")
    assert d["detach_action"] == "magic" and d["throw_range_feet"] == 60 and has(b, "You can take a Magic action to detach a bead and throw it up to 60 feet away")
    assert (d["spell"], d["spell_level"], d["save_dc"]) == ("fireball", 3, 15) and has(b, "the bead detonates as a level 3 Fireball (save DC 15)")
    tied([d["multiple_beads"]], {"can-hurl-multiple-beads-or-whole-necklace-at-one-time": "You can hurl multiple beads, or even the whole necklace, at one time"}, b, "fireballs multiple")
    assert roll(d["extra_damage_per_bead_after_first"]) == (1, 6, 0) and roll(d["maximum_damage"]) == (12, 6, 0)
    need(b, "fireballs damage", "increase the damage of the Fireball by 1d6 for each bead after the first (maximum 12d6)")
    assert charges["Necklace of Fireballs"] == none

    d = detail["Necklace of Prayer Beads"]["necklace_of_prayer_beads"]
    b = block["Necklace of Prayer Beads"]
    tied(d["attunement_by_any_of"], {"cleric": "Requires Attunement by a Cleric", "druid": "Druid", "paladin": "or Paladin"}, b, "beads attunement")
    assert roll(d["magic_beads"]) == (1, 4, 2) and has(b, "This necklace has 1d4 + 2 magic beads")
    tied(d["magic_bead_materials"], {"aquamarine": "made from aquamarine", "black-pearl": "black pearl", "topaz": "or topaz"}, b, "beads materials")
    tied(d["nonmagical_bead_examples"], {"amber": "stones such as amber", "bloodstone": "bloodstone", "citrine": "citrine", "coral": "coral", "jade": "jade", "pearl": "pearl", "quartz": "or quartz"}, b, "beads nonmagical")
    assert d["removed_bead_loses_magic"] is True and has(b, "If a magic bead is removed from the necklace, that bead loses its magic")
    assert d["bead_types"] == 6 and has(b, "Six types of magic beads exist")
    tied(d["type_selection"], {"gm-decides": "The GM decides the type of each bead on the necklace", "roll-on-table": "determines it randomly by rolling on the table below"}, b, "beads selection")
    assert d["necklace_may_have_repeated_types"] is True and has(b, "A necklace can have more than one bead of the same type")
    assert d["requires_wearing"] is True and has(b, "To use one, you must be wearing the necklace")
    assert d["cast_action"] == "bonus" and has(b, "a spell that you can cast from it as a Bonus Action")
    tied([d["save_dc"]], {"your-spell-save-dc": "using your spell save DC if a save is necessary"}, b, "beads dc")
    assert d["uses_per_bead"] == 1 and d["bead_recharge"] == "next-dawn" and has(b, "Once a magic bead’s spell is cast, that bead can’t be used again until the next dawn")
    assert d["table_die_sides"] == 20 and has(b, "1d20 Bead Spell")
    table = re.search(r"1d20 Bead Spell (.*)$", b).group(1)
    printed = re.findall(r"(\d+)(?:–(\d+))? (Bead of [A-Za-z]+)", table)
    # The last name wraps inside its table cell: "Bead of Wind" and, below, "Walking".
    printed = [(lo, hi, name + " Walking" if name == "Bead of Wind" else name) for lo, hi, name in printed]
    assert has(table, "Wind Walk Walking")
    got = [(r["min"], r["max"], r["bead"]) for r in d["beads_table"]]
    assert got == [(int(lo), int(hi or lo), name.lower().replace(" ", "-")) for lo, hi, name in printed], printed
    assert [r["spell"] for r in d["beads_table"]] == ["bless", "cure-wounds", "greater-restoration", "shining-smite", "guardian-of-faith", "wind-walk"]
    need(table, "beads spells", "Bead of Blessing Bless", "Bead of Curing Cure Wounds (level 2 version)", "Bead of Favor Greater Restoration", "Bead of Smiting Shining Smite", "Bead of Summons Guardian of Faith", "Bead of Wind Wind Walk Walking")
    assert [r["spell_level"] for r in d["beads_table"]] == [0, 2, 0, 0, 0, 0]
    assert charges["Necklace of Prayer Beads"] == none

    d = detail["Nine Lives Stealer"]["nine_lives_stealer"]
    b = block["Nine Lives Stealer"]
    assert (d["attack_bonus"], d["damage_bonus"]) == (2, 2) and has(b, "You gain a +2 bonus to attack rolls and damage rolls made with this magic weapon")
    assert d["property"] == "life-stealing" and has(b, "Life Stealing.")
    assert roll(d["charges"]) == (1, 8, 1) and charges["Nine Lives Stealer"] == {"mode": "limited", "maximum": 9, "per_use": 1, "recharge": "none"} and has(b, "The weapon has 1d8 + 1 charges")
    assert d["target_hit_points_below"] == 100 and d["trigger_d20_roll"] == 20 and has(b, "a creature that has fewer than 100 Hit Points with this weapon and roll a 20 on the d20 for the attack roll")
    assert (d["save_ability"], d["save_dc"]) == ("constitution", 15) and has(b, "must succeed on a DC 15 Constitution saving throw")
    tied([d["save_failure"]], {"slain-instantly": "be slain instantly as the sword tears its life force from its body"}, b, "nine lives failure")
    tied(d["save_auto_succeed"], {"constructs": "Constructs", "undead": "and Undead succeed on the save automatically"}, b, "nine lives auto")
    assert d["charges_lost"] == 1 and d["charge_lost_when"] == "creature-is-slain" and has(b, "The weapon loses 1 charge if the creature is slain")
    tied([d["when_no_charges"]], {"weapon-loses-this-property": "When the weapon has no charges remaining, it loses this property"}, b, "nine lives empty")

    d = detail["Oathbow"]["oathbow"]
    b = block["Oathbow"]
    tied([d["nocked_effect"]], {"whispers-in-elvish-swift-defeat-to-my-enemies": "When you nock an arrow on this bow, it whispers in Elvish, “Swift defeat to my enemies.”"}, b, "oathbow nock")
    tied([d["command_trigger"]], {"ranged-attack-with-this-weapon": "When you use this weapon to make a ranged attack"}, b, "oathbow trigger")
    tied(d["command_methods"], {"utter": "you can utter", "sign": "or sign the following command words"}, b, "oathbow methods")
    need(b, "oathbow words", "“Swift death to you who have wronged me.”")
    tied([d["sworn_enemy"]], {"target-of-the-attack": "The target of your attack becomes your sworn enemy"}, b, "oathbow enemy")
    tied(d["sworn_enemy_ends"], {"until-it-dies": "until it dies", "until-dawn-7-days-later": "or until dawn 7 days later"}, b, "oathbow ends")
    assert d["sworn_enemy_days"] == 7 and d["max_sworn_enemies"] == 1 and has(b, "You can have only one such sworn enemy at a time")
    tied([d["new_enemy_when"]], {"after-the-next-dawn-once-sworn-enemy-dies": "When your sworn enemy dies, you can choose a new one after the next dawn"}, b, "oathbow new")
    assert d["ranged_attack_advantage_against_sworn_enemy"] is True and has(b, "When you make a ranged attack roll with this weapon against your sworn enemy, you have Advantage on the roll")
    tied(d["cover_ignored"], {"half-cover": "your target gains no benefit from Half Cover", "three-quarters-cover": "or Three-Quarters Cover"}, b, "oathbow cover")
    assert d["no_long_range_disadvantage"] is True and has(b, "you suffer no Disadvantage due to long range")
    assert d["extra_damage"] == {"dice": 3, "sides": 6, "type": "piercing"} and has(b, "your sworn enemy takes an extra 3d6 Piercing damage")
    tied([d["extra_damage_when"]], {"attack-hits-sworn-enemy": "If the attack hits, your sworn enemy takes an extra"}, b, "oathbow damage")
    assert d["other_weapons"] == "disadvantage-on-attack-rolls" and d["other_weapons_while"] == "sworn-enemy-lives"
    need(b, "oathbow other", "While your sworn enemy lives, you have Disadvantage on attack rolls with all other weapons")
    assert charges["Oathbow"] == none

    d = detail["Oil of Etherealness"]["oil_of_etherealness"]
    b = block["Oil of Etherealness"]
    tied([d["covers"]], {"one-medium-or-smaller-creature-with-worn-and-carried-equipment": "One vial of this oil can cover one Medium or smaller creature, along with the equipment it’s wearing and carrying"}, b, "oil covers")
    assert d["extra_vials_per_size_above_medium"] == 1 and has(b, "one additional vial is required for each size category above Medium")
    assert d["apply_minutes"] == 10 and has(b, "Applying the oil takes 10 minutes")
    assert d["spell"] == "etherealness" and d["duration_hours"] == 1 and has(b, "gains the effect of the Etherealness spell for 1 hour")
    tied([d["appearance"]], {"cloudy-gray-oil-beads-form-outside-container-and-quickly-evaporate": "Beads of this cloudy, gray oil form on the outside of its container and quickly evaporate"}, b, "oil appearance")
    assert charges["Oil of Etherealness"] == none

    for name in names:
        if name not in ("Medallion of Thoughts", "Nine Lives Stealer"):
            assert charges[name] == none, name
    assert " GP" not in " ".join(block[n] for n in names if n not in ("Marvelous Pigments", "Mysterious Deck")), "a price is printed"
    print("#235 pinned PDF, twelve identities/pages/headers, 12 variants, deck and bead tables, charges and source values pass")


if __name__ == "__main__":
    main()
