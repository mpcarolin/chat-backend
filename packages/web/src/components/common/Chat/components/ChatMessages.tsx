import { useRef, useEffect } from "react";
import { Message } from "./Message";
import type { ChatMessage } from "../../../../hooks/useChat";
import "./ChatMessages.css";

interface ChatMessagesProps {
    messages: ChatMessage[];
    loading: boolean;
}

/**
 * Scrollable container that displays chat messages with auto-scroll to latest
 */
export const ChatMessages = ({ messages, loading }: ChatMessagesProps) => {
    const messageRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        messageRef.current?.scrollIntoView({ behavior: "smooth" })
    }, [messages]);

    return (
        <div className="chat-messages">
            {messages.map(
                (msg, idx) => <Message
                    key={msg.uuid}
                    position={msg.role === "system" ? "left" : "right"}
                    content={msg.content}
                    ref={idx === messages.length - 1 ? messageRef : undefined}
                />
            )}
            {loading && <Message position="left" loading />}
        </div>
    );
}
