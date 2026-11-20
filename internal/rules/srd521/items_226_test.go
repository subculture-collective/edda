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

func loadItems226(t *testing.T) ItemSourceBatch[ItemDetails226] {
	t.Helper()
	b, err := LoadItems226()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 211–213, in planning inventory order for issue #226. Bag of
// Tricks lists three cloth colours, each with its own table (p. 212); Belt of
// Giant Strength lists five rows in its Str./Rarity table (p. 213).
func TestItems226SourceIdentities(t *testing.T) {
	expected := []struct {
		name, slug string
		page       int
		category   string
		base       string
		attunement string
		variants   string // slug=rarity, in source order
	}{
		{"Arrow-Catching Shield", "arrow-catching-shield", 211, "armor", "shield", "required", "standard=rare"},
		{"Bag of Beans", "bag-of-beans", 211, "wondrous-item", "", "none", "standard=rare"},
		{"Bag of Devouring", "bag-of-devouring", 212, "wondrous-item", "", "none", "standard=very-rare"},
		{"Bag of Holding", "bag-of-holding", 212, "wondrous-item", "", "none", "standard=uncommon"},
		{"Bag of Tricks", "bag-of-tricks", 212, "wondrous-item", "", "none", "gray=uncommon rust=uncommon tan=uncommon"},
		{"Bead of Force", "bead-of-force", 212, "wondrous-item", "", "none", "standard=rare"},
		{"Bead of Nourishment", "bead-of-nourishment", 213, "wondrous-item", "", "none", "standard=common"},
		{"Belt of Dwarvenkind", "belt-of-dwarvenkind", 213, "wondrous-item", "", "required", "standard=rare"},
		{"Belt of Giant Strength", "belt-of-giant-strength", 213, "wondrous-item", "", "required", "hill=rare frost-or-stone=very-rare fire=very-rare cloud=legendary storm=legendary"},
		{"Berserker Axe", "berserker-axe", 213, "weapon", "battleaxe greataxe halberd", "required", "standard=rare"},
		{"Boots of Elvenkind", "boots-of-elvenkind", 213, "wondrous-item", "", "none", "standard=uncommon"},
		{"Boots of Levitation", "boots-of-levitation", 213, "wondrous-item", "", "required", "standard=rare"},
	}
	variantNames := []string{
		"Arrow-Catching Shield", "Bag of Beans", "Bag of Devouring", "Bag of Holding",
		"Gray Bag of Tricks", "Rust Bag of Tricks", "Tan Bag of Tricks",
		"Bead of Force", "Bead of Nourishment", "Belt of Dwarvenkind",
		"Belt of Giant Strength (hill)", "Belt of Giant Strength (frost or stone)", "Belt of Giant Strength (fire)",
		"Belt of Giant Strength (cloud)", "Belt of Giant Strength (storm)",
		"Berserker Axe", "Boots of Elvenkind", "Boots of Levitation",
	}
	b := loadItems226(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-02" || b.CatalogueIssue != 226 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #226")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 211–213") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	names := []string{}
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			if f.ID != id || f.Name != want.name || f.Kind != "magic" || f.Category != want.category || !slices.Equal(f.Base.Options, strings.Fields(want.base)) || len(f.Base.Excluded) != 0 {
				t.Fatalf("source identity mismatch: %+v", f)
			}
			if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/magic-items/"+want.slug || !slices.Equal(f.Presentations[0].Pages, []int{want.page, want.page}) {
				t.Fatal("source presentation mismatch")
			}
			got := []string{}
			for _, v := range f.Variants {
				slug := strings.TrimPrefix(v.ID, id+"/variant/")
				got = append(got, slug+"="+v.Rarity)
				names = append(names, v.Name)
				if v.Attunement.Mode != want.attunement || len(v.Attunement.Restrictions) != 0 || *v.Charges != (ItemCharges{Mode: "none", Recharge: "none"}) {
					t.Fatal("attunement or charges mismatch")
				}
				if v.Cost.Kind != "unspecified" || v.Cost.Amount != 0 || v.Cost.Currency != "none" || !slices.Equal(v.Cost.SourcePages, []int{want.page, want.page}) {
					t.Fatal("no price is listed; cost must stay unspecified")
				}
				if v.MechanicsAvailable() {
					t.Fatal("catalogue record became available")
				}
			}
			if !slices.Equal(got, strings.Fields(want.variants)) {
				t.Fatalf("source-listed variants mismatch: %v", got)
			}
		})
	}
	if !slices.Equal(names, variantNames) {
		t.Fatalf("expected the eighteen source-listed variant names, got %v", names)
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
			Issue   int    `json:"issue"`
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
		if r.Issue != 226 {
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
		t.Fatalf("inventory assigns %d families to issue #226", index)
	}
	// The planning batch names the same members, pages and issue.
	found := false
	for _, planned := range inventory.Batches {
		if planned.Key != b.Batch {
			continue
		}
		found = true
		if planned.Issue != 226 || planned.Pages != "211, 212, 213" || len(planned.Members) != 12 {
			t.Fatal("planning batch disagrees with issue #226")
		}
		for i, member := range planned.Members {
			if member.Name != expected[i].name || member.Page != expected[i].page {
				t.Fatalf("planning member %s disagrees with the batch", member.Name)
			}
		}
	}
	if !found {
		t.Fatal("planning batch data-magic-items-02 missing")
	}
	b.Families[0].Variants[0].Details.ArrowCatching.Condition = "changed"
	if fresh := loadItems226(t); fresh.Families[0].Variants[0].Details.ArrowCatching.Condition != "while-wielding" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all eighteen variants, transcribed from pinned PDF
// pp. 211–213. The Bag of Beans table spans both columns of p. 211 and Bead of
// Force continues from p. 212 onto p. 213. Table rows and d100 ranges are data
// here, never separate variants.
func TestItems226SourceValues(t *testing.T) {
	expected := []string{
		`{"arrow_catching":{"condition":"while-wielding","ac_bonus":2,"ac_bonus_against":"ranged-attack-rolls","adds_to_shield_bonus":true,"redirect_action":"reaction","redirect_trigger":"attacker-makes-ranged-attack-roll-against-target-within-range","redirect_range_feet":5,"redirect_effect":"become-target-instead"}}`,
		`{"beans":{"material":"heavy-cloth","beans_dice":3,"beans_die_sides":4,"weight_pounds_numerator":1,"weight_pounds_denominator":2,"weight_ignores_contents":true,"when_empty":"becomes-nonmagical","dump_trigger":"dump-one-or-more-beans","area_shape":"sphere","area_radius_feet":10,"area_center":"dumped-beans","dumped_beans_destroyed":true,"affects_user":true,"save":"dexterity","save_dc":15,"damage":{"dice":5,"sides":4,"type":"force"},"save_success":"half-damage","plant_steps":["remove-bean-from-bag","plant-in-dirt-or-sand","water-it"],"plant_delay_seconds":60,"bean_disappears":true,"effect_origin":"ground-where-planted","effect_selection":["gm-chooses","random-table"],"table_die_sides":100,"effects":[{"min":1,"max":1,"effect":"toadstools-sprout","creatures":[],"quantities":[{"label":"toadstools","dice":5,"sides":4,"bonus":0,"unit":"count"},{"label":"odd-roll-save-dc","dice":0,"sides":0,"bonus":15,"unit":"constitution-save-dc"},{"label":"odd-roll-failed-save-damage","dice":5,"sides":6,"bonus":0,"unit":"poison-damage"},{"label":"odd-roll-failed-save-poisoned","dice":0,"sides":0,"bonus":1,"unit":"hours"},{"label":"even-roll-temporary-hit-points","dice":5,"sides":6,"bonus":0,"unit":"temporary-hit-points"},{"label":"even-roll-temporary-hit-points-duration","dice":0,"sides":0,"bonus":1,"unit":"hours"}],"facts":["eater-rolls-any-die","odd-roll-save-or-poison-damage-and-poisoned","even-roll-gains-temporary-hit-points"]},{"min":2,"max":10,"effect":"geyser-erupts","creatures":[],"quantities":[{"label":"spout-height","dice":0,"sides":0,"bonus":30,"unit":"feet"},{"label":"duration","dice":1,"sides":4,"bonus":0,"unit":"minutes"}],"facts":["gm-chooses-liquid","liquid-water","liquid-beer","liquid-mayonnaise","liquid-tea","liquid-vinegar","liquid-wine","liquid-oil"]},{"min":11,"max":20,"effect":"treant-sprouts","creatures":["treant"],"quantities":[{"label":"treants","dice":0,"sides":0,"bonus":1,"unit":"count"}],"facts":["roll-any-die","odd-roll-chaotic-evil","even-roll-chaotic-good"]},{"min":21,"max":30,"effect":"threatening-statue-rises","creatures":[],"quantities":[{"label":"becomes-inanimate-after","dice":0,"sides":0,"bonus":24,"unit":"hours"}],"facts":["animate-but-immobile-stone-statue-in-user-likeness","makes-verbal-threats-against-user","if-user-leaves-directs-newcomers-to-find-and-attack-user","knows-user-location-on-same-plane"]},{"min":31,"max":40,"effect":"green-campfire-springs-forth","creatures":[],"quantities":[{"label":"burns-for","dice":0,"sides":0,"bonus":24,"unit":"hours"}],"facts":["green-flames","ends-early-if-extinguished"]},{"min":41,"max":50,"effect":"shrieker-fungi-sprout","creatures":["shrieker-fungus"],"quantities":[{"label":"shrieker-fungi","dice":0,"sides":0,"bonus":3,"unit":"count"}],"facts":[]},{"min":51,"max":60,"effect":"bright-pink-toads-crawl-forth","creatures":[],"quantities":[{"label":"toads","dice":1,"sides":4,"bonus":4,"unit":"count"},{"label":"monster-remains","dice":0,"sides":0,"bonus":1,"unit":"minutes"}],"facts":["touched-toad-transforms-into-large-or-smaller-monster-of-gm-choice","monster-acts-by-alignment-and-nature","monster-disappears-in-bright-pink-smoke"]},{"min":61,"max":70,"effect":"hungry-bulette-attacks","creatures":["bulette"],"quantities":[{"label":"bulettes","dice":0,"sides":0,"bonus":1,"unit":"count"}],"facts":["burrows-up-and-attacks"]},{"min":71,"max":80,"effect":"fruit-tree-grows","creatures":[],"quantities":[{"label":"fruit","dice":1,"sides":10,"bonus":20,"unit":"count"},{"label":"potion-fruit","dice":1,"sides":8,"bonus":0,"unit":"count"},{"label":"tree-vanishes-after","dice":0,"sides":0,"bonus":1,"unit":"hours"},{"label":"picked-fruit-retains-magic","dice":0,"sides":0,"bonus":30,"unit":"days"}],"facts":["potion-fruit-act-as-randomly-determined-potions","picked-fruit-remains"]},{"min":81,"max":90,"effect":"rainbow-egg-nest-springs-up","creatures":[],"quantities":[{"label":"eggs","dice":1,"sides":4,"bonus":3,"unit":"count"},{"label":"eater-save-dc","dice":0,"sides":0,"bonus":20,"unit":"constitution-save-dc"},{"label":"successful-save-lowest-ability-score-increase","dice":0,"sides":0,"bonus":1,"unit":"ability-score-points"},{"label":"failed-save-damage","dice":10,"sides":6,"bonus":0,"unit":"force-damage"}],"facts":["any-creature-that-eats-an-egg-saves","increase-is-permanent","randomly-choose-among-equally-low-scores","failed-save-internal-explosion"]},{"min":91,"max":95,"effect":"pyramid-bursts-upward","creatures":["mummy","mummy-lord"],"quantities":[{"label":"square-base-side","dice":0,"sides":0,"bonus":60,"unit":"feet"}],"facts":["burial-chamber-inside","contains-mummy-mummy-lord-or-other-undead-of-gm-choice","sarcophagus-contains-treasure-of-gm-choice"]},{"min":96,"max":100,"effect":"giant-beanstalk-sprouts","creatures":[],"quantities":[{"label":"beanstalks","dice":0,"sides":0,"bonus":1,"unit":"count"}],"facts":["height-of-gm-choice","top-leads-where-gm-chooses","example-great-view","example-cloud-giant-castle","example-another-plane-of-existence"]}]}}`,
		`{"devouring":{"resembles":"bag-of-holding","nature":"feeding-orifice-of-gigantic-extradimensional-creature","closes_orifice":"turning-bag-inside-out","creature_senses":"whatever-is-placed-inside","devours":["animal-matter-placed-wholly-inside","vegetable-matter-placed-wholly-inside"],"devoured_outcome":"devoured-and-lost-forever","pull_trigger":"part-of-living-creature-placed-inside","pull_chance_percent":50,"escape_action":"action","escape_check":{"ability":"strength","skill":"athletics","dc":15},"rescue_action":"action","rescue_check":{"ability":"strength","skill":"athletics","dc":20},"rescue_condition":"puller-not-pulled-inside-first","turn_start_inside":"creature-devoured-body-destroyed","stores":"inanimate-objects","capacity_cubic_feet":1,"swallow_frequency":"once-each-day","swallow_effect":"objects-spat-out-into-another-plane-of-existence","swallow_determined_by":"gm-determines-time-and-plane","destroyed_by":["pierced","torn"],"destroyed_effect":"contents-transported-to-random-astral-plane-location"}}`,
		`{"holding":{"interior_square_side_feet":2,"interior_depth_feet":4,"interior_approximate":true,"capacity_pounds":500,"capacity_cubic_feet":64,"weight_pounds":5,"weight_ignores_contents":true,"retrieve_action":"utilize","destroyed_by":["overloaded","pierced","torn"],"destroyed_effect":"contents-scattered-in-astral-plane","inside_out_effect":"contents-spill-unharmed","inside_out_reuse":"must-be-put-right-before-reuse","air_minutes":10,"air_divided_by":"breathing-creatures-inside","nesting_spaces":["handy-haversack","portable-hole","similar-item"],"nesting_effects":["both-items-destroyed","gate-to-astral-plane-opens"],"gate_origin":"where-one-item-was-placed-inside-the-other","gate_shape":"sphere","gate_radius_feet":10,"gate_effect":"creatures-sucked-to-random-astral-plane-location","gate_closes":true,"gate_one_way":true,"gate_reopenable":false}}`,
		`{"tricks":{"color":"gray","cloth_colors":["gray","rust","tan"],"appearance":"appears-empty","contents":"small-fuzzy-object","action":"magic","throw_feet":20,"transforms":"on-landing-into-creature-from-color-table","table_die_sides":8,"creatures":[{"roll":1,"creature":"weasel"},{"roll":2,"creature":"giant-rat"},{"roll":3,"creature":"badger"},{"roll":4,"creature":"boar"},{"roll":5,"creature":"panther"},{"roll":6,"creature":"giant-badger"},{"roll":7,"creature":"dire-wolf"},{"roll":8,"creature":"giant-elk"}],"stat_block_source":"monsters","vanishes":["next-dawn","reduced-to-0-hit-points"],"attitude":"friendly-to-user-and-allies","initiative":"immediately-after-user","command_action":"bonus-action","command_scope":"movement-and-action-on-next-turn","without_orders":"acts-appropriate-to-nature","uses":3,"uses_reset":"next-dawn"}}`,
		`{"tricks":{"color":"rust","cloth_colors":["gray","rust","tan"],"appearance":"appears-empty","contents":"small-fuzzy-object","action":"magic","throw_feet":20,"transforms":"on-landing-into-creature-from-color-table","table_die_sides":8,"creatures":[{"roll":1,"creature":"rat"},{"roll":2,"creature":"owl"},{"roll":3,"creature":"mastiff"},{"roll":4,"creature":"goat"},{"roll":5,"creature":"giant-goat"},{"roll":6,"creature":"giant-boar"},{"roll":7,"creature":"lion"},{"roll":8,"creature":"brown-bear"}],"stat_block_source":"monsters","vanishes":["next-dawn","reduced-to-0-hit-points"],"attitude":"friendly-to-user-and-allies","initiative":"immediately-after-user","command_action":"bonus-action","command_scope":"movement-and-action-on-next-turn","without_orders":"acts-appropriate-to-nature","uses":3,"uses_reset":"next-dawn"}}`,
		`{"tricks":{"color":"tan","cloth_colors":["gray","rust","tan"],"appearance":"appears-empty","contents":"small-fuzzy-object","action":"magic","throw_feet":20,"transforms":"on-landing-into-creature-from-color-table","table_die_sides":8,"creatures":[{"roll":1,"creature":"jackal"},{"roll":2,"creature":"ape"},{"roll":3,"creature":"baboon"},{"roll":4,"creature":"axe-beak"},{"roll":5,"creature":"black-bear"},{"roll":6,"creature":"giant-weasel"},{"roll":7,"creature":"giant-hyena"},{"roll":8,"creature":"tiger"}],"stat_block_source":"monsters","vanishes":["next-dawn","reduced-to-0-hit-points"],"attitude":"friendly-to-user-and-allies","initiative":"immediately-after-user","command_action":"bonus-action","command_scope":"movement-and-action-on-next-turn","without_orders":"acts-appropriate-to-nature","uses":3,"uses_reset":"next-dawn"}}`,
		`{"force":{"source_pages":[212,213],"appearance":"small-black-sphere","diameter_inches_numerator":3,"diameter_inches_denominator":4,"weight_ounces":1,"found_dice":1,"found_die_sides":4,"found_bonus":4,"action":"magic","throw_feet":60,"area_shape":"sphere","area_radius_feet":10,"destroyed_on_impact":true,"save":"dexterity","save_dc":15,"damage":{"dice":5,"sides":4,"type":"force"},"save_success":"no-damage","sphere_duration_seconds":60,"trapped_if_all":["failed-save","completely-within-area"],"pushed_if_any":["succeeded-on-save","partially-within-area"],"push_effect":"pushed-from-center-until-outside","passes_through":["breathable-air"],"blocked":["attacks","other-effects"],"move_action":"utilize","move_actor":"enclosed-creature","move_distance":"half-creature-speed","can_be_picked_up":true,"sphere_weight_pounds":1,"weight_ignores_creatures":true}}`,
		`{"nourishment":{"appearance":"flavorless-gelatinous-bead","consumed_by":"dissolves-on-tongue","nourishment_days":1,"equivalent":"rations"}}`,
		`{"dwarvenkind":{"condition":"while-wearing","benefits":["dwarvish","friend-of-dwarvenkind","toughness"],"language":"dwarvish","check_ability":"charisma","check_skill":"persuasion","check_audience":["dwarves","duergar"],"constitution_increase":2,"constitution_maximum":20,"beard_condition":"while-attuned","beard_chance_percent":50,"beard_timing":"each-day-at-dawn","beard_effects":["full-beard-if-able-to-grow-one","thicker-beard-if-already-bearded"],"additional_condition":"not-dwarf-or-duergar","additional_benefits":["darkvision","resilience"],"darkvision_feet":60,"resistance":"poison","save_advantage":"avoid-or-end-poisoned-condition"}}`,
		`{"giant_strength":{"condition":"while-wearing","giant_types":["hill"],"ability":"strength","score":21,"score_determined_by":"giant-type","no_effect_if_score_at_least":21}}`,
		`{"giant_strength":{"condition":"while-wearing","giant_types":["frost","stone"],"ability":"strength","score":23,"score_determined_by":"giant-type","no_effect_if_score_at_least":23}}`,
		`{"giant_strength":{"condition":"while-wearing","giant_types":["fire"],"ability":"strength","score":25,"score_determined_by":"giant-type","no_effect_if_score_at_least":25}}`,
		`{"giant_strength":{"condition":"while-wearing","giant_types":["cloud"],"ability":"strength","score":27,"score_determined_by":"giant-type","no_effect_if_score_at_least":27}}`,
		`{"giant_strength":{"condition":"while-wearing","giant_types":["storm"],"ability":"strength","score":29,"score_determined_by":"giant-type","no_effect_if_score_at_least":29}}`,
		`{"berserker":{"attack_bonus":1,"damage_bonus":1,"hit_point_maximum_per_level":1,"hit_point_condition":"while-attuned","cursed":true,"curse_trigger":"becoming-attuned","curse_effects":["unwilling-to-part-with-weapon","keeps-weapon-within-reach","disadvantage-on-attack-rolls-with-other-weapons"],"berserk_trigger":"another-creature-damages-user-while-weapon-in-possession","berserk_save":"wisdom","berserk_save_dc":15,"berserk_save_failure":"go-berserk","berserk_ends":"start-of-turn-with-no-creatures-in-range-seen-or-heard","berserk_end_range_feet":60,"berserk_target":"nearest-creature-seen-or-heard","berserk_tie_break":"choose-one-at-random","berserk_compulsion":["move-as-close-to-target-as-possible","take-attack-action-targeting-it"],"berserk_unreachable":"turn-ends-after-all-available-movement-used","berserk_retarget_triggers":["target-dies","target-no-longer-seen-or-heard"],"berserk_retarget":"next-nearest-creature-seen-or-heard"}}`,
		`{"elvenkind":{"condition":"while-wearing","silent_steps":true,"any_surface":true,"advantage_ability":"dexterity","advantage_skill":"stealth"}}`,
		`{"levitation":{"condition":"while-wearing","spell":"levitate","target":"self"}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"wield:reaction:ranged-attack-roll-against-target-within-5-feet:item.attunement item.wield-shield defense.armor-class-bonus-against-ranged-attacks reaction.become-target-of-ranged-attack",
		"none:dump-or-plant:dump-one-or-more-beans plant-in-dirt-or-sand-then-water:item.contains-3d4-beans item.nonmagical-when-empty area.sphere-10-foot-radius save.dexterity-dc-15-half damage.force-5d4 item.planted-bean-delayed-effect item.bean-effect-table",
		"none:container::item.extradimensional-creature-orifice item.devours-animal-and-vegetable-matter chance.pulled-inside-50-percent check.strength-athletics-dc-15 check.strength-athletics-dc-20 item.devours-creature-at-turn-start item.capacity-one-cubic-foot item.daily-planar-ejection item.destroyed-if-pierced-or-torn transport.random-astral-plane-location",
		"none:utilize-action:retrieve-item:item.extradimensional-storage item.capacity-500-pounds-64-cubic-feet item.fixed-weight-5-pounds action.utilize item.destroyed-if-overloaded-pierced-or-torn item.inside-out-spills-contents item.air-supply item.extradimensional-nesting-opens-astral-gate area.sphere-10-foot-radius transport.random-astral-plane-location",
		"none:magic-action:three-objects-until-next-dawn:action.magic item.throw-object-20-feet item.creature-table-by-color summon.creature-until-dawn-or-0-hit-points summon.friendly-acts-after-user action.bonus-action item.three-uses-until-next-dawn",
		"none:magic-action:three-objects-until-next-dawn:action.magic item.throw-object-20-feet item.creature-table-by-color summon.creature-until-dawn-or-0-hit-points summon.friendly-acts-after-user action.bonus-action item.three-uses-until-next-dawn",
		"none:magic-action:three-objects-until-next-dawn:action.magic item.throw-object-20-feet item.creature-table-by-color summon.creature-until-dawn-or-0-hit-points summon.friendly-acts-after-user action.bonus-action item.three-uses-until-next-dawn",
		"none:magic-action:throw-up-to-60-feet bead-destroyed:action.magic item.throw-60-feet item.destroyed-on-use area.sphere-10-foot-radius save.dexterity-dc-15 damage.force-5d4 effect.force-sphere-one-minute effect.trapped-or-pushed-out action.utilize",
		"none:consume:consume-one:item.consume sustenance.one-day-of-rations",
		"wear:none::item.attunement item.wear language.dwarvish check.advantage-persuasion-with-dwarves-and-duergar ability.constitution-plus-2-maximum-20 chance.daily-beard-50-percent species.not-dwarf-or-duergar senses.darkvision-60-feet defense.resistance-poison save.advantage-against-poisoned",
		"wear:none::item.attunement item.wear ability.strength-score-by-giant-type ability.no-effect-if-strength-not-lower",
		"wear:none::item.attunement item.wear ability.strength-score-by-giant-type ability.no-effect-if-strength-not-lower",
		"wear:none::item.attunement item.wear ability.strength-score-by-giant-type ability.no-effect-if-strength-not-lower",
		"wear:none::item.attunement item.wear ability.strength-score-by-giant-type ability.no-effect-if-strength-not-lower",
		"wear:none::item.attunement item.wear ability.strength-score-by-giant-type ability.no-effect-if-strength-not-lower",
		"wield:weapon-attack::item.attunement item.magic-weapon attack.weapon-bonus damage.weapon-bonus hit-points.maximum-plus-1-per-level item.curse item.curse-extends-on-attunement item.curse-unwilling-to-part attack.disadvantage-with-other-weapons save.wisdom-dc-15 condition.berserk-compelled-attack",
		"wear:none::item.wear movement.silent-steps check.advantage-dexterity-stealth",
		"wear:cast-spell:levitate self-only:item.attunement item.wear spell.cast-levitate-on-self",
	}
	index := 0
	for _, f := range loadItems226(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails226
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

// Composes the Healing contract fixture, merged batch #225 and this batch.
func TestItems226CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems226(t)
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	earlier, err := LoadItems225()
	if err != nil {
		t.Fatal(err)
	}
	composed, err := ComposeItemSourceIdentities(contract.SourceIdentities(), earlier.SourceIdentities(), batch.SourceIdentities())
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 25 || variants != 38 || composed[0].Name != "Potion of Healing" || composed[1].Name != "Adamantine Armor" ||
		composed[12].Name != "Armor, +1, +2, or +3" || composed[13].Name != "Arrow-Catching Shield" || composed[24].Name != "Boots of Levitation" {
		t.Fatal("expected Potion of Healing, the twelve #225 families and the exact twelve families and eighteen variants of #226")
	}
	// Batch #225 ends on p. 211 where this batch starts; the two share no identity.
	if len(earlier.Families) != 12 || earlier.Families[10].Name != "Armor of Vulnerability" || earlier.Families[10].Presentations[0].Pages[0] != 211 {
		t.Fatal("batch #225 changed")
	}
	if _, err := ComposeItemSourceIdentities(earlier.SourceIdentities(), batch.SourceIdentities(), earlier.SourceIdentities()); err == nil {
		t.Fatal("repeated earlier batch accepted")
	}
	// A batch decoder accepts only its own batch.
	if _, err := DecodeItems226(items225JSON); err == nil {
		t.Fatal("batch #226 decoder accepted batch #225")
	}
	if _, err := DecodeItems225(items226JSON); err == nil {
		t.Fatal("batch #225 decoder accepted batch #226")
	}
	if len(contract.Families) != 1 || len(contract.Families[0].Variants) != 4 {
		t.Fatal("the Healing contract fixture changed")
	}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), batch.SourceIdentities()); err == nil {
		t.Fatal("duplicate batch accepted")
	}
	clash := batch.SourceIdentities()[:1]
	clash[0].FamilyID, clash[0].Name, clash[0].PresentationIDs = "srd-5.2.1:content/items/other", "Other", []string{"srd-5.2.1:content/magic-items/other"}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), clash); err == nil {
		t.Fatal("duplicate variant identity accepted")
	}
	if _, err := ComposeItemSourceIdentities(); err == nil {
		t.Fatal("empty composition accepted")
	}
	// The catalogue is descriptive: it must not promote coverage or add entries.
	register, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range register.Entries {
		for _, fixture := range entry.Fixtures {
			if strings.Contains(fixture, "items_226") {
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
func TestItems226RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items226JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 900 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items226JSON, &fresh); err != nil {
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
		if _, err := DecodeItems226(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems226RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails226]
	// Family indexes follow inventory order; see TestItems226SourceIdentities.
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-01" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 225 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) { b.Families = append(b.Families, b.Families[0]) }},
		{"unknown family", func(b *batch) { b.Families[11].ID = "srd-5.2.1:content/items/boots-of-speed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Arrow Catching Shield" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"widened base", func(b *batch) { b.Families[9].Base.Options = append(b.Families[9].Base.Options, "handaxe") }},
		{"narrowed base", func(b *batch) { b.Families[9].Base.Options = b.Families[9].Base.Options[:2] }},
		{"invented base", func(b *batch) { b.Families[1].Base.Options = []string{"pouch"} }},
		{"invented exclusion", func(b *batch) { b.Families[0].Base.Excluded = []string{"shield"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{212, 212} }},
		{"page range", func(b *batch) { b.Families[5].Presentations[0].Pages = []int{212, 213} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"duplicate presentation", func(b *batch) {
			b.Families[0].Presentations = append(b.Families[0].Presentations, b.Families[0].Presentations[0])
		}},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"missing colour variant", func(b *batch) { b.Families[4].Variants = b.Families[4].Variants[:2] }},
		{"missing belt row", func(b *batch) { b.Families[8].Variants = b.Families[8].Variants[:4] }},
		{"reordered belt rows", func(b *batch) {
			b.Families[8].Variants[0], b.Families[8].Variants[1] = b.Families[8].Variants[1], b.Families[8].Variants[0]
		}},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"duplicate variant id", func(b *batch) { b.Families[4].Variants[1].ID = b.Families[4].Variants[0].ID }},
		{"missing variant id", func(b *batch) { b.Families[0].Variants[0].ID = "" }},
		{"variant name", func(b *batch) { b.Families[8].Variants[1].Name = "Belt of Giant Strength (frost)" }},
		{"rarity", func(b *batch) { b.Families[8].Variants[2].Rarity = "legendary" }},
		{"unknown rarity", func(b *batch) { b.Families[0].Variants[0].Rarity = "mythic" }},
		{"dropped attunement", func(b *batch) { b.Families[11].Variants[0].Attunement.Mode = "none" }},
		{"invented attunement", func(b *batch) { b.Families[10].Variants[0].Attunement.Mode = "required" }},
		{"attunement restriction", func(b *batch) { b.Families[0].Variants[0].Attunement.Restrictions = []string{"spellcaster"} }},
		{"missing charges", func(b *batch) { b.Families[0].Variants[0].Charges = nil }},
		{"invented charges", func(b *batch) { b.Families[4].Variants[0].Charges.Maximum = 3 }},
		{"recharge", func(b *batch) { b.Families[4].Variants[0].Charges.Recharge = "dawn" }},
		{"equip mode", func(b *batch) { b.Families[0].Variants[0].Equip.Mode = "wear" }},
		{"equip restriction", func(b *batch) { b.Families[0].Variants[0].Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { b.Families[0].Variants[0].Use.Mode = "free-action" }},
		{"dropped use restriction", func(b *batch) { b.Families[4].Variants[0].Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { b.Families[0].Variants[0].Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { b.Families[3].Variants[0].Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { b.Families[0].Variants[0].Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { b.Families[0].Variants[0].Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) {
			b.Families[0].Variants[0].RequiredMechanicKeys = b.Families[0].Variants[0].RequiredMechanicKeys[:1]
		}},
		{"no mechanics", func(b *batch) { b.Families[0].Variants[0].RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { b.Families[0].Variants[0].RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"duplicate mechanic", func(b *batch) {
			b.Families[0].Variants[0].RequiredMechanicKeys[1] = b.Families[0].Variants[0].RequiredMechanicKeys[0]
		}},
		{"missing details", func(b *batch) { b.Families[0].Variants[0].Details = nil }},
		{"empty details", func(b *batch) { b.Families[0].Variants[0].Details = &ItemDetails226{} }},
		{"mismatched payload", func(b *batch) { b.Families[0].Variants[0].Details = b.Families[11].Variants[0].Details }},
		{"second payload", func(b *batch) {
			b.Families[0].Variants[0].Details.Levitation = b.Families[11].Variants[0].Details.Levitation
		}},
		{"blank key", func(b *batch) { b.Families[0].Variants[0].Details.ArrowCatching.Condition = " " }},
		{"prose value", func(b *batch) {
			b.Families[0].Variants[0].Details.ArrowCatching.Condition = "While you wield this Shield"
		}},
		{"shield bonus", func(b *batch) { b.Families[0].Variants[0].Details.ArrowCatching.ACBonus = 0 }},
		{"shield negative range", func(b *batch) { b.Families[0].Variants[0].Details.ArrowCatching.RedirectRangeFeet = -5 }},
		{"bean table gap", func(b *batch) { b.Families[1].Variants[0].Details.Beans.Effects[1].Min = 3 }},
		{"bean table short", func(b *batch) {
			b.Families[1].Variants[0].Details.Beans.Effects = b.Families[1].Variants[0].Details.Beans.Effects[:11]
		}},
		{"bean table end", func(b *batch) { b.Families[1].Variants[0].Details.Beans.Effects[11].Max = 99 }},
		{"bean duplicate effect", func(b *batch) { b.Families[1].Variants[0].Details.Beans.Effects[2].Effect = "geyser-erupts" }},
		{"bean quantity without die", func(b *batch) { b.Families[1].Variants[0].Details.Beans.Effects[0].Quantities[0].Sides = 0 }},
		{"bean empty quantity", func(b *batch) { b.Families[1].Variants[0].Details.Beans.Effects[2].Quantities[0].Bonus = 0 }},
		{"bean quantities dropped", func(b *batch) { b.Families[1].Variants[0].Details.Beans.Effects[5].Quantities = []ItemBeanQuantity{} }},
		{"bean duplicate fact", func(b *batch) {
			b.Families[1].Variants[0].Details.Beans.Effects[2].Facts[1] = b.Families[1].Variants[0].Details.Beans.Effects[2].Facts[0]
		}},
		{"bean save dc", func(b *batch) { b.Families[1].Variants[0].Details.Beans.SaveDC = 0 }},
		{"bean damage dice", func(b *batch) { b.Families[1].Variants[0].Details.Beans.Damage.Dice = -5 }},
		{"bean weight", func(b *batch) { b.Families[1].Variants[0].Details.Beans.WeightPoundsDenominator = 0 }},
		{"devouring chance", func(b *batch) { b.Families[2].Variants[0].Details.Devouring.PullChancePercent = 150 }},
		{"devouring escape check", func(b *batch) { b.Families[2].Variants[0].Details.Devouring.EscapeCheck.DC = 0 }},
		{"devouring rescue check", func(b *batch) { b.Families[2].Variants[0].Details.Devouring.RescueCheck.Skill = "" }},
		{"devouring destroyed by", func(b *batch) { b.Families[2].Variants[0].Details.Devouring.DestroyedBy = nil }},
		{"holding capacity", func(b *batch) { b.Families[3].Variants[0].Details.Holding.CapacityPounds = 0 }},
		{"holding duplicate destroyer", func(b *batch) { b.Families[3].Variants[0].Details.Holding.DestroyedBy[2] = "pierced" }},
		{"holding gate radius", func(b *batch) { b.Families[3].Variants[0].Details.Holding.GateRadiusFeet = 0 }},
		{"tricks colour", func(b *batch) { b.Families[4].Variants[0].Details.Tricks.Color = "rust" }},
		{"tricks swapped table", func(b *batch) {
			b.Families[4].Variants[0].Details, b.Families[4].Variants[1].Details = b.Families[4].Variants[1].Details, b.Families[4].Variants[0].Details
		}},
		{"tricks table short", func(b *batch) {
			b.Families[4].Variants[1].Details.Tricks.Creatures = b.Families[4].Variants[1].Details.Tricks.Creatures[:7]
		}},
		{"tricks table order", func(b *batch) { b.Families[4].Variants[2].Details.Tricks.Creatures[3].Roll = 5 }},
		{"tricks duplicate creature", func(b *batch) { b.Families[4].Variants[0].Details.Tricks.Creatures[7].Creature = "weasel" }},
		{"tricks uses", func(b *batch) { b.Families[4].Variants[0].Details.Tricks.Uses = 0 }},
		{"force continuation page", func(b *batch) { b.Families[5].Variants[0].Details.Force.SourcePages = []int{212, 212} }},
		{"force save dc", func(b *batch) { b.Families[5].Variants[0].Details.Force.SaveDC = 0 }},
		{"force damage", func(b *batch) { b.Families[5].Variants[0].Details.Force.Damage.Sides = 0 }},
		{"force duration", func(b *batch) { b.Families[5].Variants[0].Details.Force.SphereDurationSeconds = 0 }},
		{"force blocked", func(b *batch) { b.Families[5].Variants[0].Details.Force.Blocked = []string{} }},
		{"nourishment days", func(b *batch) { b.Families[6].Variants[0].Details.Nourishment.NourishmentDays = 0 }},
		{"dwarvenkind maximum", func(b *batch) { b.Families[7].Variants[0].Details.Dwarvenkind.ConstitutionMaximum = 31 }},
		{"dwarvenkind chance", func(b *batch) { b.Families[7].Variants[0].Details.Dwarvenkind.BeardChancePercent = 0 }},
		{"dwarvenkind duplicate benefit", func(b *batch) { b.Families[7].Variants[0].Details.Dwarvenkind.Benefits[1] = "dwarvish" }},
		{"beans save ability", func(b *batch) { b.Families[1].Variants[0].Details.Beans.Save = "wisdom" }},
		{"beans save dc literal", func(b *batch) { b.Families[1].Variants[0].Details.Beans.SaveDC = 20 }},
		{"devouring escape dc literal", func(b *batch) { b.Families[2].Variants[0].Details.Devouring.EscapeCheck.DC = 20 }},
		{"devouring rescue dc literal", func(b *batch) { b.Families[2].Variants[0].Details.Devouring.RescueCheck.DC = 25 }},
		{"tricks gray creature literal", func(b *batch) { b.Families[4].Variants[0].Details.Tricks.Creatures[0].Creature = "rat" }},
		{"tricks invented creature", func(b *batch) { b.Families[4].Variants[0].Details.Tricks.Creatures[0].Creature = "dragon" }},
		{"tricks swapped tables", func(b *batch) {
			rust, tan := b.Families[4].Variants[1].Details.Tricks, b.Families[4].Variants[2].Details.Tricks
			rust.Creatures, tan.Creatures = tan.Creatures, rust.Creatures
		}},
		{"tricks reordered rows", func(b *batch) {
			c := b.Families[4].Variants[2].Details.Tricks.Creatures
			c[6].Creature, c[7].Creature = c[7].Creature, c[6].Creature
		}},
		{"belt score", func(b *batch) { b.Families[8].Variants[0].Details.GiantStrength.Score = 23 }},
		{"belt threshold", func(b *batch) { b.Families[8].Variants[4].Details.GiantStrength.NoEffectIfScoreAtLeast = 30 }},
		{"belt giant types", func(b *batch) { b.Families[8].Variants[1].Details.GiantStrength.GiantTypes = []string{"frost"} }},
		{"belt swapped rows", func(b *batch) {
			b.Families[8].Variants[2].Details, b.Families[8].Variants[3].Details = b.Families[8].Variants[3].Details, b.Families[8].Variants[2].Details
		}},
		{"axe bonus", func(b *batch) { b.Families[9].Variants[0].Details.Berserker.DamageBonus = 2 }},
		{"axe uncursed", func(b *batch) { b.Families[9].Variants[0].Details.Berserker.Cursed = false }},
		{"axe save dc", func(b *batch) { b.Families[9].Variants[0].Details.Berserker.BerserkSaveDC = 0 }},
		{"elvenkind prose", func(b *batch) { b.Families[10].Variants[0].Details.Elvenkind.AdvantageSkill = "Dexterity (Stealth)" }},
		{"levitation spell", func(b *batch) { b.Families[11].Variants[0].Details.Levitation.Spell = "" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems226(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems226(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems226RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/arrow-catching-shield/variant/standard"`)
	if !bytes.Contains(items226JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"current_charges":3,`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items226JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems226(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"arrow_catching": {`, `"arrow_catching_plus": {`},
		"unknown detail":     {`"adds_to_shield_bonus"`, `"replaces_shield_bonus"`},
		"string for number":  {`"save_dc": 15`, `"save_dc": "15"`},
		"fractional number":  {`"save_dc": 15`, `"save_dc": 15.5`},
		"number for boolean": {`"adds_to_shield_bonus": true`, `"adds_to_shield_bonus": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items226JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items226JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems226(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items226JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems226(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// The Healing contract decoder must not accept a batch, nor the reverse.
	if _, err := DecodeItemCatalogue(items226JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems226(itemFamiliesJSON); err == nil {
		t.Fatal("batch decoder accepted the contract fixture")
	}
}
