import { useEffect, useState } from 'react';
import { useParams, Link, useSearchParams } from 'react-router-dom';
import { brandsApi } from '../api/brands';
import { devicesApi } from '../api/devices';
import type { Brand, Device, PaginationMeta } from '../types';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import Pagination from '../components/ui/Pagination';
import { ChevronRight, Smartphone, ArrowLeft, Heart } from 'lucide-react';
import { resolveDeviceImageUrl } from '../utils/resolveDeviceImageUrl';
import { useFavorites } from '../contexts/FavoritesContext';

const DEFAULT_LIMIT = 12;

const parsePage = (value: string | null) => {
  const parsed = Number(value);
  if (!Number.isFinite(parsed) || parsed < 1) return 1;
  return Math.floor(parsed);
};

export default function BrandDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { isFavorite, toggleFavorite } = useFavorites();
  const [searchParams, setSearchParams] = useSearchParams();
  const [brand, setBrand] = useState<Brand | null>(null);
  const [devices, setDevices] = useState<Device[]>([]);
  const [pagination, setPagination] = useState<PaginationMeta | null>(null);
  const [loadingBrand, setLoadingBrand] = useState(true);
  const [loadingDevices, setLoadingDevices] = useState(true);

  const currentPage = parsePage(searchParams.get('page'));
  const totalDevices = pagination?.total ?? devices.length;

  useEffect(() => {
    if (!id) {
      setBrand(null);
      setLoadingBrand(false);
      return;
    }

    let cancelled = false;
    setLoadingBrand(true);

    brandsApi
      .getById(id)
      .then(({ data }) => {
        if (!cancelled) {
          setBrand(data.data || null);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setBrand(null);
        }
      })
      .finally(() => {
        if (!cancelled) {
          setLoadingBrand(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [id]);

  useEffect(() => {
    if (!id) {
      setDevices([]);
      setPagination(null);
      setLoadingDevices(false);
      return;
    }

    let cancelled = false;
    setLoadingDevices(true);
    setDevices([]);
    setPagination(null);

    devicesApi
      .getByBrandId(id, currentPage, DEFAULT_LIMIT)
      .then(({ data }) => {
        if (cancelled) return;

        const nextDevices = data.data?.devices || [];
        const nextPagination = data.data?.pagination || null;

        if (nextPagination?.total_pages && currentPage > nextPagination.total_pages) {
          const nextParams = new URLSearchParams(searchParams);
          nextParams.set('page', String(nextPagination.total_pages));
          setSearchParams(nextParams, { replace: true });
          return;
        }

        setDevices(nextDevices);
        setPagination(nextPagination);
      })
      .catch(() => {
        if (!cancelled) {
          setDevices([]);
          setPagination(null);
        }
      })
      .finally(() => {
        if (!cancelled) {
          setLoadingDevices(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [id, currentPage, searchParams, setSearchParams]);

  if (loadingBrand) return <LoadingSpinner text="Loading brand..." />;

  if (!brand) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-20 text-center">
        <p className="text-text-secondary text-lg">Brand not found</p>
        <Link to="/brands" className="mt-4 inline-flex items-center gap-2 text-sm text-primary hover:text-primary-light">
          <ArrowLeft size={16} /> Back to brands
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-1.5 text-sm text-text-muted mb-6">
        <Link to="/brands" className="hover:text-text-primary transition-colors">Brands</Link>
        <ChevronRight size={14} />
        <span className="text-text-primary font-medium">{brand.brand_name}</span>
      </nav>

      {/* Header */}
      <div className="flex items-center gap-4 mb-8">
        <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary/20 to-accent/20 flex items-center justify-center">
          <span className="text-2xl font-bold gradient-text">
            {brand.brand_name?.[0]?.toUpperCase()}
          </span>
        </div>
        <div>
          <h1 className="text-3xl font-bold text-text-primary">{brand.brand_name}</h1>
          <p className="text-text-muted text-sm">
            {totalDevices} device{totalDevices !== 1 ? 's' : ''}
          </p>
        </div>
      </div>

      {/* Devices grid */}
      {loadingDevices ? (
        <LoadingSpinner text="Loading devices..." />
      ) : devices.length === 0 ? (
        <div className="text-center py-20 glass rounded-2xl">
          <Smartphone size={48} className="mx-auto text-text-muted mb-4" />
          <p className="text-text-secondary">No devices found for this brand</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-5">
          {devices.map((device, i) => (
            <Link
              key={device.id}
              to={`/devices/${device.id}`}
              className="glass rounded-2xl overflow-hidden card-hover group animate-fadeIn relative"
              style={{ animationDelay: `${i * 0.05}s`, opacity: 0 }}
            >
              <button
                type="button"
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  toggleFavorite(device.id);
                }}
                className={`absolute top-3 right-3 z-10 p-2 rounded-full border backdrop-blur-sm transition-colors ${
                  isFavorite(device.id)
                    ? 'text-danger border-danger/40 bg-danger/15'
                    : 'text-text-muted border-border bg-surface/60 hover:text-danger'
                }`}
                aria-label="Toggle favorite"
              >
                <Heart size={14} className={isFavorite(device.id) ? 'fill-danger' : ''} />
              </button>
              <div className="aspect-square bg-gradient-to-br from-surface-lighter/50 to-surface-light flex items-center justify-center p-6 overflow-hidden">
                {device.imageUrl ? (
                  <img
                    src={resolveDeviceImageUrl(device.imageUrl)}
                    alt={device.model_name}
                    className="max-h-full w-auto object-contain group-hover:scale-105 transition-transform duration-500"
                  />
                ) : (
                  <Smartphone size={48} className="text-text-muted" />
                )}
              </div>
              <div className="p-4">
                <h3 className="text-sm font-semibold text-text-primary mb-2 truncate">
                  {device.model_name}
                </h3>
                <div className="space-y-1.5 text-xs text-text-muted">
                  {device.specifications?.platform?.chipset && (
                    <p className="flex items-center gap-1.5">
                      <span className="w-1.5 h-1.5 rounded-full bg-primary flex-shrink-0" />
                      {device.specifications.platform.chipset.split('(')[0].trim()}
                    </p>
                  )}
                  {device.specifications?.display?.size && (
                    <p className="flex items-center gap-1.5">
                      <span className="w-1.5 h-1.5 rounded-full bg-accent flex-shrink-0" />
                      {device.specifications.display.size.split('(')[0].trim()}
                    </p>
                  )}
                  {device.specifications?.battery?.type && (
                    <p className="flex items-center gap-1.5">
                      <span className="w-1.5 h-1.5 rounded-full bg-success flex-shrink-0" />
                      {device.specifications.battery.type}
                    </p>
                  )}
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}

      {pagination && pagination.total_pages > 1 && (
        <Pagination
          pagination={pagination}
          onPageChange={(page) => {
            const nextParams = new URLSearchParams(searchParams);
            nextParams.set('page', String(page));
            setSearchParams(nextParams, { replace: true });
          }}
        />
      )}
    </div>
  );
}
