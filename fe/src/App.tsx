import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import { AuthProvider } from './contexts/AuthContext';
import { FavoritesProvider } from './contexts/FavoritesContext';
import { ThemeProvider } from './contexts/ThemeContext';
import Layout from './components/layout/Layout';
import AdminLayout from './components/layout/AdminLayout';
import ProtectedRoute from './components/ProtectedRoute';
import AdminRoute from './components/AdminRoute';

import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import ForgotPasswordPage from './pages/ForgotPasswordPage';
import ChangePasswordPage from './pages/ChangePasswordPage';
import RegisterPage from './pages/RegisterPage';
import BrandsPage from './pages/BrandsPage';
import BrandDetailPage from './pages/BrandDetailPage';
import DeviceDetailPage from './pages/DeviceDetailPage';
import ComparePage from './pages/ComparePage';
import DeviceFinderPage from './pages/DeviceFinderPage';
import FavoritesPage from './pages/FavoritesPage';
import AdminDashboard from './pages/admin/AdminDashboard';
import BrandManagePage from './pages/admin/BrandManagePage';
import DeviceManagePage from './pages/admin/DeviceManagePage';

function AppContent() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <FavoritesProvider>
        <Toaster
          position="top-right"
          toastOptions={{
            style: {
              background: 'var(--tzone-surface-light)',
              color: 'var(--tzone-text-primary)',
              border: '1px solid var(--tzone-border)',
              borderRadius: '12px',
              fontSize: '14px',
            },
            success: {
              iconTheme: { primary: 'var(--tzone-success)', secondary: 'var(--tzone-text-primary)' },
            },
            error: {
              iconTheme: { primary: 'var(--tzone-danger)', secondary: 'var(--tzone-text-primary)' },
            },
          }}
        />

        <Routes>
          {/* Public routes with main layout */}
          <Route element={<Layout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/forgot-password" element={<ForgotPasswordPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route
              path="/change-password"
              element={
                <ProtectedRoute>
                  <ChangePasswordPage />
                </ProtectedRoute>
              }
            />
            <Route path="/brands" element={<BrandsPage />} />
            <Route path="/brands/:id" element={<BrandDetailPage />} />
            <Route path="/devices/:id" element={<DeviceDetailPage />} />
            <Route path="/finder" element={<DeviceFinderPage />} />
            <Route path="/compare" element={<ComparePage />} />
            <Route path="/favorites" element={<FavoritesPage />} />
          </Route>

        {/* Admin routes with admin layout */}
        <Route
          element={
            <AdminRoute>
              <AdminLayout />
            </AdminRoute>
          }
        >
          <Route path="/admin" element={<AdminDashboard />} />
          <Route path="/admin/brands" element={<BrandManagePage />} />
          <Route path="/admin/devices" element={<DeviceManagePage />} />
        </Route>
        </Routes>
        </FavoritesProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}

function App() {
  return (
    <ThemeProvider>
      <AppContent />
    </ThemeProvider>
  );
}

export default App;
