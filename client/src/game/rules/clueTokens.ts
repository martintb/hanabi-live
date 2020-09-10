// Functions related to clues: gaining clues, giving clues, applying clues

import { clueTokensRules } from '../rules';
import { ActionPlay, ActionDiscard } from '../types/actions';
import { MAX_CLUE_NUM } from '../types/constants';
import Variant from '../types/Variant';
import * as variantRules from './variant';

// Gain a clue by discarding or finishing a stack
export const gain = (
  action: ActionPlay | ActionDiscard,
  clueTokens: number,
  variant: Variant,
  playStackComplete: boolean = false,
) => {
  if (shouldGenerateClue(action, clueTokens, variant, playStackComplete)) {
    return clueTokens + 1;
  }
  return clueTokens;
};

const shouldGenerateClue = (
  action: ActionPlay | ActionDiscard,
  clueTokens: number,
  variant: Variant,
  playStackComplete: boolean,
) => {
  if (clueTokensRules.atMax(clueTokens, variant)) {
    return false;
  }

  switch (action.type) {
    case 'play': {
      // Finishing a play stack grants an extra clue
      // (but not in certain variants)
      return playStackComplete && !variantRules.isThrowItInAHole(variant);
    }

    case 'discard': {
      // Discarding a card grants an extra clue
      // But misplayed cards do not grant extra clues
      return !action.failed;
    }

    default: {
      return false;
    }
  }
};

export const getAdjusted = (clueTokens: number, variant: Variant) => {
  // In "Clue Starved" variants, each discard only grants 0.5 clue tokens
  // This is represented on the client by discards granting 1 clue token and clues costing 2 tokens
  // (to avoid having to use floating point numbers)
  if (variantRules.isClueStarved(variant)) {
    return clueTokens * 2;
  }

  return clueTokens;
};

export const atMax = (
  clueTokens: number,
  variant: Variant,
) => clueTokens >= getAdjusted(MAX_CLUE_NUM, variant);

// The value of clues gained when discarding or finishing a suit
// This function is *only* used in efficiency calculations
export const value = (variant: Variant) => {
  // In "Clue Starved" variants, each discard gives only half a clue
  if (variantRules.isClueStarved(variant)) {
    return 0.5;
  }

  return 1;
};
