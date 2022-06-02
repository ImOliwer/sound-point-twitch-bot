export const TitleBase = "Sound Point Twitch Bot |";

export type Deployed = {
  price: number;
  file_name: string;
  cooldown: number;
  last_used: number;
};

export type SoundMap = {
  [key: string]: Deployed;
};

const fixed_regex = /^-?\d+(?:\.\d{0,2})?/;

export function fixed(value: number): string {
  const check = value.toString().match(fixed_regex);
  if (check == null) {
    return ""
  }
  return check[0];
}

export function formatNumber(cooldown: number): string {
  const milliseconds = (cooldown % 1000) / 100,
        seconds      = (cooldown / 1000) % 60,
        minutes      = (cooldown / (1000 * 60)) % 60,
        hours        = ((cooldown / (1000 * 60 * 60)) % 24),
        days         = (cooldown) / (1000 * 60 * 60 * 24);
  
  if (days >= 1) {
    return `${fixed(days)}d`;
  }

  if (hours >= 1) {
    return `${fixed(hours)}h`;
  }

  if (minutes >= 1) {
    return `${fixed(minutes)}m`;
  }

  if (seconds >= 1) {
    return `${fixed(seconds)}s`;
  }
  
  return `${fixed(milliseconds)}ms`;
}

export function pageOf(map: SoundMap, page: number, max: number): string[] {
  page = Math.floor(page);
  if (page <= 0) {
    return [];
  }
  
  const keys = Object.keys(map);
  const length = keys.length;

  if (length == 0) {
    return [];
  }

  max = Math.floor(max);
  const expectedIndex = page*max - 4;

  if (length < expectedIndex) {
    return []; // pageOf(map, page-1, max); return the last known (if any)
  }

  const population: string[] = [];
  let index = 0;

  for (const key of keys) {
    index++;

    if (index < expectedIndex) {
      continue;
    }

    if (index > (expectedIndex+max)-1) {
      break;
    }

    population.push(key)
  }
  return population;
}

export function notEmptyOrElse(value: string[], orElse: () => string[]) {
  return value.length > 0 ? value : orElse();
}