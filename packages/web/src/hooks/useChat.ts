import { useState } from "react";

type Message = {
  role: "user" | "system"
  content: string
}

type UserMessage = Message & {
  role: "user"
}

/**
 *
 */
export const useChat = () => {
  const [messages, setMessages] = useState<Message[]>([]);

  const sendMessage = async (newMessage: UserMessage) => {
    const newMessages = [
      ...messages,
      newMessage
    ];
    setMessages(newMessages);

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
        { role: "system", content: json.response }
      ]));
  }

  return [messages, sendMessage]
}

