const API_BASE_URL = (process.env.REACT_APP_API_URL || "http://localhost:8080/api/v1").replace(/\/$/, "");

export async function apiRequest(path, options = {}) {
  const isFormData = typeof FormData !== "undefined" && options.body instanceof FormData;
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    credentials: "include",
    headers: isFormData ? { ...(options.headers || {}) } : { "Content-Type": "application/json", ...(options.headers || {}) },
  });
  const body = await response.json().catch(() => ({}));
  if (!response.ok) throw new Error(body.message || "Terjadi kesalahan pada server.");
  return body;
}

export { API_BASE_URL };
