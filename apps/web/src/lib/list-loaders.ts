import { extractApiErrorMessage } from "@/lib/api-error";

type CollectionResponse<T> = {
  data?: T[];
  pagination?: unknown;
};

type FetchCollectionWithStateOptions<T> = {
  fetcher: () => Promise<CollectionResponse<T>>;
  setItems: (items: T[]) => void;
  fallbackError: string;
  setLoading?: (value: boolean) => void;
  setError?: (message: string) => void;
  setPagination?: (pagination: unknown) => void;
  onErrorMessage?: (message: string) => void;
};

export async function fetchCollectionWithState<T>({
  fetcher,
  setItems,
  fallbackError,
  setLoading,
  setError,
  setPagination,
  onErrorMessage,
}: FetchCollectionWithStateOptions<T>): Promise<void> {
  try {
    setLoading?.(true);
    setError?.("");

    const response = await fetcher();
    setItems(response.data || []);

    if (setPagination && response.pagination) {
      setPagination(response.pagination);
    }
  } catch (error) {
    const message = extractApiErrorMessage(error, fallbackError);
    setError?.(message);
    onErrorMessage?.(message);
  } finally {
    setLoading?.(false);
  }
}

type LoadListWithDefaultSelectionOptions<T> = {
  fetchList: () => Promise<T[] | undefined | null>;
  setList: (items: T[]) => void;
  setSelectedId?: (id: string) => void;
  getId?: (item: T) => string;
  onEmpty?: () => void;
  onError?: () => void;
  onFinally?: () => void;
};

type LoadListEffectOptions<T> = LoadListWithDefaultSelectionOptions<T> & {
  enabled?: boolean;
  beforeLoad?: () => void;
};

export async function loadListWithDefaultSelection<T>({
  fetchList,
  setList,
  setSelectedId,
  getId,
  onEmpty,
  onError,
  onFinally,
}: LoadListWithDefaultSelectionOptions<T>): Promise<void> {
  try {
    const list = (await fetchList()) || [];
    setList(list);

    if (list.length > 0 && setSelectedId && getId) {
      setSelectedId(getId(list[0]));
    } else {
      onEmpty?.();
    }
  } catch {
    onError?.();
  } finally {
    onFinally?.();
  }
}

export async function loadListEffect<T>({
  enabled = true,
  beforeLoad,
  ...options
}: LoadListEffectOptions<T>): Promise<void> {
  if (!enabled) {
    return;
  }

  beforeLoad?.();
  await loadListWithDefaultSelection(options);
}
