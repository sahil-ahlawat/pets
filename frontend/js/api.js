/**
 * api.js
 * Centralized API client for Pet Digital Identity Platform.
 * Communicates with backend /api/v1/ endpoints with automatic JWT injection.
 */

const API_BASE = window.PET_API_BASE || (
  (window.location.port === '3000' || window.location.port === '80' || window.location.port === '' || !window.location.origin.includes('localhost'))
    ? '/api/v1'
    : 'http://localhost:8080/api/v1'
);

class ApiClient {
  constructor() {
    this.baseUrl = API_BASE;
  }

  getAuthHeader() {
    const token = localStorage.getItem('pet_access_token');
    return token ? { 'Authorization': `Bearer ${token}` } : {};
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseUrl}${endpoint}`;
    const headers = {
      'Content-Type': 'application/json',
      ...this.getAuthHeader(),
      ...options.headers,
    };

    try {
      const response = await fetch(url, { ...options, headers });

      // Handle 401 Unauthorized token refresh
      if (response.status === 401 && !endpoint.includes('/auth/login') && !endpoint.includes('/auth/refresh')) {
        const refreshed = await this.refreshToken();
        if (refreshed) {
          // Retry original request
          headers['Authorization'] = `Bearer ${localStorage.getItem('pet_access_token')}`;
          return fetch(url, { ...options, headers }).then(res => res.json());
        } else {
          localStorage.removeItem('pet_access_token');
          localStorage.removeItem('pet_refresh_token');
          window.location.href = '/pages/login.html';
          return null;
        }
      }

      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.error?.message || data.message || 'API request failed');
      }
      return data;
    } catch (err) {
      console.error(`API Error [${endpoint}]:`, err);
      throw err;
    }
  }

  async refreshToken() {
    const refreshToken = localStorage.getItem('pet_refresh_token');
    if (!refreshToken) return false;

    try {
      const res = await fetch(`${this.baseUrl}/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refreshToken })
      });
      if (!res.ok) return false;
      const data = await res.json();
      localStorage.setItem('pet_access_token', data.accessToken);
      localStorage.setItem('pet_refresh_token', data.refreshToken);
      return true;
    } catch {
      return false;
    }
  }

  // Auth
  register(payload) { return this.request('/auth/register', { method: 'POST', body: JSON.stringify(payload) }); }
  login(payload) { return this.request('/auth/login', { method: 'POST', body: JSON.stringify(payload) }); }
  getMe() { return this.request('/auth/me'); }

  // Pets
  getPets() { return this.request('/pets'); }
  getPet(id) { return this.request(`/pets/${id}`); }
  createPet(payload) { return this.request('/pets', { method: 'POST', body: JSON.stringify(payload) }); }
  updatePet(id, payload) { return this.request(`/pets/${id}`, { method: 'PATCH', body: JSON.stringify(payload) }); }
  deletePet(id) { return this.request(`/pets/${id}`, { method: 'DELETE' }); }
  markPetLost(id, payload) { return this.request(`/pets/${id}/lost`, { method: 'POST', body: JSON.stringify(payload) }); }
  markPetReunited(id) { return this.request(`/pets/${id}/reunited`, { method: 'POST' }); }

  // QR
  getPetQRCodes(petId) { return this.request(`/pets/${petId}/qr`); }
  generateQRCode(petId, tagLabel) { return this.request(`/pets/${petId}/qr/generate`, { method: 'POST', body: JSON.stringify({ tagLabel }) }); }
  revokeQRCode(qrId, reason) { return this.request(`/qr/${qrId}/revoke`, { method: 'POST', body: JSON.stringify({ reason }) }); }

  // Public QR Profile (No auth required)
  getPublicPet(token, coords) {
    const query = coords ? `?lat=${coords.lat}&lng=${coords.lng}` : '';
    return this.request(`/public/pets/${token}${query}`);
  }
  contactOwner(token, payload) { return this.request(`/public/pets/${token}/contact`, { method: 'POST', body: JSON.stringify(payload) }); }
  reportFound(token, payload) { return this.request(`/public/pets/${token}/found`, { method: 'POST', body: JSON.stringify(payload) }); }

  // Health & Reminders
  getVaccinationTypes(species) { return this.request(`/vaccination-types?species=${species || ''}`); }
  getVaccinations(petId) { return this.request(`/pets/${petId}/vaccinations`); }
  createVaccination(petId, payload) { return this.request(`/pets/${petId}/vaccinations`, { method: 'POST', body: JSON.stringify(payload) }); }
  getMedicalEvents(petId) { return this.request(`/pets/${petId}/medical-events`); }
  createMedicalEvent(petId, payload) { return this.request(`/pets/${petId}/medical-events`, { method: 'POST', body: JSON.stringify(payload) }); }
  getMedications(petId) { return this.request(`/pets/${petId}/medications`); }
  createMedication(petId, payload) { return this.request(`/pets/${petId}/medications`, { method: 'POST', body: JSON.stringify(payload) }); }
  getTimeline(petId) { return this.request(`/pets/${petId}/timeline`); }

  // Documents
  getDocuments(petId) { return this.request(`/pets/${petId}/documents`); }
  saveDocument(petId, payload) { return this.request(`/pets/${petId}/documents`, { method: 'POST', body: JSON.stringify(payload) }); }
  getPresignedUrl(payload) { return this.request('/storage/presigned-url', { method: 'POST', body: JSON.stringify(payload) }); }

  // Reminders
  getReminders() { return this.request('/reminders'); }
  createReminder(payload) { return this.request('/reminders', { method: 'POST', body: JSON.stringify(payload) }); }
  completeReminder(id) { return this.request(`/reminders/${id}/complete`, { method: 'POST' }); }

  // Vets
  getVets(query, city, emergencyOnly) {
    const params = new URLSearchParams();
    if (query) params.set('query', query);
    if (city) params.set('city', city);
    if (emergencyOnly) params.set('emergencyOnly', 'true');
    return this.request(`/vets?${params.toString()}`);
  }
  getNearbyVets(lat, lng, radiusKm = 25) {
    return this.request(`/vets/nearby?latitude=${lat}&longitude=${lng}&radiusKm=${radiusKm}`);
  }
}

window.api = new ApiClient();
