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
  response: string
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

    const send = options.stream
      ? sendMessageStreaming
      : sendMessageNonStreaming;

    return send(newMessages, setMessages, setLoading).finally(() => setLoading(false))
  }

  return { messages, sendMessage, loading }
}

const sendMessageNonStreaming = async (
  messages: ChatMessage[],
  setMessages: React.Dispatch<React.SetStateAction<ChatMessage[]>>
) => {
  return fetch("/api/chat", {
    method: "POST",
    body: JSON.stringify({ messages }),
    headers: {
      "Content-Type": "application/json",
    }
  })
    .then(res => res.json())
    .then(json => setMessages([
      ...messages,
      systemMessage(json.response)
    ]))
}

const sendMessageStreaming = async (
  messages: ChatMessage[],
  setMessages: React.Dispatch<React.SetStateAction<ChatMessage[]>>,
  setLoading: React.Dispatch<React.SetStateAction<boolean>>
) => {
  setMessages([
    ...messages,
    systemMessage("") // <-- start a new message that we will be adding to with each chunk
  ]);
  return fetchEventSource("/api/chat", {
    method: "POST",
    body: JSON.stringify({ messages, streaming: true }),
    headers: {
      "Content-Type": "application/json",
    },
    onmessage: (message) => {
      console.log("onmessage", { message })
      let value: StreamResponse;
      if (!message.data) {
        value = { response: " ", done: false };
      }
      try {
        value = JSON.parse(message.data)
      } catch (err) {
        console.log("Could not parse part of message", err)
        return;
      }

      setLoading(false); // once we've received a valid message, can disable loading ui

      if (value.done) {
        return;
      }

      setMessages(prev => {
        const last = prev[prev.length - 1];
        return [
          ...prev.slice(0, prev.length - 1),
          appendToMessage(last, value.response ?? " ")
        ];
      });
    }
  })
}
