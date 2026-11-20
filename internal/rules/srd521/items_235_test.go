package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func loadItems235(t *testing.T) ItemSourceBatch[ItemDetails235] {
	t.Helper()
	b, err := LoadItems235()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 230-233, in planning inventory order for issue #235. Every
// family prints one form, so each has the single variant "standard".
func TestItems235SourceIdentities(t *testing.T) {
	expected := []struct {
		name, slug string
		page       int
		category   string
		base       []string
		excluded   []string
		attunement string
		rarity     string
	}{
		{"Manual of Quickness of Action", "manual-of-quickness-of-action", 230, "wondrous-item", nil, nil, "none", "very-rare"},
		{"Marvelous Pigments", "marvelous-pigments", 230, "wondrous-item", nil, nil, "none", "very-rare"},
		{"Medallion of Thoughts", "medallion-of-thoughts", 230, "wondrous-item", nil, nil, "required", "uncommon"},
		{"Mirror of Life Trapping", "mirror-of-life-trapping", 230, "wondrous-item", nil, nil, "none", "very-rare"},
		{"Mithral Armor", "mithral-armor", 231, "armor", []string{"any-medium-armor", "any-heavy-armor"}, []string{"hide-armor"}, "none", "uncommon"},
		{"Mysterious Deck", "mysterious-deck", 231, "wondrous-item", nil, nil, "none", "legendary"},
		{"Necklace of Adaptation", "necklace-of-adaptation", 232, "wondrous-item", nil, nil, "required", "uncommon"},
		{"Necklace of Fireballs", "necklace-of-fireballs", 233, "wondrous-item", nil, nil, "none", "rare"},
		{"Necklace of Prayer Beads", "necklace-of-prayer-beads", 233, "wondrous-item", nil, nil, "required", "rare"},
		{"Nine Lives Stealer", "nine-lives-stealer", 233, "weapon", []string{"any-simple-weapon", "any-martial-weapon"}, nil, "required", "very-rare"},
		{"Oathbow", "oathbow", 233, "weapon", []string{"longbow", "shortbow"}, nil, "required", "very-rare"},
		{"Oil of Etherealness", "oil-of-etherealness", 233, "potion", nil, nil, "none", "rare"},
	}
	b := loadItems235(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-11" || b.CatalogueIssue != 235 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #235")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 230–233") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			if f.ID != id || f.Name != want.name || f.Kind != "magic" || f.Category != want.category ||
				!slices.Equal(f.Base.Options, append([]string{}, want.base...)) || !slices.Equal(f.Base.Excluded, append([]string{}, want.excluded...)) {
				t.Fatalf("source identity mismatch: %+v", f)
			}
			if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/magic-items/"+want.slug || !slices.Equal(f.Presentations[0].Pages, []int{want.page, want.page}) {
				t.Fatal("source presentation mismatch")
			}
			if len(f.Variants) != 1 {
				t.Fatalf("the source prints one form of %s, got %d variants", want.name, len(f.Variants))
			}
			v := f.Variants[0]
			if v.ID != id+"/variant/standard" || v.Name != want.name || v.Rarity != want.rarity || v.Attunement.Mode != want.attunement || len(v.Attunement.Restrictions) != 0 {
				t.Fatalf("variant mismatch: %+v", v)
			}
			if v.Cost.Kind != "unspecified" || v.Cost.Amount != 0 || v.Cost.Currency != "none" || !slices.Equal(v.Cost.SourcePages, []int{want.page, want.page}) {
				t.Fatal("no price is listed; cost must stay unspecified")
			}
			// Medallion of Thoughts (5, regains 1d4 at dawn) and Nine Lives Stealer
			// (1d8 + 1, so at most 9, never recharged) are the only charge pools.
			wantCharges := ItemCharges{Mode: "none", Recharge: "none"}
			switch want.slug {
			case "medallion-of-thoughts":
				wantCharges = ItemCharges{Mode: "limited", Maximum: 5, PerUse: 1, Recharge: "regain-1d4-at-dawn"}
			case "nine-lives-stealer":
				wantCharges = ItemCharges{Mode: "limited", Maximum: 9, PerUse: 1, Recharge: "none"}
			}
			if *v.Charges != wantCharges {
				t.Fatalf("charges mismatch: %+v", *v.Charges)
			}
			if v.MechanicsAvailable() {
				t.Fatal("catalogue record became available")
			}
		})
	}
	// Join every family to the decision inventory the issue cites, in its order.
	data, err := os.ReadFile("../../../docs/investigations/issue-80-spell-item-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Records []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Family string `json:"family"`
			Page   int    `json:"source_page"`
			Batch  string `json:"batch_key"`
			Issue  int    `json:"catalogue_issue"`
			Status string `json:"status"`
		} `json:"records"`
		Batches []struct {
			Key     string `json:"key"`
			Pages   string `json:"pages"`
			Members []struct {
				Name string `json:"name"`
				Page int    `json:"source_page"`
			} `json:"members"`
		} `json:"batches"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	index := 0
	for _, r := range inventory.Records {
		if r.Issue != 235 {
			continue
		}
		if index >= len(b.Families) {
			t.Fatal("inventory lists more families than the batch")
		}
		f := b.Families[index]
		if r.Family != "magic-items" || r.Batch != b.Batch || r.Status != "unsupported" || r.ID != f.Presentations[0].ID || r.Name != f.Name || r.Page != f.Presentations[0].Pages[0] {
			t.Fatalf("inventory record %s disagrees with batch family %s", r.ID, f.Name)
		}
		index++
	}
	if index != 12 {
		t.Fatalf("inventory assigns %d families to issue #235", index)
	}
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "230, 231, 232, 233" || len(batch.Members) != 12 {
			t.Fatalf("planning batch pages or member count changed: %q, %d", batch.Pages, len(batch.Members))
		}
		for i, m := range batch.Members {
			if m.Name != expected[i].name || m.Page != expected[i].page {
				t.Fatalf("planning batch member %d is %s p.%d", i, m.Name, m.Page)
			}
		}
	}
	if planned != 1 {
		t.Fatalf("expected one planning batch %s, got %d", b.Batch, planned)
	}
	b.Families[0].Variants[0].Details.Manual.Ability = "changed"
	if fresh := loadItems235(t); fresh.Families[0].Variants[0].Details.Manual.Ability != "dexterity" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all twelve variants, transcribed from pinned PDF
// pp. 230-233. The Mysterious Deck continues from p. 231 to p. 232 and its
// Avatar of Death stat block is on p. 232. A printed 00 is recorded as 100;
// the Euryale penalty is an amount with the unit "minus".
func TestItems235SourceValues(t *testing.T) {
	expected := []string{
		`{"manual_of_quickness_of_action":{"study_hours":48,"study_within_days":6,"ability":"dexterity","increase":2,"max_score":30,"loses_magic_after_use":true,"regains_magic_after_years":100}}`,
		`{"marvelous_pigments":{"container":"fine-wooden-box","pots":{"dice":1,"sides":4,"bonus":0},"pots_and_brush_weight_pounds":1,"pots_expended_per_painting":1,"cube_feet":20,"painting_minutes":10,"requires_concentration":true,"must_remain_in_cube":true,"failure_triggers":["concentration-broken","leave-cube-before-work-done"],"failure_result":"all-painted-elements-vanish-and-pot-is-wasted","completion_result":"painted-objects-and-terrain-features-become-real","object_max_value_gp":25,"total_max_value_per_pot_gp":500,"over_value_result":"look-authentic-but-are-paste-cookies-or-other-worthless-material","energy_result":"energy-dissipates-when-painting-completes-doing-no-harm"}}`,
		`{"medallion_of_thoughts":{"condition":"while-wearing","max_charges":5,"charges_per_cast":1,"spell":"detect-thoughts","save_dc":13,"recharge_amount":{"dice":1,"sides":4,"bonus":0},"recharge_time":"daily-at-dawn"}}`,
		`{"mirror_of_life_trapping":{"height_feet":4,"width_feet":2,"weight_pounds":50,"ac":11,"hp":10,"immunities":["poison","psychic"],"vulnerabilities":["bludgeoning"],"destroyed_when":"reduced-to-0-hit-points","activate_action":"magic","activate_requires":["command-word","hanging-on-vertical-surface"],"activate_range_feet":5,"deactivate_action":"magic","reflection_range_feet":30,"save_ability":"charisma","save_dc":15,"save_failure":"trapped-with-worn-and-carried-items-in-extradimensional-cell","save_advantage_when":"knows-the-mirrors-nature","save_auto_succeed":["constructs"],"cells":12,"cell_description":"infinite-expanse-of-thick-fog","cell_visibility_feet":10,"trapped_creatures_age":false,"trapped_creatures_need_food_drink_or_sleep":false,"escape_method":"magic-that-permits-planar-travel","otherwise_confined":"until-freed","full_cells_result":"frees-one-trapped-creature-at-random","freed_creature_appears":"unoccupied-space-within-sight-of-mirror-facing-away","shatter_result":"all-creatures-freed-in-unoccupied-spaces-near-mirror","name_action":"magic","name_range_feet":5,"name_result":"named-creature-or-cell-occupant-appears-as-image-and-can-communicate","free_action":"magic","free_command_word_ordinal":2,"free_result":"creature-appears-with-possessions-in-nearest-unoccupied-space-facing-away","extradimensional_items_destroyed":["bag-of-holding","portable-hole","similar-item"],"gate_to":"astral-plane","gate_radius_feet":10,"gate_affected":"creatures-not-behind-total-cover-sucked-to-random-location","gate_after":"closes","gate_one_way":true,"gate_can_reopen":false}}`,
		`{"mithral_armor":{"wearable":"under-normal-clothes","stealth_disadvantage":"not-imposed","strength_requirement":"does-not-apply"}}`,
		`{"mysterious_deck":{"source_pages":[231,232],"card_materials":["ivory","vellum"],"deck_sizes":[13,22],"thirteen_card_percent":75,"table_die_sides":100,"deck_13_cards":[{"card":"euryale","min":1,"max":8},{"card":"flames","min":9,"max":16},{"card":"jester","min":17,"max":24},{"card":"key","min":25,"max":32},{"card":"knight","min":33,"max":40},{"card":"moon","min":41,"max":48},{"card":"rogue","min":49,"max":56},{"card":"ruin","min":57,"max":64},{"card":"skull","min":65,"max":72},{"card":"star","min":73,"max":80},{"card":"sun","min":81,"max":88},{"card":"throne","min":89,"max":96},{"card":"void","min":97,"max":100}],"deck_22_cards":[{"card":"balance","min":1,"max":5},{"card":"comet","min":6,"max":10},{"card":"donjon","min":11,"max":14},{"card":"euryale","min":15,"max":18},{"card":"fates","min":19,"max":23},{"card":"flames","min":24,"max":27},{"card":"fool","min":28,"max":31},{"card":"gem","min":32,"max":36},{"card":"jester","min":37,"max":41},{"card":"key","min":42,"max":46},{"card":"knight","min":47,"max":51},{"card":"moon","min":52,"max":56},{"card":"puzzle","min":57,"max":60},{"card":"rogue","min":61,"max":64},{"card":"ruin","min":65,"max":68},{"card":"sage","min":69,"max":73},{"card":"skull","min":74,"max":77},{"card":"star","min":78,"max":82},{"card":"sun","min":83,"max":87},{"card":"talons","min":88,"max":91},{"card":"throne","min":92,"max":96},{"card":"void","min":97,"max":100}],"declare_draw_count_before_drawing":true,"excess_draws_have_no_effect":true,"max_hours_between_draws":1,"undrawn_cards_result":"remaining-cards-fly-from-deck-and-take-effect-at-once","drawn_card_disappears":true,"does_not_reappear":["fool","jester"],"reroll_when_gone":["fool","jester"],"cards":[{"card":"balance","effects":["increase-one-ability-score-while-decreasing-another","may-choose-not-to-adjust-scores-for-no-effect"],"facts":[{"key":"increase","amount":2,"unit":"points"},{"key":"maximum-score","amount":22,"unit":"score"},{"key":"decrease","amount":2,"unit":"points"},{"key":"cannot-decrease-at-or-below-score","amount":5,"unit":"score"}]},{"card":"comet","effects":["next-combat-against-hostile-creatures-select-one-as-foe-at-initiative","reducing-foe-to-0-hit-points-grants-advantage-on-death-saving-throws","no-effect-if-someone-else-reduces-foe-or-no-foe-chosen"],"facts":[{"key":"advantage-duration","amount":1,"unit":"year"}]},{"card":"donjon","effects":["entombed-in-suspended-animation-in-extradimensional-sphere","wearing-and-carrying-disappears-except-artifacts","artifacts-stay-in-space-you-occupied","imprisoned-until-found-and-removed-from-sphere","not-locatable-by-divination","wish-spell-reveals-prison-location","draw-no-more-cards"],"facts":[]},{"card":"euryale","effects":["medusa-like-visage-curses-you","penalty-to-saving-throws-while-cursed","only-a-god-or-fates-card-ends-curse"],"facts":[{"key":"saving-throw-penalty","amount":2,"unit":"minus"}]},{"card":"fates","effects":["avoid-or-erase-one-event","use-when-drawn-or-any-time-before-you-die"],"facts":[]},{"card":"flames","effects":["powerful-devil-becomes-your-enemy","devil-seeks-your-ruin-and-torments-you-before-attempting-to-slay-you","enmity-lasts-until-you-or-the-devil-dies"],"facts":[]},{"card":"fool","effects":["disadvantage-on-d20-tests","draw-another-card-not-counted-as-declared-draw"],"facts":[{"key":"duration","amount":72,"unit":"hours"},{"key":"additional-draws","amount":1,"unit":"cards"}]},{"card":"gem","effects":["jewelry-or-gems-appear-at-your-feet"],"facts":[{"key":"jewelry-pieces","amount":25,"unit":"pieces"},{"key":"jewelry-value-each","amount":2000,"unit":"gp"},{"key":"gems","amount":50,"unit":"pieces"},{"key":"gem-value-each","amount":1000,"unit":"gp"}]},{"card":"jester","effects":["advantage-on-d20-tests-or-draw-additional-cards-beyond-declared-draws"],"facts":[{"key":"duration","amount":72,"unit":"hours"},{"key":"additional-draws","amount":2,"unit":"cards"}]},{"card":"key","effects":["rare-or-rarer-magic-weapon-you-are-proficient-with-appears-on-your-person","gm-chooses-the-weapon"],"facts":[]},{"card":"knight","effects":["gain-service-of-a-knight-in-unoccupied-space-you-choose","same-alignment-as-you","serves-you-loyally-until-death","work-with-gm-on-name-and-backstory","gm-may-use-a-different-stat-block"],"facts":[{"key":"appears-within","amount":30,"unit":"feet"}]},{"card":"moon","effects":["gain-ability-to-cast-wish"],"facts":[{"key":"wish-casts","amount":1,"unit":"d3"}]},{"card":"puzzle","effects":["permanently-reduce-intelligence-or-wisdom","draw-one-additional-card-beyond-declared-draws"],"facts":[{"key":"reduction","amount":1,"unit":"d4"},{"key":"reduction-bonus","amount":1,"unit":"points"},{"key":"minimum-score","amount":1,"unit":"score"},{"key":"additional-draws","amount":1,"unit":"cards"}]},{"card":"rogue","effects":["npc-of-gm-choice-becomes-hostile-toward-you","identity-unknown-until-revealed","only-wish-or-divine-intervention-ends-hostility"],"facts":[]},{"card":"ruin","effects":["lose-all-wealth-other-than-magic-items","portable-property-vanishes","businesses-buildings-and-land-lost-in-way-that-alters-reality-least","documentation-proving-ownership-disappears"],"facts":[]},{"card":"sage","effects":["ask-one-question-in-meditation-and-mentally-receive-truthful-answer"],"facts":[{"key":"window","amount":1,"unit":"year"}]},{"card":"skull","effects":["avatar-of-death-appears-in-unoccupied-space-as-close-to-you-as-possible","avatar-targets-only-you","avatar-appears-as-ghostly-skeleton-in-tattered-black-robe-with-spectral-scythe","avatar-disappears-at-0-hit-points-or-when-you-die","ally-damaging-avatar-summons-another-avatar-that-targets-only-that-ally","each-creature-summons-only-one-avatar-from-this-draw","creature-slain-by-avatar-cannot-be-restored-to-life"],"facts":[]},{"card":"star","effects":["increase-one-ability-score"],"facts":[{"key":"increase","amount":2,"unit":"points"},{"key":"maximum-score","amount":24,"unit":"score"}]},{"card":"sun","effects":["magic-item-chosen-by-gm-appears-on-your-person","gain-temporary-hit-points-daily-at-dawn-until-you-die"],"facts":[{"key":"temporary-hit-points","amount":10,"unit":"hit-points"}]},{"card":"talons","effects":["every-magic-item-you-wear-or-carry-disintegrates","artifacts-in-your-possession-vanish-instead"],"facts":[]},{"card":"throne","effects":["gain-proficiency-and-expertise-in-history-insight-intimidation-or-persuasion","gain-rightful-ownership-of-a-small-keep","keep-is-home-to-monsters-that-must-be-cleared-before-claiming-it"],"facts":[]},{"card":"void","effects":["soul-drawn-from-body-and-contained-in-object-in-place-of-gm-choice","powerful-beings-guard-the-place","body-inert-ceases-aging-and-requires-no-food-air-or-water","wish-cannot-return-soul-but-reveals-object-location","draw-no-more-cards"],"facts":[]}],"avatar_of_death":{"size":"medium","creature_type":"undead","alignment":"neutral-evil","ac":20,"initiative_modifier":3,"initiative_score":13,"hit_points":"half-the-hp-maximum-of-its-summoner","speed_feet":60,"fly_speed_feet":60,"hover":true,"abilities":[{"ability":"strength","score":16,"mod":3,"save":3},{"ability":"dexterity","score":16,"mod":3,"save":3},{"ability":"constitution","score":16,"mod":3,"save":3},{"ability":"intelligence","score":16,"mod":3,"save":3},{"ability":"wisdom","score":16,"mod":3,"save":3},{"ability":"charisma","score":16,"mod":3,"save":3}],"damage_immunities":["necrotic","poison"],"condition_immunities":["charmed","exhaustion","frightened","paralyzed","petrified","poisoned","unconscious"],"truesight_feet":60,"passive_perception":13,"languages":"all-languages-known-to-its-summoner","challenge_rating":"none","xp":0,"proficiency_bonus":"equals-its-summoners","incorporeal_movement":"move-through-creatures-and-objects-as-difficult-terrain","incorporeal_damage_ending_turn_inside_object":{"average":5,"dice":1,"sides":10,"bonus":0,"type":"force"},"multiattack_count":"half-the-summoners-proficiency-bonus-rounded-up","multiattack_attack":"reaping-scythe","reaping_scythe_attack":"melee-attack-roll","reaping_scythe_automatic_hit":true,"reaping_scythe_reach_feet":5,"reaping_scythe_damage":{"average":7,"dice":1,"sides":8,"bonus":3,"type":"slashing"},"reaping_scythe_extra_damage":{"average":4,"dice":1,"sides":8,"bonus":0,"type":"necrotic"}}}}`,
		`{"necklace_of_adaptation":{"condition":"while-wearing","breathe_normally_in":"any-environment","advantage_on_saves_to":["avoid","end"],"advantage_against_condition":"poisoned"}}`,
		`{"necklace_of_fireballs":{"beads":{"dice":1,"sides":6,"bonus":3},"detach_action":"magic","throw_range_feet":60,"spell":"fireball","spell_level":3,"save_dc":15,"multiple_beads":"can-hurl-multiple-beads-or-whole-necklace-at-one-time","extra_damage_per_bead_after_first":{"dice":1,"sides":6,"bonus":0},"maximum_damage":{"dice":12,"sides":6,"bonus":0}}}`,
		`{"necklace_of_prayer_beads":{"attunement_by_any_of":["cleric","druid","paladin"],"magic_beads":{"dice":1,"sides":4,"bonus":2},"magic_bead_materials":["aquamarine","black-pearl","topaz"],"nonmagical_bead_examples":["amber","bloodstone","citrine","coral","jade","pearl","quartz"],"removed_bead_loses_magic":true,"bead_types":6,"type_selection":["gm-decides","roll-on-table"],"necklace_may_have_repeated_types":true,"requires_wearing":true,"cast_action":"bonus","save_dc":"your-spell-save-dc","uses_per_bead":1,"bead_recharge":"next-dawn","table_die_sides":20,"beads_table":[{"min":1,"max":6,"bead":"bead-of-blessing","spell":"bless","spell_level":0},{"min":7,"max":12,"bead":"bead-of-curing","spell":"cure-wounds","spell_level":2},{"min":13,"max":16,"bead":"bead-of-favor","spell":"greater-restoration","spell_level":0},{"min":17,"max":18,"bead":"bead-of-smiting","spell":"shining-smite","spell_level":0},{"min":19,"max":19,"bead":"bead-of-summons","spell":"guardian-of-faith","spell_level":0},{"min":20,"max":20,"bead":"bead-of-wind-walking","spell":"wind-walk","spell_level":0}]}}`,
		`{"nine_lives_stealer":{"attack_bonus":2,"damage_bonus":2,"property":"life-stealing","charges":{"dice":1,"sides":8,"bonus":1},"target_hit_points_below":100,"trigger_d20_roll":20,"save_ability":"constitution","save_dc":15,"save_failure":"slain-instantly","save_auto_succeed":["constructs","undead"],"charge_lost_when":"creature-is-slain","charges_lost":1,"when_no_charges":"weapon-loses-this-property"}}`,
		`{"oathbow":{"nocked_effect":"whispers-in-elvish-swift-defeat-to-my-enemies","command_trigger":"ranged-attack-with-this-weapon","command_methods":["utter","sign"],"sworn_enemy":"target-of-the-attack","sworn_enemy_ends":["until-it-dies","until-dawn-7-days-later"],"sworn_enemy_days":7,"max_sworn_enemies":1,"new_enemy_when":"after-the-next-dawn-once-sworn-enemy-dies","ranged_attack_advantage_against_sworn_enemy":true,"cover_ignored":["half-cover","three-quarters-cover"],"no_long_range_disadvantage":true,"extra_damage":{"dice":3,"sides":6,"type":"piercing"},"extra_damage_when":"attack-hits-sworn-enemy","other_weapons":"disadvantage-on-attack-rolls","other_weapons_while":"sworn-enemy-lives"}}`,
		`{"oil_of_etherealness":{"covers":"one-medium-or-smaller-creature-with-worn-and-carried-equipment","extra_vials_per_size_above_medium":1,"apply_minutes":10,"spell":"etherealness","duration_hours":1,"appearance":"cloudy-gray-oil-beads-form-outside-container-and-quickly-evaporate"}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"none:study:48-hours-within-6-days:study.48-hours-within-6-days ability.dexterity-increase-2 ability.dexterity-maximum-30 item.loses-magic-after-use item.regains-magic-after-a-century",
		"none:paint:expend-one-pot concentration remain-in-cube:item.pots-1d4 action.expend-one-pot area.cube-20-foot time.10-minutes concentration.required effect.paintings-become-real object.value-limit-25-gp object.total-value-limit-500-gp effect.energy-dissipates",
		"wear:cast-spell:expend-one-charge:item.attunement item.wear item.charges spell.cast-detect-thoughts spell.detect-thoughts-save-dc-13",
		"none:magic-action:command-word hanging-on-vertical-surface within-5-feet:action.magic item.command-word save.charisma-dc-15 effect.trapped-in-extradimensional-cell area.cells-12 effect.free-one-trapped-creature effect.destroys-bag-of-holding-and-portable-hole effect.gate-to-astral-plane",
		"wear:none::item.wear item.worn-under-clothes armor.no-stealth-disadvantage armor.no-strength-requirement",
		"none:draw-card:declare-draw-count draw-within-1-hour-of-previous:item.deck-13-or-22-cards table.mysterious-deck-1d100 action.declare-draws card.effects-22 card.fool-and-jester-leave-deck creature.avatar-of-death",
		"wear:none::item.attunement item.wear breath.any-environment save.advantage-poisoned",
		"none:magic-action:detach-bead throw-within-60-feet:item.beads-1d6-plus-3 action.magic spell.fireball-level-3 save.dexterity-dc-15 damage.extra-1d6-per-extra-bead damage.maximum-12d6",
		"wear:cast-spell:bonus-action once-per-bead-until-next-dawn:item.attunement-by-cleric-druid-or-paladin item.wear item.beads-1d4-plus-2 table.bead-type-1d20 spell.cast-from-bead item.bead-loses-magic-when-removed",
		"wield:weapon-attack:natural-20-against-creature-under-100-hit-points:item.attunement item.charges-1d8-plus-1 bonus.attack-and-damage-2 save.constitution-dc-15 effect.slain-instantly effect.constructs-and-undead-succeed effect.loses-property-without-charges",
		"wield:weapon-attack:command-words-on-ranged-attack:item.attunement effect.sworn-enemy attack.advantage-against-sworn-enemy cover.ignored-against-sworn-enemy damage.extra-3d6 attack.disadvantage-with-other-weapons",
		"none:apply:10-minutes one-vial-covers-medium-or-smaller:item.vials-per-size-above-medium time.apply-10-minutes spell.etherealness duration.1-hour",
	}
	index := 0
	for _, f := range loadItems235(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails235
				d := json.NewDecoder(strings.NewReader(expected[index]))
				d.DisallowUnknownFields()
				if err := d.Decode(&want); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(v.Details, &want) {
					t.Fatalf("PDF p.%d values mismatch: %+v", f.Presentations[0].Pages[0], v.Details)
				}
				got := fmt.Sprintf("%s:%s:%s:%s", v.Equip.Mode, v.Use.Mode, strings.Join(v.Use.Restrictions, " "), strings.Join(v.RequiredMechanicKeys, " "))
				if got != requirements[index] || len(v.Equip.Restrictions) != 0 {
					t.Fatalf("source requirements mismatch: %s", got)
				}
			})
			index++
		}
	}
	if index != len(expected) || index != len(requirements) {
		t.Fatal("fixture count disagrees with the batch")
	}
}

// The Mysterious Deck prints 22 cards in alphabetical order; the 13-card
// column lists thirteen of them. Each column covers the 1d100 exactly.
func TestItems235MysteriousDeckCompleteness(t *testing.T) {
	deck := loadItems235(t).Families[5].Variants[0].Details.Deck
	if len(deck.Cards) != 22 || len(deck.Deck22) != 22 || len(deck.Deck13) != 13 {
		t.Fatalf("deck sizes: %d cards, %d and %d rows", len(deck.Cards), len(deck.Deck22), len(deck.Deck13))
	}
	want13 := "euryale flames jester key knight moon rogue ruin skull star sun throne void"
	got13 := []string{}
	for _, r := range deck.Deck13 {
		got13 = append(got13, r.Card)
	}
	if strings.Join(got13, " ") != want13 || deck.Deck13[12].Min != 97 || deck.Deck13[12].Max != 100 || deck.Deck22[21].Max != 100 {
		t.Fatalf("13-card column: %v", got13)
	}
	for i, c := range deck.Cards {
		if c.Card != deck.Deck22[i].Card {
			t.Fatalf("card %d is %s, table row is %s", i, c.Card, deck.Deck22[i].Card)
		}
	}
	if len(deck.Avatar.Abilities) != 6 || deck.Avatar.ScytheDamage.Average != 7 || deck.Avatar.ScytheExtraDamage.Average != 4 {
		t.Fatal("Avatar of Death stat block incomplete")
	}
}

func TestItems235CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems235(t)
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	identities := [][]ItemSourceIdentity{contract.SourceIdentities()}
	loaders := []func() ([]ItemSourceIdentity, error){
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems225(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems226(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems227(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems228(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems229(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems230(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems231(); return b.SourceIdentities(), err },
	}
	for _, load := range loaders {
		ids, err := load()
		if err != nil {
			t.Fatal(err)
		}
		identities = append(identities, ids)
	}
	composed, err := ComposeItemSourceIdentities(append(identities, batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 97 || variants != 135 || composed[0].Name != "Potion of Healing" || composed[84].Name != "Gem of Brightness" || composed[85].Name != "Manual of Quickness of Action" || composed[96].Name != "Oil of Etherealness" {
		t.Fatal("expected Potion of Healing, batches #225 to #231 and exactly the twelve families of #235")
	}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), batch.SourceIdentities()); err == nil {
		t.Fatal("duplicate batch accepted")
	}
	clash := batch.SourceIdentities()[:1]
	clash[0].FamilyID, clash[0].Name, clash[0].PresentationIDs = "srd-5.2.1:content/items/other", "Other", []string{"srd-5.2.1:content/magic-items/other"}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), clash); err == nil {
		t.Fatal("duplicate variant identity accepted")
	}
	// The catalogue is descriptive: it must not promote coverage or add entries.
	register, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range register.Entries {
		for _, fixture := range entry.Fixtures {
			if strings.Contains(fixture, "items_235") {
				t.Fatal("catalogue batch must not be a coverage fixture")
			}
		}
	}
	if report := register.Coverage(); report.PlannedSpellItemScopes != 753 || report.FullSupport {
		t.Fatal("catalogue data changed the coverage boundary")
	}
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems235RejectMissingFields(t *testing.T) {
	var paths [][]any
	var walk func(node any, path []any)
	walk = func(node any, path []any) {
		switch n := node.(type) {
		case map[string]any:
			for key, child := range n {
				p := append(slices.Clone(path), key)
				paths = append(paths, p)
				walk(child, p)
			}
		case []any:
			for i, child := range n {
				walk(child, append(slices.Clone(path), i))
			}
		}
	}
	var root any
	if err := json.Unmarshal(items235JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 900 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items235JSON, &fresh); err != nil {
			t.Fatal(err)
		}
		node := fresh
		for _, step := range path[:len(path)-1] {
			switch s := step.(type) {
			case string:
				node = node.(map[string]any)[s]
			case int:
				node = node.([]any)[s]
			}
		}
		last, isKey := path[len(path)-1].(string)
		if !isKey {
			continue // list elements are covered by the list-length tests
		}
		delete(node.(map[string]any), last)
		data, err := json.Marshal(fresh)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := DecodeItems235(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems235RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails235]
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails235] { return &b.Families[family].Variants[0] }
	const (
		manual, pigments, medallion, mirror, mithral, deck, adaptation, fireballs, beads, nine, oathbow, oil = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-07" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 231 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/oil-of-sharpness"
			b.Families[12].Name = "Oil of Sharpness"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/manual-of-gainful-exercise" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Manual of Action" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[oil].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"widened base", func(b *batch) {
			b.Families[oathbow].Base.Options = append(b.Families[oathbow].Base.Options, "crossbow")
		}},
		{"dropped base", func(b *batch) { b.Families[oathbow].Base.Options = []string{"longbow"} }},
		{"armor exclusion dropped", func(b *batch) { b.Families[mithral].Base.Excluded = []string{} }},
		{"armor light", func(b *batch) {
			b.Families[mithral].Base.Options = append(b.Families[mithral].Base.Options, "any-light-armor")
		}},
		{"nine lives base", func(b *batch) { b.Families[nine].Base.Options = []string{"any-weapon"} }},
		{"invented base", func(b *batch) { b.Families[deck].Base.Options = []string{"any-armor"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{231, 231} }},
		{"page range", func(b *batch) { b.Families[deck].Presentations[0].Pages = []int{231, 232} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"split deck by size", func(b *batch) {
			b.Families[deck].Variants = append(b.Families[deck].Variants, b.Families[deck].Variants[0])
		}},
		{"variant id", func(b *batch) { v(b, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0).Name = "Manual" }},
		{"rarity", func(b *batch) { v(b, 0).Rarity = "rare" }},
		{"deck rarity", func(b *batch) { v(b, deck).Rarity = "artifact" }},
		{"unknown rarity", func(b *batch) { v(b, 0).Rarity = "mythic" }},
		{"dropped attunement", func(b *batch) { v(b, medallion).Attunement.Mode = "none" }},
		{"invented attunement", func(b *batch) { v(b, mirror).Attunement.Mode = "required" }},
		{"attunement restriction", func(b *batch) { v(b, beads).Attunement.Restrictions = []string{"cleric"} }},
		{"missing charges", func(b *batch) { v(b, medallion).Charges = nil }},
		{"dropped charges", func(b *batch) { *v(b, medallion).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"medallion maximum", func(b *batch) { v(b, medallion).Charges.Maximum = 6 }},
		{"medallion recharge", func(b *batch) { v(b, medallion).Charges.Recharge = "regain-all-at-dawn" }},
		{"nine lives maximum", func(b *batch) { v(b, nine).Charges.Maximum = 8 }},
		{"nine lives recharges", func(b *batch) { v(b, nine).Charges.Recharge = "regain-1d4-at-dawn" }},
		{"beads are not charges", func(b *batch) {
			*v(b, beads).Charges = ItemCharges{Mode: "limited", Maximum: 6, PerUse: 1, Recharge: "regain-all-at-dawn"}
		}},
		{"fireball beads are not charges", func(b *batch) {
			*v(b, fireballs).Charges = ItemCharges{Mode: "limited", Maximum: 9, PerUse: 1, Recharge: "none"}
		}},
		{"pots are not charges", func(b *batch) { v(b, pigments).Charges.Maximum = 4 }},
		{"equip mode", func(b *batch) { v(b, medallion).Equip.Mode = "none" }},
		{"equip restriction", func(b *batch) { v(b, nine).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, nine).Use.Mode = "magic-action" }},
		{"dropped use restriction", func(b *batch) { v(b, mirror).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0).Cost.Amount = 500 }},
		{"gem value is not a price", func(b *batch) {
			*v(b, deck).Cost = ItemCost{Kind: "listed", Amount: 200000, Currency: "cp", SourcePages: []int{231, 231}}
		}},
		{"cost currency", func(b *batch) { v(b, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys = v(b, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails235{} }},
		{"mismatched payload", func(b *batch) { v(b, 0).Details = v(b, 1).Details }},
		{"second payload", func(b *batch) { v(b, 0).Details.Pigments = v(b, 1).Details.Pigments }},
		{"blank key", func(b *batch) { v(b, 0).Details.Manual.Ability = " " }},
		{"prose value", func(b *batch) { v(b, 0).Details.Manual.Ability = "Dexterity score" }},
		{"negative value", func(b *batch) { v(b, 0).Details.Manual.Increase = -2 }},
		{"manual hours", func(b *batch) { v(b, manual).Details.Manual.StudyHours = 40 }},
		{"manual maximum", func(b *batch) { v(b, manual).Details.Manual.MaxScore = 20 }},
		{"manual century", func(b *batch) { v(b, manual).Details.Manual.RegainsAfterYears = 10 }},
		{"pigments cube", func(b *batch) { v(b, pigments).Details.Pigments.CubeFeet = 10 }},
		{"pigments value", func(b *batch) { v(b, pigments).Details.Pigments.ObjectMaxValueGP = 0 }},
		{"pigments total", func(b *batch) { v(b, pigments).Details.Pigments.TotalMaxValueGP = 25 }},
		{"pigments pot dice", func(b *batch) { v(b, pigments).Details.Pigments.Pots.Sides = 0 }},
		{"pigments empty triggers", func(b *batch) { v(b, pigments).Details.Pigments.FailureTriggers = []string{} }},
		{"medallion dc", func(b *batch) { v(b, medallion).Details.Medallion.SaveDC = 14 }},
		{"medallion charges", func(b *batch) { v(b, medallion).Details.Medallion.MaxCharges = 6 }},
		{"medallion recharge dice", func(b *batch) { v(b, medallion).Details.Medallion.Recharge.Dice = 0 }},
		{"mirror dc", func(b *batch) { v(b, mirror).Details.Mirror.SaveDC = 14 }},
		{"mirror cells", func(b *batch) { v(b, mirror).Details.Mirror.Cells = 10 }},
		{"mirror second command word", func(b *batch) { v(b, mirror).Details.Mirror.FreeCommandOrdinal = 1 }},
		{"mirror duplicate immunity", func(b *batch) { v(b, mirror).Details.Mirror.Immunities[1] = "poison" }},
		{"mirror empty vulnerabilities", func(b *batch) { v(b, mirror).Details.Mirror.Vulnerabilities = []string{} }},
		{"mirror gate radius", func(b *batch) { v(b, mirror).Details.Mirror.GateRadiusFeet = 0 }},
		{"mirror destroyed items", func(b *batch) { v(b, mirror).Details.Mirror.ExtradimensionalItems = []string{} }},
		{"mithral blank", func(b *batch) { v(b, mithral).Details.MithralArmor.Wearable = "" }},
		{"deck pages", func(b *batch) { v(b, deck).Details.Deck.SourcePages = []int{231} }},
		{"deck percent", func(b *batch) { v(b, deck).Details.Deck.ThirteenCardPercent = 50 }},
		{"deck sizes", func(b *batch) { v(b, deck).Details.Deck.DeckSizes = []int{13, 21} }},
		{"deck die", func(b *batch) { v(b, deck).Details.Deck.TableDieSides = 20 }},
		{"deck 13 gap", func(b *batch) { v(b, deck).Details.Deck.Deck13[3].Min = 26 }},
		{"deck 13 end", func(b *batch) { v(b, deck).Details.Deck.Deck13[12].Max = 99 }},
		{"deck 13 missing card", func(b *batch) { v(b, deck).Details.Deck.Deck13 = v(b, deck).Details.Deck.Deck13[:12] }},
		{"deck 13 unknown card", func(b *batch) { v(b, deck).Details.Deck.Deck13[0].Card = "hermit" }},
		{"deck 13 reordered", func(b *batch) {
			r := v(b, deck).Details.Deck.Deck13
			r[0].Card, r[1].Card = r[1].Card, r[0].Card
		}},
		{"deck 22 gap", func(b *batch) { v(b, deck).Details.Deck.Deck22[5].Min = 25 }},
		{"deck 22 missing card", func(b *batch) { v(b, deck).Details.Deck.Deck22 = v(b, deck).Details.Deck.Deck22[:21] }},
		{"deck 22 duplicate card", func(b *batch) { v(b, deck).Details.Deck.Deck22[1].Card = "balance" }},
		{"deck missing effects card", func(b *batch) { v(b, deck).Details.Deck.Cards = v(b, deck).Details.Deck.Cards[:21] }},
		{"deck card order", func(b *batch) {
			c := v(b, deck).Details.Deck.Cards
			c[0], c[1] = c[1], c[0]
		}},
		{"deck card no effects", func(b *batch) { v(b, deck).Details.Deck.Cards[0].Effects = []string{} }},
		{"deck card duplicate effect", func(b *batch) {
			v(b, deck).Details.Deck.Cards[0].Effects[1] = v(b, deck).Details.Deck.Cards[0].Effects[0]
		}},
		{"deck card duplicate fact", func(b *batch) {
			v(b, deck).Details.Deck.Cards[0].Facts[1].Key = v(b, deck).Details.Deck.Cards[0].Facts[0].Key
		}},
		{"deck hours between draws", func(b *batch) { v(b, deck).Details.Deck.MaxHoursBetween = 0 }},
		{"deck no vanishing cards", func(b *batch) { v(b, deck).Details.Deck.DoesNotReappear = []string{} }},
		{"avatar ac", func(b *batch) { v(b, deck).Details.Deck.Avatar.AC = 19 }},
		{"avatar initiative", func(b *batch) { v(b, deck).Details.Deck.Avatar.InitiativeScore = 3 }},
		{"avatar abilities", func(b *batch) {
			v(b, deck).Details.Deck.Avatar.Abilities = v(b, deck).Details.Deck.Avatar.Abilities[:5]
		}},
		{"avatar ability order", func(b *batch) {
			a := v(b, deck).Details.Deck.Avatar.Abilities
			a[0].Ability, a[1].Ability = a[1].Ability, a[0].Ability
		}},
		{"avatar immunities", func(b *batch) { v(b, deck).Details.Deck.Avatar.DamageImmunities = []string{} }},
		{"avatar scythe dice", func(b *batch) { v(b, deck).Details.Deck.Avatar.ScytheDamage.Sides = 0 }},
		{"avatar reach", func(b *batch) { v(b, deck).Details.Deck.Avatar.ScytheReachFeet = 0 }},
		{"adaptation saves", func(b *batch) { v(b, adaptation).Details.Adaptation.AdvantageSaves = []string{"avoid", "avoid"} }},
		{"fireballs level", func(b *batch) { v(b, fireballs).Details.Fireballs.SpellLevel = 4 }},
		{"fireballs dc", func(b *batch) { v(b, fireballs).Details.Fireballs.SaveDC = 14 }},
		{"fireballs range", func(b *batch) { v(b, fireballs).Details.Fireballs.ThrowRangeFeet = 30 }},
		{"fireballs maximum", func(b *batch) { v(b, fireballs).Details.Fireballs.MaximumDamage.Dice = 0 }},
		{"beads count", func(b *batch) { v(b, beads).Details.PrayerBeads.MagicBeads.Dice = 0 }},
		{"beads types", func(b *batch) { v(b, beads).Details.PrayerBeads.BeadTypes = 5 }},
		{"beads die", func(b *batch) { v(b, beads).Details.PrayerBeads.TableDieSides = 100 }},
		{"beads row gap", func(b *batch) { v(b, beads).Details.PrayerBeads.Beads[3].Min = 18 }},
		{"beads row end", func(b *batch) { v(b, beads).Details.PrayerBeads.Beads[5].Max = 19 }},
		{"beads missing row", func(b *batch) { v(b, beads).Details.PrayerBeads.Beads = v(b, beads).Details.PrayerBeads.Beads[:5] }},
		{"beads duplicate bead", func(b *batch) { v(b, beads).Details.PrayerBeads.Beads[1].Bead = "bead-of-blessing" }},
		{"beads no classes", func(b *batch) { v(b, beads).Details.PrayerBeads.AttunementBy = []string{} }},
		{"beads duplicate class", func(b *batch) { v(b, beads).Details.PrayerBeads.AttunementBy[1] = "cleric" }},
		{"beads no materials", func(b *batch) { v(b, beads).Details.PrayerBeads.MagicMaterials = []string{} }},
		{"nine lives bonus", func(b *batch) { v(b, nine).Details.NineLives.AttackBonus = 3 }},
		{"nine lives dc", func(b *batch) { v(b, nine).Details.NineLives.SaveDC = 14 }},
		{"nine lives hit points", func(b *batch) { v(b, nine).Details.NineLives.TargetHPBelow = 50 }},
		{"nine lives roll", func(b *batch) { v(b, nine).Details.NineLives.TriggerRoll = 19 }},
		{"nine lives dice", func(b *batch) { v(b, nine).Details.NineLives.Charges.Sides = 0 }},
		{"nine lives auto succeed", func(b *batch) { v(b, nine).Details.NineLives.AutoSucceed = []string{} }},
		{"oathbow days", func(b *batch) { v(b, oathbow).Details.Oathbow.SwornEnemyDays = 8 }},
		{"oathbow enemies", func(b *batch) { v(b, oathbow).Details.Oathbow.MaxSwornEnemies = 0 }},
		{"oathbow damage", func(b *batch) { v(b, oathbow).Details.Oathbow.ExtraDamage.Dice = 0 }},
		{"oathbow cover", func(b *batch) { v(b, oathbow).Details.Oathbow.CoverIgnored[1] = "half-cover" }},
		{"oathbow methods", func(b *batch) { v(b, oathbow).Details.Oathbow.CommandMethods = []string{} }},
		{"oil minutes", func(b *batch) { v(b, oil).Details.Etherealness.ApplyMinutes = 1 }},
		{"oil duration", func(b *batch) { v(b, oil).Details.Etherealness.DurationHours = 8 }},
		{"oil vials", func(b *batch) { v(b, oil).Details.Etherealness.ExtraVialsPerSizeAbove = 0 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems235(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems235(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems235RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/manual-of-quickness-of-action/variant/standard"`)
	if !bytes.Contains(items235JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items235JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems235(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"manual_of_quickness_of_action": {`, `"manual_of_quickness_plus": {`},
		"unknown detail":     {`"study_hours"`, `"study_hours_all"`},
		"string for number":  {`"save_dc": 13`, `"save_dc": "13"`},
		"fractional number":  {`"save_dc": 13`, `"save_dc": 13.5`},
		"number for boolean": {`"loses_magic_after_use": true`, `"loses_magic_after_use": 1`},
		"null fact list":     {`"facts": []`, `"facts": null`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items235JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items235JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems235(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items235JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems235(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor another batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items235JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems231(items235JSON); err == nil {
		t.Fatal("#231 decoder accepted the #235 batch")
	}
	if _, err := DecodeItems235(items231JSON); err == nil {
		t.Fatal("#235 decoder accepted the #231 batch")
	}
}
