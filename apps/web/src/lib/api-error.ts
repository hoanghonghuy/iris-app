type ApiErrorPayload = {
  error?: string;
};

type ApiLikeError = {
  response?: {
    data?: ApiErrorPayload;
  };
};

export function extractApiErrorRawMessage(error: unknown): string | undefined {
  if (typeof error === "object" && error !== null && "response" in error) {
    const response = (error as ApiLikeError).response;
    return response?.data?.error;
  }

  return undefined;
}

export function extractApiErrorMessage(error: unknown, fallback: string): string {
  return extractApiErrorRawMessage(error) || fallback;
}
