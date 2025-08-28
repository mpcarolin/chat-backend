/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_CHAT_PROVIDER?: string
  readonly VITE_API_URL?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
