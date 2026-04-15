import { useGlobalState } from '~/store'

export const useIsAdmin = () => {
  const global = useGlobalState()
  return computed(() => global.value.userinfo.id === 1)
}
