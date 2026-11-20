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

func loadItems231(t *testing.T) ItemSourceBatch[ItemDetails231] {
	t.Helper()
	b, err := LoadItems231()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 220-223, in planning inventory order for issue #231. Feather
// Token and Figurine of Wondrous Power print "Rarity Varies": each printed kind
// is a variant with the rarity printed after its name.
func TestItems231SourceIdentities(t *testing.T) {
	type variant struct{ slug, name, rarity string }
	expected := []struct {
		name, slug string
		page       int
		category   string
		base       []string
		attunement string
		variants   []variant
	}{
		{"Energy Bow", "energy-bow", 220, "weapon", []string{"longbow", "shortbow"}, "required", []variant{{"standard", "Energy Bow", "very-rare"}}},
		{"Eversmoking Bottle", "eversmoking-bottle", 220, "wondrous-item", nil, "none", []variant{{"standard", "Eversmoking Bottle", "uncommon"}}},
		{"Eyes of Charming", "eyes-of-charming", 221, "wondrous-item", nil, "required", []variant{{"standard", "Eyes of Charming", "uncommon"}}},
		{"Eyes of Minute Seeing", "eyes-of-minute-seeing", 221, "wondrous-item", nil, "none", []variant{{"standard", "Eyes of Minute Seeing", "uncommon"}}},
		{"Eyes of the Eagle", "eyes-of-the-eagle", 221, "wondrous-item", nil, "none", []variant{{"standard", "Eyes of the Eagle", "uncommon"}}},
		{"Feather Token", "feather-token", 221, "wondrous-item", nil, "none", []variant{
			{"anchor", "Anchor", "uncommon"}, {"bird", "Bird", "rare"}, {"fan", "Fan", "uncommon"}, {"swan-boat", "Swan Boat", "rare"}, {"tree", "Tree", "uncommon"}, {"whip", "Whip", "rare"}}},
		{"Figurine of Wondrous Power", "figurine-of-wondrous-power", 221, "wondrous-item", nil, "none", []variant{
			{"bronze-griffon", "Bronze Griffon", "rare"}, {"ebony-fly", "Ebony Fly", "rare"}, {"golden-lions", "Golden Lions", "rare"}, {"ivory-goats", "Ivory Goats", "rare"},
			{"marble-elephant", "Marble Elephant", "rare"}, {"obsidian-steed", "Obsidian Steed", "very-rare"}, {"onyx-dog", "Onyx Dog", "rare"},
			{"serpentine-owl", "Serpentine Owl", "rare"}, {"silver-raven", "Silver Raven", "uncommon"}}},
		{"Flame Tongue", "flame-tongue", 223, "weapon", []string{"any-melee-weapon"}, "required", []variant{{"standard", "Flame Tongue", "rare"}}},
		{"Folding Boat", "folding-boat", 223, "wondrous-item", nil, "none", []variant{{"standard", "Folding Boat", "rare"}}},
		{"Frost Brand", "frost-brand", 223, "weapon", []string{"glaive", "greatsword", "longsword", "rapier", "scimitar", "shortsword"}, "required", []variant{{"standard", "Frost Brand", "very-rare"}}},
		{"Gauntlets of Ogre Power", "gauntlets-of-ogre-power", 223, "wondrous-item", nil, "required", []variant{{"standard", "Gauntlets of Ogre Power", "uncommon"}}},
		{"Gem of Brightness", "gem-of-brightness", 223, "wondrous-item", nil, "none", []variant{{"standard", "Gem of Brightness", "uncommon"}}},
	}
	b := loadItems231(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-07" || b.CatalogueIssue != 231 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #231")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 220–223") {
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
				// Only Eyes of Charming (3, regains all at dawn) and the Gem of Brightness (50, none) list charges.
				wantCharges := ItemCharges{Mode: "none", Recharge: "none"}
				switch want.slug {
				case "eyes-of-charming":
					wantCharges = ItemCharges{Mode: "limited", Maximum: 3, PerUse: 1, Recharge: "regain-all-at-dawn"}
				case "gem-of-brightness":
					wantCharges = ItemCharges{Mode: "limited", Maximum: 50, PerUse: 1, Recharge: "none"}
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
	if total != 25 {
		t.Fatalf("expected 25 source-listed variants, got %d", total)
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
		if r.Issue != 231 {
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
		t.Fatalf("inventory assigns %d families to issue #231", index)
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "220, 221, 223" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.EnergyBow.ArrowDamageType = "changed"
	if fresh := loadItems231(t); fresh.Families[0].Variants[0].Details.EnergyBow.ArrowDamageType != "force" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 25 variants, transcribed from pinned PDF pp.
// 220-223. The Eversmoking Bottle continues from p. 220 to p. 221 and the
// Figurine of Wondrous Power from p. 221 to p. 222. A printed 00 is recorded as
// 100; the Giant Fly's minus modifiers are amounts with the unit "minus".
func TestItems231SourceValues(t *testing.T) {
	expected := []string{
		`{"energy_bow":{"attack_bonus":1,"damage_bonus":1,"has_string":false,"arrow_appears":"when-arm-pulled-back-in-firing-motion","arrow_damage_type":"force","replaces_damage_type":"piercing","arrow_ends":"disappears-after-hit-or-miss","bright_light_feet":20,"dim_light_additional_feet":20,"arrow_of_restraint":{"trigger":"ranged-attack-instead-of-damage","save_ability":"strength","save_dc":15,"condition":"restrained","duration_seconds":60,"escape_action":"action","escape_check":{"ability":"strength","skill":"athletics","dc":20},"escape_ends":"effect-on-itself"},"arrow_of_transport":{"action":"magic","range_feet":60,"targets":["willing-medium-or-smaller-creature","unworn-uncarried-object"],"object_max_cube_feet":5,"destination":"unoccupied-space-you-can-see","destination_range_feet":10},"energy_ladder":{"action":"magic","wall_range_feet":60,"ladder_length_feet":60,"duration_seconds":60,"ends_with":"disappears"}}}`,
		`{"eversmoking_bottle":{"source_pages":[220,221],"action":"magic","operations":["open","close"],"emanation_feet":60,"obscurement":"heavily-obscured","growth_feet_per_minute":10,"max_emanation_feet":120,"closed_state":"cloud-fixed-in-place","closed_disperses_after_seconds":600,"strong_wind_disperses_after_seconds":60,"strong_wind_example":"gust-of-wind"}}`,
		`{"eyes_of_charming":{"appearance":"crystal-lenses-over-eyes","condition":"while-wearing","max_charges":3,"spell":"charm-person","save_dc":13,"base_spell_level":1,"base_charges":1,"levels_per_additional_charge":1,"recharge":"all-expended-charges-daily-at-dawn"}}`,
		`{"eyes_of_minute_seeing":{"appearance":"crystal-lenses-over-eyes","condition":"while-wearing","range_feet":1,"darkvision_within_range":true,"check_ability":"intelligence","check_skill":"investigation","check_mode":"advantage","check_subject":"examine-something-within-range","check_range_feet":1}}`,
		`{"eyes_of_the_eagle":{"appearance":"crystal-lenses-over-eyes","condition":"while-wearing","check_ability":"wisdom","check_skill":"perception","check_mode":"advantage","check_requirement":"relies-on-sight","detail_condition":"clear-visibility","detail_size_feet":2,"detail_subjects":"extremely-distant-creatures-and-objects"}}`,
		`{"feather_token":{"token":"anchor","table_die_sides":100,"table_min":1,"table_max":20,"use_limit":"single-use","action":"magic","placement":"touch-token-to-boat-or-ship","effect":"vessel-cannot-be-moved-by-any-means","dismiss_action":"none","facts":[{"key":"duration","amount":24,"unit":"hours"}],"traits":["touching-vessel-again-ends-effect"],"ends_on":["duration-elapsed","token-touched-to-vessel-again"]}}`,
		`{"feather_token":{"token":"bird","table_die_sides":100,"table_min":21,"table_max":35,"use_limit":"single-use","action":"magic","placement":"toss-token-into-air","effect":"enormous-multicolored-bird-replaces-token","dismiss_action":"magic","facts":[{"key":"toss-height","amount":5,"unit":"feet"},{"key":"carry-weight-at-maximum-speed","amount":500,"unit":"pounds"},{"key":"carry-weight-at-half-speed","amount":1000,"unit":"pounds"},{"key":"maximum-flying-speed","amount":16,"unit":"miles-per-hour"},{"key":"maximum-distance-per-day","amount":144,"unit":"miles"},{"key":"flying-hours-per-rest","amount":3,"unit":"hours"},{"key":"rest-hours","amount":1,"unit":"hours"}],"traits":["statistics-of-roc","cannot-attack","obeys-simple-commands"],"ends_on":["maximum-distance-flown-for-a-day","drops-to-0-hit-points"]}}`,
		`{"feather_token":{"token":"fan","table_die_sides":100,"table_min":36,"table_max":50,"use_limit":"single-use","action":"magic","placement":"toss-token-into-air-while-on-boat-or-ship","effect":"giant-flapping-fan-replaces-token","dismiss_action":"magic","facts":[{"key":"toss-height","amount":10,"unit":"feet"},{"key":"ships-with-sails-filled","amount":1,"unit":"ship"},{"key":"speed-increase","amount":5,"unit":"miles-per-hour"},{"key":"duration","amount":8,"unit":"hours"}],"traits":["floats","creates-strong-wind"],"ends_on":["duration-elapsed"]}}`,
		`{"feather_token":{"token":"swan-boat","table_die_sides":100,"table_min":51,"table_max":65,"use_limit":"single-use","action":"magic","placement":"touch-token-to-body-of-water","effect":"swan-shaped-boat-replaces-token","dismiss_action":"magic","facts":[{"key":"minimum-water-diameter","amount":60,"unit":"feet"},{"key":"boat-length","amount":50,"unit":"feet"},{"key":"boat-width","amount":20,"unit":"feet"},{"key":"speed","amount":6,"unit":"miles-per-hour"},{"key":"maximum-turn","amount":90,"unit":"degrees"},{"key":"duration","amount":24,"unit":"hours"}],"traits":["self-propelled","moves-across-water","commanded-with-magic-action-while-aboard"],"ends_on":["duration-elapsed"]}}`,
		`{"feather_token":{"token":"tree","table_die_sides":100,"table_min":66,"table_max":90,"use_limit":"single-use","action":"magic","placement":"touch-token-to-unoccupied-ground-space","effect":"nonmagical-oak-tree-replaces-token","dismiss_action":"none","facts":[{"key":"tree-height","amount":60,"unit":"feet"},{"key":"trunk-diameter","amount":5,"unit":"feet"},{"key":"branch-spread-radius","amount":20,"unit":"feet"}],"traits":["must-be-outdoors","nonmagical-oak-tree"],"ends_on":[]}}`,
		`{"feather_token":{"token":"whip","table_die_sides":100,"table_min":91,"table_max":100,"use_limit":"single-use","action":"magic","placement":"throw-token-to-point-within-10-feet","effect":"floating-whip-replaces-token","dismiss_action":"magic","facts":[{"key":"placement-range","amount":10,"unit":"feet"},{"key":"attack-bonus","amount":9,"unit":"plus"},{"key":"damage-dice","amount":1,"unit":"d6"},{"key":"damage-bonus","amount":5,"unit":"force-damage"},{"key":"attack-range","amount":10,"unit":"feet"},{"key":"flight-distance","amount":20,"unit":"feet"},{"key":"duration","amount":1,"unit":"hours"}],"traits":["melee-spell-attack","attack-as-bonus-action","fly-and-repeat-attack-as-bonus-action"],"ends_on":["duration-elapsed","dismissed-with-magic-action","you-die","you-are-incapacitated"]}}`,
		`{"figurine":{"source_pages":[221,222],"action":"magic","throw_range_feet":60,"throw_target":"point-on-ground","no_space_result":"figurine-does-not-become-creature","friendly_to":["you","your-allies"],"understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-owner","no_command_behavior":"defends-itself-and-takes-no-other-actions","revert_triggers":["duration-elapsed","creature-form-at-0-hit-points","magic-action-while-touching-creature"],"revert_action":"magic","figurines_in_set":1,"forms":[{"key":"bronze-griffon","creature":"griffon","duration_hours":6,"cooldown_days":5,"facts":[],"traits":[]}]}}`,
		`{"figurine":{"source_pages":[221,222],"action":"magic","throw_range_feet":60,"throw_target":"point-on-ground","no_space_result":"figurine-does-not-become-creature","friendly_to":["you","your-allies"],"understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-owner","no_command_behavior":"defends-itself-and-takes-no-other-actions","revert_triggers":["duration-elapsed","creature-form-at-0-hit-points","magic-action-while-touching-creature"],"revert_action":"magic","figurines_in_set":1,"forms":[{"key":"ebony-fly","creature":"giant-fly","duration_hours":12,"cooldown_days":2,"facts":[{"key":"armor-class","amount":11,"unit":"ac"},{"key":"initiative-bonus","amount":1,"unit":"plus"},{"key":"initiative-score","amount":11,"unit":"score"},{"key":"hit-points","amount":19,"unit":"hp"},{"key":"hit-point-dice","amount":3,"unit":"d10"},{"key":"hit-point-bonus","amount":3,"unit":"points"},{"key":"speed","amount":30,"unit":"feet"},{"key":"fly-speed","amount":60,"unit":"feet"},{"key":"strength","amount":14,"unit":"score"},{"key":"strength-modifier","amount":2,"unit":"plus"},{"key":"strength-save","amount":2,"unit":"plus"},{"key":"dexterity","amount":13,"unit":"score"},{"key":"dexterity-modifier","amount":1,"unit":"plus"},{"key":"dexterity-save","amount":1,"unit":"plus"},{"key":"constitution","amount":13,"unit":"score"},{"key":"constitution-modifier","amount":1,"unit":"plus"},{"key":"constitution-save","amount":1,"unit":"plus"},{"key":"intelligence","amount":2,"unit":"score"},{"key":"intelligence-modifier","amount":4,"unit":"minus"},{"key":"intelligence-save","amount":4,"unit":"minus"},{"key":"wisdom","amount":10,"unit":"score"},{"key":"wisdom-modifier","amount":0,"unit":"plus"},{"key":"wisdom-save","amount":0,"unit":"plus"},{"key":"charisma","amount":3,"unit":"score"},{"key":"charisma-modifier","amount":4,"unit":"minus"},{"key":"charisma-save","amount":4,"unit":"minus"},{"key":"darkvision","amount":60,"unit":"feet"},{"key":"passive-perception","amount":10,"unit":"score"},{"key":"challenge-rating","amount":0,"unit":"cr"},{"key":"experience-points","amount":0,"unit":"xp"},{"key":"proficiency-bonus","amount":2,"unit":"plus"}],"traits":["can-be-ridden-as-mount","stat-block-accompanies-item","large-beast","unaligned","languages-none"]}]}}`,
		`{"figurine":{"source_pages":[221,222],"action":"magic","throw_range_feet":60,"throw_target":"point-on-ground","no_space_result":"figurine-does-not-become-creature","friendly_to":["you","your-allies"],"understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-owner","no_command_behavior":"defends-itself-and-takes-no-other-actions","revert_triggers":["duration-elapsed","creature-form-at-0-hit-points","magic-action-while-touching-creature"],"revert_action":"magic","figurines_in_set":2,"forms":[{"key":"golden-lions","creature":"lion","duration_hours":1,"cooldown_days":7,"facts":[{"key":"figurines-usable-simultaneously","amount":2,"unit":"figurines"}],"traits":["cooldown-per-lion"]}]}}`,
		`{"figurine":{"source_pages":[221,222],"action":"magic","throw_range_feet":60,"throw_target":"point-on-ground","no_space_result":"figurine-does-not-become-creature","friendly_to":["you","your-allies"],"understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-owner","no_command_behavior":"defends-itself-and-takes-no-other-actions","revert_triggers":["duration-elapsed","creature-form-at-0-hit-points","magic-action-while-touching-creature"],"revert_action":"magic","figurines_in_set":3,"forms":[{"key":"goat-of-terror","creature":"giant-goat","duration_hours":3,"cooldown_days":15,"facts":[{"key":"horn-weapons","amount":2,"unit":"horns"},{"key":"lance-bonus","amount":1,"unit":"plus"},{"key":"longsword-bonus","amount":2,"unit":"plus"},{"key":"fear-emanation","amount":30,"unit":"feet"},{"key":"save-dc","amount":15,"unit":"dc"},{"key":"frightened-duration","amount":60,"unit":"seconds"},{"key":"immunity-after-success","amount":24,"unit":"hours"}],"traits":["cannot-attack","horns-removable-as-weapons-with-magic-action","horn-weapons-disappear-when-goat-reverts","frightens-hostile-creatures-starting-turn-in-emanation-while-ridden","wisdom-save","frightened-ends-when-no-longer-riding-or-goat-reverts","repeats-save-at-end-of-each-turn"]},{"key":"goat-of-traveling","creature":"large-goat","duration_hours":0,"cooldown_days":7,"facts":[{"key":"charges","amount":24,"unit":"charges"},{"key":"charge-cost-per-hour-or-portion","amount":1,"unit":"charges"}],"traits":["statistics-of-riding-horse","usable-as-often-as-wished-while-charged","reverts-when-out-of-charges","regains-all-charges-after-cooldown"]},{"key":"goat-of-travail","creature":"giant-goat","duration_hours":3,"cooldown_days":30,"facts":[],"traits":[]}]}}`,
		`{"figurine":{"source_pages":[221,222],"action":"magic","throw_range_feet":60,"throw_target":"point-on-ground","no_space_result":"figurine-does-not-become-creature","friendly_to":["you","your-allies"],"understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-owner","no_command_behavior":"defends-itself-and-takes-no-other-actions","revert_triggers":["duration-elapsed","creature-form-at-0-hit-points","magic-action-while-touching-creature"],"revert_action":"magic","figurines_in_set":1,"forms":[{"key":"marble-elephant","creature":"elephant","duration_hours":24,"cooldown_days":7,"facts":[],"traits":[]}]}}`,
		`{"figurine":{"source_pages":[221,222],"action":"magic","throw_range_feet":60,"throw_target":"point-on-ground","no_space_result":"figurine-does-not-become-creature","friendly_to":["you","your-allies"],"understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-owner","no_command_behavior":"defends-itself-and-takes-no-other-actions","revert_triggers":["duration-elapsed","creature-form-at-0-hit-points","magic-action-while-touching-creature"],"revert_action":"magic","figurines_in_set":1,"forms":[{"key":"obsidian-steed","creature":"nightmare","duration_hours":24,"cooldown_days":5,"facts":[{"key":"ignore-orders-chance","amount":10,"unit":"percent"}],"traits":["fights-only-to-defend-itself","may-ignore-revert-command","mounted-while-ignoring-orders-transports-to-random-location-on-plane-of-hades","reverts-to-figurine-on-plane-of-hades"]}]}}`,
		`{"figurine":{"source_pages":[221,222],"action":"magic","throw_range_feet":60,"throw_target":"point-on-ground","no_space_result":"figurine-does-not-become-creature","friendly_to":["you","your-allies"],"understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-owner","no_command_behavior":"defends-itself-and-takes-no-other-actions","revert_triggers":["duration-elapsed","creature-form-at-0-hit-points","magic-action-while-touching-creature"],"revert_action":"magic","figurines_in_set":1,"forms":[{"key":"onyx-dog","creature":"mastiff","duration_hours":6,"cooldown_days":7,"facts":[{"key":"intelligence","amount":8,"unit":"score"},{"key":"blindsight","amount":60,"unit":"feet"}],"traits":["can-speak-common"]}]}}`,
		`{"figurine":{"source_pages":[221,222],"action":"magic","throw_range_feet":60,"throw_target":"point-on-ground","no_space_result":"figurine-does-not-become-creature","friendly_to":["you","your-allies"],"understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-owner","no_command_behavior":"defends-itself-and-takes-no-other-actions","revert_triggers":["duration-elapsed","creature-form-at-0-hit-points","magic-action-while-touching-creature"],"revert_action":"magic","figurines_in_set":1,"forms":[{"key":"serpentine-owl","creature":"giant-owl","duration_hours":8,"cooldown_days":2,"facts":[],"traits":["telepathy-with-owner-at-any-range-on-same-plane"]}]}}`,
		`{"figurine":{"source_pages":[221,222],"action":"magic","throw_range_feet":60,"throw_target":"point-on-ground","no_space_result":"figurine-does-not-become-creature","friendly_to":["you","your-allies"],"understands_owner_languages":true,"obeys_owner_commands":true,"initiative":"immediately-after-owner","no_command_behavior":"defends-itself-and-takes-no-other-actions","revert_triggers":["duration-elapsed","creature-form-at-0-hit-points","magic-action-while-touching-creature"],"revert_action":"magic","figurines_in_set":1,"forms":[{"key":"silver-raven","creature":"raven","duration_hours":12,"cooldown_days":2,"facts":[],"traits":["raven-form-grants-owner-animal-messenger-cast-on-it"]}]}}`,
		`{"flame_tongue":{"condition":"while-holding","action":"bonus","activation":"command-word","effect":"flames-engulf-damage-dealing-part","bright_light_feet":40,"dim_light_additional_feet":40,"extra_damage":{"dice":2,"sides":6,"type":"fire"},"extra_damage_condition":"on-hit-while-ablaze","ends_on":["bonus-action-command-again","drop-weapon","stow-weapon","sheathe-weapon"]}}`,
		`{"folding_boat":{"length_inches":12,"width_inches":6,"depth_inches":6,"weight_pounds":4,"floats":true,"stores_items_when_open":true,"command_action":"magic","command_words":[{"ordinal":1,"result":"rowboat","requires":[]},{"ordinal":2,"result":"keelboat","requires":[]},{"ordinal":3,"result":"folds-back-into-box","requires":["no-creatures-aboard"]}],"objects_that_fit":"stored-inside-box","objects_that_do_not_fit":"remain-outside-box","vessel_weight":"normal-weight-of-vessel-its-size","stored_items_when_vessel":"remain-in-boat","vessel_statistics_in":["equipment"],"destroyed_when":"vessel-reduced-to-0-hit-points"}}`,
		`{"frost_brand":{"extra_damage":{"dice":1,"sides":6,"type":"cold"},"extra_damage_trigger":"hit-with-attack-roll-using-weapon","resistance_condition":"while-holding","resistance":"fire","light_condition":"freezing-temperatures","bright_light_feet":10,"dim_light_additional_feet":10,"extinguish_trigger":"draw-weapon","extinguish_targets":"nonmagical-flames","extinguish_radius_feet":30,"extinguish_uses":1,"extinguish_cooldown_hours":1}}`,
		`{"gauntlets_of_ogre_power":{"condition":"while-wearing","strength_score":19,"no_effect_at_or_above_score":19}}`,
		`{"gem_of_brightness":{"condition":"while-holding","action":"magic","command_words":3,"max_charges":50,"first_command_word":{"charge_cost":0,"bright_light_feet":30,"dim_light_additional_feet":30,"ends_on":["bonus-action-repeat-command-word","another-function-used"]},"second_command_word":{"charge_cost":1,"target":"one-creature-you-can-see","range_feet":60,"save_ability":"constitution","save_dc":15,"condition":"blinded","duration_seconds":60,"repeat_save":"end-of-each-of-its-turns","success_ends":"effect-on-itself"},"third_command_word":{"charge_cost":5,"shape":"cone","size_feet":30,"targets":"each-creature-in-cone","save":"as-struck-by-second-command-word-beam"},"when_depleted":"becomes-nonmagical-jewel","depleted_value_gp":50}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"wield:weapon-attack::item.attunement bonus.attack-and-damage-1 weapon.no-string weapon.energy-arrow-force-damage light.bright-20-dim-20 property.arrow-of-restraint save.strength-dc-15 check.strength-athletics-dc-20 property.arrow-of-transport action.magic property.energy-ladder",
		"none:magic-action:open-or-close:action.magic area.emanation-60-foot obscurement.heavily-obscured area.grows-10-feet-per-minute-to-120 effect.fixed-cloud-disperses-after-10-minutes effect.strong-wind-disperses-after-1-minute",
		"wear:cast-spell:expend-one-or-more-charges:item.attunement item.wear item.charges spell.cast-charm-person spell.charm-person-save-dc-13 spell.level-increases-per-extra-charge",
		"wear:none::item.wear vision.darkvision-1-foot check.investigation-advantage-within-1-foot",
		"wear:none::item.wear check.perception-sight-advantage vision.make-out-details-2-feet-across",
		"none:magic-action:single-use:action.magic item.single-use item.token-disappears item.kind-by-gm-or-1d100-table",
		"none:magic-action:single-use:action.magic item.single-use item.token-disappears item.kind-by-gm-or-1d100-table",
		"none:magic-action:single-use:action.magic item.single-use item.token-disappears item.kind-by-gm-or-1d100-table",
		"none:magic-action:single-use:action.magic item.single-use item.token-disappears item.kind-by-gm-or-1d100-table",
		"none:magic-action:single-use:action.magic item.single-use item.token-disappears item.kind-by-gm-or-1d100-table",
		"none:magic-action:single-use:action.magic item.single-use item.token-disappears item.kind-by-gm-or-1d100-table",
		"none:magic-action:throw-to-point-within-60-feet use-cooldown:action.magic item.throw-60-feet creature.friendly creature.obeys-commands creature.initiative-after-owner creature.reverts-to-figurine item.use-cooldown",
		"none:magic-action:throw-to-point-within-60-feet use-cooldown:action.magic item.throw-60-feet creature.friendly creature.obeys-commands creature.initiative-after-owner creature.reverts-to-figurine item.use-cooldown",
		"none:magic-action:throw-to-point-within-60-feet use-cooldown:action.magic item.throw-60-feet creature.friendly creature.obeys-commands creature.initiative-after-owner creature.reverts-to-figurine item.use-cooldown",
		"none:magic-action:throw-to-point-within-60-feet use-cooldown:action.magic item.throw-60-feet creature.friendly creature.obeys-commands creature.initiative-after-owner creature.reverts-to-figurine item.use-cooldown",
		"none:magic-action:throw-to-point-within-60-feet use-cooldown:action.magic item.throw-60-feet creature.friendly creature.obeys-commands creature.initiative-after-owner creature.reverts-to-figurine item.use-cooldown",
		"none:magic-action:throw-to-point-within-60-feet use-cooldown:action.magic item.throw-60-feet creature.friendly creature.obeys-commands creature.initiative-after-owner creature.reverts-to-figurine item.use-cooldown",
		"none:magic-action:throw-to-point-within-60-feet use-cooldown:action.magic item.throw-60-feet creature.friendly creature.obeys-commands creature.initiative-after-owner creature.reverts-to-figurine item.use-cooldown",
		"none:magic-action:throw-to-point-within-60-feet use-cooldown:action.magic item.throw-60-feet creature.friendly creature.obeys-commands creature.initiative-after-owner creature.reverts-to-figurine item.use-cooldown",
		"none:magic-action:throw-to-point-within-60-feet use-cooldown:action.magic item.throw-60-feet creature.friendly creature.obeys-commands creature.initiative-after-owner creature.reverts-to-figurine item.use-cooldown",
		"wield:bonus-action:command-word:item.attunement action.bonus item.command-word damage.fire-2d6 light.bright-40-dim-40",
		"none:magic-action:command-word:action.magic item.command-words-three item.floats item.container vessel.rowboat vessel.keelboat item.destroyed-with-vessel",
		"wield:weapon-attack:extinguish-once-per-hour:item.attunement damage.cold-1d6 defense.resistance-fire light.bright-10-dim-10 effect.extinguish-nonmagical-flames-30-feet item.once-per-hour",
		"wear:none::item.attunement item.wear ability.strength-score-19",
		"hold:magic-action:command-word:action.magic item.charges item.command-words-three light.bright-30-dim-30 save.constitution-dc-15 condition.blinded area.cone-30-foot item.nonmagical-jewel-when-expended",
	}
	index := 0
	for _, f := range loadItems231(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails231
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

func TestItems231CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems231(t)
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	first, err := LoadItems225()
	if err != nil {
		t.Fatal(err)
	}
	second, err := LoadItems226()
	if err != nil {
		t.Fatal(err)
	}
	third, err := LoadItems227()
	if err != nil {
		t.Fatal(err)
	}
	fourth, err := LoadItems228()
	if err != nil {
		t.Fatal(err)
	}
	fifth, err := LoadItems229()
	if err != nil {
		t.Fatal(err)
	}
	sixth, err := LoadItems230()
	if err != nil {
		t.Fatal(err)
	}
	composed, err := ComposeItemSourceIdentities(contract.SourceIdentities(), first.SourceIdentities(), second.SourceIdentities(), third.SourceIdentities(), fourth.SourceIdentities(), fifth.SourceIdentities(), sixth.SourceIdentities(), batch.SourceIdentities())
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 85 || variants != 123 || composed[0].Name != "Potion of Healing" || composed[73].Name != "Energy Bow" || composed[84].Name != "Gem of Brightness" {
		t.Fatal("expected Potion of Healing, batches #225 to #230 and exactly the twelve families of #231")
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
			if strings.Contains(fixture, "items_231") {
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
func TestItems231RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items231JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 1000 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items231JSON, &fresh); err != nil {
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
		if _, err := DecodeItems231(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems231RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails231]
	// v addresses a variant by family and variant index; see TestItems231SourceIdentities.
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails231] {
		return &b.Families[family].Variants[variant]
	}
	const (
		bow, bottle, charming, minute, eagle, feather, figurine, tongue, boat, brand, gauntlets, gem = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-06" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 230 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/gem-of-seeing"
			b.Families[12].Name = "Gem of Seeing"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/giant-slayer" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Energy Longbow" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"widened base", func(b *batch) { b.Families[bow].Base.Options = append(b.Families[bow].Base.Options, "crossbow") }},
		{"dropped base", func(b *batch) { b.Families[bow].Base.Options = []string{"longbow"} }},
		{"flame tongue base", func(b *batch) { b.Families[tongue].Base.Options = []string{"any-weapon"} }},
		{"invented base", func(b *batch) { b.Families[bottle].Base.Options = []string{"any-armor"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{221, 221} }},
		{"page range", func(b *batch) { b.Families[bottle].Presentations[0].Pages = []int{220, 221} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"missing token variant", func(b *batch) { b.Families[feather].Variants = b.Families[feather].Variants[:5] }},
		{"merged token variants", func(b *batch) { b.Families[feather].Variants = b.Families[feather].Variants[:1] }},
		{"reordered token variants", func(b *batch) {
			s := b.Families[feather].Variants
			s[0], s[1] = s[1], s[0]
		}},
		{"missing figurine variant", func(b *batch) { b.Families[figurine].Variants = b.Families[figurine].Variants[:8] }},
		{"invented figurine variant", func(b *batch) {
			b.Families[figurine].Variants = append(b.Families[figurine].Variants, b.Families[figurine].Variants[0])
		}},
		{"split goats", func(b *batch) {
			f := &b.Families[figurine]
			f.Variants = append(f.Variants[:3], append([]ItemSourceVariant[ItemDetails231]{f.Variants[3], f.Variants[3]}, f.Variants[4:]...)...)
		}},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Bow" }},
		{"token name", func(b *batch) { v(b, feather, 3).Name = "Swan" }},
		{"rarity", func(b *batch) { v(b, 0, 0).Rarity = "rare" }},
		{"token rarity", func(b *batch) { v(b, feather, 0).Rarity = "rare" }},
		{"figurine rarity", func(b *batch) { v(b, figurine, 5).Rarity = "rare" }},
		{"raven rarity", func(b *batch) { v(b, figurine, 8).Rarity = "rare" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"dropped attunement", func(b *batch) { v(b, bow, 0).Attunement.Mode = "none" }},
		{"invented attunement", func(b *batch) { v(b, bottle, 0).Attunement.Mode = "required" }},
		{"token attunement", func(b *batch) { v(b, feather, 0).Attunement.Mode = "required" }},
		{"attunement restriction", func(b *batch) { v(b, charming, 0).Attunement.Restrictions = []string{"spellcaster"} }},
		{"missing charges", func(b *batch) { v(b, charming, 0).Charges = nil }},
		{"dropped charges", func(b *batch) { *v(b, charming, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"charge maximum", func(b *batch) { v(b, charming, 0).Charges.Maximum = 4 }},
		{"charge recharge", func(b *batch) { v(b, charming, 0).Charges.Recharge = "none" }},
		{"gem charge maximum", func(b *batch) { v(b, gem, 0).Charges.Maximum = 5 }},
		{"gem recharges", func(b *batch) { v(b, gem, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"goat charges are not variant charges", func(b *batch) {
			*v(b, figurine, 3).Charges = ItemCharges{Mode: "limited", Maximum: 24, PerUse: 1, Recharge: "none"}
		}},
		{"invented charges", func(b *batch) { v(b, bottle, 0).Charges.Maximum = 10 }},
		{"equip mode", func(b *batch) { v(b, tongue, 0).Equip.Mode = "none" }},
		{"equip restriction", func(b *batch) { v(b, tongue, 0).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, tongue, 0).Use.Mode = "magic-action" }},
		{"dropped use restriction", func(b *batch) { v(b, boat, 0).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0, 0).Cost.Amount = 500 }},
		{"jewel value is not a price", func(b *batch) {
			*v(b, gem, 0).Cost = ItemCost{Kind: "listed", Amount: 5000, Currency: "cp", SourcePages: []int{223, 223}}
		}},
		{"cost currency", func(b *batch) { v(b, 0, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = v(b, 0, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails231{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.Eversmoking = v(b, 1, 0).Details.Eversmoking }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.EnergyBow.ArrowDamageType = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.EnergyBow.ArrowDamageType = "Force damage" }},
		{"negative value", func(b *batch) { v(b, 0, 0).Details.EnergyBow.AttackBonus = -1 }},
		{"bow bonus", func(b *batch) { v(b, bow, 0).Details.EnergyBow.AttackBonus = 0 }},
		{"bow restraint dc", func(b *batch) { v(b, bow, 0).Details.EnergyBow.Restraint.SaveDC = 14 }},
		{"bow escape dc", func(b *batch) { v(b, bow, 0).Details.EnergyBow.Restraint.EscapeCheck.DC = 15 }},
		{"bow transport range", func(b *batch) { v(b, bow, 0).Details.EnergyBow.Transport.RangeFeet = 0 }},
		{"bow duplicate target", func(b *batch) {
			v(b, bow, 0).Details.EnergyBow.Transport.Targets[1] = "willing-medium-or-smaller-creature"
		}},
		{"bow empty targets", func(b *batch) { v(b, bow, 0).Details.EnergyBow.Transport.Targets = []string{} }},
		{"bow ladder length", func(b *batch) { v(b, bow, 0).Details.EnergyBow.Ladder.LadderLengthFeet = 0 }},
		{"bottle pages", func(b *batch) { v(b, bottle, 0).Details.Eversmoking.SourcePages = []int{220} }},
		{"bottle maximum", func(b *batch) { v(b, bottle, 0).Details.Eversmoking.MaxEmanationFeet = 60 }},
		{"bottle growth", func(b *batch) { v(b, bottle, 0).Details.Eversmoking.GrowthFeetPerMinute = 0 }},
		{"bottle duplicate operation", func(b *batch) { v(b, bottle, 0).Details.Eversmoking.Operations[1] = "open" }},
		{"charming dc", func(b *batch) { v(b, charming, 0).Details.EyesCharming.SaveDC = 14 }},
		{"charming charges", func(b *batch) { v(b, charming, 0).Details.EyesCharming.MaxCharges = 4 }},
		{"charming level step", func(b *batch) { v(b, charming, 0).Details.EyesCharming.LevelsPerExtraCharge = 0 }},
		{"minute range", func(b *batch) { v(b, minute, 0).Details.EyesMinute.RangeFeet = 0 }},
		{"eagle detail size", func(b *batch) { v(b, eagle, 0).Details.EyesEagle.DetailSizeFeet = 0 }},
		{"token row", func(b *batch) { v(b, feather, 1).Details.FeatherToken.TableMin = 20 }},
		{"token row gap", func(b *batch) { v(b, feather, 5).Details.FeatherToken.TableMax = 99 }},
		{"token die", func(b *batch) { v(b, feather, 0).Details.FeatherToken.TableDieSides = 20 }},
		{"token kind", func(b *batch) { v(b, feather, 0).Details.FeatherToken.Token = "bird" }},
		{"token no facts", func(b *batch) { v(b, feather, 0).Details.FeatherToken.Facts = []ItemFact231{} }},
		{"token duplicate fact", func(b *batch) {
			v(b, feather, 1).Details.FeatherToken.Facts[1].Key = v(b, feather, 1).Details.FeatherToken.Facts[0].Key
		}},
		{"token duplicate end", func(b *batch) { v(b, feather, 1).Details.FeatherToken.EndsOn[1] = "maximum-distance-flown-for-a-day" }},
		{"token no traits", func(b *batch) { v(b, feather, 1).Details.FeatherToken.Traits = []string{} }},
		{"figurine pages", func(b *batch) { v(b, figurine, 0).Details.Figurine.SourcePages = []int{221} }},
		{"figurine throw range", func(b *batch) { v(b, figurine, 0).Details.Figurine.ThrowRangeFeet = 0 }},
		{"figurine no forms", func(b *batch) { v(b, figurine, 0).Details.Figurine.Forms = []ItemFigurineForm{} }},
		{"figurine one goat", func(b *batch) {
			v(b, figurine, 3).Details.Figurine.Forms = v(b, figurine, 3).Details.Figurine.Forms[:1]
		}},
		{"figurine two forms", func(b *batch) {
			f := v(b, figurine, 0).Details.Figurine
			f.Forms = append(f.Forms, f.Forms[0])
		}},
		{"figurine cooldown", func(b *batch) { v(b, figurine, 0).Details.Figurine.Forms[0].CooldownDays = 0 }},
		{"figurine duration", func(b *batch) { v(b, figurine, 0).Details.Figurine.Forms[0].DurationHours = 0 }},
		{"figurine set size", func(b *batch) { v(b, figurine, 2).Details.Figurine.FigurinesInSet = 0 }},
		{"figurine duplicate fact", func(b *batch) { v(b, figurine, 1).Details.Figurine.Forms[0].Facts[1].Key = "armor-class" }},
		{"figurine duplicate trait", func(b *batch) { v(b, figurine, 1).Details.Figurine.Forms[0].Traits[1] = "can-be-ridden-as-mount" }},
		{"figurine duplicate revert", func(b *batch) { v(b, figurine, 0).Details.Figurine.RevertTriggers[1] = "duration-elapsed" }},
		{"figurine duplicate friend", func(b *batch) { v(b, figurine, 0).Details.Figurine.FriendlyTo[1] = "you" }},
		{"flame tongue dice", func(b *batch) { v(b, tongue, 0).Details.FlameTongue.ExtraDamage.Dice = 0 }},
		{"flame tongue light", func(b *batch) { v(b, tongue, 0).Details.FlameTongue.BrightLightFeet = 0 }},
		{"flame tongue ends", func(b *batch) { v(b, tongue, 0).Details.FlameTongue.EndsOn = []string{} }},
		{"flame tongue duplicate end", func(b *batch) { v(b, tongue, 0).Details.FlameTongue.EndsOn[2] = "drop-weapon" }},
		{"boat size", func(b *batch) { v(b, boat, 0).Details.FoldingBoat.LengthInches = 0 }},
		{"boat missing command", func(b *batch) {
			v(b, boat, 0).Details.FoldingBoat.Commands = v(b, boat, 0).Details.FoldingBoat.Commands[:2]
		}},
		{"boat command order", func(b *batch) {
			c := v(b, boat, 0).Details.FoldingBoat.Commands
			c[0], c[1] = c[1], c[0]
		}},
		{"boat duplicate requirement", func(b *batch) { v(b, boat, 0).Details.FoldingBoat.Commands[2].Requires = []string{"a", "a"} }},
		{"boat no vessel statistics", func(b *batch) { v(b, boat, 0).Details.FoldingBoat.VesselStatistics = []string{} }},
		{"brand dice", func(b *batch) { v(b, brand, 0).Details.FrostBrand.ExtraDamage.Sides = 0 }},
		{"brand radius", func(b *batch) { v(b, brand, 0).Details.FrostBrand.ExtinguishRadiusFeet = 0 }},
		{"brand cooldown", func(b *batch) { v(b, brand, 0).Details.FrostBrand.ExtinguishCooldownHrs = 0 }},
		{"gauntlets score", func(b *batch) { v(b, gauntlets, 0).Details.Gauntlets.StrengthScore = 20 }},
		{"gauntlets threshold", func(b *batch) { v(b, gauntlets, 0).Details.Gauntlets.NoEffectAtOrAbove = 0 }},
		{"gem commands", func(b *batch) { v(b, gem, 0).Details.GemBrightness.CommandWords = 2 }},
		{"gem charges", func(b *batch) { v(b, gem, 0).Details.GemBrightness.MaxCharges = 5 }},
		{"gem first cost", func(b *batch) { v(b, gem, 0).Details.GemBrightness.First.ChargeCost = 1 }},
		{"gem second cost", func(b *batch) { v(b, gem, 0).Details.GemBrightness.Second.ChargeCost = 2 }},
		{"gem third cost", func(b *batch) { v(b, gem, 0).Details.GemBrightness.Third.ChargeCost = 1 }},
		{"gem save", func(b *batch) { v(b, gem, 0).Details.GemBrightness.Second.SaveDC = 0 }},
		{"gem cone", func(b *batch) { v(b, gem, 0).Details.GemBrightness.Third.SizeFeet = 0 }},
		{"gem ends", func(b *batch) { v(b, gem, 0).Details.GemBrightness.First.EndsOn = []string{} }},
		{"gem jewel value", func(b *batch) { v(b, gem, 0).Details.GemBrightness.DepletedValueGP = 0 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems231(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems231(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems231RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/energy-bow/variant/standard"`)
	if !bytes.Contains(items231JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items231JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems231(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"energy_bow": {`, `"energy_bow_plus": {`},
		"unknown detail":     {`"arrow_ends"`, `"arrow_ends_all"`},
		"string for number":  {`"save_dc": 15`, `"save_dc": "15"`},
		"fractional number":  {`"save_dc": 15`, `"save_dc": 15.5`},
		"number for boolean": {`"has_string": false`, `"has_string": 0`},
		"null fact list":     {`"facts": []`, `"facts": null`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items231JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items231JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems231(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items231JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems231(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items231JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems228(items231JSON); err == nil {
		t.Fatal("#228 decoder accepted the #231 batch")
	}
	if _, err := DecodeItems231(items228JSON); err == nil {
		t.Fatal("#231 decoder accepted the #228 batch")
	}
}
