import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { Heart, Smartphone, Trash2 } from 'lucide-react';
import { devicesApi } from '../api/devices';
import { useFavorites } from '../contexts/FavoritesContext';
import type { Device } from '../types';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import { resolveDeviceImageUrl } from '../utils/resolveDeviceImageUrl';

export default function FavoritesPage() {
  const { favoriteIds, isLoading: loadingFavorites, removeFavorite } = useFavorites();
  const [devices, setDevices] = useState<Device[]>([]);
  const [loadingDevices, setLoadingDevices] = useState(false);

  useEffect(() => {
    let active = true;

    const fetchFavorites = async () => {
      if (favoriteIds.length === 0) {
        setDevices([]);
        return;
      }

      setLoadingDevices(true);
      try {
        const results = await Promise.allSettled(
          favoriteIds.map((deviceId) => devicesApi.getById(deviceId)),
        );

        if (!active) return;

        const fetchedDevices: Device[] = [];
        results.forEach((result) => {
          if (result.status === 'fulfilled' && result.value.data.data) {
            fetchedDevices.push(result.value.data.data);
          }
        });

        setDevices(fetchedDevices);
      } finally {
        if (active) {
          setLoadingDevices(false);
        }
      }
    };

    fetchFavorites();

    return () => {
      active = false;
    };
  }, [favoriteIds]);

  const deviceMap = useMemo(() => {
    const map = new Map<string, Device>();
    devices.forEach((device) => {
      if (device.id) {
        map.set(device.id, device);
      }
    });
    return map;
  }, [devices]);

  if (loadingFavorites || loadingDevices) {
    return <LoadingSpinner text="Loading favorites..." />;
  }

  if (favoriteIds.length === 0) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-16 text-center glass rounded-2xl mt-10">
        <Heart size={44} className="mx-auto text-text-muted mb-4" />
        <h1 className="text-2xl font-bold text-text-primary mb-2">Your favorites list is empty</h1>
        <p className="text-text-secondary">Tap the heart icon on any device to add it here.</p>
        <Link
          to="/brands"
          className="inline-flex items-center gap-2 mt-6 px-5 py-2.5 rounded-xl text-sm font-semibold text-white btn-gradient"
        >
          Browse devices
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div className="flex items-center gap-3 mb-6">
        <Heart size={20} className="text-danger" />
        <h1 className="text-2xl sm:text-3xl font-bold text-text-primary">Favorite Devices</h1>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-5">
        {favoriteIds.map((deviceId) => {
          const device = deviceMap.get(deviceId);

          if (!device) {
            return (
              <div key={deviceId} className="glass rounded-2xl p-4 flex flex-col gap-3">
                <div className="h-40 rounded-xl bg-surface-light flex items-center justify-center">
                  <Smartphone size={36} className="text-text-muted" />
                </div>
                <p className="text-sm text-text-secondary">This device is no longer available.</p>
                <button
                  type="button"
                  onClick={() => removeFavorite(deviceId)}
                  className="inline-flex items-center justify-center gap-2 px-3 py-2 rounded-lg text-xs font-semibold text-danger hover:bg-surface-light"
                >
                  <Trash2 size={14} />
                  Remove
                </button>
              </div>
            );
          }

          return (
            <div key={deviceId} className="glass rounded-2xl overflow-hidden card-hover">
              <Link to={`/devices/${deviceId}`} className="block">
                <div className="aspect-square bg-gradient-to-br from-surface-lighter/50 to-surface-light flex items-center justify-center p-6 overflow-hidden">
                  {device.imageUrl ? (
                    <img
                      src={resolveDeviceImageUrl(device.imageUrl)}
                      alt={device.model_name}
                      className="max-h-full w-auto object-contain"
                    />
                  ) : (
                    <Smartphone size={48} className="text-text-muted" />
                  )}
                </div>
                <div className="p-4">
                  <h2 className="text-sm font-semibold text-text-primary truncate">{device.model_name}</h2>
                  <p className="text-xs text-text-muted mt-1 truncate">
                    {device.specifications?.platform?.chipset?.split('(')[0].trim() || 'Device'}
                  </p>
                </div>
              </Link>

              <div className="px-4 pb-4">
                <button
                  type="button"
                  onClick={() => removeFavorite(deviceId)}
                  className="w-full inline-flex items-center justify-center gap-2 px-3 py-2 rounded-lg text-xs font-semibold text-danger hover:bg-surface-light"
                >
                  <Trash2 size={14} />
                  Remove from favorites
                </button>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

