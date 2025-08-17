import { useState } from "react";

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

/**
 *
 */
export const useChat = ({ initialMessages }: { initialMessages?: ChatMessage[] } = {}) => {
  const [messages, setMessages] = useState<ChatMessage[]>(initialMessages ?? []);

  // can we do something fancier with react suspense?
  const [loading, setLoading] = useState(false);

  const sendMessage = async (newMessage: UserMessage) => {
    const newMessages = [
      ...messages,
      newMessage
    ];

    setMessages(newMessages);
    setLoading(true);

    return fetch("/api/chat", {
      method: "POST",
      body: JSON.stringify({ messages: newMessages }),
      headers: {
        "Content-Type": "application/json"
      }
    })
      .then(res => res.json())
      .then(json => setMessages([
        ...newMessages,
        systemMessage(json.response)
      ]))
      .finally(() => setLoading(false))
  }

  return { messages, sendMessage, loading }
}

