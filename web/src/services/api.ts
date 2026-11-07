import type { City, Overview } from '../types/weather';
const get = async <T,>(url: string): Promise<T> => { const response = await fetch(url); if (!response.ok) { const body = await response.json().catch(() => ({})) as { error?: string }; throw new Error(body.error ?? `Request failed (${response.status})`); } return response.json() as Promise<T>; };
export const searchCities = (name: string) => get<City[]>(`/api/cities/search?name=${encodeURIComponent(name)}`);
export const getOverview = (city: City) => get<Overview>(`/api/weather/overview?latitude=${city.latitude}&longitude=${city.longitude}`);
