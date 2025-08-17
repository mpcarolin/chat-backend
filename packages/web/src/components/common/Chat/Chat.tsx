import { useState, useRef, useEffect } from "react";
import { systemMessage, useChat, userMessage } from "../../../hooks/useChat";
import "./Chat.css";

const STOP_COMMANDS = [
    "quit",
    "q",
    "end",
    "stop",
    "goodbye",
    "bye",
    "finish",
    "finished",
    "i'm finished",
    "i'm done",
    "done",
    "all done",
    "complete",
];

export const Chat = () => {
    const { messages, sendMessage, loading } = useChat({
        initialMessages: [
            systemMessage("How can I help you?"),
        ]
    });
    const [pendingMessage, setPendingMessage] = useState("")

    const penultimateMessageRef = useRef<HTMLDivElement>(null);
    useEffect(() => {
        penultimateMessageRef.current?.scrollIntoView({ behavior: "smooth" })
    }, [messages]);

    return (
        <div className="chat-root">
            <div className="chat-messages">
                {messages.map(
                    (msg, idx) => <Message
                        key={msg.uuid}
                        position={msg.role === "system" ? "left" : "right"}
                        content={msg.content}
                        ref={idx === messages.length - 2 ? penultimateMessageRef : undefined}
                    />
                )}
                {loading && <Message position="left" loading />}
            </div>
            <form className="chat-form" onSubmit={e => {
                e.preventDefault();
                if (!pendingMessage) return;
                sendMessage(userMessage(pendingMessage))
                setPendingMessage("")
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
})

const SendMessageIcon = () => {
    return (
        <svg width="24px" height="24px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path fillRule="evenodd" clipRule="evenodd" d="M3.3938 2.20468C3.70395 1.96828 4.12324 1.93374 4.4679 2.1162L21.4679 11.1162C21.7953 11.2895 22 11.6296 22 12C22 12.3704 21.7953 12.7105 21.4679 12.8838L4.4679 21.8838C4.12324 22.0662 3.70395 22.0317 3.3938 21.7953C3.08365 21.5589 2.93922 21.1637 3.02382 20.7831L4.97561 12L3.02382 3.21692C2.93922 2.83623 3.08365 2.44109 3.3938 2.20468ZM6.80218 13L5.44596 19.103L16.9739 13H6.80218ZM16.9739 11H6.80218L5.44596 4.89699L16.9739 11Z" fill="#000000" />
        </svg>
    )
}


