import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { devicesApi } from '../api/devices';
import type { Device } from '../types';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import { resolveDeviceImageUrl } from '../utils/resolveDeviceImageUrl';
import { useFavorites } from '../contexts/FavoritesContext';
import {
  ChevronRight, Smartphone, Wifi, Battery, Cpu, Monitor,
  Camera, Volume2, Ruler, Fingerprint, Info, ArrowLeft, BarChart3, Heart
} from 'lucide-react';

export default function DeviceDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { isFavorite, toggleFavorite } = useFavorites();
  const [device, setDevice] = useState<Device | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    devicesApi.getById(id)
      .then(({ data }) => setDevice(data.data || null))
      .catch(() => setDevice(null))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <LoadingSpinner text="Loading device..." />;

  if (!device) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-20 text-center">
        <p className="text-text-secondary text-lg">Device not found</p>
        <Link to="/" className="mt-4 inline-flex items-center gap-2 text-sm text-primary">
          <ArrowLeft size={16} /> Back to home
        </Link>
      </div>
    );
  }

  const specs = device.specifications;

  const specSections = [
    {
      title: 'Network',
      icon: Wifi,
      items: specs?.network ? [
        ['Technology', specs.network.technology],
        ['2G Bands', specs.network.bands_2g],
        ['3G Bands', specs.network.bands_3g],
        ['4G Bands', specs.network.bands_4g],
        ['5G Bands', specs.network.bands_5g],
        ['Speed', specs.network.speed],
      ] : [],
    },
    {
      title: 'Launch',
      icon: Info,
      items: specs?.launch ? [
        ['Announced', specs.launch.announced],
        ['Status', specs.launch.status],
      ] : [],
    },
    {
      title: 'Body',
      icon: Ruler,
      items: specs?.body ? [
        ['Dimensions', specs.body.dimensions],
        ['Weight', specs.body.weight],
        ['Build', specs.body.build],
        ['SIM', specs.body.sim],
        ['IP Rating', specs.body.ip_rating],
      ] : [],
    },
    {
      title: 'Display',
      icon: Monitor,
      items: specs?.display ? [
        ['Type', specs.display.type],
        ['Size', specs.display.size],
        ['Resolution', specs.display.resolution],
      ] : [],
    },
    {
      title: 'Platform',
      icon: Cpu,
      items: specs?.platform ? [
        ['OS', specs.platform.os],
        ['Chipset', specs.platform.chipset],
        ['CPU', specs.platform.cpu],
        ['GPU', specs.platform.gpu],
      ] : [],
    },
    {
      title: 'Memory',
      icon: Cpu,
      items: specs?.memory ? [
        ['Card Slot', specs.memory.card_lot],
        ['Internal', specs.memory.internal],
      ] : [],
    },
    {
      title: 'Main Camera',
      icon: Camera,
      items: specs?.mainCamera ? [
        ['Triple', specs.mainCamera.triple],
        ['Single', specs.mainCamera.single],
        ['Features', specs.mainCamera.features],
        ['Video', specs.mainCamera.video],
      ] : [],
    },
    {
      title: 'Selfie Camera',
      icon: Camera,
      items: specs?.selfieCamera ? [
        ['Single', specs.selfieCamera.single],
        ['Video', specs.selfieCamera.video],
      ] : [],
    },
    {
      title: 'Sound',
      icon: Volume2,
      items: specs?.sound ? [
        ['Loudspeaker', specs.sound.loudspeaker],
        ['3.5mm Jack', specs.sound['jack_3.5mm']],
      ] : [],
    },
    {
      title: 'Comms',
      icon: Wifi,
      items: specs?.comms ? [
        ['WLAN', specs.comms.wlan],
        ['Bluetooth', specs.comms.bluetooth],
        ['Positioning', specs.comms.positioning],
        ['NFC', specs.comms.nfc],
        ['Radio', specs.comms.radio],
        ['USB', specs.comms.usb],
      ] : [],
    },
    {
      title: 'Features',
      icon: Fingerprint,
      items: specs?.features ? [
        ['Sensors', specs.features.sensors],
      ] : [],
    },
    {
      title: 'Battery',
      icon: Battery,
      items: specs?.battery ? [
        ['Type', specs.battery.type],
        ['Charging', specs.battery.charging],
      ] : [],
    },
    {
      title: 'Misc',
      icon: Info,
      items: specs?.misc ? [
        ['Colors', specs.misc.colors],
        ['Models', specs.misc.models],
        ['Price', specs.misc.price],
      ] : [],
    },
  ];

  return (
    <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-1.5 text-sm text-text-muted mb-6">
        <Link to="/" className="hover:text-text-primary transition-colors">Home</Link>
        <ChevronRight size={14} />
        <span className="text-text-primary font-medium truncate">{device.model_name}</span>
      </nav>

      {/* Hero */}
      <div className="glass rounded-2xl p-6 md:p-8 mb-8">
        <div className="flex flex-col md:flex-row gap-8 items-center">
          <div className="w-48 h-48 md:w-64 md:h-64 bg-gradient-to-br from-surface-lighter/50 to-surface-light rounded-2xl flex items-center justify-center p-6 flex-shrink-0">
            {device.imageUrl ? (
              <img
                src={resolveDeviceImageUrl(device.imageUrl)}
                alt={device.model_name}
                className="max-h-full w-auto object-contain"
              />
            ) : (
              <Smartphone size={64} className="text-text-muted" />
            )}
          </div>
          <div className="flex-1 text-center md:text-left">
            <h1 className="text-2xl md:text-3xl font-bold text-text-primary">{device.model_name}</h1>

            {/* Quick specs */}
            <div className="mt-4 flex flex-wrap gap-2 justify-center md:justify-start">
              {specs?.platform?.chipset && (
                <span className="text-xs px-3 py-1 rounded-full bg-primary/10 text-primary font-medium">
                  {specs.platform.chipset.split('(')[0].trim()}
                </span>
              )}
              {specs?.display?.size && (
                <span className="text-xs px-3 py-1 rounded-full bg-accent/10 text-accent font-medium">
                  {specs.display.size.split('(')[0].trim()}
                </span>
              )}
              {specs?.battery?.type && (
                <span className="text-xs px-3 py-1 rounded-full bg-success/10 text-success font-medium">
                  {specs.battery.type}
                </span>
              )}
              {specs?.memory?.internal && (
                <span className="text-xs px-3 py-1 rounded-full bg-warning/10 text-warning font-medium">
                  {specs.memory.internal.split(',')[0].trim()}
                </span>
              )}
            </div>

            <div className="mt-6">
              <div className="flex flex-wrap gap-3 justify-center md:justify-start">
                <Link
                  to={`/compare?device=${id}`}
                  className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl text-sm font-semibold text-white btn-gradient"
                >
                  <BarChart3 size={16} />
                  Add to Compare
                </Link>
                <button
                  type="button"
                  onClick={() => toggleFavorite(device.id)}
                  className={`inline-flex items-center gap-2 px-5 py-2.5 rounded-xl text-sm font-semibold border transition-colors ${
                    isFavorite(device.id)
                      ? 'text-danger border-danger/40 bg-danger/10'
                      : 'text-text-secondary border-border hover:text-text-primary hover:bg-surface-light'
                  }`}
                >
                  <Heart size={16} className={isFavorite(device.id) ? 'fill-danger' : ''} />
                  {isFavorite(device.id) ? 'Favorited' : 'Add to favorites'}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Specification Tables */}
      <div className="space-y-4">
        {specSections.map((section) => {
          const validItems = section.items.filter(([, val]) => val);
          if (validItems.length === 0) return null;

          const Icon = section.icon;
          return (
            <div key={section.title} className="glass rounded-2xl overflow-hidden">
              <div className="flex items-center gap-3 px-5 py-3.5 border-b border-border bg-surface-light/30">
                <Icon size={18} className="text-primary" />
                <h2 className="text-sm font-semibold text-text-primary">{section.title}</h2>
              </div>
              <div className="spec-table">
                {validItems.map(([label, value], idx) => (
                  <div
                    key={idx}
                    className="flex border-b border-border last:border-0"
                  >
                    <div className="w-32 sm:w-40 flex-shrink-0 px-5 py-3 text-xs font-medium text-text-muted bg-surface-light/20">
                      {label}
                    </div>
                    <div className="flex-1 px-5 py-3 text-xs text-text-secondary leading-relaxed">
                      {value}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
