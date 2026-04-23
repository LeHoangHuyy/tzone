import { useEffect, useState } from 'react';
import type { FormEvent } from 'react';
import { useParams, Link } from 'react-router-dom';
import { isAxiosError } from 'axios';
import { devicesApi } from '../api/devices';
import { aiApi } from '../api/ai';
import { reviewsApi } from '../api/reviews';
import type { AIVideoReview, Device, Review } from '../types';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import { resolveDeviceImageUrl } from '../utils/resolveDeviceImageUrl';
import { pushRecentlyViewedId } from '../utils/recentlyViewed';
import { useFavorites } from '../contexts/FavoritesContext';
import { useAuth } from '../contexts/AuthContext';

function extractYouTubeId(url: string): string | null {
  try {
    const urlObj = new URL(url);

    if (urlObj.hostname.includes('youtube.com')) {
      return urlObj.searchParams.get('v');
    } else if (urlObj.hostname.includes('youtu.be')) {
      return urlObj.pathname.slice(1);
    }
  } catch {
    // URL parsing failed, try regex
  }

  // Fallback regex matching
  const regexMatch = url.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&\n?#]+)/);
  return regexMatch ? regexMatch[1] : null;
}
import {
  ChevronRight,
  Smartphone,
  Wifi,
  Battery,
  Cpu,
  Monitor,
  Camera,
  Volume2,
  Ruler,
  Fingerprint,
  Info,
  ArrowLeft,
  BarChart3,
  Heart,
  PlayCircle,
  MessageSquare,
  Star,
  ExternalLink,
} from 'lucide-react';

const COMMENT_PAGE_SIZE = 5;

function getAiErrorMessage(error: unknown): string {
  if (isAxiosError(error)) {
    const message = error.response?.data?.message;
    const fieldError = error.response?.data?.errors?.[0]?.error;
    return message || fieldError || 'Unable to load video reviews right now.';
  }
  return 'Unable to load video reviews right now.';
}

function getReviewErrorMessage(error: unknown): string {
  if (isAxiosError(error)) {
    const message = error.response?.data?.message;
    const fieldError = error.response?.data?.errors?.[0]?.error;
    return message || fieldError || 'Unable to process your input right now.';
  }
  return 'Unable to process your input right now.';
}

function renderStars(rating: number): string {
  if (rating <= 0) {
    return 'No rating yet';
  }
  const safe = Math.max(1, Math.min(5, rating));
  return '★'.repeat(safe) + '☆'.repeat(5 - safe);
}

export default function DeviceDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { isFavorite, toggleFavorite } = useFavorites();
  const { user, isAuthenticated } = useAuth();

  const [device, setDevice] = useState<Device | null>(null);
  const [loading, setLoading] = useState(true);

  const [videoReviews, setVideoReviews] = useState<AIVideoReview[]>([]);
  const [videoReply, setVideoReply] = useState('');
  const [videoLoading, setVideoLoading] = useState(false);
  const [videoError, setVideoError] = useState('');

  const [reviews, setReviews] = useState<Review[]>([]);
  const [reviewsLoading, setReviewsLoading] = useState(false);
  const [reviewsError, setReviewsError] = useState('');

  const [commentPage, setCommentPage] = useState(1);
  const [commentTotalPages, setCommentTotalPages] = useState(1);
  const [commentTotal, setCommentTotal] = useState(0);

  const [ratingAverage, setRatingAverage] = useState<number | null>(null);
  const [ratingCount, setRatingCount] = useState(0);

  const [ratingSaving, setRatingSaving] = useState(false);
  const [commentSaving, setCommentSaving] = useState(false);
  const [ratingInput, setRatingInput] = useState(5);
  const [commentInput, setCommentInput] = useState('');

  const [editingCommentId, setEditingCommentId] = useState<string | null>(null);
  const [editingCommentText, setEditingCommentText] = useState('');

  const isAdmin = user?.email?.toLowerCase() === 'admin@tzone.com';

  const loadComments = async (deviceId: string, page: number) => {
    setReviewsLoading(true);
    setReviewsError('');
    try {
      const { data } = await reviewsApi.getByDeviceId(deviceId, page, COMMENT_PAGE_SIZE);
      const payload = data.data;
      setReviews(payload?.reviews || []);
      setCommentTotal(payload?.total || 0);
      setCommentTotalPages(Math.max(payload?.pagination?.total_pages || 1, 1));
      setRatingAverage(payload?.rating_summary?.count ? payload.rating_summary.average : null);
      setRatingCount(payload?.rating_summary?.count || 0);
    } catch (error) {
      setReviews([]);
      setCommentTotal(0);
      setCommentTotalPages(1);
      setRatingAverage(null);
      setRatingCount(0);
      setReviewsError(getReviewErrorMessage(error));
    } finally {
      setReviewsLoading(false);
    }
  };

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    devicesApi.getById(id)
      .then(({ data }) => setDevice(data.data || null))
      .catch(() => setDevice(null))
      .finally(() => setLoading(false));
  }, [id]);

  useEffect(() => {
    if (!id) return;
    loadComments(id, commentPage);
  }, [id, commentPage]);

  useEffect(() => {
    if (device?.id) {
      pushRecentlyViewedId(device.id);
    }
  }, [device?.id]);

  useEffect(() => {
    const deviceName = device?.model_name?.trim();
    if (!deviceName) {
      setVideoReviews([]);
      setVideoReply('');
      setVideoError('');
      setVideoLoading(false);
      return;
    }

    let active = true;
    setVideoLoading(true);
    setVideoError('');

    aiApi.videoReviews({ device_name: deviceName, limit: 1 })
      .then(({ data }) => {
        if (!active) return;
        const payload = data.data;
        setVideoReviews(payload?.videos || []);
        setVideoReply(payload?.reply || '');
      })
      .catch((error) => {
        if (!active) return;
        setVideoReviews([]);
        setVideoReply('');
        setVideoError(getAiErrorMessage(error));
      })
      .finally(() => {
        if (!active) return;
        setVideoLoading(false);
      });

    return () => {
      active = false;
    };
  }, [device?.model_name]);


  const handleSaveRating = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!id || !isAuthenticated) return;

    setRatingSaving(true);
    setReviewsError('');
    try {
      await reviewsApi.setRating(id, ratingInput);
      await loadComments(id, commentPage);
    } catch (error) {
      setReviewsError(getReviewErrorMessage(error));
    } finally {
      setRatingSaving(false);
    }
  };

  const handleSaveComment = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!id || !isAuthenticated) return;

    const comment = commentInput.trim();
    if (!comment) {
      setReviewsError('Comment is required.');
      return;
    }

    setCommentSaving(true);
    setReviewsError('');
    try {
      await reviewsApi.setComment(id, comment);
      setCommentPage(1);
      await loadComments(id, 1);
    } catch (error) {
      setReviewsError(getReviewErrorMessage(error));
    } finally {
      setCommentSaving(false);
    }
  };

  const handleStartEditComment = (review: Review) => {
    setEditingCommentId(review.id);
    setEditingCommentText(review.comment);
    window.location.hash = 'comments';
  };

  const handleUpdateComment = async (reviewId: string) => {
    if (!editingCommentText.trim()) {
      setReviewsError('Comment is required.');
      return;
    }

    setCommentSaving(true);
    setReviewsError('');
    try {
      await reviewsApi.updateComment(reviewId, editingCommentText.trim());
      setEditingCommentId(null);
      setEditingCommentText('');
      if (id) {
        await loadComments(id, commentPage);
      }
    } catch (error) {
      setReviewsError(getReviewErrorMessage(error));
    } finally {
      setCommentSaving(false);
    }
  };

  const handleDeleteReview = async (reviewId: string) => {
    if (!id) return;
    setCommentSaving(true);
    setReviewsError('');
    try {
      await reviewsApi.remove(reviewId);
      const nextPage = reviews.length === 1 && commentPage > 1 ? commentPage - 1 : commentPage;
      setCommentPage(nextPage);
      await loadComments(id, nextPage);
    } catch (error) {
      setReviewsError(getReviewErrorMessage(error));
    } finally {
      setCommentSaving(false);
    }
  };

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
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <nav className="flex items-center gap-1.5 text-sm text-text-muted mb-6">
        <Link to="/" className="hover:text-text-primary transition-colors">Home</Link>
        <ChevronRight size={14} />
        <span className="text-text-primary font-medium truncate">{device.model_name}</span>
      </nav>

      <div className="grid grid-cols-1 lg:grid-cols-[220px_1fr] gap-6">
        <aside className="lg:sticky lg:top-24 h-fit glass rounded-2xl p-4">
          <p className="text-xs uppercase tracking-wide text-text-muted mb-3">Sections</p>
          <div className="space-y-2 text-sm">
            <a href="#device-overview" className="block text-text-secondary hover:text-primary">Device</a>
            <a href="#specifications" className="block text-text-secondary hover:text-primary">Specifications</a>
            <a href="#video-reviews" className="block text-text-secondary hover:text-primary">Video Reviews</a>
            <a href="#comments" className="block text-text-secondary hover:text-primary">Comments</a>
          </div>
        </aside>

        <div>
          <section id="device-overview" className="glass rounded-2xl p-6 md:p-8 mb-8 scroll-mt-24">
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
          </section>

          <section id="specifications" className="space-y-4 scroll-mt-24 mb-8">
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
          </section>

          <section id="video-reviews" className="glass rounded-2xl p-6 md:p-8 mb-8 scroll-mt-24">
            <div className="flex items-center gap-2 mb-3">
              <PlayCircle size={18} className="text-primary" />
              <h2 className="text-lg font-semibold text-text-primary">Video Reviews</h2>
            </div>

            {videoReply && <p className="text-sm text-text-secondary mb-4">{videoReply}</p>}
            {videoError && <p className="text-sm text-danger mb-4">{videoError}</p>}

            {videoLoading ? (
              <p className="text-sm text-text-muted">Loading video review...</p>
            ) : videoReviews.length === 0 ? (
              <p className="text-sm text-text-muted">No video review available yet.</p>
            ) : (
              <div className="space-y-4">
                {videoReviews.map((video, index) => {
                  const videoId = extractYouTubeId(video.url);
                  const isYouTubeSearchLink = video.url.includes('search_query');
                  
                  if (isYouTubeSearchLink || !videoId) {
                    return (
                      <a
                        key={`${video.url}-${index}`}
                        href={video.url}
                        target="_blank"
                        rel="noreferrer noopener"
                        className="group block rounded-xl border border-border bg-surface-light overflow-hidden hover:border-primary/40 transition-colors"
                      >
                        <div className="relative aspect-video bg-gradient-to-br from-surface-lighter to-surface-light flex items-center justify-center">
                          <div className="flex flex-col items-center gap-3">
                            <div className="w-16 h-16 rounded-full bg-primary/20 flex items-center justify-center group-hover:bg-primary/30 transition-colors">
                              <PlayCircle size={36} className="text-primary" />
                            </div>
                            <span className="text-xs text-text-muted">Open on YouTube</span>
                          </div>
                        </div>
                        <div className="flex items-center gap-2 p-4">
                          <ExternalLink size={14} className="text-primary shrink-0" />
                          <p className="text-sm font-medium text-text-primary truncate">{video.title}</p>
                        </div>
                      </a>
                    );
                  }
                  
                  return (
                    <div key={`${video.url}-${index}`} className="rounded-xl border border-border bg-surface-light overflow-hidden">
                      <div className="aspect-video overflow-hidden bg-black">
                        <iframe
                          width="100%"
                          height="100%"
                          src={`https://www.youtube.com/embed/${videoId}`}
                          title={video.title}
                          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                          allowFullScreen
                          className="w-full h-full"
                        />
                      </div>
                      <div className="flex items-center gap-2 p-4">
                        <a
                          href={video.url}
                          target="_blank"
                          rel="noreferrer noopener"
                          className="text-sm font-medium text-text-primary hover:text-primary transition-colors truncate"
                        >
                          {video.title}
                        </a>
                        <a
                          href={video.url}
                          target="_blank"
                          rel="noreferrer noopener"
                          className="shrink-0 text-text-muted hover:text-primary transition-colors"
                          title="Open on YouTube"
                        >
                          <ExternalLink size={14} />
                        </a>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </section>

          <section id="comments" className="glass rounded-2xl p-6 md:p-8 mb-8 scroll-mt-24">
            <div className="flex items-center justify-between gap-3 mb-4">
              <div className="flex items-center gap-2">
                <MessageSquare size={18} className="text-primary" />
                <h2 className="text-lg font-semibold text-text-primary">Comments</h2>
              </div>
              <div className="text-sm text-text-muted">
                {ratingAverage !== null ? `${ratingAverage.toFixed(1)}/5 (${ratingCount} ratings)` : 'No ratings yet'}
              </div>
            </div>

            {isAuthenticated ? (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-5">
                <form onSubmit={handleSaveRating} className="rounded-xl border border-border bg-surface-light p-4">
                  <p className="text-sm font-medium text-text-primary mb-2">Rate this device</p>
                  <div className="flex gap-2 mb-3">
                    <select
                      value={ratingInput}
                      onChange={(event) => setRatingInput(Number(event.target.value))}
                      className="px-3 py-2 rounded-lg border border-border bg-surface text-text-primary"
                      disabled={ratingSaving}
                    >
                      {[5, 4, 3, 2, 1].map((value) => (
                        <option key={value} value={value}>{value} Star{value > 1 ? 's' : ''}</option>
                      ))}
                    </select>
                    <input
                      value={renderStars(ratingInput)}
                      readOnly
                      className="flex-1 px-3 py-2 rounded-lg border border-border bg-background text-warning"
                    />
                  </div>
                  <button
                    type="submit"
                    disabled={ratingSaving}
                    className="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-white btn-gradient disabled:opacity-60"
                  >
                    <Star size={14} />
                    {ratingSaving ? 'Saving...' : 'Save rating'}
                  </button>
                </form>

                <form onSubmit={handleSaveComment} className="rounded-xl border border-border bg-surface-light p-4">
                  <p className="text-sm font-medium text-text-primary mb-2">Write your comment</p>
                  <textarea
                    value={commentInput}
                    onChange={(event) => setCommentInput(event.target.value)}
                    rows={4}
                    maxLength={2000}
                    placeholder="Share your hands-on experience with this device"
                    className="w-full px-3 py-2 rounded-lg border border-border bg-surface text-text-primary focus:outline-none"
                    disabled={commentSaving}
                  />
                  <button
                    type="submit"
                    disabled={commentSaving || !commentInput.trim()}
                    className="mt-3 inline-flex items-center gap-2 px-4 py-2 rounded-lg text-white btn-gradient disabled:opacity-60"
                  >
                    {commentSaving ? 'Saving...' : 'Save comment'}
                  </button>
                </form>
              </div>
            ) : (
              <p className="text-sm text-text-muted mb-5">Please log in to submit a rating or comment.</p>
            )}

            {reviewsError && <p className="text-sm text-danger mb-4">{reviewsError}</p>}

            {reviewsLoading ? (
              <p className="text-sm text-text-muted">Loading comments...</p>
            ) : reviews.length === 0 ? (
              <p className="text-sm text-text-muted">No comments yet.</p>
            ) : (
              <div className="space-y-3">
                {reviews.map((review) => {
                  const canEditComment = !!user && (review.user_id === user.id || isAdmin);
                  const canDelete = !!user && isAdmin;

                  return (
                    <div key={review.id} className="rounded-xl border border-border bg-surface-light p-4">
                      <div className="flex flex-wrap items-center justify-between gap-2">
                        <div>
                          <p className="text-sm font-semibold text-text-primary">{review.user_email || 'Unknown user'}</p>
                          <p className="text-xs text-warning">{renderStars(review.rating)}</p>
                        </div>
                        <p className="text-xs text-text-muted">{new Date(review.updated_at).toLocaleString()}</p>
                      </div>

                      {editingCommentId === review.id ? (
                        <div className="mt-3">
                          <textarea
                            value={editingCommentText}
                            onChange={(event) => setEditingCommentText(event.target.value)}
                            rows={3}
                            className="w-full px-3 py-2 rounded-lg border border-border bg-surface text-text-primary"
                          />
                          <div className="flex gap-2 mt-2">
                            <button
                              type="button"
                              onClick={() => handleUpdateComment(review.id)}
                              className="px-3 py-1.5 rounded-md text-xs text-white btn-gradient"
                              disabled={commentSaving || !editingCommentText.trim()}
                            >
                              Update comment
                            </button>
                            <button
                              type="button"
                              onClick={() => {
                                setEditingCommentId(null);
                                setEditingCommentText('');
                              }}
                              className="px-3 py-1.5 rounded-md border border-border text-xs text-text-secondary"
                            >
                              Cancel
                            </button>
                          </div>
                        </div>
                      ) : (
                        <p className="text-sm text-text-secondary mt-2 whitespace-pre-wrap">{review.comment}</p>
                      )}

                      {(canEditComment || canDelete) && editingCommentId !== review.id && (
                        <div className="flex gap-2 mt-3">
                          {canEditComment && (
                            <button
                              type="button"
                              onClick={() => handleStartEditComment(review)}
                              className="px-3 py-1.5 rounded-md border border-border text-xs text-text-secondary hover:text-text-primary"
                            >
                              Edit comment
                            </button>
                          )}
                          {canDelete && (
                            <button
                              type="button"
                              onClick={() => handleDeleteReview(review.id)}
                              className="px-3 py-1.5 rounded-md border border-danger/30 bg-danger/10 text-xs text-danger"
                              disabled={commentSaving}
                            >
                              Delete
                            </button>
                          )}
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            )}

            {commentTotalPages > 1 && (
              <div className="mt-5 flex items-center justify-between gap-3 text-sm">
                <p className="text-text-muted">{commentTotal} comments total</p>
                <div className="flex items-center gap-2">
                  <button
                    type="button"
                    onClick={() => setCommentPage((prev) => Math.max(1, prev - 1))}
                    disabled={commentPage <= 1}
                    className="px-3 py-1.5 rounded-md border border-border text-text-secondary disabled:opacity-50"
                  >
                    Prev
                  </button>
                  <span className="text-text-primary">Page {commentPage} / {commentTotalPages}</span>
                  <button
                    type="button"
                    onClick={() => setCommentPage((prev) => Math.min(commentTotalPages, prev + 1))}
                    disabled={commentPage >= commentTotalPages}
                    className="px-3 py-1.5 rounded-md border border-border text-text-secondary disabled:opacity-50"
                  >
                    Next
                  </button>
                </div>
              </div>
            )}
          </section>
        </div>
      </div>
    </div>
  );
}
