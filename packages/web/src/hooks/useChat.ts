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

/**
 *
 */
export const useChat = ({ initialMessages, stream }: { initialMessages?: ChatMessage[], stream?: boolean } = {}) => {
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

    const fetchFn = stream ? sendMessageStreaming : sendMessageNonStreaming

    return fetchFn(newMessages, setMessages).finally(() => setLoading(false))
  }

  return { messages, sendMessage, loading }
}

const sendMessageBase = (messages: ChatMessage[], streaming?: boolean) => {
  return fetch("/api/chat", {
    method: "POST",
    body: JSON.stringify({ messages, streaming }),
    headers: {
      "Content-Type": "application/json"
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

  const reader = res.body.getReader();
  const decoder = new TextDecoder("utf-8");
  let buffer = "";

  setMessages([
    ...messages,
    systemMessage("") // <-- start a new message that we will be adding to with each chunk
  ])

  while (true) {
    const { value, done } = await reader.read();
    console.log({ value, done })
    if (done) break;

    // convert bytes to string
    buffer += decoder.decode(value, { stream: true });

    // split by newline
    const lines = buffer.split("\n");
    console.log({ lines, buffer })
    buffer = lines.pop() ?? ""; // incomplete line stays in buffer

    for (const line of lines) {
      if (!line.trim()) continue;
      try {
        const chunk = JSON.parse(line);
        console.log({ chunk })

        // incrementally build this last message
        setMessages((prev) => {
          const [last, ...messages] = prev.reverse();
          return [
            ...messages,
            appendToMessage(last, chunk)
          ]
        });
      } catch (err) {
        console.error("Failed to pase chunk", err, line);
      }
    }
  }
}


