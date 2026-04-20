import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { useAuth } from '../../contexts/AuthContext';


import { Menu, X, Smartphone, LogOut, Shield, ChevronDown, KeyRound } from 'lucide-react';
import DeviceSearchBar from './DeviceSearchBar';
import ThemeToggle from '../ui/ThemeToggle';

export default function Navbar() {
  const { isAuthenticated, user, logout } = useAuth();
  const [mobileOpen, setMobileOpen] = useState(false);
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();

  const navLinks = [
    { to: '/', label: 'Home' },
    { to: '/brands', label: 'Brands' },
    { to: '/finder', label: 'Finder' },
    { to: '/compare', label: 'Compare' },
    { to: '/favorites', label: 'Favorites' },
  ];

  const isActive = (path: string) => location.pathname === path;

  const handleLogout = async () => {
    await logout();
    setUserMenuOpen(false);
    navigate('/');
  };

  return (
    <nav className="glass-strong sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center gap-2 group">
            <div className="w-9 h-9 rounded-lg btn-gradient flex items-center justify-center">
              <Smartphone size={20} className="text-white" />
            </div>
            <span className="text-xl font-bold gradient-text">TZone</span>
          </Link>

          {/* Desktop Nav */}
          <div className="hidden md:flex items-center gap-1">
            {navLinks.map((link) => (
              <Link
                key={link.to}
                to={link.to}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${
                  isActive(link.to)
                    ? 'bg-primary/10 text-primary'
                    : 'text-text-secondary hover:text-text-primary hover:bg-surface-light'
                }`}
              >
                {link.label}
              </Link>
            ))}
          </div>

          {/* Header Search */}
          <div className="hidden md:block w-72 lg:w-80">
            <DeviceSearchBar />
          </div>

          {/* Desktop Auth */}
          <div className="hidden md:flex items-center gap-3">
            <ThemeToggle />
            {isAuthenticated ? (
              <div className="relative">
                <button
                  onClick={() => setUserMenuOpen(!userMenuOpen)}
                  className="flex items-center gap-2 px-3 py-2 rounded-lg text-sm text-text-secondary hover:text-text-primary hover:bg-surface-light transition-all"
                >
                  <div className="w-7 h-7 rounded-full bg-gradient-to-br from-primary to-accent flex items-center justify-center text-white text-xs font-bold">
                    {user?.email?.[0]?.toUpperCase()}
                  </div>
                  <span className="max-w-[120px] truncate">{user?.email}</span>
                  <ChevronDown size={14} className={`transition-transform ${userMenuOpen ? 'rotate-180' : ''}`} />
                </button>
                {userMenuOpen && (
                  <div className="absolute right-0 mt-2 w-48 glass rounded-xl shadow-2xl py-2 animate-fadeIn">
                    {user?.role === 'admin' && (
                      <>
                        <Link
                          to="/admin"
                          onClick={() => setUserMenuOpen(false)}
                          className="flex items-center gap-2 px-4 py-2.5 text-sm text-text-secondary hover:text-text-primary hover:bg-surface-light transition-colors"
                        >
                          <Shield size={16} />
                          Admin Panel
                        </Link>
                        <hr className="border-border my-1" />
                      </>
                    )}
                    <Link
                      to="/change-password"
                      onClick={() => setUserMenuOpen(false)}
                      className="flex items-center gap-2 px-4 py-2.5 text-sm text-text-secondary hover:text-text-primary hover:bg-surface-light transition-colors"
                    >
                      <KeyRound size={16} />
                      Change Password
                    </Link>
                    <hr className="border-border my-1" />
                    <button
                      onClick={handleLogout}
                      className="flex items-center gap-2 w-full px-4 py-2.5 text-sm text-danger hover:bg-surface-light transition-colors"
                    >
                      <LogOut size={16} />
                      Logout
                    </button>
                  </div>
                )}
              </div>
            ) : (
              <>
                <Link
                  to="/login"
                  className="px-4 py-2 rounded-lg text-sm font-medium text-text-secondary hover:text-text-primary hover:bg-surface-light transition-all"
                >
                  Login
                </Link>
                <Link
                  to="/register"
                  className="px-5 py-2 rounded-lg text-sm font-semibold text-white btn-gradient"
                >
                  Register
                </Link>
              </>
            )}
          </div>

          {/* Mobile menu button */}
          <button
            onClick={() => setMobileOpen(!mobileOpen)}
            className="md:hidden p-2 rounded-lg text-text-secondary hover:text-text-primary hover:bg-surface-light transition-colors"
          >
            {mobileOpen ? <X size={22} /> : <Menu size={22} />}
          </button>
        </div>
      </div>

      {/* Mobile menu */}
      {mobileOpen && (
        <div className="md:hidden border-t border-border animate-fadeIn">
          <div className="px-4 py-4 space-y-1">
            <div className="pb-3">
              <DeviceSearchBar placeholder="Search devices in header..." />
            </div>
            {navLinks.map((link) => (
              <Link
                key={link.to}
                to={link.to}
                onClick={() => setMobileOpen(false)}
                className={`block px-4 py-3 rounded-lg text-sm font-medium transition-all ${
                  isActive(link.to)
                    ? 'bg-primary/10 text-primary'
                    : 'text-text-secondary hover:text-text-primary hover:bg-surface-light'
                }`}
              >
                {link.label}
              </Link>
            ))}
            <hr className="border-border my-2" />
            <ThemeToggle className="w-full justify-center border border-border" />
            <hr className="border-border my-2" />
            {isAuthenticated ? (
              <>
                {user?.role === 'admin' && (
                  <Link
                    to="/admin"
                    onClick={() => setMobileOpen(false)}
                    className="flex items-center gap-2 px-4 py-3 rounded-lg text-sm text-text-secondary hover:text-text-primary hover:bg-surface-light"
                  >
                    <Shield size={16} />
                    Admin Panel
                  </Link>
                )}
                <Link
                  to="/change-password"
                  onClick={() => setMobileOpen(false)}
                  className="flex items-center gap-2 px-4 py-3 rounded-lg text-sm text-text-secondary hover:text-text-primary hover:bg-surface-light"
                >
                  <KeyRound size={16} />
                  Change Password
                </Link>
                <button
                  onClick={() => { handleLogout(); setMobileOpen(false); }}
                  className="flex items-center gap-2 w-full px-4 py-3 rounded-lg text-sm text-danger hover:bg-surface-light"
                >
                  <LogOut size={16} />
                  Logout
                </button>
              </>
            ) : (
              <div className="flex gap-2 pt-2">
                <Link
                  to="/login"
                  onClick={() => setMobileOpen(false)}
                  className="flex-1 text-center px-4 py-2.5 rounded-lg text-sm font-medium text-text-secondary border border-border hover:bg-surface-light transition-all"
                >
                  Login
                </Link>
                <Link
                  to="/register"
                  onClick={() => setMobileOpen(false)}
                  className="flex-1 text-center px-4 py-2.5 rounded-lg text-sm font-semibold text-white btn-gradient"
                >
                  Register
                </Link>
              </div>
            )}
          </div>
        </div>
      )}
    </nav>
  );
}
