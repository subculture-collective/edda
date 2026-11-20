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

func loadItems233(t *testing.T) ItemSourceBatch[ItemDetails233] {
	t.Helper()
	b, err := LoadItems233()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 225-227 (Ioun Stone continues on p. 228), in planning inventory
// order for issue #233. The Horn of Valhalla and the Ioun Stone print "Rarity
// Varies" or several rarities: each printed metal or stone type is a variant
// with the rarity printed after its name.
func TestItems233SourceIdentities(t *testing.T) {
	type variant struct{ slug, name, rarity string }
	expected := []struct {
		name, slug string
		page       int
		category   string
		base       []string
		attunement string
		variants   []variant
	}{
		{"Helm of Brilliance", "helm-of-brilliance", 225, "wondrous-item", nil, "required", []variant{{"standard", "Helm of Brilliance", "very-rare"}}},
		{"Helm of Comprehending Languages", "helm-of-comprehending-languages", 226, "wondrous-item", nil, "none", []variant{{"standard", "Helm of Comprehending Languages", "uncommon"}}},
		{"Helm of Telepathy", "helm-of-telepathy", 226, "wondrous-item", nil, "required", []variant{{"standard", "Helm of Telepathy", "uncommon"}}},
		{"Helm of Teleportation", "helm-of-teleportation", 226, "wondrous-item", nil, "required", []variant{{"standard", "Helm of Teleportation", "rare"}}},
		{"Holy Avenger", "holy-avenger", 226, "weapon", []string{"any-simple-weapon", "any-martial-weapon"}, "required", []variant{{"standard", "Holy Avenger", "legendary"}}},
		{"Horn of Blasting", "horn-of-blasting", 226, "wondrous-item", nil, "none", []variant{{"standard", "Horn of Blasting", "rare"}}},
		{"Horn of Valhalla", "horn-of-valhalla", 226, "wondrous-item", nil, "none", []variant{
			{"silver", "Silver", "rare"}, {"brass", "Brass", "rare"}, {"bronze", "Bronze", "very-rare"}, {"iron", "Iron", "legendary"}}},
		{"Horseshoes of Speed", "horseshoes-of-speed", 227, "wondrous-item", nil, "none", []variant{{"standard", "Horseshoes of Speed", "rare"}}},
		{"Horseshoes of a Zephyr", "horseshoes-of-a-zephyr", 226, "wondrous-item", nil, "none", []variant{{"standard", "Horseshoes of a Zephyr", "very-rare"}}},
		{"Immovable Rod", "immovable-rod", 227, "rod", nil, "none", []variant{{"standard", "Immovable Rod", "uncommon"}}},
		{"Instant Fortress", "instant-fortress", 227, "wondrous-item", nil, "required", []variant{{"standard", "Instant Fortress", "rare"}}},
		{"Ioun Stone", "ioun-stone", 227, "wondrous-item", nil, "required", []variant{
			{"absorption", "Absorption", "very-rare"}, {"agility", "Agility", "very-rare"}, {"awareness", "Awareness", "rare"}, {"fortitude", "Fortitude", "very-rare"},
			{"greater-absorption", "Greater Absorption", "legendary"}, {"insight", "Insight", "very-rare"}, {"intellect", "Intellect", "very-rare"}, {"leadership", "Leadership", "very-rare"},
			{"mastery", "Mastery", "legendary"}, {"protection", "Protection", "rare"}, {"regeneration", "Regeneration", "legendary"}, {"reserve", "Reserve", "rare"},
			{"strength", "Strength", "very-rare"}, {"sustenance", "Sustenance", "rare"}}},
	}
	b := loadItems233(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-09" || b.CatalogueIssue != 233 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #233")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 225–227") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	total := 0
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			if f.ID != id || f.Name != want.name || f.Kind != "magic" || f.Category != want.category || len(f.Base.Excluded) != 0 || !slices.Equal(f.Base.Options, append([]string{}, want.base...)) {
				t.Fatalf("source identity mismatch: %+v", f)
			}
			if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/magic-items/"+want.slug || !slices.Equal(f.Presentations[0].Pages, []int{want.page, want.page}) {
				t.Fatal("source presentation mismatch")
			}
			if len(f.Variants) != len(want.variants) {
				t.Fatalf("the source lists %d forms of %s, got %d variants", len(want.variants), want.name, len(f.Variants))
			}
			for j, wv := range want.variants {
				v := f.Variants[j]
				if v.ID != id+"/variant/"+wv.slug || v.Name != wv.name || v.Rarity != wv.rarity || v.Attunement.Mode != want.attunement || len(v.Attunement.Restrictions) != 0 {
					t.Fatalf("variant mismatch: %+v", v)
				}
				if v.Cost.Kind != "unspecified" || v.Cost.Amount != 0 || v.Cost.Currency != "none" || !slices.Equal(v.Cost.SourcePages, []int{want.page, want.page}) {
					t.Fatal("no price is listed; cost must stay unspecified")
				}
				// Only the Helm of Teleportation lists a charge pool (3, regains 1d3 at dawn).
				wantCharges := ItemCharges{Mode: "none", Recharge: "none"}
				if want.slug == "helm-of-teleportation" {
					wantCharges = ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-1d3-at-dawn"}
				}
				if *v.Charges != wantCharges {
					t.Fatalf("charges mismatch: %+v", *v.Charges)
				}
				if v.MechanicsAvailable() {
					t.Fatal("catalogue record became available")
				}
				total++
			}
		})
	}
	if total != 28 {
		t.Fatalf("expected 28 source-listed variants, got %d", total)
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
		if r.Issue != 233 {
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
		t.Fatalf("inventory assigns %d families to issue #233", index)
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "225, 226, 227" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.HelmBrilliance.Gems[0].Gem = "changed"
	if fresh := loadItems233(t); fresh.Families[0].Variants[0].Details.HelmBrilliance.Gems[0].Gem != "diamond" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 28 variants, transcribed from pinned PDF pp.
// 225-228. The Helm of Brilliance continues from p. 225 to p. 226 and the Ioun
// Stone from p. 227 to p. 228. A printed 00 is recorded as 100.
func TestItems233SourceValues(t *testing.T) {
	expected := []string{
		`{"helm_of_brilliance":{"source_pages":[225,226],"gems":[{"gem":"diamond","dice":1,"sides":10},{"gem":"ruby","dice":2,"sides":10},{"gem":"fire-opal","dice":3,"sides":10},{"gem":"opal","dice":4,"sides":10}],"gems_crumble_when_pried":true,"loses_magic_when_all_gems_gone":true,"condition":"while-wearing","diamond_light":{"requires":"at-least-one-diamond","emanation_feet":30,"dim_light_when":"undead-within-emanation","damage_trigger":"undead-starts-turn-in-emanation","damage":{"dice":1,"sides":6,"type":"radiant"}},"fire_opal_flames":{"requires":"at-least-one-fire-opal","action":"magic","target":"one-weapon-you-are-holding","bright_light_feet":10,"dim_light_additional_feet":10,"harmless_to":["you","the-weapon"],"extra_damage":{"dice":1,"sides":6,"type":"fire"},"extra_damage_trigger":"hit-with-attack-using-blazing-weapon","ends_on":["bonus-action-extinguish","drop-weapon","stow-weapon"]},"ruby_resistance":{"requires":"at-least-one-ruby","resistance":"fire"},"spells":{"save_dc":18,"component":"one-helm-gem-of-specified-type","gem_destroyed_when_cast":true,"gem_disappears_from_helm":true,"spells":[{"spell":"daylight","gem":"opal"},{"spell":"fireball","gem":"fire-opal"},{"spell":"prismatic-spray","gem":"diamond"},{"spell":"wall-of-fire","gem":"ruby"}]},"taking_fire_damage":{"trigger":"wearing-and-fire-damage-from-failed-save-against-spell","roll_die_sides":20,"destroyed_on_roll":1,"effect":"beams-of-light-from-remaining-gems","helm_destroyed":true,"emanation_feet":60,"origin":"you","save_ability":"dexterity","save_dc":17,"damage_type":"radiant","damage":"number-of-gems-in-helm"}}}`,
		`{"helm_of_comprehending_languages":{"condition":"while-wearing","spell":"comprehend-languages"}}`,
		`{"helm_of_telepathy":{"condition":"while-wearing","telepathy_range_feet":30,"spells":["detect-thoughts","suggestion"],"save_dc_spell":"suggestion","save_dc":13,"spell_reuse":"that-spell-cannot-be-cast-again","spell_reuse_window":"until-next-dawn"}}`,
		`{"helm_of_teleportation":{"condition":"while-wearing","max_charges":3,"spell":"teleport","charge_cost":1,"regain_dice":{"dice":1,"sides":3,"type":"charges"},"regain_timing":"dawn"}}`,
		`{"holy_avenger":{"attunement_by_any_of":["paladin"],"attack_bonus":3,"damage_bonus":3,"extra_damage":{"dice":2,"sides":10,"type":"radiant"},"extra_damage_against":["fiend","undead"],"extra_damage_trigger":"hit-with-weapon","aura":{"condition":"while-holding-drawn-weapon","emanation_feet":10,"origin":"you","affected":["you","creatures-friendly-to-you-in-emanation"],"benefit":"advantage-on-saving-throws","against":["spells","other-magical-effects"],"larger_emanation_min_paladin_levels":17,"larger_emanation_feet":30}}}`,
		`{"horn_of_blasting":{"action":"magic","cone_feet":30,"audible_feet":600,"save_ability":"constitution","save_dc":15,"damage":{"dice":5,"sides":8,"type":"thunder"},"condition":"deafened","condition_duration_seconds":60,"success_result":"half-damage-only","object_targets":"glass-or-crystal-in-cone-not-worn-or-carried","object_damage":{"dice":10,"sides":8,"type":"thunder"},"explosion_chance_percent":20,"explosion_damage":{"dice":10,"sides":6,"type":"force"},"explosion_target":"user","explosion_destroys_horn":true}}`,
		`{"horn_of_valhalla":{"action":"magic","spawn_range_feet":60,"spirits_from_plane":"ysgard","spirit_stat_block":"berserker","spirit_appearance":"living-breathing-warriors","duration_hours":1,"return_triggers":["duration-elapsed","drops-to-0-hit-points"],"spirit_immunities":["charmed","frightened"],"cooldown_days":7,"types_known":4,"type":"silver","table_die_sides":100,"table_min":1,"table_max":40,"spirits":2,"requirement":"none","unmet_requirement_result":"spirits-attack-you","met_requirement_result":["friendly-to-you-and-your-allies","follow-your-commands"]}}`,
		`{"horn_of_valhalla":{"action":"magic","spawn_range_feet":60,"spirits_from_plane":"ysgard","spirit_stat_block":"berserker","spirit_appearance":"living-breathing-warriors","duration_hours":1,"return_triggers":["duration-elapsed","drops-to-0-hit-points"],"spirit_immunities":["charmed","frightened"],"cooldown_days":7,"types_known":4,"type":"brass","table_die_sides":100,"table_min":41,"table_max":75,"spirits":3,"requirement":"proficiency-with-all-simple-weapons","unmet_requirement_result":"spirits-attack-you","met_requirement_result":["friendly-to-you-and-your-allies","follow-your-commands"]}}`,
		`{"horn_of_valhalla":{"action":"magic","spawn_range_feet":60,"spirits_from_plane":"ysgard","spirit_stat_block":"berserker","spirit_appearance":"living-breathing-warriors","duration_hours":1,"return_triggers":["duration-elapsed","drops-to-0-hit-points"],"spirit_immunities":["charmed","frightened"],"cooldown_days":7,"types_known":4,"type":"bronze","table_die_sides":100,"table_min":76,"table_max":90,"spirits":4,"requirement":"training-with-all-medium-armor","unmet_requirement_result":"spirits-attack-you","met_requirement_result":["friendly-to-you-and-your-allies","follow-your-commands"]}}`,
		`{"horn_of_valhalla":{"action":"magic","spawn_range_feet":60,"spirits_from_plane":"ysgard","spirit_stat_block":"berserker","spirit_appearance":"living-breathing-warriors","duration_hours":1,"return_triggers":["duration-elapsed","drops-to-0-hit-points"],"spirit_immunities":["charmed","frightened"],"cooldown_days":7,"types_known":4,"type":"iron","table_die_sides":100,"table_min":91,"table_max":100,"spirits":5,"requirement":"proficiency-with-all-martial-weapons","unmet_requirement_result":"spirits-attack-you","met_requirement_result":["friendly-to-you-and-your-allies","follow-your-commands"]}}`,
		`{"horseshoes_of_speed":{"set_size":4,"affix_action":"magic","affix_target":"hoof-of-horse-or-similar-creature","remove_action":"magic","requires_all_affixed_to_same_creature":true,"speed_bonus_feet":30}}`,
		`{"horseshoes_of_a_zephyr":{"set_size":4,"affix_action":"magic","affix_target":"hoof-of-horse-or-similar-creature","remove_action":"magic","requires_all_affixed":true,"float_inches":4,"moves_normally":true,"cross_surfaces":"nonsolid-or-unstable","surface_examples":["water","lava"],"leaves_tracks":false,"ignores_difficult_terrain":true,"travel_hours_per_day":12,"no_exhaustion_from":"extended-travel"}}`,
		`{"immovable_rod":{"material":"iron","button_on":"one-end","button_action":"utilize","fixed_effect":"magically-fixed-in-place","holds_against_gravity":true,"released_by":["you","another-creature"],"max_weight_pounds":8000,"overload_result":"deactivates-and-falls","move_action":"utilize","move_check":{"ability":"strength","skill":"athletics","dc":30},"move_feet":10,"move_on":"success"}}`,
		`{"instant_fortress":{"statuette_inches":1,"statuette_material":"adamantine","action":"magic","activation":"command-word","placed_on":"ground","tower":{"shape":"square","side_feet":20,"height_feet":30,"features":["arrow-slits-on-all-sides","battlement-atop"],"floors":2,"floor_connection_your_choice":["ladder","staircase","ramp"],"roof_access":"trapdoor","doors":1,"door_location":"ground-level-facing-you","door_opens_by":"your-command","door_open_action":"bonus","door_immune_to":"knock-and-similar-magic","cannot_be_tipped":true,"durable_parts":["roof","door","walls"],"armor_class":20,"hit_points":100,"immune_to":["bludgeoning","piercing","slashing"],"immunity_exception":"damage-dealt-by-siege-equipment","resistant_to":"all-other-damage"},"revert_activation":"repeat-command-word","revert_condition":"tower-is-empty","displaced_creatures":"pushed-to-unoccupied-space-outside-next-to-tower","displaced_objects":"unworn-uncarried-objects-pushed-clear","damage_repair":{"shrinking_repairs_damage":false,"only_spell":"wish","counts_as_spell_level_at_most":8,"each_casting_restores":"all-hit-points"}}}`,
		`{"ioun_stone":{"source_pages":[227],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"pale-lavender","shape":"ellipsoid","benefit":"cancel-spell","action":"reaction","ability":"none","facts":[{"key":"max-spell-level","amount":4,"unit":"level"},{"key":"burns-out-after-levels","amount":20,"unit":"levels"}],"traits":["target-spell-cast-by-creature-you-can-see","canceled-spell-has-no-effect","resources-used-to-cast-are-wasted","burns-out-turns-dull-gray-loses-magic"]}}`,
		`{"ioun_stone":{"source_pages":[227],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"deep-red","shape":"sphere","benefit":"ability-score-increase","action":"none","ability":"dexterity","facts":[{"key":"increase","amount":2,"unit":"points"},{"key":"maximum-score","amount":20,"unit":"score"}],"traits":["while-orbiting-head"]}}`,
		`{"ioun_stone":{"source_pages":[227],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"dark-blue","shape":"rhomboid","benefit":"advantage-on-rolls","action":"none","ability":"none","facts":[],"traits":["advantage-on-initiative-rolls","advantage-on-wisdom-perception-checks","while-orbiting-head"]}}`,
		`{"ioun_stone":{"source_pages":[227],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"pink","shape":"rhomboid","benefit":"ability-score-increase","action":"none","ability":"constitution","facts":[{"key":"increase","amount":2,"unit":"points"},{"key":"maximum-score","amount":20,"unit":"score"}],"traits":["while-orbiting-head"]}}`,
		`{"ioun_stone":{"source_pages":[227],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"marbled-lavender-and-green","shape":"ellipsoid","benefit":"cancel-spell","action":"reaction","ability":"none","facts":[{"key":"max-spell-level","amount":8,"unit":"level"},{"key":"burns-out-after-levels","amount":20,"unit":"levels"}],"traits":["target-spell-cast-by-creature-you-can-see","canceled-spell-has-no-effect","resources-used-to-cast-are-wasted","burns-out-turns-dull-gray-loses-magic"]}}`,
		`{"ioun_stone":{"source_pages":[227],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"incandescent-blue","shape":"sphere","benefit":"ability-score-increase","action":"none","ability":"wisdom","facts":[{"key":"increase","amount":2,"unit":"points"},{"key":"maximum-score","amount":20,"unit":"score"}],"traits":["while-orbiting-head"]}}`,
		`{"ioun_stone":{"source_pages":[227],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"marbled-scarlet-and-blue","shape":"sphere","benefit":"ability-score-increase","action":"none","ability":"intelligence","facts":[{"key":"increase","amount":2,"unit":"points"},{"key":"maximum-score","amount":20,"unit":"score"}],"traits":["while-orbiting-head"]}}`,
		`{"ioun_stone":{"source_pages":[227],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"marbled-pink-and-green","shape":"sphere","benefit":"ability-score-increase","action":"none","ability":"charisma","facts":[{"key":"increase","amount":2,"unit":"points"},{"key":"maximum-score","amount":20,"unit":"score"}],"traits":["while-orbiting-head"]}}`,
		`{"ioun_stone":{"source_pages":[228],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"pale-green","shape":"prism","benefit":"proficiency-bonus-increase","action":"none","ability":"none","facts":[{"key":"increase","amount":1,"unit":"points"}],"traits":["while-orbiting-head"]}}`,
		`{"ioun_stone":{"source_pages":[228],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"dusty-rose","shape":"prism","benefit":"armor-class-bonus","action":"none","ability":"none","facts":[{"key":"bonus","amount":1,"unit":"points"}],"traits":["while-orbiting-head"]}}`,
		`{"ioun_stone":{"source_pages":[228],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"pearly-white","shape":"spindle","benefit":"regain-hit-points","action":"none","ability":"none","facts":[{"key":"regain","amount":15,"unit":"hit-points"},{"key":"minimum-hit-points","amount":1,"unit":"hit-points"}],"traits":["end-of-each-hour","while-orbiting-head"]}}`,
		`{"ioun_stone":{"source_pages":[228],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"vibrant-purple","shape":"prism","benefit":"store-spells","action":"none","ability":"none","facts":[{"key":"capacity","amount":4,"unit":"levels"},{"key":"initial-levels-dice","amount":1,"unit":"dice"},{"key":"initial-levels-sides","amount":4,"unit":"sides"},{"key":"cast-into-min-level","amount":1,"unit":"level"},{"key":"cast-into-max-level","amount":4,"unit":"level"}],"traits":["initial-spells-chosen-by-gm","any-creature-can-cast-into-stone-by-touching-it","stored-spell-has-no-other-effect","spell-expended-without-effect-if-stone-cannot-hold-it","slot-level-determines-space-used","you-can-cast-stored-spell-while-orbiting","uses-original-caster-slot-level-save-dc-attack-bonus-and-spellcasting-ability","cast-spell-is-no-longer-stored"]}}`,
		`{"ioun_stone":{"source_pages":[228],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"pale-blue","shape":"rhomboid","benefit":"ability-score-increase","action":"none","ability":"strength","facts":[{"key":"increase","amount":2,"unit":"points"},{"key":"maximum-score","amount":20,"unit":"score"}],"traits":["while-orbiting-head"]}}`,
		`{"ioun_stone":{"source_pages":[228],"shared_rules":{"size":"roughly-marble","activation_action":"magic","method":"toss-into-the-air","orbit_distance_dice":1,"orbit_distance_sides":3,"max_orbiting":3,"orbiting_stone_counts_as":"object-you-are-wearing","orbit_traits":["avoids-contact-with-creatures-and-objects","adjusts-orbit-to-avoid-collisions","thwarts-attacks-and-snatching"],"stow_action":"utilize","stow_count":"any-number-orbiting","attunement_ends_result":"stone-falls-as-if-dropped"},"color":"clear","shape":"spindle","benefit":"no-need-to-eat-or-drink","action":"none","ability":"none","facts":[],"traits":["while-orbiting-head"]}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"wear:magic-action:gem-required:item.attunement item.wear item.gems-1d10-diamonds-2d10-rubies-3d10-fire-opals-4d10-opals item.loses-magic-without-gems light.emanation-30-foot damage.radiant-1d6-undead action.magic damage.fire-1d6-blazing-weapon light.bright-10-dim-10 defense.resistance-fire spell.save-dc-18 spell.daylight-fireball-prismatic-spray-wall-of-fire item.gem-destroyed-on-cast save.dexterity-dc-17 damage.radiant-per-gem",
		"wear:cast-spell::item.wear spell.cast-comprehend-languages",
		"wear:cast-spell:each-spell-once-until-next-dawn:item.attunement item.wear sense.telepathy-30-feet spell.cast-detect-thoughts spell.cast-suggestion spell.suggestion-save-dc-13 item.once-until-next-dawn",
		"wear:cast-spell:expend-one-charge:item.attunement item.wear item.charges spell.cast-teleport item.regain-1d3-charges-at-dawn",
		"wield:weapon-attack::item.attunement attunement.paladin bonus.attack-and-damage-3 damage.radiant-2d10-fiend-undead aura.emanation-10-foot save.advantage-spells-and-magical-effects aura.emanation-30-foot-paladin-17",
		"none:magic-action:explosion-chance-20-percent:action.magic area.cone-30-foot sound.audible-600-feet save.constitution-dc-15 damage.thunder-5d8 condition.deafened damage.thunder-10d8-glass-crystal damage.force-10d6-explosion item.destroyed-by-explosion",
		"none:magic-action:use-cooldown:action.magic creature.spirits-berserker-stat-block creature.returns-to-ysgard immunity.charmed-and-frightened item.use-cooldown item.type-by-gm-or-1d100-table creature.attack-you-without-requirement creature.friendly-with-requirement",
		"none:magic-action:use-cooldown:action.magic creature.spirits-berserker-stat-block creature.returns-to-ysgard immunity.charmed-and-frightened item.use-cooldown item.type-by-gm-or-1d100-table creature.attack-you-without-requirement creature.friendly-with-requirement",
		"none:magic-action:use-cooldown:action.magic creature.spirits-berserker-stat-block creature.returns-to-ysgard immunity.charmed-and-frightened item.use-cooldown item.type-by-gm-or-1d100-table creature.attack-you-without-requirement creature.friendly-with-requirement",
		"none:magic-action:use-cooldown:action.magic creature.spirits-berserker-stat-block creature.returns-to-ysgard immunity.charmed-and-frightened item.use-cooldown item.type-by-gm-or-1d100-table creature.attack-you-without-requirement creature.friendly-with-requirement",
		"none:magic-action:touch-to-hoof:action.magic item.set-of-four item.affixes-to-hoof speed.bonus-30-feet-with-all-four",
		"none:magic-action:touch-to-hoof:action.magic item.set-of-four item.affixes-to-hoof movement.float-4-inches terrain.cross-nonsolid-surfaces movement.leaves-no-tracks terrain.ignores-difficult travel.12-hours-no-exhaustion",
		"none:utilize-action:press-button:action.utilize item.fixed-in-place item.holds-8000-pounds check.strength-athletics-dc-30 movement.rod-10-feet",
		"none:magic-action:command-word:item.attunement action.magic item.command-word object.tower-20-foot-30-high defense.armor-class-20 defense.hit-points-100 defense.immune-nonsiege-weapon-damage defense.resistance-other-damage spell.wish-repairs-tower",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
		"none:magic-action:toss-into-air up-to-three-orbiting:item.attunement action.magic action.utilize item.orbits-head item.max-three-orbiting item.counts-as-worn item.type-determines-rarity-and-effects item.attunement-end-drops-stone",
	}
	index := 0
	for _, f := range loadItems233(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails233
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

func TestItems233CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems233(t)
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	identities := [][]ItemSourceIdentity{contract.SourceIdentities()}
	for _, load := range []func() ([]ItemSourceIdentity, error){
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems225(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems226(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems227(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems228(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems229(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems230(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems231(); return b.SourceIdentities(), err },
	} {
		ids, err := load()
		if err != nil {
			t.Fatal(err)
		}
		identities = append(identities, ids)
	}
	composed, err := ComposeItemSourceIdentities(append(slices.Clone(identities), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 97 || variants != 151 || composed[0].Name != "Potion of Healing" || composed[85].Name != "Helm of Brilliance" || composed[96].Name != "Ioun Stone" {
		t.Fatal("expected Potion of Healing, batches #225 to #231 and exactly the twelve families of #233")
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
			if strings.Contains(fixture, "items_233") {
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
func TestItems233RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items233JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 1000 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items233JSON, &fresh); err != nil {
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
		delete(node.(map[string]any), path[len(path)-1].(string))
		data, err := json.Marshal(fresh)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := DecodeItems233(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems233RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails233]
	// v addresses a variant by family and variant index; see TestItems233SourceIdentities.
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails233] {
		return &b.Families[family].Variants[variant]
	}
	const (
		helm, languages, telepathy, teleport, avenger, blasting, valhalla, speed, zephyr, rod, fortress, ioun = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-08" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 232 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/iron-bands"
			b.Families[12].Name = "Iron Bands"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/helm-of-seeing" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Helm of Radiance" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"rod category", func(b *batch) { b.Families[rod].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"narrowed base", func(b *batch) { b.Families[avenger].Base.Options = []string{"any-martial-weapon"} }},
		{"widened base", func(b *batch) {
			b.Families[avenger].Base.Options = append(b.Families[avenger].Base.Options, "any-ammunition")
		}},
		{"invented base", func(b *batch) { b.Families[helm].Base.Options = []string{"any-armor"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{226, 226} }},
		{"page range", func(b *batch) { b.Families[helm].Presentations[0].Pages = []int{225, 226} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"missing horn", func(b *batch) { b.Families[valhalla].Variants = b.Families[valhalla].Variants[:3] }},
		{"merged horns", func(b *batch) { b.Families[valhalla].Variants = b.Families[valhalla].Variants[:1] }},
		{"reordered horns", func(b *batch) {
			s := b.Families[valhalla].Variants
			s[0], s[1] = s[1], s[0]
		}},
		{"missing stone", func(b *batch) { b.Families[ioun].Variants = b.Families[ioun].Variants[:13] }},
		{"invented stone", func(b *batch) {
			b.Families[ioun].Variants = append(b.Families[ioun].Variants, b.Families[ioun].Variants[0])
		}},
		{"merged absorption stones", func(b *batch) {
			f := &b.Families[ioun]
			f.Variants = append(f.Variants[:4], f.Variants[5:]...)
		}},
		{"reordered stones", func(b *batch) {
			s := b.Families[ioun].Variants
			s[0], s[1] = s[1], s[0]
		}},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Helm" }},
		{"horn name", func(b *batch) { v(b, valhalla, 0).Name = "Silver Horn" }},
		{"rarity", func(b *batch) { v(b, 0, 0).Rarity = "rare" }},
		{"horn rarity", func(b *batch) { v(b, valhalla, 2).Rarity = "rare" }},
		{"iron rarity", func(b *batch) { v(b, valhalla, 3).Rarity = "very-rare" }},
		{"brass rarity", func(b *batch) { v(b, valhalla, 1).Rarity = "very-rare" }},
		{"stone rarity", func(b *batch) { v(b, ioun, 4).Rarity = "very-rare" }},
		{"absorption rarity", func(b *batch) { v(b, ioun, 0).Rarity = "rare" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"dropped attunement", func(b *batch) { v(b, helm, 0).Attunement.Mode = "none" }},
		{"invented attunement", func(b *batch) { v(b, languages, 0).Attunement.Mode = "required" }},
		{"horn attunement", func(b *batch) { v(b, valhalla, 0).Attunement.Mode = "required" }},
		{"class attunement is a payload fact", func(b *batch) { v(b, avenger, 0).Attunement.Restrictions = []string{"paladin"} }},
		{"missing charges", func(b *batch) { v(b, teleport, 0).Charges = nil }},
		{"dropped charges", func(b *batch) { *v(b, teleport, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"charge maximum", func(b *batch) { v(b, teleport, 0).Charges.Maximum = 4 }},
		{"charge recharge", func(b *batch) { v(b, teleport, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"invented charges", func(b *batch) { v(b, helm, 0).Charges.Maximum = 10 }},
		{"stone charges", func(b *batch) {
			*v(b, ioun, 11).Charges = ItemCharges{Mode: "limited", Maximum: 4, PerUse: 1, Recharge: "none"}
		}},
		{"equip mode", func(b *batch) { v(b, helm, 0).Equip.Mode = "none" }},
		{"equip restriction", func(b *batch) { v(b, helm, 0).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, avenger, 0).Use.Mode = "magic-action" }},
		{"dropped use restriction", func(b *batch) { v(b, telepathy, 0).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = v(b, 0, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails233{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.HelmLanguages = v(b, 1, 0).Details.HelmLanguages }},
		{"blank key", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.Condition = " " }},
		{"prose value", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.Condition = "While wearing the helm" }},
		{"negative value", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.Spells.SaveDC = -1 }},
		{"helm pages", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.SourcePages = []int{225} }},
		{"helm gem dice", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.Gems[1].Dice = 1 }},
		{"helm missing gem", func(b *batch) {
			v(b, helm, 0).Details.HelmBrilliance.Gems = v(b, helm, 0).Details.HelmBrilliance.Gems[:3]
		}},
		{"helm unknown gem", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.Gems[3].Gem = "emerald" }},
		{"helm spell dc", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.Spells.SaveDC = 17 }},
		{"helm mishap dc", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.FireDamage.SaveDC = 18 }},
		{"helm mishap die", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.FireDamage.RollDieSides = 0 }},
		{"helm emanation", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.DiamondLight.EmanationFee = 0 }},
		{"helm duplicate spell", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.Spells.SpellsWithGems[1].Spell = "daylight" }},
		{"helm missing spell", func(b *batch) {
			v(b, helm, 0).Details.HelmBrilliance.Spells.SpellsWithGems = v(b, helm, 0).Details.HelmBrilliance.Spells.SpellsWithGems[:3]
		}},
		{"helm duplicate end", func(b *batch) { v(b, helm, 0).Details.HelmBrilliance.FireOpalFlames.EndsOn[2] = "drop-weapon" }},
		{"languages spell", func(b *batch) { v(b, languages, 0).Details.HelmLanguages.Spell = "tongues" }},
		{"telepathy dc", func(b *batch) { v(b, telepathy, 0).Details.HelmTelepathy.SaveDC = 14 }},
		{"telepathy range", func(b *batch) { v(b, telepathy, 0).Details.HelmTelepathy.TelepathyFeet = 0 }},
		{"telepathy dc spell", func(b *batch) { v(b, telepathy, 0).Details.HelmTelepathy.SaveDCSpell = "fireball" }},
		{"telepathy duplicate spell", func(b *batch) { v(b, telepathy, 0).Details.HelmTelepathy.Spells[1] = "detect-thoughts" }},
		{"teleport charges", func(b *batch) { v(b, teleport, 0).Details.HelmTeleport.MaxCharges = 4 }},
		{"teleport regain", func(b *batch) { v(b, teleport, 0).Details.HelmTeleport.RegainDice.Sides = 0 }},
		{"avenger bonus", func(b *batch) { v(b, avenger, 0).Details.HolyAvenger.AttackBonus = 2 }},
		{"avenger damage bonus", func(b *batch) { v(b, avenger, 0).Details.HolyAvenger.DamageBonus = 0 }},
		{"avenger dice", func(b *batch) { v(b, avenger, 0).Details.HolyAvenger.ExtraDamage.Dice = 0 }},
		{"avenger no class", func(b *batch) { v(b, avenger, 0).Details.HolyAvenger.AttunementBy = []string{} }},
		{"avenger duplicate target", func(b *batch) { v(b, avenger, 0).Details.HolyAvenger.ExtraDamageAgainst[1] = "fiend" }},
		{"avenger larger aura", func(b *batch) { v(b, avenger, 0).Details.HolyAvenger.Aura.LargerFeet = 10 }},
		{"avenger paladin levels", func(b *batch) { v(b, avenger, 0).Details.HolyAvenger.Aura.LargerMinLevels = 0 }},
		{"blasting dc", func(b *batch) { v(b, blasting, 0).Details.HornBlasting.SaveDC = 14 }},
		{"blasting cone", func(b *batch) { v(b, blasting, 0).Details.HornBlasting.ConeFeet = 0 }},
		{"blasting explosion", func(b *batch) { v(b, blasting, 0).Details.HornBlasting.ExplosionPercent = 10 }},
		{"blasting object dice", func(b *batch) { v(b, blasting, 0).Details.HornBlasting.ObjectDamage.Dice = 0 }},
		{"valhalla row", func(b *batch) { v(b, valhalla, 1).Details.HornValhalla.TableMin = 40 }},
		{"valhalla row end", func(b *batch) { v(b, valhalla, 3).Details.HornValhalla.TableMax = 99 }},
		{"valhalla die", func(b *batch) { v(b, valhalla, 0).Details.HornValhalla.TableDieSides = 20 }},
		{"valhalla type", func(b *batch) { v(b, valhalla, 0).Details.HornValhalla.Type = "brass" }},
		{"valhalla spirits", func(b *batch) { v(b, valhalla, 2).Details.HornValhalla.Spirits = 3 }},
		{"valhalla types known", func(b *batch) { v(b, valhalla, 0).Details.HornValhalla.TypesKnown = 3 }},
		{"valhalla cooldown", func(b *batch) { v(b, valhalla, 0).Details.HornValhalla.CooldownDays = 0 }},
		{"valhalla duplicate immunity", func(b *batch) { v(b, valhalla, 0).Details.HornValhalla.Immunities[1] = "charmed" }},
		{"valhalla no returns", func(b *batch) { v(b, valhalla, 0).Details.HornValhalla.ReturnTriggers = []string{} }},
		{"speed bonus", func(b *batch) { v(b, speed, 0).Details.HorseshoesFast.SpeedBonusFeet = 20 }},
		{"speed set", func(b *batch) { v(b, speed, 0).Details.HorseshoesFast.SetSize = 0 }},
		{"zephyr float", func(b *batch) { v(b, zephyr, 0).Details.HorseshoesAir.FloatInches = 0 }},
		{"zephyr hours", func(b *batch) { v(b, zephyr, 0).Details.HorseshoesAir.TravelHours = 8 }},
		{"zephyr set", func(b *batch) { v(b, zephyr, 0).Details.HorseshoesAir.SetSize = 3 }},
		{"zephyr duplicate surface", func(b *batch) { v(b, zephyr, 0).Details.HorseshoesAir.SurfaceExamples[1] = "water" }},
		{"rod weight", func(b *batch) { v(b, rod, 0).Details.ImmovableRod.MaxWeightPounds = 4000 }},
		{"rod dc", func(b *batch) { v(b, rod, 0).Details.ImmovableRod.MoveCheck.DC = 20 }},
		{"rod move", func(b *batch) { v(b, rod, 0).Details.ImmovableRod.MoveFeet = 0 }},
		{"rod duplicate release", func(b *batch) { v(b, rod, 0).Details.ImmovableRod.ReleasedBy[1] = "you" }},
		{"fortress ac", func(b *batch) { v(b, fortress, 0).Details.Fortress.Tower.ArmorClass = 19 }},
		{"fortress hp", func(b *batch) { v(b, fortress, 0).Details.Fortress.Tower.HitPoints = 0 }},
		{"fortress side", func(b *batch) { v(b, fortress, 0).Details.Fortress.Tower.SideFeet = 30 }},
		{"fortress height", func(b *batch) { v(b, fortress, 0).Details.Fortress.Tower.HeightFeet = 20 }},
		{"fortress floors", func(b *batch) { v(b, fortress, 0).Details.Fortress.Tower.Floors = 0 }},
		{"fortress statuette", func(b *batch) { v(b, fortress, 0).Details.Fortress.StatuetteInches = 0 }},
		{"fortress wish level", func(b *batch) { v(b, fortress, 0).Details.Fortress.DamageRepairedBy.SpellLevelAtMost = 0 }},
		{"fortress duplicate connection", func(b *batch) { v(b, fortress, 0).Details.Fortress.Tower.FloorConnection[1] = "ladder" }},
		{"fortress duplicate immunity", func(b *batch) { v(b, fortress, 0).Details.Fortress.Tower.ImmuneTo[2] = "piercing" }},
		{"fortress no features", func(b *batch) { v(b, fortress, 0).Details.Fortress.Tower.Features = []string{} }},
		{"stone page", func(b *batch) { v(b, ioun, 0).Details.IounStone.SourcePages = []int{228} }},
		{"stone p228 page", func(b *batch) { v(b, ioun, 8).Details.IounStone.SourcePages = []int{227} }},
		{"stone max orbiting", func(b *batch) { v(b, ioun, 0).Details.IounStone.Shared.MaxOrbiting = 4 }},
		{"stone orbit dice", func(b *batch) { v(b, ioun, 0).Details.IounStone.Shared.OrbitDiceSides = 0 }},
		{"stone no orbit traits", func(b *batch) { v(b, ioun, 0).Details.IounStone.Shared.OrbitTraits = []string{} }},
		{"stone duplicate fact", func(b *batch) { v(b, ioun, 1).Details.IounStone.Facts[1].Key = "increase" }},
		{"stone no traits", func(b *batch) { v(b, ioun, 1).Details.IounStone.Traits = []string{} }},
		{"stone duplicate trait", func(b *batch) {
			t := v(b, ioun, 0).Details.IounStone.Traits
			t[1] = t[0]
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems233(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems233(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems233RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/helm-of-brilliance/variant/standard"`)
	if !bytes.Contains(items233JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items233JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems233(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"helm_of_brilliance": {`, `"helm_of_brilliance_plus": {`},
		"unknown detail":     {`"gems_crumble_when_pried"`, `"gems_crumble_when_removed"`},
		"string for number":  {`"save_dc": 18`, `"save_dc": "18"`},
		"fractional number":  {`"save_dc": 18`, `"save_dc": 18.5`},
		"number for boolean": {`"gems_crumble_when_pried": true`, `"gems_crumble_when_pried": 1`},
		"null fact list":     {`"facts": []`, `"facts": null`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items233JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items233JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems233(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items233JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems233(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items233JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems231(items233JSON); err == nil {
		t.Fatal("#231 decoder accepted the #233 batch")
	}
	if _, err := DecodeItems233(items231JSON); err == nil {
		t.Fatal("#233 decoder accepted the #231 batch")
	}
}
