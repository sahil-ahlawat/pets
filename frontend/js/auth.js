/**
 * auth.js
 * Authentication and session management helper.
 */

class AuthService {
  isAuthenticated() {
    return !!localStorage.getItem('pet_access_token');
  }

  getUser() {
    const raw = localStorage.getItem('pet_user');
    return raw ? JSON.parse(raw) : null;
  }

  setUser(user) {
    localStorage.setItem('pet_user', JSON.stringify(user));
  }

  saveSession(authResponse) {
    if (authResponse && authResponse.tokens) {
      localStorage.setItem('pet_access_token', authResponse.tokens.accessToken);
      localStorage.setItem('pet_refresh_token', authResponse.tokens.refreshToken);
      this.setUser(authResponse.user);
    }
  }

  logout() {
    localStorage.removeItem('pet_access_token');
    localStorage.removeItem('pet_refresh_token');
    localStorage.removeItem('pet_user');
    localStorage.removeItem('active_pet_id');
    window.location.href = '/pages/login.html';
  }

  requireAuth() {
    if (!this.isAuthenticated()) {
      window.location.href = '/pages/login.html';
    }
  }

  requireGuest() {
    if (this.isAuthenticated()) {
      window.location.href = '/pages/dashboard.html';
    }
  }
}

window.auth = new AuthService();
