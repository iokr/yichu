export function formatDate(value?: string | null): string {
  if (!value) return '未填写'
  return value.slice(0, 10)
}

export function formatMonth(value?: string | null): string {
  if (!value) return '—'
  const [year, month] = value.slice(0, 10).split('-')
  return `${year}/${month}`
}

export function formatDateTime(value?: string | null): string {
  if (!value) return '—'
  const date = new Date(value.replace(' ', 'T'))
  if (Number.isNaN(date.getTime())) return value
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}/${pad(date.getMonth() + 1)}/${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

export function today(): string {
  const date = new Date()
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
}

export function wearText(count: number): string {
  return count > 0 ? `穿过 ${count} 次` : '还没穿过'
}

export function relativeDays(value?: string | null): string {
  if (!value) return '从未穿着'
  const target = new Date(value.replace(' ', 'T'))
  if (Number.isNaN(target.getTime())) return value
  const diff = Math.floor((Date.now() - target.getTime()) / 86400000)
  if (diff <= 0) return '今天穿过'
  if (diff === 1) return '昨天穿过'
  if (diff < 30) return `${diff} 天前穿过`
  if (diff < 365) return `${Math.floor(diff / 30)} 个月前穿过`
  return `${Math.floor(diff / 365)} 年前穿过`
}

export function percentText(value: number): string {
  return `${value.toFixed(1)}%`
}
