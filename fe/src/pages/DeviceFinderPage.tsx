import { useEffect, useMemo, useState } from 'react';
import type { FormEvent } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { Heart, Search, SlidersHorizontal, Smartphone } from 'lucide-react';
import { devicesApi } from '../api/devices';
import { brandsApi } from '../api/brands';
import type { Brand, Device, DeviceFinderParams, PaginationMeta } from '../types';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import Pagination from '../components/ui/Pagination';
import { resolveDeviceImageUrl } from '../utils/resolveDeviceImageUrl';
import { finderDropdownOptions } from '../constants/finderDropdownOptions';
import { useFavorites } from '../contexts/FavoritesContext';

const DEFAULT_LIMIT = 12;

const FILTER_DROPDOWN_FIELDS = [
  { key: 'os', label: 'OS: any' },
  { key: 'chipset', label: 'Chipset: any' },
  { key: 'cpu', label: 'CPU: any' },
  { key: 'gpu', label: 'GPU: any' },
  { key: 'memory', label: 'Memory: any' },
  { key: 'display_size', label: 'Display size: any' },
  { key: 'battery', label: 'Battery: any' },
  { key: 'nfc', label: 'NFC: any' },
] as const;

const parsePositiveInt = (value: string | null, fallback: number) => {
  const parsed = Number(value);
  if (!Number.isFinite(parsed) || parsed < 1) return fallback;
  return Math.floor(parsed);
};

export default function DeviceFinderPage() {
  const { isFavorite, toggleFavorite } = useFavorites();
  const [searchParams, setSearchParams] = useSearchParams();
  const [brands, setBrands] = useState<Brand[]>([]);
  const [devices, setDevices] = useState<Device[]>([]);
  const [pagination, setPagination] = useState<PaginationMeta | null>(null);
  const [loading, setLoading] = useState(true);

  const page = parsePositiveInt(searchParams.get('page'), 1);
  const limit = parsePositiveInt(searchParams.get('limit'), DEFAULT_LIMIT);

  const filters = useMemo(
    () => ({
      brand_id: searchParams.get('brand_id') || '',
      os: searchParams.get('os') || '',
      chipset: searchParams.get('chipset') || '',
      cpu: searchParams.get('cpu') || '',
      gpu: searchParams.get('gpu') || '',
      memory: searchParams.get('memory') || '',
      display_size: searchParams.get('display_size') || '',
      battery: searchParams.get('battery') || '',
      nfc: searchParams.get('nfc') || '',
    }),
    [searchParams],
  );
  const [draftFilters, setDraftFilters] = useState(filters);

  useEffect(() => {
    setDraftFilters(filters);
  }, [filters]);

  useEffect(() => {
    const fetchBrands = async () => {
      try {
        const { data } = await brandsApi.getAll(1, 200);
        setBrands(data.data?.brands || []);
      } catch {
        setBrands([]);
      }
    };

    fetchBrands();
  }, []);

  useEffect(() => {
    const fetchDevices = async () => {
      setLoading(true);
      try {
        const params: DeviceFinderParams = {
          ...filters,
          page,
          limit,
        };

        const { data } = await devicesApi.finder(params);
        setDevices(data.data?.devices || []);
        setPagination(data.data?.pagination || null);
      } catch {
        setDevices([]);
        setPagination(null);
      } finally {
        setLoading(false);
      }
    };

    fetchDevices();
  }, [filters, page, limit]);

  const updateQuery = (next: Record<string, string>, resetPage = false) => {
    const params = new URLSearchParams(searchParams);
    Object.entries(next).forEach(([key, value]) => {
      const trimmed = value.trim();
      if (trimmed) {
        params.set(key, trimmed);
      } else {
        params.delete(key);
      }
    });

    if (resetPage) {
      params.set('page', '1');
    }

    params.set('limit', String(limit));
    setSearchParams(params, { replace: true });
  };

  const activeFilterCount = useMemo(
    () => Object.values(filters).filter((value) => value.trim().length > 0).length,
    [filters],
  );

  const handleSearchSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    updateQuery(draftFilters, true);
  };

  const handleResetFilters = () => {
    setDraftFilters({
      brand_id: '',
      os: '',
      chipset: '',
      cpu: '',
      gpu: '',
      memory: '',
      display_size: '',
      battery: '',
      nfc: '',
    });
    setSearchParams({ page: '1', limit: String(limit) }, { replace: true });
  };

  const updateDraftFilter = (key: keyof typeof draftFilters, value: string) => {
    setDraftFilters((prev) => ({ ...prev, [key]: value }));
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-6">
        <div className="flex items-center gap-3 mb-1">
          <SlidersHorizontal size={22} className="text-primary" />
          <h1 className="text-2xl sm:text-3xl font-bold text-text-primary">Phone Finder</h1>
        </div>
        <p className="text-sm text-text-muted">Search devices by brand and common hardware specs</p>
      </div>

      <form onSubmit={handleSearchSubmit} className="grid grid-cols-1 lg:grid-cols-[280px_1fr] gap-5">
        <aside className="glass rounded-2xl p-4 h-fit">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-sm font-semibold text-text-primary">Filters</h2>
            <button
              type="button"
              onClick={handleResetFilters}
              className="text-xs text-primary hover:text-primary-light transition-colors"
            >
              Reset
            </button>
          </div>

          <div className="space-y-3">
            <select
              value={draftFilters.brand_id}
              onChange={(e) => updateDraftFilter('brand_id', e.target.value)}
              className="finder-select w-full px-3 py-2 rounded-lg border text-sm focus:outline-none"
            >
              <option value="">All brands</option>
              {brands.map((brand) => (
                <option key={brand.id} value={brand.id}>
                  {brand.brand_name}
                </option>
              ))}
            </select>

            {FILTER_DROPDOWN_FIELDS.map(({ key, label }) => (
              <select
                key={key}
                value={draftFilters[key]}
                onChange={(e) => updateDraftFilter(key, e.target.value)}
                className="finder-select w-full px-3 py-2 rounded-lg border text-sm focus:outline-none"
              >
                <option value="">{label}</option>
                {finderDropdownOptions[key].map((option) => (
                  <option key={option} value={option}>
                    {option}
                  </option>
                ))}
              </select>
            ))}

            <button
              type="submit"
              className="w-full inline-flex items-center justify-center gap-2 py-2 rounded-lg text-sm font-semibold text-white btn-gradient"
            >
              <Search size={16} />
              Find phones
            </button>
          </div>
        </aside>

        <section className="glass rounded-2xl p-4 sm:p-5">
          <div className="flex items-center justify-between gap-3 border-b border-border pb-3 mb-4">
            <p className="text-sm text-text-secondary">
              {pagination?.total ?? devices.length} result(s)
            </p>
            <p className="text-xs text-text-muted">
              {activeFilterCount > 0 ? `${activeFilterCount} filter(s) active` : 'No filters applied'}
            </p>
          </div>

          {loading ? (
            <LoadingSpinner text="Finding devices..." />
          ) : devices.length === 0 ? (
            <div className="text-center py-20">
              <Smartphone size={48} className="mx-auto text-text-muted mb-4" />
              <p className="text-text-secondary">No devices match the current filters</p>
            </div>
          ) : (
            <>
              <div className="divide-y divide-border">
                {devices.map((device) => (
                  <Link
                    key={device.id}
                    to={`/devices/${device.id}`}
                    className="py-4 first:pt-0 last:pb-0 flex gap-4 hover:bg-surface-light/40 transition-colors rounded-lg px-2 relative"
                  >
                    <button
                      type="button"
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        toggleFavorite(device.id);
                      }}
                      className={`absolute right-2 top-3 p-2 rounded-full border transition-colors ${
                        isFavorite(device.id)
                          ? 'text-danger border-danger/40 bg-danger/10'
                          : 'text-text-muted border-border hover:text-danger hover:bg-surface-light'
                      }`}
                      aria-label="Toggle favorite"
                    >
                      <Heart size={14} className={isFavorite(device.id) ? 'fill-danger' : ''} />
                    </button>
                    <div className="w-24 h-24 rounded-xl bg-surface-light flex items-center justify-center overflow-hidden shrink-0">
                      {device.imageUrl ? (
                        <img
                          src={resolveDeviceImageUrl(device.imageUrl)}
                          alt={device.model_name}
                          className="max-h-full w-auto object-contain"
                        />
                      ) : (
                        <Smartphone size={32} className="text-text-muted" />
                      )}
                    </div>

                    <div className="min-w-0 flex-1">
                      <h3 className="text-base font-semibold text-text-primary truncate">{device.model_name}</h3>
                      <div className="mt-1 space-y-1 text-xs text-text-muted">
                        {device.specifications?.platform?.os && <p><span className="text-text-secondary">OS:</span> {device.specifications.platform.os}</p>}
                        {device.specifications?.platform?.chipset && <p><span className="text-text-secondary">Chipset:</span> {device.specifications.platform.chipset}</p>}
                        {device.specifications?.memory?.internal && <p><span className="text-text-secondary">Memory:</span> {device.specifications.memory.internal}</p>}
                        {device.specifications?.display?.size && <p><span className="text-text-secondary">Display:</span> {device.specifications.display.size}</p>}
                        {device.specifications?.battery?.type && <p><span className="text-text-secondary">Battery:</span> {device.specifications.battery.type}</p>}
                      </div>
                    </div>
                  </Link>
                ))}
              </div>

              {pagination && (
                <Pagination
                  pagination={pagination}
                  onPageChange={(next) => updateQuery({ page: String(next) })}
                />
              )}
            </>
          )}
        </section>
      </form>
    </div>
  );
}

