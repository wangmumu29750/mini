import { request } from '@/api/http'
import type { SystemSetting } from '@/types/domain'

export function fetchSystemSettings() {
  return request<SystemSetting[]>({
    method: 'GET',
    url: '/admin/settings',
  })
}

export function updateSystemSettings(settings: Pick<SystemSetting, 'key' | 'value'>[]) {
  return request<SystemSetting[]>({
    method: 'PUT',
    url: '/admin/settings',
    data: { settings },
  })
}
