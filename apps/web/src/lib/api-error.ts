type ApiErrorPayload = {
  error?: string;
};

type ApiLikeError = {
  response?: {
    data?: ApiErrorPayload;
  };
};

export function extractApiErrorMessage(error: unknown, fallback: string): string {
  if (typeof error === "object" && error !== null && "response" in error) {
    const response = (error as ApiLikeError).response;
    return response?.data?.error || fallback;
  }
  return fallback;
}
