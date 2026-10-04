/**
 * ui.js
 * Toast notifications, modal windows, QR rendering, and UI utility helpers.
 */

class UIManager {
  constructor() {
    this.createToastContainer();
  }

  createToastContainer() {
    if (!document.querySelector('.toast-container')) {
      const container = document.createElement('div');
      container.className = 'toast-container';
      document.body.appendChild(container);
    }
  }

  showToast(message, type = 'info', duration = 3500) {
    const container = document.querySelector('.toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;

    let icon = 'ℹ️';
    if (type === 'success') icon = '✓';
    if (type === 'error') icon = '✕';

    toast.innerHTML = `<span>${icon}</span> <span>${message}</span>`;
    container.appendChild(toast);

    setTimeout(() => {
      toast.style.opacity = '0';
      toast.style.transform = 'translateY(10px)';
      toast.style.transition = 'all 0.3s ease';
      setTimeout(() => toast.remove(), 300);
    }, duration);
  }

  openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
      modal.classList.add('active');
    }
  }

  closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
      modal.classList.remove('active');
    }
  }

  formatDate(dateStr) {
    if (!dateStr) return 'N/A';
    try {
      const d = new Date(dateStr);
      return d.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
    } catch {
      return dateStr;
    }
  }

  formatDateTime(dateStr) {
    if (!dateStr) return 'N/A';
    try {
      const d = new Date(dateStr);
      return d.toLocaleString('en-US', {
        month: 'short',
        day: 'numeric',
        hour: 'numeric',
        minute: '2-digit',
      });
    } catch {
      return dateStr;
    }
  }

  getStatusBadge(status) {
    switch (status) {
      case 'ACTIVE':
        return '<span class="badge badge-active">● Active</span>';
      case 'LOST':
        return '<span class="badge badge-lost">🚨 LOST</span>';
      case 'FOUND':
        return '<span class="badge badge-found">🔍 Found</span>';
      case 'REUNITED':
        return '<span class="badge badge-reunited">🎉 Reunited</span>';
      default:
        return `<span class="badge">${status}</span>`;
    }
  }

  getSpeciesIcon(species) {
    switch (species?.toUpperCase()) {
      case 'DOG': return '🐶';
      case 'CAT': return '🐱';
      case 'BIRD': return '🦜';
      case 'RABBIT': return '🐰';
      case 'HORSE': return '🐴';
      default: return '🐾';
    }
  }
}

window.ui = new UIManager();
