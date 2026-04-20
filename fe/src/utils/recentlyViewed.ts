const RECENTLY_VIEWED_KEY = 'recently_viewed_devices';
const MAX_RECENTLY_VIEWED = 8;

const normalizeIds = (value: unknown): string[] => {
  if (!Array.isArray(value)) return [];

  const unique = new Set<string>();
  const result: string[] = [];

  value.forEach((item) => {
    if (typeof item !== 'string') return;
    const trimmed = item.trim();
    if (!trimmed || unique.has(trimmed)) return;
    unique.add(trimmed);
    result.push(trimmed);
  });

  return result;
};

export const getRecentlyViewedIds = (): string[] => {
  try {
    const raw = localStorage.getItem(RECENTLY_VIEWED_KEY);
    if (!raw) return [];

    return normalizeIds(JSON.parse(raw));
  } catch {
    return [];
  }
};

export const pushRecentlyViewedId = (deviceId?: string) => {
  const id = deviceId?.trim();
  if (!id) return;

  const current = getRecentlyViewedIds().filter((item) => item !== id);
  const next = [id, ...current].slice(0, MAX_RECENTLY_VIEWED);

  localStorage.setItem(RECENTLY_VIEWED_KEY, JSON.stringify(next));
};

