import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react';
import { favoritesApi } from '../api/favorites';
import { useAuth } from './AuthContext';

interface FavoritesContextType {
  favoriteIds: string[];
  isLoading: boolean;
  isFavorite: (deviceId?: string) => boolean;
  toggleFavorite: (deviceId?: string) => Promise<void>;
  removeFavorite: (deviceId?: string) => Promise<void>;
}

const GUEST_FAVORITES_KEY = 'favorites';

const FavoritesContext = createContext<FavoritesContextType | undefined>(undefined);

const normalizeGuestFavorites = (raw: unknown): string[] => {
  if (!Array.isArray(raw)) {
    return [];
  }

  const result: string[] = [];
  const seen = new Set<string>();

  raw.forEach((value) => {
    if (typeof value !== 'string') return;
    const trimmed = value.trim();
    if (!trimmed || seen.has(trimmed)) return;
    seen.add(trimmed);
    result.push(trimmed);
  });

  return result;
};

const getGuestFavorites = (): string[] => {
  try {
    const stored = localStorage.getItem(GUEST_FAVORITES_KEY);
    if (!stored) return [];
    return normalizeGuestFavorites(JSON.parse(stored));
  } catch {
    return [];
  }
};

const saveGuestFavorites = (deviceIds: string[]) => {
  localStorage.setItem(GUEST_FAVORITES_KEY, JSON.stringify(deviceIds));
};

export function FavoritesProvider({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth();
  const [favoriteIds, setFavoriteIds] = useState<string[]>(() => getGuestFavorites());
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    let active = true;

    const bootstrapFavorites = async () => {
      if (!isAuthenticated) {
        if (active) {
          setFavoriteIds(getGuestFavorites());
          setIsLoading(false);
        }
        return;
      }

      setIsLoading(true);
      try {
        const guestFavorites = getGuestFavorites();

        if (guestFavorites.length > 0) {
          const { data } = await favoritesApi.sync({ device_ids: guestFavorites });
          if (!active) return;
          setFavoriteIds(data.data?.device_ids || []);
          localStorage.removeItem(GUEST_FAVORITES_KEY);
        } else {
          const { data } = await favoritesApi.getAll();
          if (!active) return;
          setFavoriteIds(data.data?.device_ids || []);
        }
      } catch {
        if (active) {
          setFavoriteIds([]);
        }
      } finally {
        if (active) {
          setIsLoading(false);
        }
      }
    };

    bootstrapFavorites();

    return () => {
      active = false;
    };
  }, [isAuthenticated]);

  const isFavorite = useCallback(
    (deviceId?: string) => {
      if (!deviceId) return false;
      return favoriteIds.includes(deviceId);
    },
    [favoriteIds],
  );

  const removeFavorite = useCallback(
    async (deviceId?: string) => {
      const normalized = deviceId?.trim();
      if (!normalized) return;

      if (isAuthenticated) {
        const { data } = await favoritesApi.remove(normalized);
        setFavoriteIds(data.data?.device_ids || []);
        return;
      }

      setFavoriteIds((prev) => {
        const next = prev.filter((id) => id !== normalized);
        saveGuestFavorites(next);
        return next;
      });
    },
    [isAuthenticated],
  );

  const toggleFavorite = useCallback(
    async (deviceId?: string) => {
      const normalized = deviceId?.trim();
      if (!normalized) return;

      if (isAuthenticated) {
        if (favoriteIds.includes(normalized)) {
          const { data } = await favoritesApi.remove(normalized);
          setFavoriteIds(data.data?.device_ids || []);
        } else {
          const { data } = await favoritesApi.add({ device_id: normalized });
          setFavoriteIds(data.data?.device_ids || []);
        }
        return;
      }

      setFavoriteIds((prev) => {
        const exists = prev.includes(normalized);
        const next = exists ? prev.filter((id) => id !== normalized) : [normalized, ...prev];
        saveGuestFavorites(next);
        return next;
      });
    },
    [favoriteIds, isAuthenticated],
  );

  const value = useMemo(
    () => ({
      favoriteIds,
      isLoading,
      isFavorite,
      toggleFavorite,
      removeFavorite,
    }),
    [favoriteIds, isLoading, isFavorite, toggleFavorite, removeFavorite],
  );

  return <FavoritesContext.Provider value={value}>{children}</FavoritesContext.Provider>;
}

export function useFavorites() {
  const context = useContext(FavoritesContext);
  if (!context) {
    throw new Error('useFavorites must be used within a FavoritesProvider');
  }

  return context;
}

