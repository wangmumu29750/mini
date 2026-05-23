export interface ApiResponse<T> {
  code: string
  message: string
  data: T
  traceId?: string
}

export interface PageResult<T> {
  items: T[]
  page: number
  pageSize: number
  total: number
}

export interface ApiErrorPayload {
  code: string
  message: string
  traceId?: string
  status?: number
}

