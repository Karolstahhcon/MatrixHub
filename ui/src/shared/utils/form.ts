/**
 * Extract the first error message from a TanStack Form field.
 * Standard Schema validators (e.g. Zod) return issue objects with a `.message`
 * property, while custom validators may return plain strings.
 */
export function fieldError(field: { state: { meta: { errors: unknown[]
  isTouched: boolean } } }): string | undefined {
  if (!field.state.meta.isTouched) {
    return undefined
  }

  const first = field.state.meta.errors[0]

