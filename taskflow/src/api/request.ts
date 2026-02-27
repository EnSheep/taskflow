const BASE_URL = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080'

export interface ApiError {
  error: string
}

export interface RequestOptions {
  method?: string
  body?: object
  headers?: HeadersInit
  /** GET 请求的查询参数，会拼到 path 上 */
  params?: Record<string, string>
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, headers: customHeaders, params } = options
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(customHeaders as Record<string, string>),
  }
  let url = `${BASE_URL}${path}`
  if (params && Object.keys(params).length > 0) {
    const search = new URLSearchParams(params).toString()
    url += (path.includes('?') ? '&' : '?') + search
  }
  const res = await fetch(url, {
    method,
    headers,
    body: body != null ? JSON.stringify(body) : undefined,
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    const msg = (data as ApiError).error ?? res.statusText ?? '请求失败'
    throw new Error(msg)
  }
  return data as T
}

export { BASE_URL, request }
