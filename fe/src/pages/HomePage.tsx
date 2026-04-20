import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { brandsApi } from '../api/brands';
import { devicesApi } from '../api/devices';
import type { Brand, Device } from '../types';
import { ArrowRight, Smartphone, Layers, BarChart3, Zap, ChevronRight } from 'lucide-react';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import { resolveDeviceImageUrl } from '../utils/resolveDeviceImageUrl';
import { getRecentlyViewedIds } from '../utils/recentlyViewed';

export default function HomePage() {
  const [brands, setBrands] = useState<Brand[]>([]);
  const [devices, setDevices] = useState<Device[]>([]);
  const [recentDevices, setRecentDevices] = useState<Device[]>([]);
  const [recentLoading, setRecentLoading] = useState(true);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [brandsRes, devicesRes] = await Promise.all([
          brandsApi.getAll(1, 8),
          devicesApi.getAll(1, 8),
        ]);
        setBrands(brandsRes.data.data?.brands || []);
        setDevices(devicesRes.data.data?.devices || []);
      } catch (err) {
        console.error('Failed to fetch homepage data', err);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  useEffect(() => {
    let active = true;

    const fetchRecentlyViewed = async () => {
      const ids = getRecentlyViewedIds().slice(0, 5);
      if (ids.length === 0) {
        if (active) {
          setRecentDevices([]);
          setRecentLoading(false);
        }
        return;
      }

      try {
        const results = await Promise.allSettled(ids.map((deviceId) => devicesApi.getById(deviceId)));
        if (!active) return;

        const byId = new Map<string, Device>();
        results.forEach((result) => {
          if (result.status === 'fulfilled') {
            const item = result.value.data.data;
            if (item?.id) {
              byId.set(item.id, item);
            }
          }
        });

        const ordered = ids.map((deviceId) => byId.get(deviceId)).filter((item): item is Device => Boolean(item));
        setRecentDevices(ordered);
      } finally {
        if (active) {
          setRecentLoading(false);
        }
      }
    };

    fetchRecentlyViewed();

    return () => {
      active = false;
    };
  }, []);

  return (
    <div>
      {/* Hero Section */}
      <section className="hero-gradient relative overflow-hidden">
        <div className="absolute inset-0">
          <div className="absolute top-20 left-10 w-72 h-72 bg-primary/5 rounded-full blur-3xl" />
          <div className="absolute bottom-10 right-10 w-96 h-96 bg-accent/5 rounded-full blur-3xl" />
        </div>

        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 lg:py-36">
          <div className="text-center max-w-3xl mx-auto">
            <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full glass text-xs font-medium text-primary mb-6 animate-fadeIn">
              <Zap size={14} />
              Compare smartphones like a pro
            </div>
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-extrabold tracking-tight animate-fadeIn stagger-1" style={{ opacity: 0 }}>
              Find Your Perfect{' '}
              <span className="gradient-text">Smartphone</span>
            </h1>
            <p className="mt-6 text-lg text-text-secondary max-w-2xl mx-auto animate-fadeIn stagger-2" style={{ opacity: 0 }}>
              Compare detailed specifications across hundreds of devices. Make informed decisions with side-by-side comparisons of display, camera, battery, and more.
            </p>
            <div className="mt-8 flex flex-col sm:flex-row gap-3 justify-center animate-fadeIn stagger-3" style={{ opacity: 0 }}>
              <Link
                to="/brands"
                className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-xl text-sm font-semibold text-white btn-gradient"
              >
                Browse Brands
                <ArrowRight size={18} />
              </Link>
              <Link
                to="/compare"
                className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-xl text-sm font-semibold text-text-primary border border-border hover:bg-surface-light transition-all"
              >
                <BarChart3 size={18} />
                Compare Devices
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[
            {
              icon: Smartphone,
              title: 'Detailed Specs',
              desc: 'Access over 15 specification categories for each device, from network bands to sensor data.',
              color: 'from-primary/20 to-primary/5',
            },
            {
              icon: BarChart3,
              title: 'Side-by-Side',
              desc: 'Compare up to 3 devices simultaneously with highlighted differences for quick analysis.',
              color: 'from-accent/20 to-accent/5',
            },
            {
              icon: Layers,
              title: 'All Major Brands',
              desc: 'Browse devices from all leading manufacturers, organized by brand for easy navigation.',
              color: 'from-success/20 to-success/5',
            },
          ].map((feat, i) => {
            const Icon = feat.icon;
            return (
              <div
                key={i}
                className="glass rounded-2xl p-6 card-hover"
              >
                <div className={`w-11 h-11 rounded-xl bg-gradient-to-br ${feat.color} flex items-center justify-center mb-4`}>
                  <Icon size={22} className="text-text-primary" />
                </div>
                <h3 className="text-base font-semibold text-text-primary mb-2">{feat.title}</h3>
                <p className="text-sm text-text-muted leading-relaxed">{feat.desc}</p>
              </div>
            );
          })}
        </div>
      </section>

      {/* stream + sidebar */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 pb-24">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
          <div className="lg:col-span-8">
            <div className="glass rounded-2xl overflow-hidden">
              <div className="px-5 py-4 border-b border-border bg-surface-light/40 flex items-center justify-between">
                <div>
                  <h2 className="text-xl font-bold text-text-primary">Latest Devices</h2>
                  <p className="text-xs text-text-muted mt-0.5">Quick browsing list</p>
                </div>
                <Link
                  to="/finder"
                  className="inline-flex items-center gap-1 text-xs font-semibold text-primary hover:text-primary-light"
                >
                  Open Finder <ChevronRight size={14} />
                </Link>
              </div>

              {loading ? (
                <div className="p-6">
                  <LoadingSpinner text="Loading devices..." />
                </div>
              ) : devices.length === 0 ? (
                <div className="p-8 text-center text-text-secondary text-sm">No devices available.</div>
              ) : (
                <div className="divide-y divide-border">
                  {devices.map((device) => (
                    <Link
                      key={device.id}
                      to={`/devices/${device.id}`}
                      className="flex gap-4 p-4 sm:p-5 hover:bg-surface-light/40 transition-colors"
                    >
                      <div className="w-20 h-24 sm:w-24 sm:h-28 shrink-0 rounded-lg bg-surface-light flex items-center justify-center overflow-hidden border border-border/60">
                        {device.imageUrl ? (
                          <img
                            src={resolveDeviceImageUrl(device.imageUrl)}
                            alt={device.model_name}
                            className="max-h-full w-auto object-contain"
                          />
                        ) : (
                          <Smartphone size={28} className="text-text-muted" />
                        )}
                      </div>

                      <div className="min-w-0 flex-1">
                        <h3 className="text-base font-semibold text-text-primary truncate">{device.model_name}</h3>
                        <div className="mt-2 space-y-1 text-xs text-text-muted">
                          {device.specifications?.platform?.os && (
                            <p><span className="text-text-secondary">OS:</span> {device.specifications.platform.os}</p>
                          )}
                          {device.specifications?.platform?.chipset && (
                            <p>
                              <span className="text-text-secondary">Chipset:</span> {device.specifications.platform.chipset.split('(')[0].trim()}
                            </p>
                          )}
                          {device.specifications?.display?.size && (
                            <p><span className="text-text-secondary">Display:</span> {device.specifications.display.size.split('(')[0].trim()}</p>
                          )}
                          {device.specifications?.battery?.type && (
                            <p><span className="text-text-secondary">Battery:</span> {device.specifications.battery.type}</p>
                          )}
                        </div>
                      </div>

                      <div className="hidden sm:flex items-center text-primary">
                        <ChevronRight size={18} />
                      </div>
                    </Link>
                  ))}
                </div>
              )}
            </div>
          </div>

          <aside className="lg:col-span-4 space-y-5">
            <div className="glass rounded-2xl overflow-hidden">
              <div className="px-4 py-3 border-b border-border bg-surface-light/40">
                <h2 className="text-sm font-semibold text-text-primary">Recently Viewed</h2>
              </div>

              {recentLoading ? (
                <div className="p-4">
                  <LoadingSpinner text="Loading recently viewed..." />
                </div>
              ) : recentDevices.length === 0 ? (
                <div className="p-4 text-sm text-text-secondary">No recently viewed devices yet.</div>
              ) : (
                <div className="divide-y divide-border">
                  {recentDevices.map((device) => (
                    <Link
                      key={device.id}
                      to={`/devices/${device.id}`}
                      className="px-4 py-3 flex items-center gap-3 hover:bg-surface-light/40 transition-colors"
                    >
                      <div className="w-11 h-11 rounded-lg bg-surface-light border border-border/60 flex items-center justify-center overflow-hidden shrink-0">
                        {device.imageUrl ? (
                          <img
                            src={resolveDeviceImageUrl(device.imageUrl)}
                            alt={device.model_name}
                            className="max-h-full w-auto object-contain"
                          />
                        ) : (
                          <Smartphone size={16} className="text-text-muted" />
                        )}
                      </div>
                      <div className="min-w-0 flex-1">
                        <p className="text-sm font-medium text-text-primary truncate">{device.model_name}</p>
                        <p className="text-xs text-text-muted truncate">
                          {device.specifications?.platform?.chipset?.split('(')[0].trim() || 'Device'}
                        </p>
                      </div>
                      <ChevronRight size={14} className="text-text-muted" />
                    </Link>
                  ))}
                </div>
              )}
            </div>

            <div className="glass rounded-2xl overflow-hidden">
              <div className="px-4 py-3 border-b border-border bg-surface-light/40 flex items-center justify-between">
                <h2 className="text-sm font-semibold text-text-primary">Brand Directory</h2>
                <Link
                  to="/brands"
                  className="inline-flex items-center gap-1 text-xs font-semibold text-primary hover:text-primary-light"
                >
                  View all <ChevronRight size={14} />
                </Link>
              </div>

              {loading ? (
                <div className="p-4">
                  <LoadingSpinner text="Loading brands..." />
                </div>
              ) : brands.length === 0 ? (
                <div className="p-4 text-sm text-text-secondary">No brands found.</div>
              ) : (
                <div className="divide-y divide-border">
                  {brands.map((brand) => (
                    <Link
                      key={brand.id}
                      to={`/brands/${brand.id}`}
                      className="px-4 py-3 flex items-center gap-3 hover:bg-surface-light/40 transition-colors"
                    >
                      <div className="w-9 h-9 rounded-lg bg-gradient-to-br from-surface-lighter to-surface-light flex items-center justify-center border border-border/60">
                        <span className="text-sm font-bold gradient-text">{brand.brand_name?.[0]?.toUpperCase()}</span>
                      </div>
                      <div className="min-w-0 flex-1">
                        <p className="text-sm font-medium text-text-primary truncate">{brand.brand_name}</p>
                      </div>
                      <ChevronRight size={14} className="text-text-muted" />
                    </Link>
                  ))}
                </div>
              )}
            </div>

            <div className="glass rounded-2xl p-4">
              <h3 className="text-sm font-semibold text-text-primary mb-2">Quick Actions</h3>
              <div className="space-y-2">
                <Link
                  to="/finder"
                  className="flex items-center justify-between text-sm text-text-secondary hover:text-text-primary px-3 py-2 rounded-lg hover:bg-surface-light"
                >
                  Device Finder
                  <ChevronRight size={14} />
                </Link>
                <Link
                  to="/compare"
                  className="flex items-center justify-between text-sm text-text-secondary hover:text-text-primary px-3 py-2 rounded-lg hover:bg-surface-light"
                >
                  Compare Devices
                  <ChevronRight size={14} />
                </Link>
              </div>
            </div>
          </aside>
        </div>
      </section>
    </div>
  );
}
