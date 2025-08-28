import { Chat } from "../../common/Chat/Chat"
import { systemMessage, useChat, userMessage } from "../../../hooks/useChat"
import "./ChatDemoContainer.css"

/**
 * Container that manages chat state and provides message handling for the chat demo
 */
export const ChatDemoContainer = () => {
    const { messages, sendMessage, loading } = useChat({
        initialMessages: [
            systemMessage("How can I help you?"),
        ],
    });

    const handleMessageSubmit = (message: string) => {
        sendMessage(
            userMessage(message),
            // stream can only be enabled for providers that support it, e.g. OpenAI
            { stream: false }
        );
    };

    return (
        <div className="chat-demo-container">
            <Chat
                messages={messages}
                onMessageSubmit={handleMessageSubmit}
                loading={loading}
            />
        </div>
    )
}
