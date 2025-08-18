import { useState } from "react";
import { fetchEventSource } from '@microsoft/fetch-event-source';


export const userMessage = (content: string): UserMessage => ({
  content,
  role: "user",
  uuid: crypto.randomUUID()
});

export const systemMessage = (content: string): SystemMessage => ({
  content,
  role: "system",
  uuid: crypto.randomUUID()
});

export const appendToMessage = (msg: ChatMessage, content: string) => {
  return {
    ...msg,
    content: msg.content + content
  }
}

export type ChatMessage = {
  role: "user" | "system"
  content: string,
  uuid: string
}

export type UserMessage = ChatMessage & {
  role: "user"
}

export type SystemMessage = ChatMessage & {
  role: "system"
}

export type StreamResponse = {
  done: boolean;
  response: ChatMessage
}

/**
 *
 */
export const useChat = ({ initialMessages }: { initialMessages?: ChatMessage[] } = {}) => {
  const [messages, setMessages] = useState<ChatMessage[]>(initialMessages ?? []);

  // can we do something fancier with react suspense?
  const [loading, setLoading] = useState(false);

  const sendMessage = async (newMessage: UserMessage, options: { stream?: boolean } = {}) => {
    const newMessages = [
      ...messages,
      newMessage
    ];

    setMessages(newMessages);
    setLoading(true);

    const fetchFn = options.stream ? sendMessageStreaming : sendMessageNonStreaming

    return fetchFn(newMessages, setMessages).finally(() => setLoading(false))
  }

  return { messages, sendMessage, loading }
}

const sendMessageBase = (messages: ChatMessage[], streaming?: boolean) => {
  return fetch("/api/chat", {
    method: "POST",
    body: JSON.stringify({ messages, streaming }),
    headers: {
      "Content-Type": "application/json",
    }
  })
}

const sendMessageNonStreaming = async (messages: ChatMessage[], setMessages: React.Dispatch<React.SetStateAction<ChatMessage[]>>) => {
  return sendMessageBase(messages, true)
    .then(res => res.json())
    .then(json => setMessages([
      ...messages,
      systemMessage(json.response)
    ]))
}

const sendMessageStreaming = async (messages: ChatMessage[], setMessages: React.Dispatch<React.SetStateAction<ChatMessage[]>>) => {
  const res = await sendMessageBase(messages, true)
  if (!res.body) throw new Error("Uh oh");

  setMessages([
    ...messages,
    systemMessage("") // <-- start a new message that we will be adding to with each chunk
  ]);

  const decoder = new TextDecoder();

  return readChunks(res.body, (chunk) => {
    const raw = decoder.decode(chunk);
    console.log({ raw, chunk })
    const streamResponse = JSON.parse(raw) as StreamResponse;
    setMessages(prev => {
      const last = prev[prev.length - 1];
      return [
        ...prev.slice(0, prev.length - 1),
        appendToMessage(last, streamResponse.response.content)
      ];
    });
  })

}

async function readChunks<T>(stream: ReadableStream<T>, onChunkReceived: (chunk: T) => void) {
  const reader = stream.getReader();
  const chunks = [];

  let done, value;
  while (!done) {
    ({ value, done } = await reader.read());
    if (done) {
      break;
    }
    if (value !== undefined) {
      onChunkReceived(value);
    }

    chunks.push(value);
  }

  return chunks;
}


