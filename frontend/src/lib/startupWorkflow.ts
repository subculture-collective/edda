import type { CampaignProfile, CharacterProfile, CharacterSpawnPackage, OpeningSceneResponse, StarterItem } from '../api/types';

export interface GuidedCharacterAttributes {
  readonly name: string;
  readonly race: string;
  readonly class: string;
  readonly background: string;
  readonly alignment: string;
  readonly traits?: readonly string[];
}

export interface StartupPlaySeed {
  readonly campaignName: string;
  readonly campaignSummary: string;
  readonly openingScene: OpeningSceneResponse;
  readonly seededAt: string;
}

const CLASS_MOTIVATIONS: Record<string, readonly string[]> = {
  Barbarian: ['Prove strength in battle', 'Protect the tribe'],
  Bard: ['Collect epic tales', 'Inspire others through art'],
  Cleric: ['Serve the divine', 'Heal the wounded'],
  Druid: ['Preserve the natural order', 'Guard the wilds'],
  Fighter: ['Protect the weak', 'Master the art of combat'],
  Monk: ['Achieve inner perfection', 'Uphold monastic traditions'],
  Paladin: ['Uphold a sacred oath', 'Smite the wicked'],
  Ranger: ['Guard the frontier', 'Hunt dangerous beasts'],
  Rogue: ['Seek fortune and thrills', 'Uncover hidden truths'],
  Sorcerer: ['Understand innate power', 'Prove worthy of the gift'],
  Warlock: ['Fulfill the pact', 'Gain forbidden knowledge'],
  Wizard: ['Seek arcane knowledge', 'Unlock the mysteries of magic'],
};

const RACE_STRENGTHS: Record<string, readonly string[]> = {
  Human: ['Versatile adaptability'],
  Elf: ['Keen perception', 'Grace under pressure'],
  Dwarf: ['Unyielding resilience', 'Stonework intuition'],
  Halfling: ['Remarkable luck', 'Nimble evasion'],
  Gnome: ['Inventive ingenuity', 'Arcane curiosity'],
  'Half-Elf': ['Social adaptability', 'Dual heritage insight'],
  'Half-Orc': ['Savage endurance', 'Intimidating presence'],
  Tiefling: ['Infernal resilience', 'Force of personality'],
  Dragonborn: ['Draconic breath weapon', 'Commanding presence'],
};

const CLASS_STRENGTHS: Record<string, readonly string[]> = {
  Barbarian: ['Raw physical power'],
  Bard: ['Silver tongue'],
  Cleric: ['Divine channeling'],
  Druid: ['Nature magic'],
  Fighter: ['Combat expertise'],
  Monk: ['Unarmed discipline'],
  Paladin: ['Righteous resolve'],
  Ranger: ['Wilderness tracking'],
  Rogue: ['Stealth and cunning'],
  Sorcerer: ['Innate spellcasting'],
  Warlock: ['Eldritch invocations'],
  Wizard: ['Scholarly spellcraft'],
};

const BACKGROUND_WEAKNESSES: Record<string, readonly string[]> = {
  Acolyte: ['Naïve idealism', 'Blind faith'],
  Charlatan: ['Compulsive dishonesty', 'Difficulty forming trust'],
  Criminal: ['Trust issues', 'Haunted by past deeds'],
  Entertainer: ['Craves attention', 'Fear of being forgotten'],
  'Folk Hero': ['Stubborn pride', 'Unrealistic expectations'],
  'Guild Artisan': ['Obsessive perfectionism', 'Material attachment'],
  Hermit: ['Social awkwardness', 'Distrust of authority'],
  Noble: ['Sense of entitlement', 'Sheltered worldview'],
  Outlander: ['Distrust of civilization', 'Blunt to a fault'],
  Sage: ['Absent-minded', 'Overthinks simple problems'],
  Sailor: ['Restless on land', 'Rough manners'],
  Soldier: ['Rigid discipline', 'Haunted by war'],
  Urchin: ['Deep-seated insecurity', 'Hoarding instinct'],
};

const DEFAULT_MOTIVATIONS = ['Seek adventure', 'Find a place in the world'] as const;
const DEFAULT_STRENGTHS = ['Determined spirit'] as const;
const DEFAULT_WEAKNESSES = ['Untested resolve'] as const;

const CLASS_STARTER_ITEMS: Record<string, StarterItem> = {
  barbarian: {
    name: 'Notched Handaxe',
    description: 'A hard-used handaxe balanced for rough travel and sudden violence.',
    item_type: 'weapon',
    equipped: true,
  },
  bard: {
    name: 'Weathered Lute',
    description: 'A travel-scarred instrument with names and tavern marks scratched inside the case.',
    item_type: 'tool',
  },
  cleric: {
    name: 'Traveling Shrine Kit',
    description: 'A compact bundle of candles, chalk, and a worn holy symbol for roadside rites.',
    item_type: 'tool',
  },
  druid: {
    name: 'Herbalist Wrap',
    description: 'Pressed leaves, twine, and field salves gathered from familiar wild places.',
    item_type: 'tool',
  },
  fighter: {
    name: 'Well-Kept Blade',
    description: 'A practical weapon kept sharp by habit rather than ceremony.',
    item_type: 'weapon',
    equipped: true,
  },
  monk: {
    name: 'Prayer Bead Cord',
    description: 'A cord of smooth beads used to focus breath, discipline, and memory.',
    item_type: 'trinket',
  },
  paladin: {
    name: 'Oathbound Emblem',
    description: 'A polished emblem that marks the bearer as someone sworn to a difficult promise.',
    item_type: 'trinket',
  },
  ranger: {
    name: 'Waxed Route Map',
    description: 'A folded map treated against rain, annotated with trails, hazards, and old campsite marks.',
    item_type: 'tool',
  },
  rogue: {
    name: 'Concealed Lockpick Roll',
    description: 'A slim leather roll of picks and shims hidden inside a mending kit.',
    item_type: 'tool',
  },
  sorcerer: {
    name: 'Cracked Focus Stone',
    description: 'A small stone that warms when the bearer’s innate magic stirs.',
    item_type: 'arcane_focus',
  },
  warlock: {
    name: 'Sealed Pact Token',
    description: 'A token whose seal should probably not be broken without need.',
    item_type: 'trinket',
  },
  wizard: {
    name: 'Annotated Spell Primer',
    description: 'A battered primer crowded with marginal notes, warnings, and half-solved diagrams.',
    item_type: 'book',
  },
};

const BACKGROUND_STARTER_ITEMS: Record<string, StarterItem> = {
  acolyte: {
    name: 'Pilgrim Tokens',
    description: 'Small tokens from temples, shrines, or teachers that still open a few doors.',
    item_type: 'trinket',
  },
  charlatan: {
    name: 'False-Seal Papers',
    description: 'Convincing blank papers, enough for one risky lie if used carefully.',
    item_type: 'document',
  },
  criminal: {
    name: 'Underworld Contact Mark',
    description: 'A discreet sign that may be recognized by people who prefer not to use names.',
    item_type: 'trinket',
  },
  entertainer: {
    name: 'Performance Handbill',
    description: 'A faded notice from a prior show, useful proof that someone has seen the road.',
    item_type: 'document',
  },
  'folk hero': {
    name: 'Village Favor Token',
    description: 'A handmade token from people who believe the bearer once did something worth remembering.',
    item_type: 'trinket',
  },
  'guild artisan': {
    name: 'Guild Letter of Passage',
    description: 'A stamped letter asking honest workshops to offer professional courtesy.',
    item_type: 'document',
  },
  hermit: {
    name: 'Cryptic Field Notes',
    description: 'Private observations that make sense only after the world proves them true.',
    item_type: 'document',
  },
  noble: {
    name: 'Signet of Standing',
    description: 'A minor signet that still carries weight with people impressed by bloodlines.',
    item_type: 'trinket',
  },
  outlander: {
    name: 'Trail Ration Bundle',
    description: 'Hard rations, firestarter, and a weatherproof wrap packed by someone who knows hunger.',
    item_type: 'supply',
  },
  sage: {
    name: 'Indexed Research Notes',
    description: 'A tidy stack of observations cross-referenced with questions not yet answered.',
    item_type: 'document',
  },
  sailor: {
    name: 'Salt-Stained Compass',
    description: 'A compass with unreliable polish and reliable instincts for direction.',
    item_type: 'tool',
  },
  soldier: {
    name: 'Campaign Medal',
    description: 'A battered medal from an old campaign, meaningful to veterans and dangerous around enemies.',
    item_type: 'trinket',
  },
  urchin: {
    name: 'Hidden Coin Cache',
    description: 'A tiny emergency stash sewn into a hem for the day everything goes wrong.',
    item_type: 'supply',
  },
};

export function summarizeCampaignProfile(profile: CampaignProfile): string {
  const parts: string[] = [];

  if (profile.tone !== '' || profile.genre !== '') {
    parts.push(`${profile.tone} ${profile.genre}`.trim());
  }
  if (profile.themes.length > 0) {
    parts.push(`themes of ${profile.themes.join(', ')}`);
  }
  if (profile.world_type !== '') {
    parts.push(`set in a ${profile.world_type} world`);
  }
  if (profile.danger_level !== '') {
    parts.push(`with ${profile.danger_level} danger`);
  }
  if (profile.political_complexity !== '') {
    parts.push(`and ${profile.political_complexity} politics`);
  }

  if (parts.length === 0) {
    return 'A new adventure awaits.';
  }

  return `${parts.join(' ').trim()}.`;
}

export function buildCharacterProfileFromGuidedAttributes(attributes: GuidedCharacterAttributes): CharacterProfile {
  const concept = `${attributes.race.toLowerCase()} ${attributes.class.toLowerCase()}`.trim();
  const traits = (attributes.traits ?? []).filter((trait) => trait.trim().length > 0);

  let personality = attributes.alignment;
  if (traits.length > 0) {
    personality = `${personality}; ${traits.join('; ')}`;
  }

  return {
    name: attributes.name.trim(),
    concept,
    background: attributes.background,
    personality,
    motivations: lookupStrings(CLASS_MOTIVATIONS, attributes.class, DEFAULT_MOTIVATIONS),
    strengths: mergeStrengths(attributes.race, attributes.class),
    weaknesses: lookupStrings(BACKGROUND_WEAKNESSES, attributes.background, DEFAULT_WEAKNESSES),
  };
}

export function buildCharacterSpawnPackage(
  characterProfile: CharacterProfile,
  campaignProfile: CampaignProfile,
  campaignSummary: string,
): CharacterSpawnPackage {
  const items = uniqueItems([
    starterItemForClass(characterProfile.concept),
    starterItemForBackground(characterProfile.background),
  ]);

  return {
    items,
    known_facts: [
      {
        fact: `${characterProfile.name} begins with personal context from a ${characterProfile.background.toLowerCase()} background: ${characterProfile.motivations[0] ?? 'seek adventure'}.`,
        category: 'character',
      },
      {
        fact: `This campaign is framed as ${summarizeCampaignProfile(campaignProfile).toLowerCase()} ${campaignSummary}`.trim(),
        category: 'world',
      },
    ],
  };
}

function lookupStrings(
  source: Record<string, readonly string[]>,
  key: string,
  fallback: readonly string[],
): string[] {
  return [...(source[key] ?? fallback)];
}

function mergeStrengths(race: string, characterClass: string): string[] {
  return [
    ...lookupStrings(RACE_STRENGTHS, race, DEFAULT_STRENGTHS),
    ...lookupStrings(CLASS_STRENGTHS, characterClass, DEFAULT_STRENGTHS),
  ];
}

function starterItemForClass(concept: string): StarterItem {
  const normalizedConcept = concept.toLowerCase();
  const classKey = Object.keys(CLASS_STARTER_ITEMS).find((candidate) => normalizedConcept.includes(candidate));
  return withDefaults(classKey ? CLASS_STARTER_ITEMS[classKey] : {
    name: 'Traveler’s Kit',
    description: 'A practical bundle of road-worn gear for a newly begun journey.',
    item_type: 'supply',
  });
}

function starterItemForBackground(background: string): StarterItem {
  const backgroundKey = background.toLowerCase();
  return withDefaults(BACKGROUND_STARTER_ITEMS[backgroundKey] ?? {
    name: 'Personal Keepsake',
    description: 'A small object tied to the character’s life before the adventure began.',
    item_type: 'trinket',
  });
}

function withDefaults(item: StarterItem): StarterItem {
  return {
    rarity: 'common',
    quantity: 1,
    properties: { source: 'startup_character_profile' },
    ...item,
  };
}

function uniqueItems(items: StarterItem[]): StarterItem[] {
  const seen = new Set<string>();
  return items.filter((item) => {
    const key = item.name.toLowerCase();
    if (seen.has(key)) {
      return false;
    }
    seen.add(key);
    return true;
  });
}
