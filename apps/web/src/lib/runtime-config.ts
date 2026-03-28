const DEFAULT_API_BASE_URL = "http://localhost:8080/api/v1";

function parseBoolEnv(value: string | undefined, defaultValue = false): boolean {
  if (!value) return defaultValue;
  const normalized = value.trim().toLowerCase();
  if (["1", "true", "yes", "y", "on"].includes(normalized)) return true;
  if (["0", "false", "no", "n", "off"].includes(normalized)) return false;
  return defaultValue;
}

export function getApiBaseUrl(): string {
  return process.env.NEXT_PUBLIC_API_URL || DEFAULT_API_BASE_URL;
}

export function getWsBaseUrl(): string {
  return getApiBaseUrl().replace(/^http/, "ws");
}
