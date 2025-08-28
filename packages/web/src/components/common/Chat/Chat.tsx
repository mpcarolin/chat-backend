import { useState } from "react";
import type { ChatMessage } from "../../../hooks/useChat"
import { ChatMessages, SendMessageIcon } from "./components";
import "./Chat.css";

interface ChatProps {
    messages: ChatMessage[];
    onMessageSubmit: (message: string) => void;
    loading: boolean;
}

/**
 * Main chat interface component that displays messages and handles user input
 */
export const Chat = ({ messages, onMessageSubmit, loading }: ChatProps) => {
    const [pendingMessage, setPendingMessage] = useState("")

    return (
        <div className="chat-root">
            <ChatMessages messages={messages} loading={loading} />
            <form className="chat-form" onSubmit={e => {
                e.preventDefault();
                if (!pendingMessage) return;
                onMessageSubmit(pendingMessage);
                setPendingMessage("");
            }}>
                <input
                    name="userinput"
                    type="text"
                    className="chat-input"
                    placeholder="Ask a question..."
                    value={pendingMessage}
                    onChange={e => setPendingMessage(e.currentTarget.value)}
                />
                <button type="submit" disabled={!pendingMessage || loading} className="submit-button">
                    <SendMessageIcon />
                </button>
            </form>
        </div>
    )
}


