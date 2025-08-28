/**
 * Gets the chat provider from environment variables
 */
export const getChatProvider = (): string => {
  return import.meta.env.VITE_CHAT_PROVIDER || 'mock';
};

const STREAM_PROVIDERS = [
  "openai",
  "ollama"
];

/**
 * Determines if streaming should be enabled based on the chat provider
 * Streaming is supported by OpenAI and Ollama providers
 */
export const shouldEnableStreaming = (): boolean => {
  const provider = getChatProvider();
  return STREAM_PROVIDERS.includes(provider);
};
