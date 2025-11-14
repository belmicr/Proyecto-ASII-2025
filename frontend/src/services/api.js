// Configuración de URLs de microservicios
const API_URLS = {
  users: "http://localhost:8080",        // users-api
  search: "http://localhost:8082",       // search-api
  hotels: "http://localhost:8081",       // hotels-api
  reservations: "http://localhost:8086", // reservations-api
};

// Helper para hacer requests con token (Bearer)
const fetchWithAuth = async (url, options = {}) => {
  const token = localStorage.getItem("token");

  const config = {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token && { Authorization: `Bearer ${token}` }),
      ...(options.headers || {}),
    },
  };

  const response = await fetch(url, config);

  if (response.status === 401) {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    window.location.href = "/login";
  }

  return response;
};

/* =======================
 *  AUTH (users-api)
 * ======================= */
export const authService = {
  // POST http://localhost:8080/users/login
  login: async (email, password) => {
    const response = await fetch(`${API_URLS.users}/users/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    return response.json();
  },

  // POST http://localhost:8080/users
  register: async (userData) => {
    const response = await fetch(`${API_URLS.users}/users`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(userData),
    });
    return response.json();
  },

  logout: () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
  },
};

/* =======================
 *  SEARCH (search-api)
 * ======================= */
export const searchService = {
  // GET http://localhost:8082/search?q=...&page=...
  search: async (query = "", page = 1) => {
    const response = await fetchWithAuth(
      `${API_URLS.search}/search?q=${encodeURIComponent(
        query
      )}&page=${page}`
    );
    return response.json();
  },
};

/* =======================
 *  HOTELS / ROOMS (hotels-api)
 * ======================= */
export const hotelService = {
  // GET http://localhost:8081/hotels
  getAll: async () => {
    const response = await fetchWithAuth(`${API_URLS.hotels}/hotels`);
    return response.json();
  },

  // GET http://localhost:8081/hotels/:id
  getById: async (id) => {
    const response = await fetchWithAuth(`${API_URLS.hotels}/hotels/${id}`);
    return response.json();
  },

  // Para admin:
  // POST http://localhost:8081/hotels
  create: async (hotelData) => {
    const response = await fetchWithAuth(`${API_URLS.hotels}/hotels`, {
      method: "POST",
      body: JSON.stringify(hotelData),
    });
    return response.json();
  },

  // PUT http://localhost:8081/hotels/:id
  update: async (id, hotelData) => {
    const response = await fetchWithAuth(`${API_URLS.hotels}/hotels/${id}`, {
      method: "PUT",
      body: JSON.stringify(hotelData),
    });
    return response.json();
  },

  // DELETE http://localhost:8081/hotels/:id
  delete: async (id) => {
    const response = await fetchWithAuth(`${API_URLS.hotels}/hotels/${id}`, {
      method: "DELETE",
    });
    return response.json();
  },
};

/* =======================
 *  RESERVATIONS (reservations-api)
 * ======================= */
export const reservationService = {
  // POST http://localhost:8086/reservations
  create: async (reservationData) => {
    const response = await fetchWithAuth(
      `${API_URLS.reservations}/reservations`,
      {
        method: "POST",
        body: JSON.stringify(reservationData),
      }
    );
    return response.json();
  },

  // GET http://localhost:8086/reservations/my
  getMyReservations: async () => {
    const response = await fetchWithAuth(
      `${API_URLS.reservations}/reservations/my`
    );
    return response.json();
  },

  // DELETE http://localhost:8086/reservations/:id
  delete: async (id) => {
    const response = await fetchWithAuth(
      `${API_URLS.reservations}/reservations/${id}`,
      {
        method: "DELETE",
      }
    );
    return response.json();
  },

  // Para admin:
  // GET http://localhost:8086/reservations
  getAll: async () => {
    const response = await fetchWithAuth(
      `${API_URLS.reservations}/reservations`
    );
    return response.json();
  },

  // PUT http://localhost:8086/reservations/:id
  update: async (id, data) => {
    const response = await fetchWithAuth(
      `${API_URLS.reservations}/reservations/${id}`,
      {
        method: "PUT",
        body: JSON.stringify(data),
      }
    );
    return response.json();
  },
};

export default API_URLS;
