import { SLICE_NAMES } from "../../constants/enums"

type RequestOptions = {
  headers?: Record<string, string>
  params?: Record<string, string | number | boolean>
  body?: unknown
}

class FetchClient {
  private baseURL: string

  constructor(baseURL = '') {
    this.baseURL = baseURL
  }

  private buildURL(
    url: string,
    params?: Record<string, string | number | boolean>
  ) {
    const fullUrl = new URL(url, this.baseURL)

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        fullUrl.searchParams.append(key, String(value))
      })
    }

    return fullUrl.toString()
  }

  private async request<T>(
    method: string,
    url: string,
    options: RequestOptions = {}
  ): Promise<T> {
    const token = localStorage.getItem(SLICE_NAMES.USER)

    const response = await fetch(
      this.buildURL(url, options.params),
      {
        method,
        headers: {
          'Content-Type': 'application/json',
          ...(token && {
            Authorization: `Bearer ${token}`,
          }),
          ...options.headers,
        },
        body:
          options.body !== undefined
            ? JSON.stringify(options.body)
            : undefined,
        credentials: 'include',
      }
    )

    let data

    const contentType = response.headers.get('content-type')

    if (contentType?.includes('application/json')) {
      data = await response.json()
    } else {
      data = await response.text()
    }

    if (!response.ok) {
      throw {
        status: response.status,
        data,
      }
    }

    return data
  }

  get<T>(url: string, options?: RequestOptions) {
    return this.request<T>('GET', url, options)
  }

  post<T>(url: string, body?: unknown, options?: RequestOptions) {
    return this.request<T>('POST', url, {
      ...options,
      body,
    })
  }

  put<T>(url: string, body?: unknown, options?: RequestOptions) {
    return this.request<T>('PUT', url, {
      ...options,
      body,
    })
  }

  patch<T>(url: string, body?: unknown, options?: RequestOptions) {
    return this.request<T>('PATCH', url, {
      ...options,
      body,
    })
  }

  delete<T>(url: string, options?: RequestOptions) {
    return this.request<T>('DELETE', url, options)
  }
}

export const api = new FetchClient()