package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// SpellSourceDetails retains the batch's additional source shapes as offline
// metadata. These fields never select or implement a runtime resolver.
type SpellSourceDetails struct {
	Antipathy         *SpellAntipathy       `json:"antipathy,omitempty"`
	Eye               *SpellEye             `json:"eye,omitempty"`
	Hand              *SpellHand            `json:"hand,omitempty"`
	Lock              *SpellLock            `json:"lock,omitempty"`
	Sword             *SpellSword           `json:"sword,omitempty"`
	MagicAura         *SpellMagicAura       `json:"magic_aura,omitempty"`
	Astral            *SpellAstral          `json:"astral,omitempty"`
	Augury            *SpellAugury          `json:"augury,omitempty"`
	Life              *SpellLife            `json:"life,omitempty"`
	Awaken            *SpellAwaken          `json:"awaken,omitempty"`
	Bane              *SpellBane            `json:"bane,omitempty"`
	Banishment        *SpellBanishment      `json:"banishment,omitempty"`
	Barkskin          *SpellBarkskin        `json:"barkskin,omitempty"`
	Beacon            *SpellBeacon          `json:"beacon,omitempty"`
	Befuddlement      *SpellBefuddlement    `json:"befuddlement,omitempty"`
	Curse             *SpellCurse           `json:"curse,omitempty"`
	Tentacles         *SpellTentacles       `json:"tentacles,omitempty"`
	Barrier           *SpellBarrier         `json:"barrier,omitempty"`
	Bless             *SpellBless           `json:"bless,omitempty"`
	Blight            *SpellBlight          `json:"blight,omitempty"`
	Blindness         *SpellBlindness       `json:"blindness,omitempty"`
	Blink             *SpellBlink           `json:"blink,omitempty"`
	Blur              *SpellBlur            `json:"blur,omitempty"`
	Burning           *SpellBurning         `json:"burning,omitempty"`
	Lightning         *SpellLightning       `json:"lightning,omitempty"`
	Calm              *SpellCalm            `json:"calm,omitempty"`
	Chain             *SpellChain           `json:"chain,omitempty"`
	Monster           *SpellCharm           `json:"monster,omitempty"`
	Person            *SpellCharm           `json:"person,omitempty"`
	Chill             *SpellChill           `json:"chill,omitempty"`
	Orb               *SpellOrb             `json:"orb,omitempty"`
	Death             *SpellDeath           `json:"death,omitempty"`
	Clairvoyance      *SpellClairvoyance    `json:"clairvoyance,omitempty"`
	Clone             *SpellClone           `json:"clone,omitempty"`
	Cloudkill         *SpellCloudkill       `json:"cloudkill,omitempty"`
	Spray             *SpellSpray           `json:"spray,omitempty"`
	Command           *SpellCommand         `json:"command,omitempty"`
	Commune           *SpellCommune         `json:"commune,omitempty"`
	Nature            *SpellNature          `json:"nature,omitempty"`
	Comprehend        *SpellComprehend      `json:"comprehend,omitempty"`
	Compulsion        *SpellCompulsion      `json:"compulsion,omitempty"`
	Cone              *SpellCone            `json:"cone,omitempty"`
	Confusion         *SpellConfusion       `json:"confusion,omitempty"`
	Animals           *SpellAnimals         `json:"animals,omitempty"`
	Celestial         *SpellCelestial       `json:"celestial,omitempty"`
	Elemental         *SpellElemental       `json:"elemental,omitempty"`
	Fey               *SpellFey             `json:"fey,omitempty"`
	Minor             *SpellMinor           `json:"minor,omitempty"`
	Woodland          *SpellWoodland        `json:"woodland,omitempty"`
	OtherPlane        *SpellOtherPlane      `json:"other_plane,omitempty"`
	Contagion         *SpellContagion       `json:"contagion,omitempty"`
	Contingency       *SpellContingency     `json:"contingency,omitempty"`
	Flame             *SpellFlame           `json:"flame,omitempty"`
	WaterControl      *SpellWaterControl    `json:"water_control,omitempty"`
	Weather           *SpellWeather         `json:"weather,omitempty"`
	Counter           *SpellCounter         `json:"counter,omitempty"`
	Provisions        *SpellProvisions      `json:"provisions,omitempty"`
	Undead            *SpellUndead          `json:"undead,omitempty"`
	WaterCreate       *SpellWaterCreate     `json:"water_create,omitempty"`
	Creation          *SpellCreation        `json:"creation,omitempty"`
	CastingSeconds    int                   `json:"casting_seconds"`
	CastingTriggers   []string              `json:"casting_triggers"`
	EndTiming         string                `json:"end_timing"`
	Rays              *SpellRays            `json:"rays,omitempty"`
	Scrying           *SpellScrying         `json:"scrying,omitempty"`
	Smite             *SpellSmite           `json:"smite,omitempty"`
	Chest             *SpellChest           `json:"chest,omitempty"`
	Vision            *SpellVision          `json:"vision,omitempty"`
	Seeming           *SpellSeeming         `json:"seeming,omitempty"`
	Sending           *SpellSending         `json:"sending,omitempty"`
	Sequester         *SpellSequester       `json:"sequester,omitempty"`
	Shapechange       *SpellShapechange     `json:"shapechange,omitempty"`
	Shatter           *SpellShatter         `json:"shatter,omitempty"`
	Shield            *SpellShield          `json:"shield,omitempty"`
	Faith             *SpellFaith           `json:"faith,omitempty"`
	Cure              *SpellCure            `json:"cure,omitempty"`
	Lights            *SpellLights          `json:"lights,omitempty"`
	Darkness          *SpellDarkness        `json:"darkness,omitempty"`
	Darkvision        *SpellDarkvision      `json:"darkvision,omitempty"`
	Daylight          *SpellDaylight        `json:"daylight,omitempty"`
	DeathWard         *SpellDeathWard       `json:"death_ward,omitempty"`
	Fireball          *SpellFireball        `json:"fireball,omitempty"`
	Demiplane         *SpellDemiplane       `json:"demiplane,omitempty"`
	DetectEvil        *SpellDetectEvil      `json:"detect_evil,omitempty"`
	DetectMagic       *SpellDetectMagic     `json:"detect_magic,omitempty"`
	DetectPoison      *SpellDetectPoison    `json:"detect_poison,omitempty"`
	DetectThoughts    *SpellDetectThoughts  `json:"detect_thoughts,omitempty"`
	DimensionDoor     *SpellDimensionDoor   `json:"dimension_door,omitempty"`
	DisguiseSelf      *SpellDisguiseSelf    `json:"disguise_self,omitempty"`
	Disintegrate      *SpellDisintegrate    `json:"disintegrate,omitempty"`
	DispelEvil        *SpellDispelEvil      `json:"dispel_evil,omitempty"`
	DispelMagic       *SpellDispelMagic     `json:"dispel_magic,omitempty"`
	DissonantWhispers *SpellWhispers        `json:"dissonant_whispers,omitempty"`
	Divination        *SpellDivination      `json:"divination,omitempty"`
	DivineFavor       *SpellDivineFavor     `json:"divine_favor,omitempty"`
	DivineSmite       *SpellDivineSmite     `json:"divine_smite,omitempty"`
	DivineWord        *SpellDivineWord      `json:"divine_word,omitempty"`
	DominateBeast     *SpellDominate        `json:"dominate_beast,omitempty"`
	DominateMonster   *SpellDominate        `json:"dominate_monster,omitempty"`
	Dominate          *SpellDominatePerson  `json:"dominate,omitempty"`
	Breath            *SpellDragonsBreath   `json:"breath,omitempty"`
	Dream             *SpellDream           `json:"dream,omitempty"`
	Druidcraft        *SpellDruidcraft      `json:"druidcraft,omitempty"`
	Earthquake        *SpellEarthquake      `json:"earthquake,omitempty"`
	Blast             *SpellEldritchBlast   `json:"blast,omitempty"`
	Elementalism      *SpellElementalism    `json:"elementalism,omitempty"`
	Enhance           *SpellEnhanceAbility  `json:"enhance,omitempty"`
	Enlarge           *SpellEnlargeReduce   `json:"enlarge,omitempty"`
	Ensnaring         *SpellEnsnaringStrike `json:"ensnaring,omitempty"`
	Entangle          *SpellEntangle        `json:"entangle,omitempty"`
	Enthrall          *SpellEnthrall        `json:"enthrall,omitempty"`

	// Issue #205 payloads.
	Etherealness *SpellEtherealness `json:"etherealness,omitempty"`
	Retreat      *SpellRetreat      `json:"retreat,omitempty"`
	Eyebite      *SpellEyebite      `json:"eyebite,omitempty"`
	Fabricate    *SpellFabricate    `json:"fabricate,omitempty"`
	FaerieFire   *SpellFaerieFire   `json:"faerie_fire,omitempty"`
	Hound        *SpellHound        `json:"hound,omitempty"`
	FalseLife    *SpellFalseLife    `json:"false_life,omitempty"`
	Fear         *SpellFear         `json:"fear,omitempty"`
	FeatherFall  *SpellFeatherFall  `json:"feather_fall,omitempty"`
	Familiar     *SpellFamiliar     `json:"familiar,omitempty"`
	Steed        *SpellSteed        `json:"steed,omitempty"`
	Traps        *SpellTraps        `json:"traps,omitempty"`

	// Issue #206 payloads.
	FindPath      *SpellFindPath      `json:"find_path,omitempty"`
	FingerOfDeath *SpellFingerOfDeath `json:"finger_of_death,omitempty"`
	FireBolt      *SpellFireBolt      `json:"fire_bolt,omitempty"`
	FireShield    *SpellFireShield    `json:"fire_shield,omitempty"`
	FireStorm     *SpellFireStorm     `json:"fire_storm,omitempty"`
	FireballBurst *SpellFireballBurst `json:"fireball_burst,omitempty"`
	FlameBlade    *SpellFlameBlade    `json:"flame_blade,omitempty"`
	FlameStrike   *SpellFlameStrike   `json:"flame_strike,omitempty"`
	FlamingSphere *SpellFlamingSphere `json:"flaming_sphere,omitempty"`
	FleshToStone  *SpellFleshToStone  `json:"flesh_to_stone,omitempty"`
	FloatingDisk  *SpellFloatingDisk  `json:"floating_disk,omitempty"`
	Fly           *SpellFly           `json:"fly,omitempty"`
	Fog           *SpellFog           `json:"fog,omitempty"`
	Forbiddance   *SpellForbiddance   `json:"forbiddance,omitempty"`
	Forcecage     *SpellForcecage     `json:"forcecage,omitempty"`
	Foresight     *SpellForesight     `json:"foresight,omitempty"`
	Freedom       *SpellFreedom       `json:"freedom,omitempty"`
	Freezing      *SpellFreezing      `json:"freezing,omitempty"`
	Gaseous       *SpellGaseous       `json:"gaseous,omitempty"`
	Gate          *SpellGate          `json:"gate,omitempty"`
	Geas          *SpellGeas          `json:"geas,omitempty"`
	Repose        *SpellRepose        `json:"repose,omitempty"`
	Insect        *SpellInsect        `json:"insect,omitempty"`
	Glibness      *SpellGlibness      `json:"glibness,omitempty"`

	// Issue #196 payloads. Acid Arrow needs only the common record fields.
	Splash     *SpellSplash         `json:"splash,omitempty"`
	Aid        *SpellAid            `json:"aid,omitempty"`
	Alarm      *SpellAlarm          `json:"alarm,omitempty"`
	Alter      *SpellAlterSelf      `json:"alter,omitempty"`
	Friendship *SpellFriendship     `json:"friendship,omitempty"`
	Messenger  *SpellMessenger      `json:"messenger,omitempty"`
	Shapes     *SpellShapes         `json:"shapes,omitempty"`
	Dead       *SpellAnimateDead    `json:"dead,omitempty"`
	Objects    *SpellAnimateObjects `json:"objects,omitempty"`
	Antilife   *SpellAntilife       `json:"antilife,omitempty"`
	Antimagic  *SpellAntimagic      `json:"antimagic,omitempty"`

	// Issue #208 payloads.
	Globe               *SpellGlobe               `json:"globe,omitempty"`
	Glyph               *SpellGlyph               `json:"glyph,omitempty"`
	Goodberry           *SpellGoodberry           `json:"goodberry,omitempty"`
	Grease              *SpellGrease              `json:"grease,omitempty"`
	GreaterInvisibility *SpellGreaterInvisibility `json:"greater_invisibility,omitempty"`
	GreaterRestoration  *SpellGreaterRestoration  `json:"greater_restoration,omitempty"`
	Guardian            *SpellGuardian            `json:"guardian,omitempty"`
	Guards              *SpellGuards              `json:"guards,omitempty"`
	Guidance            *SpellGuidance            `json:"guidance,omitempty"`
	GuidingBolt         *SpellGuidingBolt         `json:"guiding_bolt,omitempty"`
	Gust                *SpellGust                `json:"gust,omitempty"`
	Hallow              *SpellHallow              `json:"hallow,omitempty"`

	// Issue #209 payloads.
	HallucinatoryTerrain *SpellHallucinatoryTerrain `json:"hallucinatory_terrain,omitempty"`
	Harm                 *SpellHarm                 `json:"harm,omitempty"`
	Haste                *SpellHaste                `json:"haste,omitempty"`
	Heal                 *SpellHeal                 `json:"heal,omitempty"`
	HealingWord          *SpellHealingWord          `json:"healing_word,omitempty"`
	HeatMetal            *SpellHeatMetal            `json:"heat_metal,omitempty"`
	HellishRebuke        *SpellHellishRebuke        `json:"hellish_rebuke,omitempty"`
	HeroesFeast          *SpellHeroesFeast          `json:"heroes_feast,omitempty"`
	Heroism              *SpellHeroism              `json:"heroism,omitempty"`
	Hex                  *SpellHex                  `json:"hex,omitempty"`
	HideousLaughter      *SpellHideousLaughter      `json:"hideous_laughter,omitempty"`
	HoldMonster          *SpellHoldMonster          `json:"hold_monster,omitempty"`

	// Issue #210 payloads.
	HoldPerson      *SpellHoldPerson      `json:"hold_person,omitempty"`
	HolyAura        *SpellHolyAura        `json:"holy_aura,omitempty"`
	HuntersMark     *SpellHuntersMark     `json:"hunters_mark,omitempty"`
	HypnoticPattern *SpellHypnoticPattern `json:"hypnotic_pattern,omitempty"`
	IceKnife        *SpellIceKnife        `json:"ice_knife,omitempty"`
	IceStorm        *SpellIceStorm        `json:"ice_storm,omitempty"`
	Identify        *SpellIdentify        `json:"identify,omitempty"`
	IllusoryScript  *SpellIllusoryScript  `json:"illusory_script,omitempty"`
	Imprisonment    *SpellImprisonment    `json:"imprisonment,omitempty"`
	IncendiaryCloud *SpellIncendiaryCloud `json:"incendiary_cloud,omitempty"`
	InflictWounds   *SpellInflictWounds   `json:"inflict_wounds,omitempty"`
	InsectPlague    *SpellInsectPlague    `json:"insect_plague,omitempty"`

	// Issue #211 payloads.
	InstantSummons        *SpellInstantSummons        `json:"instant_summons,omitempty"`
	Invisibility          *SpellInvisibility          `json:"invisibility,omitempty"`
	IrresistibleDance     *SpellIrresistibleDance     `json:"irresistible_dance,omitempty"`
	Jump                  *SpellJump                  `json:"jump,omitempty"`
	Knock                 *SpellKnock                 `json:"knock,omitempty"`
	LegendLore            *SpellLegendLore            `json:"legend_lore,omitempty"`
	LesserRestoration     *SpellLesserRestoration     `json:"lesser_restoration,omitempty"`
	Levitate              *SpellLevitate              `json:"levitate,omitempty"`
	Light                 *SpellLight                 `json:"light,omitempty"`
	LightningBolt         *SpellLightningBolt         `json:"lightning_bolt,omitempty"`
	LocateAnimalsOrPlants *SpellLocateAnimalsOrPlants `json:"locate_animals_or_plants,omitempty"`
	LocateCreature        *SpellLocateCreature        `json:"locate_creature,omitempty"`

	// Issue #212 payloads.
	LocateObject       *SpellLocateObject       `json:"locate_object,omitempty"`
	Longstrider        *SpellLongstrider        `json:"longstrider,omitempty"`
	MageArmor          *SpellMageArmor          `json:"mage_armor,omitempty"`
	MageHand           *SpellMageHand           `json:"mage_hand,omitempty"`
	MagicCircle        *SpellMagicCircle        `json:"magic_circle,omitempty"`
	MagicJar           *SpellMagicJar           `json:"magic_jar,omitempty"`
	MagicMissile       *SpellMagicMissile       `json:"magic_missile,omitempty"`
	MagicMouth         *SpellMagicMouth         `json:"magic_mouth,omitempty"`
	MagicWeapon        *SpellMagicWeapon        `json:"magic_weapon,omitempty"`
	MagnificentMansion *SpellMagnificentMansion `json:"magnificent_mansion,omitempty"`
	MajorImage         *SpellMajorImage         `json:"major_image,omitempty"`
	MassCureWounds     *SpellMassCureWounds     `json:"mass_cure_wounds,omitempty"`

	// Issue #213 payloads.
	MassHeal        *SpellMassHeal        `json:"mass_heal,omitempty"`
	MassHealingWord *SpellMassHealingWord `json:"mass_healing_word,omitempty"`
	MassSuggestion  *SpellMassSuggestion  `json:"mass_suggestion,omitempty"`
	Maze            *SpellMaze            `json:"maze,omitempty"`
	MeldIntoStone   *SpellMeldIntoStone   `json:"meld_into_stone,omitempty"`
	Mending         *SpellMending         `json:"mending,omitempty"`
	Message         *SpellMessage         `json:"message,omitempty"`
	MeteorSwarm     *SpellMeteorSwarm     `json:"meteor_swarm,omitempty"`
	MindBlank       *SpellMindBlank       `json:"mind_blank,omitempty"`
	MindSpike       *SpellMindSpike       `json:"mind_spike,omitempty"`
	MinorIllusion   *SpellMinorIllusion   `json:"minor_illusion,omitempty"`
	MirageArcane    *SpellMirageArcane    `json:"mirage_arcane,omitempty"`

	// Issue #214 payloads.
	MirrorImage      *SpellMirrorImage      `json:"mirror_image,omitempty"`
	Mislead          *SpellMislead          `json:"mislead,omitempty"`
	MistyStep        *SpellMistyStep        `json:"misty_step,omitempty"`
	ModifyMemory     *SpellModifyMemory     `json:"modify_memory,omitempty"`
	Moonbeam         *SpellMoonbeam         `json:"moonbeam,omitempty"`
	MoveEarth        *SpellMoveEarth        `json:"move_earth,omitempty"`
	Nondetection     *SpellNondetection     `json:"nondetection,omitempty"`
	PassWithoutTrace *SpellPassWithoutTrace `json:"pass_without_trace,omitempty"`
	Passwall         *SpellPasswall         `json:"passwall,omitempty"`
	PhantasmalForce  *SpellPhantasmalForce  `json:"phantasmal_force,omitempty"`
	PhantasmalKiller *SpellPhantasmalKiller `json:"phantasmal_killer,omitempty"`
	PhantomSteed     *SpellPhantomSteed     `json:"phantom_steed,omitempty"`

	// Issue #215 payloads.
	PlanarAlly       *SpellPlanarAlly       `json:"planar_ally,omitempty"`
	PlanarBinding    *SpellPlanarBinding    `json:"planar_binding,omitempty"`
	PlaneShift       *SpellPlaneShift       `json:"plane_shift,omitempty"`
	PlantGrowth      *SpellPlantGrowth      `json:"plant_growth,omitempty"`
	PoisonSpray      *SpellPoisonSpray      `json:"poison_spray,omitempty"`
	Polymorph        *SpellPolymorph        `json:"polymorph,omitempty"`
	PowerWordHeal    *SpellPowerWordHeal    `json:"power_word_heal,omitempty"`
	PowerWordKill    *SpellPowerWordKill    `json:"power_word_kill,omitempty"`
	PowerWordStun    *SpellPowerWordStun    `json:"power_word_stun,omitempty"`
	PrayerOfHealing  *SpellPrayerOfHealing  `json:"prayer_of_healing,omitempty"`
	Prestidigitation *SpellPrestidigitation `json:"prestidigitation,omitempty"`
	PrismaticSpray   *SpellPrismaticSpray   `json:"prismatic_spray,omitempty"`

	// Issue #216 payloads.
	PrismaticWall             *SpellPrismaticWall             `json:"prismatic_wall,omitempty"`
	PrivateSanctum            *SpellPrivateSanctum            `json:"private_sanctum,omitempty"`
	ProduceFlame              *SpellProduceFlame              `json:"produce_flame,omitempty"`
	ProgrammedIllusion        *SpellProgrammedIllusion        `json:"programmed_illusion,omitempty"`
	ProjectImage              *SpellProjectImage              `json:"project_image,omitempty"`
	ProtectionFromEnergy      *SpellProtectionFromEnergy      `json:"protection_from_energy,omitempty"`
	ProtectionFromEvilAndGood *SpellProtectionFromEvilAndGood `json:"protection_from_evil_and_good,omitempty"`
	ProtectionFromPoison      *SpellProtectionFromPoison      `json:"protection_from_poison,omitempty"`
	PurifyFoodAndDrink        *SpellPurifyFoodAndDrink        `json:"purify_food_and_drink,omitempty"`
	RaiseDead                 *SpellRaiseDead                 `json:"raise_dead,omitempty"`
	RayOfEnfeeblement         *SpellRayOfEnfeeblement         `json:"ray_of_enfeeblement,omitempty"`
	RayOfFrost                *SpellRayOfFrost                `json:"ray_of_frost,omitempty"`

	// Issue #217 payloads.
	RayOfSickness   *SpellRayOfSickness   `json:"ray_of_sickness,omitempty"`
	Regenerate      *SpellRegenerate      `json:"regenerate,omitempty"`
	Reincarnate     *SpellReincarnate     `json:"reincarnate,omitempty"`
	RemoveCurse     *SpellRemoveCurse     `json:"remove_curse,omitempty"`
	ResilientSphere *SpellResilientSphere `json:"resilient_sphere,omitempty"`
	Resistance      *SpellResistance      `json:"resistance,omitempty"`
	Resurrection    *SpellResurrection    `json:"resurrection,omitempty"`
	ReverseGravity  *SpellReverseGravity  `json:"reverse_gravity,omitempty"`
	Revivify        *SpellRevivify        `json:"revivify,omitempty"`
	RopeTrick       *SpellRopeTrick       `json:"rope_trick,omitempty"`
	SacredFlame     *SpellSacredFlame     `json:"sacred_flame,omitempty"`
	Sanctuary       *SpellSanctuary       `json:"sanctuary,omitempty"`

	// Issue #219 payloads.
	Shillelagh       *SpellShillelagh       `json:"shillelagh,omitempty"`
	ShiningSmite     *SpellShiningSmite     `json:"shining_smite,omitempty"`
	ShockingGrasp    *SpellShockingGrasp    `json:"shocking_grasp,omitempty"`
	Silence          *SpellSilence          `json:"silence,omitempty"`
	SilentImage      *SpellSilentImage      `json:"silent_image,omitempty"`
	Simulacrum       *SpellSimulacrum       `json:"simulacrum,omitempty"`
	Sleep            *SpellSleep            `json:"sleep,omitempty"`
	SleetStorm       *SpellSleetStorm       `json:"sleet_storm,omitempty"`
	Slow             *SpellSlow             `json:"slow,omitempty"`
	SorcerousBurst   *SpellSorcerousBurst   `json:"sorcerous_burst,omitempty"`
	SpareTheDying    *SpellSpareTheDying    `json:"spare_the_dying,omitempty"`
	SpeakWithAnimals *SpellSpeakWithAnimals `json:"speak_with_animals,omitempty"`

	// Issue #220 payloads.
	SpeakWithDead    *SpellSpeakWithDead    `json:"speak_with_dead,omitempty"`
	SpeakWithPlants  *SpellSpeakWithPlants  `json:"speak_with_plants,omitempty"`
	SpiderClimb      *SpellSpiderClimb      `json:"spider_climb,omitempty"`
	SpikeGrowth      *SpellSpikeGrowth      `json:"spike_growth,omitempty"`
	SpiritGuardians  *SpellSpiritGuardians  `json:"spirit_guardians,omitempty"`
	SpiritualWeapon  *SpellSpiritualWeapon  `json:"spiritual_weapon,omitempty"`
	StarryWisp       *SpellStarryWisp       `json:"starry_wisp,omitempty"`
	StinkingCloud    *SpellStinkingCloud    `json:"stinking_cloud,omitempty"`
	StoneShape       *SpellStoneShape       `json:"stone_shape,omitempty"`
	Stoneskin        *SpellStoneskin        `json:"stoneskin,omitempty"`
	StormOfVengeance *SpellStormOfVengeance `json:"storm_of_vengeance,omitempty"`
	Suggestion       *SpellSuggestion       `json:"suggestion,omitempty"`

	// Issue #221 payloads.
	SummonDragon        *SpellSummonDragon        `json:"summon_dragon,omitempty"`
	Sunbeam             *SpellSunbeam             `json:"sunbeam,omitempty"`
	Sunburst            *SpellSunburst            `json:"sunburst,omitempty"`
	Symbol              *SpellSymbol              `json:"symbol,omitempty"`
	Telekinesis         *SpellTelekinesis         `json:"telekinesis,omitempty"`
	TelepathicBond      *SpellTelepathicBond      `json:"telepathic_bond,omitempty"`
	Teleport            *SpellTeleport            `json:"teleport,omitempty"`
	TeleportationCircle *SpellTeleportationCircle `json:"teleportation_circle,omitempty"`
	Thaumaturgy         *SpellThaumaturgy         `json:"thaumaturgy,omitempty"`
	Thunderwave         *SpellThunderwave         `json:"thunderwave,omitempty"`
	TimeStop            *SpellTimeStop            `json:"time_stop,omitempty"`
	TinyHut             *SpellTinyHut             `json:"tiny_hut,omitempty"`

	// Issue #222 payloads.
	Tongues            *SpellTongues            `json:"tongues,omitempty"`
	TransportViaPlants *SpellTransportViaPlants `json:"transport_via_plants,omitempty"`
	TreeStride         *SpellTreeStride         `json:"tree_stride,omitempty"`
	TruePolymorph      *SpellTruePolymorph      `json:"true_polymorph,omitempty"`
	TrueResurrection   *SpellTrueResurrection   `json:"true_resurrection,omitempty"`
	TrueSeeing         *SpellTrueSeeing         `json:"true_seeing,omitempty"`
	TrueStrike         *SpellTrueStrike         `json:"true_strike,omitempty"`
	Tsunami            *SpellTsunami            `json:"tsunami,omitempty"`
	UnseenServant      *SpellUnseenServant      `json:"unseen_servant,omitempty"`
	VampiricTouch      *SpellVampiricTouch      `json:"vampiric_touch,omitempty"`
	ViciousMockery     *SpellViciousMockery     `json:"vicious_mockery,omitempty"`
	VitriolicSphere    *SpellVitriolicSphere    `json:"vitriolic_sphere,omitempty"`

	// Issue #223 payloads.
	WallOfFire     *SpellWallOfFire     `json:"wall_of_fire,omitempty"`
	WallOfForce    *SpellWallOfForce    `json:"wall_of_force,omitempty"`
	WallOfIce      *SpellWallOfIce      `json:"wall_of_ice,omitempty"`
	WallOfStone    *SpellWallOfStone    `json:"wall_of_stone,omitempty"`
	WallOfThorns   *SpellWallOfThorns   `json:"wall_of_thorns,omitempty"`
	WardingBond    *SpellWardingBond    `json:"warding_bond,omitempty"`
	WaterBreathing *SpellWaterBreathing `json:"water_breathing,omitempty"`
	WaterWalk      *SpellWaterWalk      `json:"water_walk,omitempty"`
	Web            *SpellWeb            `json:"web,omitempty"`
	Weird          *SpellWeird          `json:"weird,omitempty"`
	WindWalk       *SpellWindWalk       `json:"wind_walk,omitempty"`
	WindWall       *SpellWindWall       `json:"wind_wall,omitempty"`

	// Issue #224 payloads.
	Wish         *SpellWish         `json:"wish,omitempty"`
	WordOfRecall *SpellWordOfRecall `json:"word_of_recall,omitempty"`
	ZoneOfTruth  *SpellZoneOfTruth  `json:"zone_of_truth,omitempty"`
}

// SpellRays describes independently aimed ranged spell attacks.
type SpellRays struct {
	AdditionalPerSlot int      `json:"additional_per_slot"`
	Count             int      `json:"count"`
	Distribution      []string `json:"distribution"`
	AttackPerRay      bool     `json:"attack_per_ray"`
}

// SpellSaveModifier is a source-listed knowledge or physical-connection choice.
type SpellSaveModifier struct {
	Choice   string `json:"choice"`
	Modifier int    `json:"modifier"`
}

// SpellScrying preserves the creature and stationary-location variants.
type SpellScrying struct {
	Variants               []string            `json:"variants"`
	CreatureSamePlane      bool                `json:"creature_same_plane"`
	Save                   string              `json:"save"`
	SaveAwareness          string              `json:"save_awareness"`
	Knowledge              []SpellSaveModifier `json:"knowledge"`
	Connections            []SpellSaveModifier `json:"connections"`
	SuccessRetryHours      int                 `json:"success_retry_hours"`
	SensorDistanceFeet     int                 `json:"sensor_distance_feet"`
	SensorInvisible        bool                `json:"sensor_invisible"`
	SensorIntangible       bool                `json:"sensor_intangible"`
	Senses                 []string            `json:"senses"`
	SensorFollowsCreature  bool                `json:"sensor_follows_creature"`
	VisibleSensor          string              `json:"visible_sensor"`
	LocationPreviouslySeen bool                `json:"location_previously_seen"`
	LocationSensorMoves    bool                `json:"location_sensor_moves"`
	FocusExamples          []string            `json:"focus_examples"`
}

// SpellSmite describes damage followed by a repeating end save.
type SpellSmite struct {
	Save             string `json:"save"`
	SaveTiming       string `json:"save_timing"`
	SuccessEnds      bool   `json:"success_ends"`
	FailureContinues bool   `json:"failure_continues"`
}

// SpellChest preserves both material objects and the recall/return lifecycle.
type SpellChest struct {
	DecayTiming            string   `json:"decay_timing"`
	ChestSizeFeet          []int    `json:"chest_size_feet"`
	ChestMinimumCostCP     int      `json:"chest_minimum_cost_cp"`
	ReplicaMinimumCostCP   int      `json:"replica_minimum_cost_cp"`
	ReplicaSize            string   `json:"replica_size"`
	Materials              string   `json:"materials"`
	BothTouched            bool     `json:"both_touched"`
	CapacityCubicFeet      int      `json:"capacity_cubic_feet"`
	Contents               string   `json:"contents"`
	Plane                  string   `json:"plane"`
	RecallAction           string   `json:"recall_action"`
	RecallTouch            string   `json:"recall_touch"`
	RecallDistanceFeet     int      `json:"recall_distance_feet"`
	RecallSpace            string   `json:"recall_space"`
	ReturnAction           string   `json:"return_action"`
	ReturnTouch            []string `json:"return_touch"`
	DecayAfterDays         int      `json:"decay_after_days"`
	DailyCumulativePercent int      `json:"daily_cumulative_percent"`
	EndTriggers            []string `json:"end_triggers"`
	EndsOnEthereal         string   `json:"ends_on_ethereal"`
}

// SpellVision preserves invisible and Ethereal perception without granting it.
type SpellVision struct {
	Subjects           []string `json:"subjects"`
	InvisibleAsVisible bool     `json:"invisible_as_visible"`
	SeeEthereal        bool     `json:"see_ethereal"`
	EtherealAppearance string   `json:"ethereal_appearance"`
}

// SpellSeeming preserves each target's appearance choices and discovery checks.
type SpellSeeming struct {
	Save                  string   `json:"save"`
	SaveEligibility       string   `json:"save_eligibility"`
	SaveSuccess           string   `json:"save_success"`
	AppearanceChoices     []string `json:"appearance_choices"`
	Affects               []string `json:"affects"`
	HeightDeltaFeet       int      `json:"height_delta_feet"`
	WeightChoices         []string `json:"weight_choices"`
	RetainLimbArrangement bool     `json:"retain_limb_arrangement"`
	OtherAppearance       string   `json:"other_appearance"`
	PhysicalInspection    string   `json:"physical_inspection"`
	DiscoveryAction       string   `json:"discovery_action"`
	DiscoveryCheck        string   `json:"discovery_check"`
	DiscoveryDC           string   `json:"discovery_dc"`
	DiscoverySuccess      string   `json:"discovery_success"`
}

// SpellSending records addressing alternatives, reply and recipient blocking.
type SpellSending struct {
	ReplyMaxWords            int      `json:"reply_max_words"`
	BlockTiming              string   `json:"block_timing"`
	AddressChoices           []string `json:"address_choices"`
	MaxWords                 int      `json:"max_words"`
	Delivery                 string   `json:"delivery"`
	RecognizeKnownSender     bool     `json:"recognize_known_sender"`
	UnderstandMeaning        bool     `json:"understand_meaning"`
	Reply                    string   `json:"reply"`
	CrossPlaneFailurePercent int      `json:"cross_plane_failure_percent"`
	SenderKnowsFailure       bool     `json:"sender_knows_failure"`
	RecipientMayBlock        bool     `json:"recipient_may_block"`
	BlockHours               int      `json:"block_hours"`
	BlockedAttempt           string   `json:"blocked_attempt"`
}

// SpellSequester preserves protections, creature suspension and early endings.
type SpellSequester struct {
	TargetChoices       []string `json:"target_choices"`
	Condition           string   `json:"condition"`
	Protections         []string `json:"protections"`
	CreatureCondition   string   `json:"creature_condition"`
	CreatureAges        bool     `json:"creature_ages"`
	UnneededSustenance  []string `json:"unneeded_sustenance"`
	EarlyEndCondition   string   `json:"early_end_condition"`
	EarlyEndRadiusMiles int      `json:"early_end_radius_miles"`
	EarlyEndExamples    []string `json:"early_end_examples"`
	DamageEnds          bool     `json:"damage_ends"`
}

// SpellShapechange preserves eligible forms, retained statistics and equipment choices.
type SpellShapechange struct {
	FormCRLimit                 string   `json:"form_cr_limit"`
	PreviouslySeen              bool     `json:"previously_seen"`
	ExcludedTypes               []string `json:"excluded_types"`
	ChangeAction                string   `json:"change_action"`
	TemporaryHP                 string   `json:"temporary_hp"`
	TemporaryHPGranted          string   `json:"temporary_hp_granted"`
	TemporaryHPEndsWithSpell    bool     `json:"temporary_hp_ends_with_spell"`
	ReplaceStatistics           string   `json:"replace_statistics"`
	Retained                    []string `json:"retained"`
	RetainSpellcastingIfPresent bool     `json:"retain_spellcasting_if_present"`
	EquipmentChoices            []string `json:"equipment_choices"`
}

// SpellShatter preserves area, saves and the separate object eligibility rule.
type SpellShatter struct {
	Center            string   `json:"center"`
	RadiusFeet        int      `json:"radius_feet"`
	Save              string   `json:"save"`
	SaveSuccess       string   `json:"save_success"`
	DisadvantageTypes []string `json:"disadvantage_types"`
	Objects           []string `json:"objects"`
	ObjectDamage      string   `json:"object_damage"`
	Sound             string   `json:"sound"`
}

// SpellShield describes the triggering attack and Magic Missile protection.
type SpellShield struct {
	IncludesTriggeringAttack bool     `json:"includes_triggering_attack"`
	NoDamageFrom             []string `json:"no_damage_from"`
	Appearance               string   `json:"appearance"`
}

// SpellFaith describes the full source target, separately from the self fixture.
type SpellFaith struct {
	TargetChoice string `json:"target_choice"`
	Appearance   string `json:"appearance"`
}

// UnmarshalJSON requires explicit fields within every selected typed payload.
// This prevents omitted false/zero source facts silently taking Go defaults.
func (s *SpellSourceDetails) UnmarshalJSON(data []byte) error {
	type plain SpellSourceDetails
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	var decoded plain
	if err := d.Decode(&decoded); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	for _, key := range []string{"casting_seconds", "casting_triggers", "end_timing"} {
		if len(fields[key]) == 0 || bytes.Equal(bytes.TrimSpace(fields[key]), []byte("null")) {
			return fmt.Errorf("source details require %s", key)
		}
	}
	templates := map[string]any{"rays": SpellRays{}, "scrying": SpellScrying{}, "smite": SpellSmite{}, "chest": SpellChest{}, "vision": SpellVision{}, "seeming": SpellSeeming{}, "sending": SpellSending{}, "sequester": SpellSequester{}, "shapechange": SpellShapechange{}, "shatter": SpellShatter{}, "shield": SpellShield{}, "faith": SpellFaith{}}
	for _, batch := range []map[string]any{spell196Templates(), spell197Templates(), spell198Templates(), spell199Templates(), spell200Templates(), spell201Templates(), spell202Templates(), spell208Templates(), spell207Templates(), spell206Templates(), spell205Templates(), spell204Templates(), spell203Templates(), spell209Templates(), spell210Templates(), spell211Templates(), spell212Templates(), spell213Templates(), spell214Templates(), spell215Templates(), spell216Templates(), spell217Templates(), spell219Templates(), spell220Templates(), spell221Templates(), spell222Templates(), spell223Templates(), spell224Templates()} {
		for key, template := range batch {
			if raw, ok := fields[key]; ok {
				if err := requireSpell198Fields(raw, reflect.TypeOf(template)); err != nil {
					return fmt.Errorf("%s: %w", key, err)
				}
			}
		}
	}
	for key, template := range templates {
		raw, ok := fields[key]
		if !ok {
			continue
		}
		b, err := json.Marshal(template)
		if err != nil {
			return err
		}
		var required, actual map[string]json.RawMessage
		if err := json.Unmarshal(b, &required); err != nil {
			return err
		}
		if err := json.Unmarshal(raw, &actual); err != nil {
			return err
		}
		for field := range required {
			value := bytes.TrimSpace(actual[field])
			if len(value) == 0 || bytes.Equal(value, []byte("null")) {
				return fmt.Errorf("%s requires %s", key, field)
			}
			var content any
			if err := json.Unmarshal(value, &content); err != nil {
				return err
			}
			switch content := content.(type) {
			case string:
				if strings.TrimSpace(content) == "" {
					return fmt.Errorf("%s has invalid %s", key, field)
				}
			case []any:
				if len(content) == 0 {
					return fmt.Errorf("%s has empty %s", key, field)
				}
				for _, item := range content {
					text, isText := item.(string)
					if item == nil || (isText && strings.TrimSpace(text) == "") {
						return fmt.Errorf("%s has invalid %s entry", key, field)
					}
				}
			case float64:
				if content < 0 {
					return fmt.Errorf("%s has invalid %s", key, field)
				}
			}
		}
	}
	if decoded.Scrying != nil {
		var raw map[string]json.RawMessage
		if err := json.Unmarshal(fields["scrying"], &raw); err != nil {
			return err
		}
		for _, group := range []string{"knowledge", "connections"} {
			var rows []map[string]json.RawMessage
			if err := json.Unmarshal(raw[group], &rows); err != nil {
				return err
			}
			for _, row := range rows {
				var choice string
				if err := json.Unmarshal(row["choice"], &choice); err != nil || strings.TrimSpace(choice) == "" || len(row["modifier"]) == 0 || bytes.Equal(bytes.TrimSpace(row["modifier"]), []byte("null")) {
					return fmt.Errorf("incomplete scrying choice")
				}
			}
		}
	}
	*s = SpellSourceDetails(decoded)
	return nil
}

func validSpellTarget(t SpellTarget) bool {
	if validSpell197Target(t) || validSpell200Target(t) || validSpell202Target(t) || validSpell203Target(t) || validSpell204Target(t) || validSpell205Target(t) || validSpell206Target(t) || validSpell207Target(t) || validSpell208Target(t) || validSpell209Target(t) || validSpell210Target(t) || validSpell211Target(t) || validSpell212Target(t) || validSpell213Target(t) || validSpell214Target(t) || validSpell215Target(t) || validSpell216Target(t) || validSpell217Target(t) || validSpell219Target(t) || validSpell220Target(t) || validSpell221Target(t) || validSpell222Target(t) || validSpell223Target(t) || validSpell224Target(t) {
		return true
	}
	switch t.Shape {
	case "square":
		return t.Count == 0 && t.Eligibility == "creatures-in-visible-ground-area"
	case "wall":
		return t.Count == 0 && t.Eligibility == "creatures-in-wall"
	case "cone":
		return t.Count == 0 && (t.Eligibility == "creatures-and-flammable-unattended-objects" || t.Eligibility == "creatures-in-cone")
	case "cube":
		return t.Count == 0 && t.Eligibility == "water-in-cube"
	case "cylinder":
		return t.Count == 0 && t.Eligibility == "creatures-within-5-feet-of-chosen-point-under-cloud"
	case "chain":
		return t.Count == 1 && t.Eligibility == "visible-creature-or-object-in-range"
	case "location":
		return t.Count == 1 && t.Eligibility == "familiar-or-obvious-location-in-range"
	case "up-to-count":
		return (t.Count == 3 && t.Eligibility == "creatures-in-range") || (t.Count == 1 && (t.Eligibility == "visible-creatures-in-range" || t.Eligibility == "visible-humanoids-in-range")) || (t.Count == 3 && t.Eligibility == "corpses-of-medium-or-small-humanoids-in-range")
	case "single":
		return t.Count == 1 && slices.Contains([]string{"target-in-range", "creature-in-range", "conscious-self", "self", "hit-target", "object-or-willing-creature", "chest-and-replica", "willing-creature-touch", "creature-touch", "visible-creature-in-range", "visible-beast-in-range", "visible-tiny-beast-in-range", "target-within-reach", "creature-or-flesh-touch", "object-touch"}, t.Eligibility)
	case "rays":
		return t.Count == 3 && t.Eligibility == "targets-in-range"
	case "choice":
		return t.Count == 1 && slices.Contains([]string{"same-plane-creature-or-seen-location", "visible-creature-or-nonmagical-noncreature-plant", "door-window-or-area-in-range", "bones-or-corpse-of-medium-or-small-humanoid-in-range", "ground-or-containers-in-range", "container-or-cube-in-range"}, t.Eligibility)
	case "selected-creatures":
		return t.Count == 0 && slices.Contains([]string{"visible-creatures-in-range", "creatures-in-range", "visible-willing-creatures-in-range"}, t.Eligibility)
	case "sphere":
		return t.Count == 0 && (t.Eligibility == "creatures-and-eligible-objects" || t.Eligibility == "creatures-in-sphere" || t.Eligibility == "humanoids-in-sphere")
	case "selected-objects":
		return t.Count == 0 && t.Eligibility == "eligible-nonmagical-objects-in-range"
	case "emanation":
		return t.Count == 0 && (t.Eligibility == "creatures-other-than-constructs-and-undead" || t.Eligibility == "everything-in-aura" || t.Eligibility == "visible-creatures-in-emanation")
	default:
		return false
	}
}

func (s SpellRecord) validateSourceDetails() error {
	// The bounded Light fixture carries no source payload; its pinned values
	// live in the complete batch 211 record.
	if s.Scope == lightFixtureScope {
		if s.SourceDetails != nil {
			return fmt.Errorf("bounded Light fixture cannot carry source details")
		}
		return nil
	}
	if handled, err := s.validateSpell196(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell197(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell198(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell199(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell200(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell201(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell202(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell203(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell204(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell205(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell206(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell207(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell208(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell209(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell210(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell211(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell212(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell213(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell214(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell215(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell216(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell217(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell219(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell220(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell221(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell222(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell223(); handled || err != nil {
		return err
	}
	if handled, err := s.validateSpell224(); handled || err != nil {
		return err
	}
	d := s.SourceDetails
	if d == nil {
		if s.Scope == "complete-source-record" && slices.Contains([]string{"Scorching Ray", "Scrying", "Searing Smite", "Secret Chest", "See Invisibility", "Seeming", "Sending", "Sequester", "Shapechange", "Shatter", "Shield", "Shield of Faith"}, s.Name) {
			return fmt.Errorf("batch source details required")
		}
		if s.CastingTime == "reaction" || s.CastingTime == "10-minutes" || s.CastingTime == "1-minute" || s.CastingTime == "12-hours" || spell197CastingSeconds[s.CastingTime] != 0 || s.Duration == "until-dispelled" || s.Duration == "special" || s.Range.Kind == "unlimited" || s.Range.Kind == "special" || s.Range.Kind == "sight" || s.Target.Shape != "single" {
			return fmt.Errorf("source details required for extended shape")
		}
		return nil
	}
	if s.Scope != "complete-source-record" || d.EndTiming == "" || d.CastingSeconds < 0 || d.CastingTriggers == nil {
		return fmt.Errorf("invalid source details")
	}
	if (s.CastingTime == "10-minutes") != (d.CastingSeconds == 600) || (s.CastingTime == "reaction" && len(d.CastingTriggers) == 0) {
		return fmt.Errorf("invalid casting details")
	}
	payloads := []bool{d.Rays != nil, d.Scrying != nil, d.Smite != nil, d.Chest != nil, d.Vision != nil, d.Seeming != nil, d.Sending != nil, d.Sequester != nil, d.Shapechange != nil, d.Shatter != nil, d.Shield != nil, d.Faith != nil}
	names := []string{"scorching-ray", "scrying", "searing-smite", "secret-chest", "see-invisibility", "seeming", "sending", "sequester", "shapechange", "shatter", "shield", "shield-of-faith"}
	count := 0
	for i, present := range payloads {
		if present {
			count++
			if s.ID != prefix+"spells/"+names[i] {
				return fmt.Errorf("mismatched source payload")
			}
		}
	}
	if count != 1 {
		return fmt.Errorf("exactly one batch payload required")
	}
	if d.Scrying != nil {
		if len(d.Scrying.Variants) != 2 || len(d.Scrying.Knowledge) != 3 || len(d.Scrying.Connections) != 3 || d.Scrying.Save != "wisdom" || d.Scrying.SensorDistanceFeet <= 0 || d.Scrying.SuccessRetryHours <= 0 {
			return fmt.Errorf("invalid scrying choices")
		}
	}
	if d.Chest != nil && (len(d.Chest.ChestSizeFeet) != 3 || d.Chest.ChestMinimumCostCP <= 0 || d.Chest.ReplicaMinimumCostCP <= 0 || d.Chest.DailyCumulativePercent <= 0 || d.Chest.DailyCumulativePercent > 100) {
		return fmt.Errorf("invalid chest dimensions, costs or decay")
	}
	if d.Shapechange != nil && (len(d.Shapechange.ExcludedTypes) != 2 || len(d.Shapechange.Retained) == 0 || len(d.Shapechange.EquipmentChoices) != 2) {
		return fmt.Errorf("invalid form choices")
	}
	if d.Shatter != nil && (d.Shatter.RadiusFeet <= 0 || d.Shatter.Save != "constitution") {
		return fmt.Errorf("invalid shatter area or save")
	}
	if d.Sending != nil && (d.Sending.MaxWords <= 0 || d.Sending.CrossPlaneFailurePercent < 0 || d.Sending.CrossPlaneFailurePercent > 100 || d.Sending.BlockHours <= 0) {
		return fmt.Errorf("invalid sending limits")
	}
	if d.Rays != nil && (d.Rays.Count != 3 || d.Rays.AdditionalPerSlot <= 0 || !slices.Equal(d.Rays.Distribution, []string{"one-target", "several-targets"}) || !d.Rays.AttackPerRay) {
		return fmt.Errorf("invalid ray choices")
	}
	if d.Scrying != nil && (!slices.Equal(d.Scrying.Variants, []string{"creature", "seen-location"}) || !d.Scrying.CreatureSamePlane || !d.Scrying.LocationPreviouslySeen || d.Scrying.LocationSensorMoves) {
		return fmt.Errorf("invalid scrying variants")
	}
	if d.Smite != nil && (d.Smite.Save != "constitution" || d.Smite.SaveTiming != "after-start-of-turn-damage" || !d.Smite.SuccessEnds || !d.Smite.FailureContinues || len(d.CastingTriggers) != 2) {
		return fmt.Errorf("invalid smite save or trigger")
	}
	if d.Seeming != nil && (d.Seeming.HeightDeltaFeet < 0 || d.Seeming.Save != "charisma" || len(d.Seeming.AppearanceChoices) != 2 || len(d.Seeming.WeightChoices) != 2) {
		return fmt.Errorf("invalid seeming choices")
	}
	if d.Sequester != nil && (d.Sequester.EarlyEndRadiusMiles <= 0 || !slices.Equal(d.Sequester.TargetChoices, []string{"object", "willing-creature"}) || d.Sequester.CreatureAges || !slices.Equal(d.Sequester.UnneededSustenance, []string{"food", "water", "air"})) {
		return fmt.Errorf("invalid sequester choices")
	}
	if d.Shield != nil && (!slices.Equal(d.CastingTriggers, []string{"hit-by-attack-roll", "targeted-by-magic-missile"}) || d.EndTiming != "start-of-caster-next-turn") {
		return fmt.Errorf("invalid shield reaction")
	}

	return nil
}
