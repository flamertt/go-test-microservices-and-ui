import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import Navbar from './components/Navbar';
import ProtectedRoute from './components/ProtectedRoute';
import HomePage from './pages/HomePage';
import BooksPage from './pages/BooksPage';
import AuthorsPage from './pages/AuthorsPage';
import GenresPage from './pages/GenresPage';
import BookDetailPage from './pages/BookDetailPage';
import AuthorDetailPage from './pages/AuthorDetailPage';
import GenreDetailPage from './pages/GenreDetailPage';
import RecommendationsPage from './pages/RecommendationsPage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ProfilePage from './pages/ProfilePage';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Navbar />
        <div className="app-container">
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            
            {/* Korumalı Sayfalar */}
            <Route 
              path="/books" 
              element={
                <ProtectedRoute>
                  <BooksPage />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="/books/:id" 
              element={
                <ProtectedRoute>
                  <BookDetailPage />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="/authors" 
              element={
                <ProtectedRoute>
                  <AuthorsPage />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="/authors/:name" 
              element={
                <ProtectedRoute>
                  <AuthorDetailPage />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="/genres" 
              element={
                <ProtectedRoute>
                  <GenresPage />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="/genres/:name" 
              element={
                <ProtectedRoute>
                  <GenreDetailPage />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="/recommendations" 
              element={
                <ProtectedRoute>
                  <RecommendationsPage />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="/profile" 
              element={
                <ProtectedRoute>
                  <ProfilePage />
                </ProtectedRoute>
              } 
            />
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;
