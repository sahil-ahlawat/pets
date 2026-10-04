/**
 * app.js
 * Common application initialization, navbar injection, and active pet switcher.
 */

document.addEventListener('DOMContentLoaded', () => {
  renderNavbar();
  setupGlobalPetSwitcher();
});

function renderNavbar() {
  const navContainer = document.getElementById('navbar-container');
  if (!navContainer) return;

  const isAuth = window.auth.isAuthenticated();
  const user = window.auth.getUser();

  const currentPath = window.location.pathname;

  navContainer.innerHTML = `
    <nav class="navbar">
      <div class="container nav-wrapper">
        <a href="/pages/index.html" class="brand-logo">
          <div class="logo-icon">🐾</div>
          <span>PetIdentity</span>
        </a>

        <ul class="nav-links">
          ${isAuth ? `
            <li><a href="/pages/dashboard.html" class="nav-link ${currentPath.includes('dashboard') ? 'active' : ''}">Dashboard</a></li>
            <li><a href="/pages/pet.html" class="nav-link ${currentPath.includes('pet.html') ? 'active' : ''}">My Pets</a></li>
            <li><a href="/pages/vaccinations.html" class="nav-link ${currentPath.includes('vaccinations') ? 'active' : ''}">Vaccinations</a></li>
            <li><a href="/pages/reminders.html" class="nav-link ${currentPath.includes('reminders') ? 'active' : ''}">Reminders</a></li>
            <li><a href="/pages/lost-found.html" class="nav-link ${currentPath.includes('lost-found') ? 'active' : ''}">Lost & Found</a></li>
            <li><a href="/pages/vets.html" class="nav-link ${currentPath.includes('vets') ? 'active' : ''}">Find Vets</a></li>
          ` : `
            <li><a href="/pages/features.html" class="nav-link ${currentPath.includes('features') ? 'active' : ''}">Features</a></li>
            <li><a href="/pages/pricing.html" class="nav-link ${currentPath.includes('pricing') ? 'active' : ''}">Smart Tags</a></li>
            <li><a href="/pages/lost-found.html" class="nav-link ${currentPath.includes('lost-found') ? 'active' : ''}">Lost Pets</a></li>
            <li><a href="/pages/vets.html" class="nav-link ${currentPath.includes('vets') ? 'active' : ''}">Vets</a></li>
            <li><a href="/pages/blog.html" class="nav-link ${currentPath.includes('blog') ? 'active' : ''}">Articles</a></li>
          `}
        </ul>

        <div class="nav-actions">
          ${isAuth ? `
            <a href="/pages/pet-new.html" class="btn btn-primary btn-sm">+ Add Pet</a>
            <a href="/pages/profile.html" class="btn btn-secondary btn-sm">👤 ${user ? user.firstName : 'Profile'}</a>
            <button onclick="window.auth.logout()" class="btn btn-outline btn-sm">Logout</button>
          ` : `
            <a href="/pages/login.html" class="btn btn-secondary btn-sm">Sign In</a>
            <a href="/pages/register.html" class="btn btn-primary btn-sm">Get Started</a>
          `}
        </div>
      </div>
    </nav>
  `;
}

async function setupGlobalPetSwitcher() {
  const switcher = document.getElementById('global-pet-switcher');
  if (!switcher || !window.auth.isAuthenticated()) return;

  try {
    const res = await window.api.getPets();
    const pets = res.pets || [];
    if (pets.length === 0) return;

    let activePetId = localStorage.getItem('active_pet_id');
    if (!activePetId || !pets.some(p => p.id === activePetId)) {
      activePetId = pets[0].id;
      localStorage.setItem('active_pet_id', activePetId);
    }

    switcher.innerHTML = pets.map(p => `
      <option value="${p.id}" ${p.id === activePetId ? 'selected' : ''}>
        ${window.ui.getSpeciesIcon(p.species)} ${p.name}
      </option>
    `).join('');

    switcher.addEventListener('change', (e) => {
      localStorage.setItem('active_pet_id', e.target.value);
      window.location.reload();
    });
  } catch (err) {
    console.error('Failed to load pets into switcher:', err);
  }
}
