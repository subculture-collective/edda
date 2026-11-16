package srd521

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func invocations02Catalogue(t *testing.T) (CharacterCatalogue, CharacterCatalogue) {
	t.Helper()
	raw, err := os.ReadFile("testdata/character-invocations-02-composition.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Catalogues []string `json:"catalogues"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(manifest.Catalogues, []string{"character-warlock.json", "character-origin-feats.json", "character-invocations-02.json"}) {
		t.Fatal("composition changed")
	}
	var batch, composed CharacterCatalogue
	for _, file := range manifest.Catalogues {
		raw, err = os.ReadFile("testdata/" + file)
		if err != nil {
			t.Fatal(err)
		}
		c, err := DecodeCharacterCatalogue(strings.NewReader(string(raw)))
		if err != nil {
			t.Fatal(err)
		}
		if file == "character-invocations-02.json" {
			batch = c
		}
		composed.SchemaVersion, composed.Source = c.SchemaVersion, c.Source
		composed.Records = append(composed.Records, c.Records...)
	}
	return batch, composed
}

func invocations02SourceError(c CharacterCatalogue) error {
	want := invocations02SourceRecords()
	if len(c.Records) != len(want) {
		return fmt.Errorf("record count %d, want %d", len(c.Records), len(want))
	}
	for i, r := range c.Records {
		if !reflect.DeepEqual(r, want[i]) {
			return fmt.Errorf("pp73–74 metadata %s:\n got %+v\nwant %+v", r.ID, r, want[i])
		}
	}
	return nil
}

func TestInvocations02SourceCompositionAndUnavailable(t *testing.T) {
	batch, c := invocations02Catalogue(t)
	repo := os.DirFS("../../..")
	if err := c.Validate(repo); err != nil {
		t.Fatal(err)
	}
	if err := invocations02SourceError(batch); err != nil {
		t.Fatal(err)
	}
	if batch.Validate(repo) == nil {
		t.Fatal("missing actual owners accepted")
	}
	if c.ValidateReadiness(repo) == nil {
		t.Fatal("metadata playable")
	}
	TestWarlockSourceAndBoundaries(t)
	TestCharacterCatalogueDoesNotPromoteFixedBuilds(t)
	if !reflect.DeepEqual(c.Records[:len(warlockCatalogue(t).Records)], warlockCatalogue(t).Records) {
		t.Fatal("canonical Warlock changed")
	}
	for _, r := range batch.Records {
		for _, group := range [][]string{r.RequiredMechanics, r.OptionKeys} {
			for _, k := range group {
				if qualifiedCharacterMechanic(k) {
					t.Fatalf("qualified %s", k)
				}
			}
		}
	}
	for _, records := range [][]CharacterRecord{batch.Records, warlockCatalogue(t).Records} {
		duplicate := c
		duplicate.Records = append(append([]CharacterRecord{}, c.Records...), records...)
		if duplicate.Validate(repo) == nil {
			t.Fatal("duplicate batch/Warlock accepted")
		}
	}
}

func TestInvocations02ContrarySourceBothProofs(t *testing.T) {
	batch, _ := invocations02Catalogue(t)
	repo := os.DirFS("../../..")
	var controls [][]CharacterRecord
	check := func(name string, mutate func(*CharacterCatalogue)) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			b, c := invocations02Catalogue(t)
			mutate(&b)
			c.Records = append(c.Records[:len(c.Records)-len(b.Records)], b.Records...)
			if err := c.Validate(repo); err != nil {
				t.Fatalf("contrary graph must remain structurally valid: %v", err)
			}
			if invocations02SourceError(b) == nil {
				t.Fatal("Go source proof accepted contrary graph")
			}
			if c.ValidateReadiness(repo) == nil {
				t.Fatal("contrary graph playable")
			}
			controls = append(controls, b.Records)
		})
	}
	for i, r := range batch.Records {
		for j := range r.RequiredMechanics {
			check(fmt.Sprintf("condition-or-source-variant-%d-%d", i, j), func(b *CharacterCatalogue) {
				q := &b.Records[i]
				old := q.RequiredMechanics[j]
				key := "srd-5.2.1:mechanic/invocations/contrary-variant"
				q.RequiredMechanics[j] = key
				if v, ok := q.Parameters[old]; ok {
					delete(q.Parameters, old)
					q.Parameters[key] = v
				}
				if q.Requirement != nil && q.Requirement.Key == old {
					q.Requirement.Key = key
				}
			})
		}
		for key := range r.Parameters {
			check(fmt.Sprintf("number-cost-timing-%d-%s", i, key), func(b *CharacterCatalogue) { b.Records[i].Parameters[key]++ })
		}
		if r.Requirement != nil {
			for key := range r.Requirement.Parameters {
				check(fmt.Sprintf("predicate-number-%d-%s", i, key), func(b *CharacterCatalogue) { b.Records[i].Requirement.Parameters[key]++ })
			}
			for j := range r.Requirement.RecordRefs {
				check(fmt.Sprintf("known-wrong-operand-%d-%d", i, j), func(b *CharacterCatalogue) {
					b.Records[i].Requirement.RecordRefs[j] = "srd-5.2.1:character/choice/warlock/skills"
				})
			}
		}
		if r.Kind == "choice" {
			check(fmt.Sprintf("selection-count-%d", i), func(b *CharacterCatalogue) { b.Records[i].Select++ })
			for j := range r.OptionKeys {
				check(fmt.Sprintf("wrong-alternative-%d-%d", i, j), func(b *CharacterCatalogue) { b.Records[i].OptionKeys[j] = "srd-5.2.1:mechanic/damage/force" })
				if len(r.OptionKeys) > 1 {
					check(fmt.Sprintf("missing-alternative-%d-%d", i, j), func(b *CharacterCatalogue) {
						q := &b.Records[i]
						q.OptionKeys = append(q.OptionKeys[:j], q.OptionKeys[j+1:]...)
					})
				}
			}
		}
		if r.Kind == "invocation" {
			check(fmt.Sprintf("invented-scaling-or-recovery-%d", i), func(b *CharacterCatalogue) {
				q := &b.Records[i]
				key := "srd-5.2.1:mechanic/invocations/scale-per-warlock-level-recover-short-rest"
				q.RequiredMechanics = append(q.RequiredMechanics, key)
				if q.Parameters == nil {
					q.Parameters = map[string]int{}
				}
				q.Parameters[key] = 2
			})
			check(fmt.Sprintf("missing-condition-%d", i), func(b *CharacterCatalogue) { b.Records[i].RequiredMechanics = b.Records[i].RequiredMechanics[1:] })
		}
	}
	raw, err := json.Marshal(controls)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "contrary.json")
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("python3", "../../../scripts/verify_srd_invocations_02.py", "--controls", path)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("standalone semantic proof: %v\n%s", err, out)
	}
	t.Logf("%d structurally valid contrary graphs rejected by both proofs and unavailable", len(controls))
}

func TestInvocations02SchemaUnknownAndPDF(t *testing.T) {
	batch, c := invocations02Catalogue(t)
	raw, err := json.Marshal(batch)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{
		strings.Replace(string(raw), `"kind":`, `"playable":true,"kind":`, 1),
		strings.Replace(string(raw), `"requirement":{`, `"requirement":{"execute":true,`, 1),
		strings.Replace(string(raw), `"schema_version":1`, `"schema_version":2`, 1),
		strings.Replace(string(raw), batch.Source.SHA256, strings.Repeat("0", 64), 1),
		strings.Replace(string(raw), `"version":"5.2.1"`, `"version":"5.1"`, 1),
		string(raw) + ` {}`,
	} {
		b, err := DecodeCharacterCatalogue(strings.NewReader(bad))
		if err == nil {
			b.Records = append(append([]CharacterRecord{}, c.Records[:len(c.Records)-len(batch.Records)]...), b.Records...)
			if b.Validate(os.DirFS("../../..")) == nil {
				t.Fatal("invalid schema/source accepted")
			}
		}
	}
	for _, mutate := range []func(*CharacterCatalogue){
		func(b *CharacterCatalogue) { b.Records[0].ParentID = "srd-5.2.1:character/feat/alert" },
		func(b *CharacterCatalogue) { b.Records[0].ID += "-wrong" },
		func(b *CharacterCatalogue) { b.Records[0].Sources[0].Page = 74 },
		func(b *CharacterCatalogue) {
			b.Records[0].Sources[0].InventoryID = "srd-5.2.1:planning/feats/alert"
			b.Records[0].Sources[0].Page = 87
		},
	} {
		_, b := invocations02Catalogue(t)
		start := len(b.Records) - len(batch.Records)
		isolated := CharacterCatalogue{SchemaVersion: b.SchemaVersion, Source: b.Source, Records: b.Records[start:]}
		mutate(&isolated)
		if b.Validate(os.DirFS("../../..")) == nil {
			t.Fatal("invalid invocation guard accepted")
		}
	}
	key := "srd-5.2.1:mechanic/invocations/unknown"
	c.Records[0].RequiredMechanics = append([]string{key}, c.Records[0].RequiredMechanics...)
	if err := c.Validate(os.DirFS("../../..")); err != nil {
		t.Fatal(err)
	}
	if err := c.ValidateReadiness(os.DirFS("../../..")); err == nil || !strings.Contains(err.Error(), key) {
		t.Fatal("unknown key available")
	}
	pdf := os.Getenv("EDDA_SRD_PDF")
	if pdf == "" {
		t.Skip("set EDDA_SRD_PDF for offline pinned PDF proof")
	}
	cmd := exec.Command("python3", "../../../scripts/verify_srd_invocations_02.py", pdf)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("PDF proof: %v\n%s", err, out)
	}
}
