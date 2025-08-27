import { useState, useRef, useEffect } from "react";
import { systemMessage, useChat, userMessage } from "../../../hooks/useChat";
import "./Chat.css";

export const Chat = () => {
    const { messages, sendMessage, loading } = useChat({
        initialMessages: [
            systemMessage("How can I help you?"),
        ],
    });
    const [pendingMessage, setPendingMessage] = useState("")

    const messageRef = useRef<HTMLDivElement>(null);
    useEffect(() => {
        messageRef.current?.scrollIntoView({ behavior: "smooth" })
    }, [messages]);

    return (
        <div className="chat-root">
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
            <form className="chat-form" onSubmit={e => {
                e.preventDefault();
                if (!pendingMessage) return;
                sendMessage(
                    userMessage(pendingMessage),
                    { stream: false }
                );
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

const Message = (({ position, content, ref, loading }: {
    position: "left" | "right",
    content?: React.ReactNode,
    ref?: React.RefObject<HTMLDivElement | null>
    loading?: boolean
}) => {
    const messageClass = position === "left" ? "system-message" : "user-message"
    const loadingClass = loading ? "message-shimmer" : ""
    return (
        <div ref={ref} className={`message ${messageClass} ${loadingClass}`}>
            {!loading && <p>{content}</p>}
        </div>
    )
});

const SendMessageIcon = () => {
    return (
        <svg
            width="32"
            height="32"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
        >
            <path
                d="M14.8285 11.9481L16.2427 10.5339L12 6.29122L7.7574 10.5339L9.17161 11.9481L11 10.1196V17.6568H13V10.1196L14.8285 11.9481Z"
                fill="currentColor"
            />
            <path
                fill-rule="evenodd"
                clip-rule="evenodd"
                d="M19.7782 4.22183C15.4824 -0.0739415 8.51759 -0.0739422 4.22183 4.22183C-0.0739415 8.51759 -0.0739422 15.4824 4.22183 19.7782C8.51759 24.0739 15.4824 24.0739 19.7782 19.7782C24.0739 15.4824 24.0739 8.51759 19.7782 4.22183ZM18.364 5.63604C14.8492 2.12132 9.15076 2.12132 5.63604 5.63604C2.12132 9.15076 2.12132 14.8492 5.63604 18.364C9.15076 21.8787 14.8492 21.8787 18.364 18.364C21.8787 14.8492 21.8787 9.15076 18.364 5.63604Z"
                fill="currentColor"
            />
        </svg>
    );
}


